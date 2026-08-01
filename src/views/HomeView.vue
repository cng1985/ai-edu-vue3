<script setup>
import { computed } from 'vue'
import { courses, getCourse, totalChapterCount } from '../data/courses'
import { communityPosts, leaderboard, hotTopics, dailyQuote } from '../data/community'
import { useLearningStore } from '../stores/learning'
import { useAuthStore } from '../stores/auth'
import CourseCard from '../components/CourseCard.vue'
import ProgressRing from '../components/ProgressRing.vue'

const learning = useLearningStore()
const auth = useAuthStore()

const greeting = computed(() => {
  const hour = new Date().getHours()
  if (hour < 6) return '夜深了'
  if (hour < 12) return '早上好'
  if (hour < 18) return '下午好'
  return '晚上好'
})

const displayName = computed(() => {
  if (!auth.user) return '探索者'
  return auth.isGuest ? '游客' : auth.user.nickname
})

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
  {
    label: '已完成章节',
    value: learning.completedCount,
    suffix: `/ ${totalChapterCount}`,
    icon: '⚡',
    tone: 'indigo'
  },
  { label: '学习笔记', value: learning.noteCount, suffix: '篇', icon: '🗒️', tone: 'cyan' },
  {
    label: '测验平均分',
    value: learning.quizAverageScore === null ? '—' : learning.quizAverageScore,
    suffix: learning.quizAverageScore === null ? '' : '分',
    icon: '🏆',
    tone: 'amber'
  },
  { label: '今日共学', value: 1284, suffix: '人在线', icon: '🔥', tone: 'rose' }
])
</script>

<template>
  <div class="page home">
    <!-- ============ Hero：深色前卫风 ============ -->
    <section class="hero">
      <div class="hero__glow hero__glow--a"></div>
      <div class="hero__glow hero__glow--b"></div>
      <div class="hero__glow hero__glow--c"></div>
      <div class="hero__grid"></div>

      <div class="hero__inner">
        <div class="hero__text">
          <div class="hero__badge">
            <span class="hero__badge-dot"></span>
            AI 学习社区 · 1,284 位学习者正在共学
          </div>
          <h1 class="hero__title">
            {{ greeting }}，{{ displayName }}
            <span class="hero__title-accent">从写代码，到调度智能</span>
          </h1>
          <p class="hero__desc">
            提示词工程 → RAG 实战 → AI 原生架构，跟随课程循序渐进，
            和社区伙伴一起打卡、讨论、分享笔记，完成认知跃迁。
          </p>
          <div class="hero__actions">
            <router-link
              v-if="continueTarget"
              :to="`/courses/${continueTarget.course.id}/${continueTarget.chapter.id}`"
              class="hero-btn hero-btn--primary"
            >
              ▶ 继续学习：{{ continueTarget.chapter.title }}
            </router-link>
            <router-link v-else to="/courses" class="hero-btn hero-btn--primary">
              🚀 开始学习之旅
            </router-link>
            <router-link to="/chat" class="hero-btn hero-btn--glass">💬 问问 AI 助手</router-link>
          </div>
        </div>

        <div class="hero__panel">
          <ProgressRing :percent="learning.overallProgress" :size="118" :stroke="9" color="#818cf8" />
          <div class="hero__panel-label">总体学习进度</div>
          <div class="hero__panel-sub">
            {{ learning.completedCount }} / {{ totalChapterCount }} 章节
          </div>
        </div>
      </div>
    </section>

    <!-- ============ 数据看板 ============ -->
    <section class="stats">
      <div v-for="s in stats" :key="s.label" class="stat" :class="`stat--${s.tone}`">
        <span class="stat__icon">{{ s.icon }}</span>
        <div class="stat__body">
          <div class="stat__value">
            {{ s.value }}
            <small>{{ s.suffix }}</small>
          </div>
          <div class="stat__label">{{ s.label }}</div>
        </div>
      </div>
    </section>

    <!-- ============ 主体两栏 ============ -->
    <div class="home__columns">
      <div class="home__main">
        <!-- 课程体系 -->
        <section class="block">
          <div class="block__head">
            <h2><span class="block__marker"></span>课程体系</h2>
            <router-link to="/courses" class="block__more">查看全部 →</router-link>
          </div>
          <div class="course-grid">
            <CourseCard v-for="course in courses" :key="course.id" :course="course" />
          </div>
        </section>

        <!-- 社区动态 -->
        <section class="block">
          <div class="block__head">
            <h2><span class="block__marker block__marker--cyan"></span>社区动态</h2>
            <span class="block__hint">来自共学伙伴的实战分享</span>
          </div>

          <router-link to="/chat" class="composer card">
            <span
              class="composer__avatar"
              :style="{ background: auth.user?.avatarColor || '#6366f1' }"
            >
              {{ auth.user?.avatar || '你' }}
            </span>
            <span class="composer__placeholder">分享你的学习心得，或者向 AI 助手提问…</span>
            <span class="composer__btn">✨ 发布</span>
          </router-link>

          <article v-for="post in communityPosts" :key="post.id" class="post card">
            <div class="post__head">
              <span class="post__avatar" :style="{ background: post.avatarColor }">
                {{ post.avatar }}
              </span>
              <div class="post__who">
                <div class="post__author">
                  {{ post.author }}
                  <em class="post__level">{{ post.level }}</em>
                  <em v-if="post.hot" class="post__hot">🔥 热议</em>
                </div>
                <div class="post__time">{{ post.time }}</div>
              </div>
            </div>
            <h3 class="post__title">{{ post.title }}</h3>
            <p class="post__excerpt">{{ post.excerpt }}</p>
            <div class="post__foot">
              <div class="post__tags">
                <span v-for="tag in post.tags" :key="tag" class="post__tag"># {{ tag }}</span>
              </div>
              <div class="post__meta">
                <span>👍 {{ post.likes }}</span>
                <span>💬 {{ post.comments }}</span>
              </div>
            </div>
          </article>
        </section>
      </div>

      <!-- ============ 右侧社区栏 ============ -->
      <aside class="home__rail">
        <section class="rail-card card">
          <div class="rail-card__head">
            <h3>🏅 学霸排行榜</h3>
            <span class="rail-card__sub">本周</span>
          </div>
          <ol class="board">
            <li v-for="(u, i) in leaderboard" :key="u.name" class="board__row">
              <span class="board__rank" :class="`board__rank--${i + 1}`">{{ i + 1 }}</span>
              <span class="board__avatar" :style="{ background: u.avatarColor }">{{ u.avatar }}</span>
              <span class="board__name">{{ u.name }}</span>
              <span class="board__score">{{ u.chapters }} 章 · {{ u.streak }} 天</span>
            </li>
          </ol>
        </section>

        <section class="rail-card card">
          <div class="rail-card__head">
            <h3>📈 热门话题</h3>
            <span class="rail-card__sub">24h</span>
          </div>
          <div class="topics">
            <router-link v-for="t in hotTopics" :key="t.name" to="/courses" class="topic">
              <span class="topic__name"># {{ t.name }}</span>
              <span class="topic__meta">{{ t.posts }} 帖 <em>{{ t.trend }}</em></span>
            </router-link>
          </div>
        </section>

        <section class="quote">
          <div class="quote__mark">"</div>
          <p class="quote__text">{{ dailyQuote.text }}</p>
          <div class="quote__from">—— {{ dailyQuote.from }}</div>
        </section>
      </aside>
    </div>
  </div>
</template>

<style scoped>
.home {
  max-width: 1160px;
}

/* ================= Hero ================= */

.hero {
  position: relative;
  overflow: hidden;
  border-radius: 22px;
  padding: 42px 40px;
  margin-bottom: 22px;
  background:
    radial-gradient(ellipse 90% 130% at 8% -20%, #312e81 0%, transparent 55%),
    linear-gradient(135deg, #0b1020 0%, #131a33 52%, #101426 100%);
  border: 1px solid rgba(129, 140, 248, 0.22);
  box-shadow: 0 24px 60px -22px rgba(30, 27, 75, 0.55);
  color: #e2e8f0;
  isolation: isolate;
}

.hero__glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(70px);
  opacity: 0.55;
  z-index: -1;
  animation: hero-float 9s ease-in-out infinite alternate;
}

.hero__glow--a {
  width: 380px;
  height: 380px;
  top: -180px;
  right: -60px;
  background: radial-gradient(circle, rgba(99, 102, 241, 0.85), transparent 65%);
}

.hero__glow--b {
  width: 300px;
  height: 300px;
  bottom: -170px;
  left: 28%;
  background: radial-gradient(circle, rgba(14, 165, 233, 0.6), transparent 65%);
  animation-delay: -3s;
}

.hero__glow--c {
  width: 240px;
  height: 240px;
  top: -110px;
  left: -80px;
  background: radial-gradient(circle, rgba(217, 70, 239, 0.45), transparent 65%);
  animation-delay: -6s;
}

@keyframes hero-float {
  from {
    transform: translate3d(0, 0, 0) scale(1);
  }
  to {
    transform: translate3d(26px, 22px, 0) scale(1.12);
  }
}

.hero__grid {
  position: absolute;
  inset: 0;
  z-index: -1;
  background-image:
    linear-gradient(rgba(148, 163, 184, 0.07) 1px, transparent 1px),
    linear-gradient(90deg, rgba(148, 163, 184, 0.07) 1px, transparent 1px);
  background-size: 34px 34px;
  mask-image: radial-gradient(ellipse 75% 90% at 50% 0%, #000 30%, transparent 100%);
}

.hero__inner {
  display: flex;
  align-items: center;
  gap: 36px;
}

.hero__text {
  flex: 1;
  min-width: 0;
}

.hero__badge {
  display: inline-flex;
  align-items: center;
  gap: 8px;
  padding: 5px 14px;
  margin-bottom: 18px;
  border-radius: 999px;
  border: 1px solid rgba(129, 140, 248, 0.35);
  background: rgba(99, 102, 241, 0.14);
  color: #c7d2fe;
  font-size: 12.5px;
  font-weight: 600;
  letter-spacing: 0.02em;
  backdrop-filter: blur(6px);
}

.hero__badge-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background: #34d399;
  box-shadow: 0 0 0 0 rgba(52, 211, 153, 0.6);
  animation: pulse-dot 2s ease-out infinite;
}

@keyframes pulse-dot {
  from {
    box-shadow: 0 0 0 0 rgba(52, 211, 153, 0.6);
  }
  to {
    box-shadow: 0 0 0 9px rgba(52, 211, 153, 0);
  }
}

.hero__title {
  margin: 0 0 12px;
  font-size: 30px;
  line-height: 1.3;
  letter-spacing: -0.02em;
  color: #f8fafc;
}

.hero__title-accent {
  display: block;
  font-size: 22px;
  margin-top: 4px;
  background: linear-gradient(92deg, #818cf8 0%, #38bdf8 45%, #e879f9 100%);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}

.hero__desc {
  margin: 0 0 24px;
  max-width: 540px;
  color: #94a3b8;
  font-size: 14.5px;
}

.hero__actions {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.hero-btn {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  padding: 11px 22px;
  border-radius: 12px;
  font-size: 14px;
  font-weight: 600;
  transition: transform 0.15s ease, box-shadow 0.15s ease, background 0.15s ease;
  white-space: nowrap;
}

.hero-btn--primary {
  background: linear-gradient(92deg, #6366f1, #8b5cf6);
  color: #fff;
  box-shadow: 0 10px 26px -8px rgba(99, 102, 241, 0.65);
}

.hero-btn--primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 14px 32px -8px rgba(99, 102, 241, 0.8);
}

.hero-btn--glass {
  background: rgba(148, 163, 184, 0.1);
  border: 1px solid rgba(148, 163, 184, 0.28);
  color: #e2e8f0;
  backdrop-filter: blur(6px);
}

.hero-btn--glass:hover {
  background: rgba(148, 163, 184, 0.2);
  transform: translateY(-2px);
}

.hero__panel {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 6px;
  padding: 26px 34px;
  border-radius: 18px;
  background: rgba(15, 23, 42, 0.45);
  border: 1px solid rgba(129, 140, 248, 0.25);
  backdrop-filter: blur(10px);
}

.hero__panel :deep(.progress-ring__text) {
  fill: #f8fafc;
}

.hero__panel-label {
  font-size: 13px;
  font-weight: 600;
  color: #cbd5e1;
}

.hero__panel-sub {
  font-size: 12px;
  color: #64748b;
}

/* ================= 数据看板 ================= */

.stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 14px;
  margin-bottom: 30px;
}

.stat {
  display: flex;
  align-items: center;
  gap: 13px;
  padding: 16px 18px;
  border-radius: 16px;
  background: var(--surface);
  border: 1px solid var(--border);
  box-shadow: var(--shadow);
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.stat:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(15, 23, 42, 0.08), 0 18px 42px rgba(15, 23, 42, 0.07);
}

.stat__icon {
  width: 44px;
  height: 44px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  border-radius: 13px;
  flex-shrink: 0;
}

.stat--indigo .stat__icon { background: linear-gradient(135deg, #eef2ff, #e0e7ff); }
.stat--cyan .stat__icon { background: linear-gradient(135deg, #ecfeff, #cffafe); }
.stat--amber .stat__icon { background: linear-gradient(135deg, #fffbeb, #fef3c7); }
.stat--rose .stat__icon { background: linear-gradient(135deg, #fff1f2, #ffe4e6); }

.stat__value {
  font-size: 21px;
  font-weight: 800;
  letter-spacing: -0.02em;
  line-height: 1.2;
}

.stat__value small {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-3);
  margin-left: 2px;
}

.stat__label {
  font-size: 12.5px;
  color: var(--text-2);
}

/* ================= 两栏布局 ================= */

.home__columns {
  display: grid;
  grid-template-columns: minmax(0, 1fr) 300px;
  gap: 22px;
  align-items: start;
}

.home__main {
  min-width: 0;
}

.block {
  margin-bottom: 34px;
}

.block__head {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 16px;
}

.block__head h2 {
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 0;
  font-size: 19px;
  letter-spacing: -0.01em;
}

.block__marker {
  width: 5px;
  height: 20px;
  border-radius: 3px;
  background: linear-gradient(180deg, #6366f1, #a78bfa);
}

.block__marker--cyan {
  background: linear-gradient(180deg, #0ea5e9, #22d3ee);
}

.block__more {
  font-size: 13.5px;
  font-weight: 600;
}

.block__hint {
  font-size: 12.5px;
  color: var(--text-3);
}

.course-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(270px, 1fr));
  gap: 16px;
}

/* ================= 社区动态 ================= */

.composer {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 13px 16px;
  margin-bottom: 14px;
  color: var(--text-3);
  transition: border-color 0.15s ease, box-shadow 0.15s ease;
}

.composer:hover {
  border-color: var(--primary);
  box-shadow: 0 4px 14px rgba(99, 102, 241, 0.14);
}

.composer__avatar {
  width: 34px;
  height: 34px;
  min-width: 34px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: #fff;
  font-size: 13px;
  font-weight: 700;
}

.composer__placeholder {
  flex: 1;
  font-size: 13.5px;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.composer__btn {
  padding: 7px 16px;
  border-radius: 999px;
  background: linear-gradient(92deg, #6366f1, #8b5cf6);
  color: #fff;
  font-size: 12.5px;
  font-weight: 700;
  flex-shrink: 0;
}

.post {
  padding: 18px 20px;
  margin-bottom: 14px;
  transition: transform 0.15s ease, box-shadow 0.15s ease;
}

.post:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 16px rgba(15, 23, 42, 0.08), 0 18px 42px rgba(15, 23, 42, 0.07);
}

.post__head {
  display: flex;
  align-items: center;
  gap: 11px;
  margin-bottom: 12px;
}

.post__avatar {
  width: 38px;
  height: 38px;
  min-width: 38px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: #fff;
  font-size: 14px;
  font-weight: 700;
}

.post__author {
  display: flex;
  align-items: center;
  gap: 7px;
  font-size: 14px;
  font-weight: 700;
}

.post__level {
  font-style: normal;
  font-size: 10.5px;
  font-weight: 700;
  padding: 1px 7px;
  border-radius: 999px;
  background: var(--primary-soft);
  color: var(--primary-strong);
}

.post__hot {
  font-style: normal;
  font-size: 10.5px;
  font-weight: 700;
  padding: 1px 7px;
  border-radius: 999px;
  background: #fff1f2;
  color: #e11d48;
}

.post__time {
  font-size: 12px;
  color: var(--text-3);
}

.post__title {
  margin: 0 0 7px;
  font-size: 15.5px;
  font-weight: 700;
  letter-spacing: -0.01em;
}

.post__excerpt {
  margin: 0 0 14px;
  font-size: 13.5px;
  line-height: 1.75;
  color: var(--text-2);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.post__foot {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.post__tags {
  display: flex;
  gap: 7px;
  flex-wrap: wrap;
}

.post__tag {
  font-size: 12px;
  font-weight: 600;
  padding: 2px 10px;
  border-radius: 999px;
  background: var(--surface-2);
  border: 1px solid var(--border);
  color: var(--text-2);
}

.post__meta {
  display: flex;
  gap: 14px;
  font-size: 12.5px;
  color: var(--text-3);
  white-space: nowrap;
}

/* ================= 右侧栏 ================= */

.home__rail {
  display: flex;
  flex-direction: column;
  gap: 16px;
  position: sticky;
  top: 24px;
}

.rail-card {
  padding: 18px;
}

.rail-card__head {
  display: flex;
  justify-content: space-between;
  align-items: baseline;
  margin-bottom: 14px;
}

.rail-card__head h3 {
  margin: 0;
  font-size: 15px;
}

.rail-card__sub {
  font-size: 11.5px;
  font-weight: 700;
  padding: 1px 8px;
  border-radius: 999px;
  background: var(--surface-2);
  color: var(--text-3);
}

.board {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.board__row {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 7px 8px;
  border-radius: 10px;
  transition: background 0.15s ease;
}

.board__row:hover {
  background: var(--surface-2);
}

.board__rank {
  width: 20px;
  text-align: center;
  font-size: 13px;
  font-weight: 800;
  color: var(--text-3);
  font-style: italic;
}

.board__rank--1 { color: #f59e0b; }
.board__rank--2 { color: #94a3b8; }
.board__rank--3 { color: #d97706; }

.board__avatar {
  width: 28px;
  height: 28px;
  min-width: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  color: #fff;
  font-size: 11.5px;
  font-weight: 700;
}

.board__name {
  flex: 1;
  min-width: 0;
  font-size: 13px;
  font-weight: 600;
  overflow: hidden;
  white-space: nowrap;
  text-overflow: ellipsis;
}

.board__score {
  font-size: 11.5px;
  color: var(--text-3);
  white-space: nowrap;
}

.topics {
  display: flex;
  flex-direction: column;
  gap: 3px;
}

.topic {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 10px;
  padding: 8px 9px;
  border-radius: 10px;
  color: var(--text);
  transition: background 0.15s ease;
}

.topic:hover {
  background: var(--primary-soft);
}

.topic__name {
  font-size: 13px;
  font-weight: 600;
}

.topic__meta {
  font-size: 11.5px;
  color: var(--text-3);
  white-space: nowrap;
}

.topic__meta em {
  font-style: normal;
  font-weight: 700;
  color: var(--success);
  margin-left: 4px;
}

.quote {
  position: relative;
  padding: 22px 20px 18px;
  border-radius: 16px;
  background: linear-gradient(135deg, #1e1b4b 0%, #172554 100%);
  border: 1px solid rgba(129, 140, 248, 0.3);
  color: #e0e7ff;
  overflow: hidden;
}

.quote__mark {
  position: absolute;
  top: -18px;
  right: 10px;
  font-size: 110px;
  font-family: Georgia, serif;
  color: rgba(129, 140, 248, 0.16);
  pointer-events: none;
}

.quote__text {
  margin: 0 0 12px;
  font-size: 14px;
  line-height: 1.8;
  font-weight: 500;
}

.quote__from {
  font-size: 11.5px;
  color: #818cf8;
}

/* ================= 响应式 ================= */

@media (max-width: 1024px) {
  .home__columns {
    grid-template-columns: 1fr;
  }

  .home__rail {
    position: static;
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  }
}

@media (max-width: 860px) {
  .stats {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 720px) {
  .hero {
    padding: 30px 24px;
  }

  .hero__inner {
    flex-direction: column-reverse;
    text-align: center;
  }

  .hero__title {
    font-size: 25px;
  }

  .hero__title-accent {
    font-size: 19px;
  }

  .hero__actions {
    justify-content: center;
  }

  .hero__desc {
    margin-inline: auto;
  }

  .stats {
    grid-template-columns: 1fr;
  }
}
</style>
