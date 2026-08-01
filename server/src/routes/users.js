import { Router } from 'express'
import { loadDb, saveDb, genId } from '../utils/db.js'
import { hashPassword } from '../utils/auth.js'
import { ok, fail } from '../middleware/response.js'
import { requireRole } from '../middleware/auth.js'

const router = Router()

router.get('/', (req, res) => {
  const db = loadDb()
  const { keyword = '', role = '', status = '', page = 1, pageSize = 20 } = req.query
  let list = db.users.map(({ passwordHash, ...u }) => u)

  if (keyword) {
    const kw = keyword.toLowerCase()
    list = list.filter(
      (u) => u.username.toLowerCase().includes(kw) || u.nickname.toLowerCase().includes(kw)
    )
  }
  if (role) list = list.filter((u) => u.role === role)
  if (status) list = list.filter((u) => u.status === status)

  const total = list.length
  const start = (Number(page) - 1) * Number(pageSize)
  list = list.slice(start, start + Number(pageSize))

  return ok(res, { list, total, page: Number(page), pageSize: Number(pageSize) })
})

router.get('/:id', (req, res) => {
  const db = loadDb()
  const user = db.users.find((u) => u.id === req.params.id)
  if (!user) return fail(res, 404, '用户不存在')
  const { passwordHash, ...safe } = user
  return ok(res, safe)
})

router.post('/', requireRole('admin'), (req, res) => {
  const { username, nickname, password, role = 'learner', avatar, avatarColor } = req.body || {}
  if (!username || !password) return fail(res, 400, '用户名和密码必填')

  const db = loadDb()
  if (db.users.some((u) => u.username.toLowerCase() === username.toLowerCase())) {
    return fail(res, 400, '用户名已存在')
  }

  const user = {
    id: genId('u'),
    username: username.trim(),
    nickname: nickname?.trim() || username,
    passwordHash: hashPassword(password),
    role,
    status: 'active',
    avatar: avatar || username[0].toUpperCase(),
    avatarColor: avatarColor || '#6366f1',
    joinedAt: Date.now()
  }
  db.users.push(user)
  saveDb()

  const { passwordHash, ...safe } = user
  return ok(res, safe, '创建成功')
})

router.put('/:id', requireRole('admin'), (req, res) => {
  const db = loadDb()
  const idx = db.users.findIndex((u) => u.id === req.params.id)
  if (idx === -1) return fail(res, 404, '用户不存在')

  const { nickname, role, status, avatar, avatarColor, password } = req.body || {}
  const user = db.users[idx]
  if (nickname !== undefined) user.nickname = nickname
  if (role !== undefined) user.role = role
  if (status !== undefined) user.status = status
  if (avatar !== undefined) user.avatar = avatar
  if (avatarColor !== undefined) user.avatarColor = avatarColor
  if (password) user.passwordHash = hashPassword(password)

  saveDb()
  const { passwordHash, ...safe } = user
  return ok(res, safe, '更新成功')
})

router.delete('/:id', requireRole('admin'), (req, res) => {
  const db = loadDb()
  const idx = db.users.findIndex((u) => u.id === req.params.id)
  if (idx === -1) return fail(res, 404, '用户不存在')
  if (db.users[idx].role === 'admin' && db.users.filter((u) => u.role === 'admin').length <= 1) {
    return fail(res, 400, '不能删除最后一个管理员')
  }
  db.users.splice(idx, 1)
  saveDb()
  return ok(res, null, '删除成功')
})

export default router
