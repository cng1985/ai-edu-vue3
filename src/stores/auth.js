import { defineStore } from 'pinia'

// 前端演示版账号体系:用户注册信息与会话均保存在 localStorage。
// 接入真实后端时,替换 register/login/logout 中的存储逻辑为 API 调用即可。

const USERS_KEY = 'ai-learning-system:users'
const SESSION_KEY = 'ai-learning-system:session'

export const AVATARS = ['🧑‍💻', '👩‍💻', '🦊', '🐼', '🦉', '🐯', '🚀', '🤖', '🐱', '🌟']

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

// 演示用途的简易散列,避免明文存储;生产环境必须由后端完成密码处理
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
    register({ username, nickname, password, avatar }) {
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

      const user = {
        id: 'u' + Date.now().toString(36),
        username,
        nickname,
        avatar: avatar || AVATARS[0],
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

    /** 游客登录:无需账号密码,创建临时会话,可立即体验全部功能 */
    loginAsGuest() {
      const stamp = Date.now().toString(36).slice(-4)
      const user = {
        id: 'guest-' + Date.now().toString(36),
        username: 'guest_' + stamp,
        nickname: '游客',
        avatar: '👋',
        isGuest: true,
        joinedAt: Date.now()
      }
      this.setSession(user)
      return user
    },

    setSession(user) {
      // 会话中不保留密码散列
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
