<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { ecoApi } from '../api'
import { useAuthStore } from '../stores/auth'
import { timeAgo } from '../utils/eco'
import MarkdownRenderer from '../components/MarkdownRenderer.vue'
import SkillLevel from '../components/SkillLevel.vue'
import Icon from '../components/Icon.vue'

const route = useRoute()
const auth = useAuthStore()
const project = ref(null)
const activeId = ref(route.query.task || '')
const form = ref({ content: '', repoUrl: '' })
const submitting = ref(false)
const error = ref('')
const lastResult = ref(null)

const active = computed(() => project.value?.tasks.find((t) => t.id === activeId.value))
const taskSubs = computed(() => (project.value?.submissions || []).filter((s) => s.taskId === activeId.value))
const canWrite = computed(() => auth.hasPermission('growth:write'))

async function load() {
  project.value = await ecoApi.project(route.params.id)
  if (!activeId.value) activeId.value = project.value.tasks.find((t) => t.bestScore < 60)?.id || project.value.tasks[0]?.id
}

function skillName(id) {
  return active.value?.match?.items.find((it) => it.skillId === id)?.skillName || id
}

function pick(id) {
  activeId.value = id
  lastResult.value = null
  error.value = ''
}

async function submit() {
  error.value = ''
  submitting.value = true
  try {
    lastResult.value = await ecoApi.submit(project.value.id, activeId.value, form.value)
    form.value = { content: '', repoUrl: '' }
    await load()
  } catch (e) {
    error.value = e.message
  } finally {
    submitting.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="page pd">
    <router-link to="/projects" class="back-link"><Icon name="arrowLeft" :size="14" /> 项目实践</router-link>
    <template v-if="project">
      <header class="card panel pd__head" :style="{ borderTop: `4px solid ${project.color}` }">
        <div>
          <span class="tag tag--neutral">{{ project.domain }} · 难度 {{ '★'.repeat(project.difficulty) }}</span>
          <h1>{{ project.title }}</h1>
          <p>{{ project.summary }}</p>
        </div>
        <div class="pd__progress">
          <strong class="num">{{ project.progress }}%</strong>
          <small>完成 {{ project.completed }}/{{ project.tasks.length }} 个任务</small>
        </div>
      </header>

      <div class="pd__body">
        <aside class="stack">
          <section class="card panel">
            <h3 class="side-title">任务链</h3>
            <button
              v-for="(t, i) in project.tasks"
              :key="t.id"
              class="task"
              :class="{ 'task--active': t.id === activeId, 'task--done': t.bestScore >= 60 }"
              @click="pick(t.id)"
            >
              <span class="task__idx">
                <Icon v-if="t.bestScore >= 60" name="check" :size="12" :stroke="3" />
                <template v-else>{{ i + 1 }}</template>
              </span>
              <span class="task__body">
                <strong>{{ t.title }}</strong>
                <small>{{ t.submitted ? `最高 ${t.bestScore} 分` : `约 ${t.estimatedHours} 小时` }} · 匹配 {{ t.match?.score ?? 0 }}%</small>
              </span>
            </button>
          </section>
          <section class="card panel">
            <MarkdownRenderer :source="project.description" />
          </section>
        </aside>

        <div v-if="active" class="stack">
          <section class="card panel">
            <div class="panel-head">
              <div><h2>{{ active.title }}</h2><p>交付物：{{ active.deliverable }}</p></div>
              <span class="tag" :class="active.match?.qualified ? 'tag--success' : 'tag--warning'">{{ active.match?.qualified ? '技能已达标' : '有挑战的任务' }}</span>
            </div>
            <p class="pd__desc">{{ active.description }}</p>
            <h4 class="side-title">技能要求</h4>
            <div v-for="it in active.match?.items || []" :key="it.skillId" class="req">
              <span>{{ it.skillName }}</span>
              <SkillLevel :level="it.current" :required="it.required" />
            </div>
          </section>

          <section class="card panel">
            <div class="panel-head"><div><h2>提交成果</h2><p>描述方案、关键代码与验证结果，越完整评分越高</p></div></div>
            <form v-if="canWrite" class="stack" @submit.prevent="submit">
              <textarea v-model="form.content" class="textarea" rows="10" placeholder="支持 Markdown：## 方案设计 / 关键代码 ``` / 测试与性能验证 / 取舍与反思"></textarea>
              <input v-model="form.repoUrl" class="input" placeholder="代码仓库链接（可选）https://github.com/..." />
              <p v-if="error" class="error-text">{{ error }}</p>
              <div class="row row--between">
                <span class="muted small">提交后 Review Agent 将从方案完整性、技术正确性、工程质量与表达四个维度评审</span>
                <button class="btn btn--primary" :disabled="submitting || form.content.trim().length < 20">
                  <Icon name="send" :size="14" /> {{ submitting ? 'Review Agent 评审中…' : '提交评审' }}
                </button>
              </div>
            </form>
            <p v-else class="muted small">登录正式账号后可提交项目任务。</p>
          </section>

          <section v-if="lastResult" class="card panel result">
            <div class="result__score">
              <strong class="num">{{ lastResult.score }}</strong>
              <small>{{ lastResult.evaluator === 'ai' ? 'Review Agent 评分' : '规则评审评分' }}</small>
            </div>
            <div class="result__body">
              <MarkdownRenderer :source="lastResult.feedback" />
              <div class="row row--wrap">
                <span v-for="e in lastResult.evidence" :key="e.id" class="chip chip--ok"><Icon name="trophy" :size="11" /> 新增能力证据 · {{ skillName(e.skillId) }}</span>
              </div>
            </div>
          </section>

          <section v-if="taskSubs.length" class="card panel">
            <h3 class="side-title">历史提交</h3>
            <details v-for="s in taskSubs" :key="s.id" class="sub">
              <summary>
                <b class="num">{{ s.score }} 分</b>
                <span>{{ s.evaluator === 'ai' ? 'AI 评审' : '规则评审' }}</span>
                <small>{{ timeAgo(s.createdAt) }}</small>
              </summary>
              <MarkdownRenderer :source="s.feedback" />
            </details>
          </section>
        </div>
      </div>
    </template>
  </div>
</template>

<style scoped>
.pd { display: flex; flex-direction: column; gap: 18px; }
.pd .back-link { margin-bottom: 0; }
.pd__head { display: flex; justify-content: space-between; gap: 24px; }
.pd__head h1 { margin: 10px 0 6px; font-size: 28px; }
.pd__head p { margin: 0; color: var(--text-2); }
.pd__progress { display: flex; flex-direction: column; align-items: flex-end; justify-content: center; }
.pd__progress strong { font-size: 36px; color: var(--primary-strong); }
.pd__progress small { color: var(--text-3); }
.pd__body { display: grid; grid-template-columns: 320px minmax(0, 1fr); gap: 18px; align-items: start; }
.side-title { margin: 0 0 10px; font-size: 13.5px; }
.task { display: flex; align-items: center; gap: 10px; width: 100%; margin-bottom: 6px; padding: 10px 12px; border: 1px solid transparent; border-radius: 12px; background: transparent; font: inherit; text-align: left; cursor: pointer; }
.task:hover { background: var(--surface-2); }
.task--active { border-color: var(--primary-soft-2); background: var(--primary-soft); }
.task__idx { display: grid; place-items: center; width: 26px; height: 26px; flex: 0 0 26px; border-radius: 50%; background: var(--surface-3); color: var(--text-2); font-size: 12px; font-weight: 700; }
.task--done .task__idx { background: var(--success); color: #fff; }
.task__body { display: flex; flex-direction: column; line-height: 1.4; }
.task__body strong { font-size: 13.5px; }
.task__body small { color: var(--text-3); font-size: 11.5px; }
.pd__desc { margin: 0 0 16px; color: var(--text-2); }
.req { display: flex; align-items: center; justify-content: space-between; padding: 8px 0; border-bottom: 1px dashed var(--border); font-size: 13.5px; }
.result { display: flex; gap: 22px; border-color: var(--primary-soft-2); }
.result__score { display: flex; flex-direction: column; align-items: center; justify-content: center; align-self: flex-start; min-width: 110px; padding: 16px 12px; border-radius: 16px; background: var(--primary-soft); }
.result__score strong { font-size: 40px; color: var(--primary-strong); line-height: 1.1; }
.result__score small { color: var(--text-3); font-size: 12px; }
.result__body { flex: 1; min-width: 0; }
.result__body :deep(.markdown-body) { font-size: 14px; }
.sub { padding: 10px 0; border-bottom: 1px dashed var(--border); }
.sub summary { display: flex; align-items: center; gap: 10px; cursor: pointer; font-size: 13.5px; }
.sub summary small { margin-left: auto; color: var(--text-3); }
.sub :deep(.markdown-body) { margin-top: 8px; font-size: 13.5px; }
@media (max-width: 960px) {
  .pd__body { grid-template-columns: 1fr; }
  .pd__head { flex-direction: column; }
}
</style>
