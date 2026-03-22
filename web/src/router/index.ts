import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/login',
      name: 'Login',
      component: () => import('@/pages/Login.vue'),
      meta: { requiresAuth: false },
    },
    {
      path: '/',
      component: () => import('@/components/Layout/MainLayout.vue'),
      meta: { requiresAuth: true },
      children: [
        { path: '', name: 'Dashboard', component: () => import('@/pages/Dashboard.vue') },
        { path: 'proxies', name: 'Proxies', component: () => import('@/pages/Proxies.vue') },
        { path: 'rules', name: 'Rules', component: () => import('@/pages/Rules.vue') },
        { path: 'connections', name: 'Connections', component: () => import('@/pages/Connections.vue') },
        { path: 'logs', name: 'Logs', component: () => import('@/pages/Logs.vue') },
        { path: 'subscriptions', name: 'Subscriptions', component: () => import('@/pages/Subscriptions.vue') },
        { path: 'settings', name: 'Settings', component: () => import('@/pages/Settings.vue') },
      ],
    },
  ],
})

router.beforeEach(async (to, _from, next) => {
  if (to.meta.requiresAuth === false) {
    next()
    return
  }

  const authStore = useAuthStore()

  if (authStore.user) {
    next()
    return
  }

  if (authStore.checked) {
    next('/login')
    return
  }

  try {
    await authStore.fetchUser()
    next()
  } catch {
    authStore.setChecked()
    next('/login')
  }
})

export default router
