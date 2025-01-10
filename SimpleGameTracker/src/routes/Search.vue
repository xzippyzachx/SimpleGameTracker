<script setup>
import { ref } from 'vue'

import GameCard from '../components/GameCard.vue'

const searchName = ref("")
const searchedGames = ref([])

function SearchGame() {
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
      <input v-model="searchName" @keyup.enter="SearchGame()" placeholder="Game Name"></input>
      <button @click="SearchGame()"><font-awesome-icon icon="fa-solid fa-search" /></button>
    </div>
  </header>

  <main>
    <div v-for="game in searchedGames">
      <GameCard :id="game.id" :name="game.name" :cover="game.cover"/>
    </div>
  </main>
</template>

<style scoped>
  header button {
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