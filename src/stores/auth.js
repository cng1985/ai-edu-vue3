import { defineStore } from 'pinia'

// 前端演示版账号体系:用户注册信息与会话均保存在 localStorage。
// 接入真实后端时,替换 register/login/logout 中的存储逻辑为 API 调用即可。

const USERS_KEY = 'ai-learning-system:users'
const SESSION_KEY = 'ai-learning-system:session'

/** 社区字母头像配色（替代 emoji，更接近知乎/掘金常见样式） */
export const AVATAR_PRESETS = [
  { char: 'A', color: '#1772f6' },
  { char: 'B', color: '#0ea5e9' },
  { char: 'C', color: '#10b981' },
  { char: 'D', color: '#f59e0b' },
  { char: 'E', color: '#ef4444' },
  { char: 'F', color: '#8b5cf6' },
  { char: 'G', color: '#14b8a6' },
  { char: 'H', color: '#f97316' }
]

function loadUsers() {
  try {
    const raw = localStorage.getItem(USERS_KEY)
    if (raw) return JSON.parse(raw)
  } catch (e) { /* 忽略损坏数据 */ }
  return []
}

function saveUsers(users) {
  localStorage.setItem(USERS_KEY, JSON.stringify(users))
}

function loadSession() {
  try {
    const raw = localStorage.getItem(SESSION_KEY)
    if (raw) return JSON.parse(raw)
  } catch (e) { /* 忽略损坏数据 */ }
  return null
}

function hash(text) {
  let h = 5381
  for (let i = 0; i < text.length; i++) {
    h = ((h << 5) + h + text.charCodeAt(i)) >>> 0
  }
  return 'h' + h.toString(36) + ':' + text.length
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: loadSession()
  }),

  getters: {
    isLoggedIn: (state) => Boolean(state.user),
    isGuest: (state) => Boolean(state.user?.isGuest)
  },

  actions: {
    register({ username, nickname, password, avatar, avatarColor }) {
      username = username.trim()
      nickname = nickname.trim()

      if (!/^[a-zA-Z0-9_-]{3,20}$/.test(username)) {
        throw new Error('用户名需为 3~20 位字母、数字、下划线或短横线')
      }
      if (!nickname) throw new Error('请填写昵称')
      if (nickname.length > 16) throw new Error('昵称最长 16 个字符')
      if (password.length < 6) throw new Error('密码至少 6 位')

      const users = loadUsers()
      if (users.some((u) => u.username.toLowerCase() === username.toLowerCase())) {
        throw new Error('该用户名已被注册')
      }

      const preset = AVATAR_PRESETS[0]
      const user = {
        id: 'u' + Date.now().toString(36),
        username,
        nickname,
        avatar: avatar || preset.char,
        avatarColor: avatarColor || preset.color,
        passwordHash: hash(password),
        isGuest: false,
        joinedAt: Date.now()
      }
      users.push(user)
      saveUsers(users)
      this.setSession(user)
      return user
    },

    login(username, password) {
      const users = loadUsers()
      const user = users.find(
        (u) => u.username.toLowerCase() === username.trim().toLowerCase()
      )
      if (!user || user.passwordHash !== hash(password)) {
        throw new Error('用户名或密码错误')
      }
      this.setSession(user)
      return user
    },

    loginAsGuest() {
      const stamp = Date.now().toString(36).slice(-4)
      const user = {
        id: 'guest-' + Date.now().toString(36),
        username: 'guest_' + stamp,
        nickname: '游客',
        avatar: '访',
        avatarColor: '#94a3b8',
        isGuest: true,
        joinedAt: Date.now()
      }
      this.setSession(user)
      return user
    },

    setSession(user) {
      const { passwordHash, ...safe } = user
      this.user = safe
      localStorage.setItem(SESSION_KEY, JSON.stringify(safe))
    },

    logout() {
      this.user = null
      localStorage.removeItem(SESSION_KEY)
    }
  }
})
