<script setup>
import { computed } from 'vue'
import { useRouter } from 'vue-router'
import { useLearningStore } from '../stores/learning'
import { useAuthStore } from '../stores/auth'
import { useGrowthStore } from '../stores/growth'
import { PERM } from '../constants/permissions'
import Icon from './Icon.vue'

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

const navGroups = [
  {
    title: '目标成长',
    items: [
      { to: '/', icon: 'home', label: '学习驾驶舱', exact: true },
      { to: '/career', icon: 'compass', label: '职业与目标' },
      { to: '/path', icon: 'map', label: '学习路径' },
      { to: '/review', icon: 'refresh', label: '复习与补强' }
    ]
  },
  {
    title: '课程与练习',
    items: [
      { to: '/courses', icon: 'book', label: '全部课程', permission: PERM.COURSE_READ },
      { to: '/quiz', icon: 'clipboard', label: '知识测验', permission: PERM.QUIZ_READ },
      { to: '/chat', icon: 'sparkles', label: 'AI 学习助手', permission: PERM.AI_CHAT },
      { to: '/support', icon: 'headset', label: '客户咨询', permission: PERM.CUSTOMER_CHAT }
    ]
  },
  {
    title: '数据与激励',
    items: [
      { to: '/stats', icon: 'chart', label: '达成度报告' },
      { to: '/incentives', icon: 'medal', label: '成长激励' }
    ]
  }
]

const visibleGroups = computed(() =>
  navGroups
    .map((group) => ({
      ...group,
      items: group.items.filter((item) => !item.permission || auth.hasPermission(item.permission))
    }))
    .filter((group) => group.items.length)
)

const progress = computed(() => (growth.hasGoal ? growth.achievement : learning.overallProgress))
</script>

<template>
  <aside class="sidebar" :class="{ 'sidebar--open': open }">
    <router-link to="/" class="sidebar__brand" @click="$emit('navigate')">
      <span class="sidebar__logo">
        <Icon name="brain" :size="20" :stroke="2.2" />
      </span>
      <span class="sidebar__brand-text">
        <strong>AI 学习系统</strong>
        <small>目标驱动 · 数据评估</small>
      </span>
    </router-link>

    <nav class="sidebar__nav">
      <div v-for="group in visibleGroups" :key="group.title" class="sidebar__group">
        <span class="sidebar__group-title">{{ group.title }}</span>
        <router-link
          v-for="item in group.items"
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
          <Icon :name="item.icon" :size="18" />
          <span>{{ item.label }}</span>
        </router-link>
      </div>
    </nav>

    <div class="sidebar__footer">
      <div class="sidebar__progress">
        <div class="sidebar__progress-label">
          <span>{{ growth.hasGoal ? '目标达成度' : '课程进度' }}</span>
          <strong class="num">{{ progress }}%</strong>
        </div>
        <div class="sidebar__progress-track">
          <div class="sidebar__progress-fill" :style="{ width: progress + '%' }"></div>
        </div>
      </div>

      <div v-if="auth.user" class="sidebar__user">
        <span
          class="sidebar__user-avatar"
          :style="{ background: auth.user.avatarColor || '#6b5cff' }"
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
          <Icon name="logout" :size="16" />
        </button>
      </div>
      <router-link
        v-if="auth.isGuest"
        to="/register"
        class="sidebar__upgrade"
        @click="$emit('navigate')"
      >
        注册账号，保存学习进度
        <Icon name="arrowRight" :size="14" />
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
  display: flex;
  flex-direction: column;
  padding: 22px 16px 18px;
  z-index: 50;
  color: var(--ink-text);
  background:
    radial-gradient(90% 40% at 100% 0%, rgba(107, 92, 255, 0.28), transparent 70%),
    radial-gradient(70% 30% at 0% 100%, rgba(15, 185, 129, 0.16), transparent 70%),
    linear-gradient(180deg, var(--ink-2) 0%, var(--ink) 100%);
  box-shadow: 1px 0 0 rgba(255, 255, 255, 0.04);
}

.sidebar__brand {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 4px 8px 22px;
  color: #fff;
}

.sidebar__logo {
  display: grid;
  place-items: center;
  width: 40px;
  height: 40px;
  border-radius: 13px;
  color: #fff;
  background: linear-gradient(135deg, var(--primary), #a071ff);
  box-shadow: 0 8px 20px var(--primary-glow), inset 0 1px 0 rgba(255, 255, 255, 0.35);
}

.sidebar__brand-text {
  display: flex;
  flex-direction: column;
  line-height: 1.25;
}

.sidebar__brand-text strong {
  font-family: 'Sora', 'Noto Sans SC', sans-serif;
  font-size: 16px;
  font-weight: 700;
  letter-spacing: -0.01em;
}

.sidebar__brand-text small {
  font-size: 11px;
  color: rgba(201, 203, 230, 0.6);
  letter-spacing: 0.04em;
}

.sidebar__nav {
  flex: 1;
  min-height: 0;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 18px;
  padding-right: 2px;
}

.sidebar__group {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.sidebar__group-title {
  padding: 0 12px 8px;
  font-size: 11px;
  font-weight: 700;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: rgba(201, 203, 230, 0.45);
}

.sidebar__link {
  position: relative;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 10px 12px;
  border-radius: 12px;
  color: rgba(201, 203, 230, 0.82);
  font-size: 14px;
  font-weight: 500;
  transition: background var(--t-fast), color var(--t-fast), transform var(--t-fast);
}

.sidebar__link .icon {
  color: rgba(201, 203, 230, 0.6);
  transition: color var(--t-fast);
}

.sidebar__link:hover {
  background: rgba(255, 255, 255, 0.06);
  color: #fff;
}

.sidebar__link:hover .icon {
  color: #fff;
}

.sidebar__link--active {
  background: linear-gradient(90deg, rgba(107, 92, 255, 0.34), rgba(107, 92, 255, 0.12));
  color: #fff;
  font-weight: 600;
  box-shadow: inset 0 0 0 1px rgba(160, 113, 255, 0.28);
}

.sidebar__link--active::before {
  content: '';
  position: absolute;
  left: -16px;
  top: 10px;
  bottom: 10px;
  width: 3px;
  border-radius: 0 3px 3px 0;
  background: linear-gradient(180deg, #a071ff, var(--primary));
}

.sidebar__link--active .icon {
  color: #c7bfff;
}

.sidebar__footer {
  padding-top: 16px;
  margin-top: 8px;
  border-top: 1px solid rgba(255, 255, 255, 0.08);
}

.sidebar__progress {
  padding: 0 4px;
}

.sidebar__progress-label {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  font-size: 12px;
  color: rgba(201, 203, 230, 0.7);
  margin-bottom: 8px;
}

.sidebar__progress-label strong {
  color: #fff;
  font-size: 13px;
}

.sidebar__progress-track {
  height: 6px;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 999px;
  overflow: hidden;
}

.sidebar__progress-fill {
  height: 100%;
  background: linear-gradient(90deg, var(--primary), #a071ff 60%, #6ee7b7);
  border-radius: 999px;
  transition: width 0.5s var(--ease);
  box-shadow: 0 0 12px var(--primary-glow);
}

.sidebar__user {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 16px;
  padding: 10px;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.07);
  border-radius: 14px;
}

.sidebar__user-avatar {
  width: 36px;
  height: 36px;
  min-width: 36px;
  display: grid;
  place-items: center;
  font-size: 13px;
  font-weight: 700;
  color: #fff;
  border-radius: 12px;
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.3);
}

.sidebar__user-info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  line-height: 1.35;
}

.sidebar__user-info strong {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: #fff;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.sidebar__guest-badge {
  font-style: normal;
  font-size: 10px;
  font-weight: 700;
  padding: 1px 6px;
  border-radius: 999px;
  background: rgba(42, 140, 244, 0.25);
  color: #9ec9ff;
  flex-shrink: 0;
}

.sidebar__user-info span {
  font-size: 11.5px;
  color: rgba(201, 203, 230, 0.55);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.sidebar__logout {
  display: grid;
  place-items: center;
  width: 30px;
  height: 30px;
  border: none;
  border-radius: 9px;
  background: transparent;
  color: rgba(201, 203, 230, 0.6);
  cursor: pointer;
  transition: all var(--t-fast);
}

.sidebar__logout:hover {
  background: rgba(239, 77, 99, 0.2);
  color: #ff9aa9;
}

.sidebar__upgrade {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  margin-top: 10px;
  padding: 9px 12px;
  border-radius: 11px;
  background: linear-gradient(135deg, var(--primary), #8f7bff);
  color: #fff;
  font-size: 12.5px;
  font-weight: 600;
  box-shadow: 0 6px 16px var(--primary-glow);
  transition: transform var(--t-fast), box-shadow var(--t-fast);
}

.sidebar__upgrade:hover {
  transform: translateY(-1px);
  box-shadow: 0 10px 22px var(--primary-glow);
}

@media (max-width: 860px) {
  .sidebar {
    transform: translateX(-100%);
    transition: transform 0.28s var(--ease);
    box-shadow: none;
  }

  .sidebar--open {
    transform: translateX(0);
    box-shadow: 24px 0 64px rgba(18, 20, 48, 0.4);
  }
}
</style>
