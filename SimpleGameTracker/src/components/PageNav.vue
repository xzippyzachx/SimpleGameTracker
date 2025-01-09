
<template>
  <nav>
    <RouterLink :to="{ name: 'library' }"><font-awesome-icon icon="fa-solid fa-gamepad" /></RouterLink>
    <RouterLink :to="{ name: 'search' }"><font-awesome-icon icon="fa-solid fa-plus" /></RouterLink>
    <RouterLink :to="{ name: 'library' }"><font-awesome-icon icon="fa-solid fa-bars" /></RouterLink>
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
  }

  nav a {
    width: 70px;
    padding: 5px;
    display: flex;
    justify-content: center;
    flex-direction: column;
  }

  nav a:hover {
    color: var(--vt-c-blue);
  }

  nav svg {
    height: 40px;
  }

</style>
