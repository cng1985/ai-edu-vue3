<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import AuthShell from '../components/AuthShell.vue'
import { homeFor } from '../router'
import Icon from '../components/Icon.vue'

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
  router.replace(typeof route.query.redirect === 'string' ? route.query.redirect : homeFor(auth))
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
    await auth.login(username.value, password.value)
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
    await auth.loginAsGuest()
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
        <span>用户名</span>
        <div class="control">
          <Icon name="user" :size="17" class="control__icon" />
          <input
            v-model="username"
            type="text"
            autocomplete="username"
            placeholder="请输入用户名"
          />
        </div>
      </label>

      <label class="field">
        <span>密码</span>
        <div class="control">
          <Icon name="lock" :size="17" class="control__icon" />
          <input
            v-model="password"
            :type="showPassword ? 'text' : 'password'"
            autocomplete="current-password"
            placeholder="请输入密码"
          />
          <button
            type="button"
            class="control__toggle"
            :aria-label="showPassword ? '隐藏密码' : '显示密码'"
            @click="showPassword = !showPassword"
          >
            <Icon name="eye" :size="16" />
          </button>
        </div>
      </label>

      <p v-if="error" class="error"><Icon name="alert" :size="14" /> {{ error }}</p>

      <button type="submit" class="btn btn--primary btn--lg btn--block" :disabled="loading || guestLoading">
        {{ loading ? '登录中…' : '登录' }}
      </button>
    </form>

    <div class="divider"><span>或</span></div>

    <button
      type="button"
      class="btn btn--ghost btn--block guest"
      :disabled="loading || guestLoading"
      @click="guestLogin"
    >
      <Icon name="eye" :size="16" />
      {{ guestLoading ? '进入中…' : '游客模式，先看看' }}
    </button>

    <p class="hint">演示账号：学习者 <code>demo</code> / <code>demo123</code> · 企业 <code>company</code> / <code>company123</code></p>
  </AuthShell>
</template>

<style scoped>
.form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.control {
  position: relative;
  display: flex;
  align-items: center;
}

.control__icon {
  position: absolute;
  left: 14px;
  color: var(--text-3);
  pointer-events: none;
}

.control input {
  width: 100%;
  height: 46px;
  padding: 0 44px 0 42px;
  border: 1px solid var(--border-strong);
  border-radius: var(--radius-sm);
  background: var(--surface-2);
  font: inherit;
  font-size: 14.5px;
  color: var(--text);
  outline: none;
  transition: border-color var(--t-fast), background var(--t-fast), box-shadow var(--t-fast);
}

.control input::placeholder {
  color: var(--text-3);
}

.control input:focus {
  border-color: var(--primary);
  background: var(--surface);
  box-shadow: 0 0 0 4px var(--primary-soft-2);
}

.control__toggle {
  position: absolute;
  right: 8px;
  display: grid;
  place-items: center;
  width: 32px;
  height: 32px;
  border: none;
  border-radius: 8px;
  background: transparent;
  color: var(--text-3);
  cursor: pointer;
  transition: all var(--t-fast);
}

.control__toggle:hover {
  background: var(--surface-3);
  color: var(--text);
}

.error {
  display: flex;
  align-items: center;
  gap: 6px;
  margin: 0;
  padding: 9px 12px;
  border-radius: 10px;
  background: var(--danger-soft);
  font-size: 13px;
  color: var(--danger-strong);
}

.divider {
  display: flex;
  align-items: center;
  gap: 12px;
  margin: 20px 0 14px;
  color: var(--text-3);
  font-size: 12px;
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--border);
}

.guest {
  height: 44px;
}

.hint {
  margin: 16px 0 0;
  text-align: center;
  color: var(--text-3);
  font-size: 12.5px;
}

.hint code {
  padding: 1px 6px;
  border-radius: 5px;
  background: var(--surface-3);
  color: var(--text-2);
  font-family: 'JetBrains Mono', Menlo, monospace;
  font-size: 12px;
}
</style>
