<script setup>
import { computed } from 'vue'

const props = defineProps({
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

const stateClass = computed(() => ({
  wishlist: props.state == 'wishlist',
  playing: props.state == 'playing',
  completed: props.state == 'completed',
}))

</script>

<template>
  <RouterLink :to="{ name: 'game', params: { gameId: id } }">
    <div class="game-card" :class="stateClass">
      <div class="img" :style="{ backgroundImage: 'url(' + cover + ')'}"></div>
      <h3><span :class="stateClass">{{ name }}</span></h3>
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
    border-radius: 8px 8px 0px 0px;
    width: 100%;
    aspect-ratio: 0.75 / 1;
    background-size: 100%;
    background-repeat: no-repeat;
  }

  .game-card.wishlist {
    border-color: var(--vt-wishlist-c) !important;
  }
  .wishlist {
    background-color: var(--vt-wishlist-c) !important;
  }
  .game-card.playing {
    border-color: var(--vt-playing-c) !important;
  }
  .playing {
    background-color: var(--vt-playing-c) !important;
  }
  .game-card.completed {
    border-color: var(--vt-completed-c) !important;
  }
  .completed {
    background-color: var(--vt-completed-c) !important;
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
