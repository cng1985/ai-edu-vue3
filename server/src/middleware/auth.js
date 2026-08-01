import { verifyToken } from '../utils/auth.js'

export function requireAuth(req, res, next) {
  const header = req.headers.authorization
  if (!header?.startsWith('Bearer ')) {
    return res.status(401).json({ code: 401, message: '未登录', data: null })
  }
  try {
    req.user = verifyToken(header.slice(7))
    next()
  } catch {
    return res.status(401).json({ code: 401, message: '登录已过期', data: null })
  }
}

export function requireAdmin(req, res, next) {
  if (!['admin', 'reviewer', 'operator'].includes(req.user?.role)) {
    return res.status(403).json({ code: 403, message: '无权限', data: null })
  }
  next()
}

export function requireRole(...roles) {
  return (req, res, next) => {
    if (!roles.includes(req.user?.role)) {
      return res.status(403).json({ code: 403, message: '无权限', data: null })
    }
    next()
  }
}
