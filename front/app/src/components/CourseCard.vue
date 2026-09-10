<script setup>
import { computed } from 'vue'
import { useLearningStore } from '../stores/learning'
import Icon from './Icon.vue'

const props = defineProps({
  course: { type: Object, required: true }
})

const learning = useLearningStore()
const progress = computed(() => learning.courseProgress(props.course.id))
const done = computed(() => learning.courseCompletedCount(props.course.id))
</script>

<template>
  <router-link :to="`/courses/${course.id}`" class="course-card card card--hover">
    <div class="course-card__cover" :style="{ '--accent': course.accent }">
      <span class="course-card__icon">{{ course.icon }}</span>
      <span class="tag" :class="`tag--level-${course.level}`">{{ course.level }}</span>
    </div>
    <div class="course-card__body">
      <h3 class="course-card__title">{{ course.title }}</h3>
      <p class="course-card__desc">{{ course.description }}</p>
      <div class="course-card__tags">
        <span v-for="tag in course.tags" :key="tag" class="tag tag--neutral">{{ tag }}</span>
      </div>
      <div class="course-card__meta">
        <span><Icon name="layers" :size="13" /> {{ course.chapters.length }} 章节</span>
        <span><Icon name="clock" :size="13" /> 约 {{ course.estimatedMinutes }} 分钟</span>
        <span class="course-card__done num">{{ done }}/{{ course.chapters.length }}</span>
      </div>
      <div class="progress progress--thin">
        <i :style="{ width: progress + '%', background: course.accent }"></i>
      </div>
    </div>
  </router-link>
</template>

<style scoped>
.course-card {
  display: flex;
  flex-direction: column;
  overflow: hidden;
  color: var(--text);
}

.course-card__cover {
  position: relative;
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  padding: 20px 20px 0;
  height: 112px;
  background:
    radial-gradient(120% 100% at 100% 0%, color-mix(in srgb, var(--accent) 28%, transparent), transparent 60%),
    linear-gradient(180deg, color-mix(in srgb, var(--accent) 14%, #fff), #fff);
}

.course-card__cover::after {
  content: '';
  position: absolute;
  inset: auto 0 0;
  height: 1px;
  background: linear-gradient(90deg, transparent, var(--border), transparent);
}

.course-card__icon {
  display: grid;
  place-items: center;
  width: 60px;
  height: 60px;
  border-radius: 18px;
  font-size: 30px;
  background: #fff;
  box-shadow: var(--shadow-sm), inset 0 0 0 1px color-mix(in srgb, var(--accent) 22%, transparent);
  transform: translateY(6px);
  transition: transform var(--t);
}

.course-card:hover .course-card__icon {
  transform: translateY(2px) rotate(-4deg) scale(1.04);
}

.course-card__body {
  display: flex;
  flex: 1;
  flex-direction: column;
  padding: 18px 20px 20px;
}

.course-card__title {
  margin: 0 0 8px;
  font-size: 17px;
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
  margin-bottom: 16px;
}

.course-card__meta {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 12.5px;
  color: var(--text-3);
  margin-bottom: 10px;
}

.course-card__meta span {
  display: inline-flex;
  align-items: center;
  gap: 4px;
}

.course-card__done {
  margin-left: auto;
  font-weight: 600;
  color: var(--text-2);
}
</style>
