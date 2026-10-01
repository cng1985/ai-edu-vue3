<script setup>
import { computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useGrowthStore } from '../stores/growth'
import { FLYWHEEL, money, pct } from '../utils/eco'
import ProgressRing from '../components/ProgressRing.vue'
import SkillLevel from '../components/SkillLevel.vue'
import MasteryBar from '../components/MasteryBar.vue'
import Icon from '../components/Icon.vue'

const auth = useAuthStore()
const growth = useGrowthStore()
const ov = computed(() => growth.overview)

const greeting = computed(() => {
  const h = new Date().getHours()
  return h < 12 ? '早上好' : h < 18 ? '下午好' : '晚上好'
})

const chain = [
  { title: '职业目标', icon: 'flag' },
  { title: '能力分析', icon: 'target' },
  { title: 'AI 学习规划', icon: 'sparkles' },
  { title: '知识学习', icon: 'book' },
  { title: '社区交流', icon: 'users' },
  { title: '项目实践', icon: 'code' },
  { title: '真实任务', icon: 'bolt' },
  { title: '获得收益', icon: 'gift' },
  { title: '能力提升', icon: 'trend' }
]

const flywheel = computed(() => {
  const f = ov.value?.flywheel || {}
  const values = {
    learn: { value: f.learningEvents || 0, hint: `${Math.round((f.learningMinutes || 0) / 60)} 小时学习` },
    practice: { value: f.projectsPracticed || 0, hint: `${f.submissions || 0} 次任务提交` },
    output: { value: f.evidenceCount || 0, hint: '项能力证据' },
    income: { value: money(f.income), hint: `${f.tasksCompleted || 0} 个真实任务` },
    invest: { value: f.masteredCount || 0, hint: '个知识点已掌握' },
    stronger: { value: 'L' + (f.avgSkillLevel || 0), hint: '平均技能等级' }
  }
  return FLYWHEEL.map((s) => ({ ...s, ...values[s.key] }))
})

const heat = computed(() => {
  const days = ov.value?.activity || []
  const max = Math.max(1, ...days.map((d) => d.events))
  return days.map((d) => ({ ...d, level: d.events ? Math.ceil((d.events / max) * 4) : 0 }))
})

const planWeek = computed(() => {
  const out = ov.value?.latestPlan?.output
  return out?.plan?.weeks?.[0] || null
})

onMounted(() => growth.refresh())
</script>

<template>
  <div class="page growth">
    <div v-if="!ov && growth.loading" class="card empty-state"><p>正在加载你的成长数据…</p></div>
    <div v-else-if="growth.error" class="card empty-state">
      <h3>加载失败</h3>
      <p>{{ growth.error }}</p>
      <button class="btn btn--primary" @click="growth.refresh()">重试</button>
    </div>

    <template v-else-if="ov">
      <section v-if="!ov.goal" class="onboard card card--ink">
        <span class="onboard__eyebrow"><Icon name="sparkles" :size="14" /> AI 时代个人成长操作系统</span>
        <h1>{{ greeting }}，{{ auth.user?.nickname }}<br />先确定你想成为的人</h1>
        <p>选择职业目标后，平台会基于「职业 → 岗位 → 能力 → 技能 → 知识」模型分析你的差距，由 AI 规划学习，再通过项目实践与真实任务把知识变成可验证的能力与收入。</p>
        <div class="row row--wrap">
          <router-link v-if="auth.hasPermission('growth:write')" to="/career" class="btn btn--light btn--lg"><Icon name="compass" :size="18" /> 选择职业目标</router-link>
          <router-link v-else to="/register" class="btn btn--light btn--lg"><Icon name="user" :size="18" /> 注册后开始成长</router-link>
          <router-link to="/knowledge" class="btn btn--glass btn--lg">浏览知识图谱</router-link>
        </div>
        <ol class="chain">
          <li v-for="(c, i) in chain" :key="c.title">
            <span><Icon :name="c.icon" :size="15" /></span>
            <b>{{ c.title }}</b>
            <i v-if="i < chain.length - 1"><Icon name="arrowRight" :size="12" /></i>
          </li>
        </ol>
      </section>

      <template v-else>
        <section class="hero card card--gradient">
          <div class="hero__main">
            <span class="tag tag--glass"><Icon name="flag" :size="12" /> {{ ov.gap?.career?.name }} · 职业目标</span>
            <h1>{{ ov.gap?.role?.name }}</h1>
            <p>{{ greeting }}，{{ auth.user?.nickname }}。已达标 {{ ov.gap?.metCount }}/{{ ov.gap?.totalCount }} 项技能要求，每周投入 {{ ov.goal.weeklyHours }} 小时，目标 {{ ov.goal.targetWeeks }} 周。</p>
            <div class="hero__caps">
              <div v-for="cap in ov.gap?.capabilities" :key="cap.id" class="hero__cap">
                <div class="row row--between"><span>{{ cap.name }}</span><b class="num">{{ cap.readiness }}%</b></div>
                <div class="progress progress--thin progress--glass"><i :style="{ width: cap.readiness + '%' }"></i></div>
              </div>
            </div>
            <div class="row row--wrap">
              <router-link to="/kernel" class="btn btn--light"><Icon name="sparkles" :size="16" /> 运行 AI 学习内核</router-link>
              <router-link :to="`/career?role=${ov.goal.roleId}`" class="btn btn--glass">查看能力模型</router-link>
            </div>
          </div>
          <div class="hero__ring">
            <ProgressRing :percent="ov.gap?.readiness || 0" :size="150" :stroke="11" color="#fff" track="rgba(255,255,255,.2)" />
            <span>岗位达成度</span>
          </div>
        </section>

        <section class="card panel flywheel">
          <div class="panel-head">
            <div><h2>个人成长飞轮</h2><p>学习 → 实践 → 产出 → 收益 → 投资学习 → 更强能力</p></div>
            <router-link to="/profile">人才画像 →</router-link>
          </div>
          <div class="flywheel__ring">
            <div v-for="(s, i) in flywheel" :key="s.key" class="flywheel__node">
              <span class="flywheel__icon"><Icon :name="s.icon" :size="18" /></span>
              <small>{{ s.label }}</small>
              <strong class="num">{{ s.value }}</strong>
              <em>{{ s.hint }}</em>
              <i v-if="i < flywheel.length - 1" class="flywheel__arrow"><Icon name="arrowRight" :size="14" /></i>
            </div>
          </div>
        </section>

        <div class="grid-2">
          <section class="card panel">
            <div class="panel-head">
              <div><h2>关键能力差距</h2><p>按差距 × 能力权重排序</p></div>
              <router-link :to="`/career?role=${ov.goal.roleId}`">全部 →</router-link>
            </div>
            <div v-if="!ov.gap?.gaps?.length" class="muted">全部技能已达标，去任务市场接真实任务吧！</div>
            <div v-for="g in ov.gap?.gaps?.slice(0, 5)" :key="g.skillId" class="gap-row">
              <div class="list-row__body">
                <strong>{{ g.skillName }}</strong>
                <small>{{ g.capabilityName }}</small>
              </div>
              <SkillLevel :level="g.current" :required="g.required" />
            </div>
          </section>

          <section class="card panel">
            <div class="panel-head">
              <div><h2><span class="tag tag--ai">AI</span> 下一步学什么</h2><p>基于差距与知识图谱前置关系推荐</p></div>
              <router-link to="/knowledge">知识图谱 →</router-link>
            </div>
            <div class="stack stack--tight">
              <router-link v-for="r in ov.recommendations.slice(0, 5)" :key="r.knowledge.id" :to="`/knowledge/${r.knowledge.id}`" class="list-row">
                <span class="rec-icon" :class="{ 'rec-icon--blocked': !r.ready }"><Icon :name="r.ready ? 'play' : 'lock'" :size="14" /></span>
                <div class="list-row__body">
                  <strong>{{ r.knowledge.name }}</strong>
                  <small>{{ r.ready ? r.reason : '需先学习：' + r.missingPrereqs.map((p) => p.name).join('、') }}</small>
                </div>
                <MasteryBar :value="r.mastery" compact />
              </router-link>
            </div>
          </section>
        </div>

        <div class="grid-3">
          <section class="card panel">
            <div class="panel-head"><div><h3>今日复习</h3><p>按遗忘曲线到期</p></div><span class="tag tag--warning">{{ ov.dueReviews.length }}</span></div>
            <div v-if="!ov.dueReviews.length" class="muted small">暂无到期复习，保持节奏！</div>
            <router-link v-for="d in ov.dueReviews" :key="d.knowledge.id" :to="`/knowledge/${d.knowledge.id}`" class="mini-row">
              <span>{{ d.knowledge.name }}</span><b class="num">{{ pct(d.effectiveMastery) }}%</b>
            </router-link>
          </section>

          <section class="card panel">
            <div class="panel-head"><div><h3>推荐实践</h3><p>跳一跳够得着的项目任务</p></div><router-link to="/projects">项目 →</router-link></div>
            <div v-if="!ov.nextProjects.length" class="muted small">暂无推荐任务</div>
            <router-link v-for="t in ov.nextProjects" :key="t.taskId" :to="`/projects/${t.projectId}?task=${t.taskId}`" class="mini-row">
              <span>{{ t.taskTitle }}<small>{{ t.projectTitle }}</small></span><b class="num">{{ t.match.score }}%</b>
            </router-link>
          </section>

          <section class="card panel">
            <div class="panel-head"><div><h3>匹配的真实任务</h3><p>完成任务获得收益</p></div><router-link to="/market">市场 →</router-link></div>
            <div v-if="!ov.opportunities.length" class="muted small">暂无匹配任务</div>
            <router-link v-for="o in ov.opportunities" :key="o.id" :to="`/market/${o.id}`" class="mini-row">
              <span>{{ o.title }}<small>{{ o.company }} · {{ money(o.budget) }}</small></span><b class="num">{{ o.match?.score }}%</b>
            </router-link>
          </section>
        </div>

        <div class="grid-2">
          <section class="card panel">
            <div class="panel-head">
              <div><h2>本周学习计划</h2><p>{{ ov.latestPlan ? ov.latestPlan.output.plan?.summary : '由 AI 学习内核生成' }}</p></div>
              <router-link to="/kernel">{{ ov.latestPlan ? '重新规划' : '生成计划' }} →</router-link>
            </div>
            <div v-if="!planWeek" class="muted small">还没有学习计划，运行一次 AI 学习内核即可生成。</div>
            <template v-else>
              <router-link
                v-for="item in planWeek.items"
                :key="item.refId + (item.subId || '')"
                :to="item.type === 'project' ? `/projects/${item.refId}?task=${item.subId}` : `/knowledge/${item.refId}`"
                class="mini-row"
              >
                <span><Icon :name="item.type === 'project' ? 'code' : 'book'" :size="13" /> {{ item.title }}<small>{{ item.reason }}</small></span>
                <b class="num">{{ item.minutes }}′</b>
              </router-link>
            </template>
          </section>

          <section class="card panel">
            <div class="panel-head"><div><h2>近 4 周学习活跃度</h2><p>每一次学习都在积累个人成长数据</p></div></div>
            <div class="heat">
              <span v-for="d in heat" :key="d.date" :class="`heat--${d.level}`" :title="`${d.date}：${d.events} 次学习，${d.minutes} 分钟`"></span>
            </div>
            <div class="row small muted heat__legend">少 <span class="heat--0"></span><span class="heat--1"></span><span class="heat--2"></span><span class="heat--3"></span><span class="heat--4"></span> 多</div>
          </section>
        </div>
      </template>
    </template>
  </div>
</template>

<style scoped>
.growth { display: flex; flex-direction: column; gap: 18px; }
.onboard { padding: 48px; border-radius: var(--radius-xl); }
.onboard__eyebrow { display: inline-flex; align-items: center; gap: 8px; padding: 6px 12px; border: 1px solid rgba(255,255,255,.16); border-radius: 999px; background: rgba(255,255,255,.06); color: #c7bfff; font-size: 12px; font-weight: 700; letter-spacing: .08em; }
.onboard h1 { margin: 18px 0 12px; font-size: 38px; line-height: 1.25; }
.onboard p { max-width: 640px; margin: 0 0 26px; color: rgba(201,203,230,.88); font-size: 15.5px; }
.chain { display: flex; flex-wrap: wrap; gap: 8px; margin: 36px 0 0; padding: 0; list-style: none; }
.chain li { display: flex; align-items: center; gap: 8px; }
.chain span { display: grid; place-items: center; width: 30px; height: 30px; border-radius: 10px; background: rgba(107,92,255,.32); color: #d6d0ff; }
.chain b { font-size: 13px; font-weight: 600; }
.chain i { color: rgba(201,203,230,.4); }

.hero { display: flex; border-radius: var(--radius-lg); }
.hero__main { flex: 1; padding: 30px 34px; }
.hero h1 { margin: 12px 0 6px; font-size: 30px; }
.hero p { margin: 0 0 18px; color: rgba(255,255,255,.8); font-size: 13.5px; }
.hero__caps { display: grid; grid-template-columns: repeat(auto-fill, minmax(190px, 1fr)); gap: 10px 18px; margin-bottom: 20px; }
.hero__cap { font-size: 12.5px; color: rgba(255,255,255,.85); }
.hero__cap .progress { margin-top: 5px; }
.hero__ring { display: flex; width: 220px; flex-direction: column; align-items: center; justify-content: center; gap: 10px; border-left: 1px solid rgba(255,255,255,.14); background: rgba(255,255,255,.07); font-size: 13px; font-weight: 600; }
.hero__ring :deep(.progress-ring__text) { fill: #fff; font-size: 24px; }

.flywheel__ring { display: grid; grid-template-columns: repeat(6, 1fr); gap: 10px; }
.flywheel__node { position: relative; display: flex; flex-direction: column; align-items: flex-start; gap: 2px; padding: 14px; border: 1px solid var(--border); border-radius: 14px; background: var(--surface-2); }
.flywheel__icon { display: grid; place-items: center; width: 34px; height: 34px; margin-bottom: 6px; border-radius: 11px; background: var(--primary-soft); color: var(--primary); }
.flywheel__node small { color: var(--text-3); font-size: 12px; font-weight: 600; }
.flywheel__node strong { font-size: 22px; line-height: 1.2; }
.flywheel__node em { color: var(--text-3); font-size: 11.5px; font-style: normal; }
.flywheel__arrow { position: absolute; right: -11px; top: 50%; z-index: 1; display: grid; place-items: center; width: 18px; height: 18px; margin-top: -9px; border-radius: 50%; background: var(--surface); color: var(--primary); box-shadow: var(--shadow-xs); }

.gap-row { display: flex; align-items: center; gap: 12px; padding: 10px 0; border-bottom: 1px dashed var(--border); }
.gap-row:last-child { border-bottom: none; }
.stack--tight { gap: 8px; }
.rec-icon { display: grid; place-items: center; width: 30px; height: 30px; flex: 0 0 30px; border-radius: 10px; background: var(--primary-soft); color: var(--primary); }
.rec-icon--blocked { background: var(--surface-3); color: var(--text-3); }

.mini-row { display: flex; align-items: center; justify-content: space-between; gap: 10px; padding: 9px 0; border-bottom: 1px dashed var(--border); color: inherit; font-size: 13.5px; }
.mini-row:last-child { border-bottom: none; }
.mini-row span { display: flex; flex-direction: column; min-width: 0; }
.mini-row span small { color: var(--text-3); font-size: 11.5px; }
.mini-row b { color: var(--primary-strong); font-size: 13px; }
.mini-row:hover span { color: var(--primary); }

.heat { display: grid; grid-template-columns: repeat(14, 1fr); gap: 5px; }
.heat span, .heat__legend span { aspect-ratio: 1; border-radius: 5px; background: var(--surface-3); }
.heat__legend { gap: 4px; margin-top: 12px; }
.heat__legend span { width: 12px; }
.heat--1 { background: #dcd7ff !important; }
.heat--2 { background: #b3a8ff !important; }
.heat--3 { background: #8676ff !important; }
.heat--4 { background: var(--primary-strong) !important; }

@media (max-width: 960px) {
  .hero { flex-direction: column; }
  .hero__ring { width: 100%; padding: 20px; border-left: none; border-top: 1px solid rgba(255,255,255,.14); }
  .flywheel__ring { grid-template-columns: repeat(3, 1fr); }
  .flywheel__arrow { display: none; }
  .onboard { padding: 32px 24px; }
  .onboard h1 { font-size: 28px; }
}
</style>
