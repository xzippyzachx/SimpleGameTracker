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
      <h3><span :class="{wishlist: state == 'wishlist', playing: state == 'playing', completed: state == 'completed'}">{{ name }}</span></h3>
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

  .wishlist {
    border-color: var(--vt-wishlist-c) !important;
  }
  .playing {
    border-color: var(--vt-playing-c) !important;
  }
  .completed {
    border-color: var(--vt-completed-c) !important;
  }

  h3 {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    padding: 5px 10px;
  }

  h3:hover {
    overflow: visible;
    padding: 5px 0px;
  }

  h3:hover span {
    background-color: var(--color-background-mute);
    padding: 8.75px 10px;
    border: solid;
    border-width: 2.5px;
    border-radius: 10px;
    margin-left: -100%;
    margin-right: -100%;
    text-align: center;
  }
</style>
