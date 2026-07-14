<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore, AVATARS } from '../stores/auth'
import AuthShell from '../components/AuthShell.vue'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const username = ref('')
const nickname = ref('')
const password = ref('')
const confirm = ref('')
const avatar = ref(AVATARS[2])
const agreed = ref(false)
const error = ref('')
const loading = ref(false)

async function submit() {
  error.value = ''
  if (password.value !== confirm.value) {
    error.value = '两次输入的密码不一致'
    return
  }
  if (!agreed.value) {
    error.value = '请先阅读并同意社区公约'
    return
  }
  loading.value = true
  try {
    await new Promise((r) => setTimeout(r, 350))
    auth.register({
      username: username.value,
      nickname: nickname.value,
      password: password.value,
      avatar: avatar.value
    })
    router.replace(typeof route.query.redirect === 'string' ? route.query.redirect : '/')
  } catch (e) {
    error.value = e.message
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <AuthShell mode="register">
    <header class="head">
      <h2 class="font-display">创建账号</h2>
      <p>一分钟加入 AI 学习者社区</p>
    </header>

    <form class="form" @submit.prevent="submit">
      <div class="field">
        <span class="field__label">选择社区头像</span>
        <div class="avatars">
          <button
            v-for="a in AVATARS"
            :key="a"
            type="button"
            class="avatar"
            :class="{ 'avatar--active': avatar === a }"
            :aria-label="`选择头像 ${a}`"
            @click="avatar = a"
          >
            {{ a }}
          </button>
        </div>
      </div>

      <div class="row">
        <label class="field">
          <span class="field__label">用户名</span>
          <div class="field__box">
            <input
              v-model="username"
              type="text"
              autocomplete="username"
              placeholder="3~20 位字母数字"
            />
          </div>
        </label>
        <label class="field">
          <span class="field__label">昵称</span>
          <div class="field__box">
            <input v-model="nickname" type="text" placeholder="社区怎么称呼你" />
          </div>
        </label>
      </div>

      <label class="field">
        <span class="field__label">密码</span>
        <div class="field__box">
          <input
            v-model="password"
            type="password"
            autocomplete="new-password"
            placeholder="至少 6 位"
          />
        </div>
      </label>

      <label class="field">
        <span class="field__label">确认密码</span>
        <div class="field__box">
          <input
            v-model="confirm"
            type="password"
            autocomplete="new-password"
            placeholder="再次输入密码"
          />
        </div>
      </label>

      <label class="agree">
        <input v-model="agreed" type="checkbox" />
        <span>我已阅读并同意《社区公约》:友善交流、乐于分享、尊重原创</span>
      </label>

      <p v-if="error" class="error">{{ error }}</p>

      <button type="submit" class="submit" :disabled="loading">
        <span v-if="loading" class="spinner"></span>
        {{ loading ? '注册中…' : '注册并加入社区' }}
      </button>
    </form>

    <p class="switch">
      已有账号?
      <router-link :to="{ path: '/login', query: route.query }">直接登录</router-link>
      <span class="switch__dot">·</span>
      <router-link :to="{ path: '/login', query: route.query }">游客体验</router-link>
    </p>
  </AuthShell>
</template>

<style scoped>
.head {
  margin-bottom: 24px;
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
  gap: 14px;
}

.row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 0;
}

.field__label {
  font-size: 13px;
  font-weight: 600;
  color: #334155;
}

.field__box {
  display: flex;
  align-items: center;
  min-height: 48px;
  padding: 0 14px;
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

.field__box input {
  width: 100%;
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

.avatars {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.avatar {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  border: 1.5px solid #e2e8f0;
  border-radius: 50%;
  background: #f8fafc;
  cursor: pointer;
  transition: transform 0.15s ease, border-color 0.15s ease, box-shadow 0.15s ease;
}

.avatar:hover {
  transform: translateY(-2px);
  border-color: #a5b4fc;
}

.avatar--active {
  border-color: #6366f1;
  background: #eef2ff;
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.18);
}

.agree {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  font-size: 12.5px;
  color: #64748b;
  line-height: 1.6;
  cursor: pointer;
}

.agree input {
  margin-top: 3px;
  accent-color: #6366f1;
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
  min-height: 50px;
  border: none;
  border-radius: 14px;
  background: linear-gradient(135deg, #6366f1 0%, #4f46e5 100%);
  color: #fff;
  font-size: 15px;
  font-weight: 700;
  cursor: pointer;
  box-shadow: 0 10px 24px rgba(79, 70, 229, 0.28);
  transition: transform 0.15s ease, box-shadow 0.15s ease;
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

.switch {
  margin: 20px 0 0;
  text-align: center;
  font-size: 13.5px;
  color: #64748b;
}

.switch a {
  font-weight: 700;
  color: #4f46e5;
}

.switch__dot {
  margin: 0 6px;
  color: #cbd5e1;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

@media (max-width: 480px) {
  .row {
    grid-template-columns: 1fr;
  }
}
</style>
