<script setup>
import { ref, computed, useTemplateRef, watch, onMounted } from 'vue'

import GameCard from '../components/GameCard.vue'
import StateNav from '../components/StateNav.vue'

const searchName = ref("")
const libraryGames = ref([])

function GetGameData(state) {
  if (localStorage.gameData) {
    libraryGames.value = Object.values(JSON.parse(localStorage.gameData))
    .filter(game => state == 'none' || game.state == state)
    .sort((a, b) => a.name.localeCompare(b.name))
  }
}
GetGameData('none')

const stateNav = useTemplateRef('state-nav-ref')
const state = computed(() => stateNav.value?.state)
onMounted(() => {
  watch(state, async (value, lastValue) => {
    GetGameData(value)
  })
})

</script>

<template>
  <header>
    <StateNav ref="state-nav-ref"/>
    <div>
      <input v-model="searchName" placeholder="Game Name"></input>
      <button id="search"><font-awesome-icon icon="fa-solid fa-search" /></button>
    </div>
  </header>

  <main>
    <div v-for="game in libraryGames">
      <GameCard :id="game.id" :name="game.name" :cover="game.cover" :state="game.state" v-if="game.name.toLowerCase().includes(searchName.toLowerCase())"/>
    </div>
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