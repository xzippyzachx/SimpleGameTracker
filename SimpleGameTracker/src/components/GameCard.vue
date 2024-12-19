<script setup>
defineProps({
  name: {
    type: String,
    required: true,
  },
})
</script>

<template>
  <div class="game-card">
    <div class="img"></div>
    <h3>{{ name }}</h3>
  </div>
</template>

<style scoped>
.game-card {
  flex-grow: 1 1 0;
  width: 200px;
  margin: 5px;
  border: solid;
  border-width: 2.5px;
  border-radius: 10px;
  text-align: center;
  overflow: hidden;
  cursor: pointer;
}

.game-card:hover {
  scale: 0.99;
}

.game-card .img {
  background-color: var(--vt-c-indigo);
  width: 100%;
  aspect-ratio: 1 / 1;
}

h3 {
  text-wrap: nowrap;
}
</style>
