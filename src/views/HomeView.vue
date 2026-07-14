<script setup>
import { computed } from 'vue'
import { courses, getCourse, totalChapterCount } from '../data/courses'
import { useLearningStore } from '../stores/learning'
import { useAuthStore } from '../stores/auth'
import CourseCard from '../components/CourseCard.vue'
import ProgressRing from '../components/ProgressRing.vue'

const learning = useLearningStore()
const auth = useAuthStore()

const continueTarget = computed(() => {
  const last = learning.lastVisited
  if (!last) return null
  const course = getCourse(last.courseId)
  if (!course) return null
  const chapter = course.chapters.find((ch) => ch.id === last.chapterId)
  if (!chapter) return null
  return { course, chapter }
})

const stats = computed(() => [
  { label: '已完成章节', value: `${learning.completedCount} / ${totalChapterCount}`, icon: '✅' },
  { label: '学习笔记', value: learning.noteCount, icon: '🗒️' },
  {
    label: '测验平均分',
    value: learning.quizAverageScore === null ? '—' : learning.quizAverageScore + ' 分',
    icon: '🏆'
  }
])
</script>

<template>
  <div class="page">
    <section class="hero card">
      <div class="hero__text">
        <h1>
          {{ auth.user ? `${auth.user.avatar} 欢迎回来,${auth.user.nickname}` : '欢迎来到 AI 学习系统' }}
        </h1>
        <p>
          从提示词工程到 AI 原生架构的系统化学习平台。跟随课程循序渐进，
          配合 AI 学习助手答疑与章节测验巩固，完成从"写代码"到"调度智能"的跃迁。
        </p>
        <div class="hero__actions">
          <router-link
            v-if="continueTarget"
            :to="`/courses/${continueTarget.course.id}/${continueTarget.chapter.id}`"
            class="btn btn--primary"
          >
            ▶ 继续学习：{{ continueTarget.chapter.title }}
          </router-link>
          <router-link v-else to="/courses" class="btn btn--primary">🚀 开始学习</router-link>
          <router-link to="/chat" class="btn btn--ghost">💬 问问 AI 助手</router-link>
        </div>
      </div>
      <div class="hero__ring">
        <ProgressRing :percent="learning.overallProgress" :size="110" :stroke="9" />
        <span class="hero__ring-label">总体进度</span>
      </div>
    </section>

    <section class="stats">
      <div v-for="s in stats" :key="s.label" class="stats__item card">
        <span class="stats__icon">{{ s.icon }}</span>
        <div>
          <div class="stats__value">{{ s.value }}</div>
          <div class="stats__label">{{ s.label }}</div>
        </div>
      </div>
    </section>

    <section>
      <div class="section-head">
        <h2>课程体系</h2>
        <router-link to="/courses">查看全部 →</router-link>
      </div>
      <div class="course-grid">
        <CourseCard v-for="course in courses" :key="course.id" :course="course" />
      </div>
    </section>
  </div>
</template>

<style scoped>
.hero {
  display: flex;
  align-items: center;
  gap: 28px;
  padding: 34px 36px;
  margin-bottom: 22px;
  background: linear-gradient(120deg, #eef2ff 0%, #ffffff 55%, #f0f9ff 100%);
}

.hero__text {
  flex: 1;
}

.hero__text h1 {
  margin: 0 0 10px;
  font-size: 27px;
  letter-spacing: -0.02em;
}

.hero__text p {
  margin: 0 0 20px;
  color: var(--text-2);
  max-width: 560px;
}

.hero__actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.hero__ring {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.hero__ring-label {
  font-size: 13px;
  color: var(--text-2);
}

.stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-bottom: 34px;
}

.stats__item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 18px 20px;
}

.stats__icon {
  font-size: 26px;
}

.stats__value {
  font-size: 19px;
  font-weight: 700;
}

.stats__label {
  font-size: 13px;
  color: var(--text-2);
}

.section-head {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 16px;
}

.section-head h2 {
  margin: 0;
  font-size: 19px;
}

.course-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(290px, 1fr));
  gap: 18px;
}

@media (max-width: 720px) {
  .hero {
    flex-direction: column-reverse;
    text-align: center;
    padding: 26px 22px;
  }

  .hero__actions {
    justify-content: center;
  }

  .stats {
    grid-template-columns: 1fr;
  }
}
</style>
