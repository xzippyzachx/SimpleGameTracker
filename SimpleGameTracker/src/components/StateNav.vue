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
    <button id="shelved-btn" :class="{ selected: state == 'shelved' }" @click="ChangeState('shelved')">Shelved</button>
  </nav>
</template>

<style scoped>
  nav {
    display: flex;
    justify-content: center;
    flex-wrap: wrap;
  }

  button.selected {
    color: var(--color-btn);
    border-style: solid;
    border-width: 2.5px;
    margin-top: 0px;
    margin-bottom: 0px;
  }

  button {
    flex-grow: 1;
    width: inherit;
    max-width: 150px;
  }

  @media only screen and (max-width: 550px) {
    button {
      margin: 2.5px 2.5px;
      padding: 10px 5px;
    }
  }

  #wishlist-btn {
    background-color: var(--vt-wishlist-c);
  }
  #playing-btn {
    background-color: var(--vt-playing-c);
  }
  #completed-btn {
    background-color: var(--vt-completed-c);
  }
  #shelved-btn {
    background-color: var(--vt-shelved-c);
  }
</style>
