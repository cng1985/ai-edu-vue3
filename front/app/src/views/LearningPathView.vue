<script setup>
import { computed } from 'vue'
import { frontendPath, microUnits } from '../data/careerPath'
import { useGrowthStore } from '../stores/growth'
import ProgressRing from '../components/ProgressRing.vue'
import Icon from '../components/Icon.vue'

const growth = useGrowthStore()
const unitsByDomain = computed(() => frontendPath.competencies.map((domain) => ({
  ...domain,
  progress: growth.competencyProgress.find((item) => item.id === domain.id)?.progress || 0,
  units: domain.points.map((point) => ({
    ...microUnits.find((unit) => unit.id === point.unitId),
    pointName: point.name
  }))
})))
const totalUnits = computed(() => microUnits.length)
</script>

<template>
  <div class="page path-page">
    <header class="page-header page-header--split">
      <div>
        <span class="eyebrow">个性化标准路径</span>
        <h1>我的学习路径</h1>
        <p v-if="growth.hasGoal">{{ growth.goal.name }} · 每周 {{ growth.goal.weeklyHours }} 小时 · 已完成 {{ growth.completedUnitCount }}/{{ totalUnits }} 个微单元</p>
        <p v-else>设定职业目标后，系统将按知识依赖生成可执行路径。</p>
      </div>
      <router-link v-if="!growth.hasGoal" to="/career" class="btn btn--primary btn--lg">创建目标 <Icon name="arrowRight" :size="16" /></router-link>
      <div v-else class="achievement card">
        <ProgressRing :percent="growth.achievement" :size="64" :stroke="7" />
        <div><strong>目标达成度</strong><span>达标线 75%</span></div>
      </div>
    </header>

    <div v-if="growth.hasGoal" class="path-timeline">
      <section v-for="(domain, domainIndex) in unitsByDomain" :key="domain.id" class="domain-section">
        <div class="domain-marker" :style="{ background: domain.color }">{{ domainIndex + 1 }}</div>
        <div class="domain-content">
          <div class="domain-head">
            <div>
              <h2>{{ domain.name }}</h2>
              <span class="tag tag--neutral">能力权重 {{ domain.weight }}%</span>
            </div>
            <strong class="num" :style="{ color: domain.color }">{{ domain.progress }}%</strong>
          </div>
          <div class="progress progress--thin domain-progress">
            <div :style="{ width: domain.progress + '%', background: domain.color }"></div>
          </div>
          <div class="unit-grid">
            <router-link
              v-for="(unit, unitIndex) in domain.units"
              :key="unit.id"
              :to="`/micro/${unit.id}`"
              class="unit-card card card--hover"
              :class="{ 'unit-card--done': growth.isUnitCompleted(unit.id), 'unit-card--current': growth.nextUnit?.id === unit.id }"
            >
              <div class="unit-order">
                <Icon v-if="growth.isUnitCompleted(unit.id)" name="check" :size="15" :stroke="3" />
                <span v-else>{{ unitIndex + 1 }}</span>
              </div>
              <div class="unit-body">
                <small><Icon name="clock" :size="11" /> {{ unit.duration }} 分钟 · {{ unit.difficulty }}</small>
                <h3>{{ unit.title }}</h3>
                <span v-if="growth.isUnitCompleted(unit.id)" class="status status--done">已掌握 · {{ growth.unitScores[unit.id] }} 分</span>
                <span v-else-if="growth.nextUnit?.id === unit.id" class="status status--current">当前任务 <Icon name="arrowRight" :size="12" /></span>
                <span v-else class="status status--pending">待学习</span>
              </div>
              <span v-if="growth.nextUnit?.id === unit.id" class="unit-pulse"></span>
            </router-link>
          </div>
        </div>
      </section>
    </div>

    <div v-else class="empty-state card">
      <div class="empty-state__icon"><Icon name="compass" :size="32" /></div>
      <h2>先确认方向，再开始学习</h2>
      <p>AI 将根据你的基础、时间和目标拆解能力域、知识点与里程碑。</p>
      <router-link to="/career" class="btn btn--primary">开始职业规划</router-link>
    </div>
  </div>
</template>

<style scoped>
.path-page { max-width: 1080px; }
.achievement { display: flex; align-items: center; gap: 14px; padding: 12px 18px 12px 12px; }
.achievement div { display: flex; flex-direction: column; line-height: 1.35; }
.achievement strong { font-size: 14px; }
.achievement span { color: var(--text-3); font-size: 12px; }

.path-timeline { position: relative; padding-top: 6px; }
.path-timeline::before { content: ''; position: absolute; left: 22px; top: 30px; bottom: 40px; width: 2px; background: linear-gradient(180deg, var(--border-strong), var(--border) 90%, transparent); }
.domain-section { position: relative; display: flex; gap: 24px; margin-bottom: 36px; }
.domain-marker {
  z-index: 1;
  display: grid;
  place-items: center;
  width: 46px;
  height: 46px;
  flex: 0 0 46px;
  border-radius: 15px;
  color: white;
  font-family: 'Sora', sans-serif;
  font-weight: 800;
  box-shadow: 0 0 0 6px var(--bg), inset 0 1px 0 rgba(255, 255, 255, 0.3), 0 8px 18px rgba(22, 24, 44, 0.12);
}
.domain-content { flex: 1; min-width: 0; }
.domain-head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 10px; }
.domain-head > div { display: flex; align-items: center; gap: 10px; }
.domain-head h2 { margin: 0; font-size: 20px; }
.domain-head strong { font-size: 20px; }
.domain-progress { margin-bottom: 14px; }
.unit-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.unit-card { position: relative; display: flex; gap: 14px; padding: 16px 18px; color: inherit; }
.unit-card--current { border-color: var(--primary); box-shadow: 0 0 0 4px var(--primary-soft); }
.unit-card--done { background: linear-gradient(140deg, #f2fcf8, #fff 60%); border-color: #bfeedd; }
.unit-order { display: grid; place-items: center; width: 34px; height: 34px; flex: 0 0 34px; border-radius: 11px; background: var(--primary-soft); color: var(--primary-strong); font-weight: 700; }
.unit-card--done .unit-order { background: var(--success); color: white; }
.unit-card--current .unit-order { background: linear-gradient(135deg, var(--primary), var(--primary-strong)); color: #fff; }
.unit-body { display: flex; flex-direction: column; min-width: 0; }
.unit-body small { display: inline-flex; align-items: center; gap: 4px; color: var(--text-3); font-size: 12px; }
.unit-body h3 { margin: 3px 0 6px; font-size: 14.5px; }
.status { display: inline-flex; align-items: center; gap: 4px; font-size: 12px; font-weight: 600; }
.status--current { color: var(--primary); }
.status--done { color: var(--success-strong); }
.status--pending { color: var(--text-3); }
.unit-pulse { position: absolute; top: 14px; right: 14px; width: 9px; height: 9px; border-radius: 50%; background: var(--primary); box-shadow: 0 0 0 0 var(--primary-glow); animation: pulse 1.8s infinite; }
@keyframes pulse { 0% { box-shadow: 0 0 0 0 var(--primary-glow); } 70% { box-shadow: 0 0 0 10px transparent; } 100% { box-shadow: 0 0 0 0 transparent; } }

@media (max-width: 760px) {
  .unit-grid { grid-template-columns: 1fr; }
  .path-timeline::before { display: none; }
  .domain-section { flex-direction: column; gap: 12px; }
  .domain-head { flex-direction: column; align-items: flex-start; gap: 6px; }
}
</style>
