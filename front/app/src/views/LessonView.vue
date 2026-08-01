<script setup>
import { ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getChapter } from '../data/courses'
import { useLearningStore } from '../stores/learning'
import MarkdownRenderer from '../components/MarkdownRenderer.vue'

const route = useRoute()
const learning = useLearningStore()

const lesson = computed(() =>
  getChapter(route.params.courseId, route.params.chapterId)
)

const completed = computed(() =>
  lesson.value
    ? learning.isChapterCompleted(lesson.value.course.id, lesson.value.chapter.id)
    : false
)

const noteDraft = ref('')
const noteSaved = ref(false)

watch(
  lesson,
  (val) => {
    if (val) {
      learning.recordVisit(val.course.id, val.chapter.id)
      noteDraft.value = learning.noteFor(val.course.id, val.chapter.id)
      noteSaved.value = false
    }
  },
  { immediate: true }
)

function toggleCompleted() {
  learning.toggleChapterCompleted(lesson.value.course.id, lesson.value.chapter.id)
}

function saveNote() {
  learning.saveNote(lesson.value.course.id, lesson.value.chapter.id, noteDraft.value)
  noteSaved.value = true
  setTimeout(() => (noteSaved.value = false), 1600)
}
</script>

<template>
  <div class="lesson" v-if="lesson">
    <aside class="lesson__toc card">
      <router-link :to="`/courses/${lesson.course.id}`" class="lesson__course">
        {{ lesson.course.icon }} {{ lesson.course.title }}
      </router-link>
      <nav>
        <router-link
          v-for="(ch, i) in lesson.course.chapters"
          :key="ch.id"
          :to="`/courses/${lesson.course.id}/${ch.id}`"
          class="lesson__toc-item"
          :class="{ 'lesson__toc-item--active': ch.id === lesson.chapter.id }"
        >
          <span
            class="lesson__toc-dot"
            :class="{
              'lesson__toc-dot--done': learning.isChapterCompleted(lesson.course.id, ch.id)
            }"
          >
            {{ learning.isChapterCompleted(lesson.course.id, ch.id) ? '✓' : i + 1 }}
          </span>
          {{ ch.title }}
        </router-link>
      </nav>
    </aside>

    <div class="lesson__main">
      <article class="lesson__content card fade-up" :key="lesson.chapter.id">
        <MarkdownRenderer :source="lesson.chapter.content" />

        <div class="lesson__complete">
          <button
            class="btn"
            :class="completed ? 'btn--ghost' : 'btn--primary'"
            @click="toggleCompleted"
          >
            {{ completed ? '✓ 已完成本章（点击取消）' : '完成本章学习' }}
          </button>
        </div>
      </article>

      <section class="lesson__notes card">
        <h3>🗒️ 本章笔记</h3>
        <textarea
          v-model="noteDraft"
          rows="5"
          placeholder="记录你的理解、疑问或延伸思考…（笔记保存在本地浏览器中）"
        ></textarea>
        <div class="lesson__notes-actions">
          <span v-if="noteSaved" class="lesson__notes-saved">✓ 已保存</span>
          <button class="btn btn--primary" @click="saveNote">保存笔记</button>
        </div>
      </section>

      <nav class="lesson__pager">
        <router-link
          v-if="lesson.prev"
          :to="`/courses/${lesson.course.id}/${lesson.prev.id}`"
          class="lesson__pager-link card"
        >
          <span class="lesson__pager-dir">← 上一章</span>
          <span class="lesson__pager-title">{{ lesson.prev.title }}</span>
        </router-link>
        <span v-else></span>
        <router-link
          v-if="lesson.next"
          :to="`/courses/${lesson.course.id}/${lesson.next.id}`"
          class="lesson__pager-link lesson__pager-link--next card"
        >
          <span class="lesson__pager-dir">下一章 →</span>
          <span class="lesson__pager-title">{{ lesson.next.title }}</span>
        </router-link>
        <router-link v-else to="/quiz" class="lesson__pager-link lesson__pager-link--next card">
          <span class="lesson__pager-dir">课程完结 🎉</span>
          <span class="lesson__pager-title">去做课程测验</span>
        </router-link>
      </nav>
    </div>
  </div>

  <div class="page" v-else>
    <div class="card" style="padding: 48px; text-align: center">
      <p>未找到该章节。</p>
      <router-link to="/courses" class="btn btn--primary">返回课程列表</router-link>
    </div>
  </div>
</template>

<style scoped>
.lesson {
  display: flex;
  gap: 22px;
  max-width: 1240px;
  margin: 0 auto;
  padding: 28px 32px 64px;
  align-items: flex-start;
}

.lesson__toc {
  width: 256px;
  min-width: 256px;
  padding: 18px 14px;
  position: sticky;
  top: 24px;
  max-height: calc(100vh - 48px);
  overflow-y: auto;
}

.lesson__course {
  display: block;
  font-size: 14px;
  font-weight: 700;
  color: var(--text);
  padding: 4px 10px 14px;
  border-bottom: 1px solid var(--border);
  margin-bottom: 10px;
}

.lesson__toc-item {
  display: flex;
  align-items: center;
  gap: 9px;
  padding: 8px 10px;
  border-radius: var(--radius-sm);
  font-size: 13px;
  color: var(--text-2);
  line-height: 1.4;
  transition: background 0.15s ease;
}

.lesson__toc-item:hover {
  background: var(--surface-2);
}

.lesson__toc-item--active {
  background: var(--primary-soft);
  color: var(--primary-strong);
  font-weight: 600;
}

.lesson__toc-dot {
  width: 21px;
  height: 21px;
  min-width: 21px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: var(--surface-2);
  border: 1px solid var(--border);
  font-size: 11px;
  font-weight: 600;
}

.lesson__toc-dot--done {
  background: var(--success-soft);
  border-color: var(--success);
  color: var(--success);
}

.lesson__main {
  flex: 1;
  min-width: 0;
}

.lesson__content {
  padding: 36px 42px;
}

.lesson__complete {
  margin-top: 32px;
  padding-top: 22px;
  border-top: 1px solid var(--border);
  text-align: center;
}

.lesson__notes {
  margin-top: 20px;
  padding: 22px 26px;
}

.lesson__notes h3 {
  margin: 0 0 12px;
  font-size: 15.5px;
}

.lesson__notes textarea {
  width: 100%;
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  padding: 12px 14px;
  font-family: inherit;
  font-size: 14px;
  line-height: 1.7;
  resize: vertical;
  outline: none;
  transition: border-color 0.15s ease;
}

.lesson__notes textarea:focus {
  border-color: var(--primary);
}

.lesson__notes-actions {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 12px;
  margin-top: 10px;
}

.lesson__notes-saved {
  color: var(--success);
  font-size: 13px;
  font-weight: 600;
}

.lesson__pager {
  display: flex;
  justify-content: space-between;
  gap: 16px;
  margin-top: 20px;
}

.lesson__pager-link {
  display: flex;
  flex-direction: column;
  gap: 3px;
  padding: 14px 20px;
  max-width: 46%;
  color: var(--text);
  transition: transform 0.15s ease;
}

.lesson__pager-link:hover {
  transform: translateY(-2px);
}

.lesson__pager-link--next {
  text-align: right;
  margin-left: auto;
}

.lesson__pager-dir {
  font-size: 12px;
  color: var(--text-3);
}

.lesson__pager-title {
  font-size: 14px;
  font-weight: 600;
}

@media (max-width: 1000px) {
  .lesson {
    flex-direction: column;
    padding: 20px 16px 48px;
  }

  .lesson__toc {
    width: 100%;
    min-width: 0;
    position: static;
    max-height: none;
  }

  .lesson__content {
    padding: 24px 20px;
  }
}
</style>
