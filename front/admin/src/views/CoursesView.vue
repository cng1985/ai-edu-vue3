<template>
  <div>
    <div class="page-header">
      <h2>课程管理</h2>
      <el-button v-permission="PERM.COURSE_WRITE" type="primary" :icon="Plus" @click="openDialog()">新增课程</el-button>
    </div>

    <el-card shadow="never">
      <div style="display: flex; gap: 12px; margin-bottom: 16px">
        <el-input v-model="filters.keyword" placeholder="搜索课程名称" clearable style="width: 220px" @keyup.enter="loadData" />
        <el-select v-model="filters.status" placeholder="状态" clearable style="width: 120px" @change="loadData">
          <el-option label="已发布" value="published" />
          <el-option label="草稿" value="draft" />
        </el-select>
        <el-button type="primary" @click="loadData">查询</el-button>
      </div>

      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column label="课程" min-width="260">
          <template #default="{ row }">
            <div style="display: flex; align-items: center; gap: 10px">
              <span style="font-size: 24px">{{ row.icon }}</span>
              <div>
                <div style="font-weight: 600">{{ row.title }}</div>
                <div style="font-size: 12px; color: #9ca3af">{{ row.id }}</div>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="level" label="难度" width="80" />
        <el-table-column prop="chapterCount" label="章节数" width="80" />
        <el-table-column prop="estimatedMinutes" label="时长(分)" width="90" />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 'published' ? 'success' : 'info'" size="small">
              {{ row.status === 'published' ? '已发布' : '草稿' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="$router.push(`/courses/${row.id}`)">编辑</el-button>
            <el-button v-permission="PERM.COURSE_WRITE" link @click="openDialog(row)">设置</el-button>
            <el-popconfirm v-if="auth.hasPermission(PERM.COURSE_DELETE)" title="确定删除该课程？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button link type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

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
import { ref, reactive, onMounted } from 'vue'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { coursesApi } from '../api'
import { useAuthStore } from '../stores/auth'
import { PERM } from '../constants/permissions'

const auth = useAuthStore()
const loading = ref(false)
const saving = ref(false)
const list = ref([])
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
