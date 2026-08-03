<script setup>
import { computed, onMounted, ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useGrowthStore } from '../stores/growth'
import { useLearningStore } from '../stores/learning'
import { communityPosts } from '../data/community'
import { aiApi } from '../api'
import ProgressRing from '../components/ProgressRing.vue'

const auth = useAuthStore()
const growth = useGrowthStore()
const learning = useLearningStore()
const aiSuggestions = ref([])

const greeting = computed(() => {
  const hour = new Date().getHours()
  return hour < 12 ? '早上好' : hour < 18 ? '下午好' : '晚上好'
})
const daysLeft = computed(() => {
  if (!growth.goal) return 0
  return Math.max(0, Math.ceil((new Date(growth.goal.deadline) - new Date()) / 86400000))
})
const weakDomain = computed(() => [...growth.competencyProgress].sort((a, b) => a.progress - b.progress)[0])
const todayTasks = computed(() => {
  if (!growth.hasGoal) return []
  const tasks = []
  if (growth.dueReviewCount) {
    tasks.push({ title: `复习 ${growth.dueReviewCount} 个到期知识点`, meta: '遗忘曲线复习 · 约 5 分钟', to: '/review' })
  }
  if (growth.nextUnit) {
    tasks.push({ title: growth.nextUnit.title, meta: `微单元 · ${growth.nextUnit.duration} 分钟`, to: `/micro/${growth.nextUnit.id}` })
  }
  tasks.push({ title: '完成今日学习打卡', meta: '打卡获得 10 成长积分', action: 'checkin', done: growth.checkedInToday })
  return tasks.slice(0, 3)
})

onMounted(async () => {
  if (!growth.hasGoal) return
  try {
    const result = await aiApi.learningSuggest({
      competencyProgress: growth.competencyProgress.map((c) => ({ name: c.name, progress: c.progress })),
      nextMilestone: growth.nextMilestone?.name || '',
      streak: growth.streak
    })
    aiSuggestions.value = result.suggestions || []
  } catch {
    aiSuggestions.value = []
  }
})
</script>

<template>
  <div class="page dashboard">
    <section v-if="!growth.hasGoal" class="onboarding">
      <div class="onboarding__glow"></div>
      <div class="onboarding__copy">
        <span>目标驱动 · AI 辅助 · 数据评估</span>
        <h1>{{ greeting }}，{{ auth.user?.nickname || '学习者' }}<br />先确定你想成为的人</h1>
        <p>完成 3 步职业规划，系统会根据你的基础和时间生成能力图谱、学习任务与里程碑。</p>
        <router-link to="/career" class="hero-button">✨ 开始 AI 职业规划</router-link>
      </div>
      <div class="loop">
        <div v-for="(item, index) in ['目标确认', '路径分解', '学习执行', '达成评估']" :key="item">
          <b>{{ index + 1 }}</b><span>{{ item }}</span>
        </div>
      </div>
    </section>

    <template v-else>
      <section class="welcome">
        <div>
          <span class="welcome__date">今日学习驾驶舱</span>
          <h1>{{ greeting }}，{{ auth.user?.nickname }}！</h1>
          <p>连续学习 <strong>{{ growth.streak }}</strong> 天 · 今日还有 {{ todayTasks.filter(item => !item.done).length }} 项任务</p>
        </div>
        <router-link :to="growth.nextUnit ? `/micro/${growth.nextUnit.id}` : '/path'" class="btn btn--primary">
          ▶ {{ growth.nextUnit ? '开始今日学习' : '查看完成路径' }}
        </router-link>
      </section>

      <section class="goal-card">
        <div class="goal-card__content">
          <span class="goal-tag">当前职业目标 · 进行中</span>
          <h2>{{ growth.goal.name }}</h2>
          <p>剩余 {{ daysLeft }} 天 · 每周 {{ growth.goal.weeklyHours }} 小时 · {{ growth.goal.baseLevel }}起步</p>
          <div class="goal-track"><div :style="{ width: growth.achievement + '%' }"></div></div>
          <div class="goal-labels"><strong>总达成度 {{ growth.achievement }}%</strong><span>达标线 75%</span></div>
          <div v-if="growth.nextMilestone" class="milestone">
            <span>🏁 下一里程碑</span>
            <strong>{{ growth.nextMilestone.name }}</strong>
            <small>第 {{ growth.nextMilestone.week }} 周 · {{ growth.nextMilestone.standard }}</small>
          </div>
        </div>
        <div class="goal-card__score">
          <ProgressRing :percent="growth.achievement" :size="132" :stroke="10" color="#fff" />
          <router-link to="/stats">查看达成报告 →</router-link>
        </div>
      </section>

      <section class="quick-stats">
        <router-link to="/path" class="quick-stat card">
          <span>⚡</span><div><strong>{{ growth.completedUnitCount }}</strong><small>已完成微单元</small></div>
        </router-link>
        <router-link to="/review" class="quick-stat card">
          <span>🔍</span><div><strong>{{ growth.dueReviewCount }}</strong><small>今日待复习</small></div>
        </router-link>
        <router-link to="/incentives" class="quick-stat card">
          <span>🏅</span><div><strong>{{ growth.points }}</strong><small>成长积分 · Lv{{ growth.level.number }}</small></div>
        </router-link>
        <router-link to="/courses" class="quick-stat card">
          <span>📚</span><div><strong>{{ learning.completedCount }}</strong><small>已完成课程章节</small></div>
        </router-link>
      </section>

      <div class="dashboard-grid">
        <section class="card panel">
          <div class="panel-head"><div><h2>能力域进度</h2><p>按知识点掌握度加权聚合</p></div><router-link to="/path">完整路径 →</router-link></div>
          <div v-for="domain in growth.competencyProgress" :key="domain.id" class="competency">
            <div><strong>{{ domain.name }}</strong><span>{{ domain.progress }}%</span></div>
            <div class="bar"><i :style="{ width: domain.progress + '%', background: domain.color }"></i></div>
          </div>
        </section>

        <section class="card panel">
          <div class="panel-head"><div><h2>今日待办</h2><p>根据路径、薄弱点和复习日程生成</p></div><span>{{ todayTasks.length }} 项</span></div>
          <div class="task-list">
            <component
              :is="task.to ? 'router-link' : 'button'"
              v-for="(task, index) in todayTasks"
              :key="task.title"
              :to="task.to"
              class="task"
              :class="{ done: task.done }"
              @click="task.action === 'checkin' && growth.checkIn()"
            >
              <b>{{ task.done ? '✓' : index + 1 }}</b>
              <span><strong>{{ task.title }}</strong><small>{{ task.meta }}</small></span>
              <i>{{ task.done ? '已完成' : '›' }}</i>
            </component>
          </div>
        </section>

        <section class="card panel ai-panel">
          <div class="panel-head"><div><h2>✨ AI 学习建议</h2><p>基于当前达成数据</p></div><router-link to="/chat">问 AI →</router-link></div>
          <ul v-if="aiSuggestions.length" class="ai-suggest-list">
            <li v-for="item in aiSuggestions" :key="item">{{ item }}</li>
          </ul>
          <template v-else>
            <p v-if="growth.achievement === 0">从第一个 5 分钟微单元开始。一次掌握一个小知识点，通过快测后再进入下一步。</p>
            <p v-else>当前最需要关注的是 <strong>{{ weakDomain?.name }}</strong>（{{ weakDomain?.progress }}%）。优先完成该能力域任务，并按间隔计划复习。</p>
          </template>
          <div><span>基于达成度</span><span>关键路径优先</span><span>每日可完成</span></div>
        </section>

        <section class="card panel growth-panel">
          <div class="panel-head"><div><h2>成长激励</h2><p>{{ growth.level.name }}</p></div><router-link to="/incentives">勋章墙 →</router-link></div>
          <div class="points"><strong>{{ growth.points }}</strong><span>累计积分</span><b>Lv{{ growth.level.number }}</b></div>
          <div class="level-track"><i :style="{ width: Math.min(100, growth.points / growth.level.next * 100) + '%' }"></i></div>
          <div class="badges">
            <span v-for="badge in growth.badges.slice(0, 3)" :key="badge">🏅 {{ badge }}</span>
            <span v-if="!growth.badges.length" class="muted">完成首个任务解锁勋章</span>
          </div>
        </section>
      </div>
    </template>

    <section class="community card">
      <div class="panel-head">
        <div><h2>同路人正在分享</h2><p>学习不孤单，看看社区里的实践与复盘</p></div>
        <span>学习动态</span>
      </div>
      <div class="post-grid">
        <article v-for="post in communityPosts.slice(0, 3)" :key="post.id">
          <div class="post-author"><span :style="{ background: post.avatarColor }">{{ post.avatar }}</span><b>{{ post.author }}</b><small>{{ post.time }}</small></div>
          <h3>{{ post.title }}</h3>
          <p>{{ post.excerpt }}</p>
          <div><span v-for="tag in post.tags" :key="tag">#{{ tag }}</span><small>♡ {{ post.likes }} · 💬 {{ post.comments }}</small></div>
        </article>
      </div>
    </section>
  </div>
</template>

<style scoped>
.dashboard { max-width: 1120px; }
.onboarding { position: relative; min-height: 390px; padding: 50px; overflow: hidden; border-radius: 22px; color: white; background: #11132b; }
.onboarding__glow { position: absolute; inset: -40% -10% auto 40%; height: 500px; background: radial-gradient(circle, rgba(124,58,237,.7), transparent 62%); }
.onboarding__copy { position: relative; z-index: 1; max-width: 650px; }
.onboarding__copy > span, .welcome__date { color: #a5b4fc; font-size: 12px; font-weight: 800; letter-spacing: .08em; }
.onboarding h1 { margin: 12px 0; font-size: 38px; line-height: 1.25; }
.onboarding p { max-width: 600px; color: #cbd5e1; font-size: 16px; }
.hero-button { display: inline-flex; margin-top: 8px; padding: 12px 20px; border-radius: 10px; background: white; color: #312e81; font-weight: 800; }
.loop { position: relative; z-index: 1; display: grid; grid-template-columns: repeat(4, 1fr); gap: 10px; margin-top: 50px; }
.loop div { display: flex; align-items: center; gap: 9px; padding: 12px; border: 1px solid rgba(255,255,255,.13); border-radius: 10px; background: rgba(255,255,255,.06); }
.loop b { display: grid; place-items: center; width: 27px; height: 27px; border-radius: 8px; background: #6366f1; }
.loop span { font-size: 12px; }
.welcome { display: flex; align-items: center; justify-content: space-between; margin-bottom: 17px; }
.welcome h1 { margin: 2px 0; font-size: 27px; }
.welcome p { margin: 0; color: var(--text-2); }
.welcome strong { color: var(--warning); }
.goal-card { display: flex; margin-bottom: 16px; overflow: hidden; border-radius: 17px; color: white; background: linear-gradient(120deg, #312e81, #4f46e5 62%, #7c3aed); box-shadow: 0 13px 35px rgba(49,46,129,.2); }
.goal-card__content { flex: 1; padding: 27px 31px; }
.goal-tag { padding: 3px 9px; border-radius: 99px; background: rgba(255,255,255,.13); color: #c7d2fe; font-size: 11px; }
.goal-card h2 { margin: 8px 0 2px; font-size: 22px; }
.goal-card p { margin: 0 0 13px; color: #c7d2fe; font-size: 12px; }
.goal-track { height: 7px; overflow: hidden; border-radius: 8px; background: rgba(255,255,255,.2); }
.goal-track div { height: 100%; background: linear-gradient(90deg, #34d399, #a7f3d0); }
.goal-labels { display: flex; justify-content: space-between; margin-top: 4px; color: #c7d2fe; font-size: 11px; }
.milestone { display: grid; grid-template-columns: auto 1fr; gap: 0 10px; margin-top: 15px; padding: 9px 12px; border-radius: 9px; background: rgba(255,255,255,.09); }
.milestone span { grid-row: span 2; align-self: center; color: #c7d2fe; font-size: 11px; }
.milestone small { color: #c7d2fe; }
.goal-card__score { display: flex; width: 205px; flex-direction: column; align-items: center; justify-content: center; background: rgba(255,255,255,.05); }
.goal-card__score :deep(.progress-ring__text) { fill: white; font-size: 17px; }
.goal-card__score a { color: #e0e7ff; font-size: 11px; }
.quick-stats { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; margin-bottom: 16px; }
.quick-stat { display: flex; align-items: center; gap: 11px; padding: 14px 16px; color: inherit; box-shadow: none; }
.quick-stat > span { font-size: 22px; }
.quick-stat div { display: flex; flex-direction: column; }
.quick-stat strong { font-size: 19px; line-height: 1.2; }
.quick-stat small { color: var(--text-3); }
.dashboard-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 16px; }
.panel, .community { padding: 21px 23px; box-shadow: none; }
.panel-head { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 14px; }
.panel-head h2 { margin: 0; font-size: 16px; }
.panel-head p { margin: 2px 0 0; color: var(--text-3); font-size: 10px; }
.panel-head > a, .panel-head > span { color: var(--primary); font-size: 11px; }
.competency { margin: 13px 0; }
.competency > div:first-child { display: flex; justify-content: space-between; margin-bottom: 4px; font-size: 12px; }
.competency span { color: var(--text-2); }
.bar, .level-track { height: 6px; overflow: hidden; border-radius: 6px; background: var(--border); }
.bar i, .level-track i { display: block; height: 100%; border-radius: inherit; }
.task-list { display: flex; flex-direction: column; gap: 7px; }
.task { display: flex; align-items: center; gap: 10px; width: 100%; padding: 9px; border: none; border-radius: 9px; background: var(--surface-2); color: inherit; font: inherit; text-align: left; cursor: pointer; }
.task:hover { background: var(--primary-soft); }
.task.done { opacity: .6; }
.task > b { display: grid; place-items: center; width: 27px; height: 27px; border: 1px solid var(--border); border-radius: 8px; background: white; color: var(--primary); }
.task > span { display: flex; flex: 1; flex-direction: column; }
.task strong { font-size: 12px; }
.task small { color: var(--text-3); }
.task > i { color: var(--text-3); font-size: 11px; font-style: normal; }
.ai-panel { background: linear-gradient(140deg, #faf5ff, white); }
.ai-panel > p { color: var(--text-2); font-size: 13px; }
.ai-suggest-list { margin: 0 0 12px; padding-left: 18px; color: var(--text-2); font-size: 13px; line-height: 1.7; }
.ai-panel > div:last-child { display: flex; gap: 6px; flex-wrap: wrap; }
.ai-panel > div:last-child span { padding: 3px 7px; border: 1px solid #e9d5ff; border-radius: 99px; color: #7e22ce; font-size: 9px; }
.points { display: flex; align-items: baseline; gap: 8px; margin-bottom: 7px; }
.points strong { color: var(--warning); font-size: 28px; }
.points span { color: var(--text-3); font-size: 11px; }
.points b { margin-left: auto; color: var(--primary); }
.growth-panel .level-track i { background: linear-gradient(90deg, #f59e0b, #fbbf24); }
.badges { display: flex; gap: 6px; margin-top: 13px; flex-wrap: wrap; }
.badges span { padding: 4px 7px; border-radius: 7px; background: #fffbeb; color: #92400e; font-size: 10px; }
.badges .muted { color: var(--text-3); background: var(--surface-2); }
.community { margin-top: 16px; }
.post-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px; }
.post-grid article { padding: 15px; border: 1px solid var(--border); border-radius: 11px; background: var(--surface-2); }
.post-author { display: flex; align-items: center; gap: 7px; }
.post-author > span { display: grid; place-items: center; width: 27px; height: 27px; border-radius: 50%; color: white; font-size: 10px; }
.post-author b { font-size: 11px; }
.post-author small { margin-left: auto; color: var(--text-3); font-size: 9px; }
.post-grid h3 { margin: 10px 0 5px; font-size: 13px; }
.post-grid p { display: -webkit-box; margin: 0 0 10px; overflow: hidden; color: var(--text-2); font-size: 10px; -webkit-box-orient: vertical; -webkit-line-clamp: 2; }
.post-grid article > div:last-child { display: flex; gap: 5px; color: var(--primary); font-size: 9px; }
.post-grid article > div:last-child small { margin-left: auto; color: var(--text-3); }
@media (max-width: 800px) {
  .onboarding { padding: 32px 24px; }
  .onboarding h1 { font-size: 29px; }
  .loop, .quick-stats { grid-template-columns: 1fr 1fr; }
  .goal-card { flex-direction: column; }
  .goal-card__score { width: 100%; padding: 18px; }
  .dashboard-grid { grid-template-columns: 1fr; }
  .post-grid { grid-template-columns: 1fr; }
}
@media (max-width: 480px) {
  .welcome { align-items: flex-start; gap: 10px; }
  .loop, .quick-stats { grid-template-columns: 1fr; }
}
</style>
