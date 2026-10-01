<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ecoApi } from '../api'
import { KNOWLEDGE_STATUS, RELATION_TYPES, pct } from '../utils/eco'
import KnowledgeGraph from '../components/KnowledgeGraph.vue'
import MasteryBar from '../components/MasteryBar.vue'
import Icon from '../components/Icon.vue'

const route = useRoute()
const router = useRouter()
const skills = ref([])
const graph = ref({ nodes: [], edges: [] })
const skillId = ref(route.query.skill || '')
const selected = ref(null)
const view = ref('graph')
const loading = ref(false)

const grouped = computed(() => {
  const map = {}
  for (const s of skills.value) (map[s.category] ||= []).push(s)
  return map
})

const stats = computed(() => {
  const s = { mastered: 0, learning: 0, new: 0 }
  for (const n of graph.value.nodes) s[n.status] = (s[n.status] || 0) + 1
  return s
})

const sortedNodes = computed(() => [...graph.value.nodes].sort((a, b) => a.depth - b.depth || b.mastery - a.mastery))

async function loadGraph() {
  loading.value = true
  try {
    graph.value = await ecoApi.graph(skillId.value)
    selected.value = null
  } finally {
    loading.value = false
  }
}

function pickSkill(id) {
  skillId.value = skillId.value === id ? '' : id
  router.replace({ query: skillId.value ? { skill: skillId.value } : {} })
}

watch(skillId, loadGraph)
onMounted(async () => {
  skills.value = await ecoApi.skills()
  await loadGraph()
})
</script>

<template>
  <div class="page kgv">
    <header class="page-header page-header--split">
      <div>
        <span class="eyebrow">Knowledge Graph</span>
        <h1>知识图谱</h1>
        <p>知识不是目录，而是图谱：前置、关联、相似与高级关系决定了学习顺序。节点颜色代表你的当前掌握度（含遗忘衰减）。</p>
      </div>
      <div class="segmented">
        <button :class="{ active: view === 'graph' }" @click="view = 'graph'">图谱</button>
        <button :class="{ active: view === 'list' }" @click="view = 'list'">列表</button>
      </div>
    </header>

    <section class="card panel filters">
      <div v-for="(list, cat) in grouped" :key="cat" class="filters__group">
        <span class="filters__label">{{ cat }}</span>
        <button v-for="s in list" :key="s.id" class="chip" :class="{ active: s.id === skillId }" @click="pickSkill(s.id)">{{ s.name }}</button>
      </div>
    </section>

    <div class="kgv__stats">
      <span class="tag tag--success">已掌握 {{ stats.mastered || 0 }}</span>
      <span class="tag tag--info">学习中 {{ stats.learning || 0 }}</span>
      <span class="tag tag--neutral">未学习 {{ stats.new || 0 }}</span>
      <span class="legend"><i class="legend__solid"></i>前置 / 高级</span>
      <span class="legend"><i class="legend__dash"></i>关联 / 相似</span>
    </div>

    <div class="kgv__body">
      <section class="card panel kgv__canvas">
        <div v-if="loading" class="muted">加载中…</div>
        <KnowledgeGraph v-else-if="view === 'graph'" :nodes="graph.nodes" :edges="graph.edges" @select="selected = $event" />
        <div v-else class="stack stack--tight">
          <router-link v-for="n in sortedNodes" :key="n.id" :to="`/knowledge/${n.id}`" class="list-row">
            <span class="tag" :class="`tag--${KNOWLEDGE_STATUS[n.status]?.tone}`">{{ KNOWLEDGE_STATUS[n.status]?.label }}</span>
            <div class="list-row__body"><strong>{{ n.name }}</strong><small>{{ n.domain }} · 难度 {{ n.difficulty }} · 约 {{ n.estimatedMinutes }} 分钟</small></div>
            <MasteryBar :value="n.mastery" compact />
          </router-link>
        </div>
      </section>

      <aside class="card panel kgv__side">
        <template v-if="selected">
          <span class="tag" :class="`tag--${KNOWLEDGE_STATUS[selected.status]?.tone}`">{{ KNOWLEDGE_STATUS[selected.status]?.label }}</span>
          <h3>{{ selected.name }}</h3>
          <p class="muted small">{{ selected.domain }} · 难度 {{ selected.difficulty }} · 约 {{ selected.estimatedMinutes }} 分钟</p>
          <p>{{ selected.summary }}</p>
          <MasteryBar :value="selected.mastery" />
          <div class="kgv__rels">
            <template v-for="e in graph.edges" :key="e.id">
              <span v-if="e.toId === selected.id" class="chip">{{ RELATION_TYPES[e.type] }} ← {{ graph.nodes.find((n) => n.id === e.fromId)?.name }}</span>
              <span v-else-if="e.fromId === selected.id" class="chip">{{ RELATION_TYPES[e.type] === '前置知识' ? '后续' : RELATION_TYPES[e.type] }} → {{ graph.nodes.find((n) => n.id === e.toId)?.name }}</span>
            </template>
          </div>
          <router-link :to="`/knowledge/${selected.id}`" class="btn btn--primary btn--block"><Icon name="play" :size="14" /> 开始学习与练习</router-link>
        </template>
        <div v-else class="muted small">
          <Icon name="brain" :size="28" />
          <p>点击图谱中的知识点查看详情；悬停可高亮它的前后置关系。</p>
          <p>当前展示 {{ graph.nodes.length }} 个知识点、{{ graph.edges.length }} 条关系，平均掌握度 {{ pct(graph.nodes.reduce((s, n) => s + n.mastery, 0) / Math.max(1, graph.nodes.length)) }}%。</p>
        </div>
      </aside>
    </div>
  </div>
</template>

<style scoped>
.kgv { display: flex; flex-direction: column; gap: 16px; }
.kgv .page-header { margin-bottom: 4px; }
.filters { display: flex; flex-direction: column; gap: 10px; padding: 18px 22px; }
.filters__group { display: flex; flex-wrap: wrap; align-items: center; gap: 6px; }
.filters__label { width: 72px; color: var(--text-3); font-size: 12px; font-weight: 700; }
.kgv__stats { display: flex; flex-wrap: wrap; align-items: center; gap: 8px; }
.legend { display: inline-flex; align-items: center; gap: 6px; margin-left: 8px; color: var(--text-3); font-size: 12px; }
.legend i { width: 22px; height: 0; border-top: 2px solid #b6bad0; }
.legend__dash { border-top-style: dashed !important; }
.kgv__body { display: grid; grid-template-columns: minmax(0, 1fr) 300px; gap: 16px; align-items: start; }
.kgv__canvas { min-height: 420px; overflow: hidden; }
.kgv__side { position: sticky; top: 24px; display: flex; flex-direction: column; gap: 10px; }
.kgv__side h3 { margin: 0; font-size: 18px; }
.kgv__side p { margin: 0; font-size: 13.5px; }
.kgv__rels { display: flex; flex-wrap: wrap; gap: 6px; }
.stack--tight { gap: 8px; }
@media (max-width: 960px) {
  .kgv__body { grid-template-columns: 1fr; }
  .kgv__side { position: static; }
}
</style>
