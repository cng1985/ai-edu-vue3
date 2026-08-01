import { defineStore } from 'pinia'
import { courses, totalChapterCount } from '../data/courses'

const STORAGE_KEY = 'ai-learning-system:learning'

function loadState() {
  try {
    const raw = localStorage.getItem(STORAGE_KEY)
    if (raw) return JSON.parse(raw)
  } catch (e) { /* 本地数据损坏时回退到初始状态 */ }
  return null
}

export const useLearningStore = defineStore('learning', {
  state: () => {
    const saved = loadState()
    return {
      // completed: { 'courseId/chapterId': timestamp }
      completed: saved?.completed || {},
      // notes: { 'courseId/chapterId': '笔记内容' }
      notes: saved?.notes || {},
      // quizResults: { quizId: { score, total, at, answers } }
      quizResults: saved?.quizResults || {},
      // lastVisited: { courseId, chapterId } 用于"继续学习"
      lastVisited: saved?.lastVisited || null,
      // studyLog: { 'YYYY-MM-DD': 完成章节数 } 用于活跃度统计
      studyLog: saved?.studyLog || {}
    }
  },

  getters: {
    completedCount: (state) => Object.keys(state.completed).length,

    overallProgress() {
      return totalChapterCount === 0
        ? 0
        : Math.round((this.completedCount / totalChapterCount) * 100)
    },

    isChapterCompleted: (state) => (courseId, chapterId) =>
      Boolean(state.completed[`${courseId}/${chapterId}`]),

    courseProgress: (state) => (courseId) => {
      const course = courses.find((c) => c.id === courseId)
      if (!course) return 0
      const done = course.chapters.filter(
        (ch) => state.completed[`${courseId}/${ch.id}`]
      ).length
      return Math.round((done / course.chapters.length) * 100)
    },

    courseCompletedCount: (state) => (courseId) => {
      const course = courses.find((c) => c.id === courseId)
      if (!course) return 0
      return course.chapters.filter((ch) => state.completed[`${courseId}/${ch.id}`]).length
    },

    noteFor: (state) => (courseId, chapterId) =>
      state.notes[`${courseId}/${chapterId}`] || '',

    noteCount: (state) =>
      Object.values(state.notes).filter((n) => n && n.trim()).length,

    quizAverageScore: (state) => {
      const results = Object.values(state.quizResults)
      if (results.length === 0) return null
      const sum = results.reduce((acc, r) => acc + (r.score / r.total) * 100, 0)
      return Math.round(sum / results.length)
    },

    // 最近 7 天的学习活跃度，用于统计页图表
    weeklyActivity: (state) => {
      const days = []
      for (let i = 6; i >= 0; i--) {
        const d = new Date()
        d.setDate(d.getDate() - i)
        const key = d.toISOString().slice(0, 10)
        days.push({
          date: key,
          label: `${d.getMonth() + 1}/${d.getDate()}`,
          count: state.studyLog[key] || 0
        })
      }
      return days
    }
  },

  actions: {
    persist() {
      localStorage.setItem(
        STORAGE_KEY,
        JSON.stringify({
          completed: this.completed,
          notes: this.notes,
          quizResults: this.quizResults,
          lastVisited: this.lastVisited,
          studyLog: this.studyLog
        })
      )
    },

    toggleChapterCompleted(courseId, chapterId) {
      const key = `${courseId}/${chapterId}`
      if (this.completed[key]) {
        delete this.completed[key]
      } else {
        this.completed[key] = Date.now()
        const day = new Date().toISOString().slice(0, 10)
        this.studyLog[day] = (this.studyLog[day] || 0) + 1
      }
      this.persist()
    },

    saveNote(courseId, chapterId, text) {
      const key = `${courseId}/${chapterId}`
      if (text && text.trim()) {
        this.notes[key] = text
      } else {
        delete this.notes[key]
      }
      this.persist()
    },

    recordVisit(courseId, chapterId) {
      this.lastVisited = { courseId, chapterId, at: Date.now() }
      this.persist()
    },

    saveQuizResult(quizId, score, total, answers) {
      this.quizResults[quizId] = { score, total, answers, at: Date.now() }
      this.persist()
    },

    resetAll() {
      this.completed = {}
      this.notes = {}
      this.quizResults = {}
      this.lastVisited = null
      this.studyLog = {}
      this.persist()
    }
  }
})
