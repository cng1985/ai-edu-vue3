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
      snapshots: saved?.snapshots || [],
      reviewRecords: saved?.reviewRecords || {}
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
      const levels = [
        { number: 1, name: '编程小白', min: 0, next: 100 },
        { number: 2, name: '代码萌新', min: 100, next: 500 },
        { number: 3, name: '技术学徒', min: 500, next: 2000 },
        { number: 4, name: '开发能手', min: 2000, next: 5000 },
        { number: 5, name: '技术达人', min: 5000, next: 10000 },
        { number: 6, name: '资深工程师', min: 10000, next: 30000 },
        { number: 7, name: '技术专家', min: 30000, next: 70000 },
        { number: 8, name: '技术布道师', min: 70000, next: 70000 }
      ]
      return [...levels].reverse().find((level) => state.points >= level.min)
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
    },
    reviewQueue(state) {
      const now = Date.now()
      return microUnits
        .filter((unit) => state.completedUnits[unit.id])
        .map((unit) => {
          const record = state.reviewRecords[unit.id]
          const fallbackDue = state.completedUnits[unit.id] + 86400000
          const dueAt = record?.dueAt || fallbackDue
          return {
            ...unit,
            score: state.unitScores[unit.id] || 0,
            stage: record?.stage || 0,
            dueAt,
            due: dueAt <= now
          }
        })
        .sort((a, b) => a.dueAt - b.dueAt)
    },
    dueReviewCount() {
      return this.reviewQueue.filter((unit) => unit.due).length
    },
    weakPoints(state) {
      return frontendPath.competencies.flatMap((domain) =>
        domain.points.map((point) => {
          const unit = microUnits.find((item) => item.id === point.unitId)
          const completed = Boolean(state.completedUnits[point.unitId])
          const score = state.unitScores[point.unitId] || 0
          return {
            ...point,
            unit,
            domainId: domain.id,
            domainName: domain.name,
            completed,
            score,
            mastery: completed ? Math.round(40 + score * 0.6) : 0,
            reason: !completed ? '尚未学习' : score < 70 ? '快测正确率偏低' : '需要间隔复习'
          }
        })
      ).filter((point) => !point.completed || point.score < 70)
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
        snapshots: this.snapshots.slice(-30),
        reviewRecords: this.reviewRecords
      }))
    },
    createGoal({ careerId = 'frontend', baseLevel, weeklyHours, durationWeeks }) {
      const weeks = Number(durationWeeks) || frontendPath.durationWeeks
      const deadline = new Date()
      deadline.setDate(deadline.getDate() + weeks * 7)
      this.goal = {
        id: `goal-${Date.now().toString(36)}`,
        careerId,
        name: `${weeks} 周成为初级 Web 前端工程师`,
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
      if (firstCompletion) {
        this.award(5)
        this.reviewRecords[id] = {
          stage: 0,
          dueAt: Date.now() + 86400000,
          lastAt: null
        }
      }
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
    completeReview(id) {
      if (!this.completedUnits[id]) return false
      const intervals = [1, 3, 7, 15, 30]
      const current = this.reviewRecords[id] || { stage: 0, dueAt: Date.now(), lastAt: null }
      if (current.lastAt && dateKey(new Date(current.lastAt)) === dateKey()) return false
      const nextStage = Math.min(current.stage + 1, intervals.length - 1)
      this.reviewRecords[id] = {
        stage: nextStage,
        lastAt: Date.now(),
        dueAt: Date.now() + intervals[nextStage] * 86400000
      }
      this.award(3)
      if (Object.values(this.reviewRecords).filter((item) => item.lastAt).length >= 3) {
        this.addBadge('复习达人')
      }
      this.persist()
      return true
    },
    reset() {
      this.goal = null
      this.completedUnits = {}
      this.unitScores = {}
      this.points = 0
      this.badges = []
      this.checkIns = []
      this.snapshots = []
      this.reviewRecords = {}
      this.persist()
    }
  }
})
