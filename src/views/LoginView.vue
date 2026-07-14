<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import AuthShell from '../components/AuthShell.vue'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const username = ref('')
const password = ref('')
const showPassword = ref(false)
const error = ref('')
const loading = ref(false)
const guestLoading = ref(false)

function redirectAfterAuth() {
  router.replace(typeof route.query.redirect === 'string' ? route.query.redirect : '/')
}

async function submit() {
  error.value = ''
  if (!username.value.trim() || !password.value) {
    error.value = '请输入用户名和密码'
    return
  }
  loading.value = true
  try {
    await new Promise((r) => setTimeout(r, 350))
    auth.login(username.value, password.value)
    redirectAfterAuth()
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}

async function guestLogin() {
  error.value = ''
  guestLoading.value = true
  try {
    await new Promise((r) => setTimeout(r, 280))
    auth.loginAsGuest()
    redirectAfterAuth()
  } catch (e) {
    error.value = e.message || '游客登录失败'
  } finally {
    guestLoading.value = false
  }
}
</script>

<template>
  <AuthShell mode="login">
    <header class="head">
      <h2 class="font-display">欢迎回来</h2>
      <p>登录账号,继续你的 AI 学习之旅</p>
    </header>

    <form class="form" @submit.prevent="submit">
      <label class="field">
        <span class="field__label">用户名</span>
        <div class="field__box">
          <span class="field__icon" aria-hidden="true">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
              <path d="M12 12a4.5 4.5 0 1 0-4.5-4.5A4.5 4.5 0 0 0 12 12Zm0 2.25c-3.75 0-6.75 1.9-6.75 4.25V20h13.5v-1.5c0-2.35-3-4.25-6.75-4.25Z" fill="currentColor"/>
            </svg>
          </span>
          <input
            v-model="username"
            type="text"
            autocomplete="username"
            placeholder="输入用户名"
          />
        </div>
      </label>

      <label class="field">
        <span class="field__label">密码</span>
        <div class="field__box">
          <span class="field__icon" aria-hidden="true">
            <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
              <path d="M17 9h-1V7a4 4 0 1 0-8 0v2H7a2 2 0 0 0-2 2v8a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2v-8a2 2 0 0 0-2-2Zm-6.5-2a2.5 2.5 0 0 1 5 0v2h-5V7Zm1.5 8.8a1.5 1.5 0 1 1 0-3 1.5 1.5 0 0 1 0 3Z" fill="currentColor"/>
            </svg>
          </span>
          <input
            v-model="password"
            :type="showPassword ? 'text' : 'password'"
            autocomplete="current-password"
            placeholder="输入密码"
          />
          <button
            type="button"
            class="field__toggle"
            :aria-label="showPassword ? '隐藏密码' : '显示密码'"
            @click="showPassword = !showPassword"
          >
            {{ showPassword ? '隐藏' : '显示' }}
          </button>
        </div>
      </label>

      <p v-if="error" class="error">{{ error }}</p>

      <button type="submit" class="submit" :disabled="loading || guestLoading">
        <span v-if="loading" class="spinner"></span>
        {{ loading ? '登录中…' : '登录' }}
      </button>
    </form>

    <div class="divider"><span>或者</span></div>

    <button
      type="button"
      class="guest"
      :disabled="loading || guestLoading"
      @click="guestLogin"
    >
      <span class="guest__icon">👋</span>
      <span class="guest__text">
        <strong>{{ guestLoading ? '进入中…' : '游客登录' }}</strong>
        <em>无需注册,立即体验全部课程</em>
      </span>
    </button>

    <p class="switch">
      还没有账号?
      <router-link :to="{ path: '/register', query: route.query }">免费注册</router-link>
    </p>
  </AuthShell>
</template>

<style scoped>
.head {
  margin-bottom: 28px;
}

.head h2 {
  margin: 0 0 6px;
  font-size: 28px;
  font-weight: 700;
  letter-spacing: -0.03em;
  color: #0f172a;
}

.head p {
  margin: 0;
  font-size: 14px;
  color: #64748b;
}

.form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.field__label {
  font-size: 13px;
  font-weight: 600;
  color: #334155;
}

.field__box {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 14px;
  min-height: 50px;
  border: 1.5px solid #e2e8f0;
  border-radius: 14px;
  background: #f8fafc;
  transition: border-color 0.15s ease, box-shadow 0.15s ease, background 0.15s ease;
}

.field__box:focus-within {
  border-color: #6366f1;
  background: #fff;
  box-shadow: 0 0 0 4px rgba(99, 102, 241, 0.12);
}

.field__icon {
  display: flex;
  color: #94a3b8;
  flex-shrink: 0;
}

.field__box:focus-within .field__icon {
  color: #6366f1;
}

.field__box input {
  flex: 1;
  min-width: 0;
  border: none;
  outline: none;
  background: transparent;
  font: inherit;
  font-size: 14.5px;
  color: #0f172a;
  padding: 12px 0;
}

.field__box input::placeholder {
  color: #94a3b8;
}

.field__toggle {
  border: none;
  background: transparent;
  color: #6366f1;
  font-size: 12.5px;
  font-weight: 600;
  cursor: pointer;
  padding: 4px 2px;
  white-space: nowrap;
}

.error {
  margin: 0;
  padding: 11px 14px;
  border-radius: 12px;
  background: #fef2f2;
  border: 1px solid #fecaca;
  color: #dc2626;
  font-size: 13px;
}

.submit {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-top: 4px;
  min-height: 50px;
  border: none;
  border-radius: 14px;
  background: linear-gradient(135deg, #6366f1 0%, #4f46e5 100%);
  color: #fff;
  font-size: 15px;
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 10px 24px rgba(79, 70, 229, 0.28);
  transition: transform 0.15s ease, box-shadow 0.15s ease, filter 0.15s ease;
}

.submit:hover:not(:disabled) {
  transform: translateY(-1px);
  box-shadow: 0 14px 28px rgba(79, 70, 229, 0.34);
}

.submit:disabled {
  opacity: 0.65;
  cursor: not-allowed;
}

.spinner {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.35);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

.divider {
  display: flex;
  align-items: center;
  gap: 14px;
  margin: 22px 0 16px;
  color: #94a3b8;
  font-size: 12.5px;
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: #e2e8f0;
}

.guest {
  display: flex;
  align-items: center;
  gap: 14px;
  width: 100%;
  padding: 14px 16px;
  border: 1.5px solid #e2e8f0;
  border-radius: 14px;
  background: #fff;
  cursor: pointer;
  text-align: left;
  transition: border-color 0.15s ease, background 0.15s ease, transform 0.15s ease;
}

.guest:hover:not(:disabled) {
  border-color: #a5b4fc;
  background: #f8fafc;
  transform: translateY(-1px);
}

.guest:disabled {
  opacity: 0.65;
  cursor: not-allowed;
}

.guest__icon {
  width: 42px;
  height: 42px;
  min-width: 42px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  background: linear-gradient(135deg, #eef2ff, #e0f2fe);
  font-size: 20px;
}

.guest__text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.guest__text strong {
  font-size: 14.5px;
  color: #0f172a;
}

.guest__text em {
  font-style: normal;
  font-size: 12.5px;
  color: #64748b;
}

.switch {
  margin: 22px 0 0;
  text-align: center;
  font-size: 13.5px;
  color: #64748b;
}

.switch a {
  font-weight: 700;
  color: #4f46e5;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
