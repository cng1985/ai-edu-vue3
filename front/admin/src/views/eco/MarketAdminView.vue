<template>
  <div class="page">
    <PageHeader eyebrow="任务与人才" title="IT 任务市场" subtitle="企业发布 → 技能匹配 → 推荐人才 → 录用 → 验收结算。平台可代为处理履约与治理异常任务。">
      <el-select v-model="status" placeholder="全部状态" clearable style="width: 140px" @change="load">
        <el-option v-for="(s, k) in OPPORTUNITY_STATUS" :key="k" :label="s.label" :value="k" />
      </el-select>
      <el-button :icon="Refresh" @click="load">刷新</el-button>
    </PageHeader>

    <div class="stat-grid">
      <StatCard label="任务总数" :value="list.length" icon="Suitcase" tone="primary" />
      <StatCard label="招募中" :value="list.filter((o) => o.status === 'open').length" icon="Promotion" tone="success" />
      <StatCard label="申请总数" :value="list.reduce((s, o) => s + o.applications.length, 0)" icon="User" tone="info" />
      <StatCard label="已结算金额" :value="money(paid)" icon="Money" tone="warning" />
    </div>

    <div class="panel">
      <div class="panel__body panel__body--flush">
        <el-table :data="list" v-loading="loading" row-key="id">
          <el-table-column type="expand">
            <template #default="{ row }">
              <div class="expand">
                <p class="muted">{{ row.description }}</p>
                <el-table :data="row.applications" size="small" empty-text="暂无申请">
                  <el-table-column label="人才" min-width="140"><template #default="{ row: a }">{{ a.user?.nickname }} <span class="muted mono">{{ a.userId }}</span></template></el-table-column>
                  <el-table-column label="匹配度" width="90"><template #default="{ row: a }"><span class="num">{{ a.match?.score ?? a.matchScore }}%</span></template></el-table-column>
                  <el-table-column label="申请说明" prop="message" min-width="200" show-overflow-tooltip />
                  <el-table-column label="状态" width="100"><template #default="{ row: a }"><el-tag size="small" :type="APPLICATION_STATUS[a.status]?.type">{{ APPLICATION_STATUS[a.status]?.label }}</el-tag></template></el-table-column>
                  <el-table-column label="评价 / 结算" width="160"><template #default="{ row: a }"><span v-if="a.status === 'completed'">{{ '★'.repeat(a.rating) }} · {{ money(a.payout) }}</span></template></el-table-column>
                  <el-table-column label="操作" width="200" align="right">
                    <template #default="{ row: a }">
                      <template v-if="a.status === 'applied'">
                        <el-button link type="primary" @click="decide(a, 'accept')">录用</el-button>
                        <el-button link @click="decide(a, 'reject')">婉拒</el-button>
                      </template>
                      <el-button v-else-if="a.status === 'accepted'" link type="primary" @click="openSettle(a, row)">验收结算</el-button>
                    </template>
                  </el-table-column>
                </el-table>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="任务" min-width="240">
            <template #default="{ row }"><strong>{{ row.title }}</strong><div class="muted small">{{ row.company }} · {{ row.publisher?.nickname }}</div></template>
          </el-table-column>
          <el-table-column label="技能要求" min-width="240">
            <template #default="{ row }">
              <el-tag v-for="r in row.requirements" :key="r.skillId" size="small" effect="plain" class="tag-gap">{{ skillName(r.skillId) }} L{{ r.level }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="预算" width="110"><template #default="{ row }">{{ money(row.budget) }}</template></el-table-column>
          <el-table-column label="申请" width="70"><template #default="{ row }">{{ row.applications.length }}</template></el-table-column>
          <el-table-column label="状态" width="130">
            <template #default="{ row }">
              <el-select :model-value="row.status" size="small" @change="(v) => changeStatus(row, v)">
                <el-option v-for="(s, k) in OPPORTUNITY_STATUS" :key="k" :label="s.label" :value="k" />
              </el-select>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="160" align="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="showCandidates(row)">推荐人才</el-button>
              <el-popconfirm title="确定删除该任务？" @confirm="remove(row.id)">
                <template #reference><el-button link type="danger">删除</el-button></template>
              </el-popconfirm>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </div>

    <el-drawer v-model="drawer.visible" :title="`推荐人才 · ${drawer.title}`" size="520px">
      <el-empty v-if="!drawer.list.length" description="暂无匹配人才" />
      <div v-for="c in drawer.list" :key="c.user.id" class="cand">
        <div class="cand__head">
          <strong>{{ c.user.nickname }}</strong>
          <span class="muted">{{ c.goal || '未设定目标' }}</span>
          <el-tag :type="c.match.qualified ? 'success' : 'primary'" class="cand__score">{{ c.match.score }}%</el-tag>
        </div>
        <div v-for="it in c.match.items" :key="it.skillId" class="cand__req">
          <span>{{ it.skillName }}</span>
          <span :class="it.gap ? 'gap' : 'ok'">L{{ it.current }} / 需 L{{ it.required }}</span>
        </div>
      </div>
    </el-drawer>

    <el-dialog v-model="settle.visible" title="验收结算" width="460px">
      <el-form label-width="80px">
        <el-form-item label="评分"><el-rate v-model="settle.rating" /></el-form-item>
        <el-form-item label="结算金额"><el-input-number v-model="settle.payout" :min="0" :step="500" /></el-form-item>
        <el-form-item label="评价"><el-input v-model="settle.review" type="textarea" :rows="3" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="settle.visible = false">取消</el-button>
        <el-button type="primary" @click="confirmSettle">确认验收</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Refresh } from '@element-plus/icons-vue'
import PageHeader from '../../components/common/PageHeader.vue'
import StatCard from '../../components/common/StatCard.vue'
import { ecoAdminApi } from '../../api'
import { APPLICATION_STATUS, OPPORTUNITY_STATUS, money } from '../../utils/eco'

const list = ref([])
const skills = ref([])
const status = ref('')
const loading = ref(false)
const drawer = reactive({ visible: false, title: '', list: [] })
const settle = reactive({ visible: false, app: null, rating: 5, payout: 0, review: '' })

const paid = computed(() => list.value.reduce((s, o) => s + o.applications.filter((a) => a.status === 'completed').reduce((x, a) => x + a.payout, 0), 0))
const skillName = (id) => skills.value.find((s) => s.id === id)?.name || id

async function load() {
  loading.value = true
  try {
    list.value = await ecoAdminApi.opportunities({ status: status.value })
  } finally {
    loading.value = false
  }
}

async function changeStatus(row, value) {
  await ecoAdminApi.updateOpportunity(row.id, {
    title: row.title, company: row.company, description: row.description, requirements: row.requirements,
    budget: row.budget, durationDays: row.durationDays, mode: row.mode, status: value
  })
  ElMessage.success('状态已更新')
  await load()
}

async function remove(id) {
  await ecoAdminApi.removeOpportunity(id)
  ElMessage.success('已删除')
  await load()
}

async function showCandidates(row) {
  drawer.title = row.title
  drawer.list = await ecoAdminApi.candidates(row.id)
  drawer.visible = true
}

async function decide(app, action) {
  await ecoAdminApi.decide(app.id, { action })
  ElMessage.success('已处理')
  await load()
}

function openSettle(app, opp) {
  Object.assign(settle, { visible: true, app, rating: 5, payout: opp.budget, review: '' })
}

async function confirmSettle() {
  await ecoAdminApi.decide(settle.app.id, { action: 'complete', rating: settle.rating, payout: settle.payout, review: settle.review })
  ElMessage.success('已验收结算')
  settle.visible = false
  await load()
}

onMounted(async () => {
  skills.value = await ecoAdminApi.skills()
  await load()
})
</script>

<style scoped>
.page { display: flex; flex-direction: column; gap: 18px; }
.small { font-size: 12.5px; }
.tag-gap { margin: 2px 4px 2px 0; }
.expand { padding: 4px 24px 12px 48px; }
.cand { margin-bottom: 14px; padding: 12px 14px; border: 1px solid var(--border); border-radius: 12px; }
.cand__head { display: flex; align-items: center; gap: 10px; margin-bottom: 8px; }
.cand__score { margin-left: auto; }
.cand__req { display: flex; justify-content: space-between; font-size: 13px; padding: 2px 0; }
.cand__req .gap { color: var(--warning-strong); }
.cand__req .ok { color: var(--success-strong); }
</style>
