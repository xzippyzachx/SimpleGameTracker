package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/pocketbase/pocketbase/core"
)

const (
	// pbDateLayout is PocketBase's stored datetime format (locale-independent):
	// space-separated, millisecond precision, Z suffix.
	pbDateLayout = "2006-01-02 15:04:05.000Z"
)

// newUUID generates a random v4 UUID for new sync keys.
func newUUID() string {
	return uuid.NewString()
}

// syncHash is sha256(sync_uuid) as hex; the frontend computes the same value
// client-side when building share links.
func syncHash(syncUUID string) string {
	sum := sha256.Sum256([]byte(syncUUID))
	return hex.EncodeToString(sum[:])
}

func registerSyncRoutes(se *core.ServeEvent) {
	g := se.Router.Group("/api/sync")

	g.GET("/game-data", getGameData)
	g.POST("/game-data", updateGameData)
	g.GET("/game-data/share/{syncHash}", getGameDataShare)

	// Keep sync_hash in sync with sync_uuid (sha256 of the uuid as hex).
	se.App.OnRecordCreate("game_data").BindFunc(func(e *core.RecordEvent) error {
		if e.Record.GetString("sync_hash") == "" {
			e.Record.Set("sync_hash", syncHash(e.Record.GetString("sync_uuid")))
		}
		return e.Next()
	})
	se.App.OnRecordUpdate("game_data").BindFunc(func(e *core.RecordEvent) error {
		e.Record.Set("sync_hash", syncHash(e.Record.GetString("sync_uuid")))
		return e.Next()
	})
}

// normalizeDate converts a client timestamp into PocketBase's stored datetime
// format (space separator, millisecond precision, UTC) so filter comparisons
// work regardless of the client's flavor of RFC3339 ("...T...Z", "...T...+00:00")
// or an already-normalized PocketBase value echoed back from a previous sync.
func normalizeDate(ts string) string {
	if ts == "" {
		return ""
	}
	if t, err := time.Parse(time.RFC3339Nano, ts); err == nil {
		return t.UTC().Format(pbDateLayout)
	}
	if t, err := time.Parse(pbDateLayout, ts); err == nil {
		return t.UTC().Format(pbDateLayout)
	}
	// Unrecognized: pass through and let the filter report the problem.
	return ts
}

// GET /api/sync/game-data?sync_uuid=...&local_modified_at=...
// Returns the record only when the server copy is newer than the client's
// local timestamp.
func getGameData(e *core.RequestEvent) error {
	syncUUID := e.Request.URL.Query().Get("sync_uuid")
	localModifiedAt := normalizeDate(e.Request.URL.Query().Get("local_modified_at"))

	if syncUUID == "" {
		return e.JSON(http.StatusOK, []any{})
	}

	filter := "sync_uuid = {:sync_uuid}"
	params := map[string]any{"sync_uuid": syncUUID}
	if localModifiedAt != "" {
		filter += " && modified_at > {:local_modified_at}"
		params["local_modified_at"] = localModifiedAt
	}

	records, err := e.App.FindRecordsByFilter("game_data", filter, "", 1, 0, params)
	if err != nil || len(records) == 0 {
		return e.JSON(http.StatusOK, []any{})
	}

	r := records[0]
	return e.JSON(http.StatusOK, []map[string]any{{
		"game_data":   r.GetRaw("game_data"),
		"modified_at": r.GetString("modified_at"),
	}})
}

// POST /api/sync/game-data
// Body: {"device_uuid": "...", "sync_uuid": "...", "game_data": {...}}
// Update by sync_uuid, else update by created_by, else insert a new record
// with a fresh sync_uuid.
func updateGameData(e *core.RequestEvent) error {
	var body struct {
		DeviceUUID string          `json:"device_uuid"`
		SyncUUID   string          `json:"sync_uuid"`
		GameData   json.RawMessage `json:"game_data"`
	}
	if err := e.BindBody(&body); err != nil {
		return e.BadRequestError("invalid body", err)
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)

	var record *core.Record
	var err error

	if body.SyncUUID != "" {
		record, err = e.App.FindFirstRecordByFilter("game_data", "sync_uuid = {:sync_uuid}", map[string]any{"sync_uuid": body.SyncUUID})
	}
	if err != nil && record == nil && body.DeviceUUID != "" {
		record, err = e.App.FindFirstRecordByFilter("game_data", "created_by = {:device_uuid}", map[string]any{"device_uuid": body.DeviceUUID})
	}

	if record == nil {
		collection, collErr := e.App.FindCollectionByNameOrId("game_data")
		if collErr != nil {
			return e.InternalServerError("collection not found", collErr)
		}
		record = core.NewRecord(collection)
		record.Set("created_by", body.DeviceUUID)
		record.Set("sync_uuid", newUUID())
	}

	record.Set("game_data", body.GameData)
	record.Set("modified_by", body.DeviceUUID)
	record.Set("modified_at", now)

	if err := e.App.Save(record); err != nil {
		return e.InternalServerError("save failed", err)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"new_sync_uuid":   record.GetString("sync_uuid"),
		"new_modified_at": record.GetString("modified_at"),
	})
}

// GET /api/sync/game-data/share/{syncHash}: public read-only lookup by hash.
func getGameDataShare(e *core.RequestEvent) error {
	syncHash := e.Request.PathValue("syncHash")

	records, err := e.App.FindRecordsByFilter("game_data", "sync_hash = {:sync_hash}", "", 1, 0, map[string]any{"sync_hash": syncHash})
	if err != nil || len(records) == 0 {
		return e.JSON(http.StatusOK, nil)
	}

	return e.JSON(http.StatusOK, records[0].GetRaw("game_data"))
}
