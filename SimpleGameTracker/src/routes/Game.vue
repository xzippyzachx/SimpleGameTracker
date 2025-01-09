<script setup>
import { ref, computed, useTemplateRef, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'

import StateNav from '../components/StateNav.vue'

const gameData = ref({})
const route = useRoute()

function GetGameInfo() {
  fetch('https://sdekcxxvsnnzypebfpcr.supabase.co/functions/v1/game-info', {
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
  })
  .catch(err => {
    console.error(err);
  });
}

function SetGameState(state) {
  gameData.value.state = state
  CacheGameData()
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

</script>

<template>
  <header>
    <StateNav v-if="gameData && Object.keys(gameData).length > 0" ref="state-nav-ref"/>
  </header>

  <main>
    <img :src="gameData.cover"></img>
    <h1>{{ gameData.name }}</h1>
    <p>{{ gameData.summary }}</p>
  </main>
</template>

<style scoped>
  main {
    display: flex;
    flex-direction: column;
    justify-content: center;
  }

  img {
    max-width: 400px;
    border-radius: 10px;
    border: solid;
    border-width: 2.5px;
  }

</style>