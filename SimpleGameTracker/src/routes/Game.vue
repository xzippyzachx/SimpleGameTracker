<script setup>
import { ref } from 'vue'
import { useRoute } from 'vue-router'

import StateNav from '../components/StateNav.vue'

const gameData = ref({})
const route = useRoute()

function GameInfo()
{
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

GameInfo()

</script>

<template>
  <header>
    <div>
      <button @click="$router.back()">Back</button>
    </div>
    <StateNav/>
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