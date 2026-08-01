<script setup>
import { ref, computed } from 'vue'
import { courses } from '../data/courses'
import CourseCard from '../components/CourseCard.vue'

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
</script>

<template>
  <div class="page">
    <header class="page-header">
      <h1>全部课程</h1>
      <p>共 {{ courses.length }} 门课程，覆盖提示词工程、RAG 与 AI 原生应用开发全链路。</p>
    </header>

    <div class="toolbar">
      <input
        v-model="keyword"
        type="search"
        class="toolbar__search"
        placeholder="搜索课程标题、简介或标签…"
      />
      <div class="toolbar__levels">
        <button
          v-for="l in levels"
          :key="l"
          class="toolbar__level"
          :class="{ 'toolbar__level--active': level === l }"
          @click="level = l"
        >
          {{ l }}
        </button>
      </div>
    </div>

    <div v-if="filtered.length" class="course-grid">
      <CourseCard v-for="course in filtered" :key="course.id" :course="course" />
    </div>
    <div v-else class="empty card">
      <p>没有匹配的课程，换个关键词试试？</p>
    </div>
  </div>
</template>

<style scoped>
.toolbar {
  display: flex;
  gap: 14px;
  margin-bottom: 22px;
  flex-wrap: wrap;
}

.toolbar__search {
  flex: 1;
  min-width: 220px;
  padding: 10px 16px;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  font-size: 14px;
  background: var(--surface);
  outline: none;
  transition: border-color 0.15s ease;
}

.toolbar__search:focus {
  border-color: var(--primary);
}

.toolbar__levels {
  display: flex;
  gap: 8px;
}

.toolbar__level {
  padding: 9px 16px;
  border: 1px solid var(--border);
  background: var(--surface);
  border-radius: var(--radius-sm);
  font-size: 13.5px;
  color: var(--text-2);
  cursor: pointer;
  transition: all 0.15s ease;
}

.toolbar__level--active {
  background: var(--primary);
  border-color: var(--primary);
  color: #fff;
  font-weight: 600;
}

.course-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(290px, 1fr));
  gap: 18px;
}

.empty {
  padding: 48px;
  text-align: center;
  color: var(--text-2);
}
</style>
