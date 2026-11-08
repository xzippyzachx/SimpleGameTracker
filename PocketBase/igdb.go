package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/pocketbase/pocketbase/core"
)

const (
	igdbTokenURL    = "https://id.twitch.tv/oauth2/token"
	igdbGamesURL    = "https://api.igdb.com/v4/games"
	igdbCoverBase   = "https://images.igdb.com/igdb/image/upload/t_cover_big/%s.jpg"
	igdbTokenTTL    = 50 * time.Minute
	igdbHTTPTimeout = 15 * time.Second
)

var (
	igdbTokenMu     sync.Mutex
	igdbTokenValue  string
	igdbTokenExpiry time.Time
)

// igdbToken fetches a Twitch client-credentials token and caches it in memory
// for its lifetime (minus a safety margin), avoiding a round trip per request.
func igdbToken(app core.App) (string, error) {
	igdbTokenMu.Lock()
	defer igdbTokenMu.Unlock()

	if igdbTokenValue != "" && time.Now().Before(igdbTokenExpiry) {
		return igdbTokenValue, nil
	}

	clientID := os.Getenv("IGDB_CLIENT_ID")
	if clientID == "" {
		clientID = os.Getenv("CLIENT_ID")
	}
	clientSecret := os.Getenv("IGDB_CLIENT_SECRET")
	if clientSecret == "" {
		clientSecret = os.Getenv("CLIENT_SECRET")
	}
	if clientID == "" || clientSecret == "" {
		return "", fmt.Errorf("IGDB_CLIENT_ID and IGDB_CLIENT_SECRET must be set")
	}

	url := fmt.Sprintf("%s?client_id=%s&client_secret=%s&grant_type=client_credentials", igdbTokenURL, clientID, clientSecret)
	client := &http.Client{Timeout: igdbHTTPTimeout}
	resp, err := client.Post(url, "application/x-www-form-urlencoded", nil)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var token struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", err
	}
	if token.AccessToken == "" {
		return "", fmt.Errorf("empty access token from Twitch")
	}

	igdbTokenValue = token.AccessToken
	igdbTokenExpiry = time.Now().Add(time.Duration(token.ExpiresIn)*time.Second - igdbTokenTTL)
	return igdbTokenValue, nil
}

// igdbQuery posts an Apicalypse query to the IGDB games endpoint.
func igdbQuery(app core.App, body string) ([]byte, error) {
	token, err := igdbToken(app)
	if err != nil {
		return nil, err
	}

	clientID := os.Getenv("IGDB_CLIENT_ID")
	if clientID == "" {
		clientID = os.Getenv("CLIENT_ID")
	}

	req, err := http.NewRequest(http.MethodPost, igdbGamesURL, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Client-ID", clientID)
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: igdbHTTPTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("igdb status %d: %s", resp.StatusCode, string(data))
	}
	return data, nil
}

func registerIgdbRoutes(se *core.ServeEvent) {
	g := se.Router.Group("/api/igdb")
	g.POST("/search-games", igdbSearchGames)
	g.POST("/game-info", igdbGameInfo)
}

// POST /api/igdb/search-games  Body: {"name": "..."}
func igdbSearchGames(e *core.RequestEvent) error {
	var body struct {
		Name string `json:"name"`
	}
	if err := e.BindBody(&body); err != nil {
		return e.BadRequestError("invalid body", err)
	}

	query := fmt.Sprintf(
		`search "%s"; limit 100; fields name, cover.image_id; where game_type = (0,1,2,4,8,9,11) & version_parent = null & cover != null;`,
		body.Name,
	)
	data, err := igdbQuery(e.App, query)
	if err != nil {
		return e.InternalServerError("igdb query failed", err)
	}

	var games []map[string]any
	if err := json.Unmarshal(data, &games); err != nil {
		return e.InternalServerError("igdb response parse failed", err)
	}

	sorted := matchSorter(games, body.Name, "name")
	for _, game := range sorted {
		if cover, ok := game["cover"].(map[string]any); ok {
			if imageID, ok := cover["image_id"].(string); ok {
				game["cover"] = fmt.Sprintf(igdbCoverBase, imageID)
			}
		}
	}

	return e.JSON(http.StatusOK, sorted)
}

// POST /api/igdb/game-info  Body: {"gameId": 123} or {"gameId": [1,2,3]}
// Response reshaping: alternativeNames joined with "|", releaseDate,
// official/steam URLs, genres/themes flattened to name arrays.
func igdbGameInfo(e *core.RequestEvent) error {
	var body struct {
		GameID json.RawMessage `json:"gameId"`
	}
	if err := e.BindBody(&body); err != nil {
		return e.BadRequestError("invalid body", err)
	}

	// Accept string ids as before, numeric ids (the search results carry
	// numeric ids), and arrays of either.
	gameIDs := []string{}
	trimmed := bytes.TrimSpace(body.GameID)
	if bytes.HasPrefix(trimmed, []byte("[")) {
		var rawList []json.RawMessage
		if err := json.Unmarshal(body.GameID, &rawList); err != nil {
			return e.BadRequestError("invalid gameId", err)
		}
		for _, raw := range rawList {
			id, err := parseGameID(raw)
			if err != nil {
				return e.BadRequestError("invalid gameId", err)
			}
			gameIDs = append(gameIDs, id)
		}
	} else {
		id, err := parseGameID(trimmed)
		if err != nil {
			return e.BadRequestError("invalid gameId", err)
		}
		gameIDs = append(gameIDs, id)
	}
	length := len(gameIDs)
	gameIDList := strings.Join(gameIDs, ",")

	query := fmt.Sprintf(
		`where id = (%s); limit 500; fields name, cover.image_id, summary, alternative_names.name, websites.url, websites.category, genres.name, themes.name, rating, first_release_date;`,
		gameIDList,
	)
	data, err := igdbQuery(e.App, query)
	if err != nil {
		return e.InternalServerError("igdb query failed", err)
	}

	var gameData []map[string]any
	if err := json.Unmarshal(data, &gameData); err != nil {
		return e.InternalServerError("igdb response parse failed", err)
	}

	if len(gameData) != length {
		return e.String(http.StatusBadRequest, "Invalid gameID")
	}

	for _, game := range gameData {
		if cover, ok := game["cover"].(map[string]any); ok {
			if imageID, ok := cover["image_id"].(string); ok {
				game["cover"] = fmt.Sprintf(igdbCoverBase, imageID)
			}
		}

		if altNames, ok := game["alternative_names"].([]any); ok {
			names := make([]string, 0, len(altNames))
			for _, alt := range altNames {
				if altMap, ok := alt.(map[string]any); ok {
					if name, ok := altMap["name"].(string); ok {
						names = append(names, name)
					}
				}
			}
			delete(game, "alternative_names")
			game["alternativeNames"] = strings.Join(names, "|")
		}

		if releaseDate, ok := game["first_release_date"]; ok {
			delete(game, "first_release_date")
			game["releaseDate"] = releaseDate
		}

		if websites, ok := game["websites"].([]any); ok {
			delete(game, "websites")
			for _, w := range websites {
				if website, ok := w.(map[string]any); ok {
					category, _ := website["category"].(float64)
					url, _ := website["url"].(string)
					if category == 1 {
						game["official"] = url
					}
					if category == 13 {
						game["steam"] = url
					}
				}
			}
		}

		flattenNames(game, "genres")
		flattenNames(game, "themes")
	}

	if length == 1 {
		return e.JSON(http.StatusOK, gameData[0])
	}
	return e.JSON(http.StatusOK, gameData)
}

// parseGameID converts a raw JSON game id (string or number) into its
// canonical string form.
func parseGameID(raw []byte) (string, error) {
	var num float64
	if err := json.Unmarshal(raw, &num); err == nil {
		return strconv.FormatFloat(num, 'f', -1, 64), nil
	}
	var str string
	if err := json.Unmarshal(raw, &str); err != nil {
		return "", fmt.Errorf("gameId must be a string or number")
	}
	if _, err := strconv.ParseFloat(str, 64); err != nil {
		return "", fmt.Errorf("gameId must be numeric")
	}
	return str, nil
}

// flattenNames replaces a list of {name: "..."} objects with a plain array
// of names.
func flattenNames(game map[string]any, key string) {
	items, ok := game[key].([]any)
	if !ok {
		return
	}
	names := make([]string, 0, len(items))
	for _, item := range items {
		if m, ok := item.(map[string]any); ok {
			if name, ok := m["name"].(string); ok {
				names = append(names, name)
			}
		}
	}
	game[key] = names
}
