import { defineStore } from 'pinia'
import { authApi } from '../api'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('admin-token') || '',
    user: JSON.parse(localStorage.getItem('admin-user') || 'null')
  }),

  getters: {
    isLoggedIn: (state) => Boolean(state.token),
    isAdmin: (state) => state.user?.role === 'admin',
    isReviewer: (state) => ['admin', 'reviewer'].includes(state.user?.role)
  },

  actions: {
    async login(username, password) {
      const data = await authApi.login(username, password)
      this.token = data.token
      this.user = data.user
      localStorage.setItem('admin-token', data.token)
      localStorage.setItem('admin-user', JSON.stringify(data.user))
      return data.user
    },

    async fetchMe() {
      const user = await authApi.me()
      this.user = user
      localStorage.setItem('admin-user', JSON.stringify(user))
      return user
    },

    logout() {
      this.token = ''
      this.user = null
      localStorage.removeItem('admin-token')
      localStorage.removeItem('admin-user')
    }
  }
})
