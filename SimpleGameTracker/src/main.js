import './assets/main.css'
import './gameDataSync.js'

import { createApp } from 'vue'
import App from './App.vue'

import { library } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { fas } from '@fortawesome/free-solid-svg-icons'
import { far } from '@fortawesome/free-regular-svg-icons'
import { fab } from '@fortawesome/free-brands-svg-icons'

library.add(fas)
library.add(far)
library.add(fab)

import { createWebHashHistory, createRouter } from 'vue-router'
import Library from './routes/Library.vue'
import Search from './routes/Search.vue'
import Game from './routes/Game.vue'
import Settings from './routes/Settings.vue'

const routes = [
  { path: '/', name: 'library',  component: Library },
  { path: '/search', name: 'search', component: Search },
  { path: '/settings', name: 'settings', component: Settings },
  { path: '/game:gameId', name: 'game', component: Game },
  { path: '/shared:syncHash', name: 'shared', component: Library },
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

function DetectTheme() {
  let theme = 'dark'
  if (localStorage.getItem('theme')) {
    theme = localStorage.getItem('theme')
  } 
  else if (window.matchMedia && window.matchMedia('(prefers-color-scheme: light)').matches) {
    theme = 'light'
  }
  localStorage.setItem('theme', theme)
  document.documentElement.setAttribute('data-theme', theme)
}
DetectTheme()

createApp(App)
  .component('font-awesome-icon', FontAwesomeIcon)
  .use(router)
  .mount('#app')
