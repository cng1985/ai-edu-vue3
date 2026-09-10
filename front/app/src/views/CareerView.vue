<script setup>
import { computed, ref } from 'vue'
import { useRouter } from 'vue-router'
import { careers, frontendPath } from '../data/careerPath'
import { useGrowthStore } from '../stores/growth'
import { aiApi } from '../api'
import MarkdownRenderer from '../components/MarkdownRenderer.vue'
import Icon from '../components/Icon.vue'

const router = useRouter()
const growth = useGrowthStore()
const step = ref(growth.hasGoal ? 3 : 1)
const selectedId = ref(growth.goal?.careerId || 'frontend')
const baseLevel = ref(growth.goal?.baseLevel || '零基础')
const weeklyHours = ref(growth.goal?.weeklyHours || 12)
const durationWeeks = ref(growth.goal?.durationWeeks || 16)
const decomposing = ref(false)
const decomposeResult = ref(growth.goal?.aiDecompose || null)

const showInterview = ref(false)
const interviewInput = ref('')
const interviewMessages = ref([])
const interviewGenerating = ref(false)
let cancelInterview = null

const steps = ['选择职业', '确认目标', '查看分解']
const selected = computed(() => careers.find((item) => item.id === selectedId.value))
const canCreate = computed(() => selectedId.value === 'frontend')

function selectCareer(id) {
  selectedId.value = id
}

function buildInterviewHistory() {
  return interviewMessages.value
    .filter((m) => m.text)
    .map((m) => ({ role: m.role, content: m.text }))
}

function sendInterview() {
  const text = interviewInput.value.trim()
  if (!text || interviewGenerating.value) return
  interviewMessages.value.push({ role: 'user', text })
  const reply = { role: 'assistant', text: '', streaming: true }
  interviewMessages.value.push(reply)
  interviewInput.value = ''
  interviewGenerating.value = true

  cancelInterview = aiApi.careerInterview(text, buildInterviewHistory().slice(0, -1), {
    onToken: (chunk) => { reply.text += chunk },
    onDone: (result) => {
      reply.text = result.text || reply.text
      reply.streaming = false
      interviewGenerating.value = false
      cancelInterview = null
    },
    onError: () => {
      reply.text = '暂时无法连接 AI 服务，你可以直接选择职业方向继续。'
      reply.streaming = false
      interviewGenerating.value = false
    }
  })
}

async function confirmGoal() {
  decomposing.value = true
  try {
    decomposeResult.value = await aiApi.goalDecompose({
      careerId: selectedId.value,
      careerName: selected.value?.name || '初级前端工程师',
      baseLevel: baseLevel.value,
      weeklyHours: weeklyHours.value,
      durationWeeks: durationWeeks.value
    })
  } catch {
    decomposeResult.value = null
  } finally {
    decomposing.value = false
  }
  growth.createGoal({
    careerId: selectedId.value,
    baseLevel: baseLevel.value,
    weeklyHours: weeklyHours.value,
    durationWeeks: durationWeeks.value,
    aiDecompose: decomposeResult.value
  })
  step.value = 3
}
</script>

<template>
  <div class="page career-page">
    <header class="page-header">
      <span class="eyebrow">AI 职业规划</span>
      <h1>选择你想成为的人</h1>
      <p>通过 AI 职业访谈了解你的方向，再生成可执行、可评估的学习路径。</p>
    </header>

    <ol class="steps">
      <li
        v-for="(label, index) in steps"
        :key="label"
        class="step"
        :class="{ 'step--active': step === index + 1, 'step--done': step > index + 1 }"
      >
        <span class="step__index">
          <Icon v-if="step > index + 1" name="check" :size="13" :stroke="3" />
          <template v-else>{{ index + 1 }}</template>
        </span>
        <span class="step__label">{{ label }}</span>
      </li>
    </ol>

    <!-- 第一步：选择职业 -->
    <section v-if="step === 1" class="fade-up">
      <div class="ai-tip card">
        <span class="ai-tip__icon"><Icon name="sparkles" :size="22" /></span>
        <div class="ai-tip__body">
          <strong>AI 推荐：Web 前端工程师 <span class="tag tag--success">匹配度 92%</span></strong>
          <p>适合希望快速看到成果、通过作品进入 IT 行业的学习者。目前 MVP 已为该方向准备完整路径。</p>
        </div>
        <button class="btn btn--soft" @click="showInterview = !showInterview">
          <Icon name="message" :size="15" />
          {{ showInterview ? '收起职业访谈' : '开始 AI 职业访谈' }}
        </button>
      </div>

      <div v-if="showInterview" class="interview card fade-up">
        <div class="interview__messages">
          <p v-if="!interviewMessages.length" class="interview__hint">
            <Icon name="sparkles" :size="15" /> 告诉 AI 你的背景、兴趣和可投入时间，它会帮你分析适合的职业方向。
          </p>
          <div
            v-for="(msg, idx) in interviewMessages"
            :key="idx"
            class="interview__msg"
            :class="`interview__msg--${msg.role}`"
          >
            <MarkdownRenderer v-if="msg.role === 'assistant'" :source="msg.text" :live="msg.streaming" />
            <template v-else>{{ msg.text }}</template>
          </div>
        </div>
        <div class="interview__composer">
          <input
            v-model="interviewInput"
            class="input"
            placeholder="例如：在职测试，想转前端，每周 15 小时…"
            @keydown.enter="sendInterview"
          />
          <button class="btn btn--primary" :disabled="interviewGenerating || !interviewInput.trim()" @click="sendInterview">
            <Icon name="send" :size="15" /> 发送
          </button>
        </div>
      </div>

      <div class="career-grid stagger">
        <button
          v-for="career in careers"
          :key="career.id"
          class="career-card card card--hover"
          :class="{ 'career-card--selected': selectedId === career.id }"
          @click="selectCareer(career.id)"
        >
          <div class="career-card__top">
            <span class="career-card__icon">{{ career.icon }}</span>
            <span class="tag tag--success">{{ career.match }}% 匹配</span>
          </div>
          <small class="career-card__category">{{ career.category }}</small>
          <h2>{{ career.name }}</h2>
          <p>{{ career.description }}</p>
          <div class="career-meta">
            <span><Icon name="star" :size="12" /> 难度 {{ '★'.repeat(career.difficulty) }}</span>
            <span><Icon name="trend" :size="12" /> 需求 {{ career.demand }}</span>
            <span>{{ career.salary }}</span>
          </div>
          <span class="career-card__check"><Icon name="check" :size="13" :stroke="3" /></span>
        </button>
      </div>

      <div v-if="selected" class="selection-actions">
        <span v-if="!canCreate" class="selection-note"><Icon name="alert" :size="14" /> 该方向路径正在建设中，可先预览；当前可创建前端目标。</span>
        <button class="btn btn--primary btn--lg" :disabled="!canCreate" @click="step = 2">
          以“{{ selected.name }}”为目标 <Icon name="arrowRight" :size="16" />
        </button>
      </div>
    </section>

    <!-- 第二步：确认条件 -->
    <section v-else-if="step === 2" class="goal-layout fade-up">
      <div class="goal-form card">
        <h2>确认你的目标条件</h2>
        <p class="goal-form__desc">这些参数决定路径节奏与里程碑排布，之后可随时调整。</p>

        <label class="field">
          <span>当前基础</span>
          <select v-model="baseLevel" class="select">
            <option>零基础</option><option>入门</option><option>进阶</option><option>熟练</option>
          </select>
        </label>

        <label class="field">
          <span>每周可投入时间</span>
          <div class="range-row">
            <input v-model.number="weeklyHours" type="range" min="4" max="30" />
            <strong class="num">{{ weeklyHours }} 小时</strong>
          </div>
        </label>

        <div class="field">
          <span>目标周期</span>
          <div class="duration-options">
            <button
              v-for="opt in [{ v: 12, t: '冲刺' }, { v: 16, t: '标准' }, { v: 24, t: '稳健' }]"
              :key="opt.v"
              type="button"
              class="duration"
              :class="{ 'duration--on': durationWeeks === opt.v }"
              @click="durationWeeks = opt.v"
            >
              <strong class="num">{{ opt.v }} 周</strong>
              <small>{{ opt.t }}</small>
            </button>
          </div>
        </div>

        <div class="form-actions">
          <button class="btn btn--ghost" @click="step = 1"><Icon name="arrowLeft" :size="15" /> 返回</button>
          <button class="btn btn--primary" :disabled="decomposing" @click="confirmGoal">
            <Icon name="sparkles" :size="15" /> {{ decomposing ? 'AI 分解中…' : '生成学习路径' }}
          </button>
        </div>
      </div>

      <div class="commitment card card--gradient">
        <span class="tag tag--glass">目标承诺书</span>
        <h2>{{ durationWeeks }} 周成为初级前端工程师</h2>
        <p>我将每周投入 <strong>{{ weeklyHours }} 小时</strong>，通过微单元、快测和项目里程碑持续验证能力。</p>
        <ul>
          <li><Icon name="check" :size="14" :stroke="3" /> 难度评估：中等</li>
          <li><Icon name="check" :size="14" :stroke="3" /> 建议节奏：每日 25–40 分钟</li>
          <li><Icon name="check" :size="14" :stroke="3" /> 预计微单元：8 个 MVP 单元</li>
          <li><Icon name="check" :size="14" :stroke="3" /> 验收标准：目标达成度 ≥ 75%</li>
        </ul>
      </div>
    </section>

    <!-- 第三步：分解结果 -->
    <section v-else class="path-panel card fade-up">
      <div class="path-head">
        <div>
          <span class="tag tag--ai"><Icon name="sparkles" :size="12" /> AI 已完成目标分解</span>
          <h2>{{ decomposeResult?.goalName || growth.goal?.name || frontendPath.name }}</h2>
          <p v-if="decomposeResult?.aiSummary">{{ decomposeResult.aiSummary }}</p>
          <p v-else>4 个能力域 · 8 个知识点 · 4 个里程碑</p>
        </div>
        <button class="btn btn--primary" @click="router.push('/')">进入学习驾驶舱 <Icon name="arrowRight" :size="15" /></button>
      </div>

      <div v-if="decomposeResult?.stages?.length" class="ai-stages">
        <article v-for="stage in decomposeResult.stages" :key="stage.name" class="ai-stage">
          <div class="ai-stage__head">
            <strong>{{ stage.name }}</strong>
            <span class="tag tag--neutral">{{ stage.durationWeeks }} 周</span>
          </div>
          <ul>
            <li v-for="topic in stage.topics" :key="topic">{{ topic }}</li>
          </ul>
        </article>
      </div>

      <div class="domain-list">
        <article v-for="(domain, index) in frontendPath.competencies" :key="domain.id" class="domain">
          <div class="domain__number" :style="{ background: domain.color }">{{ index + 1 }}</div>
          <div class="domain__body">
            <div class="domain__title">
              <strong>{{ domain.name }}</strong>
              <span class="tag tag--neutral">权重 {{ domain.weight }}%</span>
            </div>
            <ul class="point-list">
              <li v-for="point in domain.points" :key="point.id"><Icon name="check" :size="12" :stroke="3" /> {{ point.name }}</li>
            </ul>
          </div>
        </article>
      </div>

      <div class="milestones">
        <div v-for="(milestone, index) in frontendPath.milestones" :key="milestone.id" class="milestone">
          <span class="milestone__week">第 {{ milestone.week }} 周</span>
          <strong>{{ milestone.name }}</strong>
          <small>{{ milestone.standard }}</small>
          <i class="milestone__dot" :style="{ animationDelay: index * 0.1 + 's' }"></i>
        </div>
      </div>
    </section>
  </div>
</template>

<style scoped>
.career-page { max-width: 1160px; }

/* 步骤条 */
.steps { display: flex; gap: 10px; margin: 0 0 28px; padding: 0; list-style: none; }
.step {
  flex: 1;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  border-radius: 14px;
  border: 1px solid var(--border);
  background: var(--surface);
  color: var(--text-3);
  font-size: 13.5px;
  font-weight: 600;
  transition: all var(--t);
}
.step__index { display: grid; place-items: center; width: 26px; height: 26px; border-radius: 50%; background: var(--surface-3); color: var(--text-3); font-size: 12px; font-weight: 800; }
.step--active { border-color: var(--primary); color: var(--primary-strong); box-shadow: 0 0 0 4px var(--primary-soft); }
.step--active .step__index { background: linear-gradient(135deg, var(--primary), var(--primary-strong)); color: #fff; }
.step--done { color: var(--text-2); }
.step--done .step__index { background: var(--success); color: #fff; }

/* AI 提示 */
.ai-tip {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px 22px;
  margin-bottom: 18px;
  background: linear-gradient(110deg, #f3f0ff, #fff 60%, #f0fdf8);
}
.ai-tip__icon { display: grid; place-items: center; width: 46px; height: 46px; flex: 0 0 46px; border-radius: 14px; background: linear-gradient(135deg, var(--primary), #a071ff); color: #fff; box-shadow: var(--shadow-primary); }
.ai-tip__body { flex: 1; }
.ai-tip__body strong { display: flex; align-items: center; gap: 8px; font-size: 15px; }
.ai-tip p { margin: 4px 0 0; color: var(--text-2); font-size: 13.5px; }

.interview { padding: 18px; margin-bottom: 18px; }
.interview__messages { max-height: 320px; overflow-y: auto; margin-bottom: 14px; }
.interview__hint { display: flex; align-items: center; gap: 8px; margin: 0; padding: 12px 14px; border-radius: 12px; background: var(--primary-soft); color: var(--primary-strong); font-size: 13.5px; }
.interview__msg { margin: 10px 0; padding: 11px 14px; border-radius: 14px; font-size: 14px; line-height: 1.65; }
.interview__msg--user { margin-left: 22%; background: linear-gradient(135deg, var(--primary), var(--primary-strong)); color: #fff; border-bottom-right-radius: 4px; }
.interview__msg--assistant { margin-right: 12%; background: var(--surface-2); border: 1px solid var(--border); border-bottom-left-radius: 4px; }
.interview__composer { display: flex; gap: 10px; }
.interview__composer .input { flex: 1; }

/* 职业卡片 */
.career-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 16px; }
.career-card { position: relative; display: flex; flex-direction: column; text-align: left; padding: 24px; cursor: pointer; font: inherit; color: inherit; }
.career-card--selected { border-color: var(--primary); box-shadow: 0 0 0 4px var(--primary-soft), var(--shadow); }
.career-card__top { display: flex; justify-content: space-between; align-items: center; margin-bottom: 14px; }
.career-card__icon { display: grid; place-items: center; width: 50px; height: 50px; border-radius: 15px; background: var(--surface-3); font-size: 26px; }
.career-card__category { color: var(--text-3); font-size: 12px; font-weight: 600; letter-spacing: 0.04em; }
.career-card h2 { margin: 4px 0 8px; font-size: 19px; }
.career-card p { flex: 1; margin: 0 0 14px; color: var(--text-2); font-size: 13.5px; line-height: 1.65; }
.career-meta { display: flex; flex-wrap: wrap; gap: 10px; color: var(--text-2); font-size: 12px; }
.career-meta span { display: inline-flex; align-items: center; gap: 4px; }
.career-card__check {
  position: absolute;
  top: 16px;
  right: 16px;
  display: grid;
  place-items: center;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  background: var(--primary);
  color: #fff;
  opacity: 0;
  transform: scale(0.6);
  transition: all var(--t);
}
.career-card--selected .career-card__check { opacity: 1; transform: scale(1); }
.career-card--selected .career-card__top .tag { visibility: hidden; }

.selection-actions { display: flex; justify-content: flex-end; align-items: center; gap: 16px; margin-top: 24px; }
.selection-note { display: inline-flex; align-items: center; gap: 6px; color: var(--warning-strong); font-size: 13px; }

/* 目标条件 */
.goal-layout { display: grid; grid-template-columns: 1.1fr 0.9fr; gap: 20px; }
.goal-form { padding: 30px; display: flex; flex-direction: column; gap: 20px; }
.goal-form h2 { margin: 0; font-size: 21px; }
.goal-form__desc { margin: -12px 0 0; color: var(--text-2); font-size: 13.5px; }
.range-row { display: flex; align-items: center; gap: 16px; }
.range-row input { flex: 1; }
.range-row strong { min-width: 64px; color: var(--primary-strong); font-size: 16px; text-align: right; }
.duration-options { display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; }
.duration {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  padding: 12px;
  border: 1px solid var(--border-strong);
  border-radius: 12px;
  background: var(--surface);
  color: var(--text-2);
  font: inherit;
  cursor: pointer;
  transition: all var(--t-fast);
}
.duration strong { font-size: 15px; color: var(--text); }
.duration small { font-size: 12px; }
.duration:hover { border-color: var(--primary); }
.duration--on { border-color: var(--primary); background: var(--primary-soft); box-shadow: 0 0 0 3px var(--primary-soft); }
.duration--on strong, .duration--on small { color: var(--primary-strong); }
.form-actions { display: flex; justify-content: flex-end; gap: 10px; margin-top: 4px; }

.commitment { padding: 30px; }
.commitment h2 { margin: 16px 0 8px; font-size: 22px; }
.commitment p { margin: 0 0 18px; color: rgba(255, 255, 255, 0.82); font-size: 14.5px; }
.commitment ul { margin: 0; padding: 0; list-style: none; display: flex; flex-direction: column; gap: 10px; }
.commitment li { display: flex; align-items: center; gap: 10px; padding: 10px 14px; border-radius: 12px; background: rgba(255, 255, 255, 0.12); font-size: 14px; }
.commitment li .icon { color: #6ee7b7; }

/* 分解结果 */
.path-panel { padding: 32px; }
.path-head { display: flex; justify-content: space-between; align-items: flex-start; gap: 20px; margin-bottom: 26px; }
.path-head h2 { margin: 12px 0 4px; font-size: 24px; }
.path-head p { margin: 0; color: var(--text-2); }
.ai-stages { display: grid; grid-template-columns: repeat(2, 1fr); gap: 12px; margin-bottom: 20px; }
.ai-stage { padding: 16px 18px; border: 1px solid var(--primary-soft-2); border-radius: 14px; background: linear-gradient(140deg, #f6f3ff, #fff); }
.ai-stage__head { display: flex; justify-content: space-between; align-items: center; }
.ai-stage ul { margin: 10px 0 0; padding-left: 18px; color: var(--text-2); font-size: 13.5px; }
.domain-list { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }
.domain { display: flex; gap: 16px; padding: 18px; border: 1px solid var(--border); border-radius: 16px; background: var(--surface-2); }
.domain__number { display: grid; place-items: center; width: 38px; height: 38px; flex: 0 0 38px; border-radius: 12px; color: #fff; font-family: 'Sora', sans-serif; font-weight: 800; box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.3); }
.domain__body { flex: 1; }
.domain__title { display: flex; justify-content: space-between; align-items: center; }
.domain__title strong { font-size: 15px; }
.point-list { display: flex; flex-direction: column; gap: 6px; margin: 10px 0 0; padding: 0; list-style: none; color: var(--text-2); font-size: 13.5px; }
.point-list li { display: flex; align-items: center; gap: 6px; }
.point-list .icon { color: var(--success); }
.milestones { display: grid; grid-template-columns: repeat(4, 1fr); gap: 12px; margin-top: 24px; position: relative; }
.milestones::before { content: ''; position: absolute; top: 0; left: 4%; right: 4%; height: 2px; background: linear-gradient(90deg, var(--primary-soft-2), var(--primary), var(--primary-soft-2)); }
.milestone { position: relative; display: flex; flex-direction: column; gap: 2px; padding: 22px 16px 16px; border-radius: 14px; background: var(--surface-2); }
.milestone__week { color: var(--primary); font-size: 12px; font-weight: 700; }
.milestone strong { font-size: 14px; }
.milestone small { color: var(--text-3); font-size: 12px; }
.milestone__dot { position: absolute; top: -6px; left: 50%; width: 12px; height: 12px; margin-left: -6px; border-radius: 50%; background: var(--primary); box-shadow: 0 0 0 4px var(--surface); }

@media (max-width: 900px) {
  .career-grid, .goal-layout, .domain-list, .ai-stages { grid-template-columns: 1fr; }
  .milestones { grid-template-columns: 1fr 1fr; }
  .milestones::before { display: none; }
  .steps { overflow-x: auto; }
  .step { min-width: 130px; }
  .path-head { flex-direction: column; }
  .ai-tip { flex-wrap: wrap; }
  .selection-actions { flex-direction: column; align-items: stretch; }
}
</style>
