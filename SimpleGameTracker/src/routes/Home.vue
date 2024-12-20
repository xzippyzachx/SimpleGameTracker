<script setup>
import GameCard from '../components/GameCard.vue'
import StateNav from '../components/StateNav.vue'
</script>

<template>
  <header>
    <div>
      <a href="#/search"><button id="add-game-btn">Add Game</button></a>
    </div>
    <StateNav/>
  </header>

  <main>
    <GameCard name="The Finals"/>
    <GameCard name="Helldivers 2"/>
    <GameCard name="Battlefield 2042"/>
  </main>
</template>

<style scoped>
  main {
    display: flex;
    justify-content: center;
  }

  #add-game-btn {
    background-color: var(--vt-c-white);
  }

</style>