<script setup>
import { computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useGrowthStore } from '../stores/growth'
import ProgressRing from '../components/ProgressRing.vue'

const auth = useAuthStore()
const growth = useGrowthStore()

const greeting = computed(() => {
  const hour = new Date().getHours()
  return hour < 12 ? '早上好' : hour < 18 ? '下午好' : '晚上好'
})
const daysLeft = computed(() => {
  if (!growth.goal) return 0
  return Math.max(0, Math.ceil((new Date(growth.goal.deadline) - new Date()) / 86400000))
})
const todayTasks = computed(() => {
  if (!growth.hasGoal) return []
  const tasks = []
  if (growth.nextUnit) tasks.push({
    title: growth.nextUnit.title,
    meta: `微单元 · ${growth.nextUnit.duration} 分钟`,
    to: `/micro/${growth.nextUnit.id}`,
    done: false
  })
  tasks.push({ title: '完成今日学习打卡', meta: '连续学习可获得额外奖励', action: 'checkin', done: growth.checkedInToday })
  tasks.push({ title: '向 AI 助手复述今日要点', meta: '用输出验证理解 · 5 分钟', to: '/chat', done: false })
  return tasks
})
const weakDomain = computed(() => [...growth.competencyProgress].sort((a, b) => a.progress - b.progress)[0])
</script>

<template>
  <div class="page dashboard">
    <section v-if="!growth.hasGoal" class="onboarding card">
      <div class="onboarding__copy">
        <span class="eyebrow">定目标 · 拆路径 · 学得准 · 达得成</span>
        <h1>{{ greeting }}，{{ auth.user?.nickname || '学习者' }}！先确定你的职业目标</h1>
        <p>完成 3 步 AI 职业规划，系统会根据你的基础和时间生成能力图谱、学习任务与里程碑。</p>
        <router-link to="/career" class="btn btn--primary">✨ 开始 AI 职业规划</router-link>
      </div>
      <div class="loop">
        <div><b>1</b><span>目标确认</span></div><i>→</i>
        <div><b>2</b><span>路径分解</span></div><i>→</i>
        <div><b>3</b><span>学习执行</span></div><i>→</i>
        <div><b>4</b><span>达成评估</span></div>
      </div>
    </section>

    <template v-else>
      <section class="welcome">
        <div>
          <h1>{{ greeting }}，{{ auth.user?.nickname }}！</h1>
          <p>已连续学习 <strong>{{ growth.streak }}</strong> 天，今天也向目标前进一步。</p>
        </div>
        <router-link v-if="growth.nextUnit" :to="`/micro/${growth.nextUnit.id}`" class="btn btn--primary">▶ 开始今日学习</router-link>
        <router-link v-else to="/path" class="btn btn--primary">查看已完成路径</router-link>
      </section>

      <section class="goal-card card">
        <div class="goal-card__main">
          <span class="tag">当前职业目标 · 进行中</span>
          <h2>{{ growth.goal.name }}</h2>
          <p>剩余 {{ daysLeft }} 天 · 每周 {{ growth.goal.weeklyHours }} 小时 · {{ growth.goal.baseLevel }}起步</p>
          <div class="goal-progress"><div :style="{ width: growth.achievement + '%' }"></div></div>
          <div class="goal-progress__labels">
            <strong>总达成度 {{ growth.achievement }}%</strong>
            <span>达标线 75%</span>
          </div>
          <div v-if="growth.nextMilestone" class="milestone">
            <span>🏁 下一里程碑</span>
            <strong>{{ growth.nextMilestone.name }}</strong>
            <small>第 {{ growth.nextMilestone.week }} 周 · {{ growth.nextMilestone.standard }}</small>
          </div>
        </div>
        <div class="goal-card__ring">
          <ProgressRing :percent="growth.achievement" :size="130" :stroke="10" color="#fff" />
          <router-link to="/career">调整目标 →</router-link>
        </div>
      </section>

      <div class="dashboard-grid">
        <section class="card panel competencies">
          <div class="panel-head"><h2>能力域进度</h2><router-link to="/path">完整路径 →</router-link></div>
          <div v-for="domain in growth.competencyProgress" :key="domain.id" class="competency">
            <div><strong>{{ domain.name }}</strong><span>{{ domain.progress }}%</span></div>
            <div class="bar"><div :style="{ width: domain.progress + '%', background: domain.color }"></div></div>
          </div>
        </section>

        <section class="card panel">
          <div class="panel-head"><h2>今日待办</h2><span>约 {{ growth.nextUnit ? growth.nextUnit.duration + 5 : 5 }} 分钟</span></div>
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
              <span class="task__check">{{ task.done ? '✓' : index + 1 }}</span>
              <span><strong>{{ task.title }}</strong><small>{{ task.meta }}</small></span>
              <i>{{ task.done ? '已完成' : '›' }}</i>
            </component>
          </div>
        </section>

        <section class="card panel ai-panel">
          <div class="panel-head"><h2>✨ AI 学习建议</h2><router-link to="/chat">继续追问</router-link></div>
          <p v-if="growth.achievement === 0">先完成第一个 5 分钟微单元。一次只掌握一个小知识点，比浏览大量内容更有效。</p>
          <p v-else>目前 <strong>{{ weakDomain?.name }}</strong> 达成度为 {{ weakDomain?.progress }}%，建议优先完成该能力域的微单元并通过快测。</p>
          <div class="suggestion-tags"><span>基于达成度</span><span>今日可完成</span><span>动态路径</span></div>
        </section>

        <section class="card panel incentive">
          <div class="panel-head"><h2>成长激励</h2><span>Lv{{ growth.level.number }}</span></div>
          <div class="points"><strong>{{ growth.points }}</strong><span>当前积分 · {{ growth.level.name }}</span></div>
          <div class="level-bar"><div :style="{ width: Math.min(100, growth.points / growth.level.next * 100) + '%' }"></div></div>
          <p>距下一等级还需 {{ growth.level.next - growth.points }} 积分</p>
          <div class="badges">
            <span v-for="badge in growth.badges" :key="badge">🏅 {{ badge }}</span>
            <span v-if="!growth.badges.length" class="muted">完成首个任务解锁勋章</span>
          </div>
        </section>
      </div>
    </template>
  </div>
</template>

<style scoped>
.dashboard { max-width: 1120px; }
.onboarding { padding: 42px; overflow: hidden; background: radial-gradient(circle at 90% 10%, #ddd6fe, transparent 35%), white; }
.eyebrow { color: var(--primary); font-size: 13px; font-weight: 800; letter-spacing: .05em; }
.onboarding h1 { max-width: 600px; margin: 10px 0; font-size: 30px; }
.onboarding p { max-width: 630px; color: var(--text-2); }
.loop { display: flex; align-items: center; gap: 12px; margin-top: 38px; }
.loop div { flex: 1; display: flex; align-items: center; gap: 8px; padding: 14px; border: 1px solid var(--border); border-radius: 10px; background: rgba(255,255,255,.8); }
.loop b { display: grid; place-items: center; width: 28px; height: 28px; border-radius: 8px; background: var(--primary-soft); color: var(--primary); }
.loop i { color: var(--text-3); }
.welcome { display: flex; justify-content: space-between; align-items: center; margin-bottom: 18px; }
.welcome h1 { margin: 0; font-size: 26px; }
.welcome p { margin: 2px 0 0; color: var(--text-2); }
.welcome p strong { color: var(--warning); }
.goal-card { display: flex; overflow: hidden; margin-bottom: 18px; color: white; border: none; background: linear-gradient(120deg, #312e81, #4f46e5 62%, #7c3aed); }
.goal-card__main { flex: 1; padding: 28px 32px; }
.goal-card .tag { background: rgba(255,255,255,.14); color: #e0e7ff; }
.goal-card h2 { margin: 9px 0 3px; font-size: 22px; }
.goal-card p { margin: 0 0 14px; color: #c7d2fe; font-size: 13px; }
.goal-progress { height: 8px; overflow: hidden; border-radius: 8px; background: rgba(255,255,255,.2); }
.goal-progress div { height: 100%; background: linear-gradient(90deg, #34d399, #a7f3d0); }
.goal-progress__labels { display: flex; justify-content: space-between; margin-top: 5px; font-size: 12px; color: #c7d2fe; }
.milestone { display: grid; grid-template-columns: auto 1fr; gap: 0 10px; margin-top: 18px; padding: 10px 13px; border-radius: 9px; background: rgba(255,255,255,.1); }
.milestone span { grid-row: span 2; align-self: center; color: #c7d2fe; font-size: 12px; }
.milestone small { color: #c7d2fe; }
.goal-card__ring { display: flex; width: 210px; flex-direction: column; align-items: center; justify-content: center; background: rgba(255,255,255,.06); }
.goal-card__ring :deep(.progress-ring__text) { fill: white; font-size: 17px; }
.goal-card__ring a { color: #e0e7ff; font-size: 12px; }
.dashboard-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 18px; }
.panel { padding: 22px; box-shadow: none; }
.panel-head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 16px; }
.panel-head h2 { margin: 0; font-size: 17px; }
.panel-head a, .panel-head span { font-size: 12px; }
.competency { margin: 14px 0; }
.competency > div:first-child { display: flex; justify-content: space-between; margin-bottom: 5px; font-size: 13px; }
.competency span { color: var(--text-2); }
.bar, .level-bar { height: 7px; overflow: hidden; border-radius: 8px; background: var(--border); }
.bar div, .level-bar div { height: 100%; border-radius: inherit; }
.task-list { display: flex; flex-direction: column; gap: 8px; }
.task { display: flex; align-items: center; gap: 11px; width: 100%; padding: 10px; border: none; border-radius: 9px; background: var(--surface-2); color: inherit; font: inherit; text-align: left; cursor: pointer; }
.task:hover { background: var(--primary-soft); }
.task.done { opacity: .65; }
.task__check { display: grid; place-items: center; width: 28px; height: 28px; flex: 0 0 28px; border: 1px solid var(--border); border-radius: 8px; background: white; color: var(--primary); font-weight: 700; }
.task > span:nth-child(2) { display: flex; flex: 1; flex-direction: column; }
.task strong { font-size: 13px; }
.task small { color: var(--text-3); }
.task i { color: var(--text-3); font-style: normal; font-size: 12px; }
.ai-panel { background: linear-gradient(140deg, #faf5ff, #fff); }
.ai-panel p { color: var(--text-2); }
.suggestion-tags { display: flex; gap: 6px; flex-wrap: wrap; }
.suggestion-tags span { padding: 3px 8px; border-radius: 99px; background: white; border: 1px solid #e9d5ff; color: #7e22ce; font-size: 11px; }
.points { display: flex; align-items: baseline; gap: 10px; }
.points strong { color: var(--warning); font-size: 30px; }
.points span, .incentive p { color: var(--text-3); font-size: 12px; }
.level-bar div { background: linear-gradient(90deg, var(--warning), #fbbf24); }
.badges { display: flex; gap: 7px; flex-wrap: wrap; }
.badges span { padding: 5px 9px; border-radius: 8px; background: #fffbeb; color: #92400e; font-size: 12px; }
.badges .muted { background: var(--surface-2); color: var(--text-3); }
@media (max-width: 760px) {
  .loop { flex-direction: column; align-items: stretch; }
  .loop i { display: none; }
  .welcome { align-items: flex-start; gap: 14px; }
  .goal-card { flex-direction: column; }
  .goal-card__ring { width: 100%; padding: 20px; }
  .dashboard-grid { grid-template-columns: 1fr; }
}
</style>
