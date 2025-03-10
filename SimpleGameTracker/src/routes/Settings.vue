<script setup>
import { ref, computed } from 'vue'

import { RefreshCache } from '../gameDataCache.js'
import { UpdateSyncGameData, GetSyncGameData } from '../gameDataSync.js'

const syncUUID = ref(localStorage.syncUUID)
const syncUUIDInput = ref(syncUUID.value)
const theme = ref(localStorage.getItem('theme'))

function DownloadData() {
  const gameData = JSON.stringify(JSON.parse(localStorage.gameData), null, 4)

  const blob = new Blob([gameData], { type: "application/json" });
  const jsonObjectUrl = URL.createObjectURL(blob);

  const today = (new Date()).toISOString().split('T')[0]
  const filename = `simpleGameTracker_gameData_${today}.json`;
  const anchorEl = document.createElement("a");
  anchorEl.href = jsonObjectUrl;
  anchorEl.download = filename;

  anchorEl.click();
  URL.revokeObjectURL(jsonObjectUrl);
}

function UploadData() {
  FileUpload().then((jsonFile) => {
    const gameData = JSON.parse(jsonFile)
    localStorage.gameData = JSON.stringify(gameData)
  })
}

function FileUpload() {
  return new Promise(resolve => {
    const inputFileElement = document.createElement('input')
    inputFileElement.setAttribute('type', 'file')
    inputFileElement.setAttribute('multiple', 'false')
    inputFileElement.setAttribute('accept', '.json')
    
    inputFileElement.addEventListener(
      'change',
      async (event) => {
        const { files } = event.target
        if (!files) {
          return
        }

        const filePromises = [...files].map(file => file.text())

        resolve(await Promise.all(filePromises))
      },
      false,
    )
    inputFileElement.click()
  })
}

async function Sync() {
  if (syncUUID.value) {
    console.log("Disable Sync")
    syncUUID.value = undefined
    syncUUIDInput.value = ""
    localStorage.removeItem('syncUUID')
  }
  else {
    console.log("Enable Sync")
    if (!syncUUIDInput.value) {
      localStorage.syncUUID = localStorage.deviceUUID
      await UpdateSyncGameData()
      console.log("Sync Enabled")
    }
    else if (/^[0-9A-F]{8}-[0-9A-F]{4}-[4][0-9A-F]{3}-[89AB][0-9A-F]{3}-[0-9A-F]{12}$/i.test(syncUUIDInput.value)) {
      localStorage.syncUUID = syncUUIDInput.value
      localStorage.syncModifiedAt = "1970-01-01T00:00:00.000Z"
      await GetSyncGameData()
      console.log("Sync Enabled")
    }

    syncUUID.value = localStorage.syncUUID
    syncUUIDInput.value = syncUUID.value
  }
}

const syncButton = computed(() => {
  return syncUUID.value ? 'Disable Sync' : 'Enable Sync'
})

const syncInputDisabled = computed(() => {
  return syncUUID.value != undefined
})

const syncInputType = computed(() => {
  return syncUUID.value ? "password" : "text"
})

function CopySyncUUID() {
  navigator.clipboard.writeText(syncUUID.value);
}

async function ShareLink() {
  let syncHash = await Hash(syncUUID.value)
  navigator.share({
    title: document.title,
    text: "Share",
    url: `${window.location.origin}/#/shared${syncHash}`
  })
}

async function Hash(string) {
  const utf8 = new TextEncoder().encode(string)
  const hashBuffer = await crypto.subtle.digest('SHA-256', utf8)
  const hashArray = Array.from(new Uint8Array(hashBuffer))
  const hashHex = hashArray
    .map((bytes) => bytes.toString(16).padStart(2, '0'))
    .join('')
  return hashHex
}

function Theme() {
  if (theme.value == 'light') {
    theme.value = 'dark'
    localStorage.setItem('theme', 'dark')
    document.documentElement.setAttribute('data-theme', 'dark')
  }
  else if (theme.value == 'dark') {
    theme.value = 'light'
    localStorage.setItem('theme', 'light')
    document.documentElement.setAttribute('data-theme', 'light')
  }
}

const isOnline = computed(() => {
  return navigator.onLine
})

</script>

<template>
  <main>
    <button @click="DownloadData()">Download Data</button>
    <button @click="UploadData()">Upload Data</button>
    <button @click="RefreshCache()" :disabled="!isOnline">Refresh Cache</button>

    <button @click="Theme()">Theme<font-awesome-icon :icon="theme == 'dark' ? 'fa-solid fa-moon' : 'fa-solid fa-sun'" /></button>

    <div>
      <button @click="Sync()" :disabled="!isOnline">{{ syncButton }}</button>
      <input :type="syncInputType" v-model="syncUUIDInput" @keyup.enter="Sync()" placeholder="Sync Key" :disabled="syncInputDisabled || !isOnline"></input>
      <button id="copy" @click="CopySyncUUID()" v-if="syncInputDisabled"><font-awesome-icon icon="fa-solid fa-copy" /></button>
    </div>

    <p v-if="isOnline">Enable sync and copy your <i>Sync Key</i> into another device to sync your data across multiple devices.</p>
    <p v-if="isOnline">Warning: Your <i>Sync Key</i> allows read/write access to your data. Only share with trusted devices.</p>
    <p v-if="!isOnline">No internet connection found</p>

    <button @click="ShareLink()" v-if="syncInputDisabled && isOnline">Share<font-awesome-icon icon="fa-solid fa-share" /></button>
    <p v-if="syncInputDisabled && isOnline">This link can be shared with others to view your library.</p>
  </main>
</template>

<style scoped>
  main {
    display: flex;
    flex-direction: column;
    gap: 10px;
    margin-top: 10px;

    div {
      display: flex;
      flex-wrap: wrap;
    }
  }

  @media only screen and (min-width: 550px) {
    main {
      padding-left: 50px;
    }
  }

  button {
    width: 180px;
  }

  input {
    height: 40px;
    width: 300px;
    font-size: x-large;
    padding: 0px 10px;
    margin: auto 0px;
    margin-left: 10px;
  }

  p {
    padding-left: 10px;
  }

  #copy {
    width: 40px;
    height: 40px;
    margin: auto 0px;
    padding: 0px;
    border-radius: 0px;
  }
</style>