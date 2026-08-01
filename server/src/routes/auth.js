import { loadDb } from '../utils/db.js'
import { verifyPassword, signToken } from '../utils/auth.js'
import { ok, fail } from '../middleware/response.js'

export function loginHandler(req, res) {
  const { username, password } = req.body || {}
  if (!username || !password) return fail(res, 400, '请输入用户名和密码')

  const db = loadDb()
  const user = db.users.find(
    (u) => u.username.toLowerCase() === username.trim().toLowerCase()
  )
  if (!user || !verifyPassword(password, user.passwordHash)) {
    return fail(res, 401, '用户名或密码错误')
  }
  if (user.status === 'disabled') {
    return fail(res, 403, '账号已被禁用')
  }
  if (!['admin', 'reviewer', 'operator'].includes(user.role)) {
    return fail(res, 403, '非管理端账号，无法登录')
  }

  const token = signToken({ id: user.id, username: user.username, role: user.role })
  const { passwordHash, ...safe } = user
  return ok(res, { token, user: safe })
}

export function meHandler(req, res) {
  const db = loadDb()
  const user = db.users.find((u) => u.id === req.user.id)
  if (!user) return fail(res, 404, '用户不存在')
  const { passwordHash, ...safe } = user
  return ok(res, safe)
}
