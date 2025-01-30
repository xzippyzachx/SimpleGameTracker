
<template>
  <nav>
    <RouterLink :to="{ name: 'settings' }" :class="{ selected: $route.name == 'settings' }"><font-awesome-icon icon="fa-solid fa-bars" /></RouterLink>
    <RouterLink :to="{ name: 'library' }" :class="{ selected: $route.name == 'library' }"><font-awesome-icon icon="fa-solid fa-gamepad" /></RouterLink>
    <RouterLink :to="{ name: 'search' }" :class="{ selected: $route.name == 'search' }"><font-awesome-icon icon="fa-regular fa-square-plus" /></RouterLink>
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
