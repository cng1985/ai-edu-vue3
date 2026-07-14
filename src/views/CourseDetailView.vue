<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { getCourse } from '../data/courses'
import { getQuiz } from '../data/quizzes'
import { useLearningStore } from '../stores/learning'
import ProgressRing from '../components/ProgressRing.vue'

const route = useRoute()
const learning = useLearningStore()

const course = computed(() => getCourse(route.params.courseId))
const quiz = computed(() => (course.value ? getQuiz(course.value.id) : null))
const progress = computed(() =>
  course.value ? learning.courseProgress(course.value.id) : 0
)

const firstUnfinished = computed(() => {
  if (!course.value) return null
  return (
    course.value.chapters.find(
      (ch) => !learning.isChapterCompleted(course.value.id, ch.id)
    ) || course.value.chapters[0]
  )
})
</script>

<template>
  <div class="page" v-if="course">
    <router-link to="/courses" class="back">← 返回课程列表</router-link>

    <section class="head card">
      <span class="head__icon" :style="{ background: course.accent + '1a' }">
        {{ course.icon }}
      </span>
      <div class="head__info">
        <div class="head__tags">
          <span class="tag" :class="`tag--level-${course.level}`">{{ course.level }}</span>
          <span v-for="tag in course.tags" :key="tag" class="tag">{{ tag }}</span>
        </div>
        <h1>{{ course.title }}</h1>
        <p>{{ course.description }}</p>
        <div class="head__actions">
          <router-link
            v-if="firstUnfinished"
            :to="`/courses/${course.id}/${firstUnfinished.id}`"
            class="btn btn--primary"
          >
            {{ progress > 0 ? '▶ 继续学习' : '🚀 开始学习' }}
          </router-link>
          <router-link v-if="quiz" :to="`/quiz/${quiz.id}`" class="btn btn--ghost">
            📝 课程测验
          </router-link>
        </div>
      </div>
      <div class="head__ring">
        <ProgressRing :percent="progress" :size="92" :stroke="8" :color="course.accent" />
      </div>
    </section>

    <section class="chapters card">
      <h2>章节目录（{{ course.chapters.length }} 章 · 约 {{ course.estimatedMinutes }} 分钟）</h2>
      <router-link
        v-for="(chapter, i) in course.chapters"
        :key="chapter.id"
        :to="`/courses/${course.id}/${chapter.id}`"
        class="chapter"
      >
        <span
          class="chapter__status"
          :class="{ 'chapter__status--done': learning.isChapterCompleted(course.id, chapter.id) }"
        >
          {{ learning.isChapterCompleted(course.id, chapter.id) ? '✓' : i + 1 }}
        </span>
        <span class="chapter__title">{{ chapter.title }}</span>
        <span class="chapter__minutes">约 {{ chapter.minutes }} 分钟</span>
      </router-link>
    </section>
  </div>

  <div class="page" v-else>
    <div class="card" style="padding: 48px; text-align: center">
      <p>未找到该课程。</p>
      <router-link to="/courses" class="btn btn--primary">返回课程列表</router-link>
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

.head {
  display: flex;
  gap: 22px;
  padding: 30px;
  margin-bottom: 20px;
}

.head__icon {
  width: 64px;
  height: 64px;
  min-width: 64px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 32px;
  border-radius: 16px;
}

.head__info {
  flex: 1;
}

.head__tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
  margin-bottom: 10px;
}

.head__info h1 {
  margin: 0 0 8px;
  font-size: 23px;
  letter-spacing: -0.01em;
}

.head__info p {
  margin: 0 0 18px;
  color: var(--text-2);
  font-size: 14.5px;
}

.head__actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.chapters {
  padding: 26px 30px;
}

.chapters h2 {
  margin: 0 0 16px;
  font-size: 17px;
}

.chapter {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 13px 12px;
  border-radius: var(--radius-sm);
  color: var(--text);
  transition: background 0.15s ease;
}

.chapter:hover {
  background: var(--surface-2);
}

.chapter + .chapter {
  border-top: 1px solid var(--border);
}

.chapter__status {
  width: 28px;
  height: 28px;
  min-width: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: var(--surface-2);
  border: 1px solid var(--border);
  font-size: 13px;
  font-weight: 600;
  color: var(--text-2);
}

.chapter__status--done {
  background: var(--success-soft);
  border-color: var(--success);
  color: var(--success);
}

.chapter__title {
  flex: 1;
  font-size: 14.5px;
  font-weight: 500;
}

.chapter__minutes {
  font-size: 12.5px;
  color: var(--text-3);
}

@media (max-width: 720px) {
  .head {
    flex-direction: column;
  }
}
</style>
