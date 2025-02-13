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
  let cachedGameData = JSON.parse(localStorage.gameData)
  
  fetch(`${env.VITE_SUPABASE_API}/functions/v1/game-info`, {
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
  let gameData = JSON.parse(localStorage.gameData)[gameId] || {}
  if (Object.keys(gameData).length > 0) {
    gameData.state = newState
    CacheGameData(gameData)
  }
  else {
    fetch(`${env.VITE_SUPABASE_API}/functions/v1/game-info`, {
      method: 'POST',
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        "gameId": gameId
      })
    }).then(response => {
      response.json().then(res => {
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