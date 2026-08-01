import express from 'express'
import cors from 'cors'
import fs from 'fs'
import path from 'path'
import { fileURLToPath } from 'url'
import { requireAuth, requireAdmin } from './middleware/auth.js'
import { loginHandler, meHandler } from './routes/auth.js'
import userRoutes from './routes/users.js'
import courseRoutes from './routes/courses.js'
import quizRoutes from './routes/quizzes.js'
import reviewRoutes from './routes/reviews.js'
import dashboardRoutes from './routes/dashboard.js'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const DB_PATH = path.join(__dirname, '../data/store.json')
const PORT = process.env.PORT || 3001

if (!fs.existsSync(DB_PATH)) {
  console.log('📦 首次启动，正在初始化数据...')
  await import('./seed.js')
}

const app = express()
app.use(cors())
app.use(express.json({ limit: '5mb' }))

app.get('/api/v1/health', (_req, res) => {
  res.json({ code: 0, message: 'ok', data: { status: 'healthy' } })
})

app.post('/api/v1/auth/login', loginHandler)
app.get('/api/v1/auth/me', requireAuth, requireAdmin, meHandler)
app.use('/api/v1/users', requireAuth, requireAdmin, userRoutes)
app.use('/api/v1/courses', requireAuth, requireAdmin, courseRoutes)
app.use('/api/v1/quizzes', requireAuth, requireAdmin, quizRoutes)
app.use('/api/v1/reviews', requireAuth, requireAdmin, reviewRoutes)
app.use('/api/v1/dashboard', requireAuth, requireAdmin, dashboardRoutes)

app.use((err, _req, res, _next) => {
  console.error(err)
  res.status(500).json({ code: 500, message: '服务器内部错误', data: null })
})

app.listen(PORT, () => {
  console.log(`🚀 API 服务已启动: http://localhost:${PORT}`)
  console.log(`   管理端登录: admin / admin123`)
})
