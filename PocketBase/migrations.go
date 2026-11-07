package main

import (
	"github.com/pocketbase/pocketbase/core"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/tools/types"
)

func init() {
	m.Register(func(app core.App) error {
		// GameData collection: one record per sync key.
		// sync_uuid doubles as the record id so lookups are O(1) by primary key.
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
}
