import { defineStore } from 'pinia'
import { authApi, saveSession, clearSession, loadSession } from '../api'

/** 社区字母头像配色 */
export const AVATAR_PRESETS = [
  { char: 'A', color: '#1772f6' },
  { char: 'B', color: '#0ea5e9' },
  { char: 'C', color: '#10b981' },
  { char: 'D', color: '#f59e0b' },
  { char: 'E', color: '#ef4444' },
  { char: 'F', color: '#8b5cf6' },
  { char: 'G', color: '#14b8a6' },
  { char: 'H', color: '#f97316' }
]

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: loadSession()
  }),

  getters: {
    isLoggedIn: (state) => Boolean(state.user),
    isGuest: (state) => Boolean(state.user?.isGuest || state.user?.role === 'guest')
  },

  actions: {
    async register({ username, nickname, password, avatar, avatarColor }) {
      const data = await authApi.register({ username, nickname, password, avatar, avatarColor })
      saveSession(data.token, data.user)
      this.user = { ...data.user, isGuest: false }
      return this.user
    },

    async login(username, password) {
      const data = await authApi.login(username, password)
      saveSession(data.token, data.user)
      this.user = { ...data.user, isGuest: false }
      return this.user
    },

    async loginAsGuest() {
      const data = await authApi.guest()
      saveSession(data.token, data.user)
      this.user = { ...data.user, isGuest: true }
      return this.user
    },

    async fetchMe() {
      const user = await authApi.me()
      this.user = { ...user, isGuest: user.role === 'guest' }
      localStorage.setItem('ai-learning-system:session', JSON.stringify(this.user))
      return this.user
    },

    logout() {
      this.user = null
      clearSession()
    }
  }
})
