<template>
  <div>
    <div class="page-header">
      <h2>题库管理</h2>
      <el-button type="primary" :icon="Plus" @click="openDialog()">新增测验</el-button>
    </div>

    <el-card shadow="never">
      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="title" label="测验名称" min-width="200" />
        <el-table-column prop="courseId" label="关联课程" width="180" />
        <el-table-column prop="questionCount" label="题目数" width="80" />
        <el-table-column label="状态" width="90">
          <template #default="{ row }">
            <el-tag :type="row.status === 'published' ? 'success' : 'info'" size="small">
              {{ row.status === 'published' ? '已发布' : '草稿' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="180" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="$router.push(`/quizzes/${row.id}`)">编辑题目</el-button>
            <el-popconfirm title="确定删除？" @confirm="handleDelete(row.id)">
              <template #reference>
                <el-button link type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" title="新增测验" width="480px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="测验 ID">
          <el-input v-model="form.id" placeholder="英文标识" />
        </el-form-item>
        <el-form-item label="标题">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="关联课程">
          <el-select v-model="form.courseId" style="width: 100%">
            <el-option v-for="c in courses" :key="c.id" :label="c.title" :value="c.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleCreate">创建</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { quizzesApi, coursesApi } from '../api'

const router = useRouter()
const loading = ref(false)
const saving = ref(false)
const list = ref([])
const courses = ref([])
const dialogVisible = ref(false)
const form = reactive({ id: '', title: '', courseId: '', description: '' })

async function loadData() {
  loading.value = true
  try {
    const [quizData, courseData] = await Promise.all([
      quizzesApi.list(),
      coursesApi.list()
    ])
    list.value = quizData.list
    courses.value = courseData.list
  } finally {
    loading.value = false
  }
}

function openDialog() {
  Object.assign(form, { id: '', title: '', courseId: courses.value[0]?.id || '', description: '' })
  dialogVisible.value = true
}

async function handleCreate() {
  if (!form.title || !form.courseId) return ElMessage.warning('请填写必填项')
  saving.value = true
  try {
    const quiz = await quizzesApi.create({ ...form, questions: [], status: 'draft' })
    ElMessage.success('创建成功')
    dialogVisible.value = false
    router.push(`/quizzes/${quiz.id}`)
  } finally {
    saving.value = false
  }
}

async function handleDelete(id) {
  await quizzesApi.remove(id)
  ElMessage.success('删除成功')
  loadData()
}

onMounted(loadData)
</script>
