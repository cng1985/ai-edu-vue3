<template>
  <div class="page">
    <PageHeader eyebrow="内容管理" title="内容审核" subtitle="结合 AI 评分对章节修订与新增题目进行人工审核。">
      <div class="segmented">
        <button
          v-for="opt in statusOptions"
          :key="opt.value"
          type="button"
          class="segmented__item"
          :class="{ 'segmented__item--active': filters.status === opt.value }"
          @click="filters.status = opt.value; loadData()"
        >
          {{ opt.label }}
        </button>
      </div>
    </PageHeader>

    <div class="panel">
      <div class="panel__body panel__body--flush">
        <el-table :data="list" v-loading="loading">
          <el-table-column label="内容" min-width="260">
            <template #default="{ row }">
              <div class="cell">
                <span class="cell__icon" :class="row.type === 'chapter' ? 'type-icon--chapter' : 'type-icon--quiz'">
                  <el-icon :size="18"><component :is="row.type === 'chapter' ? Notebook : EditPen" /></el-icon>
                </span>
                <div class="cell__main">
                  <div class="cell__title">{{ row.title }}</div>
                  <div class="cell__sub">{{ row.type === 'chapter' ? '章节修订' : '新增题目' }} · {{ row.submitter }}</div>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="AI 评分" width="200">
            <template #default="{ row }">
              <div v-if="row.aiScore" class="score">
                <div class="score__track">
                  <span :class="`score__bar--${scoreTone(row.aiScore)}`" :style="{ width: row.aiScore + '%' }" />
                </div>
                <span class="score__num num" :class="`score__num--${scoreTone(row.aiScore)}`">{{ row.aiScore }}</span>
              </div>
              <span v-else class="muted">—</span>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="110">
            <template #default="{ row }">
              <el-tag :type="statusType(row.status)" size="small">{{ statusLabel(row.status) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="提交时间" width="180">
            <template #default="{ row }"><span class="muted">{{ formatDate(row.createdAt) }}</span></template>
          </el-table-column>
          <el-table-column label="操作" width="200" fixed="right" align="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="openDetail(row)">查看</el-button>
              <template v-if="row.status === 'pending'">
                <el-button link type="success" @click="handleApprove(row)">通过</el-button>
                <el-button link type="danger" @click="handleReject(row)">驳回</el-button>
              </template>
            </template>
          </el-table-column>
          <template #empty>
            <el-empty description="当前筛选下没有审核项" :image-size="90" />
          </template>
        </el-table>
      </div>
    </div>

    <el-drawer v-model="drawerVisible" title="审核详情" size="560px">
      <template v-if="current">
        <div class="detail-head">
          <span class="cell__icon" :class="current.type === 'chapter' ? 'type-icon--chapter' : 'type-icon--quiz'">
            <el-icon :size="20"><component :is="current.type === 'chapter' ? Notebook : EditPen" /></el-icon>
          </span>
          <div>
            <h3>{{ current.title }}</h3>
            <p class="muted">{{ current.type === 'chapter' ? '章节修订' : '新增题目' }} · 提交人 {{ current.submitter }} · {{ formatDate(current.createdAt) }}</p>
          </div>
          <el-tag :type="statusType(current.status)" size="small">{{ statusLabel(current.status) }}</el-tag>
        </div>

        <div class="ai-card">
          <div class="ai-card__score">
            <strong class="num" :class="`score__num--${scoreTone(current.aiScore)}`">{{ current.aiScore ?? '—' }}</strong>
            <span>AI 评分</span>
          </div>
          <div class="ai-card__text">
            <div class="ai-card__title">AI 审核反馈</div>
            <p>{{ current.aiFeedback || '暂无 AI 反馈' }}</p>
          </div>
        </div>

        <div class="section-title">提交内容</div>
        <div class="markdown-preview" v-html="contentPreview"></div>

        <template v-if="current.status === 'pending'">
          <div class="section-title">审核操作</div>
          <el-input v-model="reviewComment" type="textarea" :rows="3" placeholder="审核意见（可选）" />
          <div class="actions">
            <el-button type="success" :loading="acting" @click="handleApprove(current)">
              <el-icon class="el-icon--left"><Check /></el-icon>通过
            </el-button>
            <el-button type="danger" plain :loading="acting" @click="handleReject(current)">
              <el-icon class="el-icon--left"><Close /></el-icon>驳回
            </el-button>
          </div>
        </template>

        <template v-if="current.comment">
          <div class="section-title">审核意见</div>
          <div class="hint hint--info">{{ current.comment }}</div>
        </template>
      </template>
    </el-drawer>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { marked } from 'marked'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Notebook, EditPen, Check, Close } from '@element-plus/icons-vue'
import { reviewsApi } from '../api'
import PageHeader from '../components/common/PageHeader.vue'

const statusOptions = [
  { label: '全部', value: '' },
  { label: '待审核', value: 'pending' },
  { label: '已通过', value: 'approved' },
  { label: '已驳回', value: 'rejected' }
]

function scoreTone(score) {
  if (score >= 80) return 'success'
  if (score >= 60) return 'warning'
  return 'danger'
}

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

<style scoped>
.segmented {
  display: inline-flex;
  gap: 4px;
  padding: 4px;
  border-radius: 12px;
  background: var(--surface-3);
}

.segmented__item {
  padding: 7px 14px;
  border: none;
  border-radius: 9px;
  background: transparent;
  color: var(--text-2);
  font-size: 13px;
  font-weight: 600;
  cursor: pointer;
  transition: all var(--t-fast);
}

.segmented__item:hover {
  color: var(--text);
}

.segmented__item--active {
  background: var(--surface);
  color: var(--primary-deep);
  box-shadow: var(--shadow-xs);
}

.type-icon--chapter {
  background: var(--primary-soft);
  color: var(--primary);
}

.type-icon--quiz {
  background: var(--sky-soft);
  color: var(--sky);
}

.score {
  display: flex;
  align-items: center;
  gap: 10px;
}

.score__track {
  flex: 1;
  height: 7px;
  border-radius: 999px;
  background: var(--surface-3);
  overflow: hidden;
}

.score__track span {
  display: block;
  height: 100%;
  border-radius: 999px;
}

.score__bar--success { background: var(--success); }
.score__bar--warning { background: var(--warning); }
.score__bar--danger { background: var(--danger); }

.score__num {
  width: 28px;
  text-align: right;
  font-weight: 700;
}

.score__num--success { color: var(--success-strong); }
.score__num--warning { color: var(--warning-strong); }
.score__num--danger { color: var(--danger-strong); }

.detail-head {
  display: flex;
  align-items: flex-start;
  gap: 14px;
}

.detail-head > div {
  flex: 1;
  min-width: 0;
}

.detail-head h3 {
  margin: 0;
  font-size: 17px;
}

.detail-head p {
  margin: 4px 0 0;
  font-size: 12.5px;
}

.ai-card {
  display: flex;
  gap: 18px;
  margin-top: 20px;
  padding: 16px 18px;
  border-radius: var(--radius-sm);
  background: linear-gradient(135deg, var(--primary-soft) 0%, #f6f4ff 100%);
  border: 1px solid var(--primary-soft-2);
}

.ai-card__score {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-width: 70px;
  padding-right: 18px;
  border-right: 1px solid var(--primary-soft-2);
}

.ai-card__score strong {
  font-size: 30px;
  letter-spacing: -0.02em;
  line-height: 1;
}

.ai-card__score span {
  margin-top: 4px;
  font-size: 11.5px;
  color: var(--text-3);
}

.ai-card__title {
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: var(--primary);
}

.ai-card__text p {
  margin: 6px 0 0;
  font-size: 13.5px;
  line-height: 1.7;
  color: var(--text-2);
}

.section-title {
  margin: 24px 0 10px;
  font-size: 12.5px;
  font-weight: 700;
  letter-spacing: 0.06em;
  text-transform: uppercase;
  color: var(--text-3);
}

.actions {
  display: flex;
  gap: 10px;
  margin-top: 14px;
}
</style>
