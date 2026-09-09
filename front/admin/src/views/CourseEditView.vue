<template>
  <div v-loading="loading" class="page">
    <PageHeader
      back
      eyebrow="课程管理"
      :title="course?.title || '课程编辑'"
      :subtitle="course?.description || '编辑课程章节内容与发布状态'"
      @back="$router.push('/courses')"
    >
      <template #title-extra>
        <el-tag :type="course?.status === 'published' ? 'success' : 'info'" size="small">
          {{ course?.status === 'published' ? '已发布' : '草稿' }}
        </el-tag>
      </template>
      <el-button type="primary" :icon="Plus" @click="openChapterDialog()">新增章节</el-button>
    </PageHeader>

    <div class="editor">
      <aside class="panel chapters">
        <div class="panel__head">
          <h3 class="panel__title">
            <span class="panel__title-icon" :style="{ background: (course?.accent || '#6b5cff') + '1f', color: course?.accent || 'var(--primary)' }">{{ course?.icon || '📚' }}</span>
            章节列表
          </h3>
          <span class="muted num">{{ course?.chapters?.length || 0 }} 章</span>
        </div>
        <div class="panel__body chapters__list">
          <button
            v-for="(ch, idx) in course?.chapters || []"
            :key="ch.id"
            type="button"
            class="chapter"
            :class="{ 'chapter--active': selectedChapter?.id === ch.id }"
            @click="selectChapter(ch)"
          >
            <span class="chapter__index num">{{ String(idx + 1).padStart(2, '0') }}</span>
            <span class="chapter__main">
              <span class="chapter__title">{{ ch.title }}</span>
              <span class="chapter__meta">
                <span class="chapter__dot" :class="ch.status === 'published' ? 'chapter__dot--on' : ''" />
                {{ ch.status === 'published' ? '已发布' : '草稿' }} · {{ ch.minutes }} 分钟
              </span>
            </span>
            <el-icon :size="14" class="chapter__arrow"><ArrowRight /></el-icon>
          </button>
          <el-empty v-if="!course?.chapters?.length" description="暂无章节" :image-size="80" />
        </div>
      </aside>

      <section class="panel">
        <template v-if="selectedChapter">
          <div class="panel__head">
            <h3 class="panel__title">
              <span class="panel__title-icon"><el-icon><EditPen /></el-icon></span>
              <span class="ellipsis">{{ selectedChapter.title }}</span>
            </h3>
            <div class="toolbar">
              <el-button size="small" @click="openChapterDialog(selectedChapter)">设置</el-button>
              <el-popconfirm title="确定删除该章节？" @confirm="handleDeleteChapter">
                <template #reference>
                  <el-button size="small" type="danger" plain>删除</el-button>
                </template>
              </el-popconfirm>
            </div>
          </div>

          <div class="panel__body">
            <div class="meta-row">
              <div class="field field--grow">
                <label>标题</label>
                <el-input v-model="chapterForm.title" />
              </div>
              <div class="field">
                <label>时长（分钟）</label>
                <el-input-number v-model="chapterForm.minutes" :min="1" style="width: 140px" />
              </div>
              <div class="field">
                <label>状态</label>
                <el-select v-model="chapterForm.status" style="width: 130px">
                  <el-option label="草稿" value="draft" />
                  <el-option label="已发布" value="published" />
                </el-select>
              </div>
            </div>

            <div class="split">
              <div class="field">
                <label>内容（Markdown）</label>
                <el-input v-model="chapterForm.content" type="textarea" :rows="22" placeholder="使用 Markdown 编写章节内容…" class="md-input" />
              </div>
              <div class="field">
                <label>实时预览</label>
                <div class="markdown-preview preview" v-html="previewHtml"></div>
              </div>
            </div>
          </div>

          <div class="panel__foot foot">
            <span class="muted">最后修改会即时同步到学习端（已发布章节）</span>
            <el-button type="primary" :loading="saving" @click="handleSaveChapter">
              <el-icon class="el-icon--left"><Check /></el-icon>保存章节
            </el-button>
          </div>
        </template>
        <div v-else class="panel__body empty">
          <el-empty description="请选择左侧章节进行编辑" />
        </div>
      </section>
    </div>

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
import { Plus, ArrowRight, EditPen, Check } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { marked } from 'marked'
import { coursesApi } from '../api'
import PageHeader from '../components/common/PageHeader.vue'

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
.editor {
  display: grid;
  grid-template-columns: 320px minmax(0, 1fr);
  gap: 20px;
  align-items: start;
}

.chapters {
  position: sticky;
  top: calc(var(--header-height) + 20px);
}

.chapters__list {
  padding: 12px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  max-height: calc(100vh - 240px);
  overflow-y: auto;
}

.chapter {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 10px 12px;
  border: 1px solid transparent;
  border-radius: 12px;
  background: transparent;
  text-align: left;
  cursor: pointer;
  transition: all var(--t-fast);
}

.chapter:hover {
  background: var(--surface-2);
}

.chapter--active {
  background: var(--primary-soft);
  border-color: var(--primary-soft-2);
}

.chapter__index {
  width: 30px;
  height: 30px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 9px;
  background: var(--surface-3);
  color: var(--text-3);
  font-size: 12px;
  font-weight: 700;
  flex-shrink: 0;
}

.chapter--active .chapter__index {
  background: var(--primary);
  color: #fff;
}

.chapter__main {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  line-height: 1.3;
}

.chapter__title {
  font-size: 14px;
  font-weight: 600;
  color: var(--text);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.chapter__meta {
  display: inline-flex;
  align-items: center;
  gap: 6px;
  margin-top: 3px;
  font-size: 12px;
  color: var(--text-3);
}

.chapter__dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background: var(--border-strong);
}

.chapter__dot--on {
  background: var(--success);
}

.chapter__arrow {
  color: var(--text-3);
  opacity: 0;
  transition: opacity var(--t-fast);
}

.chapter--active .chapter__arrow,
.chapter:hover .chapter__arrow {
  opacity: 1;
}

.ellipsis {
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 420px;
}

.meta-row {
  display: flex;
  gap: 14px;
  flex-wrap: wrap;
  margin-bottom: 18px;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 0;
}

.field--grow {
  flex: 1;
  min-width: 200px;
}

.field label {
  font-size: 12.5px;
  font-weight: 600;
  color: var(--text-3);
}

.split {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 16px;
}

.md-input :deep(.el-textarea__inner) {
  font-family: 'JetBrains Mono', 'Fira Code', ui-monospace, Menlo, monospace;
  font-size: 13px;
  line-height: 1.7;
  border-radius: 12px !important;
}

.preview {
  height: 100%;
  min-height: 300px;
  max-height: 520px;
}

.foot {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

.empty {
  min-height: 420px;
  display: flex;
  align-items: center;
  justify-content: center;
}

@media (max-width: 1100px) {
  .editor {
    grid-template-columns: 1fr;
  }

  .chapters {
    position: static;
  }

  .split {
    grid-template-columns: 1fr;
  }
}
</style>
