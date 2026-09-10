<template>
  <div class="page">
    <PageHeader eyebrow="知识库" title="知识库管理" subtitle="管理课程内容的向量索引，为 AI 对话提供语义检索增强。">
      <el-button
        v-if="auth.hasPermission(PERM.KNOWLEDGE_MANAGE)"
        type="primary"
        :loading="reindexing"
        :icon="Refresh"
        @click="handleReindex"
      >
        重建索引
      </el-button>
    </PageHeader>

    <div class="stat-grid stat-grid--4">
      <StatCard label="文本块数量" :value="status.chunkCount" icon="Files" tone="primary" />
      <StatCard label="覆盖课程" :value="status.courseCount" icon="Reading" tone="success" />
      <StatCard label="覆盖章节" :value="status.chapterCount" icon="Notebook" tone="sky" />
      <StatCard label="索引状态" :value="statusLabel" :icon="statusIcon" :tone="statusTone" :hint="`最后索引 ${formatTime(status.lastIndexedAt)}`" />
    </div>

    <div class="grid-2">
      <div v-loading="loading" class="panel">
        <div class="panel__head">
          <h3 class="panel__title">
            <span class="panel__title-icon"><el-icon><Setting /></el-icon></span>
            索引配置
          </h3>
        </div>
        <div class="panel__body">
          <div class="kv">
            <div class="kv__item">
              <span class="kv__label">嵌入模型</span>
              <span class="kv__value mono">{{ status.embedModel || '—' }}</span>
            </div>
            <div class="kv__item">
              <span class="kv__label">嵌入来源</span>
              <span class="kv__value">
                <el-tag size="small" :type="status.embedSource === 'api' ? 'success' : 'info'">
                  {{ status.embedSource === 'api' ? 'API 嵌入' : '本地哈希嵌入' }}
                </el-tag>
              </span>
            </div>
            <div class="kv__item">
              <span class="kv__label">向量维度</span>
              <span class="kv__value num">{{ status.dimensions || '—' }}</span>
            </div>
            <div class="kv__item">
              <span class="kv__label">最后索引时间</span>
              <span class="kv__value">{{ formatTime(status.lastIndexedAt) }}</span>
            </div>
          </div>
          <div class="hint hint--info" style="margin-top: 16px">
            <el-icon :size="16" style="margin-top: 2px"><InfoFilled /></el-icon>
            <span>知识库使用 SQLite 嵌入式向量存储，结合向量相似度与关键词混合检索。配置 <code class="mono">EMBEDDING_API_KEY</code> 或 <code class="mono">LLM_API_KEY</code> 可启用 API 嵌入，否则使用本地哈希嵌入（适合开发环境）。</span>
          </div>
        </div>
      </div>

      <div class="panel">
        <div class="panel__head">
          <h3 class="panel__title">
            <span class="panel__title-icon"><el-icon><Search /></el-icon></span>
            检索测试
          </h3>
          <span class="muted">验证混合检索效果</span>
        </div>
        <div class="panel__body">
          <div class="search-bar">
            <el-input v-model="searchQuery" size="large" placeholder="输入一个问题，测试知识库召回…" clearable :prefix-icon="Search" @keyup.enter="handleSearch" />
            <el-button type="primary" size="large" :loading="searching" @click="handleSearch">检索</el-button>
          </div>
          <div v-if="searchResults.length" class="results">
            <div v-for="(r, i) in searchResults" :key="i" class="result">
              <div class="result__head">
                <span class="result__rank num">#{{ i + 1 }}</span>
                <span class="result__path">{{ r.chunk.courseTitle }} <el-icon :size="12"><ArrowRight /></el-icon> {{ r.chunk.chapterTitle }}</span>
                <span class="result__score num">{{ r.score.toFixed(3) }}</span>
              </div>
              <div class="result__heading">{{ r.chunk.heading }}</div>
              <p class="result__text">{{ r.chunk.text }}</p>
              <div class="result__meta">
                <span>向量 <b class="num">{{ r.vectorScore.toFixed(3) }}</b></span>
                <span>关键词 <b class="num">{{ r.keywordScore.toFixed(3) }}</b></span>
              </div>
            </div>
          </div>
          <div v-else-if="searched" class="muted empty-text">未召回任何文本块，试试其他问法。</div>
        </div>
      </div>
    </div>

    <div class="panel">
      <div class="panel__head">
        <h3 class="panel__title">
          <span class="panel__title-icon"><el-icon><Files /></el-icon></span>
          索引文本块
        </h3>
        <span class="muted num">共 {{ total }} 块</span>
      </div>
      <div class="panel__body panel__body--flush">
        <el-table :data="chunks" v-loading="chunksLoading">
          <el-table-column label="课程" prop="courseTitle" width="170" />
          <el-table-column label="章节" prop="chapterTitle" width="170" />
          <el-table-column label="标题" prop="heading" width="190" />
          <el-table-column label="内容摘要" prop="text" show-overflow-tooltip />
          <el-table-column label="嵌入模型" width="150">
            <template #default="{ row }"><span class="mono muted">{{ row.embedModel }}</span></template>
          </el-table-column>
          <el-table-column label="更新时间" width="180">
            <template #default="{ row }"><span class="muted">{{ formatTime(row.updatedAt) }}</span></template>
          </el-table-column>
        </el-table>
      </div>
      <div class="panel__foot pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          background
          layout="total, prev, pager, next"
          @current-change="loadChunks"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Refresh, Setting, Search, Files, InfoFilled, ArrowRight, CircleCheck, Loading, CircleClose, Clock } from '@element-plus/icons-vue'
import { useAuthStore } from '../stores/auth'
import { PERM } from '../constants/permissions'
import { knowledgeApi } from '../api/knowledge.js'
import PageHeader from '../components/common/PageHeader.vue'
import StatCard from '../components/common/StatCard.vue'

const auth = useAuthStore()
const loading = ref(false)
const reindexing = ref(false)
const chunksLoading = ref(false)
const searching = ref(false)
const searched = ref(false)
const searchQuery = ref('')
const searchResults = ref([])
const chunks = ref([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const status = reactive({
  chunkCount: 0,
  courseCount: 0,
  chapterCount: 0,
  embedModel: '',
  embedSource: '',
  dimensions: 0,
  lastIndexedAt: 0,
  indexStatus: 'idle'
})

const statusLabel = computed(() => {
  const map = { ready: '就绪', indexing: '索引中', failed: '失败', idle: '未索引' }
  return map[status.indexStatus] || status.indexStatus
})

const statusTone = computed(() => {
  if (status.indexStatus === 'ready') return 'success'
  if (status.indexStatus === 'failed') return 'danger'
  if (status.indexStatus === 'indexing') return 'warning'
  return 'info'
})

const statusIcon = computed(() => {
  if (status.indexStatus === 'ready') return CircleCheck
  if (status.indexStatus === 'failed') return CircleClose
  if (status.indexStatus === 'indexing') return Loading
  return Clock
})

function formatTime(ts) {
  if (!ts) return '-'
  return new Date(ts).toLocaleString('zh-CN')
}

async function loadStatus() {
  loading.value = true
  try {
    const data = await knowledgeApi.status()
    Object.assign(status, data)
  } finally {
    loading.value = false
  }
}

async function loadChunks() {
  chunksLoading.value = true
  try {
    const data = await knowledgeApi.listChunks({ page: page.value, pageSize: pageSize.value })
    chunks.value = data.list
    total.value = data.total
  } finally {
    chunksLoading.value = false
  }
}

async function handleReindex() {
  await ElMessageBox.confirm('将清空并重建全部知识库向量索引，是否继续？', '重建索引', { type: 'warning' })
  reindexing.value = true
  try {
    const data = await knowledgeApi.reindex()
    Object.assign(status, data)
    ElMessage.success('知识库索引重建完成')
    page.value = 1
    await loadChunks()
  } finally {
    reindexing.value = false
  }
}

async function handleSearch() {
  if (!searchQuery.value.trim()) {
    ElMessage.warning('请输入检索关键词')
    return
  }
  searching.value = true
  try {
    searchResults.value = await knowledgeApi.search(searchQuery.value.trim())
    searched.value = true
  } finally {
    searching.value = false
  }
}

onMounted(async () => {
  await loadStatus()
  await loadChunks()
})
</script>

<style scoped>
.stat-grid--4 {
  grid-template-columns: repeat(4, minmax(0, 1fr));
}

.grid-2 {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 20px;
  align-items: start;
}

.search-bar {
  display: flex;
  gap: 10px;
}

.search-bar .el-input {
  flex: 1;
}

.results {
  display: flex;
  flex-direction: column;
  gap: 10px;
  margin-top: 16px;
  max-height: 420px;
  overflow-y: auto;
  padding-right: 4px;
}

.result {
  padding: 14px 16px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  background: var(--surface-2);
}

.result__head {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 12px;
  color: var(--text-3);
}

.result__rank {
  font-weight: 700;
  color: var(--primary);
}

.result__path {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.result__score {
  padding: 2px 8px;
  border-radius: 999px;
  background: var(--primary-soft);
  color: var(--primary-deep);
  font-weight: 700;
}

.result__heading {
  margin-top: 6px;
  font-weight: 700;
  color: var(--text);
}

.result__text {
  margin: 6px 0 0;
  font-size: 13px;
  line-height: 1.6;
  color: var(--text-2);
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.result__meta {
  display: flex;
  gap: 14px;
  margin-top: 8px;
  font-size: 12px;
  color: var(--text-3);
}

.result__meta b {
  color: var(--text-2);
}

.empty-text {
  margin-top: 14px;
  font-size: 13px;
}

.pagination {
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 1100px) {
  .stat-grid--4 {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .grid-2 {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 700px) {
  .stat-grid--4 {
    grid-template-columns: 1fr;
  }
}
</style>
