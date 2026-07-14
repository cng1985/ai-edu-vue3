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
const error = ref('')
const loading = ref(false)

async function submit() {
  error.value = ''
  if (!username.value.trim() || !password.value) {
    error.value = '请输入用户名和密码'
    return
  }
  loading.value = true
  try {
    // 模拟网络延迟,贴近真实登录体验
    await new Promise((r) => setTimeout(r, 350))
    auth.login(username.value, password.value)
    router.replace(typeof route.query.redirect === 'string' ? route.query.redirect : '/')
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <AuthShell>
    <header class="form-head">
      <h2>欢迎回来 👋</h2>
      <p>登录后继续你的 AI 学习之旅</p>
    </header>

    <form class="form" @submit.prevent="submit">
      <label class="field">
        <span>用户名</span>
        <input
          v-model="username"
          type="text"
          autocomplete="username"
          placeholder="你的用户名"
        />
      </label>

      <label class="field">
        <span>密码</span>
        <input
          v-model="password"
          type="password"
          autocomplete="current-password"
          placeholder="你的密码"
        />
      </label>

      <p v-if="error" class="form-error">⚠ {{ error }}</p>

      <button type="submit" class="btn btn--primary form-submit" :disabled="loading">
        {{ loading ? '登录中…' : '登 录' }}
      </button>
    </form>

    <p class="form-switch">
      还没有账号?
      <router-link :to="{ path: '/register', query: route.query }">立即注册</router-link>
    </p>
  </AuthShell>
</template>

<style scoped>
.form-head {
  margin-bottom: 24px;
}

.form-head h2 {
  margin: 0 0 5px;
  font-size: 22px;
  letter-spacing: -0.01em;
}

.form-head p {
  margin: 0;
  font-size: 13.5px;
  color: var(--text-2);
}

.form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.field span {
  font-size: 13px;
  font-weight: 600;
  color: var(--text-2);
}

.field input {
  padding: 11px 14px;
  border: 1.5px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: 14.5px;
  font-family: inherit;
  outline: none;
  transition: border-color 0.15s ease;
}

.field input:focus {
  border-color: var(--primary);
}

.form-error {
  margin: 0;
  font-size: 13px;
  color: var(--danger);
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: var(--radius-sm);
  padding: 9px 13px;
}

.form-submit {
  padding: 12px;
  font-size: 15px;
  margin-top: 4px;
}

.form-switch {
  margin: 20px 0 0;
  text-align: center;
  font-size: 13.5px;
  color: var(--text-2);
}

.form-switch a {
  font-weight: 600;
}
</style>
