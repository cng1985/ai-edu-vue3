<script setup>
import { ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getQuiz } from '../data/quizzes'
import { getCourse } from '../data/courses'
import { useLearningStore } from '../stores/learning'
import ProgressRing from '../components/ProgressRing.vue'
import Icon from '../components/Icon.vue'

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

const scoreTone = computed(() => (scorePercent.value >= 80 ? 'success' : scorePercent.value >= 60 ? 'warning' : 'danger'))
const scoreColor = computed(() => ({ success: '#0fb981', warning: '#f59e0b', danger: '#ef4d63' }[scoreTone.value]))

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
  <div class="page quiz-page" v-if="quiz">
    <router-link to="/quiz" class="back-link"><Icon name="arrowLeft" :size="15" /> 返回测验列表</router-link>

    <header class="page-header">
      <div class="quiz-tags">
        <span v-if="course" class="tag tag--neutral">{{ course.icon }} {{ course.title }}</span>
        <span class="tag">{{ quiz.questions.length }} 道题</span>
      </div>
      <h1>{{ quiz.title }}</h1>
      <p>{{ quiz.description }}</p>
    </header>

    <section v-if="submitted" class="result card fade-up">
      <ProgressRing :percent="scorePercent" :size="128" :stroke="11" :color="scoreColor" />
      <div class="result__info">
        <span class="tag" :class="`tag--${scoreTone}`">{{ scorePercent >= 80 ? '优秀' : scorePercent >= 60 ? '及格' : '需加强' }}</span>
        <h2>
          {{ scorePercent >= 80 ? '太棒了，掌握得很扎实！' : scorePercent >= 60 ? '不错，再接再厉' : '别灰心，回顾后再来一次' }}
        </h2>
        <p>
          答对 <strong class="num">{{ score }}</strong> / {{ quiz.questions.length }} 题。
          {{ scorePercent < 80 && course ? '建议回顾课程后重新测验。' : '可以继续挑战其他课程的测验。' }}
        </p>
        <div class="result__actions">
          <button class="btn btn--primary" @click="retry"><Icon name="refresh" :size="15" /> 重新测验</button>
          <router-link v-if="course" :to="`/courses/${course.id}`" class="btn btn--ghost">
            <Icon name="book" :size="15" /> 回顾课程
          </router-link>
        </div>
      </div>
    </section>

    <div v-else class="progress-bar card">
      <span class="progress-bar__label">已作答 <b class="num">{{ answeredCount }}</b> / {{ quiz.questions.length }}</span>
      <div class="progress progress--thin">
        <div :style="{ width: (answeredCount / quiz.questions.length) * 100 + '%' }"></div>
      </div>
      <div class="progress-bar__dots">
        <i v-for="(a, i) in answers" :key="i" :class="{ on: a !== null }"></i>
      </div>
    </div>

    <section
      v-for="(question, qi) in quiz.questions"
      :key="qi"
      class="quiz-question card"
      :class="{ 'quiz-question--right': submitted && answers[qi] === question.answer, 'quiz-question--miss': submitted && answers[qi] !== question.answer }"
    >
      <h3>
        <span class="quiz-question__no num">{{ qi + 1 }}</span>
        <span>{{ question.text }}</span>
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
          <span class="quiz-option__text">{{ opt }}</span>
          <Icon v-if="submitted && oi === question.answer" name="check" :size="16" :stroke="3" class="quiz-option__mark" />
          <Icon
            v-else-if="submitted && answers[qi] === oi && oi !== question.answer"
            name="x"
            :size="16"
            :stroke="3"
            class="quiz-option__mark quiz-option__mark--wrong"
          />
        </button>
      </div>
      <div v-if="submitted" class="quiz-question__explanation">
        <Icon name="sparkles" :size="15" />
        <div><strong>解析</strong>{{ question.explanation }}</div>
      </div>
    </section>

    <div v-if="!submitted" class="submit-bar">
      <button class="btn btn--primary btn--lg submit-bar__btn" :disabled="!allAnswered" @click="submit">
        {{ allAnswered ? '提交答卷' : `还有 ${quiz.questions.length - answeredCount} 题未作答` }}
        <Icon v-if="allAnswered" name="arrowRight" :size="16" />
      </button>
    </div>
  </div>

  <div class="page" v-else>
    <div class="empty-state card">
      <div class="empty-state__icon"><Icon name="alert" :size="30" /></div>
      <h2>未找到该测验</h2>
      <router-link to="/quiz" class="btn btn--primary">返回测验列表</router-link>
    </div>
  </div>
</template>

<style scoped>
.quiz-page { max-width: 880px; }

.quiz-tags { display: flex; gap: 6px; margin-bottom: 12px; }

.progress-bar {
  position: sticky;
  top: 16px;
  z-index: 5;
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 14px 20px;
  margin-bottom: 20px;
  font-size: 13.5px;
  color: var(--text-2);
  white-space: nowrap;
  backdrop-filter: blur(10px);
  background: rgba(255, 255, 255, 0.9);
}

.progress-bar__label b { color: var(--text); }
.progress-bar .progress { flex: 1; }
.progress-bar__dots { display: flex; gap: 4px; }
.progress-bar__dots i { width: 8px; height: 8px; border-radius: 50%; background: var(--surface-3); transition: background var(--t-fast); }
.progress-bar__dots i.on { background: var(--primary); }

.result {
  display: flex;
  align-items: center;
  gap: 30px;
  padding: 32px 36px;
  margin-bottom: 24px;
  border-radius: var(--radius-lg);
  background: linear-gradient(120deg, #f6f3ff, #fff 60%);
}

.result__info h2 { margin: 10px 0 6px; font-size: 22px; }
.result__info p { margin: 0 0 16px; color: var(--text-2); font-size: 14.5px; }
.result__info p strong { color: var(--text); font-size: 16px; }
.result__actions { display: flex; gap: 10px; flex-wrap: wrap; }

.quiz-question {
  padding: 26px 30px;
  margin-bottom: 16px;
  border-left: 4px solid transparent;
  transition: border-color var(--t);
}
.quiz-question--right { border-left-color: var(--success); }
.quiz-question--miss { border-left-color: var(--danger); }

.quiz-question h3 {
  display: flex;
  gap: 14px;
  margin: 0 0 18px;
  font-size: 16px;
  line-height: 1.6;
}

.quiz-question__no {
  min-width: 30px;
  height: 30px;
  display: grid;
  place-items: center;
  background: var(--primary-soft);
  color: var(--primary-strong);
  border-radius: 9px;
  font-size: 13px;
  font-weight: 700;
}

.quiz-question__options {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.quiz-option {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 13px 16px;
  border: 1.5px solid var(--border);
  border-radius: 13px;
  background: var(--surface);
  font: inherit;
  font-size: 14.5px;
  text-align: left;
  cursor: pointer;
  transition: all var(--t-fast);
  color: var(--text);
}

.quiz-option:hover { border-color: var(--primary); background: var(--primary-soft); }

.quiz-option__letter {
  min-width: 28px;
  height: 28px;
  display: grid;
  place-items: center;
  background: var(--surface-3);
  border-radius: 8px;
  font-size: 12.5px;
  font-weight: 700;
  color: var(--text-2);
}

.quiz-option__text { flex: 1; }

.quiz-option--chosen { border-color: var(--primary); background: var(--primary-soft); font-weight: 600; }
.quiz-option--chosen .quiz-option__letter { background: var(--primary); color: #fff; }
.quiz-option--correct { border-color: var(--success); background: var(--success-soft); font-weight: 600; }
.quiz-option--correct .quiz-option__letter { background: var(--success); color: #fff; }
.quiz-option--wrong { border-color: var(--danger); background: var(--danger-soft); }
.quiz-option--wrong .quiz-option__letter { background: var(--danger); color: #fff; }

.quiz-option__mark { color: var(--success); }
.quiz-option__mark--wrong { color: var(--danger); }

.quiz-question__explanation {
  display: flex;
  gap: 10px;
  margin-top: 16px;
  padding: 14px 16px;
  background: var(--surface-2);
  border-radius: 13px;
  font-size: 13.5px;
  color: var(--text-2);
  line-height: 1.7;
}

.quiz-question__explanation .icon { flex-shrink: 0; margin-top: 4px; color: var(--primary); }
.quiz-question__explanation strong { display: block; margin-bottom: 2px; color: var(--text); font-size: 12.5px; letter-spacing: 0.04em; }

.submit-bar {
  position: sticky;
  bottom: 20px;
  text-align: center;
  margin-top: 24px;
}

.submit-bar__btn {
  padding: 14px 44px;
  box-shadow: var(--shadow-lg);
}

@media (max-width: 720px) {
  .result { flex-direction: column; align-items: flex-start; }
  .progress-bar__dots { display: none; }
}
</style>
