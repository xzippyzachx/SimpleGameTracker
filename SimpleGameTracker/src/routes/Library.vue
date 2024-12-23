<script setup>
import GameCard from '../components/GameCard.vue'
import StateNav from '../components/StateNav.vue'
</script>

<template>
  <header>
    <div>
      <RouterLink :to="{ name: 'search' }">
        <button id="add-game-btn">Add Game</button>
      </RouterLink>
    </div>
    <StateNav state="playing"/>
  </header>

  <main>
    <GameCard :id="1" name="The Finals" cover="/"/>
    <GameCard :id="2" name="Helldivers 2" cover="/"/>
    <GameCard :id="3" name="Battlefield 2042" cover="/"/>
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