import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import HomeView from './views/HomeView.vue'
import CalendarView from './views/CalendarView.vue'
import BookshelfView from './views/BookshelfView.vue'
import i18n from './i18n'

import './style.css' // We might need to create this for basic styles

const router = createRouter({
  history: createWebHistory(),
  routes: [
    { path: '/', component: HomeView },
    { path: '/calendar', component: CalendarView },
    { path: '/bookshelf', component: BookshelfView },
  ]
})

const app = createApp(App)
app.use(router)
app.use(i18n)
app.mount('#app')
