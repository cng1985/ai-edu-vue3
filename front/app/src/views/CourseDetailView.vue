<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { getCourse } from '../data/courses'
import { getQuiz } from '../data/quizzes'
import { useLearningStore } from '../stores/learning'
import ProgressRing from '../components/ProgressRing.vue'
import Icon from '../components/Icon.vue'

const route = useRoute()
const learning = useLearningStore()

const course = computed(() => getCourse(route.params.courseId))
const quiz = computed(() => (course.value ? getQuiz(course.value.id) : null))
const progress = computed(() =>
  course.value ? learning.courseProgress(course.value.id) : 0
)
const doneCount = computed(() => (course.value ? learning.courseCompletedCount(course.value.id) : 0))

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
    <router-link to="/courses" class="back-link"><Icon name="arrowLeft" :size="15" /> 返回课程列表</router-link>

    <section class="head card" :style="{ '--accent': course.accent }">
      <div class="head__cover">
        <span class="head__icon">{{ course.icon }}</span>
      </div>
      <div class="head__info">
        <div class="head__tags">
          <span class="tag" :class="`tag--level-${course.level}`">{{ course.level }}</span>
          <span v-for="tag in course.tags" :key="tag" class="tag tag--neutral">{{ tag }}</span>
        </div>
        <h1>{{ course.title }}</h1>
        <p>{{ course.description }}</p>
        <div class="head__meta">
          <span><Icon name="layers" :size="14" /> {{ course.chapters.length }} 章节</span>
          <span><Icon name="clock" :size="14" /> 约 {{ course.estimatedMinutes }} 分钟</span>
          <span><Icon name="check" :size="14" /> 已完成 {{ doneCount }} 章</span>
        </div>
        <div class="head__actions">
          <router-link
            v-if="firstUnfinished"
            :to="`/courses/${course.id}/${firstUnfinished.id}`"
            class="btn btn--primary btn--lg"
          >
            <Icon name="play" :size="16" /> {{ progress > 0 ? '继续学习' : '开始学习' }}
          </router-link>
          <router-link v-if="quiz" :to="`/quiz/${quiz.id}`" class="btn btn--ghost btn--lg">
            <Icon name="clipboard" :size="16" /> 课程测验
          </router-link>
        </div>
      </div>
      <div class="head__ring">
        <ProgressRing :percent="progress" :size="112" :stroke="10" :color="course.accent" />
        <span>课程进度</span>
      </div>
    </section>

    <section class="chapters card">
      <div class="panel-head">
        <div>
          <h2>章节目录</h2>
          <p>{{ course.chapters.length }} 章 · 约 {{ course.estimatedMinutes }} 分钟 · 按顺序学习效果最佳</p>
        </div>
      </div>
      <ol class="chapter-list">
        <li v-for="(chapter, i) in course.chapters" :key="chapter.id">
          <router-link
            :to="`/courses/${course.id}/${chapter.id}`"
            class="chapter"
            :class="{ 'chapter--done': learning.isChapterCompleted(course.id, chapter.id), 'chapter--next': firstUnfinished?.id === chapter.id && progress < 100 }"
          >
            <span class="chapter__status">
              <Icon v-if="learning.isChapterCompleted(course.id, chapter.id)" name="check" :size="14" :stroke="3" />
              <template v-else>{{ i + 1 }}</template>
            </span>
            <span class="chapter__title">{{ chapter.title }}</span>
            <span v-if="firstUnfinished?.id === chapter.id && progress < 100" class="tag">继续</span>
            <span class="chapter__minutes"><Icon name="clock" :size="12" /> {{ chapter.minutes }} 分钟</span>
            <Icon name="arrowRight" :size="15" class="chapter__arrow" />
          </router-link>
        </li>
      </ol>
    </section>
  </div>

  <div class="page" v-else>
    <div class="empty-state card">
      <div class="empty-state__icon"><Icon name="alert" :size="30" /></div>
      <h2>未找到该课程</h2>
      <router-link to="/courses" class="btn btn--primary">返回课程列表</router-link>
    </div>
  </div>
</template>

<style scoped>
.head {
  display: flex;
  gap: 26px;
  padding: 30px;
  margin-bottom: 20px;
  overflow: hidden;
  position: relative;
}

.head::before {
  content: '';
  position: absolute;
  inset: 0 auto 0 0;
  width: 40%;
  background: radial-gradient(80% 100% at 0% 50%, color-mix(in srgb, var(--accent) 14%, transparent), transparent 70%);
  pointer-events: none;
}

.head__cover {
  position: relative;
  display: grid;
  place-items: center;
  width: 108px;
  height: 108px;
  flex: 0 0 108px;
  border-radius: 28px;
  background: linear-gradient(145deg, color-mix(in srgb, var(--accent) 22%, #fff), color-mix(in srgb, var(--accent) 8%, #fff));
  box-shadow: inset 0 0 0 1px color-mix(in srgb, var(--accent) 24%, transparent);
}

.head__icon {
  font-size: 48px;
  filter: drop-shadow(0 6px 10px rgba(22, 24, 44, 0.15));
}

.head__info {
  position: relative;
  flex: 1;
  min-width: 0;
}

.head__tags {
  display: flex;
  gap: 6px;
  flex-wrap: wrap;
  margin-bottom: 12px;
}

.head__info h1 {
  margin: 0 0 10px;
  font-size: 26px;
}

.head__info p {
  margin: 0 0 14px;
  color: var(--text-2);
  font-size: 15px;
  line-height: 1.7;
}

.head__meta {
  display: flex;
  flex-wrap: wrap;
  gap: 16px;
  margin-bottom: 20px;
  color: var(--text-3);
  font-size: 13px;
}

.head__meta span {
  display: inline-flex;
  align-items: center;
  gap: 5px;
}

.head__actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.head__ring {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 8px;
  padding: 0 10px;
}

.head__ring span {
  color: var(--text-3);
  font-size: 12px;
  font-weight: 600;
}

.chapters {
  padding: 26px 30px;
}

.chapter-list {
  list-style: none;
  margin: 0;
  padding: 0;
}

.chapter {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 14px 14px;
  border-radius: 13px;
  color: var(--text);
  transition: background var(--t-fast);
}

.chapter:hover {
  background: var(--surface-2);
}

.chapter:hover .chapter__arrow {
  opacity: 1;
  transform: translateX(0);
}

.chapter-list li + li .chapter {
  border-top: 1px solid var(--border);
  border-radius: 0 0 13px 13px;
}

.chapter--next {
  background: var(--primary-soft);
}

.chapter__status {
  width: 30px;
  height: 30px;
  min-width: 30px;
  display: grid;
  place-items: center;
  border-radius: 10px;
  background: var(--surface-3);
  font-size: 13px;
  font-weight: 700;
  color: var(--text-2);
}

.chapter--done .chapter__status {
  background: var(--success);
  color: #fff;
}

.chapter--next .chapter__status {
  background: linear-gradient(135deg, var(--primary), var(--primary-strong));
  color: #fff;
}

.chapter__title {
  flex: 1;
  font-size: 14.5px;
  font-weight: 500;
}

.chapter--done .chapter__title {
  color: var(--text-2);
}

.chapter__minutes {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 12.5px;
  color: var(--text-3);
}

.chapter__arrow {
  color: var(--primary);
  opacity: 0;
  transform: translateX(-4px);
  transition: all var(--t-fast);
}

@media (max-width: 760px) {
  .head {
    flex-direction: column;
  }

  .head__ring {
    flex-direction: row;
    justify-content: flex-start;
  }
}
</style>
