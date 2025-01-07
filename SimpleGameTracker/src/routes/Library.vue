<script setup>
import { ref, computed, useTemplateRef, watch, onMounted} from 'vue'

import GameCard from '../components/GameCard.vue'
import StateNav from '../components/StateNav.vue'

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
    <div>
      <RouterLink :to="{ name: 'search' }">
        <button id="add-game-btn">Add Game</button>
      </RouterLink>
    </div>
    <StateNav ref="state-nav-ref"/>
  </header>

  <main>
    <div v-for="game in libraryGames">
      <GameCard :id="game.id" :name="game.name" :cover="game.cover" :state="game.state"/>
    </div>
  </main>
</template>

<style scoped>
  main {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
  }

  #add-game-btn {
    background-color: var(--vt-c-white);
  }

</style>