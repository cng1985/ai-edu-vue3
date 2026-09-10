<template>
  <div class="page dashboard">
    <section class="hero">
      <div class="hero__glow hero__glow--a" aria-hidden="true" />
      <div class="hero__glow hero__glow--b" aria-hidden="true" />
      <div class="hero__text">
        <p class="hero__eyebrow">{{ greeting }} · {{ today }}</p>
        <h2>{{ auth.user?.nickname || '管理员' }}，欢迎回到运营看板</h2>
        <p class="hero__lede">
          今日共有 <strong class="num">{{ stats?.reviewStats?.pending || 0 }}</strong> 条内容等待审核，
          平台累计 <strong class="num">{{ stats?.userStats?.total || 0 }}</strong> 位用户与
          <strong class="num">{{ stats?.courseStats?.total || 0 }}</strong> 门课程。
        </p>
        <div class="hero__actions">
          <el-button v-if="auth.hasPermission(PERM.REVIEW_READ)" type="primary" @click="$router.push('/reviews')">
            <el-icon class="el-icon--left"><DocumentChecked /></el-icon>处理审核队列
          </el-button>
          <el-button :icon="Refresh" :loading="loading" @click="loadStats">刷新数据</el-button>
        </div>
      </div>
      <div class="hero__aside">
        <div class="hero__ring">
          <svg viewBox="0 0 120 120">
            <circle cx="60" cy="60" r="52" class="hero__ring-track" />
            <circle
              cx="60" cy="60" r="52"
              class="hero__ring-bar"
              :stroke-dasharray="ringCircumference"
              :stroke-dashoffset="ringOffset"
            />
          </svg>
          <div class="hero__ring-text">
            <strong class="num">{{ publishRate }}%</strong>
            <span>课程发布率</span>
          </div>
        </div>
        <div class="hero__mini">
          <div><span class="muted">已发布</span><strong class="num">{{ stats?.courseStats?.published || 0 }}</strong></div>
          <div><span class="muted">草稿</span><strong class="num">{{ stats?.courseStats?.draft || 0 }}</strong></div>
        </div>
      </div>
    </section>

    <section v-loading="loading" class="stat-grid stat-grid--4">
      <StatCard
        v-for="item in userCards"
        :key="item.label"
        :label="item.label"
        :value="item.value"
        :icon="item.icon"
        :tone="item.tone"
        :hint="item.hint"
      />
    </section>

    <section class="stat-grid stat-grid--4">
      <StatCard
        v-for="item in courseCards"
        :key="item.label"
        :label="item.label"
        :value="item.value"
        :icon="item.icon"
        :tone="item.tone"
        :hint="item.hint"
      />
    </section>

    <section class="dashboard__grid">
      <div class="panel">
        <div class="panel__head">
          <h3 class="panel__title">
            <span class="panel__title-icon"><el-icon><Collection /></el-icon></span>
            内容概览
          </h3>
          <span class="muted">题库与审核状态</span>
        </div>
        <div class="panel__body">
          <div class="overview">
            <div class="overview__item">
              <span class="overview__label">测验套数</span>
              <strong class="overview__value num">{{ stats?.quizStats?.total || 0 }}</strong>
            </div>
            <div class="overview__item">
              <span class="overview__label">题目总数</span>
              <strong class="overview__value num">{{ stats?.quizStats?.questions || 0 }}</strong>
            </div>
            <div class="overview__item overview__item--warning">
              <span class="overview__label">待审核</span>
              <strong class="overview__value num">{{ stats?.reviewStats?.pending || 0 }}</strong>
            </div>
            <div class="overview__item overview__item--success">
              <span class="overview__label">已通过</span>
              <strong class="overview__value num">{{ stats?.reviewStats?.approved || 0 }}</strong>
            </div>
            <div class="overview__item overview__item--danger">
              <span class="overview__label">已驳回</span>
              <strong class="overview__value num">{{ stats?.reviewStats?.rejected || 0 }}</strong>
            </div>
          </div>

          <div class="review-bar">
            <div class="review-bar__label">
              <span>审核分布</span>
              <span class="muted num">{{ reviewTotal }} 条</span>
            </div>
            <div class="review-bar__track">
              <span class="review-bar__seg review-bar__seg--warning" :style="{ width: pct(stats?.reviewStats?.pending) }" />
              <span class="review-bar__seg review-bar__seg--success" :style="{ width: pct(stats?.reviewStats?.approved) }" />
              <span class="review-bar__seg review-bar__seg--danger" :style="{ width: pct(stats?.reviewStats?.rejected) }" />
            </div>
            <div class="review-bar__legend">
              <span><i class="review-bar__dot review-bar__dot--warning" />待审核</span>
              <span><i class="review-bar__dot review-bar__dot--success" />已通过</span>
              <span><i class="review-bar__dot review-bar__dot--danger" />已驳回</span>
            </div>
          </div>
        </div>
      </div>

      <div class="panel">
        <div class="panel__head">
          <h3 class="panel__title">
            <span class="panel__title-icon"><el-icon><Promotion /></el-icon></span>
            快捷操作
          </h3>
          <span class="muted">按权限展示</span>
        </div>
        <div class="panel__body">
          <div class="quick">
            <button
              v-for="q in quickActions"
              :key="q.path"
              type="button"
              class="quick__item"
              :class="`quick__item--${q.tone}`"
              @click="$router.push(q.path)"
            >
              <span class="quick__icon"><el-icon :size="20"><component :is="q.icon" /></el-icon></span>
              <span class="quick__text">
                <strong>{{ q.title }}</strong>
                <span>{{ q.desc }}</span>
              </span>
              <el-icon class="quick__arrow" :size="16"><ArrowRight /></el-icon>
            </button>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import {
  Refresh, User, UserFilled, Avatar, CircleCheck, Reading, Promotion, Document, Notebook,
  Collection, DocumentChecked, EditPen, ArrowRight, Service, Cpu
} from '@element-plus/icons-vue'
import { dashboardApi } from '../api'
import { useAuthStore } from '../stores/auth'
import { PERM } from '../constants/permissions'
import StatCard from '../components/common/StatCard.vue'

const auth = useAuthStore()
const loading = ref(false)
const stats = ref(null)

const hour = new Date().getHours()
const greeting = hour < 6 ? '夜深了' : hour < 12 ? '早上好' : hour < 18 ? '下午好' : '晚上好'
const today = new Date().toLocaleDateString('zh-CN', { month: 'long', day: 'numeric', weekday: 'long' })

const userCards = computed(() => [
  { label: '用户总数', value: stats.value?.userStats?.total || 0, icon: User, tone: 'primary', hint: '含全部角色' },
  { label: '学员数', value: stats.value?.userStats?.learners || 0, icon: UserFilled, tone: 'success', hint: '学习端注册用户' },
  { label: '管理 / 运营', value: stats.value?.userStats?.admins || 0, icon: Avatar, tone: 'info', hint: '后台账号' },
  { label: '活跃账号', value: stats.value?.userStats?.active || 0, icon: CircleCheck, tone: 'sky', hint: '状态为正常' }
])

const courseCards = computed(() => [
  { label: '课程总数', value: stats.value?.courseStats?.total || 0, icon: Reading, tone: 'primary', hint: '全部课程' },
  { label: '已发布', value: stats.value?.courseStats?.published || 0, icon: Promotion, tone: 'success', hint: '学员可见' },
  { label: '草稿', value: stats.value?.courseStats?.draft || 0, icon: Document, tone: 'warning', hint: '尚未发布' },
  { label: '章节总数', value: stats.value?.courseStats?.chapters || 0, icon: Notebook, tone: 'rose', hint: '所有课程章节' }
])

const publishRate = computed(() => {
  const total = stats.value?.courseStats?.total || 0
  const published = stats.value?.courseStats?.published || 0
  return total ? Math.round((published / total) * 100) : 0
})

const ringCircumference = 2 * Math.PI * 52
const ringOffset = computed(() => ringCircumference * (1 - publishRate.value / 100))

const reviewTotal = computed(() => {
  const r = stats.value?.reviewStats || {}
  return (r.pending || 0) + (r.approved || 0) + (r.rejected || 0)
})

function pct(v) {
  if (!reviewTotal.value) return '0%'
  return `${((v || 0) / reviewTotal.value) * 100}%`
}

const quickActions = computed(() => [
  auth.hasPermission(PERM.COURSE_READ) && { path: '/courses', title: '管理课程', desc: '新增、编辑与发布课程', icon: Reading, tone: 'primary' },
  auth.hasPermission(PERM.QUIZ_READ) && { path: '/quizzes', title: '管理题库', desc: '维护测验与题目', icon: EditPen, tone: 'sky' },
  auth.hasPermission(PERM.REVIEW_READ) && { path: '/reviews', title: '审核队列', desc: '处理待审核内容', icon: DocumentChecked, tone: 'warning' },
  auth.hasPermission(PERM.USER_READ) && { path: '/users', title: '用户管理', desc: '账号、角色与状态', icon: User, tone: 'success' },
  auth.hasPermission(PERM.CUSTOMER_READ) && { path: '/customers', title: '客户咨询', desc: '实时回复学员工单', icon: Service, tone: 'rose' },
  auth.hasPermission(PERM.AI_MODEL_READ) && { path: '/ai-models', title: '大模型配置', desc: '路由与厂商密钥', icon: Cpu, tone: 'ink' }
].filter(Boolean))

async function loadStats() {
  loading.value = true
  try {
    stats.value = await dashboardApi.stats()
  } finally {
    loading.value = false
  }
}

onMounted(loadStats)
</script>

<style scoped>
.hero {
  position: relative;
  display: grid;
  grid-template-columns: minmax(0, 1fr) auto;
  gap: 32px;
  align-items: center;
  padding: 32px 36px;
  border-radius: var(--radius-lg);
  background: linear-gradient(135deg, var(--ink-2) 0%, var(--ink) 70%);
  color: var(--ink-text);
  overflow: hidden;
  box-shadow: var(--shadow-lg);
}

.hero__glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(50px);
  pointer-events: none;
}

.hero__glow--a {
  width: 420px;
  height: 420px;
  top: -200px;
  left: -80px;
  background: rgba(107, 92, 255, 0.55);
}

.hero__glow--b {
  width: 320px;
  height: 320px;
  bottom: -180px;
  right: 10%;
  background: rgba(15, 185, 129, 0.3);
}

.hero__text {
  position: relative;
  min-width: 0;
}

.hero__eyebrow {
  margin: 0 0 8px;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  color: #a396ff;
}

.hero h2 {
  margin: 0;
  color: #fff;
  font-size: 26px;
}

.hero__lede {
  margin: 12px 0 0;
  font-size: 14.5px;
  line-height: 1.8;
  color: var(--ink-text);
  opacity: 0.9;
}

.hero__lede strong {
  color: #fff;
  font-size: 16px;
  padding: 0 2px;
}

.hero__actions {
  display: flex;
  gap: 10px;
  margin-top: 22px;
  flex-wrap: wrap;
}

.hero__actions :deep(.el-button--default) {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.14);
  color: #fff;
}

.hero__actions :deep(.el-button--default:hover) {
  background: rgba(255, 255, 255, 0.16);
  border-color: rgba(255, 255, 255, 0.25);
  color: #fff;
}

.hero__aside {
  position: relative;
  display: flex;
  align-items: center;
  gap: 22px;
  padding: 18px 22px;
  border-radius: var(--radius);
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.08);
  backdrop-filter: blur(10px);
}

.hero__ring {
  position: relative;
  width: 124px;
  height: 124px;
}

.hero__ring svg {
  width: 100%;
  height: 100%;
  transform: rotate(-90deg);
}

.hero__ring-track {
  fill: none;
  stroke: rgba(255, 255, 255, 0.1);
  stroke-width: 10;
}

.hero__ring-bar {
  fill: none;
  stroke: #9a8cff;
  stroke-width: 10;
  stroke-linecap: round;
  transition: stroke-dashoffset 0.8s var(--ease);
}

.hero__ring-text {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 2px;
}

.hero__ring-text strong {
  font-size: 24px;
  color: #fff;
  letter-spacing: -0.02em;
}

.hero__ring-text span {
  font-size: 11px;
  color: var(--ink-text-2);
}

.hero__mini {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.hero__mini div {
  display: flex;
  flex-direction: column;
  line-height: 1.2;
}

.hero__mini .muted {
  font-size: 12px;
  color: var(--ink-text-2);
}

.hero__mini strong {
  font-size: 22px;
  color: #fff;
}

.stat-grid--4 {
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.dashboard__grid {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 20px;
}

.overview {
  display: grid;
  grid-template-columns: repeat(5, minmax(0, 1fr));
  gap: 10px;
}

.overview__item {
  --tone: var(--primary);
  --tone-soft: var(--primary-soft);
  display: flex;
  flex-direction: column;
  gap: 6px;
  padding: 14px 12px;
  border-radius: var(--radius-sm);
  background: var(--tone-soft);
  text-align: center;
}

.overview__item--warning { --tone: var(--warning-strong); --tone-soft: var(--warning-soft); }
.overview__item--success { --tone: var(--success-strong); --tone-soft: var(--success-soft); }
.overview__item--danger { --tone: var(--danger-strong); --tone-soft: var(--danger-soft); }

.overview__label {
  font-size: 12px;
  color: var(--text-2);
  font-weight: 600;
}

.overview__value {
  font-size: 22px;
  color: var(--tone);
  letter-spacing: -0.02em;
}

.review-bar {
  margin-top: 22px;
}

.review-bar__label {
  display: flex;
  justify-content: space-between;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-2);
  margin-bottom: 10px;
}

.review-bar__track {
  display: flex;
  height: 12px;
  border-radius: 999px;
  background: var(--surface-3);
  overflow: hidden;
  gap: 2px;
}

.review-bar__seg {
  display: block;
  height: 100%;
  transition: width 0.6s var(--ease);
}

.review-bar__seg--warning { background: var(--warning); }
.review-bar__seg--success { background: var(--success); }
.review-bar__seg--danger { background: var(--danger); }

.review-bar__legend {
  display: flex;
  gap: 16px;
  margin-top: 10px;
  font-size: 12px;
  color: var(--text-3);
}

.review-bar__legend span {
  display: inline-flex;
  align-items: center;
  gap: 6px;
}

.review-bar__dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.review-bar__dot--warning { background: var(--warning); }
.review-bar__dot--success { background: var(--success); }
.review-bar__dot--danger { background: var(--danger); }

.quick {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 10px;
}

.quick__item {
  --tone: var(--primary);
  --tone-soft: var(--primary-soft);
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 14px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  background: var(--surface);
  text-align: left;
  cursor: pointer;
  transition: all var(--t-fast);
}

.quick__item--sky { --tone: var(--sky); --tone-soft: var(--sky-soft); }
.quick__item--warning { --tone: var(--warning); --tone-soft: var(--warning-soft); }
.quick__item--success { --tone: var(--success); --tone-soft: var(--success-soft); }
.quick__item--rose { --tone: var(--rose); --tone-soft: var(--rose-soft); }
.quick__item--ink { --tone: var(--ink-3); --tone-soft: var(--surface-3); }

.quick__item:hover {
  border-color: var(--tone);
  transform: translateY(-2px);
  box-shadow: var(--shadow);
}

.quick__icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 42px;
  height: 42px;
  flex-shrink: 0;
  border-radius: 12px;
  background: var(--tone-soft);
  color: var(--tone);
}

.quick__text {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-width: 0;
  line-height: 1.3;
}

.quick__text strong {
  font-size: 14px;
  color: var(--text);
}

.quick__text span {
  font-size: 12px;
  color: var(--text-3);
}

.quick__arrow {
  color: var(--text-3);
  opacity: 0;
  transform: translateX(-4px);
  transition: all var(--t-fast);
}

.quick__item:hover .quick__arrow {
  opacity: 1;
  transform: translateX(0);
  color: var(--tone);
}

@media (max-width: 1100px) {
  .stat-grid--4 {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .dashboard__grid {
    grid-template-columns: 1fr;
  }

  .overview {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}

@media (max-width: 760px) {
  .hero {
    grid-template-columns: 1fr;
    padding: 24px 22px;
  }

  .hero__aside {
    justify-content: center;
  }

  .stat-grid--4 {
    grid-template-columns: 1fr;
  }

  .quick {
    grid-template-columns: 1fr;
  }
}
</style>
