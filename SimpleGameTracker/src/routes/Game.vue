<script setup>
import { ref, computed, useTemplateRef, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'

import StateNav from '../components/StateNav.vue'

import { CacheGameData } from '../gameDataCache.js'

const env = import.meta.env

const gameData = ref({})
const route = useRoute()
const searching = ref(false)
let firstState = true

const Shared = (() => sessionStorage.gameDataShared != undefined || route.name == 'shared')

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
    CacheGameData(gameData.value)
  firstState = false
}

if (sessionStorage.gameDataShared)
{
  gameData.value = JSON.parse(sessionStorage.gameDataShared)[route.params.gameId] || {}
}
else if (localStorage.gameData) {
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

  if (gameData.value.state && stateNav.value)
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

const GetTags = computed(() => {
  let tags = []
  if (gameData.value.genres)
    tags.push(...gameData.value.genres)

  if (gameData.value.themes)
    tags.push(...gameData.value.themes)

  return tags
})

</script>

<template>
  <header v-if="!Shared()">
    <StateNav v-if="gameData && Object.keys(gameData).length > 0" ref="state-nav-ref"/>
  </header>

  <main>
    <span class="loader" v-if="searching"></span>
    <div id="content" v-if="!searching">
      <img :src="gameData.cover" :class="stateClass"></img>
      <div id="info">
        <h1>{{ gameData.name }}</h1>
        <span v-if="gameData.rating" class="info"><font-awesome-icon icon="fa-solid fa-star" /> {{ ((gameData.rating / 100) * 5).toFixed(2)}}/5</span>
        <span v-if="gameData.releaseDate" class="info"><font-awesome-icon icon="fa-solid fa-calendar-days" /> {{ new Date(gameData.releaseDate * 1000).toLocaleDateString("en", { year: 'numeric', month: 'long', day: 'numeric' }) }}</span>

        <div id="tags">
          <span class="tag" v-for="tag in GetTags">{{ tag }}</span>
        </div>

        <div id="links">
          <a v-if="gameData.official" :href="gameData.official" target="empty" class="link"><font-awesome-icon icon="fa-solid fa-globe" /> Official</a>
          <a v-if="gameData.steam" :href="gameData.steam" target="empty" class="link"><font-awesome-icon icon="fa-brands fa-steam" /> Steam</a>
        </div>

        <p>{{ gameData.summary }}</p>
      </div>
    </div>
  </main>
</template>

<style scoped>
  main {
    padding: 20px 40px;
  }

  #content {
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    gap: 20px;
  }

  #info {
    flex-basis: 780px;
    display: flex;
    flex-direction: column;
    gap: 10px;
  }

  #tags, #links {
    display: flex;
    flex-direction: row;
    flex-wrap: wrap;
    gap: 2.5px;
  }

  img {
    flex-basis: 400px;
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

  .tag, .link, .info {
    background-color: var(--color-background);
    padding: 2.5px 10px;
    border-radius: 10px;
    text-wrap: nowrap;
  }

  .info {
    margin-right: auto;
  }

  .link:hover {
    background-color: var(--color-btn-hover);
  }

</style>