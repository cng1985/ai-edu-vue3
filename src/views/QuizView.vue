<script setup>
import { ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getQuiz } from '../data/quizzes'
import { getCourse } from '../data/courses'
import { useLearningStore } from '../stores/learning'

const route = useRoute()
const learning = useLearningStore()

const quiz = computed(() => getQuiz(route.params.quizId))
const course = computed(() => (quiz.value ? getCourse(quiz.value.courseId) : null))

// answers[i] = 选项下标 或 null
const answers = ref([])
const submitted = ref(false)

watch(
  quiz,
  (q) => {
    answers.value = q ? q.questions.map(() => null) : []
    submitted.value = false
  },
  { immediate: true }
)

const answeredCount = computed(() => answers.value.filter((a) => a !== null).length)
const allAnswered = computed(
  () => quiz.value && answeredCount.value === quiz.value.questions.length
)

const score = computed(() => {
  if (!quiz.value) return 0
  return quiz.value.questions.reduce(
    (sum, q, i) => sum + (answers.value[i] === q.answer ? 1 : 0),
    0
  )
})

const scorePercent = computed(() =>
  quiz.value ? Math.round((score.value / quiz.value.questions.length) * 100) : 0
)

function choose(qIndex, optIndex) {
  if (submitted.value) return
  answers.value[qIndex] = optIndex
}

function submit() {
  if (!allAnswered.value) return
  submitted.value = true
  learning.saveQuizResult(
    quiz.value.id,
    score.value,
    quiz.value.questions.length,
    [...answers.value]
  )
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function retry() {
  answers.value = quiz.value.questions.map(() => null)
  submitted.value = false
  window.scrollTo({ top: 0 })
}

function optionClass(qIndex, optIndex) {
  const question = quiz.value.questions[qIndex]
  const chosen = answers.value[qIndex] === optIndex
  if (!submitted.value) {
    return { 'quiz-option--chosen': chosen }
  }
  return {
    'quiz-option--correct': optIndex === question.answer,
    'quiz-option--wrong': chosen && optIndex !== question.answer
  }
}
</script>

<template>
  <div class="page" v-if="quiz">
    <router-link to="/quiz" class="back">← 返回测验列表</router-link>

    <header class="page-header">
      <h1>{{ quiz.title }}</h1>
      <p>{{ quiz.description }}</p>
    </header>

    <section v-if="submitted" class="result card fade-up">
      <div class="result__score" :class="{
        'result__score--good': scorePercent >= 80,
        'result__score--mid': scorePercent >= 60 && scorePercent < 80,
        'result__score--bad': scorePercent < 60
      }">
        {{ scorePercent }}
        <span>分</span>
      </div>
      <div class="result__info">
        <h2>
          {{ scorePercent >= 80 ? '🎉 优秀！' : scorePercent >= 60 ? '👍 及格，再接再厉' : '💪 继续努力' }}
        </h2>
        <p>
          答对 {{ score }} / {{ quiz.questions.length }} 题。
          {{ scorePercent < 80 && course ? '建议回顾课程后重新测验。' : '' }}
        </p>
        <div class="result__actions">
          <button class="btn btn--primary" @click="retry">重新测验</button>
          <router-link v-if="course" :to="`/courses/${course.id}`" class="btn btn--ghost">
            回顾课程
          </router-link>
        </div>
      </div>
    </section>

    <div v-else class="progress-bar card">
      <span>已作答 {{ answeredCount }} / {{ quiz.questions.length }}</span>
      <div class="progress-bar__track">
        <div
          class="progress-bar__fill"
          :style="{ width: (answeredCount / quiz.questions.length) * 100 + '%' }"
        ></div>
      </div>
    </div>

    <section
      v-for="(question, qi) in quiz.questions"
      :key="qi"
      class="quiz-question card"
    >
      <h3>
        <span class="quiz-question__no">{{ qi + 1 }}</span>
        {{ question.text }}
      </h3>
      <div class="quiz-question__options">
        <button
          v-for="(opt, oi) in question.options"
          :key="oi"
          class="quiz-option"
          :class="optionClass(qi, oi)"
          @click="choose(qi, oi)"
        >
          <span class="quiz-option__letter">{{ 'ABCD'[oi] }}</span>
          <span>{{ opt }}</span>
          <span v-if="submitted && oi === question.answer" class="quiz-option__mark">✓</span>
          <span
            v-else-if="submitted && answers[qi] === oi && oi !== question.answer"
            class="quiz-option__mark quiz-option__mark--wrong"
          >✗</span>
        </button>
      </div>
      <div v-if="submitted" class="quiz-question__explanation">
        <strong>解析：</strong>{{ question.explanation }}
      </div>
    </section>

    <div v-if="!submitted" class="submit-bar">
      <button class="btn btn--primary submit-bar__btn" :disabled="!allAnswered" @click="submit">
        {{ allAnswered ? '提交答卷' : `还有 ${quiz.questions.length - answeredCount} 题未作答` }}
      </button>
    </div>
  </div>

  <div class="page" v-else>
    <div class="card" style="padding: 48px; text-align: center">
      <p>未找到该测验。</p>
      <router-link to="/quiz" class="btn btn--primary">返回测验列表</router-link>
    </div>
  </div>
</template>

<style scoped>
.back {
  display: inline-block;
  margin-bottom: 16px;
  font-size: 14px;
  color: var(--text-2);
}

.back:hover {
  color: var(--primary);
}

.progress-bar {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 14px 20px;
  margin-bottom: 20px;
  font-size: 13.5px;
  color: var(--text-2);
  white-space: nowrap;
}

.progress-bar__track {
  flex: 1;
  height: 8px;
  background: var(--border);
  border-radius: 999px;
  overflow: hidden;
}

.progress-bar__fill {
  height: 100%;
  background: var(--primary);
  border-radius: 999px;
  transition: width 0.3s ease;
}

.result {
  display: flex;
  align-items: center;
  gap: 28px;
  padding: 30px 34px;
  margin-bottom: 24px;
}

.result__score {
  font-size: 46px;
  font-weight: 800;
  line-height: 1;
}

.result__score span {
  font-size: 17px;
  font-weight: 600;
}

.result__score--good { color: var(--success); }
.result__score--mid { color: var(--warning); }
.result__score--bad { color: var(--danger); }

.result__info h2 {
  margin: 0 0 6px;
  font-size: 19px;
}

.result__info p {
  margin: 0 0 14px;
  color: var(--text-2);
  font-size: 14px;
}

.result__actions {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
}

.quiz-question {
  padding: 24px 28px;
  margin-bottom: 16px;
}

.quiz-question h3 {
  display: flex;
  gap: 12px;
  margin: 0 0 16px;
  font-size: 15.5px;
  line-height: 1.6;
}

.quiz-question__no {
  min-width: 26px;
  height: 26px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--primary-soft);
  color: var(--primary-strong);
  border-radius: 8px;
  font-size: 13px;
  font-weight: 700;
}

.quiz-question__options {
  display: flex;
  flex-direction: column;
  gap: 9px;
}

.quiz-option {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border: 1.5px solid var(--border);
  border-radius: var(--radius-sm);
  background: var(--surface);
  font-size: 14px;
  text-align: left;
  cursor: pointer;
  transition: all 0.15s ease;
  color: var(--text);
}

.quiz-option:hover {
  border-color: var(--primary);
}

.quiz-option--chosen {
  border-color: var(--primary);
  background: var(--primary-soft);
  font-weight: 600;
}

.quiz-option--correct {
  border-color: var(--success);
  background: var(--success-soft);
  font-weight: 600;
}

.quiz-option--wrong {
  border-color: var(--danger);
  background: #fef2f2;
}

.quiz-option__letter {
  min-width: 24px;
  height: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--surface-2);
  border-radius: 7px;
  font-size: 12px;
  font-weight: 700;
  color: var(--text-2);
}

.quiz-option__mark {
  margin-left: auto;
  font-weight: 800;
  color: var(--success);
}

.quiz-option__mark--wrong {
  color: var(--danger);
}

.quiz-question__explanation {
  margin-top: 14px;
  padding: 12px 16px;
  background: var(--surface-2);
  border-radius: var(--radius-sm);
  font-size: 13.5px;
  color: var(--text-2);
  line-height: 1.7;
}

.submit-bar {
  position: sticky;
  bottom: 20px;
  text-align: center;
  margin-top: 24px;
}

.submit-bar__btn {
  padding: 13px 42px;
  font-size: 15px;
  box-shadow: var(--shadow);
}
</style>
