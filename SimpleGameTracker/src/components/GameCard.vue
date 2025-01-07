<script setup>
defineProps({
  id: {
    type: Number,
    required: true,
  },
  name: {
    type: String,
    required: true,
  },
  cover: {
    type: String,
    required: true,
  },
  state: {
    type: String,
    required: false,
    default: 'none',
  },
})
</script>

<template>
  <RouterLink :to="{ name: 'game', params: { gameId: id } }">
    <div class="game-card" :class="{wishlist: state == 'wishlist', playing: state == 'playing', completed: state == 'completed'}">
      <div class="img" :style="{ backgroundImage: 'url(' + cover + ')'}"></div>
      <h3>{{ name }}</h3>
    </div>
  </RouterLink>
</template>

<style scoped>
  .game-card {
    flex-grow: 1 1 0;
    width: 220px;
    margin: 5px;
    border: solid;
    border-width: 2.5px;
    border-radius: 10px;
    text-align: center;
    cursor: pointer;
  }

  .game-card:hover {
    scale: 0.99;
  }

  .game-card .img {
    background-color: var(--vt-c-white);
    width: 100%;
    aspect-ratio: 0.75 / 1;
    background-size: 100%;
    background-repeat: no-repeat;
  }

  .game-card.wishlist {
    border-color: var(--vt-wishlist-c);
  }
  .game-card.playing {
    border-color: var(--vt-playing-c);
  }
  .game-card.completed {
    border-color: var(--vt-completed-c);
  }

  h3 {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    font-size: 100%;
    padding: 2.5px 5px;
  }
</style>
