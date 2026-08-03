<template>
  <div>
    <div class="page-header">
      <h2>知识库管理</h2>
      <p>管理课程内容的向量索引，支持语义检索增强 AI 对话。</p>
    </div>

    <el-row :gutter="16" class="stats-row">
      <el-col :span="6">
        <el-card shadow="never">
          <div class="stat-label">文本块数量</div>
          <div class="stat-value">{{ status.chunkCount }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never">
          <div class="stat-label">覆盖课程</div>
          <div class="stat-value">{{ status.courseCount }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never">
          <div class="stat-label">覆盖章节</div>
          <div class="stat-value">{{ status.chapterCount }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="never">
          <div class="stat-label">索引状态</div>
          <div class="stat-value">
            <el-tag :type="statusTagType" size="small">{{ statusLabel }}</el-tag>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card v-loading="loading" shadow="never" class="info-card">
      <template #header>
        <div class="card-head">
          <span>索引配置</span>
          <el-button
            v-if="auth.hasPermission(PERM.KNOWLEDGE_MANAGE)"
            type="primary"
            :loading="reindexing"
            @click="handleReindex"
          >
            重建索引
          </el-button>
        </div>
      </template>
      <el-descriptions :column="2" border>
        <el-descriptions-item label="嵌入模型">{{ status.embedModel || '-' }}</el-descriptions-item>
        <el-descriptions-item label="嵌入来源">
          <el-tag size="small" :type="status.embedSource === 'api' ? 'success' : 'info'">
            {{ status.embedSource === 'api' ? 'API 嵌入' : '本地哈希嵌入' }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="向量维度">{{ status.dimensions || '-' }}</el-descriptions-item>
        <el-descriptions-item label="最后索引时间">{{ formatTime(status.lastIndexedAt) }}</el-descriptions-item>
      </el-descriptions>
      <el-alert
        class="hint-alert"
        title="说明"
        type="info"
        :closable="false"
        show-icon
        description="知识库使用 SQLite 嵌入式向量存储，结合向量相似度与关键词混合检索。配置 EMBEDDING_API_KEY 或 LLM_API_KEY 可启用 API 嵌入，否则使用本地哈希嵌入（适合开发环境）。"
      />
    </el-card>

    <el-card shadow="never" class="search-card">
      <template #header>
        <span>检索测试</span>
      </template>
      <div class="search-bar">
        <el-input v-model="searchQuery" placeholder="输入问题测试知识库检索效果" clearable @keyup.enter="handleSearch" />
        <el-button type="primary" :loading="searching" @click="handleSearch">检索</el-button>
      </div>
      <el-table v-if="searchResults.length" :data="searchResults" stripe style="margin-top: 16px">
        <el-table-column label="课程" prop="chunk.courseTitle" width="160" />
        <el-table-column label="章节" prop="chunk.chapterTitle" width="160" />
        <el-table-column label="标题" prop="chunk.heading" width="180" />
        <el-table-column label="内容摘要" prop="chunk.text" show-overflow-tooltip />
        <el-table-column label="综合分" width="80">
          <template #default="{ row }">{{ row.score.toFixed(3) }}</template>
        </el-table-column>
        <el-table-column label="向量分" width="80">
          <template #default="{ row }">{{ row.vectorScore.toFixed(3) }}</template>
        </el-table-column>
        <el-table-column label="关键词分" width="90">
          <template #default="{ row }">{{ row.keywordScore.toFixed(3) }}</template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-card shadow="never" class="chunks-card">
      <template #header>
        <span>索引文本块</span>
      </template>
      <el-table :data="chunks" stripe v-loading="chunksLoading">
        <el-table-column label="课程" prop="courseTitle" width="160" />
        <el-table-column label="章节" prop="chapterTitle" width="160" />
        <el-table-column label="标题" prop="heading" width="180" />
        <el-table-column label="内容摘要" prop="text" show-overflow-tooltip />
        <el-table-column label="嵌入模型" prop="embedModel" width="140" />
        <el-table-column label="更新时间" width="170">
          <template #default="{ row }">{{ formatTime(row.updatedAt) }}</template>
        </el-table-column>
      </el-table>
      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="loadChunks"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useAuthStore } from '../stores/auth'
import { PERM } from '../constants/permissions'
import { knowledgeApi } from '../api/knowledge.js'

const auth = useAuthStore()
const loading = ref(false)
const reindexing = ref(false)
const chunksLoading = ref(false)
const searching = ref(false)
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

const statusTagType = computed(() => {
  if (status.indexStatus === 'ready') return 'success'
  if (status.indexStatus === 'failed') return 'danger'
  if (status.indexStatus === 'indexing') return 'warning'
  return 'info'
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
.page-header { margin-bottom: 20px; }
.page-header h2 { margin: 0 0 8px; }
.page-header p { margin: 0; color: #64748b; }
.stats-row { margin-bottom: 16px; }
.stat-label { font-size: 13px; color: #64748b; margin-bottom: 4px; }
.stat-value { font-size: 24px; font-weight: 600; color: #1e1b4b; }
.card-head { display: flex; justify-content: space-between; align-items: center; }
.info-card, .search-card, .chunks-card { margin-bottom: 16px; }
.hint-alert { margin-top: 16px; }
.search-bar { display: flex; gap: 12px; }
.search-bar .el-input { flex: 1; }
.pagination { margin-top: 16px; display: flex; justify-content: flex-end; }
</style>
