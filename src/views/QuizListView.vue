<script setup>
import { quizzes } from '../data/quizzes'
import { getCourse } from '../data/courses'
import { useLearningStore } from '../stores/learning'

const learning = useLearningStore()

function courseOf(quiz) {
  return getCourse(quiz.courseId)
}

function resultOf(quiz) {
  return learning.quizResults[quiz.id] || null
}

function percent(result) {
  return Math.round((result.score / result.total) * 100)
}
</script>

<template>
  <div class="page">
    <header class="page-header">
      <h1>知识测验</h1>
      <p>每门课程配套一套测验题，检验学习效果并附详细解析。</p>
    </header>

    <div class="quiz-grid">
      <div v-for="quiz in quizzes" :key="quiz.id" class="quiz-card card fade-up">
        <div class="quiz-card__top">
          <span
            class="quiz-card__icon"
            :style="{ background: courseOf(quiz).accent + '1a' }"
          >
            {{ courseOf(quiz).icon }}
          </span>
          <span v-if="resultOf(quiz)" class="quiz-card__score" :class="{
            'quiz-card__score--good': percent(resultOf(quiz)) >= 80,
            'quiz-card__score--mid': percent(resultOf(quiz)) >= 60 && percent(resultOf(quiz)) < 80,
            'quiz-card__score--bad': percent(resultOf(quiz)) < 60
          }">
            {{ percent(resultOf(quiz)) }} 分
          </span>
          <span v-else class="quiz-card__score quiz-card__score--none">未测验</span>
        </div>
        <h3>{{ quiz.title }}</h3>
        <p>{{ quiz.description }}</p>
        <div class="quiz-card__meta">
          <span>{{ quiz.questions.length }} 道题</span>
          <span v-if="resultOf(quiz)">
            上次：{{ resultOf(quiz).score }}/{{ resultOf(quiz).total }} 正确
          </span>
        </div>
        <router-link :to="`/quiz/${quiz.id}`" class="btn btn--primary">
          {{ resultOf(quiz) ? '重新测验' : '开始测验' }}
        </router-link>
      </div>
    </div>
  </div>
</template>

<style scoped>
.quiz-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(290px, 1fr));
  gap: 18px;
}

.quiz-card {
  display: flex;
  flex-direction: column;
  padding: 22px;
}

.quiz-card__top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 14px;
}

.quiz-card__icon {
  width: 46px;
  height: 46px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  border-radius: 12px;
}

.quiz-card__score {
  font-size: 13px;
  font-weight: 700;
  padding: 3px 12px;
  border-radius: 999px;
}

.quiz-card__score--good { background: var(--success-soft); color: #047857; }
.quiz-card__score--mid { background: #fffbeb; color: #b45309; }
.quiz-card__score--bad { background: #fef2f2; color: var(--danger); }
.quiz-card__score--none { background: var(--surface-2); color: var(--text-3); }

.quiz-card h3 {
  margin: 0 0 8px;
  font-size: 16.5px;
}

.quiz-card p {
  margin: 0 0 14px;
  font-size: 13.5px;
  color: var(--text-2);
  flex: 1;
}

.quiz-card__meta {
  display: flex;
  justify-content: space-between;
  font-size: 12.5px;
  color: var(--text-3);
  margin-bottom: 14px;
}
</style>
