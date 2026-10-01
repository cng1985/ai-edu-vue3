import { defineStore } from 'pinia'
import { meApi } from '../api'

/** 个人成长总览（服务端计算），供侧栏与成长中心共享 */
export const useGrowthStore = defineStore('growth', {
  state: () => ({
    overview: null,
    loading: false,
    error: ''
  }),

  getters: {
    goal: (state) => state.overview?.goal || null,
    hasGoal: (state) => Boolean(state.overview?.goal),
    gap: (state) => state.overview?.gap || null,
    readiness: (state) => state.overview?.gap?.readiness || 0,
    roleName: (state) => state.overview?.gap?.role?.name || ''
  },

  actions: {
    async refresh() {
      this.loading = true
      this.error = ''
      try {
        this.overview = await meApi.overview()
      } catch (e) {
        this.error = e.message
      } finally {
        this.loading = false
      }
      return this.overview
    },

    reset() {
      this.overview = null
    }
  }
})
