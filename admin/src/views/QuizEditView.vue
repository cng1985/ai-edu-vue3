<template>
  <div v-loading="loading">
    <div class="page-header">
      <div style="display: flex; align-items: center; gap: 12px">
        <el-button :icon="ArrowLeft" @click="$router.push('/quizzes')">返回</el-button>
        <h2>{{ quiz?.title || '测验编辑' }}</h2>
        <el-select v-model="quizStatus" style="width: 120px" @change="handleStatusChange">
          <el-option label="草稿" value="draft" />
          <el-option label="已发布" value="published" />
        </el-select>
      </div>
      <el-button type="primary" :icon="Plus" @click="openQuestionDialog()">新增题目</el-button>
    </div>

    <el-card shadow="never">
      <el-table :data="quiz?.questions || []" stripe>
        <el-table-column type="index" label="#" width="50" />
        <el-table-column prop="text" label="题目" min-width="280" show-overflow-tooltip />
        <el-table-column label="选项数" width="80">
          <template #default="{ row }">{{ row.options?.length || 0 }}</template>
        </el-table-column>
        <el-table-column label="正确答案" width="100">
          <template #default="{ row }">{{ String.fromCharCode(65 + row.answer) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="140" fixed="right">
          <template #default="{ row, $index }">
            <el-button link type="primary" @click="openQuestionDialog(row, $index)">编辑</el-button>
            <el-button link type="danger" @click="handleDeleteQuestion($index)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="dialogVisible" :title="editingIndex >= 0 ? '编辑题目' : '新增题目'" width="640px">
      <el-form :model="qForm" label-width="80px">
        <el-form-item label="题干">
          <el-input v-model="qForm.text" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item v-for="(opt, i) in qForm.options" :key="i" :label="`选项 ${String.fromCharCode(65 + i)}`">
          <div style="display: flex; gap: 8px; width: 100%">
            <el-radio v-model="qForm.answer" :value="i" />
            <el-input v-model="qForm.options[i]" style="flex: 1" />
            <el-button v-if="qForm.options.length > 2" :icon="Delete" circle @click="qForm.options.splice(i, 1)" />
          </div>
        </el-form-item>
        <el-form-item>
          <el-button v-if="qForm.options.length < 6" @click="qForm.options.push('')">添加选项</el-button>
        </el-form-item>
        <el-form-item label="解析">
          <el-input v-model="qForm.explanation" type="textarea" :rows="2" />
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
import { ArrowLeft, Plus, Delete } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { quizzesApi } from '../api'

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
