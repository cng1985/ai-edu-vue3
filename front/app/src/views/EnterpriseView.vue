<script setup>
import { computed, onMounted, ref } from 'vue'
import { ecoApi, enterpriseApi } from '../api'
import { useAuthStore } from '../stores/auth'
import { APPLICATION_STATUS, OPPORTUNITY_STATUS, money, timeAgo } from '../utils/eco'
import SkillLevel from '../components/SkillLevel.vue'
import UserAvatar from '../components/UserAvatar.vue'
import Icon from '../components/Icon.vue'

const auth = useAuthStore()
const list = ref([])
const skills = ref([])
const selectedId = ref('')
const candidates = ref([])
const showForm = ref(false)
const editingId = ref('')
const form = ref(emptyForm())
const error = ref('')
const toast = ref('')
const settle = ref(null)

const selected = computed(() => list.value.find((o) => o.id === selectedId.value))
const stats = computed(() => ({
  open: list.value.filter((o) => o.status === 'open').length,
  applicants: list.value.reduce((s, o) => s + o.applications.length, 0),
  inProgress: list.value.filter((o) => o.status === 'in_progress').length,
  paid: list.value.reduce((s, o) => s + o.applications.filter((a) => a.status === 'completed').reduce((x, a) => x + a.payout, 0), 0)
}))

function emptyForm() {
  return { title: '', company: '', description: '', budget: 5000, durationDays: 14, mode: 'remote', requirements: [{ skillId: '', level: 3 }] }
}

function flash(t) {
  toast.value = t
  setTimeout(() => (toast.value = ''), 2200)
}

async function load() {
  list.value = await enterpriseApi.published()
  if (!selectedId.value && list.value.length) select(list.value[0].id)
}

async function select(id) {
  selectedId.value = id
  candidates.value = []
  candidates.value = await enterpriseApi.candidates(id).catch(() => [])
}

function openCreate() {
  editingId.value = ''
  form.value = emptyForm()
  form.value.company = auth.user?.nickname || ''
  error.value = ''
  showForm.value = true
}

function openEdit(o) {
  editingId.value = o.id
  form.value = { title: o.title, company: o.company, description: o.description, budget: o.budget, durationDays: o.durationDays, mode: o.mode, status: o.status, requirements: o.requirements.map((r) => ({ ...r })) }
  error.value = ''
  showForm.value = true
}

async function save() {
  error.value = ''
  const body = { ...form.value, requirements: form.value.requirements.filter((r) => r.skillId) }
  try {
    if (editingId.value) await enterpriseApi.update(editingId.value, body)
    else {
      const o = await enterpriseApi.publish(body)
      selectedId.value = o.id
    }
    showForm.value = false
    await load()
    await select(selectedId.value)
    flash(editingId.value ? '任务已更新' : '任务已发布，系统已为你匹配人才')
  } catch (e) {
    error.value = e.message
  }
}

async function remove(o) {
  if (!confirm(`确认删除任务「${o.title}」？`)) return
  await enterpriseApi.remove(o.id)
  selectedId.value = ''
  await load()
}

async function decide(app, action) {
  if (action === 'complete') {
    settle.value = { app, rating: 5, review: '', payout: selected.value.budget }
    return
  }
  try {
    await enterpriseApi.decide(app.id, { action })
    await load()
    flash(action === 'accept' ? '已录用，任务进入进行中' : '已婉拒')
  } catch (e) {
    flash(e.message)
  }
}

async function confirmSettle() {
  const s = settle.value
  try {
    await enterpriseApi.decide(s.app.id, { action: 'complete', rating: s.rating, review: s.review, payout: s.payout })
    settle.value = null
    await load()
    flash('已验收结算，人才获得收益与任务表现证据')
  } catch (e) {
    flash(e.message)
  }
}

onMounted(async () => {
  skills.value = await ecoApi.skills()
  await load()
})
</script>

<template>
  <div class="page ent">
    <header class="page-header page-header--split">
      <div>
        <span class="eyebrow">Enterprise Workspace</span>
        <h1>企业工作台</h1>
        <p>发布技术任务、按技能模型精准匹配人才、验收结算。比简历更准确：候选人的每项技能都由学习数据、项目证据与任务表现共同支撑。</p>
      </div>
      <button class="btn btn--primary" @click="openCreate"><Icon name="plus" :size="15" /> 发布新任务</button>
    </header>

    <div class="grid-4">
      <div class="card stat"><span class="stat__icon stat__icon--success"><Icon name="target" :size="20" /></span><div class="stat__body"><strong class="stat__value">{{ stats.open }}</strong><small class="stat__label">招募中任务</small></div></div>
      <div class="card stat"><span class="stat__icon stat__icon--info"><Icon name="users" :size="20" /></span><div class="stat__body"><strong class="stat__value">{{ stats.applicants }}</strong><small class="stat__label">累计申请</small></div></div>
      <div class="card stat"><span class="stat__icon stat__icon--warning"><Icon name="clock" :size="20" /></span><div class="stat__body"><strong class="stat__value">{{ stats.inProgress }}</strong><small class="stat__label">进行中</small></div></div>
      <div class="card stat"><span class="stat__icon"><Icon name="gift" :size="20" /></span><div class="stat__body"><strong class="stat__value">{{ money(stats.paid) }}</strong><small class="stat__label">已结算</small></div></div>
    </div>

    <div class="ent__body">
      <aside class="card panel ent__list">
        <h3 class="side-title">我发布的任务</h3>
        <p v-if="!list.length" class="muted small">还没有发布任务</p>
        <button v-for="o in list" :key="o.id" class="ent__item" :class="{ active: o.id === selectedId }" @click="select(o.id)">
          <strong>{{ o.title }}</strong>
          <small><span class="tag" :class="`tag--${OPPORTUNITY_STATUS[o.status]?.tone}`">{{ OPPORTUNITY_STATUS[o.status]?.label }}</span> {{ o.applications.length }} 人申请 · {{ money(o.budget) }}</small>
        </button>
      </aside>

      <div v-if="selected" class="stack">
        <section class="card panel">
          <div class="panel-head">
            <div><h2>{{ selected.title }}</h2><p>{{ selected.company }} · {{ timeAgo(selected.createdAt) }} · 工期 {{ selected.durationDays }} 天</p></div>
            <div class="row">
              <router-link :to="`/market/${selected.id}`" class="btn btn--sm btn--ghost">预览</router-link>
              <button class="btn btn--sm btn--ghost" @click="openEdit(selected)">编辑</button>
              <button class="btn btn--sm btn--danger" @click="remove(selected)">删除</button>
            </div>
          </div>
          <p class="muted">{{ selected.description }}</p>
          <div class="row row--wrap">
            <span v-for="r in selected.requirements" :key="r.skillId" class="chip">{{ skills.find((s) => s.id === r.skillId)?.name || r.skillId }} · L{{ r.level }}</span>
          </div>
        </section>

        <section class="card panel">
          <div class="panel-head"><div><h2>申请人</h2><p>按技能模型实时计算匹配度</p></div></div>
          <p v-if="!selected.applications.length" class="muted small">暂无申请，可以参考下方的 AI 推荐人才。</p>
          <div v-for="a in selected.applications" :key="a.id" class="app">
            <UserAvatar :user="a.user" />
            <div class="list-row__body">
              <strong>{{ a.user?.nickname }} <router-link :to="`/talents/${a.userId}`" class="small">查看人才画像</router-link></strong>
              <small>{{ a.message || '未填写申请说明' }} · {{ timeAgo(a.createdAt) }}</small>
            </div>
            <span class="app__score num">{{ a.match?.score ?? a.matchScore }}%</span>
            <span class="tag" :class="`tag--${APPLICATION_STATUS[a.status]?.tone}`">{{ APPLICATION_STATUS[a.status]?.label }}</span>
            <div class="row">
              <template v-if="a.status === 'applied'">
                <button class="btn btn--sm btn--primary" @click="decide(a, 'accept')">录用</button>
                <button class="btn btn--sm btn--ghost" @click="decide(a, 'reject')">婉拒</button>
              </template>
              <button v-else-if="a.status === 'accepted'" class="btn btn--sm btn--primary" @click="decide(a, 'complete')">验收结算</button>
              <span v-else-if="a.status === 'completed'" class="small muted">{{ '★'.repeat(a.rating) }} · {{ money(a.payout) }}</span>
            </div>
          </div>
        </section>

        <section class="card panel">
          <div class="panel-head"><div><h2><span class="tag tag--ai">AI</span> 推荐人才</h2><p>Opportunity Agent 基于人才画像匹配</p></div></div>
          <p v-if="!candidates.length" class="muted small">暂无匹配的人才</p>
          <div class="grid-2">
            <div v-for="c in candidates" :key="c.user.id" class="cand">
              <div class="row">
                <UserAvatar :user="c.user" />
                <div class="list-row__body"><strong>{{ c.user.nickname }}</strong><small>{{ c.goal || '未设定目标' }}</small></div>
                <span class="app__score num" :class="{ ok: c.match.qualified }">{{ c.match.score }}%</span>
              </div>
              <div v-for="it in c.match.items" :key="it.skillId" class="cand__req">
                <span>{{ it.skillName }}</span><SkillLevel :level="it.current" :required="it.required" />
              </div>
              <router-link :to="`/talents/${c.user.id}`" class="btn btn--sm btn--ghost btn--block">查看人才画像</router-link>
            </div>
          </div>
        </section>
      </div>
    </div>

    <div v-if="showForm" class="modal-mask" @click.self="showForm = false">
      <form class="modal" @submit.prevent="save">
        <h2>{{ editingId ? '编辑任务' : '发布 IT 任务' }}</h2>
        <p class="muted small">清晰的技能要求能让系统更精准地匹配人才，例如：Spring Boot L3、MySQL L3、Redis L2。</p>
        <div class="form-grid">
          <label class="field field--full"><span>任务标题</span><input v-model="form.title" class="input" placeholder="例如：开发库存查询模块" /></label>
          <label class="field"><span>企业名称</span><input v-model="form.company" class="input" /></label>
          <label class="field"><span>预算（元）</span><input v-model.number="form.budget" type="number" min="0" class="input" /></label>
          <label class="field"><span>工期（天）</span><input v-model.number="form.durationDays" type="number" min="1" class="input" /></label>
          <label class="field"><span>工作方式</span>
            <select v-model="form.mode" class="select"><option value="remote">远程</option><option value="onsite">驻场</option></select>
          </label>
          <label v-if="editingId" class="field"><span>状态</span>
            <select v-model="form.status" class="select"><option value="open">招募中</option><option value="in_progress">进行中</option><option value="closed">已结束</option></select>
          </label>
          <label class="field field--full"><span>需求描述</span><textarea v-model="form.description" class="textarea" rows="4"></textarea></label>
          <div class="field field--full">
            <span>技能要求</span>
            <div v-for="(r, i) in form.requirements" :key="i" class="req-row">
              <select v-model="r.skillId" class="select">
                <option value="">选择技能</option>
                <option v-for="s in skills" :key="s.id" :value="s.id">{{ s.category }} · {{ s.name }}</option>
              </select>
              <select v-model.number="r.level" class="select req-row__level">
                <option v-for="l in 5" :key="l" :value="l">L{{ l }}</option>
              </select>
              <button type="button" class="btn btn--sm btn--ghost" @click="form.requirements.splice(i, 1)"><Icon name="x" :size="13" /></button>
            </div>
            <button type="button" class="btn btn--sm btn--soft" @click="form.requirements.push({ skillId: '', level: 2 })"><Icon name="plus" :size="13" /> 添加技能要求</button>
          </div>
        </div>
        <p v-if="error" class="error-text">{{ error }}</p>
        <div class="modal__foot">
          <button type="button" class="btn btn--ghost" @click="showForm = false">取消</button>
          <button type="submit" class="btn btn--primary">{{ editingId ? '保存' : '发布' }}</button>
        </div>
      </form>
    </div>

    <div v-if="settle" class="modal-mask" @click.self="settle = null">
      <form class="modal" @submit.prevent="confirmSettle">
        <h2>验收结算 · {{ settle.app.user?.nickname }}</h2>
        <p class="muted small">评价将转化为任务表现证据，影响对方相关技能的等级；结算金额计入对方收入记录。</p>
        <div class="form-grid">
          <label class="field"><span>评分</span>
            <select v-model.number="settle.rating" class="select"><option v-for="n in 5" :key="n" :value="n">{{ '★'.repeat(n) }}</option></select>
          </label>
          <label class="field"><span>结算金额（元）</span><input v-model.number="settle.payout" type="number" min="0" class="input" /></label>
          <label class="field field--full"><span>评价</span><textarea v-model="settle.review" class="textarea" placeholder="交付质量、沟通与专业度"></textarea></label>
        </div>
        <div class="modal__foot">
          <button type="button" class="btn btn--ghost" @click="settle = null">取消</button>
          <button type="submit" class="btn btn--primary">确认验收</button>
        </div>
      </form>
    </div>
    <div v-if="toast" class="toast">{{ toast }}</div>
  </div>
</template>

<style scoped>
.ent { display: flex; flex-direction: column; gap: 18px; }
.ent .page-header { margin-bottom: 4px; }
.ent__body { display: grid; grid-template-columns: 300px minmax(0, 1fr); gap: 18px; align-items: start; }
.ent__list { position: sticky; top: 24px; }
.side-title { margin: 0 0 10px; font-size: 14px; }
.ent__item { display: flex; flex-direction: column; gap: 4px; width: 100%; margin-bottom: 6px; padding: 10px 12px; border: 1px solid transparent; border-radius: 12px; background: transparent; font: inherit; text-align: left; cursor: pointer; }
.ent__item:hover { background: var(--surface-2); }
.ent__item.active { border-color: var(--primary-soft-2); background: var(--primary-soft); }
.ent__item strong { font-size: 13.5px; }
.ent__item small { display: flex; align-items: center; gap: 6px; color: var(--text-3); font-size: 11.5px; }
.app { display: flex; align-items: center; gap: 12px; padding: 12px 0; border-bottom: 1px dashed var(--border); }
.app:last-child { border-bottom: none; }
.app__score { color: var(--primary-strong); font-size: 18px; font-weight: 700; }
.app__score.ok { color: var(--success-strong); }
.cand { display: flex; flex-direction: column; gap: 6px; padding: 14px; border: 1px solid var(--border); border-radius: 14px; background: var(--surface-2); }
.cand__req { display: flex; align-items: center; justify-content: space-between; font-size: 12.5px; }
.cand .btn { margin-top: 6px; }
.req-row { display: flex; gap: 8px; margin-bottom: 8px; }
.req-row__level { width: 90px; flex: 0 0 90px; }
@media (max-width: 960px) {
  .ent__body { grid-template-columns: 1fr; }
  .ent__list { position: static; }
  .app { flex-wrap: wrap; }
}
</style>
