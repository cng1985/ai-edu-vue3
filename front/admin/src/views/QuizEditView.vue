<template>
  <div v-loading="loading" class="page">
    <PageHeader
      back
      eyebrow="题库管理"
      :title="quiz?.title || '测验编辑'"
      :subtitle="quiz?.description || '维护该测验下的题目、选项与解析'"
      @back="$router.push('/quizzes')"
    >
      <template #title-extra>
        <el-tag :type="quizStatus === 'published' ? 'success' : 'info'" size="small">
          {{ quizStatus === 'published' ? '已发布' : '草稿' }}
        </el-tag>
      </template>
      <el-select v-model="quizStatus" style="width: 130px" @change="handleStatusChange">
        <el-option label="草稿" value="draft" />
        <el-option label="已发布" value="published" />
      </el-select>
      <el-button type="primary" :icon="Plus" @click="openQuestionDialog()">新增题目</el-button>
    </PageHeader>

    <div class="panel">
      <div class="panel__head">
        <h3 class="panel__title">
          <span class="panel__title-icon"><el-icon><EditPen /></el-icon></span>
          题目列表
        </h3>
        <span class="muted num">{{ quiz?.questions?.length || 0 }} 题</span>
      </div>
      <div class="panel__body questions">
        <div
          v-for="(q, idx) in quiz?.questions || []"
          :key="idx"
          class="question"
        >
          <span class="question__index num">{{ idx + 1 }}</span>
          <div class="question__main">
            <div class="question__text">{{ q.text }}</div>
            <div class="question__options">
              <span
                v-for="(opt, i) in q.options"
                :key="i"
                class="option"
                :class="{ 'option--correct': i === q.answer }"
              >
                <b>{{ String.fromCharCode(65 + i) }}</b>{{ opt }}
              </span>
            </div>
            <div v-if="q.explanation" class="question__explain">
              <el-icon :size="14"><InfoFilled /></el-icon>{{ q.explanation }}
            </div>
          </div>
          <div class="question__actions">
            <el-button link type="primary" @click="openQuestionDialog(q, idx)">编辑</el-button>
            <el-button link type="danger" @click="handleDeleteQuestion(idx)">删除</el-button>
          </div>
        </div>
        <el-empty v-if="!quiz?.questions?.length" description="暂无题目，点击右上角新增" :image-size="90" />
      </div>
    </div>

    <el-dialog v-model="dialogVisible" :title="editingIndex >= 0 ? '编辑题目' : '新增题目'" width="680px">
      <el-form :model="qForm" label-position="top">
        <el-form-item label="题干">
          <el-input v-model="qForm.text" type="textarea" :rows="3" placeholder="请输入题目内容" />
        </el-form-item>
        <el-form-item label="选项（点击左侧圆点标记正确答案）">
          <div class="option-editor">
            <div v-for="(opt, i) in qForm.options" :key="i" class="option-row" :class="{ 'option-row--correct': qForm.answer === i }">
              <button type="button" class="option-row__pick" :title="qForm.answer === i ? '正确答案' : '设为正确答案'" @click="qForm.answer = i">
                <el-icon v-if="qForm.answer === i" :size="14"><Check /></el-icon>
                <span v-else>{{ String.fromCharCode(65 + i) }}</span>
              </button>
              <el-input v-model="qForm.options[i]" :placeholder="`选项 ${String.fromCharCode(65 + i)}`" />
              <el-button v-if="qForm.options.length > 2" :icon="Delete" text @click="removeOption(i)" />
            </div>
            <el-button v-if="qForm.options.length < 6" text type="primary" :icon="Plus" @click="qForm.options.push('')">添加选项</el-button>
          </div>
        </el-form-item>
        <el-form-item label="解析（可选）">
          <el-input v-model="qForm.explanation" type="textarea" :rows="2" placeholder="答题后展示给学员的解析" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSaveQuestion">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { Plus, Delete, EditPen, InfoFilled, Check } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { quizzesApi } from '../api'
import PageHeader from '../components/common/PageHeader.vue'

const route = useRoute()
const loading = ref(false)
const saving = ref(false)
const quiz = ref(null)
const quizStatus = ref('draft')
const dialogVisible = ref(false)
const editingIndex = ref(-1)
const qForm = reactive({ text: '', options: ['', '', '', ''], answer: 0, explanation: '' })

async function loadQuiz() {
  loading.value = true
  try {
    quiz.value = await quizzesApi.get(route.params.id)
    quizStatus.value = quiz.value.status
  } finally {
    loading.value = false
  }
}

function openQuestionDialog(row, index) {
  editingIndex.value = index ?? -1
  if (row) {
    Object.assign(qForm, {
      text: row.text,
      options: [...row.options],
      answer: row.answer,
      explanation: row.explanation || ''
    })
  } else {
    Object.assign(qForm, { text: '', options: ['', '', '', ''], answer: 0, explanation: '' })
  }
  dialogVisible.value = true
}

function removeOption(i) {
  qForm.options.splice(i, 1)
  if (qForm.answer >= qForm.options.length) qForm.answer = qForm.options.length - 1
  else if (qForm.answer > i) qForm.answer -= 1
}

async function handleSaveQuestion() {
  if (!qForm.text || qForm.options.some((o) => !o.trim())) {
    return ElMessage.warning('请填写完整题目和选项')
  }
  saving.value = true
  try {
    const questions = [...(quiz.value.questions || [])]
    const item = { text: qForm.text, options: [...qForm.options], answer: qForm.answer, explanation: qForm.explanation }
    if (editingIndex.value >= 0) {
      questions[editingIndex.value] = item
    } else {
      questions.push(item)
    }
    await quizzesApi.update(route.params.id, { questions })
    ElMessage.success('保存成功')
    dialogVisible.value = false
    await loadQuiz()
  } finally {
    saving.value = false
  }
}

async function handleDeleteQuestion(index) {
  const questions = [...quiz.value.questions]
  questions.splice(index, 1)
  await quizzesApi.update(route.params.id, { questions })
  ElMessage.success('已删除')
  await loadQuiz()
}

async function handleStatusChange(status) {
  await quizzesApi.update(route.params.id, { status })
  ElMessage.success('状态已更新')
}

onMounted(loadQuiz)
</script>

<style scoped>
.questions {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.question {
  display: flex;
  gap: 14px;
  padding: 16px 18px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  background: var(--surface);
  transition: border-color var(--t-fast), box-shadow var(--t-fast);
}

.question:hover {
  border-color: var(--primary-soft-2);
  box-shadow: var(--shadow-sm);
}

.question__index {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  flex-shrink: 0;
  border-radius: 9px;
  background: var(--primary-soft);
  color: var(--primary-deep);
  font-weight: 700;
  font-size: 13px;
}

.question__main {
  flex: 1;
  min-width: 0;
}

.question__text {
  font-size: 14.5px;
  font-weight: 600;
  color: var(--text);
  line-height: 1.6;
}

.question__options {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 6px 12px;
  margin-top: 10px;
}

.option {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 6px 10px;
  border-radius: 9px;
  background: var(--surface-2);
  font-size: 13px;
  color: var(--text-2);
}

.option b {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 20px;
  height: 20px;
  border-radius: 6px;
  background: var(--surface);
  color: var(--text-3);
  font-size: 11px;
  flex-shrink: 0;
}

.option--correct {
  background: var(--success-soft);
  color: var(--success-strong);
  font-weight: 600;
}

.option--correct b {
  background: var(--success);
  color: #fff;
}

.question__explain {
  display: flex;
  align-items: flex-start;
  gap: 6px;
  margin-top: 10px;
  font-size: 12.5px;
  color: var(--text-3);
  line-height: 1.6;
}

.question__explain .el-icon {
  margin-top: 3px;
  color: var(--primary);
}

.question__actions {
  display: flex;
  align-items: flex-start;
  gap: 4px;
  flex-shrink: 0;
}

.option-editor {
  display: flex;
  flex-direction: column;
  gap: 8px;
  width: 100%;
}

.option-row {
  display: flex;
  align-items: center;
  gap: 8px;
}

.option-row__pick {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 34px;
  height: 34px;
  flex-shrink: 0;
  border: 1px solid var(--border-strong);
  border-radius: 10px;
  background: var(--surface);
  color: var(--text-3);
  font-size: 12px;
  font-weight: 700;
  cursor: pointer;
  transition: all var(--t-fast);
}

.option-row__pick:hover {
  border-color: var(--success);
  color: var(--success);
}

.option-row--correct .option-row__pick {
  background: var(--success);
  border-color: var(--success);
  color: #fff;
}

@media (max-width: 760px) {
  .question__options {
    grid-template-columns: 1fr;
  }
}
</style>
