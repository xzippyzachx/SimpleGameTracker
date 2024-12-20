<script setup>
import { ref, computed } from 'vue'
import Home from './routes/Home.vue'
import Search from './routes/Search.vue'

const routes = {
  '/': Home,
  '/search': Search
}

const currentPath = ref(window.location.hash)

window.addEventListener('hashchange', () => {
  currentPath.value = window.location.hash
})

const currentView = computed(() => {
  return routes[currentPath.value.slice(1) || '/'] || NotFound
})

</script>

<template>
  <component :is="currentView" />
</template>
