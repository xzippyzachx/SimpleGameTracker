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
})
</script>

<template>
  <div class="game-card">
    <div class="img" :style="{ backgroundImage: 'url(' + cover + ')'}"></div>
    <h3>{{ name }}</h3>
  </div>
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

h3 {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  font-size: 100%;
  padding: 2.5px 5px;
}
</style>
