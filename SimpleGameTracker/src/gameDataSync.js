import { ApplyServerGames } from './gameDataCache.js'

const env = import.meta.env

const PB_API = env.VITE_PB_API

const deviceUUID = localStorage.deviceUUID ? localStorage.deviceUUID : crypto.randomUUID()
if (!localStorage.deviceUUID) {
  localStorage.deviceUUID = deviceUUID
}

let syncingGameData = false
export async function GetSyncGameData() {
  if (!navigator.onLine)
    return

  let syncUUID = localStorage.syncUUID
  if (syncUUID != undefined) {
    syncingGameData = true
    let data = null
    let failed = false
    try {
      const response = await fetch(`${PB_API}/api/sync/game-data?sync_uuid=${encodeURIComponent(syncUUID)}&local_modified_at=${encodeURIComponent(localStorage.syncModifiedAt ?? '')}`)
      failed = !response.ok
      if (!failed) {
        try {
          data = await response.json()
        }
        catch {
          failed = true
        }
      }
    }
    catch (err) {
      console.error("get_game_data failed:", err)
      failed = true
    }

    if (!failed && data != null && data.length > 0) {
      // The server only returns the collection when its modified date is
      // newer than the client's last sync, so the server version wins here.
      await ApplyServerGames(data[0].games)
      localStorage.syncModifiedAt = data[0].modified_at
      console.log("Get synced game data")
    }
    else {
      if (failed)
        console.log("get_game_data failed:", "request error")
      if (localStorage.syncModifiedAt == "1970-01-01T00:00:00.000Z")
        localStorage.removeItem('syncUUID')
      console.log("No synced game data")
    }
    syncingGameData = false
  }
}
GetSyncGameData()

export async function UpdateSyncGameData() {
  if (!navigator.onLine)
    return

  let syncUUID = localStorage.syncUUID
  if (syncUUID != undefined) {
    let data = null
    let failed = false
    try {
      const response = await fetch(`${PB_API}/api/sync/game-data`, {
        method: 'POST',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({
          "device_uuid": deviceUUID,
          "sync_uuid": syncUUID,
          "games": Object.values(JSON.parse(localStorage.gameData ?? '{}'))
            .filter(game => {
              const id = Number(game.id)
              return Number.isFinite(id) && game.state != undefined && game.state != 'none'
            })
            .map(game => ({
              "game_id": Number(game.id),
              "state": game.state
            }))
        })
      })
      failed = !response.ok
      if (!failed) {
        try {
          data = await response.json()
        }
        catch {
          failed = true
        }
      }
    }
    catch (err) {
      console.error("update_game_data failed:", err)
      failed = true
    }

    if (!failed && data != null) {
      localStorage.syncUUID = data.new_sync_uuid
      localStorage.syncModifiedAt = data.new_modified_at
      console.log("Update sync game data")
    }
    else {
      localStorage.removeItem('syncUUID')
      console.log("No synced game data")
    }
  }
}

export async function BusyWait(test) {
  const delayMs = 500;
  while(!test()) await new Promise(resolve => setTimeout(resolve, delayMs));
}

let loadingGameDataShared = false
export async function GetSharedGameData(syncHash) {
  loadingGameDataShared = true
  if (sessionStorage.gameDataShared == undefined) {
    let data = null
    try {
      const response = await fetch(`${PB_API}/api/sync/game-data/share/${encodeURIComponent(syncHash)}`)
      if (response.ok) {
        try {
          data = await response.json()
        }
        catch {
          data = null
        }
      }
    }
    catch (err) {
      console.error("get_game_data_share failed:", err)
    }

    if (data != null && data.length > 0) {
      // The share endpoint only returns ids + states; fetch the full IGDB
      // data for those ids and build the game objects the UI expects.
      const gameIds = data.map(game => game.game_id)
      let igdbData = []
      try {
        const response = await fetch(`${PB_API}/api/igdb/game-info`, {
          method: 'POST',
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
          },
          body: JSON.stringify({
            "gameId": gameIds
          })
        })
        if (response.ok) {
          igdbData = await response.json()
        }
      }
      catch (err) {
        console.error("get_game_data_share igdb fetch failed:", err)
      }

      const gameDataShared = {}
      for (let game of igdbData) {
        gameDataShared[game.id] = game
      }
      for (let game of data) {
        if (gameDataShared[game.game_id] != undefined) {
          gameDataShared[game.game_id].state = game.state
        }
      }
      sessionStorage.gameDataShared = JSON.stringify(gameDataShared)
      console.log("Get shared game data")
    }
    else {
      console.log("No shared game data")
    }
  }
  loadingGameDataShared = false
}

export const loadingGameData = (() => syncingGameData || loadingGameDataShared )