<template>
  <div class="page">
    <PageHeader eyebrow="内容管理" title="课程管理" subtitle="维护课程基础信息、章节内容与发布状态。">
      <el-button v-permission="PERM.COURSE_WRITE" type="primary" :icon="Plus" @click="openDialog()">新增课程</el-button>
    </PageHeader>

    <div class="stat-grid stat-grid--3">
      <StatCard label="课程总数" :value="list.length" icon="Reading" tone="primary" />
      <StatCard label="已发布" :value="publishedCount" icon="Promotion" tone="success" />
      <StatCard label="草稿" :value="list.length - publishedCount" icon="Document" tone="warning" />
    </div>

    <div class="panel">
      <div class="panel__head">
        <div class="toolbar">
          <el-input v-model="filters.keyword" placeholder="搜索课程名称" clearable :prefix-icon="Search" style="width: 240px" @keyup.enter="loadData" @clear="loadData" />
          <el-select v-model="filters.status" placeholder="全部状态" clearable style="width: 130px" @change="loadData">
            <el-option label="已发布" value="published" />
            <el-option label="草稿" value="draft" />
          </el-select>
          <el-button @click="loadData">查询</el-button>
        </div>
        <span class="muted">共 {{ list.length }} 门课程</span>
      </div>

      <div class="panel__body panel__body--flush">
        <el-table :data="list" v-loading="loading">
          <el-table-column label="课程" min-width="280">
            <template #default="{ row }">
              <div class="cell">
                <span class="cell__icon" :style="{ background: (row.accent || '#6b5cff') + '1f' }">{{ row.icon }}</span>
                <div class="cell__main">
                  <div class="cell__title">{{ row.title }}</div>
                  <div class="cell__sub mono">{{ row.id }}</div>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="难度" width="100">
            <template #default="{ row }">
              <span class="level" :class="`level--${row.level}`">{{ row.level }}</span>
            </template>
          </el-table-column>
          <el-table-column label="章节" width="90">
            <template #default="{ row }"><span class="num">{{ row.chapterCount }}</span></template>
          </el-table-column>
          <el-table-column label="预计时长" width="110">
            <template #default="{ row }"><span class="num">{{ row.estimatedMinutes }}</span> 分钟</template>
          </el-table-column>
          <el-table-column label="标签" min-width="160">
            <template #default="{ row }">
              <div class="tags">
                <el-tag v-for="t in (row.tags || []).slice(0, 3)" :key="t" size="small" type="info" effect="plain">{{ t }}</el-tag>
                <span v-if="!row.tags?.length" class="muted">—</span>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 'published' ? 'success' : 'info'" size="small">
                {{ row.status === 'published' ? '已发布' : '草稿' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="200" fixed="right" align="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="$router.push(`/courses/${row.id}`)">编辑章节</el-button>
              <el-button v-permission="PERM.COURSE_WRITE" link @click="openDialog(row)">设置</el-button>
              <el-popconfirm v-if="auth.hasPermission(PERM.COURSE_DELETE)" title="确定删除该课程？" @confirm="handleDelete(row.id)">
                <template #reference>
                  <el-button link type="danger">删除</el-button>
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
          <template #empty>
            <el-empty description="暂无课程，点击右上角新增" :image-size="90" />
          </template>
        </el-table>
      </div>
    </div>

    <el-dialog v-model="dialogVisible" :title="editing ? '编辑课程' : '新增课程'" width="560px">
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="100px">
        <el-form-item v-if="!editing" label="课程 ID" prop="id">
          <el-input v-model="form.id" placeholder="英文标识，如 my-course" />
        </el-form-item>
        <el-form-item label="标题" prop="title">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="难度">
          <el-select v-model="form.level" style="width: 100%">
            <el-option label="入门" value="入门" />
            <el-option label="进阶" value="进阶" />
            <el-option label="高级" value="高级" />
          </el-select>
        </el-form-item>
        <el-form-item label="图标">
          <el-input v-model="form.icon" style="width: 80px" />
        </el-form-item>
        <el-form-item label="主题色">
          <el-color-picker v-model="form.accent" />
        </el-form-item>
        <el-form-item label="预计时长">
          <el-input-number v-model="form.estimatedMinutes" :min="1" />
          <span style="margin-left: 8px; color: #9ca3af">分钟</span>
        </el-form-item>
        <el-form-item label="标签">
          <el-select v-model="form.tags" multiple filterable allow-create style="width: 100%" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="form.status" style="width: 100%">
            <el-option label="草稿" value="draft" />
            <el-option label="已发布" value="published" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { Plus, Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { coursesApi } from '../api'
import { useAuthStore } from '../stores/auth'
import { PERM } from '../constants/permissions'
import PageHeader from '../components/common/PageHeader.vue'
import StatCard from '../components/common/StatCard.vue'

const auth = useAuthStore()
const loading = ref(false)
const saving = ref(false)
const list = ref([])
const publishedCount = computed(() => list.value.filter((c) => c.status === 'published').length)
const filters = reactive({ keyword: '', status: '' })

const dialogVisible = ref(false)
const editing = ref(false)
const editingId = ref('')
const formRef = ref()
const form = reactive({
  id: '', title: '', description: '', level: '入门', icon: '📚',
  accent: '#6366f1', estimatedMinutes: 60, tags: [], status: 'draft'
})
const formRules = {
  id: [{ required: true, message: '请输入课程 ID', trigger: 'blur' }],
  title: [{ required: true, message: '请输入标题', trigger: 'blur' }]
}

async function loadData() {
  loading.value = true
  try {
    const data = await coursesApi.list(filters)
    list.value = data.list
  } finally {
    loading.value = false
  }
}

function openDialog(row) {
  editing.value = Boolean(row)
  editingId.value = row?.id || ''
  Object.assign(form, {
    id: row?.id || '',
    title: row?.title || '',
    description: row?.description || '',
    level: row?.level || '入门',
    icon: row?.icon || '📚',
    accent: row?.accent || '#6366f1',
    estimatedMinutes: row?.estimatedMinutes || 60,
    tags: row?.tags ? [...row.tags] : [],
    status: row?.status || 'draft'
  })
  dialogVisible.value = true
}

async function handleSave() {
  await formRef.value.validate()
  saving.value = true
  try {
    if (editing.value) {
      const { id, ...payload } = form
      await coursesApi.update(editingId.value, payload)
    } else {
      await coursesApi.create(form)
    }
    ElMessage.success('保存成功')
    dialogVisible.value = false
    loadData()
  } finally {
    saving.value = false
  }
}

async function handleDelete(id) {
  await coursesApi.remove(id)
  ElMessage.success('删除成功')
  loadData()
}

onMounted(loadData)
</script>

<style scoped>
.stat-grid--3 {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
}

.level {
  display: inline-flex;
  align-items: center;
  padding: 2px 10px;
  border-radius: 999px;
  font-size: 12px;
  font-weight: 700;
  background: var(--success-soft);
  color: var(--success-strong);
}

.level--进阶 {
  background: var(--warning-soft);
  color: var(--warning-strong);
}

.level--高级 {
  background: var(--danger-soft);
  color: var(--danger-strong);
}

@media (max-width: 900px) {
  .stat-grid--3 {
    grid-template-columns: 1fr;
  }
}
</style>
