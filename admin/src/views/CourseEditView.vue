<template>
  <div v-loading="loading">
    <div class="page-header">
      <div style="display: flex; align-items: center; gap: 12px">
        <el-button :icon="ArrowLeft" @click="$router.push('/courses')">返回</el-button>
        <h2>{{ course?.icon }} {{ course?.title || '课程编辑' }}</h2>
        <el-tag :type="course?.status === 'published' ? 'success' : 'info'" size="small">
          {{ course?.status === 'published' ? '已发布' : '草稿' }}
        </el-tag>
      </div>
      <el-button type="primary" :icon="Plus" @click="openChapterDialog()">新增章节</el-button>
    </div>

    <el-row :gutter="16">
      <el-col :span="8">
        <el-card shadow="never">
          <template #header><span>章节列表 ({{ course?.chapters?.length || 0 }})</span></template>
          <div
            v-for="(ch, idx) in course?.chapters || []"
            :key="ch.id"
            class="chapter-item"
            :class="{ active: selectedChapter?.id === ch.id }"
            @click="selectChapter(ch)"
          >
            <div class="chapter-item__title">{{ idx + 1 }}. {{ ch.title }}</div>
            <div class="chapter-item__meta">
              <el-tag size="small" :type="ch.status === 'published' ? 'success' : 'info'">
                {{ ch.status === 'published' ? '已发布' : '草稿' }}
              </el-tag>
              <span>{{ ch.minutes }} 分钟</span>
            </div>
          </div>
          <el-empty v-if="!course?.chapters?.length" description="暂无章节" />
        </el-card>
      </el-col>

      <el-col :span="16">
        <el-card v-if="selectedChapter" shadow="never">
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center">
              <span>编辑章节：{{ selectedChapter.title }}</span>
              <div>
                <el-button size="small" @click="openChapterDialog(selectedChapter)">设置</el-button>
                <el-popconfirm title="确定删除该章节？" @confirm="handleDeleteChapter">
                  <template #reference>
                    <el-button size="small" type="danger">删除</el-button>
                  </template>
                </el-popconfirm>
              </div>
            </div>
          </template>

          <el-form label-width="80px">
            <el-form-item label="标题">
              <el-input v-model="chapterForm.title" />
            </el-form-item>
            <el-form-item label="时长">
              <el-input-number v-model="chapterForm.minutes" :min="1" />
              <span style="margin-left: 8px; color: #9ca3af">分钟</span>
            </el-form-item>
            <el-form-item label="状态">
              <el-select v-model="chapterForm.status" style="width: 140px">
                <el-option label="草稿" value="draft" />
                <el-option label="已发布" value="published" />
              </el-select>
            </el-form-item>
            <el-form-item label="内容">
              <el-input v-model="chapterForm.content" type="textarea" :rows="16" placeholder="Markdown 格式" />
            </el-form-item>
            <el-form-item>
              <el-button type="primary" :loading="saving" @click="handleSaveChapter">保存章节</el-button>
            </el-form-item>
          </el-form>

          <el-divider>内容预览</el-divider>
          <div class="markdown-preview" v-html="previewHtml"></div>
        </el-card>
        <el-card v-else shadow="never">
          <el-empty description="请选择左侧章节进行编辑" />
        </el-card>
      </el-col>
    </el-row>

    <el-dialog v-model="chapterDialogVisible" :title="chapterEditing ? '编辑章节设置' : '新增章节'" width="480px">
      <el-form :model="chapterMeta" label-width="80px">
        <el-form-item v-if="!chapterEditing" label="章节 ID">
          <el-input v-model="chapterMeta.id" placeholder="英文标识" />
        </el-form-item>
        <el-form-item label="标题">
          <el-input v-model="chapterMeta.title" />
        </el-form-item>
        <el-form-item label="时长">
          <el-input-number v-model="chapterMeta.minutes" :min="1" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="chapterDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="handleCreateChapter">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ArrowLeft, Plus } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { marked } from 'marked'
import { coursesApi } from '../api'

const route = useRoute()
const loading = ref(false)
const saving = ref(false)
const course = ref(null)
const selectedChapter = ref(null)
const chapterForm = reactive({ title: '', minutes: 10, content: '', status: 'draft' })

const chapterDialogVisible = ref(false)
const chapterEditing = ref(false)
const chapterMeta = reactive({ id: '', title: '', minutes: 10 })

const previewHtml = computed(() => {
  try {
    return marked.parse(chapterForm.content || '')
  } catch {
    return '<p>预览解析失败</p>'
  }
})

async function loadCourse() {
  loading.value = true
  try {
    course.value = await coursesApi.get(route.params.id)
    if (course.value.chapters?.length) {
      selectChapter(course.value.chapters[0])
    }
  } finally {
    loading.value = false
  }
}

function selectChapter(ch) {
  selectedChapter.value = ch
  Object.assign(chapterForm, {
    title: ch.title,
    minutes: ch.minutes,
    content: ch.content || '',
    status: ch.status || 'draft'
  })
}

function openChapterDialog(ch) {
  chapterEditing.value = Boolean(ch)
  Object.assign(chapterMeta, { id: ch?.id || '', title: ch?.title || '', minutes: ch?.minutes || 10 })
  chapterDialogVisible.value = true
}

async function handleCreateChapter() {
  if (!chapterMeta.title) return ElMessage.warning('请输入章节标题')
  if (chapterEditing.value) {
    await coursesApi.updateChapter(route.params.id, chapterMeta.id, {
      title: chapterMeta.title,
      minutes: chapterMeta.minutes
    })
  } else {
    await coursesApi.addChapter(route.params.id, chapterMeta)
  }
  ElMessage.success('保存成功')
  chapterDialogVisible.value = false
  await loadCourse()
}

async function handleSaveChapter() {
  saving.value = true
  try {
    await coursesApi.updateChapter(route.params.id, selectedChapter.value.id, { ...chapterForm })
    ElMessage.success('章节已保存')
    await loadCourse()
    const updated = course.value.chapters.find((c) => c.id === selectedChapter.value.id)
    if (updated) selectChapter(updated)
  } finally {
    saving.value = false
  }
}

async function handleDeleteChapter() {
  await coursesApi.removeChapter(route.params.id, selectedChapter.value.id)
  ElMessage.success('章节已删除')
  selectedChapter.value = null
  await loadCourse()
}

onMounted(loadCourse)
</script>

<style scoped>
.chapter-item {
  padding: 12px;
  border-radius: 8px;
  cursor: pointer;
  margin-bottom: 4px;
  transition: background 0.15s;
}

.chapter-item:hover {
  background: #f3f4f6;
}

.chapter-item.active {
  background: #eef2ff;
  border-left: 3px solid #6366f1;
}

.chapter-item__title {
  font-weight: 500;
  font-size: 14px;
}

.chapter-item__meta {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-top: 4px;
  font-size: 12px;
  color: #9ca3af;
}
</style>
