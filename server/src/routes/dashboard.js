import { Router } from 'express'
import { loadDb } from '../utils/db.js'
import { ok } from '../middleware/response.js'

const router = Router()

router.get('/stats', (req, res) => {
  const db = loadDb()

  const userStats = {
    total: db.users.length,
    learners: db.users.filter((u) => u.role === 'learner').length,
    admins: db.users.filter((u) => ['admin', 'operator', 'reviewer'].includes(u.role)).length,
    active: db.users.filter((u) => u.status === 'active').length
  }

  const courseStats = {
    total: db.courses.length,
    published: db.courses.filter((c) => c.status === 'published').length,
    draft: db.courses.filter((c) => c.status === 'draft').length,
    chapters: db.courses.reduce((sum, c) => sum + (c.chapters?.length || 0), 0)
  }

  const quizStats = {
    total: db.quizzes.length,
    questions: db.quizzes.reduce((sum, q) => sum + (q.questions?.length || 0), 0)
  }

  const reviewStats = {
    pending: db.reviews.filter((r) => r.status === 'pending').length,
    approved: db.reviews.filter((r) => r.status === 'approved').length,
    rejected: db.reviews.filter((r) => r.status === 'rejected').length
  }

  return ok(res, { userStats, courseStats, quizStats, reviewStats })
})

export default router
