import { Router } from 'express'
import { loadDb, saveDb, genId } from '../utils/db.js'
import { ok, fail } from '../middleware/response.js'
import { requireRole } from '../middleware/auth.js'

const router = Router()

router.get('/', (req, res) => {
  const db = loadDb()
  const { keyword = '', courseId = '', status = '', page = 1, pageSize = 20 } = req.query
  let list = db.quizzes.map((q) => ({
    ...q,
    questionCount: q.questions?.length || 0
  }))

  if (keyword) {
    const kw = keyword.toLowerCase()
    list = list.filter((q) => q.title.toLowerCase().includes(kw))
  }
  if (courseId) list = list.filter((q) => q.courseId === courseId)
  if (status) list = list.filter((q) => q.status === status)

  const total = list.length
  const start = (Number(page) - 1) * Number(pageSize)
  list = list.slice(start, start + Number(pageSize))

  return ok(res, { list, total, page: Number(page), pageSize: Number(pageSize) })
})

router.get('/:id', (req, res) => {
  const db = loadDb()
  const quiz = db.quizzes.find((q) => q.id === req.params.id)
  if (!quiz) return fail(res, 404, '测验不存在')
  return ok(res, quiz)
})

router.post('/', requireRole('admin', 'operator'), (req, res) => {
  const body = req.body || {}
  if (!body.title || !body.courseId) return fail(res, 400, '标题和关联课程必填')

  const db = loadDb()
  const id = body.id || genId('quiz')
  if (db.quizzes.some((q) => q.id === id)) return fail(res, 400, '测验 ID 已存在')

  const quiz = {
    id,
    courseId: body.courseId,
    title: body.title,
    description: body.description || '',
    questions: body.questions || [],
    status: body.status || 'draft',
    createdAt: Date.now(),
    updatedAt: Date.now()
  }
  db.quizzes.push(quiz)
  saveDb()
  return ok(res, quiz, '创建成功')
})

router.put('/:id', requireRole('admin', 'operator'), (req, res) => {
  const db = loadDb()
  const idx = db.quizzes.findIndex((q) => q.id === req.params.id)
  if (idx === -1) return fail(res, 404, '测验不存在')

  const allowed = ['title', 'description', 'courseId', 'questions', 'status']
  const quiz = db.quizzes[idx]
  for (const key of allowed) {
    if (req.body[key] !== undefined) quiz[key] = req.body[key]
  }
  quiz.updatedAt = Date.now()
  saveDb()
  return ok(res, quiz, '更新成功')
})

router.delete('/:id', requireRole('admin'), (req, res) => {
  const db = loadDb()
  const idx = db.quizzes.findIndex((q) => q.id === req.params.id)
  if (idx === -1) return fail(res, 404, '测验不存在')
  db.quizzes.splice(idx, 1)
  saveDb()
  return ok(res, null, '删除成功')
})

export default router
