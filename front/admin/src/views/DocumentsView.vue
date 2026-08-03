<template>
  <div>
    <div class="page-header">
      <h2>单据管理</h2>
      <div style="display: flex; gap: 8px">
        <el-button v-permission="PERM.DOCUMENT_EXPORT" :icon="Download" @click="handleExport">导出 Excel</el-button>
        <el-button v-permission="PERM.DOCUMENT_IMPORT" :icon="Upload" @click="openImportDialog">导入 Excel</el-button>
        <el-button v-permission="PERM.DOCUMENT_WRITE" type="primary" :icon="Plus" @click="openDialog()">新增单据</el-button>
      </div>
    </div>

    <el-card shadow="never">
      <div style="display: flex; gap: 12px; margin-bottom: 16px; flex-wrap: wrap">
        <el-input v-model="filters.keyword" placeholder="搜索编号/标题" clearable style="width: 220px" @clear="loadData" @keyup.enter="loadData" />
        <el-select v-model="filters.type" placeholder="类型" clearable style="width: 140px" @change="loadData">
          <el-option v-for="t in typeOptions" :key="t.value" :label="t.label" :value="t.value" />
        </el-select>
        <el-select v-model="filters.status" placeholder="状态" clearable style="width: 120px" @change="loadData">
          <el-option v-for="s in statusOptions" :key="s.value" :label="s.label" :value="s.value" />
        </el-select>
        <el-button type="primary" @click="loadData">查询</el-button>
      </div>

      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="docNo" label="单据编号" width="160" />
        <el-table-column prop="title" label="标题" min-width="180" show-overflow-tooltip />
        <el-table-column label="类型" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ typeLabel(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="金额" width="120" align="right">
          <template #default="{ row }">¥{{ row.amount.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="statusTagType(row.status)" size="small">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdBy" label="创建人" width="100" />
        <el-table-column label="创建时间" width="170">
          <template #default="{ row }">{{ formatDate(row.createdAt) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openDialog(row)">编辑</el-button>
            <el-popconfirm v-if="auth.hasPermission(PERM.DOCUMENT_DELETE)" title="确定删除该单据？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button link type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        v-if="total > pageSize"
        style="margin-top: 16px; justify-content: flex-end"
        layout="total, prev, pager, next"
        :total="total"
        :page-size="pageSize"
        v-model:current-page="page"
        @current-change="loadData"
      />
    </el-card>

    <!-- 新增/编辑对话框 -->
    <el-dialog v-model="dialogVisible" :title="editing ? '编辑单据' : '新增单据'" width="520px">
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="90px">
        <el-form-item label="单据编号" prop="docNo">
          <el-input v-model="form.docNo" :disabled="editing" placeholder="如 DOC-2026-001" />
        </el-form-item>
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="类型" prop="type">
          <el-select v-model="form.type" style="width: 100%">
            <el-option v-for="t in typeOptions" :key="t.value" :label="t.label" :value="t.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="金额" prop="amount">
          <el-input-number v-model="form.amount" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="form.status" style="width: 100%">
            <el-option v-for="s in statusOptions" :key="s.value" :label="s.label" :value="s.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>

    <!-- 导入对话框 -->
    <el-dialog v-model="importVisible" title="导入 Excel" width="520px" :close-on-click-modal="!importing" :close-on-press-escape="!importing" :show-close="!importing">
      <div v-if="!importing && !importDone">
        <el-alert type="info" :closable="false" style="margin-bottom: 16px">
          请上传 .xlsx 格式的 Excel 文件。可先
          <el-link type="primary" @click="downloadTemplate">下载导入模板</el-link>
          查看格式要求。
        </el-alert>
        <el-upload
          ref="uploadRef"
          drag
          :auto-upload="false"
          :limit="1"
          accept=".xlsx,.xls"
          :on-change="onFileChange"
          :on-exceed="() => ElMessage.warning('只能上传一个文件')"
        >
          <el-icon :size="48" style="color: #c0c4cc"><Upload /></el-icon>
          <div>将文件拖到此处，或<em>点击上传</em></div>
          <template #tip>
            <div style="color: #909399; font-size: 12px">仅支持 .xlsx / .xls 格式</div>
          </template>
        </el-upload>
      </div>

      <div v-if="importing || importDone">
        <div style="margin-bottom: 16px; text-align: center">
          <el-progress
            :percentage="importProgress"
            :status="importDone ? (importResult?.failed > 0 ? 'warning' : 'success') : undefined"
            :stroke-width="18"
            striped
            striped-flow
          />
        </div>
        <el-descriptions :column="2" border size="small">
          <el-descriptions-item label="总行数">{{ importResult?.total || 0 }}</el-descriptions-item>
          <el-descriptions-item label="已处理">{{ importResult?.current || 0 }}</el-descriptions-item>
          <el-descriptions-item label="成功">
            <span style="color: #67c23a">{{ importResult?.success || 0 }}</span>
          </el-descriptions-item>
          <el-descriptions-item label="失败">
            <span style="color: #f56c6c">{{ importResult?.failed || 0 }}</span>
          </el-descriptions-item>
        </el-descriptions>
        <div v-if="importResult?.errors?.length" style="margin-top: 12px">
          <el-alert type="error" :closable="false" title="错误详情">
            <div style="max-height: 120px; overflow-y: auto; font-size: 12px">
              <div v-for="(err, i) in importResult.errors" :key="i">{{ err }}</div>
            </div>
          </el-alert>
        </div>
        <div v-if="importDone && importResult?.message" style="margin-top: 12px; text-align: center; color: #606266">
          {{ importResult.message }}
        </div>
      </div>

      <template #footer>
        <el-button v-if="!importing" @click="closeImportDialog">{{ importDone ? '关闭' : '取消' }}</el-button>
        <el-button v-if="!importing && !importDone" type="primary" :disabled="!selectedFile" @click="startImport">开始导入</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, onUnmounted } from 'vue'
import { Plus, Upload, Download } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { documentsApi } from '../api'
import { useAuthStore } from '../stores/auth'
import { PERM } from '../constants/permissions'

const auth = useAuthStore()
const loading = ref(false)
const saving = ref(false)
const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = 20
const filters = reactive({ keyword: '', type: '', status: '' })

const typeOptions = [
  { label: '采购单', value: 'purchase' },
  { label: '销售单', value: 'sales' },
  { label: '报销单', value: 'expense' },
  { label: '其他', value: 'other' }
]
const statusOptions = [
  { label: '草稿', value: 'draft' },
  { label: '待审核', value: 'pending' },
  { label: '已通过', value: 'approved' },
  { label: '已驳回', value: 'rejected' }
]

const typeMap = Object.fromEntries(typeOptions.map(t => [t.value, t.label]))
const statusMap = Object.fromEntries(statusOptions.map(s => [s.value, s.label]))
function typeLabel(v) { return typeMap[v] || v }
function statusLabel(v) { return statusMap[v] || v }
function statusTagType(v) {
  if (v === 'approved') return 'success'
  if (v === 'pending') return 'warning'
  if (v === 'rejected') return 'danger'
  return 'info'
}
function formatDate(ts) { return new Date(ts).toLocaleString('zh-CN') }

// CRUD
const dialogVisible = ref(false)
const editing = ref(false)
const editingId = ref('')
const formRef = ref()
const form = reactive({ docNo: '', title: '', type: 'purchase', amount: 0, status: 'draft', remark: '' })
const formRules = {
  docNo: [{ required: true, message: '请输入单据编号', trigger: 'blur' }],
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }]
}

async function loadData() {
  loading.value = true
  try {
    const data = await documentsApi.list({ ...filters, page: page.value, pageSize })
    list.value = data.list
    total.value = data.total
  } finally {
    loading.value = false
  }
}

function openDialog(row) {
  editing.value = Boolean(row)
  editingId.value = row?.id || ''
  Object.assign(form, {
    docNo: row?.docNo || '',
    title: row?.title || '',
    type: row?.type || 'purchase',
    amount: row?.amount || 0,
    status: row?.status || 'draft',
    remark: row?.remark || ''
  })
  dialogVisible.value = true
}

async function handleSave() {
  await formRef.value.validate()
  saving.value = true
  try {
    if (editing.value) {
      await documentsApi.update(editingId.value, { ...form })
    } else {
      await documentsApi.create(form)
    }
    ElMessage.success('保存成功')
    dialogVisible.value = false
    loadData()
  } finally {
    saving.value = false
  }
}

async function handleDelete(id) {
  await documentsApi.remove(id)
  ElMessage.success('删除成功')
  loadData()
}

// Export
async function handleExport() {
  try {
    await documentsApi.exportExcel(filters)
    ElMessage.success('导出成功')
  } catch { /* handled by interceptor */ }
}

async function downloadTemplate() {
  await documentsApi.downloadTemplate()
}

// Import
const importVisible = ref(false)
const importing = ref(false)
const importDone = ref(false)
const selectedFile = ref(null)
const uploadRef = ref()
const importResult = ref(null)
let pollTimer = null

const importProgress = computed(() => {
  if (!importResult.value || !importResult.value.total) return 0
  return Math.round((importResult.value.current / importResult.value.total) * 100)
})

function openImportDialog() {
  importVisible.value = true
  importing.value = false
  importDone.value = false
  selectedFile.value = null
  importResult.value = null
  uploadRef.value?.clearFiles()
}

function closeImportDialog() {
  importVisible.value = false
  if (importDone.value) loadData()
}

function onFileChange(file) {
  selectedFile.value = file.raw
}

async function startImport() {
  if (!selectedFile.value) return
  importing.value = true
  importDone.value = false
  importResult.value = { total: 0, current: 0, success: 0, failed: 0 }

  try {
    const { taskId } = await documentsApi.importExcel(selectedFile.value)
    pollTimer = setInterval(async () => {
      try {
        const progress = await documentsApi.importProgress(taskId)
        importResult.value = progress
        if (progress.status === 'completed' || progress.status === 'failed') {
          clearInterval(pollTimer)
          pollTimer = null
          importing.value = false
          importDone.value = true
        }
      } catch {
        clearInterval(pollTimer)
        pollTimer = null
        importing.value = false
      }
    }, 300)
  } catch {
    importing.value = false
  }
}

onMounted(loadData)
onUnmounted(() => { if (pollTimer) clearInterval(pollTimer) })
</script>
