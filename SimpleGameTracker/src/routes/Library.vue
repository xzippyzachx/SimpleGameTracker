<script setup>
import { ref, computed, useTemplateRef, watch, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'

import GameCard from '../components/GameCard.vue'
import StateNav from '../components/StateNav.vue'

import { loadingGameData, BusyWait, GetSharedGameData } from '../gameDataSync.js'

const searchName = ref("")
const libraryGames = ref([])
const loading = ref(false)

const searchNameQuery = computed(() => route.query.searchName ?? '')
const stateQuery = computed(() => route.query.state ?? '')

function UpdateRouteQuery() {
  router.replace({ query: { searchName: searchName.value, state: state.value } })
}

const router = useRouter()
const route = useRoute()

if (route.params.syncHash != undefined)
{
  GetSharedGameData(route.params.syncHash)
}

const Shared = (() => sessionStorage.gameDataShared != undefined || route.name == 'shared')

async function GetGameData() {
  loading.value = true
  await BusyWait(() => loadingGameData() == false)

  let gameData = []
  if (sessionStorage.gameDataShared != undefined) {
    gameData = Object.values(JSON.parse(sessionStorage.gameDataShared))
  }
  else if (localStorage.gameData) {
    gameData = Object.values(JSON.parse(localStorage.gameData))
  }
  libraryGames.value = gameData.sort((a, b) => a.name.localeCompare(b.name))
  loading.value = false
}
GetGameData()

const stateNav = useTemplateRef('state-nav-ref')
const state = computed(() => stateNav.value?.state)
onMounted(() => {
  watch(state, async () => {
    UpdateRouteQuery()
  })

  if (searchNameQuery.value)
  {
    searchName.value = searchNameQuery.value
    SearchGame()
  }

  if (stateQuery.value)
  {
    stateNav.value.state = stateQuery.value
  }
})

function SearchGame() {
  document.getElementById('search-input').blur()
  UpdateRouteQuery()
}

function LibraryFiltered() {
  return libraryGames.value
  .filter(game => game.name.toLowerCase().includes(searchName.value.toLowerCase()) || game.alternativeNames?.toLowerCase().includes(searchName.value.toLowerCase()))
  .filter(game => state.value == 'none' || game.state == state.value)
}

function AddGame() {
  router.push({ path: '/search', query: { searchName: searchName.value } })
}

</script>

<template>
  <header>
    <StateNav ref="state-nav-ref"/>
    <div>
      <input id="search-input" v-model="searchName" @keyup.enter="SearchGame()" v-on:blur="SearchGame()" placeholder="Game Name"></input>
      <button id="search" @click="SearchGame()"><font-awesome-icon icon="fa-solid fa-search" /></button>
    </div>
  </header>

  <main>
    <span class="loader" v-if="loading"></span>
    <div v-for="game in LibraryFiltered()">
      <GameCard :id="game.id" :name="game.name" :cover="game.cover" :state="game.state"/>
    </div>
    <button v-if="!loading && !Shared() && searchName != '' && LibraryFiltered().length == 0" @click="AddGame()">Add Game<font-awesome-icon icon="fa-regular fa-square-plus" /></button>
  </main>
</template>

<style scoped>
  main {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
  }

  header div {
    padding-top: 10px;
  }

  #search {
    width: 40px;
    margin: 0px;
    padding: 0px;

    border-radius: 0px;
  }

  input {
    height: 40px;
    width: 300px;
    margin: auto 0px auto 0px;
    font-size: x-large;
    padding: 0px 10px;
  }

</style>