<template>
  <div class="page">
    <PageHeader eyebrow="成长生态" title="生态看板" subtitle="知识生产 → 知识传播 → 能力培养 → 项目验证 → 商业价值 → 人才成长 的全链路运营数据。">
      <el-button :icon="Refresh" :loading="loading" @click="load">刷新</el-button>
    </PageHeader>

    <template v-if="d">
      <section class="chain">
        <div v-for="(step, i) in chain" :key="step.label" class="chain__step">
          <div class="chain__icon"><el-icon :size="18"><component :is="step.icon" /></el-icon></div>
          <div class="chain__label">{{ step.label }}</div>
          <div class="chain__value num">{{ step.value }}</div>
          <div class="chain__hint">{{ step.hint }}</div>
          <el-icon v-if="i < chain.length - 1" class="chain__arrow"><Right /></el-icon>
        </div>
      </section>

      <div class="stat-grid">
        <StatCard label="职业 / 岗位" :value="`${d.careers} / ${d.roles}`" icon="Guide" tone="primary" hint="职业体系" />
        <StatCard label="技能模型" :value="d.skills" icon="Medal" tone="info" hint="L0-L5 等级" />
        <StatCard label="知识点 / 关系" :value="`${d.knowledge} / ${d.relations}`" icon="Share" tone="sky" hint="知识图谱" />
        <StatCard label="设定目标的用户" :value="d.goals" icon="Flag" tone="success" hint="进入成长飞轮" />
      </div>

      <div class="grid">
        <section class="panel">
          <div class="panel__head"><div><h3 class="panel__title">技能等级分布</h3><p class="panel__sub">全平台人才的技能状态（L1 及以上）</p></div></div>
          <div class="panel__body">
            <div v-for="(n, lvl) in d.levelDistribution" v-show="lvl > 0" :key="lvl" class="bar">
              <span class="bar__label"><b :style="{ background: LEVEL_COLORS[lvl] }">L{{ lvl }}</b>{{ LEVELS[lvl] }}</span>
              <div class="bar__track"><i :style="{ width: (n / maxLevel) * 100 + '%', background: LEVEL_COLORS[lvl] }" /></div>
              <span class="num bar__value">{{ n }}</span>
            </div>
          </div>
        </section>

        <section class="panel">
          <div class="panel__head"><div><h3 class="panel__title">知识资产构成</h3><p class="panel__sub">统一 Knowledge Resource</p></div></div>
          <div class="panel__body">
            <div v-for="(label, key) in RESOURCE_TYPES" :key="key" class="bar">
              <span class="bar__label">{{ label }}</span>
              <div class="bar__track"><i :style="{ width: ((d.resources[key] || 0) / maxResource) * 100 + '%' }" /></div>
              <span class="num bar__value">{{ d.resources[key] || 0 }}</span>
            </div>
          </div>
        </section>
      </div>

      <div class="stat-grid">
        <StatCard label="项目 / 提交" :value="`${d.projects} / ${d.submissions}`" icon="Files" tone="primary" clickable @click="$router.push('/eco/projects')" />
        <StatCard label="开放任务 / 全部" :value="`${d.openOpportunities} / ${d.opportunities}`" icon="Suitcase" tone="warning" clickable @click="$router.push('/eco/market')" />
        <StatCard label="已验收任务" :value="d.completedTasks" icon="CircleCheck" tone="success" :hint="`共 ${d.applications} 份申请`" />
        <StatCard label="AI 内核运行" :value="d.pipelineRuns" icon="Connection" tone="ink" clickable @click="$router.push('/eco/kernel')" />
      </div>
    </template>
  </div>
</template>

<script setup>
import { computed, onMounted, ref } from 'vue'
import { Refresh, Right, Reading, Share, Medal, Files, Suitcase, Money } from '@element-plus/icons-vue'
import PageHeader from '../../components/common/PageHeader.vue'
import StatCard from '../../components/common/StatCard.vue'
import { ecoAdminApi } from '../../api'
import { LEVELS, LEVEL_COLORS, RESOURCE_TYPES, money } from '../../utils/eco'

const d = ref(null)
const loading = ref(false)

const maxLevel = computed(() => Math.max(1, ...(d.value?.levelDistribution || []).slice(1)))
const maxResource = computed(() => Math.max(1, ...Object.values(d.value?.resources || {})))
const resourceTotal = computed(() => Object.values(d.value?.resources || {}).reduce((s, n) => s + n, 0))

const chain = computed(() => [
  { label: '知识生产', value: d.value.posts + resourceTotal.value, hint: '社区内容与知识资产', icon: Reading },
  { label: '知识传播', value: d.value.learningEvents, hint: '学习事件', icon: Share },
  { label: '能力培养', value: d.value.levelDistribution.slice(1).reduce((s, n) => s + n, 0), hint: '技能等级记录', icon: Medal },
  { label: '项目验证', value: d.value.evidence, hint: '能力证据', icon: Files },
  { label: '商业价值', value: money(d.value.totalIncome), hint: '人才任务收入', icon: Money },
  { label: '人才成长', value: d.value.completedTasks, hint: '真实任务交付', icon: Suitcase }
])

async function load() {
  loading.value = true
  try {
    d.value = await ecoAdminApi.dashboard()
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.page { display: flex; flex-direction: column; gap: 18px; }
.chain { display: grid; grid-template-columns: repeat(6, 1fr); gap: 12px; }
.chain__step { position: relative; padding: 16px; border: 1px solid var(--border); border-radius: var(--radius); background: var(--surface); box-shadow: var(--shadow-sm); }
.chain__icon { display: inline-flex; align-items: center; justify-content: center; width: 34px; height: 34px; margin-bottom: 8px; border-radius: 10px; background: var(--primary-soft); color: var(--primary); }
.chain__label { color: var(--text-3); font-size: 12.5px; font-weight: 600; }
.chain__value { margin-top: 2px; font-size: 24px; font-weight: 800; }
.chain__hint { color: var(--text-3); font-size: 12px; }
.chain__arrow { position: absolute; top: 50%; right: -13px; z-index: 1; padding: 3px; border-radius: 50%; background: var(--surface); color: var(--primary); box-shadow: var(--shadow-sm); }
.grid { display: grid; grid-template-columns: 1fr 1fr; gap: 18px; }
.bar { display: flex; align-items: center; gap: 12px; margin-bottom: 12px; font-size: 13px; }
.bar__label { display: flex; align-items: center; gap: 8px; width: 160px; color: var(--text-2); }
.bar__label b { padding: 0 6px; border-radius: 6px; color: #fff; font-size: 11px; }
.bar__track { flex: 1; height: 8px; overflow: hidden; border-radius: 4px; background: var(--surface-3); }
.bar__track i { display: block; height: 100%; border-radius: 4px; background: var(--primary); }
.bar__value { width: 36px; text-align: right; font-weight: 700; }
@media (max-width: 1100px) {
  .chain { grid-template-columns: repeat(3, 1fr); }
  .chain__arrow { display: none; }
  .grid { grid-template-columns: 1fr; }
}
</style>
