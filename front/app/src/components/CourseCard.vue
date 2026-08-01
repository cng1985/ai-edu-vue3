<script setup>
import { computed } from 'vue'
import { useLearningStore } from '../stores/learning'

const props = defineProps({
  course: { type: Object, required: true }
})

const learning = useLearningStore()
const progress = computed(() => learning.courseProgress(props.course.id))
const done = computed(() => learning.courseCompletedCount(props.course.id))
</script>

<template>
  <router-link :to="`/courses/${course.id}`" class="course-card card fade-up">
    <div class="course-card__top">
      <span class="course-card__icon" :style="{ background: course.accent + '1a' }">
        {{ course.icon }}
      </span>
      <span class="tag" :class="`tag--level-${course.level}`">{{ course.level }}</span>
    </div>
    <h3 class="course-card__title">{{ course.title }}</h3>
    <p class="course-card__desc">{{ course.description }}</p>
    <div class="course-card__tags">
      <span v-for="tag in course.tags" :key="tag" class="tag">{{ tag }}</span>
    </div>
    <div class="course-card__meta">
      <span>{{ course.chapters.length }} 章节 · 约 {{ course.estimatedMinutes }} 分钟</span>
      <span class="course-card__done">{{ done }}/{{ course.chapters.length }} 已完成</span>
    </div>
    <div class="course-card__track">
      <div
        class="course-card__fill"
        :style="{ width: progress + '%', background: course.accent }"
      ></div>
    </div>
  </router-link>
</template>

<style scoped>
.course-card {
  display: flex;
  flex-direction: column;
  padding: 22px;
  color: var(--text);
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.course-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 4px 12px rgba(15, 23, 42, 0.08), 0 16px 40px rgba(15, 23, 42, 0.08);
}

.course-card__top {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 14px;
}

.course-card__icon {
  width: 46px;
  height: 46px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  border-radius: 12px;
}

.course-card__title {
  margin: 0 0 8px;
  font-size: 17px;
  font-weight: 700;
}

.course-card__desc {
  margin: 0 0 14px;
  font-size: 13.5px;
  color: var(--text-2);
  line-height: 1.65;
  flex: 1;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.course-card__tags {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
  margin-bottom: 14px;
}

.course-card__meta {
  display: flex;
  justify-content: space-between;
  font-size: 12.5px;
  color: var(--text-3);
  margin-bottom: 8px;
}

.course-card__done {
  font-weight: 600;
  color: var(--text-2);
}

.course-card__track {
  height: 6px;
  background: var(--border);
  border-radius: 999px;
  overflow: hidden;
}

.course-card__fill {
  height: 100%;
  border-radius: 999px;
  transition: width 0.4s ease;
}
</style>
