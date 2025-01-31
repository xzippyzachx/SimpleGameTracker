<script setup>
const env = import.meta.env

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

function RefreshCache() {
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
    })
  })
  .catch(err => {
    console.error(err)
  })
}

</script>

<template>
  <main>
    <button @click="DownloadData()">Download Data</button>
    <button @click="UploadData()">Upload Data</button>
    <button @click="RefreshCache()">Refresh Cache</button>
  </main>
</template>

<style scoped>
  main {
    display: flex;
    flex-direction: column;
    gap: 10px;
    margin-top: 10px;
  }

  button {
    width: 200px;
  }
</style>