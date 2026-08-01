<script setup>
import { computed, ref } from 'vue'
import { useGrowthStore } from '../stores/growth'
import { microUnits } from '../data/careerPath'

const growth = useGrowthStore()
const checkInMessage = ref('')

const calendarDays = computed(() => {
  const days = []
  for (let i = 34; i >= 0; i--) {
    const date = new Date()
    date.setDate(date.getDate() - i)
    const key = new Date(date.getTime() - date.getTimezoneOffset() * 60000).toISOString().slice(0, 10)
    days.push({
      key,
      day: date.getDate(),
      weekday: ['日', '一', '二', '三', '四', '五', '六'][date.getDay()],
      checked: growth.checkIns.includes(key),
      today: i === 0
    })
  }
  return days
})

const levelProgress = computed(() => {
  const { min, next } = growth.level
  if (next === min) return 100
  return Math.round((growth.points - min) / (next - min) * 100)
})

const badgeCatalog = computed(() => [
  { name: '目标启航', icon: '🚀', detail: '建立第一个学习目标', unlocked: growth.badges.includes('目标启航'), progress: growth.hasGoal ? 100 : 0 },
  { name: '第一步', icon: '👣', detail: '完成第一个微单元', unlocked: growth.badges.includes('第一步'), progress: Math.min(100, growth.completedUnitCount * 100) },
  { name: '坚持者', icon: '🔥', detail: '连续学习 7 天', unlocked: growth.badges.includes('坚持者'), progress: Math.min(100, growth.streak / 7 * 100) },
  { name: '复习达人', icon: '🧠', detail: '完成 3 个知识点复习', unlocked: growth.badges.includes('复习达人'), progress: Math.min(100, Object.values(growth.reviewRecords).filter(item => item.lastAt).length / 3 * 100) },
  { name: '碎片收集者', icon: '💎', detail: `完成 ${microUnits.length} 个微单元`, unlocked: growth.completedUnitCount >= microUnits.length, progress: growth.completedUnitCount / microUnits.length * 100 },
  { name: '目标达成者', icon: '🏆', detail: '目标达成度达到 75%', unlocked: growth.achievement >= 75, progress: Math.min(100, growth.achievement / 75 * 100) }
])

const pointRules = [
  { icon: '⚡', action: '完成微单元', points: '+5', note: '每个单元首次完成' },
  { icon: '📅', action: '每日学习打卡', points: '+10', note: '每日一次' },
  { icon: '↻', action: '完成间隔复习', points: '+3', note: '每知识点每日一次' },
  { icon: '🎯', action: '建立职业目标', points: '+20', note: '首次规划奖励' }
]

function checkIn() {
  const success = growth.checkIn()
  checkInMessage.value = success ? '打卡成功，积分 +10！' : '今天已经打过卡了。'
  window.setTimeout(() => { checkInMessage.value = '' }, 2400)
}
</script>

<template>
  <div class="page incentive-page">
    <header class="page-header">
      <span class="eyebrow">成长激励</span>
      <h1>我的成长中心</h1>
      <p>每一次学习、复习和坚持都有记录，让长期成长获得即时反馈。</p>
    </header>

    <section class="level-hero card">
      <div class="level-badge"><span>Lv</span><strong>{{ growth.level.number }}</strong></div>
      <div class="level-info">
        <span>当前等级</span>
        <h2>{{ growth.level.name }}</h2>
        <div class="level-track"><div :style="{ width: levelProgress + '%' }"></div></div>
        <small v-if="growth.level.next > growth.points">
          {{ growth.points }} / {{ growth.level.next }} 积分，距升级还差 {{ growth.level.next - growth.points }}
        </small>
        <small v-else>已达到当前最高等级</small>
      </div>
      <div class="points-total"><strong>{{ growth.points }}</strong><span>累计成长积分</span></div>
      <button class="check-btn" :disabled="growth.checkedInToday" @click="checkIn">
        <b>{{ growth.checkedInToday ? '✓' : '+' }}</b>
        <span>{{ growth.checkedInToday ? '今日已打卡' : '今日打卡' }}<small>连续 {{ growth.streak }} 天</small></span>
      </button>
    </section>
    <div v-if="checkInMessage" class="toast fade-up">{{ checkInMessage }}</div>

    <div class="two-column">
      <section class="card block">
        <div class="block-head"><div><h2>学习打卡</h2><p>最近 35 天学习记录</p></div><strong>🔥 {{ growth.streak }} 天</strong></div>
        <div class="calendar-labels"><span v-for="day in ['一','二','三','四','五','六','日']" :key="day">{{ day }}</span></div>
        <div class="calendar">
          <div
            v-for="day in calendarDays"
            :key="day.key"
            :title="`${day.key} ${day.checked ? '已打卡' : '未打卡'}`"
            :class="{ checked: day.checked, today: day.today }"
          >{{ day.day }}</div>
        </div>
        <div class="calendar-legend"><span><i></i>未打卡</span><span><i class="active"></i>已打卡</span><span>坚持 7 天可解锁“坚持者”勋章</span></div>
      </section>

      <section class="card block">
        <div class="block-head"><div><h2>积分明细</h2><p>用可验证行为积累成长值</p></div></div>
        <div class="point-rules">
          <div v-for="rule in pointRules" :key="rule.action">
            <span>{{ rule.icon }}</span>
            <div><strong>{{ rule.action }}</strong><small>{{ rule.note }}</small></div>
            <b>{{ rule.points }}</b>
          </div>
        </div>
        <router-link to="/review" class="review-cta">
          <span>今天有 {{ growth.dueReviewCount }} 项复习任务</span><b>去赚复习积分 →</b>
        </router-link>
      </section>
    </div>

    <section class="card block badges-block">
      <div class="block-head">
        <div><h2>勋章墙</h2><p>已解锁 {{ badgeCatalog.filter(item => item.unlocked).length }} / {{ badgeCatalog.length }}</p></div>
      </div>
      <div class="badge-grid">
        <article v-for="badge in badgeCatalog" :key="badge.name" :class="{ locked: !badge.unlocked }">
          <div class="badge-icon">{{ badge.icon }}<span v-if="badge.unlocked">✓</span></div>
          <h3>{{ badge.name }}</h3>
          <p>{{ badge.detail }}</p>
          <div class="badge-progress"><i :style="{ width: badge.progress + '%' }"></i></div>
          <small>{{ badge.unlocked ? '已获得' : `${Math.round(badge.progress)}%` }}</small>
        </article>
      </div>
    </section>
  </div>
</template>

<style scoped>
.incentive-page { max-width: 1060px; }
.eyebrow { color: var(--warning); font-size: 12px; font-weight: 800; letter-spacing: .08em; }
.level-hero { display: flex; align-items: center; gap: 18px; padding: 24px 27px; margin-bottom: 18px; overflow: hidden; color: white; border: none; background: linear-gradient(120deg, #312e81, #4f46e5 65%, #7c3aed); }
.level-badge { display: flex; width: 72px; height: 72px; flex: 0 0 72px; align-items: baseline; justify-content: center; padding-top: 16px; border: 2px solid rgba(255,255,255,.25); border-radius: 20px; background: rgba(255,255,255,.12); }
.level-badge span { color: #c7d2fe; font-size: 13px; font-weight: 700; }
.level-badge strong { font-size: 31px; line-height: 1; }
.level-info { flex: 1; }
.level-info > span, .level-info small { color: #c7d2fe; font-size: 11px; }
.level-info h2 { margin: 0 0 7px; font-size: 20px; }
.level-track { max-width: 360px; height: 7px; margin-bottom: 4px; overflow: hidden; border-radius: 9px; background: rgba(255,255,255,.18); }
.level-track div { height: 100%; border-radius: inherit; background: linear-gradient(90deg, #fbbf24, #fde68a); }
.points-total { display: flex; min-width: 110px; flex-direction: column; text-align: center; }
.points-total strong { color: #fde68a; font-size: 30px; line-height: 1.2; }
.points-total span { color: #c7d2fe; font-size: 10px; }
.check-btn { display: flex; align-items: center; gap: 9px; padding: 10px 15px; border: 1px solid rgba(255,255,255,.25); border-radius: 11px; background: rgba(255,255,255,.12); color: white; cursor: pointer; }
.check-btn:disabled { cursor: default; opacity: .72; }
.check-btn b { display: grid; place-items: center; width: 25px; height: 25px; border-radius: 50%; background: #fbbf24; color: #78350f; }
.check-btn span { display: flex; flex-direction: column; font-weight: 700; text-align: left; }
.check-btn small { color: #c7d2fe; font-weight: 400; }
.toast { position: fixed; top: 20px; right: 25px; z-index: 100; padding: 11px 16px; border-radius: 9px; background: #065f46; color: white; box-shadow: var(--shadow); }
.two-column { display: grid; grid-template-columns: 1.1fr .9fr; gap: 18px; margin-bottom: 18px; }
.block { padding: 22px 25px; box-shadow: none; }
.block-head { display: flex; align-items: flex-start; justify-content: space-between; margin-bottom: 16px; }
.block-head h2 { margin: 0; font-size: 17px; }
.block-head p { margin: 2px 0 0; color: var(--text-3); font-size: 11px; }
.block-head > strong { color: var(--warning); font-size: 13px; }
.calendar-labels, .calendar { display: grid; grid-template-columns: repeat(7, 1fr); gap: 6px; }
.calendar-labels span { color: var(--text-3); font-size: 9px; text-align: center; }
.calendar { margin-top: 6px; }
.calendar > div { display: grid; aspect-ratio: 1; place-items: center; border-radius: 7px; background: var(--surface-2); color: var(--text-3); font-size: 10px; }
.calendar > div.checked { background: linear-gradient(135deg, #a7f3d0, #10b981); color: #064e3b; font-weight: 700; }
.calendar > div.today { outline: 2px solid var(--primary); outline-offset: 1px; }
.calendar-legend { display: flex; gap: 13px; margin-top: 12px; color: var(--text-3); font-size: 9px; }
.calendar-legend span:last-child { margin-left: auto; }
.calendar-legend i { display: inline-block; width: 9px; height: 9px; margin-right: 4px; border-radius: 3px; background: var(--surface-2); }
.calendar-legend i.active { background: var(--success); }
.point-rules > div { display: flex; align-items: center; gap: 10px; padding: 9px 0; border-top: 1px solid var(--border); }
.point-rules > div:first-child { border-top: 0; }
.point-rules > div > span { font-size: 18px; }
.point-rules div div { display: flex; flex: 1; flex-direction: column; }
.point-rules strong { font-size: 12px; }
.point-rules small { color: var(--text-3); }
.point-rules b { color: var(--success); font-size: 13px; }
.review-cta { display: flex; justify-content: space-between; margin-top: 12px; padding: 9px 11px; border-radius: 8px; background: var(--primary-soft); font-size: 11px; }
.badge-grid { display: grid; grid-template-columns: repeat(6, 1fr); gap: 12px; }
.badge-grid article { padding: 16px 10px; border: 1px solid #fde68a; border-radius: 12px; background: #fffbeb; text-align: center; }
.badge-grid article.locked { border-color: var(--border); background: var(--surface-2); filter: grayscale(1); opacity: .65; }
.badge-icon { position: relative; width: 49px; margin: auto; font-size: 34px; }
.badge-icon span { position: absolute; right: -2px; bottom: 2px; display: grid; place-items: center; width: 17px; height: 17px; border-radius: 50%; background: var(--success); color: white; font-size: 9px; }
.badge-grid h3 { margin: 7px 0 2px; font-size: 12px; }
.badge-grid p { min-height: 30px; margin: 0; color: var(--text-3); font-size: 9px; }
.badge-progress { height: 4px; margin-top: 9px; overflow: hidden; border-radius: 4px; background: #e5e7eb; }
.badge-progress i { display: block; height: 100%; background: var(--warning); }
.badge-grid small { color: var(--text-3); font-size: 9px; }
@media (max-width: 800px) {
  .level-hero { align-items: flex-start; flex-wrap: wrap; }
  .level-info { min-width: calc(100% - 100px); }
  .points-total { flex: 1; }
  .two-column { grid-template-columns: 1fr; }
  .badge-grid { grid-template-columns: repeat(3, 1fr); }
}
@media (max-width: 480px) {
  .badge-grid { grid-template-columns: repeat(2, 1fr); }
  .calendar-legend span:last-child { display: none; }
}
</style>
