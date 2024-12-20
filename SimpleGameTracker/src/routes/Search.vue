<script setup>
import GameCard from '../components/GameCard.vue'
</script>

<template>
  <header>
    <div>
      <a href="#/"><button>Library</button></a>
    </div>
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

</style>