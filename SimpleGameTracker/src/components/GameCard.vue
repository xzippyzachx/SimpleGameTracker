<script setup>
import { computed, ref } from 'vue'

import { ChangeGameState } from '../gameDataCache.js'

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
  showStateSelect: {
    type: Boolean,
    required: false,
    default: false,
  }
})

const localState = ref('')
const state = computed(() => localState.value != '' ? localState.value : props.state)

const stateSelect = ref(false)
const stateSelectRef = ref(null);

const stateClass = computed(() => ({
  wishlist: state.value == 'wishlist',
  playing: state.value == 'playing',
  completed: state.value == 'completed',
  shelved: state.value == 'shelved',
}))

function ChangeState(newState) {
  if (state.value == newState) {
    newState = 'none'
  }
  localState.value = newState
  ChangeGameState(props.id, newState)

  stateSelect.value = false
}

function OpenStateSelect() {
  stateSelect.value = true
  stateSelectRef.value.focus()
}

</script>

<template>
  <div class="game-card" :class="stateClass" @mouseleave="stateSelect = false">
    <div class="state-select" v-if="showStateSelect" tabindex="0" @blur="stateSelect = false" ref="stateSelectRef">
      <button v-if="!stateSelect" class="state-select-bars" :class="stateClass" @click="OpenStateSelect()"><font-awesome-icon icon="fa-regular fa-bookmark"/></button>

      <button v-if="stateSelect" class="wishlist" @mousedown="ChangeState('wishlist')" title="Wishlist"><font-awesome-icon icon="fa-solid fa-minus"/></button>
      <button v-if="stateSelect" class="playing" @mousedown="ChangeState('playing')" title="Playing"><font-awesome-icon icon="fa-solid fa-minus"/></button>
      <button v-if="stateSelect" class="completed" @mousedown="ChangeState('completed')" title="Completed"><font-awesome-icon icon="fa-solid fa-minus"/></button>
      <button v-if="stateSelect" class="shelved" @mousedown="ChangeState('shelved')" title="Shelved"><font-awesome-icon icon="fa-solid fa-minus"/></button>
    </div>
    
    <RouterLink :to="{ name: 'game', params: { gameId: id } }">
      <div class="img" :style="{ backgroundImage: 'url(' + cover + ')'}"></div>
      <h4><span :class="stateClass">{{ name }}</span></h4>
    </RouterLink>
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

  @media only screen and (max-width: 500px) {
    .game-card {
      width: 170px;
    }

    .game-card .state-select {
      display: flex;
      max-width: 165px;
    }

    .game-card:hover .state-select {
      max-width: 165px !important;
    }
  }

  .game-card:hover {
    scale: 0.99;

    .state-select {
      display: flex;
      max-width: 215px;
    }
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
  .game-card.shelved {
    border-color: var(--vt-shelved-c) !important;
  }
  .shelved {
    background-color: var(--vt-shelved-c) !important;
  }


  h4 {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    padding: 5px 10px;
    text-shadow: 1px 1px 10px var(--vt-c-black);
  }

  h4:hover {
    overflow: visible;
    padding: 5px 0px;
  }

  h4:hover span {
    background-color: var(--color-background-mute);
    padding: 8.75px 10px;
    border: solid;
    border-width: 2.5px;
    border-radius: 10px;
    margin-left: -100%;
    margin-right: -100%;
    text-align: center;
  }

  .state-select-bars {
    max-width: 30px;
  }

  .state-select {
    display: none;
    position: absolute;
    height: 30px;
    padding: 5px;
    gap: 3px;

    button {
      margin: 0px !important;
      padding: 0px;
      box-shadow: 1px 1px 10px var(--vt-c-black);
    }
  }

  @media only screen and (max-width: 500px) {
    .state-select {
      height: 35px;
    }
  }
</style>
