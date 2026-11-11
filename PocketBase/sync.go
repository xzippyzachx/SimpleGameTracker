package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/core"
)

const (
	// pbDateLayout is PocketBase's stored datetime format (locale-independent):
	// space-separated, millisecond precision, Z suffix.
	pbDateLayout = "2006-01-02 15:04:05.000Z"

	// pairCodeAlphabet omits visually ambiguous characters (I, O, 0, 1) so a
	// code read off a screen can be typed without confusion. 32 characters
	// divides 256 evenly, so byte % len(alphabet) is unbiased.
	pairCodeAlphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	pairCodeLength   = 8
	pairCodeTTL      = 5 * time.Minute
)

// pairCodeEntry is a short-lived, single-use code that can be exchanged for a
// sync key on another device without putting the key itself in a URL.
type pairCodeEntry struct {
	syncUUID  string
	expiresAt time.Time
}

var (
	pairCodesMu sync.Mutex
	pairCodes   = map[string]pairCodeEntry{}
)

// generatePairCode returns a random code drawn from pairCodeAlphabet.
func generatePairCode() (string, error) {
	buf := make([]byte, pairCodeLength)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	code := make([]byte, pairCodeLength)
	for i, b := range buf {
		code[i] = pairCodeAlphabet[int(b)%len(pairCodeAlphabet)]
	}
	return string(code), nil
}

// prunePairCodesLocked drops expired entries. Callers must hold pairCodesMu.
func prunePairCodesLocked(now time.Time) {
	for code, entry := range pairCodes {
		if now.After(entry.expiresAt) {
			delete(pairCodes, code)
		}
	}
}

// validGameStates are the states a game row may hold.
var validGameStates = map[string]bool{
	"wishlist":  true,
	"playing":   true,
	"completed": true,
	"shelved":   true,
}

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
	g.POST("/game-data/pair/create", createPairCode)
	g.POST("/game-data/pair/claim", claimPairCode)

	// Keep sync_hash in sync with sync_uuid (sha256 of the uuid as hex).
	se.App.OnRecordCreate("game_collection").BindFunc(func(e *core.RecordEvent) error {
		if e.Record.GetString("sync_hash") == "" {
			e.Record.Set("sync_hash", syncHash(e.Record.GetString("sync_uuid")))
		}
		return e.Next()
	})
	se.App.OnRecordUpdate("game_collection").BindFunc(func(e *core.RecordEvent) error {
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
// Returns the collection's games only when the server copy is newer than the
// client's local timestamp.
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

	records, err := e.App.FindRecordsByFilter("game_collection", filter, "", 1, 0, params)
	if err != nil || len(records) == 0 {
		return e.JSON(http.StatusOK, []any{})
	}

	r := records[0]
	games, err := collectionGames(e, r.Id)
	if err != nil {
		return e.InternalServerError("failed to load games", err)
	}
	return e.JSON(http.StatusOK, []map[string]any{{
		"created_by":  r.GetString("created_by"),
		"games":       games,
		"modified_at": r.GetString("modified_at"),
	}})
}

// collectionGames loads the game rows for a collection as
// [{"game_id": ..., "state": ...}, ...].
func collectionGames(e *core.RequestEvent, collectionId string) ([]map[string]any, error) {
	records, err := e.App.FindAllRecords("game", dbx.HashExp{"collection_id": collectionId})
	if err != nil {
		return nil, err
	}
	games := make([]map[string]any, 0, len(records))
	for _, r := range records {
		games = append(games, map[string]any{
			"game_id": r.GetFloat("game_id"),
			"state":   r.GetString("state"),
		})
	}
	return games, nil
}

// POST /api/sync/game-data
// Body: {"device_uuid": "...", "sync_uuid": "...", "games": [{"game_id": 123, "state": "playing"}, ...]}
// Update by sync_uuid, else update by created_by, else insert a new collection
// with a fresh sync_uuid. The submitted games list fully replaces the
// collection's existing game rows.
func updateGameData(e *core.RequestEvent) error {
	var body struct {
		DeviceUUID string `json:"device_uuid"`
		SyncUUID   string `json:"sync_uuid"`
		Games      []struct {
			GameID float64 `json:"game_id"`
			State  string  `json:"state"`
		} `json:"games"`
	}
	if err := e.BindBody(&body); err != nil {
		return e.BadRequestError("invalid body", err)
	}

	for _, game := range body.Games {
		if !validGameStates[game.State] {
			return e.BadRequestError("invalid state", fmt.Errorf("unknown state %q", game.State))
		}
	}

	now := time.Now().UTC().Format(time.RFC3339Nano)

	var record *core.Record
	var err error

	if body.SyncUUID != "" {
		record, err = e.App.FindFirstRecordByFilter("game_collection", "sync_uuid = {:sync_uuid}", map[string]any{"sync_uuid": body.SyncUUID})
	}
	if err != nil && record == nil && body.DeviceUUID != "" {
		record, err = e.App.FindFirstRecordByFilter("game_collection", "created_by = {:device_uuid}", map[string]any{"device_uuid": body.DeviceUUID})
	}

	if record == nil {
		collection, collErr := e.App.FindCollectionByNameOrId("game_collection")
		if collErr != nil {
			return e.InternalServerError("collection not found", collErr)
		}
		record = core.NewRecord(collection)
		record.Set("created_by", body.DeviceUUID)
		record.Set("sync_uuid", newUUID())
	}

	record.Set("modified_by", body.DeviceUUID)
	record.Set("modified_at", now)

	if err := e.App.Save(record); err != nil {
		return e.InternalServerError("save failed", err)
	}

	gameCollection, collErr := e.App.FindCollectionByNameOrId("game")
	if collErr != nil {
		return e.InternalServerError("collection not found", collErr)
	}

	// Full replace: drop the collection's existing game rows, then insert
	// the submitted list, atomically.
	replaceErr := e.App.RunInTransaction(func(txApp core.App) error {
		existing, err := txApp.FindAllRecords("game", dbx.HashExp{"collection_id": record.Id})
		if err != nil {
			return err
		}
		for _, gameRecord := range existing {
			if err := txApp.Delete(gameRecord); err != nil {
				return err
			}
		}
		for _, game := range body.Games {
			gameRecord := core.NewRecord(gameCollection)
			gameRecord.Set("collection_id", record.Id)
			gameRecord.Set("game_id", game.GameID)
			gameRecord.Set("state", game.State)
			if err := txApp.Save(gameRecord); err != nil {
				return err
			}
		}
		return nil
	})
	if replaceErr != nil {
		return e.InternalServerError("replace failed", replaceErr)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"new_sync_uuid":   record.GetString("sync_uuid"),
		"new_modified_at": record.GetString("modified_at"),
	})
}

// GET /api/sync/game-data/share/{syncHash}: public read-only lookup by hash.
// Returns the collection's games as [{"game_id": ..., "state": ...}, ...];
// the client fetches the full IGDB data for those ids itself.
func getGameDataShare(e *core.RequestEvent) error {
	syncHash := e.Request.PathValue("syncHash")

	records, err := e.App.FindRecordsByFilter("game_collection", "sync_hash = {:sync_hash}", "", 1, 0, map[string]any{"sync_hash": syncHash})
	if err != nil || len(records) == 0 {
		return e.JSON(http.StatusOK, nil)
	}

	games, err := collectionGames(e, records[0].Id)
	if err != nil {
		return e.InternalServerError("failed to load games", err)
	}
	return e.JSON(http.StatusOK, games)
}

// POST /api/sync/game-data/pair/create
// Body: {"sync_uuid": "..."}
// Returns a short-lived, single-use code that another device can exchange for
// the sync key, so the key itself never has to travel in a URL.
func createPairCode(e *core.RequestEvent) error {
	var body struct {
		SyncUUID string `json:"sync_uuid"`
	}
	if err := e.BindBody(&body); err != nil {
		return e.BadRequestError("invalid body", err)
	}

	syncUUID := strings.TrimSpace(body.SyncUUID)
	if syncUUID == "" {
		return e.BadRequestError("sync_uuid is required", nil)
	}

	records, err := e.App.FindRecordsByFilter("game_collection", "sync_uuid = {:sync_uuid}", "", 1, 0, map[string]any{"sync_uuid": syncUUID})
	if err != nil {
		return e.InternalServerError("lookup failed", err)
	}
	if len(records) == 0 {
		return e.BadRequestError("unknown sync_uuid", nil)
	}

	code, err := generatePairCode()
	if err != nil {
		return e.InternalServerError("failed to generate code", err)
	}

	now := time.Now()
	expiresAt := now.Add(pairCodeTTL)

	pairCodesMu.Lock()
	prunePairCodesLocked(now)
	pairCodes[code] = pairCodeEntry{syncUUID: syncUUID, expiresAt: expiresAt}
	pairCodesMu.Unlock()

	return e.JSON(http.StatusOK, map[string]any{
		"code":       code,
		"expires_at": expiresAt.UTC().Format(time.RFC3339),
	})
}

// POST /api/sync/game-data/pair/claim
// Body: {"code": "ABCD2345"}
// Exchanges a valid, unexpired, unused code for the sync key. Unknown, expired
// and already-used codes all return the same error so the endpoint cannot be
// used to probe which codes exist.
func claimPairCode(e *core.RequestEvent) error {
	var body struct {
		Code string `json:"code"`
	}
	if err := e.BindBody(&body); err != nil {
		return e.BadRequestError("invalid body", err)
	}

	code := strings.ToUpper(strings.TrimSpace(body.Code))
	if code == "" {
		return e.BadRequestError("code is required", nil)
	}

	now := time.Now()

	pairCodesMu.Lock()
	prunePairCodesLocked(now)
	entry, ok := pairCodes[code]
	if ok {
		delete(pairCodes, code)
	}
	pairCodesMu.Unlock()

	if !ok {
		return e.BadRequestError("invalid or expired code", nil)
	}

	return e.JSON(http.StatusOK, map[string]any{
		"sync_uuid": entry.syncUUID,
	})
}
