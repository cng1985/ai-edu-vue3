<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { ecoApi, meApi } from '../api'
import { useAuthStore } from '../stores/auth'
import { APPLICATION_STATUS, OPPORTUNITY_STATUS, money, timeAgo } from '../utils/eco'
import SkillLevel from '../components/SkillLevel.vue'
import Icon from '../components/Icon.vue'

const auth = useAuthStore()
const list = ref([])
const status = ref('open')
const sort = ref('match')
const mine = ref([])
const skillNames = ref({})
const loading = ref(false)
const isTalent = computed(() => ['learner', 'creator'].includes(auth.user?.role))

async function load() {
  loading.value = true
  try {
    list.value = await ecoApi.opportunities({ status: status.value, sort: sort.value })
  } finally {
    loading.value = false
  }
}

watch([status, sort], load)
onMounted(async () => {
  ecoApi.skills().then((list) => (skillNames.value = Object.fromEntries(list.map((s) => [s.id, s.name]))))
  await load()
  if (isTalent.value) mine.value = await meApi.applications().catch(() => [])
})
</script>

<template>
  <div class="page market">
    <header class="page-header page-header--split">
      <div>
        <span class="eyebrow">IT Task Marketplace</span>
        <h1>任务市场</h1>
        <p>企业发布真实 IT 任务，系统按「任务要求 → 技能模型 → 人才画像」计算匹配度。完成任务获得收益，任务表现会成为最有分量的能力证据。</p>
      </div>
      <router-link v-if="auth.hasPermission('opportunity:publish')" to="/enterprise" class="btn btn--primary"><Icon name="plus" :size="15" /> 发布任务</router-link>
    </header>

    <div class="row row--between row--wrap">
      <div class="segmented">
        <button :class="{ active: status === 'open' }" @click="status = 'open'">招募中</button>
        <button :class="{ active: status === 'in_progress' }" @click="status = 'in_progress'">进行中</button>
        <button :class="{ active: status === 'closed' }" @click="status = 'closed'">已结束</button>
      </div>
      <div v-if="isTalent" class="segmented">
        <button :class="{ active: sort === 'match' }" @click="sort = 'match'">按匹配度</button>
        <button :class="{ active: sort === '' }" @click="sort = ''">最新发布</button>
      </div>
    </div>

    <div class="market__body">
      <div class="stack">
        <div v-if="loading" class="muted">加载中…</div>
        <div v-else-if="!list.length" class="card empty-state"><h3>暂无任务</h3><p>换个筛选条件看看</p></div>
        <router-link v-for="o in list" :key="o.id" :to="`/market/${o.id}`" class="card card--hover opp">
          <div class="opp__main">
            <div class="row row--wrap">
              <span class="tag" :class="`tag--${OPPORTUNITY_STATUS[o.status]?.tone}`">{{ OPPORTUNITY_STATUS[o.status]?.label }}</span>
              <span class="chip">{{ o.mode === 'remote' ? '远程' : '驻场' }}</span>
              <span class="chip">{{ o.durationDays }} 天</span>
              <span v-if="o.myApplication" class="tag" :class="`tag--${APPLICATION_STATUS[o.myApplication.status]?.tone}`">我的申请：{{ APPLICATION_STATUS[o.myApplication.status]?.label }}</span>
            </div>
            <h3>{{ o.title }}</h3>
            <p>{{ o.description }}</p>
            <div class="opp__reqs">
              <template v-if="o.match">
                <span v-for="it in o.match.items" :key="it.skillId" class="opp__req">
                  {{ it.skillName }} <SkillLevel :level="it.current" :required="it.required" />
                </span>
              </template>
              <template v-else>
                <span v-for="r in o.requirements" :key="r.skillId" class="chip">{{ skillNames[r.skillId] || r.skillId }} L{{ r.level }}</span>
              </template>
            </div>
            <small class="muted">{{ o.company }} · {{ timeAgo(o.createdAt) }} · {{ o.applicantCount }} 人申请</small>
          </div>
          <div class="opp__side">
            <strong class="num">{{ money(o.budget) }}</strong>
            <template v-if="o.match">
              <div class="opp__match" :class="{ 'opp__match--ok': o.match.qualified }">
                <b class="num">{{ o.match.score }}%</b><small>{{ o.match.qualified ? '完全胜任' : '匹配度' }}</small>
              </div>
            </template>
          </div>
        </router-link>
      </div>

      <aside v-if="isTalent" class="card panel market__mine">
        <h3>我的任务</h3>
        <p v-if="!mine.length" class="muted small">还没有申请任务。匹配度 ≥ 80% 的任务是很好的起点。</p>
        <router-link v-for="a in mine" :key="a.id" :to="`/market/${a.opportunityId}`" class="mini">
          <span>{{ a.opportunity?.title }}<small>{{ a.opportunity?.company }} · 申请时匹配 {{ a.matchScore }}%</small></span>
          <span class="tag" :class="`tag--${APPLICATION_STATUS[a.status]?.tone}`">{{ APPLICATION_STATUS[a.status]?.label }}</span>
        </router-link>
      </aside>
    </div>
  </div>
</template>

<style scoped>
.market { display: flex; flex-direction: column; gap: 16px; }
.market .page-header { margin-bottom: 4px; }
.market__body { display: grid; grid-template-columns: minmax(0, 1fr) 300px; gap: 18px; align-items: start; }
.opp { display: flex; gap: 20px; padding: 20px 22px; color: inherit; }
.opp__main { display: flex; flex: 1; min-width: 0; flex-direction: column; gap: 8px; }
.opp__main h3 { margin: 2px 0 0; font-size: 17px; }
.opp__main p { display: -webkit-box; margin: 0; overflow: hidden; color: var(--text-2); font-size: 13.5px; -webkit-box-orient: vertical; -webkit-line-clamp: 2; }
.opp__reqs { display: flex; flex-wrap: wrap; gap: 8px 16px; }
.opp__req { display: inline-flex; align-items: center; gap: 8px; font-size: 12.5px; font-weight: 600; }
.opp__side { display: flex; width: 120px; flex-direction: column; align-items: flex-end; justify-content: space-between; gap: 12px; }
.opp__side > strong { color: var(--warning-strong); font-size: 20px; }
.opp__match { display: flex; flex-direction: column; align-items: center; padding: 10px 14px; border-radius: 14px; background: var(--primary-soft); color: var(--primary-strong); }
.opp__match--ok { background: var(--success-soft); color: var(--success-strong); }
.opp__match b { font-size: 22px; line-height: 1.1; }
.opp__match small { font-size: 11px; }
.market__mine { position: sticky; top: 24px; }
.market__mine h3 { margin: 0 0 10px; font-size: 15px; }
.mini { display: flex; align-items: center; justify-content: space-between; gap: 10px; padding: 10px 0; border-bottom: 1px dashed var(--border); color: inherit; font-size: 13.5px; }
.mini > span:first-child { display: flex; flex-direction: column; min-width: 0; }
.mini small { color: var(--text-3); font-size: 11.5px; }
@media (max-width: 960px) {
  .market__body { grid-template-columns: 1fr; }
  .market__mine { position: static; }
  .opp { flex-direction: column; }
  .opp__side { width: auto; flex-direction: row; align-items: center; }
}
</style>
