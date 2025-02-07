<script>
export default {
  methods: {
    Shared() { return sessionStorage.gameDataShared != undefined || this.$route.name == 'shared' }
  }
}
</script>

<template>
  <nav id="PageNav" v-if="!Shared()">
    <RouterLink :to="{ name: 'settings' }" :class="{ selected: $route.name == 'settings' }"><font-awesome-icon icon="fa-solid fa-bars" /></RouterLink>
    <RouterLink :to="{ name: 'library' }" :class="{ selected: $route.name == 'library' }"><font-awesome-icon icon="fa-solid fa-gamepad" /></RouterLink>
    <RouterLink :to="{ name: 'search' }" :class="{ selected: $route.name == 'search' }"><font-awesome-icon icon="fa-regular fa-square-plus" /></RouterLink>
  </nav>

  <nav id="PageNav" v-if="Shared()">
    <RouterLink :to="{ name: 'library' }" target="_blank" class="try"><font-awesome-icon icon="fa-solid fa-gamepad" /><h1>Try Simple Game Tracker</h1></RouterLink>
  </nav>
</template>

<style scoped>
  nav {
    display: flex;
    justify-content: center;
    flex-wrap: wrap;
    gap: 30px;
    height: 70px;
    background-color: var(--color-background-soft);
    border-bottom: solid;
    border-color: var(--vt-c-divider-light-1);
  }

  nav a {
    width: 70px;
    padding: 5px;
    display: flex;
    justify-content: center;
    flex-direction: column;
  }

  .try {
    width: auto;
    padding: 5px;
    gap: 10px;
    display: flex;
    flex-direction: row;
    margin: auto;
    cursor: pointer;
    font-size: 17px;
    text-wrap: nowrap;
  }

  .try svg {
    width: 40px;
  }

  .try:hover {
    color: var(--vt-c-blue);
  }

  nav a.selected{
    color: var(--vt-c-blue);
  }

  nav svg {
    height: 40px;
  }

  @media only screen and (max-width: 550px) {
    nav {
      border-bottom: none;
      border-top: solid;
      border-color: var(--vt-c-divider-light-1);
    }
  }

</style>
