<script setup>
import { ref, computed, useTemplateRef, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'

import StateNav from '../components/StateNav.vue'

import { UpdateSyncGameData } from '../gameDataSync.js'

const env = import.meta.env

const gameData = ref({})
const route = useRoute()
const searching = ref(false)
let firstState = true

function GetGameInfo() {
  searching.value = true
  fetch(`${env.VITE_SUPABASE_API}/functions/v1/game-info`, {
    method: 'POST',
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      "gameId": route.params.gameId
    })
  }).then(response => {
    response.json().then(res => {
      gameData.value = res
    })
    searching.value = false
  })
  .catch(err => {
    console.error(err)
    searching.value = false
  })
}

function SetGameState(state) {
  gameData.value.state = state
  if (!firstState)
    CacheGameData()
  firstState = false
}

function CacheGameData() {
  let allGameData = localStorage.gameData ? JSON.parse(localStorage.gameData) : {}
  if (gameData.value.state != 'none') {
    allGameData[gameData.value.id] = gameData.value
  }
  else {
    delete allGameData[gameData.value.id]
  }
  localStorage.gameData = JSON.stringify(allGameData)
  UpdateSyncGameData()
}

if (localStorage.gameData) {
  gameData.value = JSON.parse(localStorage.gameData)[route.params.gameId] || {}
}

if (Object.keys(gameData.value).length == 0) {
  GetGameInfo()
}

const stateNav = useTemplateRef('state-nav-ref')
onMounted(() => {
  const state = computed(() => stateNav.value?.state)
  watch(state, async (value, lastValue) => {
    SetGameState(value)
  })

  if (gameData.value.state)
  {
    stateNav.value.state = gameData.value.state
  }
})

const stateClass = computed(() => ({
  wishlist: gameData.value.state == 'wishlist',
  playing: gameData.value.state == 'playing',
  completed: gameData.value.state == 'completed',
  shelved: gameData.value.state == 'shelved',
}))

</script>

<template>
  <header>
    <StateNav v-if="gameData && Object.keys(gameData).length > 0" ref="state-nav-ref"/>
  </header>

  <main>
    <span class="loader" v-if="searching"></span>
    <div id="content" v-if="!searching">
      <img :src="gameData.cover" :class="stateClass"></img>
      <h1>{{ gameData.name }}</h1>
      <p>{{ gameData.summary }}</p>
    </div>
  </main>
</template>

<style scoped>
  main {
    padding: 20px 40px;
  }

  #content {
    display: flex;
    flex-direction: column;
    justify-content: center;
    gap: 10px;
  }

  img {
    max-width: 400px;
    aspect-ratio: 0.75 / 1;
    border-radius: 10px;
    border: solid;
    border-width: 2.5px;
  }

  h1 {
    line-height: normal;
  }

  .wishlist {
    border-color: var(--vt-wishlist-c);
  }
  .playing {
    border-color: var(--vt-playing-c);
  }
  .completed {
    border-color: var(--vt-completed-c);
  }
  .shelved {
    border-color: var(--vt-shelved-c);
  }

</style>