<template>
  <div class="page">
    <PageHeader eyebrow="内容管理" title="题库管理" subtitle="为课程创建测验、维护题目与发布状态。">
      <el-button type="primary" :icon="Plus" @click="openDialog()">新增测验</el-button>
    </PageHeader>

    <div class="stat-grid stat-grid--3">
      <StatCard label="测验套数" :value="list.length" icon="EditPen" tone="primary" />
      <StatCard label="题目总数" :value="questionTotal" icon="Tickets" tone="sky" />
      <StatCard label="已发布" :value="publishedCount" icon="Promotion" tone="success" />
    </div>

    <div class="panel">
      <div class="panel__head">
        <h3 class="panel__title">
          <span class="panel__title-icon"><el-icon><EditPen /></el-icon></span>
          全部测验
        </h3>
        <span class="muted">共 {{ list.length }} 套</span>
      </div>
      <div class="panel__body panel__body--flush">
        <el-table :data="list" v-loading="loading">
          <el-table-column label="测验" min-width="260">
            <template #default="{ row }">
              <div class="cell">
                <span class="cell__icon quiz-icon"><el-icon :size="18"><EditPen /></el-icon></span>
                <div class="cell__main">
                  <div class="cell__title">{{ row.title }}</div>
                  <div class="cell__sub mono">{{ row.id }}</div>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="关联课程" min-width="180">
            <template #default="{ row }">
              <span class="course-chip">{{ courseTitle(row.courseId) }}</span>
            </template>
          </el-table-column>
          <el-table-column label="题目数" width="100">
            <template #default="{ row }"><span class="num">{{ row.questionCount }}</span></template>
          </el-table-column>
          <el-table-column label="状态" width="100">
            <template #default="{ row }">
              <el-tag :type="row.status === 'published' ? 'success' : 'info'" size="small">
                {{ row.status === 'published' ? '已发布' : '草稿' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="170" fixed="right" align="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="$router.push(`/quizzes/${row.id}`)">编辑题目</el-button>
              <el-popconfirm title="确定删除？" @confirm="handleDelete(row.id)">
                <template #reference>
                  <el-button link type="danger">删除</el-button>
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
          <template #empty>
            <el-empty description="暂无测验，点击右上角新增" :image-size="90" />
          </template>
        </el-table>
      </div>
    </div>

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
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Plus, EditPen } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { quizzesApi, coursesApi } from '../api'
import PageHeader from '../components/common/PageHeader.vue'
import StatCard from '../components/common/StatCard.vue'

const router = useRouter()
const loading = ref(false)
const saving = ref(false)
const list = ref([])
const courses = ref([])

const questionTotal = computed(() => list.value.reduce((sum, q) => sum + (q.questionCount || 0), 0))
const publishedCount = computed(() => list.value.filter((q) => q.status === 'published').length)
function courseTitle(id) {
  return courses.value.find((c) => c.id === id)?.title || id || '—'
}
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

<style scoped>
.stat-grid--3 {
  grid-template-columns: repeat(3, minmax(0, 1fr));
}

.quiz-icon {
  background: var(--sky-soft);
  color: var(--sky);
}

.course-chip {
  display: inline-flex;
  align-items: center;
  padding: 3px 10px;
  border-radius: 8px;
  background: var(--surface-3);
  color: var(--text-2);
  font-size: 12.5px;
  font-weight: 600;
}

@media (max-width: 900px) {
  .stat-grid--3 {
    grid-template-columns: 1fr;
  }
}
</style>
