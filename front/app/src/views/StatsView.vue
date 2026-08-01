<script setup>
import { computed } from 'vue'
import { courses, totalChapterCount } from '../data/courses'
import { quizzes } from '../data/quizzes'
import { frontendPath, microUnits } from '../data/careerPath'
import { useLearningStore } from '../stores/learning'
import { useGrowthStore } from '../stores/growth'
import ProgressRing from '../components/ProgressRing.vue'
import AbilityRadar from '../components/AbilityRadar.vue'

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

function confirmReset() {
  if (window.confirm('确定要清空目标、微单元、积分、课程进度、笔记与测验记录吗？此操作不可恢复。')) {
    learning.resetAll()
    growth.reset()
  }
}
</script>

<template>
  <div class="page">
    <header class="page-header">
      <h1>目标达成与学习统计</h1>
      <p>从目标、能力域、里程碑和学习行为多层评估成长，并给出下一步行动。</p>
    </header>

    <section v-if="growth.hasGoal" class="achievement-hero card">
      <div class="achievement-hero__score">
        <ProgressRing :percent="growth.achievement" :size="138" :stroke="11" />
        <div>
          <span class="tag">当前目标</span>
          <h2>{{ growth.goal.name }}</h2>
          <p>达标线 75% · 已完成 {{ growth.completedUnitCount }}/{{ microUnits.length }} 个微单元</p>
        </div>
      </div>
      <div class="achievement-hero__metrics">
        <div><strong>{{ growth.points }}</strong><span>成长积分</span></div>
        <div><strong>{{ growth.streak }}</strong><span>连续学习天</span></div>
        <div><strong>{{ growth.badges.length }}</strong><span>获得勋章</span></div>
      </div>
    </section>

    <section v-if="growth.hasGoal" class="assessment-grid">
      <div class="card assessment-card radar-card">
        <div class="block-head">
          <div><h3>能力雷达</h3><p>虚线为目标能力，色块为当前能力</p></div>
          <router-link to="/path">查看路径 →</router-link>
        </div>
        <AbilityRadar :items="growth.competencyProgress" :size="280" />
      </div>

      <div class="card assessment-card domains-card">
        <div class="block-head">
          <div><h3>能力域达成度</h3><p>按知识点掌握度和能力权重聚合</p></div>
        </div>
        <div v-for="domain in growth.competencyProgress" :key="domain.id" class="domain-row">
          <div>
            <span class="domain-dot" :style="{ background: domain.color }"></span>
            <strong>{{ domain.name }}</strong>
            <small>权重 {{ domain.weight }}%</small>
            <b :class="{ weak: domain.progress < 60, good: domain.progress >= 85 }">{{ domain.progress }}%</b>
          </div>
          <div class="domain-track"><div :style="{ width: domain.progress + '%', background: domain.color }"></div></div>
        </div>
        <div class="formula">
          <strong>掌握度模型</strong>
          <span>微单元完成 40% + 即时快测 60%</span>
        </div>
      </div>
    </section>

    <section v-if="growth.hasGoal" class="card block trend-block">
      <div class="block-head">
        <div><h3>目标达成度趋势</h3><p>每次完成微单元后自动记录快照</p></div>
        <strong>{{ growth.achievement }}%</strong>
      </div>
      <div class="trend-chart">
        <div v-for="level in [100, 75, 50, 25, 0]" :key="level" class="trend-grid" :style="{ top: (100 - level) + '%' }">
          <span>{{ level }}%</span>
        </div>
        <svg viewBox="0 0 520 150" preserveAspectRatio="none">
          <defs>
            <linearGradient id="trendArea" x1="0" y1="0" x2="0" y2="1">
              <stop offset="0%" stop-color="#6366f1" stop-opacity=".28" />
              <stop offset="100%" stop-color="#6366f1" stop-opacity=".02" />
            </linearGradient>
          </defs>
          <polyline :points="`0,150 ${trendPoints} 520,150`" fill="url(#trendArea)" stroke="none" />
          <polyline :points="trendPoints" fill="none" stroke="#6366f1" stroke-width="3" vector-effect="non-scaling-stroke" />
        </svg>
        <div v-if="trendValues.length === 1" class="trend-empty">完成微单元后将在这里看到成长曲线</div>
      </div>
      <div class="trend-labels"><span>目标建立</span><span>当前</span></div>
    </section>

    <section v-if="growth.hasGoal" class="assessment-grid feedback-grid">
      <div class="card assessment-card">
        <div class="block-head">
          <div><h3>里程碑状态</h3><p>阶段验收确保学习不偏离目标</p></div>
        </div>
        <div class="milestone-list">
          <div v-for="milestone in milestoneRows" :key="milestone.id" :class="{ done: milestone.done, current: milestone.current }">
            <span>{{ milestone.done ? '✓' : milestone.current ? '●' : milestone.week }}</span>
            <div><strong>{{ milestone.name }}</strong><small>第 {{ milestone.week }} 周 · {{ milestone.standard }}</small></div>
            <b>{{ milestone.done ? '已达成' : milestone.current ? '进行中' : '待开始' }}</b>
          </div>
        </div>
      </div>

      <div class="card assessment-card review-card">
        <div class="block-head">
          <div><h3>AI 阶段复盘</h3><p>基于当前达成度动态生成</p></div>
          <span class="ai-badge">AI</span>
        </div>
        <p class="review-text">{{ reviewText }}</p>
        <div v-if="weakDomain" class="weak-card">
          <span>优先补强</span>
          <strong>{{ weakDomain.name }}</strong>
          <b>{{ weakDomain.progress }}%</b>
        </div>
        <router-link v-if="reinforcementUnit" :to="`/micro/${reinforcementUnit.id}`" class="recommend-action">
          <span><small>推荐下一步 · {{ reinforcementUnit.duration }} 分钟</small><strong>{{ reinforcementUnit.title }}</strong></span>
          <b>开始补强 →</b>
        </router-link>
        <router-link v-else to="/path" class="recommend-action">
          <span><small>路径建议</small><strong>查看下一阶段挑战</strong></span>
          <b>查看路径 →</b>
        </router-link>
      </div>
    </section>

    <section v-else class="no-goal card">
      <span>🧭</span>
      <div><h2>先设定职业目标，才能评估达成度</h2><p>系统会基于能力图谱、学习行为和快测成绩生成成长报告。</p></div>
      <router-link to="/career" class="btn btn--primary">开始职业规划</router-link>
    </section>

    <section class="overview">
      <div class="overview__ring card">
        <ProgressRing :percent="learning.overallProgress" :size="120" :stroke="10" />
        <div>
          <div class="overview__big">{{ learning.completedCount }} / {{ totalChapterCount }}</div>
          <div class="overview__label">章节完成</div>
        </div>
      </div>

      <div class="overview__chart card">
        <h3>近 7 天学习活跃度</h3>
        <div class="chart">
          <div v-for="day in weekly" :key="day.date" class="chart__col">
            <span class="chart__count" v-if="day.count">{{ day.count }}</span>
            <div
              class="chart__bar"
              :style="{
                height: (day.count / weeklyMax) * 100 + '%',
                opacity: day.count ? 1 : 0.25
              }"
            ></div>
            <span class="chart__label">{{ day.label }}</span>
          </div>
        </div>
      </div>
    </section>

    <section class="card block">
      <h3>各课程进度</h3>
      <div v-for="course in courses" :key="course.id" class="course-row">
        <span class="course-row__icon">{{ course.icon }}</span>
        <router-link :to="`/courses/${course.id}`" class="course-row__title">
          {{ course.title }}
        </router-link>
        <div class="course-row__track">
          <div
            class="course-row__fill"
            :style="{
              width: learning.courseProgress(course.id) + '%',
              background: course.accent
            }"
          ></div>
        </div>
        <span class="course-row__value">
          {{ learning.courseCompletedCount(course.id) }}/{{ course.chapters.length }}
        </span>
      </div>
    </section>

    <section class="card block">
      <h3>测验成绩</h3>
      <table class="quiz-table">
        <thead>
          <tr>
            <th>测验</th>
            <th>题量</th>
            <th>得分</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in quizRows" :key="row.id">
            <td>{{ row.title }}</td>
            <td>{{ row.total }} 题</td>
            <td>
              <strong v-if="row.percent !== null" :class="{
                'score--good': row.percent >= 80,
                'score--mid': row.percent >= 60 && row.percent < 80,
                'score--bad': row.percent < 60
              }">{{ row.percent }} 分</strong>
              <span v-else class="score--none">未测验</span>
            </td>
            <td>
              <router-link :to="`/quiz/${row.id}`">
                {{ row.result ? '重测' : '去测验' }} →
              </router-link>
            </td>
          </tr>
        </tbody>
      </table>
    </section>

    <section class="card block">
      <h3>我的笔记（{{ noteEntries.length }}）</h3>
      <p v-if="noteEntries.length === 0" class="muted">
        还没有笔记。在章节学习页底部可以随手记录你的理解与思考。
      </p>
      <div v-for="entry in noteEntries" :key="entry.key" class="note">
        <router-link :to="`/courses/${entry.courseId}/${entry.chapterId}`" class="note__source">
          {{ entry.course.icon }} {{ entry.course.title }} · {{ entry.chapter.title }}
        </router-link>
        <p class="note__text">{{ entry.text }}</p>
      </div>
    </section>

    <section class="danger-zone">
      <button class="btn btn--danger" @click="confirmReset">清空全部学习数据</button>
    </section>
  </div>
</template>

<style scoped>
.achievement-hero {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
  padding: 25px 30px;
  margin-bottom: 18px;
  background: linear-gradient(110deg, #eef2ff, #fff 58%, #f5f3ff);
}

.achievement-hero__score {
  display: flex;
  align-items: center;
  gap: 20px;
}

.achievement-hero__score h2 {
  margin: 7px 0 2px;
  font-size: 20px;
}

.achievement-hero__score p,
.block-head p {
  margin: 0;
  color: var(--text-3);
  font-size: 12px;
}

.achievement-hero__metrics {
  display: flex;
  gap: 28px;
}

.achievement-hero__metrics div {
  display: flex;
  flex-direction: column;
  text-align: center;
}

.achievement-hero__metrics strong {
  color: var(--primary);
  font-size: 23px;
}

.achievement-hero__metrics span {
  color: var(--text-3);
  font-size: 11px;
}

.assessment-grid {
  display: grid;
  grid-template-columns: .9fr 1.1fr;
  gap: 18px;
  margin-bottom: 18px;
}

.assessment-card {
  padding: 22px 25px;
  box-shadow: none;
}

.block-head {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 14px;
}

.block-head h3 {
  margin: 0;
  font-size: 15.5px;
}

.block-head > strong {
  color: var(--primary);
  font-size: 24px;
}

.block-head > a {
  font-size: 12px;
}

.domain-row {
  margin: 15px 0;
}

.domain-row > div:first-child {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 6px;
  font-size: 13px;
}

.domain-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.domain-row small {
  color: var(--text-3);
}

.domain-row b {
  margin-left: auto;
}

.domain-row b.weak { color: var(--danger); }
.domain-row b.good { color: var(--success); }

.domain-track {
  height: 7px;
  overflow: hidden;
  border-radius: 8px;
  background: var(--border);
}

.domain-track div {
  height: 100%;
  border-radius: inherit;
}

.formula {
  display: flex;
  justify-content: space-between;
  padding: 9px 11px;
  border-radius: 8px;
  background: var(--surface-2);
  color: var(--text-2);
  font-size: 11px;
}

.trend-block {
  padding-bottom: 16px;
}

.trend-chart {
  position: relative;
  height: 170px;
  margin: 15px 8px 0 38px;
}

.trend-chart svg {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 150px;
  overflow: visible;
}

.trend-grid {
  position: absolute;
  left: 0;
  right: 0;
  height: 1px;
  background: #edf0f5;
}

.trend-grid span {
  position: absolute;
  right: calc(100% + 8px);
  top: -8px;
  color: var(--text-3);
  font-size: 10px;
}

.trend-empty {
  position: absolute;
  inset: 55px 0 auto;
  color: var(--text-3);
  font-size: 12px;
  text-align: center;
}

.trend-labels {
  display: flex;
  justify-content: space-between;
  margin: -10px 6px 0 46px;
  color: var(--text-3);
  font-size: 10px;
}

.feedback-grid {
  grid-template-columns: 1fr 1fr;
}

.milestone-list > div {
  display: flex;
  align-items: center;
  gap: 11px;
  padding: 10px 0;
  border-top: 1px solid var(--border);
}

.milestone-list > div:first-child {
  border-top: 0;
}

.milestone-list > div > span {
  display: grid;
  place-items: center;
  width: 29px;
  height: 29px;
  flex: 0 0 29px;
  border-radius: 50%;
  background: var(--surface-2);
  color: var(--text-3);
  font-size: 11px;
  font-weight: 700;
}

.milestone-list > div.done > span {
  background: var(--success);
  color: white;
}

.milestone-list > div.current > span {
  background: var(--primary-soft);
  color: var(--primary);
}

.milestone-list div div {
  display: flex;
  flex: 1;
  min-width: 0;
  flex-direction: column;
}

.milestone-list strong {
  font-size: 12.5px;
}

.milestone-list small {
  overflow: hidden;
  color: var(--text-3);
  font-size: 10.5px;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.milestone-list b {
  color: var(--text-3);
  font-size: 10.5px;
  font-weight: 600;
}

.milestone-list .done b { color: var(--success); }
.milestone-list .current b { color: var(--primary); }

.ai-badge {
  padding: 2px 8px;
  border-radius: 99px;
  background: linear-gradient(120deg, #6366f1, #8b5cf6);
  color: white;
  font-size: 10px;
  font-weight: 800;
}

.review-card {
  background: linear-gradient(145deg, #faf5ff, white 65%);
}

.review-text {
  min-height: 67px;
  color: var(--text-2);
  font-size: 13px;
}

.weak-card {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  margin-bottom: 9px;
  border: 1px solid #fed7aa;
  border-radius: 9px;
  background: #fff7ed;
}

.weak-card span {
  color: #c2410c;
  font-size: 10px;
}

.weak-card strong {
  flex: 1;
  font-size: 13px;
}

.weak-card b {
  color: var(--danger);
}

.recommend-action {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 11px 13px;
  border: 1px solid #ddd6fe;
  border-radius: 9px;
  background: white;
}

.recommend-action span {
  display: flex;
  flex-direction: column;
}

.recommend-action small {
  color: var(--text-3);
}

.recommend-action strong,
.recommend-action b {
  font-size: 12px;
}

.no-goal {
  display: flex;
  align-items: center;
  gap: 18px;
  padding: 24px 28px;
  margin-bottom: 18px;
}

.no-goal > span {
  font-size: 32px;
}

.no-goal div {
  flex: 1;
}

.no-goal h2 {
  margin: 0;
  font-size: 17px;
}

.no-goal p {
  margin: 2px 0 0;
  color: var(--text-2);
  font-size: 12px;
}

.overview {
  display: grid;
  grid-template-columns: 300px 1fr;
  gap: 18px;
  margin-bottom: 18px;
}

.overview__ring {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 26px;
}

.overview__big {
  font-size: 24px;
  font-weight: 800;
}

.overview__label {
  font-size: 13px;
  color: var(--text-2);
}

.overview__chart {
  padding: 20px 26px;
}

.overview__chart h3,
.block h3 {
  margin: 0 0 16px;
  font-size: 15.5px;
}

.chart {
  display: flex;
  align-items: flex-end;
  gap: 12px;
  height: 110px;
}

.chart__col {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-end;
  height: 100%;
  gap: 5px;
}

.chart__count {
  font-size: 11.5px;
  font-weight: 700;
  color: var(--primary-strong);
}

.chart__bar {
  width: 100%;
  max-width: 34px;
  min-height: 5px;
  background: linear-gradient(180deg, #a78bfa, var(--primary));
  border-radius: 6px 6px 3px 3px;
  transition: height 0.4s ease;
}

.chart__label {
  font-size: 11px;
  color: var(--text-3);
}

.block {
  padding: 22px 26px;
  margin-bottom: 18px;
}

.course-row {
  display: flex;
  align-items: center;
  gap: 13px;
  padding: 11px 0;
}

.course-row + .course-row {
  border-top: 1px solid var(--border);
}

.course-row__icon {
  font-size: 19px;
}

.course-row__title {
  width: 240px;
  min-width: 150px;
  font-size: 14px;
  font-weight: 500;
  color: var(--text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.course-row__title:hover {
  color: var(--primary);
}

.course-row__track {
  flex: 1;
  height: 8px;
  background: var(--border);
  border-radius: 999px;
  overflow: hidden;
}

.course-row__fill {
  height: 100%;
  border-radius: 999px;
  transition: width 0.4s ease;
}

.course-row__value {
  font-size: 13px;
  color: var(--text-2);
  min-width: 42px;
  text-align: right;
}

.quiz-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}

.quiz-table th,
.quiz-table td {
  text-align: left;
  padding: 10px 8px;
  border-bottom: 1px solid var(--border);
}

.quiz-table th {
  font-size: 12.5px;
  color: var(--text-3);
  font-weight: 600;
}

.score--good { color: var(--success); }
.score--mid { color: var(--warning); }
.score--bad { color: var(--danger); }
.score--none { color: var(--text-3); font-size: 13px; }

.muted {
  color: var(--text-3);
  font-size: 13.5px;
  margin: 0;
}

.note {
  padding: 13px 0;
}

.note + .note {
  border-top: 1px solid var(--border);
}

.note__source {
  font-size: 12.5px;
  font-weight: 600;
  color: var(--text-2);
}

.note__source:hover {
  color: var(--primary);
}

.note__text {
  margin: 6px 0 0;
  font-size: 14px;
  white-space: pre-wrap;
}

.danger-zone {
  text-align: center;
  padding: 8px 0 20px;
}

@media (max-width: 780px) {
  .achievement-hero,
  .achievement-hero__score {
    align-items: flex-start;
    flex-direction: column;
  }

  .achievement-hero__metrics {
    width: 100%;
    justify-content: space-around;
  }

  .assessment-grid {
    grid-template-columns: 1fr;
  }

  .no-goal {
    align-items: flex-start;
    flex-wrap: wrap;
  }

  .no-goal div {
    min-width: calc(100% - 60px);
  }

  .overview {
    grid-template-columns: 1fr;
  }

  .course-row__title {
    width: auto;
    flex: 1;
  }

  .course-row__track {
    display: none;
  }
}
</style>
