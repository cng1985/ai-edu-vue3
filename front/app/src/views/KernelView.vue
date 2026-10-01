<script setup>
import { computed, onBeforeUnmount, onMounted, ref } from 'vue'
import { kernelApi } from '../api'
import { useAuthStore } from '../stores/auth'
import { useGrowthStore } from '../stores/growth'
import { timeAgo } from '../utils/eco'
import PipelineTrace from '../components/PipelineTrace.vue'
import MarkdownRenderer from '../components/MarkdownRenderer.vue'
import SkillLevel from '../components/SkillLevel.vue'
import Icon from '../components/Icon.vue'

const auth = useAuthStore()
const growth = useGrowthStore()
const stages = ref([])
const traces = ref([])
const response = ref('')
const output = ref(null)
const running = ref(false)
const error = ref('')
const question = ref('')
const runs = ref([])
let cancel = null

const canRun = computed(() => auth.hasPermission('growth:write'))
const layers = ['用户', 'AI Learning Kernel', 'AI Execution Kernel', 'Pipeline Runtime', 'Resource / Tool Layer']
const resources = ['Knowledge', 'Course', 'Project', 'Task', 'Skill', 'Talent']

function run() {
  if (running.value) return
  traces.value = []
  response.value = ''
  output.value = null
  error.value = ''
  running.value = true
  cancel = kernelApi.runStream(question.value, {
    stage: ({ stage }) => traces.value.push(stage),
    token: ({ content }) => (response.value += content),
    done: async ({ run }) => {
      running.value = false
      output.value = run.output
      response.value = run.output?.response || response.value
      runs.value = await kernelApi.runs()
      growth.refresh()
    },
    error: (e) => {
      running.value = false
      error.value = e.message
    }
  })
}

function showRun(r) {
  traces.value = r.trace || []
  output.value = r.output
  response.value = r.output?.response || ''
}

onMounted(async () => {
  stages.value = await kernelApi.stages()
  runs.value = await kernelApi.runs().catch(() => [])
  if (runs.value[0]) showRun(runs.value[0])
})
onBeforeUnmount(() => cancel?.())
</script>

<template>
  <div class="page kernel">
    <header class="page-header">
      <span class="eyebrow">AI Learning Kernel</span>
      <h1>AI 学习内核</h1>
      <p>一次完整的 AI 学习流程由 8 个 Stage 组成，Pipeline Runtime 按顺序执行并记录每一步的轨迹。结构化分析（差距、推荐、计划）由成长引擎计算，AI 只负责把结果讲给你听。</p>
    </header>

    <div class="kernel__top">
      <section class="card panel arch">
        <h3 class="side-title">AI 底层技术架构</h3>
        <div v-for="(l, i) in layers" :key="l" class="arch__layer" :class="{ 'arch__layer--hl': i === 3 }">{{ l }}</div>
        <div class="arch__res">
          <span v-for="r in resources" :key="r">{{ r }}</span>
        </div>
      </section>

      <section class="card panel">
        <div class="panel-head">
          <div><h2>运行学习 Pipeline</h2><p>可以附带一个问题，Learning Agent 会结合分析结果回答</p></div>
        </div>
        <textarea v-model="question" class="textarea" rows="3" placeholder="例如：我这周只有 6 小时，应该优先学什么？（可选）" :disabled="!canRun"></textarea>
        <div class="row row--between run-bar">
          <span class="muted small">{{ canRun ? '每次运行都会基于最新的知识状态（含遗忘衰减）重算技能' : '游客无法运行，请注册正式账号' }}</span>
          <button class="btn btn--primary" :disabled="running || !canRun" @click="run">
            <Icon name="play" :size="14" /> {{ running ? '运行中…' : '运行 AI 学习内核' }}
          </button>
        </div>
        <p v-if="error" class="error-text">{{ error }}</p>
      </section>
    </div>

    <div class="kernel__body">
      <section class="card panel">
        <h3 class="side-title">Pipeline 执行轨迹</h3>
        <PipelineTrace :stages="stages" :traces="traces" :running="running" />
      </section>

      <div class="stack">
        <section v-if="response" class="card panel ai-out">
          <div class="panel-head">
            <div><h2><span class="tag tag--ai">AI</span> 学习建议</h2><p>{{ output?.responseSource === 'rule' ? '规则模板生成（大模型未配置）' : 'Learning Agent 生成' }}</p></div>
          </div>
          <MarkdownRenderer :source="response" :live="running" />
        </section>

        <section v-if="output?.gap" class="card panel">
          <div class="panel-head"><div><h3>SkillGapStage · 技能差距</h3><p>{{ output.gap.role?.name }} 达成度 {{ output.gap.readiness }}%</p></div></div>
          <div v-for="g in output.gap.gaps" :key="g.skillId" class="mini">
            <span>{{ g.skillName }}<small>{{ g.capabilityName }}</small></span>
            <SkillLevel :level="g.current" :required="g.required" />
          </div>
        </section>

        <section v-if="output?.plan" class="card panel">
          <div class="panel-head"><div><h3>LearningPlanStage · 学习计划</h3><p>{{ output.plan.summary }}</p></div></div>
          <div v-for="w in output.plan.weeks" :key="w.week" class="week">
            <strong>第 {{ w.week }} 周 <small v-if="w.focus" class="muted">· 重点 {{ w.focus }}</small></strong>
            <router-link
              v-for="it in w.items"
              :key="it.refId + (it.subId || '')"
              :to="it.type === 'project' ? `/projects/${it.refId}?task=${it.subId}` : `/knowledge/${it.refId}`"
              class="mini"
            >
              <span><Icon :name="it.type === 'project' ? 'code' : 'book'" :size="12" /> {{ it.title }}<small>{{ it.reason }}</small></span>
              <b class="num">{{ it.minutes }}′</b>
            </router-link>
          </div>
        </section>

        <section v-if="output?.evaluation" class="card panel">
          <div class="panel-head"><div><h3>EvaluationStage · 规划评估</h3></div></div>
          <div class="row">
            <span class="tag tag--info">差距覆盖率 {{ output.evaluation.coverage }}%</span>
            <span class="tag tag--success">计划可行性 {{ output.evaluation.feasibility }}%</span>
          </div>
          <ul v-if="output.evaluation.notes.length" class="small muted notes">
            <li v-for="n in output.evaluation.notes" :key="n">{{ n }}</li>
          </ul>
        </section>

        <section v-if="runs.length" class="card panel">
          <h3 class="side-title">运行记录</h3>
          <button v-for="r in runs.slice(0, 8)" :key="r.id" class="mini mini--btn" @click="showRun(r)">
            <span>{{ r.status === 'ok' ? '运行成功' : '运行失败' }}<small>{{ timeAgo(r.createdAt) }} · {{ r.trace?.length }} 个 Stage</small></span>
            <b class="num">{{ r.durationMs }}ms</b>
          </button>
        </section>
      </div>
    </div>
  </div>
</template>

<style scoped>
.kernel { display: flex; flex-direction: column; gap: 18px; }
.kernel .page-header { margin-bottom: 0; }
.kernel__top { display: grid; grid-template-columns: 300px minmax(0, 1fr); gap: 18px; }
.kernel__body { display: grid; grid-template-columns: 360px minmax(0, 1fr); gap: 18px; align-items: start; }
.side-title { margin: 0 0 12px; font-size: 14px; }
.arch__layer { margin-bottom: 6px; padding: 7px 10px; border-radius: 10px; background: var(--surface-2); border: 1px solid var(--border); font-size: 12.5px; font-weight: 600; text-align: center; }
.arch__layer--hl { border-color: var(--primary); background: var(--primary-soft); color: var(--primary-strong); }
.arch__res { display: grid; grid-template-columns: repeat(3, 1fr); gap: 4px; }
.arch__res span { padding: 4px; border-radius: 8px; background: var(--ink); color: #c9cbe6; font-size: 11px; text-align: center; }
.run-bar { margin-top: 12px; }
.ai-out { border-color: var(--primary-soft-2); background: linear-gradient(160deg, #f6f3ff, #fff 55%); }
.mini { display: flex; align-items: center; justify-content: space-between; gap: 10px; padding: 8px 0; border-bottom: 1px dashed var(--border); color: inherit; font-size: 13.5px; }
.mini > span { display: flex; flex-direction: column; min-width: 0; }
.mini small { color: var(--text-3); font-size: 11.5px; }
.mini--btn { width: 100%; border-left: none; border-right: none; border-top: none; background: none; font: inherit; text-align: left; cursor: pointer; }
.mini--btn:hover > span { color: var(--primary); }
.week { margin-bottom: 14px; }
.week > strong { display: block; margin-bottom: 4px; font-size: 13.5px; }
.notes { margin: 10px 0 0; padding-left: 18px; }
@media (max-width: 960px) {
  .kernel__top, .kernel__body { grid-template-columns: 1fr; }
}
</style>
