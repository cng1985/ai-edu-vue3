import { Router } from 'express'
import { loadDb, saveDb } from '../utils/db.js'
import { ok, fail } from '../middleware/response.js'
import { requireRole } from '../middleware/auth.js'

const router = Router()

router.get('/', (req, res) => {
  const db = loadDb()
  const { status = '', type = '', page = 1, pageSize = 20 } = req.query
  let list = [...db.reviews]

  if (status) list = list.filter((r) => r.status === status)
  if (type) list = list.filter((r) => r.type === type)

  list.sort((a, b) => b.createdAt - a.createdAt)
  const total = list.length
  const start = (Number(page) - 1) * Number(pageSize)
  list = list.slice(start, start + Number(pageSize))

  return ok(res, { list, total, page: Number(page), pageSize: Number(pageSize) })
})

router.get('/:id', (req, res) => {
  const db = loadDb()
  const review = db.reviews.find((r) => r.id === req.params.id)
  if (!review) return fail(res, 404, '审核记录不存在')
  return ok(res, review)
})

router.post('/:id/approve', requireRole('admin', 'reviewer'), (req, res) => {
  const db = loadDb()
  const review = db.reviews.find((r) => r.id === req.params.id)
  if (!review) return fail(res, 404, '审核记录不存在')
  if (review.status !== 'pending') return fail(res, 400, '该记录已处理')

  review.status = 'approved'
  review.reviewerId = req.user.id
  review.reviewedAt = Date.now()
  review.comment = req.body?.comment || ''

  if (review.type === 'chapter') {
    const course = db.courses.find((c) => c.id === review.courseId)
    const chapter = course?.chapters.find((ch) => ch.id === review.targetId)
    if (chapter) {
      chapter.content = review.content
      chapter.status = 'published'
      chapter.updatedAt = Date.now()
    }
  } else if (review.type === 'quiz') {
    try {
      const question = JSON.parse(review.content)
      const quiz = db.quizzes.find((q) => q.id === review.targetId)
      if (quiz && question.text) {
        quiz.questions.push(question)
        quiz.updatedAt = Date.now()
      }
    } catch { /* 忽略解析错误 */ }
  }

  saveDb()
  return ok(res, review, '审核通过')
})

router.post('/:id/reject', requireRole('admin', 'reviewer'), (req, res) => {
  const db = loadDb()
  const review = db.reviews.find((r) => r.id === req.params.id)
  if (!review) return fail(res, 404, '审核记录不存在')
  if (review.status !== 'pending') return fail(res, 400, '该记录已处理')

  review.status = 'rejected'
  review.reviewerId = req.user.id
  review.reviewedAt = Date.now()
  review.comment = req.body?.comment || '内容不符合发布标准'
  saveDb()
  return ok(res, review, '已驳回')
})

export default router
