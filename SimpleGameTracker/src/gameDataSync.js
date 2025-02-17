import { createClient } from '@supabase/supabase-js'
const env = import.meta.env

const supabase = createClient(env.VITE_SUPABASE_API, env.VITE_SUPABASE_ANON)

const deviceUUID = localStorage.deviceUUID ? localStorage.deviceUUID : crypto.randomUUID()
if (!localStorage.deviceUUID) {
  localStorage.deviceUUID = deviceUUID
}

let syncingGameData = false
export async function GetSyncGameData() {
  let syncUUID = localStorage.syncUUID
  if (syncUUID != undefined) {
    syncingGameData = true
    const { data, error } = await supabase.rpc('get_game_data', { device_uuid: deviceUUID, sync_uuid: syncUUID, local_modifed_at: localStorage.syncModifiedAt })
    if (data != null && data.length > 0) {
      localStorage.gameData = JSON.stringify(data[0].game_data)
      localStorage.syncModifiedAt = data[0].modified_at
      console.log("Get synced game data")
    }
    else {
      if (localStorage.syncModifiedAt == "1970-01-01T00:00:00.000Z")
        localStorage.removeItem('syncUUID')
      console.log("No synced game data")
    }
    if (error) {
      localStorage.removeItem('syncUUID')
      console.log(error)
    }
    syncingGameData = false
  }
}
GetSyncGameData()

export async function UpdateSyncGameData() {
  let syncUUID = localStorage.syncUUID
  if (syncUUID != undefined) {
    const { data, error } = await supabase.rpc('update_game_data', { device_uuid: deviceUUID, sync_uuid: syncUUID, game_data: JSON.parse(localStorage.gameData) })
    if (data != null) {
      localStorage.syncUUID = data[0].new_sync_uuid
      localStorage.syncModifiedAt = data[0].new_modified_at
      console.log("Update sync game data")
    }
    else {
      localStorage.removeItem('syncUUID')
      console.log("No synced game data")
    }
    if (error) {
      localStorage.removeItem('syncUUID')
      console.log(error)
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
    const { data, error } = await supabase.rpc('get_game_data_share', { sync_hash: syncHash })
    if (data != null) {
      sessionStorage.gameDataShared = JSON.stringify(data)
      console.log("Get shared game data")
    }
    else {
      console.log("No shared game data")
    }
    if (error) {
      console.log(error)
    }
  }
  loadingGameDataShared = false
}

export const loadingGameData = (() => syncingGameData || loadingGameDataShared )