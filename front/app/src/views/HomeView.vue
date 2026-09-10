<script setup>
import { computed, onMounted, ref } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useGrowthStore } from '../stores/growth'
import { useLearningStore } from '../stores/learning'
import { communityPosts } from '../data/community'
import { aiApi } from '../api'
import ProgressRing from '../components/ProgressRing.vue'
import Icon from '../components/Icon.vue'

const auth = useAuthStore()
const growth = useGrowthStore()
const learning = useLearningStore()
const aiSuggestions = ref([])

const greeting = computed(() => {
  const hour = new Date().getHours()
  return hour < 12 ? '早上好' : hour < 18 ? '下午好' : '晚上好'
})
const today = computed(() =>
  new Date().toLocaleDateString('zh-CN', { month: 'long', day: 'numeric', weekday: 'long' })
)
const daysLeft = computed(() => {
  if (!growth.goal) return 0
  return Math.max(0, Math.ceil((new Date(growth.goal.deadline) - new Date()) / 86400000))
})
const weakDomain = computed(() => [...growth.competencyProgress].sort((a, b) => a.progress - b.progress)[0])
const todayTasks = computed(() => {
  if (!growth.hasGoal) return []
  const tasks = []
  if (growth.dueReviewCount) {
    tasks.push({ title: `复习 ${growth.dueReviewCount} 个到期知识点`, meta: '遗忘曲线复习 · 约 5 分钟', to: '/review', icon: 'refresh' })
  }
  if (growth.nextUnit) {
    tasks.push({ title: growth.nextUnit.title, meta: `微单元 · ${growth.nextUnit.duration} 分钟`, to: `/micro/${growth.nextUnit.id}`, icon: 'zap' })
  }
  tasks.push({ title: '完成今日学习打卡', meta: '打卡获得 10 成长积分', action: 'checkin', done: growth.checkedInToday, icon: 'calendar' })
  return tasks.slice(0, 3)
})
const pendingCount = computed(() => todayTasks.value.filter((item) => !item.done).length)

const loopSteps = [
  { title: '目标确认', desc: 'AI 访谈与职业匹配', icon: 'target' },
  { title: '路径分解', desc: '能力域 · 知识点 · 里程碑', icon: 'layers' },
  { title: '学习执行', desc: '5 分钟微单元与快测', icon: 'zap' },
  { title: '达成评估', desc: '雷达图 · 趋势 · AI 复盘', icon: 'chart' }
]

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
    <!-- 未设定目标：引导页 -->
    <section v-if="!growth.hasGoal" class="onboarding card card--ink">
      <div class="onboarding__copy">
        <span class="onboarding__eyebrow"><Icon name="sparkles" :size="14" /> 目标驱动 · AI 辅助 · 数据评估</span>
        <h1>{{ greeting }}，{{ auth.user?.nickname || '学习者' }}<br />先确定你想成为的人</h1>
        <p>完成 3 步职业规划，系统会根据你的基础和时间生成能力图谱、学习任务与里程碑，把每一次学习变成可验证的成长。</p>
        <div class="onboarding__actions">
          <router-link to="/career" class="btn btn--light btn--lg">
            <Icon name="compass" :size="18" /> 开始 AI 职业规划
          </router-link>
          <router-link to="/courses" class="btn btn--glass btn--lg">先逛逛课程</router-link>
        </div>
      </div>
      <ol class="loop">
        <li v-for="(item, index) in loopSteps" :key="item.title">
          <span class="loop__index">0{{ index + 1 }}</span>
          <span class="loop__icon"><Icon :name="item.icon" :size="18" /></span>
          <strong>{{ item.title }}</strong>
          <small>{{ item.desc }}</small>
        </li>
      </ol>
    </section>

    <template v-else>
      <section class="welcome">
        <div>
          <span class="welcome__date">{{ today }}</span>
          <h1>{{ greeting }}，{{ auth.user?.nickname }}</h1>
          <p>
            已连续学习 <strong class="num">{{ growth.streak }}</strong> 天 ·
            今日还有 <strong class="num">{{ pendingCount }}</strong> 项任务待完成
          </p>
        </div>
        <router-link :to="growth.nextUnit ? `/micro/${growth.nextUnit.id}` : '/path'" class="btn btn--primary btn--lg">
          <Icon name="play" :size="16" /> {{ growth.nextUnit ? '开始今日学习' : '查看完成路径' }}
        </router-link>
      </section>

      <section class="goal-card card card--gradient">
        <div class="goal-card__content">
          <span class="tag tag--glass"><Icon name="flag" :size="12" /> 当前职业目标 · 进行中</span>
          <h2>{{ growth.goal.name }}</h2>
          <p>剩余 {{ daysLeft }} 天 · 每周 {{ growth.goal.weeklyHours }} 小时 · {{ growth.goal.baseLevel }}起步</p>
          <div class="progress progress--glass"><i :style="{ width: growth.achievement + '%' }"></i></div>
          <div class="goal-labels">
            <strong>总达成度 {{ growth.achievement }}%</strong>
            <span>达标线 75%</span>
          </div>
          <div v-if="growth.nextMilestone" class="milestone">
            <span class="milestone__icon"><Icon name="flag" :size="16" /></span>
            <div>
              <small>下一里程碑 · 第 {{ growth.nextMilestone.week }} 周</small>
              <strong>{{ growth.nextMilestone.name }}</strong>
              <em>{{ growth.nextMilestone.standard }}</em>
            </div>
          </div>
        </div>
        <div class="goal-card__score">
          <ProgressRing :percent="growth.achievement" :size="148" :stroke="11" color="#fff" track="rgba(255,255,255,.18)" />
          <router-link to="/stats">查看达成报告 <Icon name="arrowRight" :size="13" /></router-link>
        </div>
      </section>

      <section class="quick-stats stagger">
        <router-link to="/path" class="stat card card--hover">
          <span class="stat__icon"><Icon name="zap" :size="20" /></span>
          <div class="stat__body"><strong class="stat__value">{{ growth.completedUnitCount }}</strong><small class="stat__label">已完成微单元</small></div>
        </router-link>
        <router-link to="/review" class="stat card card--hover">
          <span class="stat__icon stat__icon--warning"><Icon name="refresh" :size="20" /></span>
          <div class="stat__body"><strong class="stat__value">{{ growth.dueReviewCount }}</strong><small class="stat__label">今日待复习</small></div>
        </router-link>
        <router-link to="/incentives" class="stat card card--hover">
          <span class="stat__icon stat__icon--success"><Icon name="medal" :size="20" /></span>
          <div class="stat__body"><strong class="stat__value">{{ growth.points }}</strong><small class="stat__label">成长积分 · Lv{{ growth.level.number }}</small></div>
        </router-link>
        <router-link to="/courses" class="stat card card--hover">
          <span class="stat__icon stat__icon--info"><Icon name="book" :size="20" /></span>
          <div class="stat__body"><strong class="stat__value">{{ learning.completedCount }}</strong><small class="stat__label">已完成课程章节</small></div>
        </router-link>
      </section>

      <div class="dashboard-grid">
        <section class="card panel">
          <div class="panel-head">
            <div><h2>能力域进度</h2><p>按知识点掌握度加权聚合</p></div>
            <router-link to="/path">完整路径 →</router-link>
          </div>
          <div v-for="domain in growth.competencyProgress" :key="domain.id" class="competency">
            <div class="competency__head">
              <span class="competency__dot" :style="{ background: domain.color }"></span>
              <strong>{{ domain.name }}</strong>
              <span class="num">{{ domain.progress }}%</span>
            </div>
            <div class="progress progress--thin"><i :style="{ width: domain.progress + '%', background: domain.color }"></i></div>
          </div>
        </section>

        <section class="card panel">
          <div class="panel-head">
            <div><h2>今日待办</h2><p>根据路径、薄弱点和复习日程生成</p></div>
            <span class="tag tag--neutral">{{ todayTasks.length }} 项</span>
          </div>
          <div class="task-list">
            <component
              :is="task.to ? 'router-link' : 'button'"
              v-for="task in todayTasks"
              :key="task.title"
              :to="task.to"
              class="task"
              :class="{ 'task--done': task.done }"
              @click="task.action === 'checkin' && growth.checkIn()"
            >
              <span class="task__icon"><Icon :name="task.done ? 'check' : task.icon" :size="17" /></span>
              <span class="task__body"><strong>{{ task.title }}</strong><small>{{ task.meta }}</small></span>
              <span class="task__end">{{ task.done ? '已完成' : '' }}<Icon v-if="!task.done" name="arrowRight" :size="15" /></span>
            </component>
          </div>
        </section>

        <section class="card panel ai-panel">
          <div class="panel-head">
            <div><h2><span class="tag tag--ai">AI</span> 学习建议</h2><p>基于当前达成数据动态生成</p></div>
            <router-link to="/chat">问 AI →</router-link>
          </div>
          <ul v-if="aiSuggestions.length" class="ai-suggest-list">
            <li v-for="item in aiSuggestions" :key="item">{{ item }}</li>
          </ul>
          <template v-else>
            <p v-if="growth.achievement === 0" class="ai-text">从第一个 5 分钟微单元开始。一次掌握一个小知识点，通过快测后再进入下一步。</p>
            <p v-else class="ai-text">当前最需要关注的是 <strong>{{ weakDomain?.name }}</strong>（{{ weakDomain?.progress }}%）。优先完成该能力域任务，并按间隔计划复习。</p>
          </template>
          <div class="ai-chips"><span>基于达成度</span><span>关键路径优先</span><span>每日可完成</span></div>
        </section>

        <section class="card panel growth-panel">
          <div class="panel-head">
            <div><h2>成长激励</h2><p>{{ growth.level.name }}</p></div>
            <router-link to="/incentives">勋章墙 →</router-link>
          </div>
          <div class="points">
            <strong class="num">{{ growth.points }}</strong>
            <span>累计积分</span>
            <b class="num">Lv{{ growth.level.number }}</b>
          </div>
          <div class="progress progress--thin level-track"><i :style="{ width: Math.min(100, growth.points / growth.level.next * 100) + '%' }"></i></div>
          <div class="badges">
            <span v-for="badge in growth.badges.slice(0, 3)" :key="badge"><Icon name="medal" :size="12" /> {{ badge }}</span>
            <span v-if="!growth.badges.length" class="badges__muted">完成首个任务解锁勋章</span>
          </div>
        </section>
      </div>
    </template>

    <section class="community card panel">
      <div class="panel-head">
        <div><h2>同路人正在分享</h2><p>学习不孤单，看看社区里的实践与复盘</p></div>
        <span class="tag tag--neutral">学习动态</span>
      </div>
      <div class="post-grid">
        <article v-for="post in communityPosts.slice(0, 3)" :key="post.id">
          <div class="post-author">
            <span :style="{ background: post.avatarColor }">{{ post.avatar }}</span>
            <b>{{ post.author }}</b>
            <small>{{ post.time }}</small>
          </div>
          <h3>{{ post.title }}</h3>
          <p>{{ post.excerpt }}</p>
          <div class="post-foot">
            <span v-for="tag in post.tags" :key="tag">#{{ tag }}</span>
            <small><Icon name="heart" :size="12" /> {{ post.likes }} · <Icon name="message" :size="12" /> {{ post.comments }}</small>
          </div>
        </article>
      </div>
    </section>
  </div>
</template>

<style scoped>
.dashboard { max-width: 1160px; }

/* 引导 */
.onboarding {
  padding: 56px;
  border-radius: var(--radius-xl);
}
.onboarding__copy { max-width: 640px; }
.onboarding__eyebrow {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 12px;
  border-radius: 999px;
  border: 1px solid rgba(255, 255, 255, 0.16);
  background: rgba(255, 255, 255, 0.06);
  color: #c7bfff;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.08em;
}
.onboarding h1 { margin: 18px 0 14px; font-size: 42px; line-height: 1.2; }
.onboarding p { margin: 0; max-width: 560px; color: rgba(201, 203, 230, 0.85); font-size: 16px; }
.onboarding__actions { display: flex; flex-wrap: wrap; gap: 12px; margin-top: 28px; }
.loop {
  list-style: none;
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
  margin: 48px 0 0;
  padding: 0;
}
.loop li {
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 18px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(6px);
}
.loop__index { position: absolute; top: 14px; right: 16px; font-family: 'Sora', sans-serif; font-size: 12px; color: rgba(201, 203, 230, 0.4); }
.loop__icon { display: grid; place-items: center; width: 36px; height: 36px; margin-bottom: 8px; border-radius: 11px; background: rgba(107, 92, 255, 0.3); color: #c7bfff; }
.loop strong { font-size: 14.5px; }
.loop small { color: rgba(201, 203, 230, 0.65); font-size: 12px; }

/* 欢迎 */
.welcome { display: flex; align-items: flex-end; justify-content: space-between; gap: 20px; margin-bottom: 22px; }
.welcome__date { color: var(--text-3); font-size: 13px; font-weight: 500; }
.welcome h1 { margin: 4px 0 6px; font-size: 30px; }
.welcome p { margin: 0; color: var(--text-2); font-size: 14.5px; }
.welcome p strong { color: var(--primary-strong); font-size: 16px; }

/* 目标卡 */
.goal-card { display: flex; margin-bottom: 18px; border-radius: var(--radius-lg); }
.goal-card__content { flex: 1; padding: 30px 34px; }
.goal-card h2 { margin: 14px 0 4px; font-size: 24px; }
.goal-card p { margin: 0 0 18px; color: rgba(255, 255, 255, 0.75); font-size: 13.5px; }
.goal-labels { display: flex; justify-content: space-between; margin-top: 8px; color: rgba(255, 255, 255, 0.8); font-size: 12.5px; }
.goal-labels strong { color: #fff; }
.milestone {
  display: flex;
  gap: 14px;
  align-items: center;
  margin-top: 20px;
  padding: 14px 16px;
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.12);
  border: 1px solid rgba(255, 255, 255, 0.14);
}
.milestone__icon { display: grid; place-items: center; width: 38px; height: 38px; flex: 0 0 38px; border-radius: 12px; background: rgba(255, 255, 255, 0.16); }
.milestone div { display: flex; flex-direction: column; line-height: 1.4; }
.milestone small { color: rgba(255, 255, 255, 0.7); font-size: 11.5px; }
.milestone strong { font-size: 15px; }
.milestone em { font-style: normal; color: rgba(255, 255, 255, 0.75); font-size: 12.5px; }
.goal-card__score {
  display: flex;
  width: 230px;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 14px;
  background: rgba(255, 255, 255, 0.08);
  border-left: 1px solid rgba(255, 255, 255, 0.12);
}
.goal-card__score :deep(.progress-ring__text) { fill: white; font-size: 24px; }
.goal-card__score a { display: inline-flex; align-items: center; gap: 4px; color: #fff; font-size: 12.5px; font-weight: 600; opacity: 0.9; }

/* 快速统计 */
.quick-stats { display: grid; grid-template-columns: repeat(4, 1fr); gap: 14px; margin-bottom: 18px; }

/* 面板网格 */
.dashboard-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 18px; }
.competency { margin: 14px 0; }
.competency__head { display: flex; align-items: center; gap: 8px; margin-bottom: 7px; font-size: 13.5px; }
.competency__dot { width: 8px; height: 8px; border-radius: 50%; }
.competency__head strong { flex: 1; font-weight: 600; }
.competency__head span:last-child { color: var(--text-2); font-weight: 600; }

.task-list { display: flex; flex-direction: column; gap: 8px; }
.task {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 11px 12px;
  border: 1px solid transparent;
  border-radius: 12px;
  background: var(--surface-2);
  color: inherit;
  font: inherit;
  text-align: left;
  cursor: pointer;
  transition: all var(--t-fast);
}
.task:hover { background: var(--primary-soft); border-color: var(--primary-soft-2); }
.task--done { opacity: 0.6; }
.task__icon { display: grid; place-items: center; width: 34px; height: 34px; flex: 0 0 34px; border-radius: 10px; background: var(--surface); color: var(--primary); box-shadow: var(--shadow-xs); }
.task--done .task__icon { background: var(--success); color: #fff; }
.task__body { display: flex; flex: 1; min-width: 0; flex-direction: column; line-height: 1.4; }
.task__body strong { font-size: 13.5px; }
.task__body small { color: var(--text-3); font-size: 12px; }
.task__end { display: inline-flex; align-items: center; color: var(--text-3); font-size: 12px; }

.ai-panel { background: linear-gradient(150deg, #f6f3ff 0%, #fff 60%); }
.ai-panel .panel-head h2 { display: flex; align-items: center; gap: 8px; }
.ai-text { margin: 0 0 14px; color: var(--text-2); font-size: 14px; }
.ai-suggest-list { margin: 0 0 14px; padding-left: 18px; color: var(--text-2); font-size: 13.5px; line-height: 1.75; }
.ai-suggest-list li::marker { color: var(--primary); }
.ai-chips { display: flex; gap: 6px; flex-wrap: wrap; }
.ai-chips span { padding: 3px 10px; border: 1px solid var(--primary-soft-2); border-radius: 99px; background: var(--surface); color: var(--primary-strong); font-size: 11.5px; font-weight: 600; }

.points { display: flex; align-items: baseline; gap: 8px; margin-bottom: 10px; }
.points strong { color: var(--warning-strong); font-size: 32px; font-weight: 700; }
.points span { color: var(--text-3); font-size: 12px; }
.points b { margin-left: auto; padding: 3px 10px; border-radius: 99px; background: var(--primary-soft); color: var(--primary-strong); font-size: 12.5px; }
.level-track i { background: linear-gradient(90deg, #f59e0b, #fcd34d); }
.badges { display: flex; gap: 6px; margin-top: 16px; flex-wrap: wrap; }
.badges span { display: inline-flex; align-items: center; gap: 4px; padding: 5px 10px; border-radius: 9px; background: var(--warning-soft); color: var(--warning-strong); font-size: 12px; font-weight: 600; }
.badges__muted { color: var(--text-3) !important; background: var(--surface-2) !important; font-weight: 500 !important; }

/* 社区 */
.community { margin-top: 18px; }
.post-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 14px; }
.post-grid article {
  display: flex;
  flex-direction: column;
  padding: 18px;
  border: 1px solid var(--border);
  border-radius: 14px;
  background: var(--surface-2);
  transition: all var(--t-fast);
}
.post-grid article:hover { background: var(--surface); border-color: var(--border-strong); box-shadow: var(--shadow-sm); }
.post-author { display: flex; align-items: center; gap: 8px; }
.post-author > span { display: grid; place-items: center; width: 28px; height: 28px; border-radius: 9px; color: white; font-size: 11px; font-weight: 700; }
.post-author b { font-size: 12.5px; }
.post-author small { margin-left: auto; color: var(--text-3); font-size: 11px; }
.post-grid h3 { margin: 12px 0 6px; font-size: 14.5px; }
.post-grid p { display: -webkit-box; flex: 1; margin: 0 0 12px; overflow: hidden; color: var(--text-2); font-size: 12.5px; line-height: 1.6; -webkit-box-orient: vertical; -webkit-line-clamp: 2; }
.post-foot { display: flex; gap: 6px; align-items: center; color: var(--primary); font-size: 11.5px; font-weight: 600; }
.post-foot small { display: inline-flex; align-items: center; gap: 3px; margin-left: auto; color: var(--text-3); font-weight: 500; }

@media (max-width: 960px) {
  .onboarding { padding: 36px 28px; }
  .onboarding h1 { font-size: 30px; }
  .loop, .quick-stats { grid-template-columns: 1fr 1fr; }
  .goal-card { flex-direction: column; }
  .goal-card__score { width: 100%; padding: 22px; border-left: none; border-top: 1px solid rgba(255, 255, 255, 0.12); }
  .dashboard-grid { grid-template-columns: 1fr; }
  .post-grid { grid-template-columns: 1fr; }
}
@media (max-width: 520px) {
  .welcome { flex-direction: column; align-items: flex-start; gap: 12px; }
  .loop, .quick-stats { grid-template-columns: 1fr; }
}
</style>
