<script setup>
import { computed, ref } from 'vue'
import { useGrowthStore } from '../stores/growth'

const growth = useGrowthStore()
const activeTab = ref('gaps')
const openCard = ref(null)
const notice = ref('')

const gapGroups = computed(() => {
  const groups = {}
  growth.weakPoints.forEach((point) => {
    if (!groups[point.domainName]) groups[point.domainName] = []
    groups[point.domainName].push(point)
  })
  return groups
})

const reviewItems = computed(() => growth.reviewQueue)
const dueLabel = (timestamp) => {
  const days = Math.ceil((timestamp - Date.now()) / 86400000)
  if (days <= 0) return '今天到期'
  if (days === 1) return '明天复习'
  return `${days} 天后复习`
}

function finishReview(id) {
  const awarded = growth.completeReview(id)
  notice.value = awarded ? '复习完成，获得 3 积分，已安排下一次复习。' : '今天已复习过该知识点。'
  openCard.value = null
  window.setTimeout(() => { notice.value = '' }, 2600)
}
</script>

<template>
  <div class="page review-page">
    <header class="page-header review-header">
      <div>
        <span class="eyebrow">AI 查漏补缺</span>
        <h1>复习与补强中心</h1>
        <p>根据知识图谱、快测成绩和遗忘曲线，优先处理真正影响目标的知识缺口。</p>
      </div>
      <div class="review-score card">
        <strong>{{ growth.dueReviewCount }}</strong>
        <span>今日待复习</span>
      </div>
    </header>

    <div v-if="notice" class="notice fade-up">✓ {{ notice }}</div>

    <section v-if="!growth.hasGoal" class="empty card">
      <span>🧭</span>
      <h2>先建立目标，才能识别知识缺口</h2>
      <router-link to="/career" class="btn btn--primary">创建职业目标</router-link>
    </section>

    <template v-else>
      <section class="summary-grid">
        <div class="summary-card card">
          <span class="summary-icon danger">!</span>
          <div><strong>{{ growth.weakPoints.length }}</strong><small>待补知识点</small></div>
        </div>
        <div class="summary-card card">
          <span class="summary-icon warning">↻</span>
          <div><strong>{{ growth.dueReviewCount }}</strong><small>今日复习项</small></div>
        </div>
        <div class="summary-card card">
          <span class="summary-icon success">✓</span>
          <div><strong>{{ growth.completedUnitCount }}</strong><small>已学习单元</small></div>
        </div>
      </section>

      <nav class="tabs">
        <button :class="{ active: activeTab === 'gaps' }" @click="activeTab = 'gaps'">
          知识缺口 <span>{{ growth.weakPoints.length }}</span>
        </button>
        <button :class="{ active: activeTab === 'review' }" @click="activeTab = 'review'">
          遗忘曲线复习 <span>{{ reviewItems.length }}</span>
        </button>
      </nav>

      <section v-if="activeTab === 'gaps'" class="content-card card">
        <div class="section-head">
          <div><h2>补强优先级</h2><p>未学习的关键路径节点优先，低分知识点需重新学习并测评。</p></div>
          <span class="ai-tag">AI 已排序</span>
        </div>
        <div v-if="!growth.weakPoints.length" class="all-clear">
          <span>🎉</span><h3>当前没有明显知识缺口</h3><p>继续保持，并按复习计划巩固长期记忆。</p>
        </div>
        <div v-for="(points, domain) in gapGroups" :key="domain" class="gap-group">
          <h3>{{ domain }} <span>{{ points.length }} 项</span></h3>
          <router-link
            v-for="point in points"
            :key="point.id"
            :to="`/micro/${point.unitId}`"
            class="gap-row"
          >
            <span class="priority" :class="{ high: point.completed && point.score < 70 }">
              {{ point.completed ? '需补强' : '未学习' }}
            </span>
            <div>
              <strong>{{ point.name }}</strong>
              <small>{{ point.reason }} · 当前掌握度 {{ point.mastery }}%</small>
            </div>
            <div class="mini-track"><i :style="{ width: point.mastery + '%' }"></i></div>
            <b>{{ point.completed ? '重新学习' : '开始学习' }} →</b>
          </router-link>
        </div>
      </section>

      <section v-else class="content-card card">
        <div class="section-head">
          <div><h2>智能复习日程</h2><p>完成学习后按 1 / 3 / 7 / 15 / 30 天安排间隔复习。</p></div>
          <span class="ai-tag">艾宾浩斯</span>
        </div>
        <div v-if="!reviewItems.length" class="all-clear">
          <span>📅</span><h3>暂无复习任务</h3><p>完成微单元后，系统会自动生成首个复习日程。</p>
        </div>
        <article v-for="item in reviewItems" :key="item.id" class="review-item" :class="{ due: item.due }">
          <button class="review-item__head" @click="openCard = openCard === item.id ? null : item.id">
            <span class="calendar">{{ item.due ? '今' : item.stage + 1 }}</span>
            <div>
              <strong>{{ item.title }}</strong>
              <small>{{ item.competency }} · 上次快测 {{ item.score }} 分</small>
            </div>
            <em>{{ dueLabel(item.dueAt) }}</em>
            <b>{{ openCard === item.id ? '收起' : '打开复习卡' }}</b>
          </button>
          <div v-if="openCard === item.id" class="review-card fade-up">
            <p>{{ item.intro }}</p>
            <ul><li v-for="point in item.summary" :key="point">{{ point }}</li></ul>
            <div class="review-actions">
              <router-link :to="`/micro/${item.id}`" class="btn btn--ghost">重新学习</router-link>
              <button class="btn btn--primary" @click="finishReview(item.id)">我已回忆并掌握 +3</button>
            </div>
          </div>
        </article>
      </section>
    </template>
  </div>
</template>

<style scoped>
.review-page { max-width: 1020px; }
.review-header { display: flex; align-items: center; justify-content: space-between; }
.eyebrow { color: var(--primary); font-size: 12px; font-weight: 800; letter-spacing: .08em; }
.review-header h1 { margin-top: 5px; }
.review-score { display: flex; min-width: 120px; flex-direction: column; padding: 15px 22px; text-align: center; box-shadow: none; }
.review-score strong { color: var(--primary); font-size: 30px; line-height: 1.1; }
.review-score span { color: var(--text-3); font-size: 11px; }
.notice { position: fixed; top: 20px; right: 25px; z-index: 100; padding: 11px 16px; border-radius: 9px; background: #065f46; color: white; box-shadow: var(--shadow); font-size: 13px; }
.summary-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 14px; margin-bottom: 18px; }
.summary-card { display: flex; align-items: center; gap: 13px; padding: 17px 20px; box-shadow: none; }
.summary-icon { display: grid; place-items: center; width: 38px; height: 38px; border-radius: 11px; font-size: 18px; font-weight: 800; }
.summary-icon.danger { background: #fef2f2; color: var(--danger); }
.summary-icon.warning { background: #fffbeb; color: var(--warning); }
.summary-icon.success { background: var(--success-soft); color: var(--success); }
.summary-card div { display: flex; flex-direction: column; }
.summary-card strong { font-size: 21px; line-height: 1.2; }
.summary-card small { color: var(--text-3); }
.tabs { display: flex; gap: 5px; margin-bottom: 12px; padding: 4px; border-radius: 10px; background: #e9edf4; }
.tabs button { flex: 1; padding: 9px; border: none; border-radius: 8px; background: transparent; color: var(--text-2); font: inherit; cursor: pointer; }
.tabs button.active { background: white; color: var(--primary); font-weight: 700; box-shadow: 0 1px 3px rgba(15,23,42,.1); }
.tabs span { margin-left: 5px; padding: 1px 6px; border-radius: 99px; background: var(--primary-soft); font-size: 10px; }
.content-card { padding: 24px 27px; box-shadow: none; }
.section-head { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 18px; }
.section-head h2 { margin: 0; font-size: 18px; }
.section-head p { margin: 2px 0 0; color: var(--text-3); font-size: 12px; }
.ai-tag { padding: 3px 9px; border-radius: 99px; background: linear-gradient(110deg, #6366f1, #8b5cf6); color: white; font-size: 10px; font-weight: 700; }
.gap-group { margin-top: 20px; }
.gap-group h3 { margin: 0 0 7px; font-size: 13px; }
.gap-group h3 span { color: var(--text-3); font-size: 10px; }
.gap-row { display: grid; grid-template-columns: 60px minmax(150px, 1fr) 130px 75px; align-items: center; gap: 12px; padding: 11px 7px; border-top: 1px solid var(--border); color: inherit; }
.gap-row:hover { background: var(--surface-2); }
.priority { padding: 2px 6px; border-radius: 6px; background: var(--surface-2); color: var(--text-3); font-size: 10px; text-align: center; }
.priority.high { background: #fef2f2; color: var(--danger); }
.gap-row > div:nth-child(2) { display: flex; flex-direction: column; }
.gap-row small { color: var(--text-3); }
.gap-row b { color: var(--primary); font-size: 11px; text-align: right; }
.mini-track { height: 5px; overflow: hidden; border-radius: 5px; background: var(--border); }
.mini-track i { display: block; height: 100%; background: var(--warning); }
.review-item { border-top: 1px solid var(--border); }
.review-item__head { display: flex; align-items: center; gap: 12px; width: 100%; padding: 14px 4px; border: 0; background: transparent; color: inherit; font: inherit; text-align: left; cursor: pointer; }
.calendar { display: grid; place-items: center; width: 36px; height: 36px; border-radius: 10px; background: var(--surface-2); color: var(--text-2); font-weight: 800; }
.due .calendar { background: var(--primary); color: white; }
.review-item__head div { display: flex; flex: 1; flex-direction: column; }
.review-item__head small { color: var(--text-3); }
.review-item__head em { color: var(--text-3); font-size: 11px; font-style: normal; }
.due .review-item__head em { color: var(--danger); }
.review-item__head b { min-width: 70px; color: var(--primary); font-size: 11px; text-align: right; }
.review-card { margin: 0 5px 14px 48px; padding: 16px 18px; border-radius: 10px; background: var(--primary-soft); }
.review-card p { margin-top: 0; color: var(--text-2); }
.review-card li { margin: 4px 0; font-size: 13px; }
.review-actions { display: flex; justify-content: flex-end; gap: 8px; }
.all-clear, .empty { padding: 40px; text-align: center; }
.all-clear > span, .empty > span { font-size: 36px; }
.all-clear h3, .empty h2 { margin: 7px 0 2px; }
.all-clear p { margin: 0; color: var(--text-3); }
@media (max-width: 700px) {
  .review-header { align-items: flex-start; }
  .summary-grid { grid-template-columns: 1fr; }
  .gap-row { grid-template-columns: 58px 1fr auto; }
  .mini-track { display: none; }
  .review-item__head em { display: none; }
  .review-card { margin-left: 0; }
}
</style>
