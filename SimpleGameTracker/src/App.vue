<script setup>
import GameCard from './components/GameCard.vue'

import { ref } from 'vue'

const state = ref('playing')

</script>

<template>
  <header>
    <nav>
      <button id="wishlist-btn" :class="{ selected: state == 'wishlist' }" @click="state = 'wishlist'">Wishlist</button>
      <button id="playing-btn" :class="{ selected: state == 'playing' }" @click="state = 'playing'">Playing</button>
      <button id="completed-btn" :class="{ selected: state == 'completed' }" @click="state = 'completed'">Completed</button>
    </nav>
  </header>

  <main>
    <GameCard name="The Finals"/>
    <GameCard name="Helldivers 2"/>
    <GameCard name="Battlefield 2042"/>
  </main>
</template>

<style scoped>

  #wishlist-btn {
    background-color: rgb(228, 71, 202);
  }
  #playing-btn {
    background-color: rgb(79, 186, 248);
  }
  #completed-btn {
    background-color: rgb(58, 223, 71);
  }

  main {
    display: flex;
    justify-content: center;
  }

</style>