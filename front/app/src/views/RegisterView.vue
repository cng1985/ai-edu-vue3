<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore, AVATAR_PRESETS } from '../stores/auth'
import AuthShell from '../components/AuthShell.vue'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()

const username = ref('')
const nickname = ref('')
const password = ref('')
const confirm = ref('')
const avatarIndex = ref(0)
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
    await new Promise((r) => setTimeout(r, 320))
    const preset = AVATAR_PRESETS[avatarIndex.value]
    auth.register({
      username: username.value,
      nickname: nickname.value,
      password: password.value,
      avatar: preset.char,
      avatarColor: preset.color
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
    <form class="form" @submit.prevent="submit">
      <div class="field">
        <span class="label">头像</span>
        <div class="avatars">
          <button
            v-for="(preset, i) in AVATAR_PRESETS"
            :key="preset.char"
            type="button"
            class="avatar"
            :class="{ 'avatar--on': avatarIndex === i }"
            :style="{ background: preset.color }"
            :aria-label="`选择头像 ${preset.char}`"
            @click="avatarIndex = i"
          >
            {{ preset.char }}
          </button>
        </div>
      </div>

      <div class="row">
        <label class="field">
          <input v-model="username" type="text" autocomplete="username" placeholder="用户名" />
        </label>
        <label class="field">
          <input v-model="nickname" type="text" placeholder="昵称" />
        </label>
      </div>

      <label class="field">
        <input
          v-model="password"
          type="password"
          autocomplete="new-password"
          placeholder="密码（至少 6 位）"
        />
      </label>

      <label class="field">
        <input
          v-model="confirm"
          type="password"
          autocomplete="new-password"
          placeholder="确认密码"
        />
      </label>

      <label class="agree">
        <input v-model="agreed" type="checkbox" />
        <span>同意《社区公约》：友善交流、乐于分享、尊重原创</span>
      </label>

      <p v-if="error" class="error">{{ error }}</p>

      <button type="submit" class="btn-primary" :disabled="loading">
        {{ loading ? '注册中…' : '注册' }}
      </button>
    </form>
  </AuthShell>
</template>

<style scoped>
.form {
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 8px;
  min-width: 0;
}

.label {
  font-size: 13px;
  color: #646464;
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

.avatars {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
}

.avatar {
  width: 36px;
  height: 36px;
  border: 2px solid transparent;
  border-radius: 50%;
  color: #fff;
  font-size: 13px;
  font-weight: 700;
  cursor: pointer;
  transition: transform 0.12s ease, box-shadow 0.12s ease;
}

.avatar:hover {
  transform: translateY(-1px);
}

.avatar--on {
  box-shadow: 0 0 0 2px #fff, 0 0 0 4px #1772f6;
}

.agree {
  display: flex;
  align-items: flex-start;
  gap: 8px;
  font-size: 12.5px;
  color: #8590a6;
  line-height: 1.55;
  cursor: pointer;
}

.agree input {
  margin-top: 2px;
  accent-color: #1772f6;
}

.error {
  margin: 0;
  font-size: 13px;
  color: #e11d48;
}

.btn-primary {
  margin-top: 2px;
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

@media (max-width: 480px) {
  .row {
    grid-template-columns: 1fr;
  }
}
</style>
