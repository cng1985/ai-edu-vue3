<script setup>
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { careers, frontendPath } from '../data/careerPath'
import { useGrowthStore } from '../stores/growth'

const router = useRouter()
const growth = useGrowthStore()
const step = ref(growth.hasGoal ? 3 : 1)
const selectedId = ref(growth.goal?.careerId || 'frontend')
const baseLevel = ref(growth.goal?.baseLevel || '零基础')
const weeklyHours = ref(growth.goal?.weeklyHours || 12)
const durationWeeks = ref(growth.goal?.durationWeeks || 16)

const selected = computed(() => careers.find((item) => item.id === selectedId.value))
const canCreate = computed(() => selectedId.value === 'frontend')

function selectCareer(id) {
  selectedId.value = id
}

function confirmGoal() {
  growth.createGoal({
    careerId: selectedId.value,
    baseLevel: baseLevel.value,
    weeklyHours: weeklyHours.value,
    durationWeeks: durationWeeks.value
  })
  step.value = 3
}
</script>

<template>
  <div class="page career-page">
    <header class="page-header">
      <span class="eyebrow">AI 职业规划</span>
      <h1>选择你想成为的人</h1>
      <p>用兴趣、基础和可投入时间确定方向，再生成一条可执行、可评估的学习路径。</p>
    </header>

    <div class="steps">
      <div v-for="item in 3" :key="item" class="step" :class="{ active: step >= item }">
        <span>{{ item }}</span>
        {{ ['选择职业', '确认目标', '查看分解'][item - 1] }}
      </div>
    </div>

    <section v-if="step === 1" class="fade-up">
      <div class="ai-tip card">
        <span class="ai-tip__icon">✨</span>
        <div>
          <strong>AI 推荐：Web 前端工程师（匹配度 92%）</strong>
          <p>适合希望快速看到成果、通过作品进入 IT 行业的学习者。目前 MVP 已为该方向准备完整路径。</p>
        </div>
      </div>
      <div class="career-grid">
        <button
          v-for="career in careers"
          :key="career.id"
          class="career-card card"
          :class="{ selected: selectedId === career.id }"
          @click="selectCareer(career.id)"
        >
          <div class="career-card__top">
            <span class="career-icon">{{ career.icon }}</span>
            <span class="match">{{ career.match }}% 匹配</span>
          </div>
          <small>{{ career.category }}</small>
          <h2>{{ career.name }}</h2>
          <p>{{ career.description }}</p>
          <div class="career-meta">
            <span>难度 {{ '★'.repeat(career.difficulty) }}</span>
            <span>需求 {{ career.demand }}</span>
            <span>{{ career.salary }}</span>
          </div>
        </button>
      </div>
      <div v-if="selected" class="selection-actions">
        <span v-if="!canCreate">该方向路径正在建设中，可先预览；当前可创建前端目标。</span>
        <button class="btn btn--primary" :disabled="!canCreate" @click="step = 2">
          以“{{ selected.name }}”为目标 →
        </button>
      </div>
    </section>

    <section v-else-if="step === 2" class="goal-layout fade-up">
      <div class="goal-form card">
        <h2>确认你的目标条件</h2>
        <label>
          当前基础
          <select v-model="baseLevel">
            <option>零基础</option><option>入门</option><option>进阶</option><option>熟练</option>
          </select>
        </label>
        <label>
          每周可投入时间
          <div class="range-row">
            <input v-model.number="weeklyHours" type="range" min="4" max="30" />
            <strong>{{ weeklyHours }} 小时</strong>
          </div>
        </label>
        <label>
          目标周期
          <select v-model.number="durationWeeks">
            <option :value="12">12 周 · 冲刺</option>
            <option :value="16">16 周 · 标准</option>
            <option :value="24">24 周 · 稳健</option>
          </select>
        </label>
        <div class="form-actions">
          <button class="btn btn--ghost" @click="step = 1">返回</button>
          <button class="btn btn--primary" @click="confirmGoal">生成学习路径</button>
        </div>
      </div>
      <div class="commitment card">
        <span class="tag">目标承诺书</span>
        <h2>{{ durationWeeks }} 周成为初级前端工程师</h2>
        <p>我将每周投入 <strong>{{ weeklyHours }} 小时</strong>，通过微单元、快测和项目里程碑持续验证能力。</p>
        <ul>
          <li>难度评估：中等</li>
          <li>建议节奏：每日 25–40 分钟</li>
          <li>预计微单元：8 个 MVP 单元</li>
          <li>验收标准：目标达成度 ≥ 75%</li>
        </ul>
      </div>
    </section>

    <section v-else class="path-panel card fade-up">
      <div class="path-head">
        <div>
          <span class="tag">AI 已完成目标分解</span>
          <h2>{{ frontendPath.name }}</h2>
          <p>4 个能力域 · 8 个知识点 · 4 个里程碑</p>
        </div>
        <button class="btn btn--primary" @click="router.push('/')">进入学习驾驶舱</button>
      </div>
      <div class="domain-list">
        <article v-for="(domain, index) in frontendPath.competencies" :key="domain.id" class="domain">
          <div class="domain__number">{{ index + 1 }}</div>
          <div class="domain__body">
            <div class="domain__title">
              <strong>{{ domain.name }}</strong>
              <span>权重 {{ domain.weight }}%</span>
            </div>
            <div class="point-list">
              <span v-for="point in domain.points" :key="point.id">● {{ point.name }}</span>
            </div>
          </div>
        </article>
      </div>
      <div class="milestones">
        <div v-for="milestone in frontendPath.milestones" :key="milestone.id">
          <strong>第 {{ milestone.week }} 周</strong>
          <span>{{ milestone.name }}</span>
          <small>{{ milestone.standard }}</small>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.career-page { max-width: 1120px; }
.eyebrow { color: var(--primary); font-size: 13px; font-weight: 700; letter-spacing: .08em; }
.steps { display: flex; gap: 8px; margin: 24px 0; }
.step { flex: 1; display: flex; align-items: center; gap: 8px; padding: 10px 14px; border-radius: 10px; background: #e9edf4; color: var(--text-3); font-size: 13px; }
.step span { display: grid; place-items: center; width: 24px; height: 24px; border-radius: 50%; background: #cbd5e1; color: white; font-weight: 700; }
.step.active { background: var(--primary-soft); color: var(--primary-strong); }
.step.active span { background: var(--primary); }
.ai-tip { display: flex; gap: 14px; padding: 18px 20px; margin-bottom: 18px; background: linear-gradient(100deg, #eef2ff, #faf5ff); }
.ai-tip__icon { font-size: 25px; }
.ai-tip p { margin: 3px 0 0; color: var(--text-2); font-size: 13px; }
.career-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; }
.career-card { text-align: left; padding: 22px; cursor: pointer; font: inherit; color: inherit; transition: .2s; }
.career-card:hover, .career-card.selected { transform: translateY(-2px); border-color: var(--primary); box-shadow: 0 10px 30px rgba(99, 102, 241, .12); }
.career-card__top, .career-meta, .domain__title, .path-head { display: flex; justify-content: space-between; align-items: center; }
.career-icon { font-size: 28px; }
.match { color: var(--success); background: var(--success-soft); border-radius: 99px; padding: 3px 9px; font-size: 12px; font-weight: 700; }
.career-card small { color: var(--text-3); }
.career-card h2 { margin: 5px 0 8px; font-size: 19px; }
.career-card p { min-height: 52px; color: var(--text-2); font-size: 13px; }
.career-meta { font-size: 11px; color: var(--text-2); }
.selection-actions { display: flex; justify-content: flex-end; align-items: center; gap: 16px; margin-top: 20px; color: var(--warning); font-size: 13px; }
.goal-layout { display: grid; grid-template-columns: 1.1fr .9fr; gap: 20px; }
.goal-form, .commitment { padding: 26px; }
.goal-form h2, .commitment h2 { margin-top: 0; }
.goal-form label { display: block; margin: 18px 0; color: var(--text-2); font-size: 13px; font-weight: 600; }
select { width: 100%; margin-top: 7px; padding: 11px; border: 1px solid var(--border); border-radius: 9px; background: white; font: inherit; }
.range-row { display: flex; align-items: center; gap: 14px; margin-top: 8px; }
.range-row input { flex: 1; }
.range-row strong { min-width: 60px; color: var(--primary); }
.form-actions { display: flex; justify-content: flex-end; gap: 10px; }
.commitment { background: linear-gradient(145deg, #312e81, #6366f1); color: white; }
.commitment .tag { background: rgba(255,255,255,.16); color: white; }
.commitment p { color: #e0e7ff; }
.commitment li { margin: 10px 0; }
.path-panel { padding: 28px; }
.path-head { margin-bottom: 24px; }
.path-head h2 { margin: 8px 0 2px; }
.path-head p { margin: 0; color: var(--text-2); }
.domain-list { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.domain { display: flex; gap: 14px; padding: 16px; border: 1px solid var(--border); border-radius: 12px; }
.domain__number { display: grid; place-items: center; width: 34px; height: 34px; flex: 0 0 34px; border-radius: 10px; background: var(--primary-soft); color: var(--primary); font-weight: 800; }
.domain__body { flex: 1; }
.domain__title span { color: var(--text-3); font-size: 12px; }
.point-list { display: flex; flex-direction: column; margin-top: 8px; color: var(--text-2); font-size: 13px; }
.point-list span::first-letter { color: var(--success); }
.milestones { display: grid; grid-template-columns: repeat(4, 1fr); gap: 2px; margin-top: 22px; }
.milestones div { display: flex; flex-direction: column; padding: 14px; border-top: 3px solid var(--primary); background: var(--surface-2); }
.milestones span { font-size: 13px; font-weight: 600; }
.milestones small { color: var(--text-3); }
@media (max-width: 800px) {
  .career-grid, .goal-layout, .domain-list { grid-template-columns: 1fr; }
  .milestones { grid-template-columns: 1fr 1fr; }
  .steps { overflow-x: auto; }
  .step { min-width: 130px; }
  .path-head { align-items: flex-start; flex-direction: column; gap: 12px; }
}
</style>
