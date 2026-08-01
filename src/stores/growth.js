import { defineStore } from 'pinia'
import { frontendPath, microUnits } from '../data/careerPath'

const STORAGE_KEY = 'ai-learning-system:growth'

function loadState() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    return raw ? JSON.parse(raw) : null
  } catch {
    return null
  }
}

const dateKey = (date = new Date()) => {
  const local = new Date(date.getTime() - date.getTimezoneOffset() * 60000)
  return local.toISOString().slice(0, 10)
}

export const useGrowthStore = defineStore('growth', {
  state: () => {
    const saved = loadState()
    return {
      goal: saved?.goal || null,
      completedUnits: saved?.completedUnits || {},
      unitScores: saved?.unitScores || {},
      points: saved?.points || 0,
      badges: saved?.badges || [],
      checkIns: saved?.checkIns || [],
      snapshots: saved?.snapshots || []
    }
  },

  getters: {
    hasGoal: (state) => Boolean(state.goal),
    completedUnitCount: (state) => Object.keys(state.completedUnits).length,
    isUnitCompleted: (state) => (id) => Boolean(state.completedUnits[id]),
    nextUnit(state) {
      return microUnits.find((unit) => !state.completedUnits[unit.id]) || null
    },
    streak(state) {
      const dates = new Set(state.checkIns)
      let count = 0
      const cursor = new Date()
      if (!dates.has(dateKey(cursor))) cursor.setDate(cursor.getDate() - 1)
      while (dates.has(dateKey(cursor))) {
        count += 1
        cursor.setDate(cursor.getDate() - 1)
      }
      return count
    },
    checkedInToday: (state) => state.checkIns.includes(dateKey()),
    level(state) {
      if (state.points >= 500) return { number: 3, name: '技术学徒', next: 2000 }
      if (state.points >= 100) return { number: 2, name: '代码萌新', next: 500 }
      return { number: 1, name: '编程小白', next: 100 }
    },
    competencyProgress: (state) => {
      if (!state.goal) return []
      return frontendPath.competencies.map((competency) => {
        const progress = competency.points.reduce((sum, point) => {
          const completed = state.completedUnits[point.unitId] ? 1 : 0
          const score = state.unitScores[point.unitId] ?? 0
          // MVP 掌握度：完成行为 40% + 快测成绩 60%
          const mastery = completed * 40 + score * 0.6
          return sum + mastery * (point.weight / 100)
        }, 0)
        return { ...competency, progress: Math.round(progress) }
      })
    },
    achievement() {
      if (!this.hasGoal) return 0
      return Math.round(this.competencyProgress.reduce(
        (sum, competency) => sum + competency.progress * (competency.weight / 100),
        0
      ))
    },
    nextMilestone() {
      if (!this.hasGoal) return null
      const completed = this.completedUnitCount
      const index = completed < 2 ? 0 : completed < 4 ? 1 : completed < 6 ? 2 : 3
      return frontendPath.milestones[index]
    }
  },

  actions: {
    persist() {
      localStorage.setItem(STORAGE_KEY, JSON.stringify({
        goal: this.goal,
        completedUnits: this.completedUnits,
        unitScores: this.unitScores,
        points: this.points,
        badges: this.badges,
        checkIns: this.checkIns,
        snapshots: this.snapshots.slice(-30)
      }))
    },
    createGoal({ careerId = 'frontend', baseLevel, weeklyHours, durationWeeks }) {
      const weeks = Number(durationWeeks) || frontendPath.durationWeeks
      const deadline = new Date()
      deadline.setDate(deadline.getDate() + weeks * 7)
      this.goal = {
        id: `goal-${Date.now().toString(36)}`,
        careerId,
        name: frontendPath.name,
        description: frontendPath.description,
        baseLevel,
        weeklyHours: Number(weeklyHours),
        durationWeeks: weeks,
        deadline: deadline.toISOString().slice(0, 10),
        status: '进行中',
        createdAt: Date.now()
      }
      this.award(20, '目标启航')
      this.persist()
    },
    completeUnit(id, score) {
      const firstCompletion = !this.completedUnits[id]
      this.completedUnits[id] = Date.now()
      this.unitScores[id] = Math.max(this.unitScores[id] || 0, score)
      if (firstCompletion) this.award(5)
      if (Object.keys(this.completedUnits).length === 1) this.addBadge('第一步')
      this.addSnapshot()
      this.persist()
    },
    checkIn() {
      const today = dateKey()
      if (this.checkIns.includes(today)) return false
      this.checkIns.push(today)
      this.award(10)
      if (this.streak >= 7) this.addBadge('坚持者')
      this.persist()
      return true
    },
    award(amount, badge) {
      this.points += amount
      if (badge) this.addBadge(badge)
    },
    addBadge(name) {
      if (!this.badges.includes(name)) this.badges.push(name)
    },
    addSnapshot() {
      this.snapshots.push({ date: dateKey(), value: this.achievement })
    },
    reset() {
      this.goal = null
      this.completedUnits = {}
      this.unitScores = {}
      this.points = 0
      this.badges = []
      this.checkIns = []
      this.snapshots = []
      this.persist()
    }
  }
})
