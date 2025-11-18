import { createRouter, createWebHistory } from 'vue-router'
import Home from '../App.vue'
import Gallery from '../components/Gallery.vue'

const routes = [
  { path: '/', name: 'Home', component: Home },
  { path: '/gallery', name: 'Gallery', component: Gallery },
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
