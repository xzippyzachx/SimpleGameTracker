import { createClient } from '@supabase/supabase-js'
const env = import.meta.env

const supabase = createClient(env.VITE_SUPABASE_API, env.VITE_SUPABASE_ANON)

let syncUUID = localStorage.syncUUID
const deviceUUID = localStorage.deviceUUID ? localStorage.deviceUUID : crypto.randomUUID()
if (!localStorage.deviceUUID) {
  localStorage.deviceUUID = deviceUUID
}

export let syncingGameData = false
async function GetSyncGameData() {
  if (syncUUID != undefined) {
    syncingGameData = true
    const { data, error } = await supabase.rpc('get_game_data', { device_uuid: deviceUUID, sync_uuid: syncUUID })
    if (data != null) {
      localStorage.gameData = JSON.stringify(data)
      console.log("Get sync game data")
    }
    syncingGameData = false
  }
}
GetSyncGameData()

export async function UpdateSyncGameData(firstSync) {
  if (syncUUID != undefined || firstSync) {
    const { data, error } = await supabase.rpc('update_game_data', { device_uuid: deviceUUID, sync_uuid: syncUUID, game_data: JSON.parse(localStorage.gameData) })
    if (data != null) {
      syncUUID = data
      localStorage.syncUUID = syncUUID
      console.log("Update sync game data")
    }
  }
}

export async function BusyWait(test) {
  const delayMs = 500;
  while(!test()) await new Promise(resolve => setTimeout(resolve, delayMs));
}