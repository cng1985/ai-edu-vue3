<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useLearningStore } from '../stores/learning'
import { useAuthStore } from '../stores/auth'
import { useGrowthStore } from '../stores/growth'
import { PERM } from '../constants/permissions'

defineProps({
  open: { type: Boolean, default: false }
})
const emit = defineEmits(['navigate'])

const router = useRouter()
const learning = useLearningStore()
const auth = useAuthStore()
const growth = useGrowthStore()

function logout() {
  auth.logout()
  emit('navigate')
  router.replace('/login')
}

const navItems = [
  { to: '/', icon: '🏠', label: '学习驾驶舱', exact: true },
  { to: '/career', icon: '🧭', label: '职业与目标' },
  { to: '/path', icon: '🗺️', label: '学习路径' },
  { to: '/review', icon: '🔍', label: '复习与补强' },
  { to: '/courses', icon: '📚', label: '全部课程', permission: PERM.COURSE_READ },
  { to: '/chat', icon: '💬', label: 'AI 学习助手', permission: PERM.AI_CHAT },
  { to: '/quiz', icon: '📝', label: '知识测验', permission: PERM.QUIZ_READ },
  { to: '/stats', icon: '📊', label: '达成度报告' },
  { to: '/incentives', icon: '🏅', label: '成长激励' }
]

const visibleNavItems = computed(() =>
  navItems.filter((item) => !item.permission || auth.hasPermission(item.permission))
)

const progress = computed(() => growth.hasGoal ? growth.achievement : learning.overallProgress)
</script>

<template>
  <aside class="sidebar" :class="{ 'sidebar--open': open }">
    <router-link to="/" class="sidebar__brand" @click="$emit('navigate')">
      <span class="sidebar__logo">🧠</span>
      <span class="sidebar__title">AI 学习系统</span>
    </router-link>

    <nav class="sidebar__nav">
      <router-link
        v-for="item in visibleNavItems"
        :key="item.to"
        :to="item.to"
        class="sidebar__link"
        :class="{
          'sidebar__link--active': item.exact
            ? $route.path === item.to
            : $route.path.startsWith(item.to)
        }"
        @click="$emit('navigate')"
      >
        <span class="sidebar__icon">{{ item.icon }}</span>
        <span>{{ item.label }}</span>
      </router-link>
    </nav>

    <div class="sidebar__footer">
      <div class="sidebar__progress-label">
        <span>{{ growth.hasGoal ? '目标达成度' : '课程进度' }}</span>
        <strong>{{ progress }}%</strong>
      </div>
      <div class="sidebar__progress-track">
        <div class="sidebar__progress-fill" :style="{ width: progress + '%' }"></div>
      </div>

      <div v-if="auth.user" class="sidebar__user">
        <span
          class="sidebar__user-avatar"
          :style="{ background: auth.user.avatarColor || '#1772f6', color: '#fff' }"
        >
          {{ auth.user.avatar }}
        </span>
        <div class="sidebar__user-info">
          <strong>
            {{ auth.user.nickname }}
            <em v-if="auth.isGuest" class="sidebar__guest-badge">游客</em>
          </strong>
          <span>{{ auth.isGuest ? '临时体验账号' : '@' + auth.user.username }}</span>
        </div>
        <button class="sidebar__logout" :title="auth.isGuest ? '结束体验' : '退出登录'" @click="logout">
          ⎋
        </button>
      </div>
      <router-link
        v-if="auth.isGuest"
        to="/register"
        class="sidebar__upgrade"
        @click="$emit('navigate')"
      >
        注册账号,保存学习进度 →
      </router-link>
    </div>
  </aside>
</template>

<style scoped>
.sidebar {
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  width: var(--sidebar-width);
  background: var(--surface);
  border-right: 1px solid var(--border);
  display: flex;
  flex-direction: column;
  padding: 20px 14px;
  z-index: 50;
}

.sidebar__brand {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 6px 10px 20px;
  color: var(--text);
}

.sidebar__logo {
  font-size: 26px;
}

.sidebar__title {
  font-size: 17px;
  font-weight: 700;
  letter-spacing: -0.01em;
}

.sidebar__nav {
  display: flex;
  flex-direction: column;
  gap: 4px;
  flex: 1;
  min-height: 0;
  overflow-y: auto;
}

.sidebar__link {
  display: flex;
  align-items: center;
  gap: 11px;
  padding: 10px 12px;
  border-radius: var(--radius-sm);
  color: var(--text-2);
  font-size: 14.5px;
  font-weight: 500;
  transition: all 0.15s ease;
}

.sidebar__link:hover {
  background: var(--surface-2);
  color: var(--text);
}

.sidebar__link--active {
  background: var(--primary-soft);
  color: var(--primary-strong);
  font-weight: 600;
}

.sidebar__icon {
  font-size: 17px;
  width: 22px;
  text-align: center;
}

.sidebar__footer {
  padding: 14px 10px 4px;
  border-top: 1px solid var(--border);
}

.sidebar__progress-label {
  display: flex;
  justify-content: space-between;
  font-size: 12.5px;
  color: var(--text-2);
  margin-bottom: 8px;
}

.sidebar__progress-track {
  height: 7px;
  background: var(--border);
  border-radius: 999px;
  overflow: hidden;
}

.sidebar__progress-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--primary), #a78bfa);
  border-radius: 999px;
  transition: width 0.4s ease;
}

.sidebar__user {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 14px;
  padding: 9px 10px;
  background: var(--surface-2);
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
}

.sidebar__user-avatar {
  width: 34px;
  height: 34px;
  min-width: 34px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 13px;
  font-weight: 700;
  border-radius: 50%;
}

.sidebar__user-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  line-height: 1.35;
}

.sidebar__user-info strong {
  font-size: 13px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  display: flex;
  align-items: center;
  gap: 6px;
}

.sidebar__guest-badge {
  font-style: normal;
  font-size: 10.5px;
  font-weight: 700;
  padding: 1px 6px;
  border-radius: 999px;
  background: #e0f2fe;
  color: #0369a1;
  flex-shrink: 0;
}

.sidebar__user-info span {
  font-size: 11.5px;
  color: var(--text-3);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.sidebar__logout {
  width: 28px;
  height: 28px;
  border: none;
  border-radius: 8px;
  background: transparent;
  color: var(--text-3);
  font-size: 15px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.sidebar__logout:hover {
  background: #fef2f2;
  color: var(--danger);
}

.sidebar__upgrade {
  display: block;
  margin-top: 8px;
  padding: 8px 10px;
  border-radius: 8px;
  background: var(--primary-soft);
  color: var(--primary-strong);
  font-size: 12px;
  font-weight: 600;
  text-align: center;
  transition: background 0.15s ease;
}

.sidebar__upgrade:hover {
  background: #e0e7ff;
}

@media (max-width: 860px) {
  .sidebar {
    transform: translateX(-100%);
    transition: transform 0.25s ease;
  }

  .sidebar--open {
    transform: translateX(0);
  }
}
</style>
