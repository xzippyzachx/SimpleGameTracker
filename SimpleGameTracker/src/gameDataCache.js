const env = import.meta.env

import { UpdateSyncGameData } from './gameDataSync.js'

export function CacheGameData(gameData) {
  let allGameData = localStorage.gameData ? JSON.parse(localStorage.gameData) : {}
  if (gameData.state != 'none') {
    allGameData[gameData.id] = gameData
  }
  else {
    delete allGameData[gameData.id]
  }
  localStorage.gameData = JSON.stringify(allGameData)
  UpdateSyncGameData()
}

export function RefreshCache() {
  let cachedGameData = JSON.parse(localStorage.gameData ?? '{}')
  
  fetch(`${env.VITE_PB_API}/api/igdb/game-info`, {
    method: 'POST',
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      "gameId": Object.keys(cachedGameData)
    })
  }).then(response => {
    response.json().then(res => {
      let gameData = res

      for (let i = 0; i < gameData.length; i++) {
        let gameId = gameData[i].id
        let state = cachedGameData[gameId].state

        cachedGameData[gameId] = gameData[i]
        cachedGameData[gameId].state = state
      }

      localStorage.gameData = JSON.stringify(cachedGameData)
      UpdateSyncGameData()
    })
  })
  .catch(err => {
    console.error(err)
  })
}

export function ChangeGameState(gameId, newState) {
  const gameIdStr = String(gameId)
  let allGameData = localStorage.gameData ? JSON.parse(localStorage.gameData) : {}
  let gameData = (allGameData)[gameIdStr] || {}
  if (Object.keys(gameData).length > 0) {
    gameData.state = newState
    CacheGameData(gameData)
  }
  else {
    fetch(`${env.VITE_PB_API}/api/igdb/game-info`, {
      method: 'POST',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        "gameId": gameIdStr
      })
    }).then(response => {
      response.json().then(res => {
        // Ignore error payloads (e.g. invalid game id) — only cache real
        // games so the sync list doesn't get polluted.
        if (res == null || res.id == undefined || res.name == undefined)
          return

        gameData = res
        gameData.state = newState
        CacheGameData(gameData)
      })
    })
    .catch(err => {
      console.error(err)
    })
  }
}

// Apply a server games list ([{game_id, state}, ...]) to the local cache.
// The server is the source of truth: states are updated, missing games are
// fetched from IGDB, and locally-cached games absent from the list are
// removed. Only called when the server's modified date is newer than the
// client's last sync.
export async function ApplyServerGames(games) {
  let allGameData = localStorage.gameData ? JSON.parse(localStorage.gameData) : {}

  const serverStates = {}
  for (let game of games) {
    serverStates[game.game_id] = game.state
  }

  // Remove locally-cached games the server doesn't have.
  for (let gameId of Object.keys(allGameData)) {
    if (serverStates[Number(gameId)] == undefined) {
      delete allGameData[gameId]
    }
  }

  // Collect ids needing full IGDB data (not cached yet, non-none state).
  const missingIds = []
  for (let [gameId, state] of Object.entries(serverStates)) {
    if (state == 'none') {
      delete allGameData[gameId]
    }
    else if (allGameData[gameId] == undefined) {
      missingIds.push(gameId)
    }
    else {
      allGameData[gameId].state = state
    }
  }

  if (missingIds.length > 0) {
    try {
      const response = await fetch(`${env.VITE_PB_API}/api/igdb/game-info`, {
        method: 'POST',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          "gameId": missingIds
        })
      })
      if (response.ok) {
        const igdbData = await response.json()
        for (let game of igdbData) {
          game.state = serverStates[game.id]
          allGameData[game.id] = game
        }
      }
    }
    catch (err) {
      console.error(err)
    }
  }

  localStorage.gameData = JSON.stringify(allGameData)
}