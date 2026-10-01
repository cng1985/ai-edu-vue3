<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { ecoApi, kernelApi, meApi } from '../api'
import { useAuthStore } from '../stores/auth'
import { KNOWLEDGE_STATUS, RELATION_TYPES, RESOURCE_TYPES, dateText, pct } from '../utils/eco'
import MarkdownRenderer from '../components/MarkdownRenderer.vue'
import MasteryBar from '../components/MasteryBar.vue'
import AgentChat from '../components/AgentChat.vue'
import Icon from '../components/Icon.vue'

const route = useRoute()
const auth = useAuthStore()
const detail = ref(null)
const error = ref('')
const answers = ref([])
const result = ref(null)
const busy = ref(false)
const toast = ref('')
const knowledgeAgent = ref(null)
const showAgent = ref(false)

const state = computed(() => detail.value?.state)
const canWrite = computed(() => auth.hasPermission('growth:write'))
const prereqs = computed(() => (detail.value?.relations || []).filter((r) => r.direction === 'in' && (r.type === 'prerequisite' || r.type === 'advanced')))
const nexts = computed(() => (detail.value?.relations || []).filter((r) => r.direction === 'out' && (r.type === 'prerequisite' || r.type === 'advanced')))
const others = computed(() => (detail.value?.relations || []).filter((r) => r.type === 'related' || r.type === 'similar'))

async function load() {
  error.value = ''
  result.value = null
  try {
    detail.value = await ecoApi.knowledge(route.params.id)
    answers.value = detail.value.questions.map(() => null)
  } catch (e) {
    error.value = e.message
  }
}

function flash(text) {
  toast.value = text
  setTimeout(() => (toast.value = ''), 2200)
}

async function study() {
  busy.value = true
  try {
    detail.value.state = await meApi.recordEvent({ knowledgeId: detail.value.knowledge.id, type: 'study', source: 'knowledge-page' })
    flash('已记录本次学习')
  } catch (e) {
    flash(e.message)
  } finally {
    busy.value = false
  }
}

async function submitPractice() {
  busy.value = true
  try {
    result.value = await meApi.practice(detail.value.knowledge.id, answers.value.map((a) => (a == null ? -1 : a)))
    detail.value.state = result.value.state
  } catch (e) {
    flash(e.message)
  } finally {
    busy.value = false
  }
}

function retry() {
  result.value = null
  answers.value = detail.value.questions.map(() => null)
}

watch(() => route.params.id, load)
onMounted(async () => {
  await load()
  try {
    const agents = await kernelApi.agents()
    knowledgeAgent.value = agents.find((a) => a.code === 'knowledge')
  } catch { /* 无 AI 权限时不展示 */ }
})
</script>

<template>
  <div class="page kd">
    <router-link to="/knowledge" class="back-link"><Icon name="arrowLeft" :size="14" /> 知识图谱</router-link>
    <div v-if="error" class="card empty-state"><h3>{{ error }}</h3></div>

    <template v-else-if="detail">
      <header class="kd__head card panel">
        <div class="kd__title">
          <span class="tag" :class="`tag--${KNOWLEDGE_STATUS[state.status]?.tone}`">{{ KNOWLEDGE_STATUS[state.status]?.label }}</span>
          <span v-if="state.due" class="tag tag--warning">今日待复习</span>
          <h1>{{ detail.knowledge.name }}</h1>
          <p>{{ detail.knowledge.summary }}</p>
          <div class="row row--wrap">
            <span class="chip"><Icon name="layers" :size="12" /> {{ detail.knowledge.domain }}</span>
            <span class="chip"><Icon name="bolt" :size="12" /> 难度 {{ detail.knowledge.difficulty }}</span>
            <span class="chip"><Icon name="clock" :size="12" /> {{ detail.knowledge.estimatedMinutes }} 分钟</span>
            <router-link v-for="s in detail.skills" :key="s.id" :to="`/knowledge?skill=${s.id}`" class="chip">技能 · {{ s.name }}</router-link>
          </div>
        </div>
        <div class="kd__state">
          <MasteryBar :value="state.effectiveMastery" label="当前掌握度" />
          <MasteryBar :value="state.confidence" label="置信度" />
          <dl>
            <div><dt>学习</dt><dd class="num">{{ state.studyCount }}</dd></div>
            <div><dt>练习题</dt><dd class="num">{{ state.practiceCount }}</dd></div>
            <div><dt>正确率</dt><dd class="num">{{ pct(state.accuracy) }}%</dd></div>
            <div><dt>下次复习</dt><dd>{{ state.nextReviewAt ? dateText(state.nextReviewAt) : '—' }}</dd></div>
          </dl>
          <p class="small muted">学会不是 Boolean：掌握度由学习投入与练习表现计算，并随时间衰减。</p>
        </div>
      </header>

      <div class="kd__body">
        <div class="stack">
          <section class="card panel">
            <MarkdownRenderer :source="detail.knowledge.content || '暂无正文'" />
            <div v-if="canWrite" class="kd__study">
              <span class="muted small">读完正文后记录学习，再做练习巩固</span>
              <button class="btn btn--soft" :disabled="busy" @click="study"><Icon name="check" :size="14" /> 我学完了</button>
            </div>
          </section>

          <section v-if="detail.questions.length" class="card panel">
            <div class="panel-head"><div><h2>知识练习</h2><p>练习正确率是掌握度的主要依据</p></div></div>
            <div v-for="(q, i) in detail.questions" :key="i" class="quiz-q">
              <strong>{{ i + 1 }}. {{ q.text }}</strong>
              <label
                v-for="(opt, j) in q.options"
                :key="j"
                class="quiz-opt"
                :class="{
                  'quiz-opt--picked': answers[i] === j,
                  'quiz-opt--right': result && result.results[i] && answers[i] === j,
                  'quiz-opt--wrong': result && !result.results[i] && answers[i] === j
                }"
              >
                <input v-model="answers[i]" type="radio" :name="`q${i}`" :value="j" :disabled="!!result" />
                {{ opt }}
              </label>
            </div>
            <div v-if="result" class="quiz-result">
              <strong>答对 {{ result.correct }}/{{ result.total }}</strong>
              <span>掌握度更新为 {{ pct(result.state.mastery) }}%，{{ result.state.status === 'mastered' ? '已掌握 🎉' : '继续加油' }}</span>
              <button class="btn btn--sm btn--ghost" @click="retry">再练一次</button>
            </div>
            <button v-else-if="canWrite" class="btn btn--primary" :disabled="busy || answers.some((a) => a == null)" @click="submitPractice">提交练习</button>
            <p v-else class="muted small">登录正式账号后可记录练习成绩。</p>
          </section>

          <section v-if="knowledgeAgent" class="card panel">
            <div class="panel-head">
              <div><h2>问 Knowledge Agent</h2><p>结合知识图谱与课程知识库解答疑问</p></div>
              <button class="btn btn--sm btn--ghost" @click="showAgent = !showAgent">{{ showAgent ? '收起' : '展开对话' }}</button>
            </div>
            <AgentChat v-if="showAgent" :agent="knowledgeAgent" compact :initial-message="`请帮我讲解「${detail.knowledge.name}」，并说明它的前置知识`" />
          </section>
        </div>

        <aside class="stack">
          <section class="card panel">
            <h3 class="side-title">前置知识</h3>
            <p v-if="!prereqs.length" class="muted small">无前置要求，可直接学习</p>
            <router-link v-for="r in prereqs" :key="r.knowledge.id" :to="`/knowledge/${r.knowledge.id}`" class="mini">
              <span>{{ r.knowledge.name }}<small>{{ RELATION_TYPES[r.type] }}</small></span>
              <MasteryBar :value="r.mastery" compact />
            </router-link>
            <h3 class="side-title">后续知识</h3>
            <p v-if="!nexts.length" class="muted small">暂无</p>
            <router-link v-for="r in nexts" :key="r.knowledge.id" :to="`/knowledge/${r.knowledge.id}`" class="mini">
              <span>{{ r.knowledge.name }}</span>
              <MasteryBar :value="r.mastery" compact />
            </router-link>
            <template v-if="others.length">
              <h3 class="side-title">关联 / 相似</h3>
              <router-link v-for="r in others" :key="r.knowledge.id" :to="`/knowledge/${r.knowledge.id}`" class="mini">
                <span>{{ r.knowledge.name }}<small>{{ RELATION_TYPES[r.type] }}</small></span>
              </router-link>
            </template>
          </section>

          <section class="card panel">
            <h3 class="side-title">引用该知识点的课程</h3>
            <p v-if="!detail.chapters.length" class="muted small">课程结构与知识结构分离，暂无课程引用</p>
            <router-link v-for="c in detail.chapters" :key="c.courseId + c.chapterId" :to="`/courses/${c.courseId}/${c.chapterId}`" class="mini">
              <span>{{ c.chapterTitle }}<small>{{ c.courseTitle }}</small></span><Icon name="arrowRight" :size="13" />
            </router-link>
            <h3 class="side-title">知识资产</h3>
            <p v-if="!detail.resources.length" class="muted small">暂无相关资源</p>
            <router-link v-for="r in detail.resources" :key="r.id" :to="`/resources?id=${r.id}`" class="mini">
              <span>{{ r.title }}<small>{{ RESOURCE_TYPES[r.type]?.label }}</small></span><Icon :name="RESOURCE_TYPES[r.type]?.icon || 'note'" :size="13" />
            </router-link>
          </section>
        </aside>
      </div>
    </template>
    <div v-if="toast" class="toast">{{ toast }}</div>
  </div>
</template>

<style scoped>
.kd { display: flex; flex-direction: column; gap: 18px; }
.kd .back-link { margin-bottom: 0; }
.kd__head { display: grid; grid-template-columns: 1fr 300px; gap: 28px; }
.kd__title h1 { margin: 10px 0 6px; font-size: 28px; }
.kd__title p { margin: 0 0 14px; color: var(--text-2); }
.kd__state { display: flex; flex-direction: column; gap: 12px; padding: 16px; border-radius: 14px; background: var(--surface-2); }
.kd__state dl { display: grid; grid-template-columns: repeat(4, 1fr); gap: 6px; margin: 0; }
.kd__state dt { color: var(--text-3); font-size: 11.5px; }
.kd__state dd { margin: 0; font-weight: 700; font-size: 14px; }
.kd__state p { margin: 0; }
.kd__body { display: grid; grid-template-columns: minmax(0, 1fr) 300px; gap: 18px; align-items: start; }
.kd__study { display: flex; align-items: center; justify-content: space-between; gap: 12px; margin-top: 22px; padding-top: 16px; border-top: 1px dashed var(--border); }
.quiz-q { display: flex; flex-direction: column; gap: 8px; margin-bottom: 18px; }
.quiz-q strong { font-size: 14.5px; }
.quiz-opt { display: flex; align-items: center; gap: 10px; padding: 10px 14px; border: 1px solid var(--border); border-radius: 12px; cursor: pointer; font-size: 14px; transition: all var(--t-fast); }
.quiz-opt:hover { border-color: var(--primary-soft-2); background: var(--primary-soft); }
.quiz-opt--picked { border-color: var(--primary); background: var(--primary-soft); }
.quiz-opt--right { border-color: var(--success); background: var(--success-soft); }
.quiz-opt--wrong { border-color: var(--danger); background: var(--danger-soft); }
.quiz-result { display: flex; flex-wrap: wrap; align-items: center; gap: 12px; padding: 14px; border-radius: 12px; background: var(--primary-soft); color: var(--primary-deep); font-size: 14px; }
.side-title { margin: 14px 0 8px; font-size: 13.5px; }
.side-title:first-child { margin-top: 0; }
.mini { display: flex; align-items: center; justify-content: space-between; gap: 10px; padding: 8px 0; border-bottom: 1px dashed var(--border); color: inherit; font-size: 13.5px; }
.mini span { display: flex; flex-direction: column; }
.mini small { color: var(--text-3); font-size: 11.5px; }
.mini:hover > span { color: var(--primary); }
@media (max-width: 960px) {
  .kd__head, .kd__body { grid-template-columns: 1fr; }
}
</style>
