<script setup>
import { ref, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import { getChapter } from '../data/courses'
import { useLearningStore } from '../stores/learning'
import { useAuthStore } from '../stores/auth'
import { meApi } from '../api'
import MarkdownRenderer from '../components/MarkdownRenderer.vue'
import Icon from '../components/Icon.vue'

const route = useRoute()
const learning = useLearningStore()
const auth = useAuthStore()
const synced = ref([])

const lesson = computed(() =>
  getChapter(route.params.courseId, route.params.chapterId)
)

const completed = computed(() =>
  lesson.value
    ? learning.isChapterCompleted(lesson.value.course.id, lesson.value.chapter.id)
    : false
)

const chapterIndex = computed(() =>
  lesson.value ? lesson.value.course.chapters.findIndex((ch) => ch.id === lesson.value.chapter.id) : 0
)

const courseProgress = computed(() => (lesson.value ? learning.courseProgress(lesson.value.course.id) : 0))

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

async function toggleCompleted() {
  const { course, chapter } = lesson.value
  const wasCompleted = completed.value
  learning.toggleChapterCompleted(course.id, chapter.id)
  synced.value = []
  if (wasCompleted || !auth.hasPermission('growth:write')) return
  try {
    synced.value = await meApi.completeChapter({ courseId: course.id, chapterId: chapter.id })
  } catch {
    synced.value = []
  }
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
        <span class="lesson__course-icon">{{ lesson.course.icon }}</span>
        <span class="lesson__course-text">
          <strong>{{ lesson.course.title }}</strong>
          <small>{{ learning.courseCompletedCount(lesson.course.id) }}/{{ lesson.course.chapters.length }} 章已完成</small>
        </span>
      </router-link>
      <div class="progress progress--thin lesson__course-progress"><i :style="{ width: courseProgress + '%', background: lesson.course.accent }"></i></div>
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
            :class="{ 'lesson__toc-dot--done': learning.isChapterCompleted(lesson.course.id, ch.id) }"
          >
            <Icon v-if="learning.isChapterCompleted(lesson.course.id, ch.id)" name="check" :size="11" :stroke="3" />
            <template v-else>{{ i + 1 }}</template>
          </span>
          <span class="lesson__toc-title">{{ ch.title }}</span>
        </router-link>
      </nav>
    </aside>

    <div class="lesson__main">
      <div class="lesson__crumbs">
        <router-link to="/courses">课程</router-link>
        <Icon name="arrowRight" :size="12" />
        <router-link :to="`/courses/${lesson.course.id}`">{{ lesson.course.title }}</router-link>
        <Icon name="arrowRight" :size="12" />
        <span>第 {{ chapterIndex + 1 }} 章</span>
        <span class="lesson__crumbs-meta"><Icon name="clock" :size="12" /> 约 {{ lesson.chapter.minutes }} 分钟</span>
      </div>

      <article class="lesson__content card fade-up" :key="lesson.chapter.id">
        <MarkdownRenderer :source="lesson.chapter.content" />

        <div class="lesson__complete" :class="{ 'lesson__complete--done': completed }">
          <div>
            <strong>{{ completed ? '本章已完成' : '学完了吗？' }}</strong>
            <span>{{ completed ? '点击可取消完成状态' : '标记完成后将同步到知识图谱中对应知识点的学习状态' }}</span>
            <span v-if="synced.length" class="lesson__synced">
              已更新知识状态：
              <router-link v-for="k in synced" :key="k.knowledge.id" :to="`/knowledge/${k.knowledge.id}`">{{ k.knowledge.name }} {{ Math.round(k.mastery * 100) }}%</router-link>
            </span>
          </div>
          <button
            class="btn"
            :class="completed ? 'btn--ghost' : 'btn--primary'"
            @click="toggleCompleted"
          >
            <Icon :name="completed ? 'check' : 'check'" :size="16" :stroke="3" />
            {{ completed ? '已完成本章' : '完成本章学习' }}
          </button>
        </div>
      </article>

      <section class="lesson__notes card">
        <div class="panel-head">
          <div><h3><Icon name="note" :size="17" /> 本章笔记</h3><p>记录你的理解、疑问或延伸思考，保存在本地浏览器中</p></div>
          <span v-if="noteSaved" class="tag tag--success"><Icon name="check" :size="12" :stroke="3" /> 已保存</span>
        </div>
        <textarea
          v-model="noteDraft"
          class="textarea"
          rows="5"
          placeholder="写点什么…"
        ></textarea>
        <div class="lesson__notes-actions">
          <button class="btn btn--primary" @click="saveNote">保存笔记</button>
        </div>
      </section>

      <nav class="lesson__pager">
        <router-link
          v-if="lesson.prev"
          :to="`/courses/${lesson.course.id}/${lesson.prev.id}`"
          class="lesson__pager-link card card--hover"
        >
          <span class="lesson__pager-dir"><Icon name="arrowLeft" :size="13" /> 上一章</span>
          <span class="lesson__pager-title">{{ lesson.prev.title }}</span>
        </router-link>
        <span v-else></span>
        <router-link
          v-if="lesson.next"
          :to="`/courses/${lesson.course.id}/${lesson.next.id}`"
          class="lesson__pager-link lesson__pager-link--next card card--hover"
        >
          <span class="lesson__pager-dir">下一章 <Icon name="arrowRight" :size="13" /></span>
          <span class="lesson__pager-title">{{ lesson.next.title }}</span>
        </router-link>
        <router-link v-else to="/quiz" class="lesson__pager-link lesson__pager-link--next lesson__pager-link--final card">
          <span class="lesson__pager-dir">课程完结 <Icon name="trophy" :size="13" /></span>
          <span class="lesson__pager-title">去做课程测验</span>
        </router-link>
      </nav>
    </div>
  </div>

  <div class="page" v-else>
    <div class="empty-state card">
      <div class="empty-state__icon"><Icon name="alert" :size="30" /></div>
      <h2>未找到该章节</h2>
      <router-link to="/courses" class="btn btn--primary">返回课程列表</router-link>
    </div>
  </div>
</template>

<style scoped>
.lesson {
  display: flex;
  gap: 24px;
  max-width: 1280px;
  margin: 0 auto;
  padding: 32px 36px 72px;
  align-items: flex-start;
}

.lesson__toc {
  width: 272px;
  min-width: 272px;
  padding: 18px 14px;
  position: sticky;
  top: 24px;
  max-height: calc(100vh - 48px);
  overflow-y: auto;
}

.lesson__course {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 4px 6px 12px;
  color: var(--text);
}

.lesson__course-icon {
  display: grid;
  place-items: center;
  width: 40px;
  height: 40px;
  flex: 0 0 40px;
  border-radius: 12px;
  background: var(--surface-3);
  font-size: 20px;
}

.lesson__course-text {
  display: flex;
  flex-direction: column;
  min-width: 0;
  line-height: 1.35;
}

.lesson__course-text strong {
  font-size: 14px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.lesson__course-text small {
  color: var(--text-3);
  font-size: 11.5px;
}

.lesson__course-progress {
  margin: 0 6px 14px;
}

.lesson__toc-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  padding: 9px 10px;
  border-radius: 11px;
  font-size: 13.5px;
  color: var(--text-2);
  line-height: 1.45;
  transition: background var(--t-fast);
}

.lesson__toc-item:hover {
  background: var(--surface-2);
}

.lesson__toc-item--active {
  background: var(--primary-soft);
  color: var(--primary-deep);
  font-weight: 600;
}

.lesson__toc-dot {
  width: 22px;
  height: 22px;
  min-width: 22px;
  margin-top: 1px;
  display: grid;
  place-items: center;
  border-radius: 7px;
  background: var(--surface-3);
  font-size: 11px;
  font-weight: 700;
  color: var(--text-2);
}

.lesson__toc-item--active .lesson__toc-dot {
  background: var(--primary);
  color: #fff;
}

.lesson__toc-dot--done {
  background: var(--success) !important;
  color: #fff !important;
}

.lesson__main {
  flex: 1;
  min-width: 0;
}

.lesson__crumbs {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 14px;
  color: var(--text-3);
  font-size: 13px;
}

.lesson__crumbs a {
  color: var(--text-2);
  font-weight: 500;
}

.lesson__crumbs a:hover {
  color: var(--primary);
}

.lesson__crumbs-meta {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  margin-left: auto;
}

.lesson__content {
  padding: 40px 48px;
  border-radius: var(--radius-lg);
}

.lesson__complete {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  margin-top: 40px;
  padding: 18px 22px;
  border-radius: 16px;
  background: var(--surface-2);
  border: 1px solid var(--border);
}

.lesson__synced { display: flex; flex-wrap: wrap; gap: 8px; margin-top: 4px; }
.lesson__synced a { font-weight: 600; }

.lesson__complete--done {
  background: var(--success-soft);
  border-color: #bfeedd;
}

.lesson__complete div {
  display: flex;
  flex-direction: column;
}

.lesson__complete strong {
  font-size: 15px;
}

.lesson__complete span {
  color: var(--text-3);
  font-size: 13px;
}

.lesson__notes {
  margin-top: 20px;
  padding: 24px 26px;
}

.lesson__notes h3 {
  display: flex;
  align-items: center;
  gap: 8px;
}

.lesson__notes-actions {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 12px;
  margin-top: 12px;
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
  gap: 4px;
  padding: 16px 20px;
  max-width: 48%;
  color: var(--text);
}

.lesson__pager-link--next {
  text-align: right;
  margin-left: auto;
}

.lesson__pager-link--final {
  color: #fff;
  border: none;
  background: linear-gradient(135deg, var(--primary), var(--primary-strong));
  box-shadow: var(--shadow-primary);
}

.lesson__pager-dir {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--text-3);
}

.lesson__pager-link--next .lesson__pager-dir {
  justify-content: flex-end;
}

.lesson__pager-link--final .lesson__pager-dir {
  color: rgba(255, 255, 255, 0.8);
}

.lesson__pager-title {
  font-size: 14.5px;
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

  .lesson__complete {
    flex-direction: column;
    align-items: stretch;
  }

  .lesson__crumbs-meta {
    display: none;
  }
}
</style>
