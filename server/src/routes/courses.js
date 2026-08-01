import { Router } from 'express'
import { loadDb, saveDb, genId } from '../utils/db.js'
import { ok, fail } from '../middleware/response.js'
import { requireRole } from '../middleware/auth.js'

const router = Router()

router.get('/', (req, res) => {
  const db = loadDb()
  const { keyword = '', status = '', page = 1, pageSize = 20 } = req.query
  let list = db.courses.map((c) => ({
    ...c,
    chapterCount: c.chapters?.length || 0
  }))

  if (keyword) {
    const kw = keyword.toLowerCase()
    list = list.filter((c) => c.title.toLowerCase().includes(kw) || c.id.includes(kw))
  }
  if (status) list = list.filter((c) => c.status === status)

  const total = list.length
  const start = (Number(page) - 1) * Number(pageSize)
  list = list.slice(start, start + Number(pageSize))

  return ok(res, { list, total, page: Number(page), pageSize: Number(pageSize) })
})

router.get('/:id', (req, res) => {
  const db = loadDb()
  const course = db.courses.find((c) => c.id === req.params.id)
  if (!course) return fail(res, 404, '课程不存在')
  return ok(res, course)
})

router.post('/', requireRole('admin', 'operator'), (req, res) => {
  const body = req.body || {}
  if (!body.title) return fail(res, 400, '课程标题必填')

  const db = loadDb()
  const id = body.id || genId('course')
  if (db.courses.some((c) => c.id === id)) return fail(res, 400, '课程 ID 已存在')

  const course = {
    id,
    title: body.title,
    description: body.description || '',
    level: body.level || '入门',
    tags: body.tags || [],
    icon: body.icon || '📚',
    accent: body.accent || '#6366f1',
    estimatedMinutes: body.estimatedMinutes || 60,
    status: body.status || 'draft',
    chapters: body.chapters || [],
    createdAt: Date.now(),
    updatedAt: Date.now()
  }
  db.courses.push(course)
  saveDb()
  return ok(res, course, '创建成功')
})

router.put('/:id', requireRole('admin', 'operator'), (req, res) => {
  const db = loadDb()
  const idx = db.courses.findIndex((c) => c.id === req.params.id)
  if (idx === -1) return fail(res, 404, '课程不存在')

  const allowed = ['title', 'description', 'level', 'tags', 'icon', 'accent', 'estimatedMinutes', 'status', 'chapters']
  const course = db.courses[idx]
  for (const key of allowed) {
    if (req.body[key] !== undefined) course[key] = req.body[key]
  }
  course.updatedAt = Date.now()
  saveDb()
  return ok(res, course, '更新成功')
})

router.delete('/:id', requireRole('admin'), (req, res) => {
  const db = loadDb()
  const idx = db.courses.findIndex((c) => c.id === req.params.id)
  if (idx === -1) return fail(res, 404, '课程不存在')
  db.courses.splice(idx, 1)
  saveDb()
  return ok(res, null, '删除成功')
})

router.post('/:id/chapters', requireRole('admin', 'operator'), (req, res) => {
  const db = loadDb()
  const course = db.courses.find((c) => c.id === req.params.id)
  if (!course) return fail(res, 404, '课程不存在')

  const { id, title, minutes, content } = req.body || {}
  if (!title) return fail(res, 400, '章节标题必填')

  const chapterId = id || genId('ch')
  if (course.chapters.some((ch) => ch.id === chapterId)) {
    return fail(res, 400, '章节 ID 已存在')
  }

  const chapter = {
    id: chapterId,
    title,
    minutes: minutes || 10,
    content: content || '',
    status: 'draft',
    updatedAt: Date.now()
  }
  course.chapters.push(chapter)
  course.updatedAt = Date.now()
  saveDb()
  return ok(res, chapter, '章节创建成功')
})

router.put('/:id/chapters/:chapterId', requireRole('admin', 'operator'), (req, res) => {
  const db = loadDb()
  const course = db.courses.find((c) => c.id === req.params.id)
  if (!course) return fail(res, 404, '课程不存在')

  const idx = course.chapters.findIndex((ch) => ch.id === req.params.chapterId)
  if (idx === -1) return fail(res, 404, '章节不存在')

  const chapter = course.chapters[idx]
  const allowed = ['title', 'minutes', 'content', 'status']
  for (const key of allowed) {
    if (req.body[key] !== undefined) chapter[key] = req.body[key]
  }
  chapter.updatedAt = Date.now()
  course.updatedAt = Date.now()
  saveDb()
  return ok(res, chapter, '章节更新成功')
})

router.delete('/:id/chapters/:chapterId', requireRole('admin'), (req, res) => {
  const db = loadDb()
  const course = db.courses.find((c) => c.id === req.params.id)
  if (!course) return fail(res, 404, '课程不存在')

  const idx = course.chapters.findIndex((ch) => ch.id === req.params.chapterId)
  if (idx === -1) return fail(res, 404, '章节不存在')
  course.chapters.splice(idx, 1)
  course.updatedAt = Date.now()
  saveDb()
  return ok(res, null, '章节删除成功')
})

export default router
