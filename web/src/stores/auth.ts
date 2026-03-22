import { defineStore } from 'pinia'
import { ref } from 'vue'
import { authApi } from '@/api/auth'
import type { User } from '@/types/api'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const initialized = ref(false)
  const checked = ref(false)

  const checkInit = async () => {
    const res = await authApi.checkInit()
    initialized.value = res.data.data.initialized
    return initialized.value
  }

  const initPassword = async (username: string, password: string) => {
    await authApi.initPassword(username, password)
    initialized.value = true
  }

  const login = async (username: string, password: string) => {
    const res = await authApi.login(username, password)
    user.value = res.data.data.user
    checked.value = true
    return user.value
  }

  const logout = async () => {
    await authApi.logout()
    user.value = null
    checked.value = false
  }

  const fetchUser = async () => {
    const res = await authApi.getCurrentUser()
    user.value = res.data.data.user
    checked.value = true
    return user.value
  }

  const setChecked = () => {
    checked.value = true
  }

  return {
    user,
    initialized,
    checked,
    checkInit,
    initPassword,
    login,
    logout,
    fetchUser,
    setChecked,
  }
})
