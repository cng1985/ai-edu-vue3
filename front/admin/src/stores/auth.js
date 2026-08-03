import { defineStore } from 'pinia'
import { authApi } from '../api'

function loadUser() {
  try {
    return JSON.parse(localStorage.getItem('admin-user') || 'null')
  } catch {
    return null
  }
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    token: localStorage.getItem('admin-token') || '',
    user: loadUser(),
    permissions: loadUser()?.permissions || []
  }),

  getters: {
    isLoggedIn: (state) => Boolean(state.token),
    isAdminPortal: (state) => ['admin', 'reviewer', 'operator'].includes(state.user?.role),
    hasPermission: (state) => (perm) => {
      if (!perm) return true
      const perms = state.permissions.length ? state.permissions : state.user?.permissions
      return Array.isArray(perms) && perms.includes(perm)
    },
    hasAnyPermission: (state) => (list) => {
      if (!list?.length) return true
      const perms = state.permissions.length ? state.permissions : state.user?.permissions
      return list.some((p) => perms?.includes(p))
    },
    /** @deprecated 请使用 hasPermission */
    isAdmin: (state) => state.user?.role === 'admin',
    /** @deprecated 请使用 hasPermission */
    isReviewer: (state) => ['admin', 'reviewer'].includes(state.user?.role)
  },

  actions: {
    _persist(user) {
      this.user = user
      this.permissions = user?.permissions || []
      localStorage.setItem('admin-user', JSON.stringify(user))
      localStorage.setItem('admin-permissions', JSON.stringify(this.permissions))
    },

    async login(username, password) {
      const data = await authApi.login(username, password, 'admin')
      this.token = data.token
      localStorage.setItem('admin-token', data.token)
      this._persist(data.user)
      return data.user
    },

    async fetchMe() {
      const user = await authApi.me()
      this._persist(user)
      return user
    },

    async refreshPermissions() {
      const user = await authApi.refreshPermissions()
      this._persist(user)
      return user
    },

    logout() {
      this.token = ''
      this.user = null
      this.permissions = []
      localStorage.removeItem('admin-token')
      localStorage.removeItem('admin-user')
      localStorage.removeItem('admin-permissions')
    }
  }
})
