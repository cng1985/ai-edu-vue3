<script setup>
import { ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore, AVATAR_PRESETS } from '../stores/auth'
import AuthShell from '../components/AuthShell.vue'
import Icon from '../components/Icon.vue'

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
    await auth.register({
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
        <span>选择头像</span>
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
            <Icon v-if="avatarIndex === i" name="check" :size="11" :stroke="3" class="avatar__check" />
          </button>
        </div>
      </div>

      <div class="row">
        <label class="field">
          <span>用户名</span>
          <input v-model="username" class="input" type="text" autocomplete="username" placeholder="登录使用" />
        </label>
        <label class="field">
          <span>昵称</span>
          <input v-model="nickname" class="input" type="text" placeholder="社区展示名" />
        </label>
      </div>

      <label class="field">
        <span>密码</span>
        <input
          v-model="password"
          class="input"
          type="password"
          autocomplete="new-password"
          placeholder="至少 6 位"
        />
      </label>

      <label class="field">
        <span>确认密码</span>
        <input
          v-model="confirm"
          class="input"
          type="password"
          autocomplete="new-password"
          placeholder="再次输入密码"
        />
      </label>

      <label class="agree">
        <input v-model="agreed" type="checkbox" />
        <span>同意《社区公约》：友善交流、乐于分享、尊重原创</span>
      </label>

      <p v-if="error" class="error"><Icon name="alert" :size="14" /> {{ error }}</p>

      <button type="submit" class="btn btn--primary btn--lg btn--block" :disabled="loading">
        {{ loading ? '注册中…' : '创建账号' }}
      </button>
    </form>
  </AuthShell>
</template>

<style scoped>
.form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 12px;
}

.field .input {
  height: 46px;
  padding: 0 14px;
  background: var(--surface-2);
}

.field .input:focus {
  background: var(--surface);
}

.avatars {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.avatar {
  position: relative;
  width: 40px;
  height: 40px;
  border: none;
  border-radius: 13px;
  color: #fff;
  font-size: 14px;
  font-weight: 700;
  cursor: pointer;
  opacity: 0.75;
  transition: transform var(--t-fast), box-shadow var(--t-fast), opacity var(--t-fast);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.3);
}

.avatar:hover {
  transform: translateY(-2px);
  opacity: 1;
}

.avatar--on {
  opacity: 1;
  transform: scale(1.08);
  box-shadow: 0 0 0 2px var(--surface), 0 0 0 4px var(--primary);
}

.avatar__check {
  position: absolute;
  right: -5px;
  bottom: -5px;
  display: grid;
  place-items: center;
  width: 17px;
  height: 17px;
  padding: 3px;
  border-radius: 50%;
  background: var(--primary);
  color: #fff;
  box-shadow: 0 0 0 2px var(--surface);
}

.agree {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 10px;
  background: var(--surface-2);
  font-size: 12.5px;
  color: var(--text-2);
  line-height: 1.55;
  cursor: pointer;
}

.agree input {
  margin-top: 3px;
  accent-color: var(--primary);
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

@media (max-width: 480px) {
  .row {
    grid-template-columns: 1fr;
  }
}
</style>
