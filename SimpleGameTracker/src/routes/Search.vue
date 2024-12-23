<script setup>
import { ref } from 'vue'

import GameCard from '../components/GameCard.vue'

const searchName = ref("")
const searchedGames = ref([])
function SearchGame()
{
  fetch('https://sdekcxxvsnnzypebfpcr.supabase.co/functions/v1/search-games', {
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
  })
  .catch(err => {
    console.error(err);
  });
}

</script>

<template>
  <header>
    <div>
      <button @click="$router.back()">Back</button>
    </div>

    <div>
      <input v-model="searchName" placeholder="Game Name"></input>
      <button @click="SearchGame()">Search</button>
    </div>
  </header>

  <main>
    <div v-for="game in searchedGames">
      <GameCard :id="game.id" :name="game.name" :cover="game.cover"/>
    </div>
  </main>
</template>

<style scoped>
  main {
    display: flex;
    flex-wrap: wrap;
    justify-content: center;
  }

  input {
    height: 25px;
    width: 300px;
    margin: auto 10px;
  }

</style>