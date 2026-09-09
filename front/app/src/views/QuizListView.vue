<script setup>
import { computed } from 'vue'
import { quizzes } from '../data/quizzes'
import { getCourse } from '../data/courses'
import { useLearningStore } from '../stores/learning'
import Icon from '../components/Icon.vue'

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

function scoreTone(p) {
  if (p >= 80) return 'success'
  if (p >= 60) return 'warning'
  return 'danger'
}

const finishedCount = computed(() => quizzes.filter((q) => resultOf(q)).length)
const averageScore = computed(() => {
  const done = quizzes.map(resultOf).filter(Boolean)
  if (!done.length) return 0
  return Math.round(done.reduce((sum, r) => sum + percent(r), 0) / done.length)
})
</script>

<template>
  <div class="page">
    <header class="page-header page-header--split">
      <div>
        <span class="eyebrow">检验学习效果</span>
        <h1>知识测验</h1>
        <p>每门课程配套一套测验题，提交后逐题判分并展示解析。</p>
      </div>
      <div class="summary card">
        <div><strong class="num">{{ finishedCount }}/{{ quizzes.length }}</strong><span>已完成</span></div>
        <div><strong class="num">{{ averageScore }}</strong><span>平均分</span></div>
      </div>
    </header>

    <div class="quiz-grid stagger">
      <div v-for="quiz in quizzes" :key="quiz.id" class="quiz-card card card--hover" :style="{ '--accent': courseOf(quiz).accent }">
        <div class="quiz-card__top">
          <span class="quiz-card__icon">{{ courseOf(quiz).icon }}</span>
          <span v-if="resultOf(quiz)" class="quiz-card__score" :class="`quiz-card__score--${scoreTone(percent(resultOf(quiz)))}`">
            <strong class="num">{{ percent(resultOf(quiz)) }}</strong><small>分</small>
          </span>
          <span v-else class="tag tag--neutral">未测验</span>
        </div>
        <h3>{{ quiz.title }}</h3>
        <p>{{ quiz.description }}</p>
        <div class="quiz-card__meta">
          <span><Icon name="clipboard" :size="13" /> {{ quiz.questions.length }} 道题</span>
          <span v-if="resultOf(quiz)"><Icon name="check" :size="13" /> 上次 {{ resultOf(quiz).score }}/{{ resultOf(quiz).total }} 正确</span>
          <span v-else><Icon name="clock" :size="13" /> 约 {{ quiz.questions.length * 1.5 }} 分钟</span>
        </div>
        <router-link :to="`/quiz/${quiz.id}`" class="btn" :class="resultOf(quiz) ? 'btn--ghost' : 'btn--primary'">
          <Icon :name="resultOf(quiz) ? 'refresh' : 'play'" :size="15" />
          {{ resultOf(quiz) ? '重新测验' : '开始测验' }}
        </router-link>
      </div>
    </div>
  </div>
</template>

<style scoped>
.summary {
  display: flex;
  gap: 4px;
  padding: 6px;
}

.summary div {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 90px;
  padding: 8px 14px;
}

.summary div + div {
  border-left: 1px solid var(--border);
}

.summary strong {
  font-size: 20px;
  line-height: 1.1;
}

.summary span {
  color: var(--text-3);
  font-size: 11.5px;
}

.quiz-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

.quiz-card {
  position: relative;
  display: flex;
  flex-direction: column;
  padding: 24px;
  overflow: hidden;
}

.quiz-card::before {
  content: '';
  position: absolute;
  inset: 0 0 auto 0;
  height: 4px;
  background: linear-gradient(90deg, var(--accent), color-mix(in srgb, var(--accent) 30%, #fff));
}

.quiz-card__top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
}

.quiz-card__icon {
  width: 52px;
  height: 52px;
  display: grid;
  place-items: center;
  font-size: 26px;
  border-radius: 16px;
  background: color-mix(in srgb, var(--accent) 12%, #fff);
  box-shadow: inset 0 0 0 1px color-mix(in srgb, var(--accent) 22%, transparent);
}

.quiz-card__score {
  display: inline-flex;
  align-items: baseline;
  gap: 2px;
  padding: 4px 12px;
  border-radius: 999px;
  font-weight: 700;
}

.quiz-card__score strong { font-size: 17px; }
.quiz-card__score small { font-size: 11px; }
.quiz-card__score--success { background: var(--success-soft); color: var(--success-strong); }
.quiz-card__score--warning { background: var(--warning-soft); color: var(--warning-strong); }
.quiz-card__score--danger { background: var(--danger-soft); color: var(--danger-strong); }

.quiz-card h3 {
  margin: 0 0 8px;
  font-size: 17px;
}

.quiz-card p {
  margin: 0 0 16px;
  font-size: 13.5px;
  color: var(--text-2);
  line-height: 1.65;
  flex: 1;
}

.quiz-card__meta {
  display: flex;
  gap: 14px;
  font-size: 12.5px;
  color: var(--text-3);
  margin-bottom: 16px;
}

.quiz-card__meta span {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}
</style>
