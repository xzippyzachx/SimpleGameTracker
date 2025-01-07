<script setup>
const props = defineProps({
  state: {
    type: String,
    required: false,
    default: 'none',
  },
})

import { ref } from 'vue'

let state = ref(props.state)

function ChangeState(newState) {
  if (state.value == newState) {
    state.value = 'none'
  }
  else {
    state.value = newState
  }
}

defineExpose({
  state,
})

</script>

<template>
  <nav>
    <button id="wishlist-btn" :class="{ selected: state == 'wishlist' }" @click="ChangeState('wishlist')">Wishlist</button>
    <button id="playing-btn" :class="{ selected: state == 'playing' }" @click="ChangeState('playing')">Playing</button>
    <button id="completed-btn" :class="{ selected: state == 'completed' }" @click="ChangeState('completed')">Completed</button>
  </nav>
</template>

<style scoped>
  nav {
    display: flex;
    justify-content: center;
    flex-wrap: wrap;
    margin-bottom: 20px;
  }

  button.selected {
    color: var(--vt-c-white);
    border-style: solid;
    border-width: 2.5px;
  }

  #wishlist-btn {
    background-color: rgb(228, 71, 202);
  }
  #playing-btn {
    background-color: rgb(79, 186, 248);
  }
  #completed-btn {
    background-color: rgb(58, 223, 71);
  }
</style>
