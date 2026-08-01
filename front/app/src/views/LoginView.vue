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
    await new Promise((r) => setTimeout(r, 320))
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
    await new Promise((r) => setTimeout(r, 260))
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
    <form class="form" @submit.prevent="submit">
      <label class="field">
        <input
          v-model="username"
          type="text"
          autocomplete="username"
          placeholder="用户名"
        />
      </label>

      <label class="field">
        <input
          v-model="password"
          type="password"
          autocomplete="current-password"
          placeholder="密码"
        />
      </label>

      <p v-if="error" class="error">{{ error }}</p>

      <button type="submit" class="btn-primary" :disabled="loading || guestLoading">
        {{ loading ? '登录中…' : '登录' }}
      </button>
    </form>

    <div class="extra">
      <button
        type="button"
        class="guest"
        :disabled="loading || guestLoading"
        @click="guestLogin"
      >
        {{ guestLoading ? '进入中…' : '游客登录，先看看' }}
      </button>
    </div>
  </AuthShell>
</template>

<style scoped>
.form {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.field input {
  width: 100%;
  height: 44px;
  padding: 0 14px;
  border: 1px solid #e3e6eb;
  border-radius: 4px;
  background: #fafbfc;
  font: inherit;
  font-size: 14.5px;
  color: #121212;
  outline: none;
  transition: border-color 0.15s ease, background 0.15s ease, box-shadow 0.15s ease;
}

.field input::placeholder {
  color: #b0b6c0;
}

.field input:focus {
  border-color: #1772f6;
  background: #fff;
  box-shadow: 0 0 0 3px rgba(23, 114, 246, 0.12);
}

.error {
  margin: 0;
  font-size: 13px;
  color: #e11d48;
  line-height: 1.5;
}

.btn-primary {
  margin-top: 4px;
  height: 44px;
  border: none;
  border-radius: 4px;
  background: #1772f6;
  color: #fff;
  font-size: 15px;
  font-weight: 600;
  cursor: pointer;
  transition: background 0.15s ease;
}

.btn-primary:hover:not(:disabled) {
  background: #0f62e6;
}

.btn-primary:disabled {
  opacity: 0.65;
  cursor: not-allowed;
}

.extra {
  margin-top: 18px;
  text-align: center;
}

.guest {
  border: none;
  background: transparent;
  color: #8590a6;
  font-size: 13.5px;
  cursor: pointer;
  padding: 4px 8px;
  transition: color 0.15s ease;
}

.guest:hover:not(:disabled) {
  color: #1772f6;
}

.guest:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}
</style>
