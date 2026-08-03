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

function withGuestFlag(user) {
  return { ...user, isGuest: user?.role === 'guest' }
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: loadSession()
  }),

  getters: {
    isLoggedIn: (state) => Boolean(state.user),
    isGuest: (state) => Boolean(state.user?.isGuest || state.user?.role === 'guest'),
    permissions: (state) => state.user?.permissions || [],
    hasPermission: (state) => (perm) => {
      if (!perm) return true
      return state.user?.permissions?.includes(perm)
    },
    hasAnyPermission: (state) => (list) => {
      if (!list?.length) return true
      return list.some((p) => state.user?.permissions?.includes(p))
    }
  },

  actions: {
    _applyUser(user, token) {
      const enriched = withGuestFlag(user)
      if (token) saveSession(token, enriched)
      else {
        localStorage.setItem('ai-learning-system:session', JSON.stringify(enriched))
      }
      this.user = enriched
      return enriched
    },

    async register({ username, nickname, password, avatar, avatarColor }) {
      const data = await authApi.register({ username, nickname, password, avatar, avatarColor })
      return this._applyUser(data.user, data.token)
    },

    async login(username, password) {
      const data = await authApi.login(username, password)
      return this._applyUser(data.user, data.token)
    },

    async loginAsGuest() {
      const data = await authApi.guest()
      return this._applyUser(data.user, data.token)
    },

    async fetchMe() {
      const user = await authApi.me()
      return this._applyUser(user)
    },

    async refreshPermissions() {
      const data = await authApi.permissions()
      if (this.user) {
        return this._applyUser({
          ...this.user,
          role: data.role,
          roleName: data.roleName,
          permissions: data.permissions
        })
      }
      return data.permissions
    },

    logout() {
      this.user = null
      clearSession()
    }
  }
})
