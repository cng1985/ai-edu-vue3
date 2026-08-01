<template>
  <div>
    <div class="page-header">
      <h2>内容审核</h2>
      <el-radio-group v-model="filters.status" @change="loadData">
        <el-radio-button value="">全部</el-radio-button>
        <el-radio-button value="pending">待审核</el-radio-button>
        <el-radio-button value="approved">已通过</el-radio-button>
        <el-radio-button value="rejected">已驳回</el-radio-button>
      </el-radio-group>
    </div>

    <el-card shadow="never">
      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="title" label="标题" min-width="200" />
        <el-table-column label="类型" width="90">
          <template #default="{ row }">
            <el-tag size="small">{{ row.type === 'chapter' ? '章节' : '题目' }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="submitter" label="提交人" width="100" />
        <el-table-column label="AI 评分" width="90">
          <template #default="{ row }">
            <el-progress
              v-if="row.aiScore"
              :percentage="row.aiScore"
              :color="row.aiScore >= 80 ? '#10b981' : row.aiScore >= 60 ? '#f59e0b' : '#ef4444'"
              :stroke-width="6"
            />
          </template>
        </el-table-column>
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="statusType(row.status)" size="small">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="提交时间" width="170">
          <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDetail(row)">查看</el-button>
            <template v-if="row.status === 'pending'">
              <el-button link type="success" @click="handleApprove(row)">通过</el-button>
              <el-button link type="danger" @click="handleReject(row)">驳回</el-button>
            </template>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-drawer v-model="drawerVisible" title="审核详情" size="520px">
      <template v-if="current">
        <el-descriptions :column="1" border>
          <el-descriptions-item label="标题">{{ current.title }}</el-descriptions-item>
          <el-descriptions-item label="类型">{{ current.type === 'chapter' ? '章节修订' : '新增题目' }}</el-descriptions-item>
          <el-descriptions-item label="提交人">{{ current.submitter }}</el-descriptions-item>
          <el-descriptions-item label="AI 评分">{{ current.aiScore }} 分</el-descriptions-item>
          <el-descriptions-item label="AI 反馈">{{ current.aiFeedback }}</el-descriptions-item>
        </el-descriptions>

        <el-divider>提交内容</el-divider>
        <div class="markdown-preview" v-html="contentPreview"></div>

        <template v-if="current.status === 'pending'">
          <el-divider>审核操作</el-divider>
          <el-input v-model="reviewComment" type="textarea" :rows="2" placeholder="审核意见（可选）" />
          <div style="margin-top: 16px; display: flex; gap: 12px">
            <el-button type="success" :loading="acting" @click="handleApprove(current)">通过</el-button>
            <el-button type="danger" :loading="acting" @click="handleReject(current)">驳回</el-button>
          </div>
        </template>

        <template v-if="current.comment">
          <el-divider>审核意见</el-divider>
          <p>{{ current.comment }}</p>
        </template>
      </template>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { marked } from 'marked'
import { ElMessage, ElMessageBox } from 'element-plus'
import { reviewsApi } from '../api'

const loading = ref(false)
const acting = ref(false)
const list = ref([])
const filters = reactive({ status: 'pending' })
const drawerVisible = ref(false)
const current = ref(null)
const reviewComment = ref('')

const contentPreview = computed(() => {
  if (!current.value) return ''
  try {
    const content = current.value.content || ''
    if (current.value.type === 'quiz') {
      const q = JSON.parse(content)
      return marked.parse(`**${q.text}**\n\n${q.options?.map((o, i) => `${String.fromCharCode(65 + i)}. ${o}`).join('\n') || ''}`)
    }
    return marked.parse(content)
  } catch {
    return `<pre>${current.value?.content}</pre>`
  }
})

function statusLabel(s) {
  return { pending: '待审核', approved: '已通过', rejected: '已驳回' }[s] || s
}
function statusType(s) {
  return { pending: 'warning', approved: 'success', rejected: 'danger' }[s] || 'info'
}
function formatDate(ts) {
  return new Date(ts).toLocaleString('zh-CN')
}

async function loadData() {
  loading.value = true
  try {
    const data = await reviewsApi.list(filters)
    list.value = data.list
  } finally {
    loading.value = false
  }
}

function openDetail(row) {
  current.value = row
  reviewComment.value = ''
  drawerVisible.value = true
}

async function handleApprove(row) {
  acting.value = true
  try {
    await reviewsApi.approve(row.id, reviewComment.value)
    ElMessage.success('审核通过')
    drawerVisible.value = false
    loadData()
  } finally {
    acting.value = false
  }
}

async function handleReject(row) {
  try {
    const { value } = await ElMessageBox.prompt('请输入驳回原因', '驳回审核', {
      inputValue: reviewComment.value || '内容不符合发布标准',
      confirmButtonText: '确认驳回',
      cancelButtonText: '取消'
    })
    acting.value = true
    await reviewsApi.reject(row.id, value)
    ElMessage.success('已驳回')
    drawerVisible.value = false
    loadData()
  } catch {
    /* 取消 */
  } finally {
    acting.value = false
  }
}

onMounted(loadData)
</script>
