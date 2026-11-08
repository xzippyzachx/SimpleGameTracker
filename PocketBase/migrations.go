package main

import (
	"encoding/json"
	"strconv"

	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(app core.App) error {
		// GameData collection: one record per sync key.
		// sync_uuid is a separate indexed field for O(1) lookups.
		collection := core.NewBaseCollection("game_data")

		collection.Fields.Add(
			&core.TextField{Name: "sync_uuid", Required: true, Max: 36},
			&core.JSONField{Name: "game_data", Required: true, MaxSize: 5_000_000},
			&core.TextField{Name: "created_by", Required: true, Max: 36},
			&core.TextField{Name: "modified_by", Max: 36},
			&core.DateField{Name: "modified_at"},
			&core.TextField{Name: "sync_hash", Max: 64},
		)

		collection.AddIndex("idx_game_data_sync_uuid", true, "sync_uuid", "")
		collection.AddIndex("idx_game_data_created_by", true, "created_by", "")
		collection.AddIndex("idx_game_data_sync_hash", true, "sync_hash", "")

		// The sync_uuid itself is the secret (knowing it grants read access),
		// so the collection is publicly readable but writes go through the
		// custom /api/sync routes only.
		collection.ListRule = types.Pointer("")
		collection.ViewRule = types.Pointer("")
		collection.CreateRule = nil
		collection.UpdateRule = nil
		collection.DeleteRule = nil

		return app.Save(collection)
	}, func(app core.App) error {
		collection, err := app.FindCollectionByNameOrId("game_data")
		if err != nil {
			return err
		}
		return app.Delete(collection)
	}, "1700000000_create_game_data.go")

	m.Register(func(app core.App) error {
		// GameCollection: the parent record for one user's library (one per sync key).
		gameCollection := core.NewBaseCollection("game_collection")

		gameCollection.Fields.Add(
			&core.TextField{Name: "sync_uuid", Required: true, Max: 36},
			&core.TextField{Name: "created_by", Required: true, Max: 36},
			&core.TextField{Name: "modified_by", Max: 36},
			&core.DateField{Name: "modified_at"},
			&core.TextField{Name: "sync_hash", Max: 64},
		)

		gameCollection.AddIndex("idx_game_collection_sync_uuid", true, "sync_uuid", "")
		gameCollection.AddIndex("idx_game_collection_created_by", true, "created_by", "")
		gameCollection.AddIndex("idx_game_collection_sync_hash", true, "sync_hash", "")

		// The sync_uuid itself is the secret (knowing it grants read access),
		// so the collection is publicly readable but writes go through the
		// custom /api/sync routes only.
		gameCollection.ListRule = types.Pointer("")
		gameCollection.ViewRule = types.Pointer("")
		gameCollection.CreateRule = nil
		gameCollection.UpdateRule = nil
		gameCollection.DeleteRule = nil

		if err := app.Save(gameCollection); err != nil {
			return err
		}

		// Game: one row per game in a collection, holding only the IGDB id
		// and its state. Full IGDB data stays in the client's local cache.
		game := core.NewBaseCollection("game")

		game.Fields.Add(
			&core.RelationField{Name: "collection_id", Required: true, CollectionId: gameCollection.Id, CascadeDelete: true},
			&core.NumberField{Name: "game_id", Required: true},
			&core.TextField{Name: "state", Required: true, Max: 20},
		)

		game.AddIndex("idx_game_collection_id_game_id", true, "collection_id, game_id", "")

		// Game rows are only exposed through the custom /api/sync routes,
		// which gate access by sync_uuid/sync_hash.
		game.ListRule = nil
		game.ViewRule = nil
		game.CreateRule = nil
		game.UpdateRule = nil
		game.DeleteRule = nil

		if err := app.Save(game); err != nil {
			return err
		}

		// Data migration: convert each game_data JSON blob into game rows.
		gameDataRecords, err := app.FindAllRecords("game_data")
		if err != nil {
			return err
		}
		for _, record := range gameDataRecords {
			gameCollectionRecord := core.NewRecord(gameCollection)
			gameCollectionRecord.Set("sync_uuid", record.GetString("sync_uuid"))
			gameCollectionRecord.Set("created_by", record.GetString("created_by"))
			gameCollectionRecord.Set("modified_by", record.GetString("modified_by"))
			gameCollectionRecord.Set("modified_at", record.GetDateTime("modified_at"))
			gameCollectionRecord.Set("sync_hash", record.GetString("sync_hash"))
			if err := app.Save(gameCollectionRecord); err != nil {
				return err
			}

			var gameData map[string]any
			if err := json.Unmarshal([]byte(record.GetString("game_data")), &gameData); err != nil {
				return err
			}
			for gameId, value := range gameData {
				entry, ok := value.(map[string]any)
				if !ok {
					continue
				}
				state, _ := entry["state"].(string)
				if state == "" || state == "none" {
					continue
				}
				gameIdNum, err := strconv.ParseFloat(gameId, 64)
				if err != nil {
					continue
				}
				gameRecord := core.NewRecord(game)
				gameRecord.Set("collection_id", gameCollectionRecord.Id)
				gameRecord.Set("game_id", gameIdNum)
				gameRecord.Set("state", state)
				if err := app.Save(gameRecord); err != nil {
					return err
				}
			}
		}

		// Remove the old game_data collection now that data is migrated.
		gameDataCollection, err := app.FindCollectionByNameOrId("game_data")
		if err != nil {
			return err
		}
		return app.Delete(gameDataCollection)
	}, func(app core.App) error {
		game, err := app.FindCollectionByNameOrId("game")
		if err != nil {
			return err
		}
		if err := app.Delete(game); err != nil {
			return err
		}
		gameCollection, err := app.FindCollectionByNameOrId("game_collection")
		if err != nil {
			return err
		}
		return app.Delete(gameCollection)
	}, "1700000001_create_game_collection.go")
}
