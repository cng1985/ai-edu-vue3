<script setup>
import { computed, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { meApi, talentApi } from '../api'
import { APPLICATION_STATUS, FLYWHEEL, ROLE_NAMES, dateText, money, pct } from '../utils/eco'
import AbilityRadar from '../components/AbilityRadar.vue'
import SkillLevel from '../components/SkillLevel.vue'
import ProgressRing from '../components/ProgressRing.vue'
import UserAvatar from '../components/UserAvatar.vue'
import Icon from '../components/Icon.vue'

const route = useRoute()
const p = ref(null)
const error = ref('')
const isMe = computed(() => !route.params.id)

const ranked = computed(() => (p.value?.skills || []).filter((s) => s.level > 0 || s.knowledgeScore > 0))
const radar = computed(() =>
  ranked.value.slice(0, 8).map((s) => ({ id: s.skill.id, name: s.skill.name, progress: s.level * 20 }))
)
const heat = computed(() => {
  const days = p.value?.activity || []
  const max = Math.max(1, ...days.map((d) => d.events))
  return days.map((d) => ({ ...d, level: d.events ? Math.ceil((d.events / max) * 4) : 0 }))
})
const flywheel = computed(() => {
  const f = p.value?.flywheel || {}
  const values = {
    learn: f.learningEvents,
    practice: f.submissions,
    output: f.evidenceCount,
    income: money(f.income),
    invest: f.masteredCount,
    stronger: 'L' + (f.avgSkillLevel || 0)
  }
  return FLYWHEEL.map((s) => ({ ...s, value: values[s.key] ?? 0 }))
})

async function load() {
  error.value = ''
  try {
    p.value = isMe.value ? await meApi.profile() : await talentApi.profile(route.params.id)
  } catch (e) {
    error.value = e.message
  }
}

watch(() => route.params.id, load)
onMounted(load)
</script>

<template>
  <div class="page profile">
    <div v-if="error" class="card empty-state"><h3>{{ error }}</h3><p v-if="isMe">游客没有人才画像，注册正式账号后开始积累成长数据。</p></div>
    <template v-else-if="p">
      <section class="card panel head">
        <UserAvatar :user="p.user" size="lg" />
        <div class="head__main">
          <span class="eyebrow">Talent Profile · 人才画像</span>
          <h1>{{ p.user.nickname }} <span class="chip">{{ ROLE_NAMES[p.user.role] }}</span></h1>
          <p v-if="p.goal">职业目标：<b>{{ p.goalCareer?.name }} / {{ p.goalRole?.name }}</b> · 每周 {{ p.goal.weeklyHours }} 小时<span v-if="p.goal.motivation"> · “{{ p.goal.motivation }}”</span></p>
          <p v-else class="muted">尚未设定职业目标</p>
          <small class="muted">加入于 {{ new Date(p.joinedAt).toLocaleDateString('zh-CN') }} · 比传统简历更准确：每项能力都有学习数据、项目证据与任务表现支撑</small>
        </div>
        <div class="head__ring">
          <ProgressRing :percent="p.readiness" :size="110" :stroke="9" />
          <small>岗位达成度</small>
        </div>
      </section>

      <section class="card panel">
        <div class="flywheel">
          <div v-for="s in flywheel" :key="s.key" class="flywheel__node">
            <Icon :name="s.icon" :size="16" />
            <small>{{ s.label }}</small>
            <strong class="num">{{ s.value }}</strong>
          </div>
        </div>
      </section>

      <div class="grid-2">
        <section class="card panel">
          <div class="panel-head"><div><h2>技能水平</h2><p>技能 = 知识状态 + 项目证据 + 任务表现</p></div></div>
          <AbilityRadar v-if="radar.length >= 3" :items="radar" :size="280" />
          <p v-else class="muted small">技能数据不足，完成学习与项目后自动生成能力雷达。</p>
        </section>

        <section class="card panel">
          <div class="panel-head"><div><h2>技能明细</h2><p>知识 / 项目 / 任务三项得分</p></div></div>
          <div v-if="!ranked.length" class="muted small">暂无技能记录</div>
          <div v-for="s in ranked" :key="s.skill.id" class="skill">
            <div class="skill__head">
              <strong>{{ s.skill.name }}</strong>
              <SkillLevel :level="s.level" />
            </div>
            <div class="skill__bars">
              <span title="知识状态"><i :style="{ width: pct(s.knowledgeScore) + '%' }" class="k"></i></span>
              <span title="项目证据"><i :style="{ width: pct(s.projectScore) + '%' }" class="p"></i></span>
              <span title="任务表现"><i :style="{ width: pct(s.taskScore) + '%' }" class="t"></i></span>
            </div>
            <small class="muted">知识 {{ pct(s.knowledgeScore) }}% · 项目 {{ pct(s.projectScore) }}% · 任务 {{ pct(s.taskScore) }}% · {{ s.evidenceCount }} 项证据</small>
          </div>
        </section>
      </div>

      <div class="grid-3">
        <section class="card panel">
          <h3 class="side-title">知识掌握</h3>
          <div class="kstats">
            <div><strong class="num">{{ p.knowledge.mastered }}</strong><small>已掌握</small></div>
            <div><strong class="num">{{ p.knowledge.learning }}</strong><small>学习中</small></div>
            <div><strong class="num">{{ p.knowledge.due }}</strong><small>待复习</small></div>
          </div>
          <div v-for="d in p.knowledge.byDomain" :key="d.domain" class="mini">
            <span>{{ d.domain }} <small>{{ d.count }} 个</small></span>
            <b class="num">{{ pct(d.mastery) }}%</b>
          </div>
        </section>

        <section class="card panel">
          <h3 class="side-title">项目经历</h3>
          <p v-if="!p.projects.length" class="muted small">暂无项目提交</p>
          <router-link v-for="s in p.projects" :key="s.id" :to="`/projects/${s.projectId}?task=${s.taskId}`" class="mini">
            <span>{{ s.taskTitle }}<small>{{ s.projectTitle }} · {{ dateText(s.createdAt) }}</small></span>
            <b class="num">{{ s.score }}分</b>
          </router-link>
        </section>

        <section class="card panel">
          <h3 class="side-title">任务记录</h3>
          <p v-if="!p.tasks.length" class="muted small">暂无任务记录</p>
          <router-link v-for="a in p.tasks" :key="a.id" :to="`/market/${a.opportunityId}`" class="mini">
            <span>{{ a.opportunity?.title }}<small>{{ a.opportunity?.company }}<template v-if="a.rating"> · {{ '★'.repeat(a.rating) }}</template></small></span>
            <span class="tag" :class="`tag--${APPLICATION_STATUS[a.status]?.tone}`">{{ APPLICATION_STATUS[a.status]?.label }}</span>
          </router-link>
        </section>
      </div>

      <div class="grid-3">
        <section class="card panel">
          <h3 class="side-title">社区贡献</h3>
          <div class="kstats">
            <div><strong class="num">{{ p.contributions.posts }}</strong><small>发帖</small></div>
            <div><strong class="num">{{ p.contributions.answers }}</strong><small>回答</small></div>
            <div><strong class="num">{{ p.contributions.answerLikes }}</strong><small>获赞</small></div>
            <div><strong class="num">{{ p.contributions.resources }}</strong><small>知识资产</small></div>
          </div>
        </section>

        <section class="card panel">
          <h3 class="side-title">收入记录 <small class="muted">累计 {{ money(p.totalIncome) }}</small></h3>
          <p v-if="!p.income.length" class="muted small">完成任务市场的真实任务后获得收入</p>
          <div v-for="i in p.income" :key="i.id" class="mini">
            <span>{{ i.title }}<small>{{ dateText(i.createdAt) }}</small></span>
            <b class="num income">+{{ money(i.amount) }}</b>
          </div>
        </section>

        <section class="card panel">
          <h3 class="side-title">近 12 周学习记录</h3>
          <div class="heat">
            <span v-for="d in heat" :key="d.date" :class="`heat--${d.level}`" :title="`${d.date}：${d.events} 次学习`"></span>
          </div>
        </section>
      </div>

      <section class="card panel">
        <h3 class="side-title">能力证据（Evidence）</h3>
        <p v-if="!p.evidence.length" class="muted small">暂无证据，提交项目任务或完成真实任务即可获得</p>
        <div class="evidence">
          <div v-for="e in p.evidence.slice(0, 12)" :key="e.id" class="evidence__item">
            <span class="tag" :class="e.sourceType === 'opportunity' ? 'tag--warning' : 'tag--info'">{{ e.sourceType === 'opportunity' ? '真实任务' : '项目实践' }}</span>
            <strong>{{ e.title }}</strong>
            <small class="muted">{{ p.skills.find((s) => s.skill.id === e.skillId)?.skill.name }} · {{ e.score }} 分 · {{ dateText(e.createdAt) }}</small>
          </div>
        </div>
      </section>
    </template>
  </div>
</template>

<style scoped>
.profile { display: flex; flex-direction: column; gap: 18px; }
.head { display: flex; align-items: center; gap: 22px; }
.head__main { flex: 1; min-width: 0; }
.head__main h1 { display: flex; align-items: center; gap: 10px; margin: 4px 0; font-size: 26px; }
.head__main p { margin: 0 0 4px; color: var(--text-2); font-size: 14px; }
.head__ring { display: flex; flex-direction: column; align-items: center; gap: 6px; color: var(--text-3); font-size: 12px; }
.flywheel { display: grid; grid-template-columns: repeat(6, 1fr); gap: 10px; }
.flywheel__node { display: flex; flex-direction: column; align-items: center; gap: 2px; padding: 10px; border-radius: 12px; background: var(--surface-2); color: var(--primary); }
.flywheel__node small { color: var(--text-3); font-size: 12px; }
.flywheel__node strong { color: var(--text); font-size: 20px; }
.skill { padding: 10px 0; border-bottom: 1px dashed var(--border); }
.skill:last-child { border-bottom: none; }
.skill__head { display: flex; align-items: center; justify-content: space-between; margin-bottom: 6px; font-size: 13.5px; }
.skill__bars { display: grid; grid-template-columns: repeat(3, 1fr); gap: 6px; margin-bottom: 4px; }
.skill__bars span { height: 5px; overflow: hidden; border-radius: 3px; background: var(--surface-3); }
.skill__bars i { display: block; height: 100%; border-radius: 3px; }
.skill__bars .k { background: #2a8cf4; }
.skill__bars .p { background: var(--primary); }
.skill__bars .t { background: var(--warning); }
.skill small { font-size: 11.5px; }
.side-title { margin: 0 0 12px; font-size: 14px; }
.kstats { display: flex; justify-content: space-between; gap: 8px; margin-bottom: 10px; }
.kstats div { display: flex; flex-direction: column; align-items: center; }
.kstats strong { font-size: 22px; }
.kstats small { color: var(--text-3); font-size: 12px; }
.mini { display: flex; align-items: center; justify-content: space-between; gap: 10px; padding: 8px 0; border-bottom: 1px dashed var(--border); color: inherit; font-size: 13.5px; }
.mini > span { display: flex; flex-direction: column; min-width: 0; }
.mini small { color: var(--text-3); font-size: 11.5px; }
.income { color: var(--success-strong); }
.heat { display: grid; grid-template-columns: repeat(14, 1fr); gap: 4px; }
.heat span { aspect-ratio: 1; border-radius: 4px; background: var(--surface-3); }
.heat--1 { background: #dcd7ff !important; }
.heat--2 { background: #b3a8ff !important; }
.heat--3 { background: #8676ff !important; }
.heat--4 { background: var(--primary-strong) !important; }
.evidence { display: grid; grid-template-columns: repeat(auto-fill, minmax(240px, 1fr)); gap: 10px; }
.evidence__item { display: flex; flex-direction: column; align-items: flex-start; gap: 4px; padding: 12px; border: 1px solid var(--border); border-radius: 12px; background: var(--surface-2); font-size: 13.5px; }
@media (max-width: 960px) {
  .head { flex-wrap: wrap; }
  .flywheel { grid-template-columns: repeat(3, 1fr); }
}
</style>
