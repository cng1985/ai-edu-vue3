<script setup>
import { ref, computed } from 'vue'
import { courses, totalChapterCount } from '../data/courses'
import { useLearningStore } from '../stores/learning'
import CourseCard from '../components/CourseCard.vue'
import Icon from '../components/Icon.vue'

const learning = useLearningStore()
const keyword = ref('')
const level = ref('全部')
const levels = ['全部', '入门', '进阶', '高级']

const filtered = computed(() =>
  courses.filter((c) => {
    const matchLevel = level.value === '全部' || c.level === level.value
    const kw = keyword.value.trim().toLowerCase()
    const matchKw =
      !kw ||
      c.title.toLowerCase().includes(kw) ||
      c.description.toLowerCase().includes(kw) ||
      c.tags.some((t) => t.toLowerCase().includes(kw))
    return matchLevel && matchKw
  })
)

const totalMinutes = computed(() => courses.reduce((sum, c) => sum + c.estimatedMinutes, 0))
</script>

<template>
  <div class="page">
    <header class="page-header page-header--split">
      <div>
        <span class="eyebrow">课程体系</span>
        <h1>全部课程</h1>
        <p>覆盖提示词工程、RAG 与 AI 原生应用开发全链路，内容以 Markdown 结构化管理。</p>
      </div>
      <div class="summary">
        <div><strong class="num">{{ courses.length }}</strong><span>门课程</span></div>
        <div><strong class="num">{{ totalChapterCount }}</strong><span>个章节</span></div>
        <div><strong class="num">{{ Math.round(totalMinutes / 60) }}h</strong><span>学习时长</span></div>
        <div><strong class="num">{{ learning.overallProgress }}%</strong><span>总进度</span></div>
      </div>
    </header>

    <div class="toolbar">
      <label class="toolbar__search">
        <Icon name="search" :size="17" />
        <input
          v-model="keyword"
          type="search"
          placeholder="搜索课程标题、简介或标签…"
        />
      </label>
      <div class="segmented">
        <button
          v-for="l in levels"
          :key="l"
          :class="{ active: level === l }"
          @click="level = l"
        >
          {{ l }}
        </button>
      </div>
    </div>

    <div v-if="filtered.length" class="course-grid stagger">
      <CourseCard v-for="course in filtered" :key="course.id" :course="course" />
    </div>
    <div v-else class="empty-state card">
      <div class="empty-state__icon"><Icon name="search" :size="30" /></div>
      <h3>没有匹配的课程</h3>
      <p>换个关键词试试，或者清除等级筛选。</p>
      <button class="btn btn--ghost" @click="keyword = ''; level = '全部'">清除筛选</button>
    </div>
  </div>
</template>

<style scoped>
.summary {
  display: flex;
  gap: 4px;
  padding: 6px;
  border-radius: 16px;
  background: var(--surface);
  border: 1px solid var(--border);
  box-shadow: var(--shadow-sm);
}

.summary div {
  display: flex;
  flex-direction: column;
  align-items: center;
  min-width: 78px;
  padding: 8px 12px;
  border-radius: 11px;
}

.summary div + div {
  position: relative;
}

.summary div + div::before {
  content: '';
  position: absolute;
  left: 0;
  top: 25%;
  height: 50%;
  width: 1px;
  background: var(--border);
}

.summary strong {
  font-size: 20px;
  font-weight: 700;
  line-height: 1.1;
  color: var(--text);
}

.summary span {
  color: var(--text-3);
  font-size: 11.5px;
}

.toolbar {
  display: flex;
  gap: 14px;
  margin-bottom: 24px;
  flex-wrap: wrap;
  align-items: center;
}

.toolbar__search {
  flex: 1;
  min-width: 240px;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 0 16px;
  height: 46px;
  border: 1px solid var(--border-strong);
  border-radius: 14px;
  background: var(--surface);
  color: var(--text-3);
  transition: border-color var(--t-fast), box-shadow var(--t-fast);
}

.toolbar__search:focus-within {
  border-color: var(--primary);
  box-shadow: 0 0 0 4px var(--primary-soft-2);
}

.toolbar__search input {
  flex: 1;
  border: none;
  background: transparent;
  font: inherit;
  font-size: 14.5px;
  color: var(--text);
  outline: none;
}

.toolbar__search input::placeholder {
  color: var(--text-3);
}

.course-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 20px;
}

@media (max-width: 860px) {
  .summary { width: 100%; overflow-x: auto; }
}
</style>
