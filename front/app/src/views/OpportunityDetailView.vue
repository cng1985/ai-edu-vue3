<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { ecoApi } from '../api'
import { useAuthStore } from '../stores/auth'
import { APPLICATION_STATUS, OPPORTUNITY_STATUS, money, timeAgo } from '../utils/eco'
import SkillLevel from '../components/SkillLevel.vue'
import ProgressRing from '../components/ProgressRing.vue'
import UserAvatar from '../components/UserAvatar.vue'
import Icon from '../components/Icon.vue'

const route = useRoute()
const auth = useAuthStore()
const o = ref(null)
const message = ref('')
const applying = ref(false)
const error = ref('')
const skillNames = ref({})

const canApply = computed(() => auth.hasPermission('growth:write') && ['learner', 'creator'].includes(auth.user?.role))

async function load() {
  o.value = await ecoApi.opportunity(route.params.id)
}

async function apply() {
  applying.value = true
  error.value = ''
  try {
    await ecoApi.apply(o.value.id, message.value)
    await load()
  } catch (e) {
    error.value = e.message
  } finally {
    applying.value = false
  }
}

onMounted(() => {
  ecoApi.skills().then((list) => (skillNames.value = Object.fromEntries(list.map((s) => [s.id, s.name]))))
  load()
})
</script>

<template>
  <div class="page od">
    <router-link to="/market" class="back-link"><Icon name="arrowLeft" :size="14" /> 任务市场</router-link>
    <template v-if="o">
      <div class="od__body">
        <section class="card panel">
          <div class="row row--wrap">
            <span class="tag" :class="`tag--${OPPORTUNITY_STATUS[o.status]?.tone}`">{{ OPPORTUNITY_STATUS[o.status]?.label }}</span>
            <span class="chip">{{ o.mode === 'remote' ? '远程' : '驻场' }}</span>
            <span class="chip">工期 {{ o.durationDays }} 天</span>
            <span class="chip">{{ o.applicantCount }} 人申请</span>
          </div>
          <h1>{{ o.title }}</h1>
          <div class="row od__publisher">
            <UserAvatar :user="o.publisher" />
            <span><strong>{{ o.company }}</strong><small>发布于 {{ timeAgo(o.createdAt) }}</small></span>
            <strong class="od__budget num">{{ money(o.budget) }}</strong>
          </div>
          <h3>需求描述</h3>
          <p class="od__desc">{{ o.description }}</p>
          <h3>技能要求</h3>
          <div v-for="it in o.match?.items || []" :key="it.skillId" class="od__req">
            <span>{{ it.skillName }}</span>
            <SkillLevel :level="it.current" :required="it.required" show-name />
          </div>
          <div v-if="!o.match" class="row row--wrap">
            <span v-for="r in o.requirements" :key="r.skillId" class="chip">{{ skillNames[r.skillId] || r.skillId }} L{{ r.level }}</span>
          </div>
        </section>

        <aside class="stack">
          <section v-if="o.match" class="card panel od__match">
            <ProgressRing :percent="o.match.score" :size="120" :stroke="10" />
            <strong>{{ o.match.qualified ? '你完全胜任这个任务' : `已满足 ${o.match.metCount}/${o.match.items.length} 项要求` }}</strong>
            <p class="muted small">匹配度由你的技能等级与任务要求逐项比较得出；差 2 级及以上的技能会显著降低匹配度。</p>
            <router-link v-if="!o.match.qualified" :to="`/agents/opportunity`" class="btn btn--sm btn--soft"><Icon name="target" :size="13" /> 问 Opportunity Agent 怎么补齐</router-link>
          </section>

          <section v-if="o.myApplication" class="card panel">
            <h3 class="side-title">我的申请</h3>
            <span class="tag" :class="`tag--${APPLICATION_STATUS[o.myApplication.status]?.tone}`">{{ APPLICATION_STATUS[o.myApplication.status]?.label }}</span>
            <p class="small muted">申请时匹配度 {{ o.myApplication.matchScore }}% · {{ timeAgo(o.myApplication.createdAt) }}</p>
            <p v-if="o.myApplication.status === 'completed'" class="small">企业评价：{{ '★'.repeat(o.myApplication.rating) }} {{ o.myApplication.review }}<br />结算 {{ money(o.myApplication.payout) }}</p>
          </section>
          <section v-else-if="canApply && o.status === 'open'" class="card panel">
            <h3 class="side-title">申请任务</h3>
            <textarea v-model="message" class="textarea" rows="4" placeholder="简要说明你的相关经验与交付计划"></textarea>
            <p v-if="error" class="error-text">{{ error }}</p>
            <button class="btn btn--primary btn--block" :disabled="applying" @click="apply"><Icon name="send" :size="14" /> {{ applying ? '提交中…' : '提交申请' }}</button>
          </section>

          <section v-if="o.applications?.length" class="card panel">
            <h3 class="side-title">申请人（发布者可见）</h3>
            <div v-for="a in o.applications" :key="a.id" class="row od__app">
              <UserAvatar :user="a.user" />
              <span>{{ a.user?.nickname }}<small>匹配 {{ a.match?.score ?? a.matchScore }}%</small></span>
              <span class="tag" :class="`tag--${APPLICATION_STATUS[a.status]?.tone}`">{{ APPLICATION_STATUS[a.status]?.label }}</span>
            </div>
            <router-link to="/enterprise" class="btn btn--sm btn--ghost btn--block">去企业工作台处理</router-link>
          </section>
        </aside>
      </div>
    </template>
  </div>
</template>

<style scoped>
.od .back-link { margin-bottom: 14px; }
.od__body { display: grid; grid-template-columns: minmax(0, 1fr) 320px; gap: 18px; align-items: start; }
.od h1 { margin: 12px 0 14px; font-size: 26px; }
.od h3 { margin: 22px 0 8px; font-size: 15px; }
.od__publisher { padding: 12px 14px; border-radius: 14px; background: var(--surface-2); }
.od__publisher span { display: flex; flex-direction: column; line-height: 1.35; font-size: 13.5px; }
.od__publisher small { color: var(--text-3); font-size: 12px; }
.od__budget { margin-left: auto; color: var(--warning-strong); font-size: 22px; }
.od__desc { margin: 0; color: var(--text-2); white-space: pre-wrap; }
.od__req { display: flex; align-items: center; justify-content: space-between; padding: 10px 0; border-bottom: 1px dashed var(--border); font-size: 14px; font-weight: 600; }
.od__match { display: flex; flex-direction: column; align-items: center; gap: 10px; text-align: center; }
.od__match p { margin: 0; }
.side-title { margin: 0 0 10px; font-size: 14px; }
.od__app { padding: 8px 0; font-size: 13.5px; }
.od__app span:nth-child(2) { display: flex; flex: 1; flex-direction: column; line-height: 1.35; }
.od__app small { color: var(--text-3); font-size: 11.5px; }
.stack .btn--block { margin-top: 10px; }
@media (max-width: 960px) { .od__body { grid-template-columns: 1fr; } }
</style>
