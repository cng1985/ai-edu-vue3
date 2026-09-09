<script setup>
import { computed, ref } from 'vue'
import { useGrowthStore } from '../stores/growth'
import Icon from '../components/Icon.vue'

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
    <header class="page-header page-header--split">
      <div>
        <span class="eyebrow">AI 查漏补缺</span>
        <h1>复习与补强中心</h1>
        <p>根据知识图谱、快测成绩和遗忘曲线，优先处理真正影响目标的知识缺口。</p>
      </div>
      <div class="review-score card">
        <strong class="num">{{ growth.dueReviewCount }}</strong>
        <span>今日待复习</span>
      </div>
    </header>

    <div v-if="notice" class="toast fade-up">{{ notice }}</div>

    <section v-if="!growth.hasGoal" class="empty-state card">
      <div class="empty-state__icon"><Icon name="compass" :size="32" /></div>
      <h2>先建立目标，才能识别知识缺口</h2>
      <p>系统会基于能力图谱和快测结果，为你找出最值得补强的知识点。</p>
      <router-link to="/career" class="btn btn--primary">创建职业目标</router-link>
    </section>

    <template v-else>
      <section class="summary-grid stagger">
        <div class="stat card">
          <span class="stat__icon stat__icon--danger"><Icon name="alert" :size="20" /></span>
          <div class="stat__body"><strong class="stat__value">{{ growth.weakPoints.length }}</strong><small class="stat__label">待补知识点</small></div>
        </div>
        <div class="stat card">
          <span class="stat__icon stat__icon--warning"><Icon name="refresh" :size="20" /></span>
          <div class="stat__body"><strong class="stat__value">{{ growth.dueReviewCount }}</strong><small class="stat__label">今日复习项</small></div>
        </div>
        <div class="stat card">
          <span class="stat__icon stat__icon--success"><Icon name="check" :size="20" /></span>
          <div class="stat__body"><strong class="stat__value">{{ growth.completedUnitCount }}</strong><small class="stat__label">已学习单元</small></div>
        </div>
      </section>

      <nav class="segmented tabs">
        <button :class="{ active: activeTab === 'gaps' }" @click="activeTab = 'gaps'">
          知识缺口 <em>{{ growth.weakPoints.length }}</em>
        </button>
        <button :class="{ active: activeTab === 'review' }" @click="activeTab = 'review'">
          遗忘曲线复习 <em>{{ reviewItems.length }}</em>
        </button>
      </nav>

      <section v-if="activeTab === 'gaps'" class="card panel fade-up">
        <div class="panel-head">
          <div><h2>补强优先级</h2><p>未学习的关键路径节点优先，低分知识点需重新学习并测评。</p></div>
          <span class="tag tag--ai"><Icon name="sparkles" :size="12" /> AI 已排序</span>
        </div>
        <div v-if="!growth.weakPoints.length" class="empty-state">
          <div class="empty-state__icon"><Icon name="trophy" :size="30" /></div>
          <h3>当前没有明显知识缺口</h3>
          <p>继续保持，并按复习计划巩固长期记忆。</p>
        </div>
        <div v-for="(points, domain) in gapGroups" :key="domain" class="gap-group">
          <h3>{{ domain }} <span class="tag tag--neutral">{{ points.length }} 项</span></h3>
          <router-link
            v-for="point in points"
            :key="point.id"
            :to="`/micro/${point.unitId}`"
            class="gap-row"
          >
            <span class="tag" :class="point.completed && point.score < 70 ? 'tag--danger' : 'tag--neutral'">
              {{ point.completed ? '需补强' : '未学习' }}
            </span>
            <div class="gap-row__body">
              <strong>{{ point.name }}</strong>
              <small>{{ point.reason }} · 当前掌握度 {{ point.mastery }}%</small>
            </div>
            <div class="progress progress--thin mini-track"><i :style="{ width: point.mastery + '%' }"></i></div>
            <b>{{ point.completed ? '重新学习' : '开始学习' }} <Icon name="arrowRight" :size="13" /></b>
          </router-link>
        </div>
      </section>

      <section v-else class="card panel fade-up">
        <div class="panel-head">
          <div><h2>智能复习日程</h2><p>完成学习后按 1 / 3 / 7 / 15 / 30 天安排间隔复习。</p></div>
          <span class="tag tag--info">艾宾浩斯</span>
        </div>
        <div v-if="!reviewItems.length" class="empty-state">
          <div class="empty-state__icon"><Icon name="calendar" :size="30" /></div>
          <h3>暂无复习任务</h3>
          <p>完成微单元后，系统会自动生成首个复习日程。</p>
        </div>
        <article v-for="item in reviewItems" :key="item.id" class="review-item" :class="{ 'review-item--due': item.due, 'review-item--open': openCard === item.id }">
          <button class="review-item__head" @click="openCard = openCard === item.id ? null : item.id">
            <span class="calendar">{{ item.due ? '今' : item.stage + 1 }}</span>
            <div class="review-item__body">
              <strong>{{ item.title }}</strong>
              <small>{{ item.competency }} · 上次快测 {{ item.score }} 分</small>
            </div>
            <em>{{ dueLabel(item.dueAt) }}</em>
            <b>{{ openCard === item.id ? '收起' : '打开复习卡' }}</b>
          </button>
          <div v-if="openCard === item.id" class="review-card fade-up">
            <p>{{ item.intro }}</p>
            <ul><li v-for="point in item.summary" :key="point"><Icon name="check" :size="13" :stroke="3" /> {{ point }}</li></ul>
            <div class="review-actions">
              <router-link :to="`/micro/${item.id}`" class="btn btn--ghost">重新学习</router-link>
              <button class="btn btn--primary" @click="finishReview(item.id)"><Icon name="check" :size="15" :stroke="3" /> 我已回忆并掌握 +3</button>
            </div>
          </div>
        </article>
      </section>
    </template>
  </div>
</template>

<style scoped>
.review-page { max-width: 1040px; }
.review-score { display: flex; min-width: 128px; flex-direction: column; align-items: center; padding: 14px 24px; }
.review-score strong { color: var(--primary-strong); font-size: 32px; line-height: 1.1; }
.review-score span { color: var(--text-3); font-size: 12px; }

.summary-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 14px; margin-bottom: 20px; }
.tabs { display: flex; width: 100%; margin-bottom: 14px; }
.tabs button { flex: 1; display: inline-flex; align-items: center; justify-content: center; gap: 6px; }
.tabs em { font-style: normal; padding: 1px 8px; border-radius: 99px; background: var(--surface); color: var(--text-2); font-size: 11.5px; box-shadow: var(--shadow-xs); }
.tabs button.active em { background: var(--primary-soft); color: var(--primary-strong); }

.gap-group { margin-top: 22px; }
.gap-group h3 { display: flex; align-items: center; gap: 8px; margin: 0 0 8px; font-size: 14px; }
.gap-row {
  display: grid;
  grid-template-columns: 68px minmax(150px, 1fr) 140px 96px;
  align-items: center;
  gap: 14px;
  padding: 12px 10px;
  border-radius: 12px;
  color: inherit;
  transition: background var(--t-fast);
}
.gap-row + .gap-row { border-top: 1px solid var(--border); border-radius: 0 0 12px 12px; }
.gap-row:hover { background: var(--surface-2); }
.gap-row__body { display: flex; flex-direction: column; min-width: 0; }
.gap-row__body strong { font-size: 14px; }
.gap-row small { color: var(--text-3); font-size: 12px; }
.gap-row b { display: inline-flex; align-items: center; justify-content: flex-end; gap: 3px; color: var(--primary); font-size: 12.5px; }
.mini-track i { background: linear-gradient(90deg, var(--warning), #fcd34d); }

.review-item { border-radius: 14px; transition: background var(--t-fast); }
.review-item + .review-item { border-top: 1px solid var(--border); }
.review-item--open { background: var(--surface-2); }
.review-item__head { display: flex; align-items: center; gap: 14px; width: 100%; padding: 14px 10px; border: 0; background: transparent; color: inherit; font: inherit; text-align: left; cursor: pointer; }
.calendar { display: grid; place-items: center; width: 40px; height: 40px; flex: 0 0 40px; border-radius: 12px; background: var(--surface-3); color: var(--text-2); font-family: 'Sora', sans-serif; font-weight: 800; }
.review-item--due .calendar { background: linear-gradient(135deg, var(--primary), var(--primary-strong)); color: white; box-shadow: var(--shadow-primary); }
.review-item__body { display: flex; flex: 1; min-width: 0; flex-direction: column; }
.review-item__body strong { font-size: 14px; }
.review-item__body small { color: var(--text-3); font-size: 12px; }
.review-item__head em { color: var(--text-3); font-size: 12px; font-style: normal; }
.review-item--due .review-item__head em { color: var(--danger); font-weight: 600; }
.review-item__head b { min-width: 84px; color: var(--primary); font-size: 12.5px; text-align: right; }
.review-card { margin: 0 10px 14px 64px; padding: 18px 20px; border-radius: 14px; background: var(--surface); border: 1px solid var(--border); }
.review-card p { margin: 0 0 10px; color: var(--text-2); }
.review-card ul { margin: 0 0 14px; padding: 0; list-style: none; display: flex; flex-direction: column; gap: 6px; }
.review-card li { display: flex; align-items: center; gap: 8px; font-size: 13.5px; }
.review-card li .icon { color: var(--success); }
.review-actions { display: flex; justify-content: flex-end; gap: 8px; }

@media (max-width: 760px) {
  .summary-grid { grid-template-columns: 1fr; }
  .gap-row { grid-template-columns: 62px 1fr auto; }
  .mini-track { display: none; }
  .review-item__head em { display: none; }
  .review-card { margin-left: 10px; }
}
</style>
