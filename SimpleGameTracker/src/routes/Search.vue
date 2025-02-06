<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'

import GameCard from '../components/GameCard.vue'

const env = import.meta.env
const router = useRouter()
const route = useRoute()

const searchName = ref("")
const searchedGames = ref([])
const searching = ref(false)

const searchNameQuery = computed({
  get() {
    return route.query.searchName ?? ''
  },
  set(searchName) {
    router.replace({ query: { searchName } })
  }
})

function SearchGame() {
  searching.value = true
  searchedGames.value = []
  searchNameQuery.value = searchName.value
  fetch(`${env.VITE_SUPABASE_API}/functions/v1/search-games`, {
    method: 'POST',
    headers: {
      'Accept': 'application/json',
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      "name": searchName.value
    })
  }).then(response => {
    response.json().then(res => {
      searchedGames.value = res
    })
    searching.value = false
  })
  .catch(err => {
    console.error(err)
    searching.value = false
  })

  document.getElementById('search-input').blur()
}

onMounted(() => {
  if (searchNameQuery.value)
  {
    searchName.value = searchNameQuery.value
    SearchGame()
  }
})

</script>

<template>
  <header>
    <div>
      <input id="search-input" v-model="searchName" @keyup.enter="SearchGame()" placeholder="Game Name"></input>
      <button id="search" @click="SearchGame()"><font-awesome-icon icon="fa-solid fa-search" /></button>
    </div>
  </header>

  <main>
    <span class="loader" v-if="searching"></span>
    <div v-for="game in searchedGames">
      <GameCard :id="game.id" :name="game.name" :cover="game.cover"/>
    </div>
  </main>
</template>

<style scoped>
  #search {
    width: 40px;
    margin: 0px;
    padding: 0px;

    border-radius: 0px;
  }

  main {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
  }

  input {
    height: 40px;
    width: 300px;
    margin: auto 0px auto 0px;
    font-size: x-large;
    padding: 0px 10px;
  }

</style>