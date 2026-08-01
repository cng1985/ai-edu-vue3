<script setup>
import { computed } from 'vue'
import { courses, totalChapterCount } from '../data/courses'
import { quizzes } from '../data/quizzes'
import { useLearningStore } from '../stores/learning'
import { useGrowthStore } from '../stores/growth'
import ProgressRing from '../components/ProgressRing.vue'

const learning = useLearningStore()
const growth = useGrowthStore()

const weekly = computed(() => learning.weeklyActivity)
const weeklyMax = computed(() =>
  Math.max(1, ...weekly.value.map((d) => d.count))
)

const quizRows = computed(() =>
  quizzes.map((q) => {
    const r = learning.quizResults[q.id]
    return {
      id: q.id,
      title: q.title,
      total: q.questions.length,
      result: r || null,
      percent: r ? Math.round((r.score / r.total) * 100) : null
    }
  })
)

const noteEntries = computed(() =>
  Object.entries(learning.notes)
    .filter(([, text]) => text && text.trim())
    .map(([key, text]) => {
      const [courseId, chapterId] = key.split('/')
      const course = courses.find((c) => c.id === courseId)
      const chapter = course?.chapters.find((ch) => ch.id === chapterId)
      return { key, courseId, chapterId, course, chapter, text }
    })
    .filter((e) => e.course && e.chapter)
)

function confirmReset() {
  if (window.confirm('确定要清空目标、微单元、积分、课程进度、笔记与测验记录吗？此操作不可恢复。')) {
    learning.resetAll()
    growth.reset()
  }
}
</script>

<template>
  <div class="page">
    <header class="page-header">
      <h1>学习统计</h1>
      <p>你的学习进度、活跃度、测验成绩与笔记一览。</p>
    </header>

    <section class="overview">
      <div class="overview__ring card">
        <ProgressRing :percent="learning.overallProgress" :size="120" :stroke="10" />
        <div>
          <div class="overview__big">{{ learning.completedCount }} / {{ totalChapterCount }}</div>
          <div class="overview__label">章节完成</div>
        </div>
      </div>

      <div class="overview__chart card">
        <h3>近 7 天学习活跃度</h3>
        <div class="chart">
          <div v-for="day in weekly" :key="day.date" class="chart__col">
            <span class="chart__count" v-if="day.count">{{ day.count }}</span>
            <div
              class="chart__bar"
              :style="{
                height: (day.count / weeklyMax) * 100 + '%',
                opacity: day.count ? 1 : 0.25
              }"
            ></div>
            <span class="chart__label">{{ day.label }}</span>
          </div>
        </div>
      </div>
    </section>

    <section class="card block">
      <h3>各课程进度</h3>
      <div v-for="course in courses" :key="course.id" class="course-row">
        <span class="course-row__icon">{{ course.icon }}</span>
        <router-link :to="`/courses/${course.id}`" class="course-row__title">
          {{ course.title }}
        </router-link>
        <div class="course-row__track">
          <div
            class="course-row__fill"
            :style="{
              width: learning.courseProgress(course.id) + '%',
              background: course.accent
            }"
          ></div>
        </div>
        <span class="course-row__value">
          {{ learning.courseCompletedCount(course.id) }}/{{ course.chapters.length }}
        </span>
      </div>
    </section>

    <section class="card block">
      <h3>测验成绩</h3>
      <table class="quiz-table">
        <thead>
          <tr>
            <th>测验</th>
            <th>题量</th>
            <th>得分</th>
            <th>操作</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in quizRows" :key="row.id">
            <td>{{ row.title }}</td>
            <td>{{ row.total }} 题</td>
            <td>
              <strong v-if="row.percent !== null" :class="{
                'score--good': row.percent >= 80,
                'score--mid': row.percent >= 60 && row.percent < 80,
                'score--bad': row.percent < 60
              }">{{ row.percent }} 分</strong>
              <span v-else class="score--none">未测验</span>
            </td>
            <td>
              <router-link :to="`/quiz/${row.id}`">
                {{ row.result ? '重测' : '去测验' }} →
              </router-link>
            </td>
          </tr>
        </tbody>
      </table>
    </section>

    <section class="card block">
      <h3>我的笔记（{{ noteEntries.length }}）</h3>
      <p v-if="noteEntries.length === 0" class="muted">
        还没有笔记。在章节学习页底部可以随手记录你的理解与思考。
      </p>
      <div v-for="entry in noteEntries" :key="entry.key" class="note">
        <router-link :to="`/courses/${entry.courseId}/${entry.chapterId}`" class="note__source">
          {{ entry.course.icon }} {{ entry.course.title }} · {{ entry.chapter.title }}
        </router-link>
        <p class="note__text">{{ entry.text }}</p>
      </div>
    </section>

    <section class="danger-zone">
      <button class="btn btn--danger" @click="confirmReset">清空全部学习数据</button>
    </section>
  </div>
</template>

<style scoped>
.overview {
  display: grid;
  grid-template-columns: 300px 1fr;
  gap: 18px;
  margin-bottom: 18px;
}

.overview__ring {
  display: flex;
  align-items: center;
  gap: 20px;
  padding: 26px;
}

.overview__big {
  font-size: 24px;
  font-weight: 800;
}

.overview__label {
  font-size: 13px;
  color: var(--text-2);
}

.overview__chart {
  padding: 20px 26px;
}

.overview__chart h3,
.block h3 {
  margin: 0 0 16px;
  font-size: 15.5px;
}

.chart {
  display: flex;
  align-items: flex-end;
  gap: 12px;
  height: 110px;
}

.chart__col {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: flex-end;
  height: 100%;
  gap: 5px;
}

.chart__count {
  font-size: 11.5px;
  font-weight: 700;
  color: var(--primary-strong);
}

.chart__bar {
  width: 100%;
  max-width: 34px;
  min-height: 5px;
  background: linear-gradient(180deg, #a78bfa, var(--primary));
  border-radius: 6px 6px 3px 3px;
  transition: height 0.4s ease;
}

.chart__label {
  font-size: 11px;
  color: var(--text-3);
}

.block {
  padding: 22px 26px;
  margin-bottom: 18px;
}

.course-row {
  display: flex;
  align-items: center;
  gap: 13px;
  padding: 11px 0;
}

.course-row + .course-row {
  border-top: 1px solid var(--border);
}

.course-row__icon {
  font-size: 19px;
}

.course-row__title {
  width: 240px;
  min-width: 150px;
  font-size: 14px;
  font-weight: 500;
  color: var(--text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.course-row__title:hover {
  color: var(--primary);
}

.course-row__track {
  flex: 1;
  height: 8px;
  background: var(--border);
  border-radius: 999px;
  overflow: hidden;
}

.course-row__fill {
  height: 100%;
  border-radius: 999px;
  transition: width 0.4s ease;
}

.course-row__value {
  font-size: 13px;
  color: var(--text-2);
  min-width: 42px;
  text-align: right;
}

.quiz-table {
  width: 100%;
  border-collapse: collapse;
  font-size: 14px;
}

.quiz-table th,
.quiz-table td {
  text-align: left;
  padding: 10px 8px;
  border-bottom: 1px solid var(--border);
}

.quiz-table th {
  font-size: 12.5px;
  color: var(--text-3);
  font-weight: 600;
}

.score--good { color: var(--success); }
.score--mid { color: var(--warning); }
.score--bad { color: var(--danger); }
.score--none { color: var(--text-3); font-size: 13px; }

.muted {
  color: var(--text-3);
  font-size: 13.5px;
  margin: 0;
}

.note {
  padding: 13px 0;
}

.note + .note {
  border-top: 1px solid var(--border);
}

.note__source {
  font-size: 12.5px;
  font-weight: 600;
  color: var(--text-2);
}

.note__source:hover {
  color: var(--primary);
}

.note__text {
  margin: 6px 0 0;
  font-size: 14px;
  white-space: pre-wrap;
}

.danger-zone {
  text-align: center;
  padding: 8px 0 20px;
}

@media (max-width: 780px) {
  .overview {
    grid-template-columns: 1fr;
  }

  .course-row__title {
    width: auto;
    flex: 1;
  }

  .course-row__track {
    display: none;
  }
}
</style>
