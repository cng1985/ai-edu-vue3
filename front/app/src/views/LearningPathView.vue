<script setup>
import { computed } from 'vue'
import { frontendPath, microUnits } from '../data/careerPath'
import { useGrowthStore } from '../stores/growth'

const growth = useGrowthStore()
const unitsByDomain = computed(() => frontendPath.competencies.map((domain) => ({
  ...domain,
  units: domain.points.map((point) => ({
    ...microUnits.find((unit) => unit.id === point.unitId),
    pointName: point.name
  }))
})))
</script>

<template>
  <div class="page path-page">
    <header class="page-header path-header">
      <div>
        <span class="tag">个性化标准路径</span>
        <h1>我的学习路径</h1>
        <p v-if="growth.hasGoal">{{ growth.goal.name }} · 每周 {{ growth.goal.weeklyHours }} 小时</p>
        <p v-else>设定职业目标后，系统将按知识依赖生成可执行路径。</p>
      </div>
      <router-link v-if="!growth.hasGoal" to="/career" class="btn btn--primary">创建目标</router-link>
      <div v-else class="achievement"><strong>{{ growth.achievement }}%</strong><span>目标达成度</span></div>
    </header>

    <div v-if="growth.hasGoal" class="path-timeline">
      <section v-for="(domain, domainIndex) in unitsByDomain" :key="domain.id" class="domain-section">
        <div class="domain-marker" :style="{ background: domain.color }">{{ domainIndex + 1 }}</div>
        <div class="domain-content">
          <div class="domain-head">
            <div>
              <h2>{{ domain.name }}</h2>
              <span>能力权重 {{ domain.weight }}%</span>
            </div>
            <strong>{{ growth.competencyProgress.find(item => item.id === domain.id)?.progress }}%</strong>
          </div>
          <div class="progress-track">
            <div :style="{ width: growth.competencyProgress.find(item => item.id === domain.id)?.progress + '%', background: domain.color }"></div>
          </div>
          <div class="unit-grid">
            <router-link
              v-for="(unit, unitIndex) in domain.units"
              :key="unit.id"
              :to="`/micro/${unit.id}`"
              class="unit-card card"
              :class="{ done: growth.isUnitCompleted(unit.id), current: growth.nextUnit?.id === unit.id }"
            >
              <div class="unit-order">
                <span v-if="growth.isUnitCompleted(unit.id)">✓</span>
                <span v-else>{{ unitIndex + 1 }}</span>
              </div>
              <div>
                <small>{{ unit.duration }} 分钟 · {{ unit.difficulty }}</small>
                <h3>{{ unit.title }}</h3>
                <span v-if="growth.isUnitCompleted(unit.id)" class="status done-text">已掌握 · {{ growth.unitScores[unit.id] }} 分</span>
                <span v-else-if="growth.nextUnit?.id === unit.id" class="status">当前任务 →</span>
                <span v-else class="status pending">待学习</span>
              </div>
            </router-link>
          </div>
        </div>
      </section>
    </div>

    <div v-else class="empty card">
      <div>🧭</div>
      <h2>先确认方向，再开始学习</h2>
      <p>AI 将根据你的基础、时间和目标拆解能力域、知识点与里程碑。</p>
      <router-link to="/career" class="btn btn--primary">开始职业规划</router-link>
    </div>
  </div>
</template>

<style scoped>
.path-page { max-width: 1040px; }
.path-header { display: flex; align-items: center; justify-content: space-between; }
.path-header h1 { margin-top: 8px; }
.achievement { display: flex; flex-direction: column; text-align: right; }
.achievement strong { color: var(--primary); font-size: 30px; line-height: 1.1; }
.achievement span { color: var(--text-3); font-size: 12px; }
.path-timeline { position: relative; }
.path-timeline::before { content: ''; position: absolute; left: 20px; top: 24px; bottom: 30px; width: 2px; background: var(--border); }
.domain-section { position: relative; display: flex; gap: 22px; margin-bottom: 28px; }
.domain-marker { z-index: 1; display: grid; place-items: center; width: 42px; height: 42px; flex: 0 0 42px; border-radius: 50%; color: white; font-weight: 800; box-shadow: 0 0 0 6px var(--bg); }
.domain-content { flex: 1; min-width: 0; }
.domain-head { display: flex; justify-content: space-between; align-items: center; margin-bottom: 8px; }
.domain-head h2 { display: inline; margin: 0 10px 0 0; font-size: 19px; }
.domain-head span { color: var(--text-3); font-size: 12px; }
.domain-head strong { color: var(--text-2); }
.progress-track { height: 6px; margin-bottom: 12px; border-radius: 9px; background: var(--border); overflow: hidden; }
.progress-track div { height: 100%; transition: width .4s; }
.unit-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 12px; }
.unit-card { display: flex; gap: 12px; padding: 16px; color: inherit; box-shadow: none; transition: .15s; }
.unit-card:hover, .unit-card.current { border-color: var(--primary); transform: translateY(-1px); }
.unit-card.done { background: #f7fefb; border-color: #a7f3d0; }
.unit-order { display: grid; place-items: center; width: 30px; height: 30px; flex: 0 0 30px; border-radius: 8px; background: var(--primary-soft); color: var(--primary); font-weight: 700; }
.done .unit-order { background: var(--success); color: white; }
.unit-card small { color: var(--text-3); }
.unit-card h3 { margin: 2px 0 5px; font-size: 14px; }
.status { color: var(--primary); font-size: 12px; font-weight: 600; }
.done-text { color: var(--success); }
.pending { color: var(--text-3); }
.empty { margin-top: 70px; padding: 50px; text-align: center; }
.empty > div { font-size: 42px; }
.empty p { color: var(--text-2); }
@media (max-width: 700px) {
  .unit-grid { grid-template-columns: 1fr; }
  .path-header { align-items: flex-start; }
}
</style>
