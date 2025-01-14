import './assets/main.css'

import { createApp } from 'vue'
import App from './App.vue'

import { library } from '@fortawesome/fontawesome-svg-core'
import { FontAwesomeIcon } from '@fortawesome/vue-fontawesome'
import { fas } from '@fortawesome/free-solid-svg-icons'
import { far } from '@fortawesome/free-regular-svg-icons'

library.add(fas)
library.add(far)

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
]

const router = createRouter({
  history: createWebHashHistory(),
  routes,
})

createApp(App)
  .component('font-awesome-icon', FontAwesomeIcon)
  .use(router)
  .mount('#app')
