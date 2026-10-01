<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ecoApi, meApi } from '../api'
import { useAuthStore } from '../stores/auth'
import { useGrowthStore } from '../stores/growth'
import { LEVELS } from '../utils/eco'
import SkillLevel from '../components/SkillLevel.vue'
import ProgressRing from '../components/ProgressRing.vue'
import Icon from '../components/Icon.vue'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const growth = useGrowthStore()

const careers = ref([])
const careerId = ref('')
const roleId = ref('')
const detail = ref(null)
const loading = ref(false)
const goal = ref(null)
const showGoal = ref(false)
const form = ref({ weeklyHours: 10, targetWeeks: 16, motivation: '' })
const saving = ref(false)
const message = ref('')

const career = computed(() => careers.value.find((c) => c.id === careerId.value))
const isMyGoal = computed(() => goal.value?.roleId === roleId.value)
const gapBySkill = computed(() => {
  const map = {}
  for (const cap of detail.value?.gap?.capabilities || []) {
    for (const s of cap.skills) map[cap.id + s.skillId] = s
  }
  return map
})

async function load() {
  careers.value = await ecoApi.careers()
  goal.value = await meApi.goal().catch(() => null)
  const target = route.query.role || goal.value?.roleId
  const owner = careers.value.find((c) => c.roles?.some((r) => r.id === target))
  careerId.value = owner?.id || careers.value[0]?.id || ''
  roleId.value = owner ? target : career.value?.roles?.[0]?.id || ''
}

async function loadRole() {
  const id = roleId.value
  if (!id) return
  loading.value = true
  try {
    const res = await ecoApi.role(id)
    if (roleId.value === id) detail.value = res
  } finally {
    loading.value = false
  }
}

function pickCareer(c) {
  careerId.value = c.id
  roleId.value = c.roles?.[0]?.id || ''
}

function openGoal() {
  form.value = {
    weeklyHours: goal.value?.weeklyHours || 10,
    targetWeeks: goal.value?.targetWeeks || 16,
    motivation: isMyGoal.value ? goal.value?.motivation || '' : ''
  }
  showGoal.value = true
}

async function saveGoal() {
  saving.value = true
  try {
    goal.value = await meApi.setGoal({ roleId: roleId.value, ...form.value })
    showGoal.value = false
    message.value = '职业目标已更新，AI 将基于新目标规划学习'
    await growth.refresh()
    setTimeout(() => router.push('/'), 900)
  } catch (e) {
    message.value = e.message
  } finally {
    saving.value = false
  }
}

watch(roleId, loadRole)
onMounted(async () => {
  await load()
  await loadRole()
})
</script>

<template>
  <div class="page career">
    <header class="page-header page-header--split">
      <div>
        <span class="eyebrow">Career · Role · Capability · Skill</span>
        <h1>职业与能力模型</h1>
        <p>每个岗位由若干「能力」组成——能完成什么类型的问题；每项能力拆解为具体「技能」及要求等级。选择目标，系统会持续分析你的差距。</p>
      </div>
      <div class="levels card card--flat">
        <span v-for="l in LEVELS" :key="l.level"><b :style="{ background: l.color }">L{{ l.level }}</b>{{ l.name }}</span>
      </div>
    </header>

    <div class="career-tabs">
      <button v-for="c in careers" :key="c.id" class="career-tab card" :class="{ active: c.id === careerId }" @click="pickCareer(c)">
        <span class="career-tab__icon" :style="{ background: c.color }"><Icon :name="c.icon" :size="18" /></span>
        <span class="career-tab__body">
          <strong>{{ c.name }}</strong>
          <small>需求{{ c.demand }} · {{ c.salaryRange }}</small>
        </span>
      </button>
    </div>

    <section v-if="career" class="card panel">
      <p class="career-desc">{{ career.description }}</p>
      <div class="segmented">
        <button v-for="r in career.roles" :key="r.id" :class="{ active: r.id === roleId }" @click="roleId = r.id">
          {{ r.name }} <small class="muted">· {{ r.level }}</small>
          <Icon v-if="goal?.roleId === r.id" name="flag" :size="12" />
        </button>
      </div>
    </section>

    <section v-if="detail" class="role-grid">
      <div class="card panel role-summary">
        <ProgressRing :percent="detail.gap?.readiness || 0" :size="132" :stroke="10" />
        <h2>{{ detail.role.name }}</h2>
        <p>{{ detail.role.description }}</p>
        <div class="role-summary__stats">
          <div><strong class="num">{{ detail.role.capabilities.length }}</strong><small>项能力</small></div>
          <div><strong class="num">{{ detail.gap?.metCount || 0 }}/{{ detail.gap?.totalCount || 0 }}</strong><small>技能达标</small></div>
        </div>
        <template v-if="auth.hasPermission('growth:write')">
          <button v-if="!isMyGoal" class="btn btn--primary btn--block" @click="openGoal"><Icon name="flag" :size="16" /> 设为我的职业目标</button>
          <template v-else>
            <span class="tag tag--success"><Icon name="check" :size="12" /> 当前职业目标</span>
            <button class="btn btn--ghost btn--block" @click="openGoal">调整学习投入</button>
          </template>
        </template>
        <p v-if="message" class="small muted">{{ message }}</p>
      </div>

      <div class="stack">
        <article v-for="cap in detail.role.capabilities" :key="cap.id" class="card panel cap">
          <div class="panel-head">
            <div>
              <h3>{{ cap.name }} <small class="muted">权重 {{ cap.weight }}</small></h3>
              <p>{{ cap.description }}</p>
            </div>
            <span class="tag" :class="(detail.gap?.capabilities.find((c) => c.id === cap.id)?.readiness || 0) >= 100 ? 'tag--success' : 'tag--neutral'">
              达成 {{ detail.gap?.capabilities.find((c) => c.id === cap.id)?.readiness || 0 }}%
            </span>
          </div>
          <div v-for="s in cap.skills" :key="s.id" class="cap__skill">
            <div class="list-row__body">
              <strong>{{ s.skill?.name || s.skillId }}</strong>
              <small>{{ s.skill?.description }}</small>
            </div>
            <SkillLevel :level="gapBySkill[cap.id + s.skillId]?.current || 0" :required="s.requiredLevel" />
            <router-link :to="`/knowledge?skill=${s.skillId}`" class="btn btn--sm btn--soft">学习</router-link>
          </div>
        </article>
      </div>
    </section>

    <div v-if="showGoal" class="modal-mask" @click.self="showGoal = false">
      <form class="modal" @submit.prevent="saveGoal">
        <h2>目标：{{ detail?.role?.name }}</h2>
        <p class="muted small">设定后，AI 学习内核会根据每周投入时间生成学习计划。</p>
        <div class="form-grid">
          <label class="field"><span>每周可投入（小时）</span><input v-model.number="form.weeklyHours" type="number" min="1" max="80" class="input" /></label>
          <label class="field"><span>目标周期（周）</span><input v-model.number="form.targetWeeks" type="number" min="1" max="104" class="input" /></label>
          <label class="field field--full"><span>成长动机（可选）</span><textarea v-model="form.motivation" class="textarea" placeholder="例如：半年内转型 AI 应用开发，接到第一个付费项目"></textarea></label>
        </div>
        <div class="modal__foot">
          <button type="button" class="btn btn--ghost" @click="showGoal = false">取消</button>
          <button type="submit" class="btn btn--primary" :disabled="saving">{{ saving ? '保存中…' : '确认目标' }}</button>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
.career { display: flex; flex-direction: column; gap: 18px; }
.career .page-header { margin-bottom: 4px; }
.levels { display: grid; grid-template-columns: repeat(2, auto); gap: 6px 16px; padding: 14px 16px; font-size: 12px; color: var(--text-2); }
.levels b { display: inline-block; min-width: 24px; margin-right: 6px; padding: 0 5px; border-radius: 6px; color: #fff; font-size: 10.5px; text-align: center; }
.career-tabs { display: grid; grid-template-columns: repeat(auto-fill, minmax(190px, 1fr)); gap: 12px; }
.career-tab { display: flex; align-items: center; gap: 12px; padding: 14px; font: inherit; text-align: left; cursor: pointer; transition: all var(--t-fast); }
.career-tab:hover { border-color: var(--border-strong); }
.career-tab.active { border-color: var(--primary); box-shadow: 0 0 0 3px var(--primary-soft-2); }
.career-tab__icon { display: grid; place-items: center; width: 38px; height: 38px; flex: 0 0 38px; border-radius: 12px; color: #fff; }
.career-tab__body { display: flex; flex-direction: column; line-height: 1.4; }
.career-tab__body strong { font-size: 14.5px; }
.career-tab__body small { color: var(--text-3); font-size: 12px; }
.career-desc { margin: 0 0 14px; color: var(--text-2); }
.segmented { flex-wrap: wrap; }
.role-grid { display: grid; grid-template-columns: 280px 1fr; gap: 18px; align-items: start; }
.role-summary { position: sticky; top: 24px; display: flex; flex-direction: column; align-items: center; gap: 10px; text-align: center; }
.role-summary h2 { margin: 6px 0 0; font-size: 20px; }
.role-summary p { margin: 0; color: var(--text-2); font-size: 13px; }
.role-summary__stats { display: flex; gap: 24px; margin: 6px 0; }
.role-summary__stats div { display: flex; flex-direction: column; }
.role-summary__stats strong { font-size: 20px; }
.role-summary__stats small { color: var(--text-3); font-size: 12px; }
.cap h3 small { font-size: 12px; font-weight: 500; }
.cap__skill { display: flex; align-items: center; gap: 14px; padding: 12px 0; border-top: 1px dashed var(--border); }
@media (max-width: 960px) {
  .role-grid { grid-template-columns: 1fr; }
  .role-summary { position: static; }
}
</style>
