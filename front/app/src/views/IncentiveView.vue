<script setup>
import { computed, ref } from 'vue'
import { useGrowthStore } from '../stores/growth'
import { microUnits } from '../data/careerPath'
import Icon from '../components/Icon.vue'

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
  { icon: 'zap', action: '完成微单元', points: '+5', note: '每个单元首次完成' },
  { icon: 'calendar', action: '每日学习打卡', points: '+10', note: '每日一次' },
  { icon: 'refresh', action: '完成间隔复习', points: '+3', note: '每知识点每日一次' },
  { icon: 'target', action: '建立职业目标', points: '+20', note: '首次规划奖励' }
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
      <span class="eyebrow eyebrow--warning">成长激励</span>
      <h1>我的成长中心</h1>
      <p>每一次学习、复习和坚持都有记录，让长期成长获得即时反馈。</p>
    </header>

    <section class="level-hero card card--ink">
      <div class="level-badge">
        <span>Lv</span>
        <strong class="num">{{ growth.level.number }}</strong>
      </div>
      <div class="level-info">
        <span class="level-info__label">当前等级</span>
        <h2>{{ growth.level.name }}</h2>
        <div class="progress progress--glass level-track"><div :style="{ width: levelProgress + '%' }"></div></div>
        <small v-if="growth.level.next > growth.points">
          <b class="num">{{ growth.points }}</b> / {{ growth.level.next }} 积分，距升级还差 {{ growth.level.next - growth.points }}
        </small>
        <small v-else>已达到当前最高等级</small>
      </div>
      <div class="points-total">
        <strong class="num">{{ growth.points }}</strong>
        <span>累计成长积分</span>
      </div>
      <button class="check-btn" :disabled="growth.checkedInToday" @click="checkIn">
        <b><Icon :name="growth.checkedInToday ? 'check' : 'plus'" :size="16" :stroke="3" /></b>
        <span>{{ growth.checkedInToday ? '今日已打卡' : '今日打卡' }}<small><Icon name="fire" :size="11" /> 连续 {{ growth.streak }} 天</small></span>
      </button>
    </section>
    <div v-if="checkInMessage" class="toast fade-up">{{ checkInMessage }}</div>

    <div class="two-column">
      <section class="card panel">
        <div class="panel-head">
          <div><h2>学习打卡</h2><p>最近 35 天学习记录</p></div>
          <span class="tag tag--warning"><Icon name="fire" :size="12" /> {{ growth.streak }} 天</span>
        </div>
        <div class="calendar-labels"><span v-for="day in ['一','二','三','四','五','六','日']" :key="day">{{ day }}</span></div>
        <div class="calendar">
          <div
            v-for="day in calendarDays"
            :key="day.key"
            :title="`${day.key} ${day.checked ? '已打卡' : '未打卡'}`"
            :class="{ checked: day.checked, today: day.today }"
          >{{ day.day }}</div>
        </div>
        <div class="calendar-legend">
          <span><i></i>未打卡</span>
          <span><i class="active"></i>已打卡</span>
          <span class="calendar-legend__hint">坚持 7 天可解锁“坚持者”勋章</span>
        </div>
      </section>

      <section class="card panel">
        <div class="panel-head"><div><h2>积分明细</h2><p>用可验证行为积累成长值</p></div></div>
        <div class="point-rules">
          <div v-for="rule in pointRules" :key="rule.action">
            <span class="point-rules__icon"><Icon :name="rule.icon" :size="17" /></span>
            <div><strong>{{ rule.action }}</strong><small>{{ rule.note }}</small></div>
            <b class="num">{{ rule.points }}</b>
          </div>
        </div>
        <router-link to="/review" class="review-cta">
          <span>今天有 <b class="num">{{ growth.dueReviewCount }}</b> 项复习任务</span>
          <b>去赚复习积分 <Icon name="arrowRight" :size="13" /></b>
        </router-link>
      </section>
    </div>

    <section class="card panel badges-block">
      <div class="panel-head">
        <div><h2>勋章墙</h2><p>已解锁 {{ badgeCatalog.filter(item => item.unlocked).length }} / {{ badgeCatalog.length }}</p></div>
      </div>
      <div class="badge-grid stagger">
        <article v-for="badge in badgeCatalog" :key="badge.name" :class="{ locked: !badge.unlocked }">
          <div class="badge-icon">
            {{ badge.icon }}
            <span v-if="badge.unlocked"><Icon name="check" :size="10" :stroke="3.5" /></span>
          </div>
          <h3>{{ badge.name }}</h3>
          <p>{{ badge.detail }}</p>
          <div class="progress progress--thin badge-progress"><i :style="{ width: badge.progress + '%' }"></i></div>
          <small>{{ badge.unlocked ? '已获得' : `${Math.round(badge.progress)}%` }}</small>
        </article>
      </div>
    </section>
  </div>
</template>

<style scoped>
.incentive-page { max-width: 1080px; }

.level-hero { display: flex; align-items: center; gap: 22px; padding: 28px 32px; margin-bottom: 20px; border-radius: var(--radius-lg); }
.level-badge {
  display: flex;
  width: 80px;
  height: 80px;
  flex: 0 0 80px;
  align-items: baseline;
  justify-content: center;
  padding-top: 18px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 24px;
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.35), rgba(252, 211, 77, 0.15));
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.25);
}
.level-badge span { color: #fde68a; font-size: 13px; font-weight: 700; }
.level-badge strong { font-size: 36px; line-height: 1; color: #fff; }
.level-info { flex: 1; }
.level-info__label, .level-info small { color: rgba(201, 203, 230, 0.7); font-size: 12px; }
.level-info small b { color: #fff; }
.level-info h2 { margin: 2px 0 10px; font-size: 22px; }
.level-track { max-width: 380px; margin-bottom: 6px; }
.level-track div { background: linear-gradient(90deg, #f59e0b, #fde68a); }
.points-total { display: flex; min-width: 120px; flex-direction: column; text-align: center; }
.points-total strong { color: #fde68a; font-size: 34px; line-height: 1.15; }
.points-total span { color: rgba(201, 203, 230, 0.7); font-size: 11.5px; }
.check-btn {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 18px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.1);
  color: white;
  cursor: pointer;
  transition: all var(--t-fast);
}
.check-btn:hover:not(:disabled) { background: rgba(255, 255, 255, 0.18); transform: translateY(-1px); }
.check-btn:disabled { cursor: default; opacity: 0.75; }
.check-btn b { display: grid; place-items: center; width: 30px; height: 30px; border-radius: 50%; background: linear-gradient(135deg, #fbbf24, #f59e0b); color: #78350f; }
.check-btn span { display: flex; flex-direction: column; font-weight: 700; text-align: left; }
.check-btn small { display: inline-flex; align-items: center; gap: 3px; color: rgba(201, 203, 230, 0.75); font-weight: 400; font-size: 12px; }

.two-column { display: grid; grid-template-columns: 1.1fr 0.9fr; gap: 18px; margin-bottom: 18px; }
.calendar-labels, .calendar { display: grid; grid-template-columns: repeat(7, 1fr); gap: 7px; }
.calendar-labels span { color: var(--text-3); font-size: 11px; text-align: center; font-weight: 600; }
.calendar { margin-top: 8px; }
.calendar > div { display: grid; aspect-ratio: 1; place-items: center; border-radius: 10px; background: var(--surface-3); color: var(--text-3); font-size: 12px; transition: transform var(--t-fast); }
.calendar > div:hover { transform: scale(1.06); }
.calendar > div.checked { background: linear-gradient(135deg, #6ee7b7, #0fb981); color: #064e3b; font-weight: 700; box-shadow: 0 4px 10px rgba(15, 185, 129, 0.25); }
.calendar > div.today { outline: 2px solid var(--primary); outline-offset: 2px; }
.calendar-legend { display: flex; align-items: center; gap: 14px; margin-top: 14px; color: var(--text-3); font-size: 11.5px; }
.calendar-legend__hint { margin-left: auto; }
.calendar-legend i { display: inline-block; width: 10px; height: 10px; margin-right: 5px; border-radius: 3px; background: var(--surface-3); vertical-align: -1px; }
.calendar-legend i.active { background: var(--success); }

.point-rules > div { display: flex; align-items: center; gap: 12px; padding: 11px 0; }
.point-rules > div + div { border-top: 1px solid var(--border); }
.point-rules__icon { display: grid; place-items: center; width: 36px; height: 36px; flex: 0 0 36px; border-radius: 11px; background: var(--primary-soft); color: var(--primary-strong); }
.point-rules div div { display: flex; flex: 1; flex-direction: column; }
.point-rules strong { font-size: 13.5px; }
.point-rules small { color: var(--text-3); font-size: 12px; }
.point-rules > div > b { padding: 3px 10px; border-radius: 99px; background: var(--success-soft); color: var(--success-strong); font-size: 13px; }
.review-cta { display: flex; align-items: center; justify-content: space-between; margin-top: 14px; padding: 12px 14px; border-radius: 12px; background: var(--primary-soft); color: var(--primary-deep); font-size: 13px; transition: background var(--t-fast); }
.review-cta:hover { background: var(--primary-soft-2); }
.review-cta > b { display: inline-flex; align-items: center; gap: 3px; color: var(--primary-strong); }

.badge-grid { display: grid; grid-template-columns: repeat(6, 1fr); gap: 14px; }
.badge-grid article {
  position: relative;
  padding: 20px 12px 16px;
  border: 1px solid #fde68a;
  border-radius: 16px;
  background: linear-gradient(180deg, #fffbeb, #fff);
  text-align: center;
  transition: transform var(--t), box-shadow var(--t);
}
.badge-grid article:hover { transform: translateY(-3px); box-shadow: var(--shadow); }
.badge-grid article.locked { border-color: var(--border); background: var(--surface-2); }
.badge-grid article.locked .badge-icon { filter: grayscale(1); opacity: 0.5; }
.badge-icon { position: relative; width: 56px; height: 56px; margin: auto; display: grid; place-items: center; border-radius: 18px; background: rgba(255, 255, 255, 0.8); font-size: 30px; box-shadow: var(--shadow-xs); }
.badge-icon span { position: absolute; right: -4px; bottom: -4px; display: grid; place-items: center; width: 18px; height: 18px; border-radius: 50%; background: var(--success); color: white; box-shadow: 0 0 0 2px #fff; }
.badge-grid h3 { margin: 10px 0 2px; font-size: 13.5px; }
.badge-grid p { min-height: 34px; margin: 0; color: var(--text-3); font-size: 11.5px; line-height: 1.5; }
.badge-progress { margin-top: 10px; height: 4px; }
.badge-progress i { background: linear-gradient(90deg, var(--warning), #fcd34d); }
.badge-grid small { display: block; margin-top: 6px; color: var(--text-3); font-size: 11px; font-weight: 600; }

@media (max-width: 900px) {
  .level-hero { align-items: flex-start; flex-wrap: wrap; }
  .level-info { min-width: calc(100% - 110px); }
  .points-total { flex: 1; }
  .two-column { grid-template-columns: 1fr; }
  .badge-grid { grid-template-columns: repeat(3, 1fr); }
}
@media (max-width: 520px) {
  .badge-grid { grid-template-columns: repeat(2, 1fr); }
  .calendar-legend__hint { display: none; }
}
</style>
