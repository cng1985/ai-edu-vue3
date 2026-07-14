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
const avatar = ref(AVATARS[0])
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
  <AuthShell>
    <header class="form-head">
      <h2>创建账号 🚀</h2>
      <p>一分钟加入 AI 学习者社区</p>
    </header>

    <form class="form" @submit.prevent="submit">
      <div class="field">
        <span>选择你的社区头像</span>
        <div class="avatars">
          <button
            v-for="a in AVATARS"
            :key="a"
            type="button"
            class="avatar"
            :class="{ 'avatar--active': avatar === a }"
            @click="avatar = a"
          >
            {{ a }}
          </button>
        </div>
      </div>

      <div class="row">
        <label class="field">
          <span>用户名</span>
          <input
            v-model="username"
            type="text"
            autocomplete="username"
            placeholder="3~20 位字母数字"
          />
        </label>
        <label class="field">
          <span>昵称</span>
          <input v-model="nickname" type="text" placeholder="社区里怎么称呼你" />
        </label>
      </div>

      <label class="field">
        <span>密码</span>
        <input
          v-model="password"
          type="password"
          autocomplete="new-password"
          placeholder="至少 6 位"
        />
      </label>

      <label class="field">
        <span>确认密码</span>
        <input
          v-model="confirm"
          type="password"
          autocomplete="new-password"
          placeholder="再次输入密码"
        />
      </label>

      <label class="agree">
        <input v-model="agreed" type="checkbox" />
        <span>我已阅读并同意《社区公约》:友善交流、乐于分享、尊重原创</span>
      </label>

      <p v-if="error" class="form-error">⚠ {{ error }}</p>

      <button type="submit" class="btn btn--primary form-submit" :disabled="loading">
        {{ loading ? '注册中…' : '注册并加入社区' }}
      </button>
    </form>

    <p class="form-switch">
      已有账号?
      <router-link :to="{ path: '/login', query: route.query }">直接登录</router-link>
    </p>
  </AuthShell>
</template>

<style scoped>
.form-head {
  margin-bottom: 22px;
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
  gap: 6px;
  min-width: 0;
}

.field input {
  width: 100%;
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
  border: 1.5px solid var(--border);
  border-radius: 50%;
  background: var(--surface-2);
  cursor: pointer;
  transition: all 0.15s ease;
}

.avatar:hover {
  border-color: var(--primary);
  transform: translateY(-2px);
}

.avatar--active {
  border-color: var(--primary);
  background: var(--primary-soft);
  box-shadow: 0 0 0 3px rgba(99, 102, 241, 0.18);
}

.agree {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  font-size: 12.5px;
  color: var(--text-2);
  line-height: 1.6;
  cursor: pointer;
}

.agree input {
  margin-top: 3px;
  accent-color: var(--primary);
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
}

.form-switch {
  margin: 18px 0 0;
  text-align: center;
  font-size: 13.5px;
  color: var(--text-2);
}

.form-switch a {
  font-weight: 600;
}

@media (max-width: 480px) {
  .row {
    grid-template-columns: 1fr;
  }
}
</style>
