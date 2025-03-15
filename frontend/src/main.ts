import { createApp } from 'vue'
import { createRouter, createWebHistory } from 'vue-router'
import App from './App.vue'
import { MotionPlugin } from '@vueuse/motion'
import { routes } from './routes'
import './index.css'

const app = createApp(App)

const router = createRouter({
  history: createWebHistory(),
  routes: import.meta.hot ? [] : routes,
})

if (import.meta.hot) {
  let removeRoutes: any[] = []

  for (const route of routes) {
    removeRoutes.push(router.addRoute(route))
  }
  import.meta.hot?.accept('./routes.ts', ({ routes }) => {
    for (const removeRoute of removeRoutes) removeRoute()
    removeRoutes = []
    for (const route of routes) {
      removeRoutes.push(router.addRoute(route))
    }
    router.replace('')
  })
}

// Guard global: si la ruta no es "Login" y no hay token, redirige a Login
router.beforeEach((to, from, next) => {
  const token = localStorage.getItem('auth_token')
  // Redirige solo si la ruta no es 'Login' y no hay token
  if (to.name !== 'Login' && !token) {
    next({ name: 'Login' })
  } else {
    next()
  }
})


app.use(MotionPlugin)
app.use(router)
app.mount('#app')
