<script setup>
import { computed } from 'vue'
import { courses, totalChapterCount } from '../data/courses'
import { quizzes } from '../data/quizzes'
import { frontendPath, microUnits } from '../data/careerPath'
import { useLearningStore } from '../stores/learning'
import { useGrowthStore } from '../stores/growth'
import ProgressRing from '../components/ProgressRing.vue'
import AbilityRadar from '../components/AbilityRadar.vue'
import Icon from '../components/Icon.vue'

const learning = useLearningStore()
const growth = useGrowthStore()

const weekly = computed(() => learning.weeklyActivity.map((day) => {
  const microCount = Object.values(growth.completedUnits).filter(
    (timestamp) => new Date(timestamp).toISOString().slice(0, 10) === day.date
  ).length
  return { ...day, count: day.count + microCount }
}))
const weeklyMax = computed(() =>
  Math.max(1, ...weekly.value.map((d) => d.count))
)
const weeklyTotal = computed(() => weekly.value.reduce((sum, d) => sum + d.count, 0))

const quizRows = computed(() =>
  quizzes.map((q) => {
    const r = learning.quizResults[q.id]
    return {
      id: q.id,
      title: q.title,
      total: q.questions.length,
      result: r || null,
      percent: r ? Math.round((r.score / r.total) * 100) : null
    }
  })
)

const noteEntries = computed(() =>
  Object.entries(learning.notes)
    .filter(([, text]) => text && text.trim())
    .map(([key, text]) => {
      const [courseId, chapterId] = key.split('/')
      const course = courses.find((c) => c.id === courseId)
      const chapter = course?.chapters.find((ch) => ch.id === chapterId)
      return { key, courseId, chapterId, course, chapter, text }
    })
    .filter((e) => e.course && e.chapter)
)

const weakDomain = computed(() =>
  [...growth.competencyProgress].sort((a, b) => a.progress - b.progress)[0] || null
)

const reinforcementUnit = computed(() => {
  if (!weakDomain.value) return null
  const domain = frontendPath.competencies.find((item) => item.id === weakDomain.value.id)
  return domain?.points
    .map((point) => microUnits.find((unit) => unit.id === point.unitId))
    .find((unit) => unit && !growth.isUnitCompleted(unit.id)) || null
})

const milestoneRows = computed(() => frontendPath.milestones.map((milestone, index) => {
  const requiredUnits = (index + 1) * 2
  const done = growth.completedUnitCount >= requiredUnits
  const current = !done && growth.completedUnitCount >= index * 2
  return { ...milestone, done, current }
}))

const trendValues = computed(() => {
  const values = growth.snapshots.map((item) => item.value)
  return values.length ? [0, ...values] : [0]
})

const trendPoints = computed(() => {
  const values = trendValues.value
  const width = 520
  const height = 150
  return values.map((value, index) => {
    const x = values.length === 1 ? 0 : index / (values.length - 1) * width
    const y = height - value / 100 * height
    return `${x},${y}`
  }).join(' ')
})

const reviewText = computed(() => {
  if (!growth.completedUnitCount) return '尚未形成学习数据。完成首个微单元后，系统会生成达成度快照与个性化复盘。'
  if (growth.achievement >= 85) return '当前达成度优秀，可以加速路径并尝试更高难度的项目挑战。'
  if (growth.achievement >= 60) return '整体进展稳定。继续补齐低于 60% 的能力域，并在里程碑前安排一次综合复习。'
  return `当前处于基础积累期，${weakDomain.value?.name || '核心能力'}仍需补强。建议先完成推荐微单元，再通过快测验证掌握。`
})

function scoreTone(p) {
  if (p >= 80) return 'success'
  if (p >= 60) return 'warning'
  return 'danger'
}

function confirmReset() {
  if (window.confirm('确定要清空目标、微单元、积分、课程进度、笔记与测验记录吗？此操作不可恢复。')) {
    learning.resetAll()
    growth.reset()
  }
}
</script>

<template>
  <div class="page stats-page">
    <header class="page-header">
      <span class="eyebrow">数据评估</span>
      <h1>目标达成与学习统计</h1>
      <p>从目标、能力域、里程碑和学习行为多层评估成长，并给出下一步行动。</p>
    </header>

    <!-- 目标达成总览 -->
    <section v-if="growth.hasGoal" class="achievement-hero card card--ink">
      <div class="achievement-hero__score">
        <ProgressRing :percent="growth.achievement" :size="140" :stroke="12" color="#a071ff" track="rgba(255,255,255,.12)" />
        <div>
          <span class="tag tag--glass"><Icon name="flag" :size="12" /> 当前目标</span>
          <h2>{{ growth.goal.name }}</h2>
          <p>达标线 75% · 已完成 {{ growth.completedUnitCount }}/{{ microUnits.length }} 个微单元</p>
        </div>
      </div>
      <div class="achievement-hero__metrics">
        <div><Icon name="medal" :size="18" /><strong class="num">{{ growth.points }}</strong><span>成长积分</span></div>
        <div><Icon name="fire" :size="18" /><strong class="num">{{ growth.streak }}</strong><span>连续学习天</span></div>
        <div><Icon name="trophy" :size="18" /><strong class="num">{{ growth.badges.length }}</strong><span>获得勋章</span></div>
      </div>
    </section>

    <section v-if="growth.hasGoal" class="assessment-grid">
      <div class="card panel">
        <div class="panel-head">
          <div><h3>能力雷达</h3><p>虚线为目标能力，色块为当前能力</p></div>
          <router-link to="/path">查看路径 →</router-link>
        </div>
        <AbilityRadar :items="growth.competencyProgress" :size="290" />
      </div>

      <div class="card panel">
        <div class="panel-head">
          <div><h3>能力域达成度</h3><p>按知识点掌握度和能力权重聚合</p></div>
        </div>
        <div v-for="domain in growth.competencyProgress" :key="domain.id" class="domain-row">
          <div class="domain-row__head">
            <span class="domain-dot" :style="{ background: domain.color }"></span>
            <strong>{{ domain.name }}</strong>
            <small>权重 {{ domain.weight }}%</small>
            <b class="num" :class="{ weak: domain.progress < 60, good: domain.progress >= 85 }">{{ domain.progress }}%</b>
          </div>
          <div class="progress progress--thin"><div :style="{ width: domain.progress + '%', background: domain.color }"></div></div>
        </div>
        <div class="formula">
          <Icon name="layers" :size="14" />
          <strong>掌握度模型</strong>
          <span>微单元完成 40% + 即时快测 60%</span>
        </div>
      </div>
    </section>

    <section v-if="growth.hasGoal" class="card panel trend-block">
      <div class="panel-head">
        <div><h3>目标达成度趋势</h3><p>每次完成微单元后自动记录快照</p></div>
        <strong class="trend-value num">{{ growth.achievement }}%</strong>
      </div>
      <div class="trend-chart">
        <div v-for="level in [100, 75, 50, 25, 0]" :key="level" class="trend-grid" :class="{ 'trend-grid--target': level === 75 }" :style="{ top: (100 - level) + '%' }">
          <span>{{ level }}%</span>
        </div>
        <svg viewBox="0 0 520 150" preserveAspectRatio="none">
          <defs>
            <linearGradient id="trendArea" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stop-color="#6b5cff" stop-opacity=".3" />
              <stop offset="100%" stop-color="#6b5cff" stop-opacity="0" />
            </linearGradient>
            <linearGradient id="trendLine" x1="0" y1="0" x2="1" y2="0">
              <stop offset="0%" stop-color="#6b5cff" />
              <stop offset="100%" stop-color="#a071ff" />
            </linearGradient>
          </defs>
          <polyline :points="`0,150 ${trendPoints} 520,150`" fill="url(#trendArea)" stroke="none" />
          <polyline :points="trendPoints" fill="none" stroke="url(#trendLine)" stroke-width="3" stroke-linejoin="round" stroke-linecap="round" vector-effect="non-scaling-stroke" />
        </svg>
        <div v-if="trendValues.length === 1" class="trend-empty">完成微单元后将在这里看到成长曲线</div>
      </div>
      <div class="trend-labels"><span>目标建立</span><span>当前</span></div>
    </section>

    <section v-if="growth.hasGoal" class="assessment-grid feedback-grid">
      <div class="card panel">
        <div class="panel-head">
          <div><h3>里程碑状态</h3><p>阶段验收确保学习不偏离目标</p></div>
        </div>
        <div class="milestone-list">
          <div v-for="milestone in milestoneRows" :key="milestone.id" :class="{ done: milestone.done, current: milestone.current }">
            <span class="milestone-list__mark">
              <Icon v-if="milestone.done" name="check" :size="13" :stroke="3" />
              <template v-else>{{ milestone.week }}</template>
            </span>
            <div><strong>{{ milestone.name }}</strong><small>第 {{ milestone.week }} 周 · {{ milestone.standard }}</small></div>
            <b>{{ milestone.done ? '已达成' : milestone.current ? '进行中' : '待开始' }}</b>
          </div>
        </div>
      </div>

      <div class="card panel review-card">
        <div class="panel-head">
          <div><h3>AI 阶段复盘</h3><p>基于当前达成度动态生成</p></div>
          <span class="tag tag--ai"><Icon name="sparkles" :size="12" /> AI</span>
        </div>
        <p class="review-text">{{ reviewText }}</p>
        <div v-if="weakDomain" class="weak-card">
          <span class="weak-card__label"><Icon name="alert" :size="13" /> 优先补强</span>
          <strong>{{ weakDomain.name }}</strong>
          <b class="num">{{ weakDomain.progress }}%</b>
        </div>
        <router-link v-if="reinforcementUnit" :to="`/micro/${reinforcementUnit.id}`" class="recommend-action">
          <span><small>推荐下一步 · {{ reinforcementUnit.duration }} 分钟</small><strong>{{ reinforcementUnit.title }}</strong></span>
          <b>开始补强 <Icon name="arrowRight" :size="13" /></b>
        </router-link>
        <router-link v-else to="/path" class="recommend-action">
          <span><small>路径建议</small><strong>查看下一阶段挑战</strong></span>
          <b>查看路径 <Icon name="arrowRight" :size="13" /></b>
        </router-link>
      </div>
    </section>

    <section v-else class="no-goal card">
      <span class="no-goal__icon"><Icon name="compass" :size="28" /></span>
      <div><h2>先设定职业目标，才能评估达成度</h2><p>系统会基于能力图谱、学习行为和快测成绩生成成长报告。</p></div>
      <router-link to="/career" class="btn btn--primary">开始职业规划</router-link>
    </section>

    <!-- 课程学习统计 -->
    <section class="overview">
      <div class="overview__ring card">
        <ProgressRing :percent="learning.overallProgress" :size="118" :stroke="11" />
        <div>
          <div class="overview__big num">{{ learning.completedCount }} <small>/ {{ totalChapterCount }}</small></div>
          <div class="overview__label">章节已完成</div>
          <div class="overview__sub">{{ courses.length }} 门课程 · {{ quizRows.filter(r => r.result).length }} 次测验</div>
        </div>
      </div>

      <div class="overview__chart card panel">
        <div class="panel-head">
          <div><h3>近 7 天学习活跃度</h3><p>章节访问 + 微单元完成</p></div>
          <span class="tag tag--neutral">本周 {{ weeklyTotal }} 次</span>
        </div>
        <div class="chart">
          <div v-for="day in weekly" :key="day.date" class="chart__col">
            <span class="chart__count num" v-if="day.count">{{ day.count }}</span>
            <div
              class="chart__bar"
              :class="{ 'chart__bar--empty': !day.count }"
              :style="{ height: Math.max(6, (day.count / weeklyMax) * 100) + '%' }"
            ></div>
            <span class="chart__label">{{ day.label }}</span>
          </div>
        </div>
      </div>
    </section>

    <section class="card panel">
      <div class="panel-head">
        <div><h3>各课程进度</h3><p>按已完成章节数计算</p></div>
      </div>
      <div v-for="course in courses" :key="course.id" class="course-row">
        <span class="course-row__icon">{{ course.icon }}</span>
        <router-link :to="`/courses/${course.id}`" class="course-row__title">
          {{ course.title }}
        </router-link>
        <div class="progress progress--thin course-row__track">
          <div :style="{ width: learning.courseProgress(course.id) + '%', background: course.accent }"></div>
        </div>
        <span class="course-row__value num">
          {{ learning.courseCompletedCount(course.id) }}/{{ course.chapters.length }}
        </span>
      </div>
    </section>

    <section class="card panel">
      <div class="panel-head">
        <div><h3>测验成绩</h3><p>每门课程的最近一次成绩</p></div>
      </div>
      <table class="quiz-table">
        <thead>
          <tr>
            <th>测验</th>
            <th>题量</th>
            <th>得分</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in quizRows" :key="row.id">
            <td class="quiz-table__title">{{ row.title }}</td>
            <td>{{ row.total }} 题</td>
            <td>
              <span v-if="row.percent !== null" class="tag" :class="`tag--${scoreTone(row.percent)}`">{{ row.percent }} 分</span>
              <span v-else class="tag tag--neutral">未测验</span>
            </td>
            <td class="quiz-table__action">
              <router-link :to="`/quiz/${row.id}`" class="panel-link">
                {{ row.result ? '重测' : '去测验' }} →
              </router-link>
            </td>
          </tr>
        </tbody>
      </table>
    </section>

    <section class="card panel">
      <div class="panel-head">
        <div><h3>我的笔记</h3><p>在章节学习页底部可以随手记录你的理解与思考</p></div>
        <span class="tag tag--neutral">{{ noteEntries.length }} 条</span>
      </div>
      <p v-if="noteEntries.length === 0" class="muted">还没有笔记。</p>
      <div v-for="entry in noteEntries" :key="entry.key" class="note">
        <router-link :to="`/courses/${entry.courseId}/${entry.chapterId}`" class="note__source">
          <Icon name="note" :size="13" /> {{ entry.course.title }} · {{ entry.chapter.title }}
        </router-link>
        <p class="note__text">{{ entry.text }}</p>
      </div>
    </section>

    <section class="danger-zone">
      <button class="btn btn--danger" @click="confirmReset"><Icon name="trash" :size="15" /> 清空全部学习数据</button>
    </section>
  </div>
</template>

<style scoped>
.stats-page > .card { margin-bottom: 18px; }

.achievement-hero {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  padding: 28px 34px;
  border-radius: var(--radius-lg);
}

.achievement-hero__score { display: flex; align-items: center; gap: 24px; }
.achievement-hero__score :deep(.progress-ring__text) { fill: #fff; }
.achievement-hero__score h2 { margin: 12px 0 4px; font-size: 22px; }
.achievement-hero__score p { margin: 0; color: rgba(201, 203, 230, 0.75); font-size: 13px; }

.achievement-hero__metrics { display: flex; gap: 12px; }
.achievement-hero__metrics div {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  min-width: 108px;
  padding: 16px 14px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.07);
  border: 1px solid rgba(255, 255, 255, 0.1);
}
.achievement-hero__metrics .icon { color: #c7bfff; margin-bottom: 4px; }
.achievement-hero__metrics strong { font-size: 24px; line-height: 1.1; }
.achievement-hero__metrics span { color: rgba(201, 203, 230, 0.7); font-size: 11.5px; }

.assessment-grid { display: grid; grid-template-columns: 0.9fr 1.1fr; gap: 18px; margin-bottom: 18px; }
.feedback-grid { grid-template-columns: 1fr 1fr; }

.domain-row { margin: 16px 0; }
.domain-row__head { display: flex; align-items: center; gap: 8px; margin-bottom: 7px; font-size: 13.5px; }
.domain-dot { width: 9px; height: 9px; border-radius: 50%; }
.domain-row small { color: var(--text-3); }
.domain-row b { margin-left: auto; font-weight: 700; }
.domain-row b.weak { color: var(--danger); }
.domain-row b.good { color: var(--success-strong); }
.formula { display: flex; align-items: center; gap: 8px; margin-top: 18px; padding: 11px 14px; border-radius: 12px; background: var(--surface-2); color: var(--text-2); font-size: 12.5px; }
.formula .icon { color: var(--primary); }
.formula span { margin-left: auto; }

.trend-block { padding-bottom: 18px; }
.trend-value { color: var(--primary-strong); font-size: 26px; }
.trend-chart { position: relative; height: 170px; margin: 16px 8px 0 40px; }
.trend-chart svg { position: absolute; inset: 0; width: 100%; height: 150px; overflow: visible; }
.trend-grid { position: absolute; left: 0; right: 0; height: 1px; background: var(--border); }
.trend-grid--target { background: none; border-top: 1px dashed var(--success); }
.trend-grid span { position: absolute; right: calc(100% + 10px); top: -9px; color: var(--text-3); font-size: 11px; }
.trend-grid--target span { color: var(--success-strong); font-weight: 700; }
.trend-empty { position: absolute; inset: 55px 0 auto; color: var(--text-3); font-size: 13px; text-align: center; }
.trend-labels { display: flex; justify-content: space-between; margin: -8px 8px 0 48px; color: var(--text-3); font-size: 11px; }

.milestone-list > div { display: flex; align-items: center; gap: 12px; padding: 12px 0; }
.milestone-list > div + div { border-top: 1px solid var(--border); }
.milestone-list__mark { display: grid; place-items: center; width: 32px; height: 32px; flex: 0 0 32px; border-radius: 10px; background: var(--surface-3); color: var(--text-3); font-size: 12px; font-weight: 700; }
.milestone-list > div.done .milestone-list__mark { background: var(--success); color: white; }
.milestone-list > div.current .milestone-list__mark { background: linear-gradient(135deg, var(--primary), var(--primary-strong)); color: #fff; box-shadow: var(--shadow-primary); }
.milestone-list div div { display: flex; flex: 1; min-width: 0; flex-direction: column; }
.milestone-list strong { font-size: 13.5px; }
.milestone-list small { overflow: hidden; color: var(--text-3); font-size: 12px; text-overflow: ellipsis; white-space: nowrap; }
.milestone-list b { padding: 3px 9px; border-radius: 99px; background: var(--surface-3); color: var(--text-3); font-size: 11.5px; font-weight: 600; }
.milestone-list .done b { background: var(--success-soft); color: var(--success-strong); }
.milestone-list .current b { background: var(--primary-soft); color: var(--primary-strong); }

.review-card { background: linear-gradient(150deg, #f6f3ff 0%, #fff 60%); }
.review-text { min-height: 68px; margin: 0 0 14px; color: var(--text-2); font-size: 14px; }
.weak-card { display: flex; align-items: center; gap: 10px; padding: 12px 14px; margin-bottom: 10px; border: 1px solid #fde1b8; border-radius: 12px; background: var(--warning-soft); }
.weak-card__label { display: inline-flex; align-items: center; gap: 4px; color: var(--warning-strong); font-size: 12px; font-weight: 600; }
.weak-card strong { flex: 1; font-size: 14px; }
.weak-card b { color: var(--danger); font-size: 15px; }
.recommend-action {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 13px 16px;
  border: 1px solid var(--primary-soft-2);
  border-radius: 13px;
  background: var(--surface);
  color: var(--text);
  transition: all var(--t-fast);
}
.recommend-action:hover { border-color: var(--primary); box-shadow: var(--shadow-sm); }
.recommend-action span { display: flex; flex-direction: column; }
.recommend-action small { color: var(--text-3); font-size: 12px; }
.recommend-action strong { font-size: 13.5px; }
.recommend-action b { display: inline-flex; align-items: center; gap: 3px; color: var(--primary); font-size: 13px; white-space: nowrap; }

.no-goal { display: flex; align-items: center; gap: 18px; padding: 24px 28px; }
.no-goal__icon { display: grid; place-items: center; width: 56px; height: 56px; flex: 0 0 56px; border-radius: 18px; background: var(--primary-soft); color: var(--primary-strong); }
.no-goal div { flex: 1; }
.no-goal h2 { margin: 0; font-size: 18px; }
.no-goal p { margin: 3px 0 0; color: var(--text-2); font-size: 13.5px; }

.overview { display: grid; grid-template-columns: 320px 1fr; gap: 18px; margin-bottom: 18px; }
.overview__ring { display: flex; align-items: center; gap: 22px; padding: 26px; }
.overview__big { font-size: 30px; font-weight: 700; line-height: 1.1; }
.overview__big small { color: var(--text-3); font-size: 16px; font-weight: 600; }
.overview__label { margin-top: 4px; font-size: 14px; font-weight: 600; color: var(--text-2); }
.overview__sub { color: var(--text-3); font-size: 12px; }

.chart { display: flex; align-items: flex-end; gap: 12px; height: 120px; padding-top: 18px; }
.chart__col { flex: 1; display: flex; flex-direction: column; align-items: center; justify-content: flex-end; height: 100%; gap: 6px; }
.chart__count { font-size: 12px; font-weight: 700; color: var(--primary-strong); }
.chart__bar { width: 100%; max-width: 38px; background: linear-gradient(180deg, #a071ff, var(--primary)); border-radius: 8px 8px 4px 4px; transition: height 0.5s var(--ease); box-shadow: 0 6px 12px var(--primary-glow); }
.chart__bar--empty { background: var(--surface-3); box-shadow: none; }
.chart__label { font-size: 11.5px; color: var(--text-3); }

.course-row { display: flex; align-items: center; gap: 14px; padding: 12px 0; }
.course-row + .course-row { border-top: 1px solid var(--border); }
.course-row__icon { display: grid; place-items: center; width: 36px; height: 36px; border-radius: 11px; background: var(--surface-3); font-size: 18px; }
.course-row__title { width: 260px; min-width: 150px; font-size: 14px; font-weight: 500; color: var(--text); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.course-row__title:hover { color: var(--primary); }
.course-row__track { flex: 1; }
.course-row__value { font-size: 13px; color: var(--text-2); min-width: 44px; text-align: right; font-weight: 600; }

.quiz-table { width: 100%; border-collapse: collapse; font-size: 14px; }
.quiz-table th, .quiz-table td { text-align: left; padding: 12px 10px; border-bottom: 1px solid var(--border); }
.quiz-table th { font-size: 12px; color: var(--text-3); font-weight: 600; letter-spacing: 0.04em; text-transform: uppercase; }
.quiz-table tbody tr:last-child td { border-bottom: none; }
.quiz-table tbody tr:hover td { background: var(--surface-2); }
.quiz-table__title { font-weight: 500; }
.quiz-table__action { text-align: right; }

.muted { color: var(--text-3); font-size: 13.5px; margin: 0; }
.note { padding: 14px 0; }
.note + .note { border-top: 1px solid var(--border); }
.note__source { display: inline-flex; align-items: center; gap: 5px; font-size: 12.5px; font-weight: 600; color: var(--text-2); }
.note__source:hover { color: var(--primary); }
.note__text { margin: 6px 0 0; padding: 12px 14px; border-radius: 10px; background: var(--surface-2); font-size: 14px; white-space: pre-wrap; }

.danger-zone { text-align: center; padding: 8px 0 20px; }

@media (max-width: 860px) {
  .achievement-hero, .achievement-hero__score { align-items: flex-start; flex-direction: column; }
  .achievement-hero__metrics { width: 100%; }
  .achievement-hero__metrics div { flex: 1; min-width: 0; }
  .assessment-grid { grid-template-columns: 1fr; }
  .no-goal { align-items: flex-start; flex-wrap: wrap; }
  .no-goal div { min-width: calc(100% - 80px); }
  .overview { grid-template-columns: 1fr; }
  .course-row__title { width: auto; flex: 1; }
  .course-row__track { display: none; }
}
</style>
