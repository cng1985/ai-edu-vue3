<template>
  <div class="page">
    <PageHeader eyebrow="任务与人才" title="人才库" subtitle="基于个人成长数据的动态人才画像：职业目标、技能水平、知识掌握、项目经历、任务记录、社区贡献、收入记录。">
      <el-input v-model="keyword" placeholder="搜索人才" clearable style="width: 220px" @keyup.enter="load" @clear="load" />
      <el-button @click="load">查询</el-button>
    </PageHeader>

    <div class="panel">
      <div class="panel__body panel__body--flush">
        <el-table :data="list" v-loading="loading" @row-click="open">
          <el-table-column label="人才" min-width="180">
            <template #default="{ row }">
              <div class="cell">
                <span class="avatar" :style="{ background: row.user.avatarColor }">{{ row.user.avatar }}</span>
                <div><strong>{{ row.user.nickname }}</strong><div class="muted small mono">{{ row.user.id }}</div></div>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="职业目标" min-width="160"><template #default="{ row }">{{ row.goal || '—' }}</template></el-table-column>
          <el-table-column label="岗位达成度" width="160">
            <template #default="{ row }"><el-progress :percentage="row.readiness" :stroke-width="8" /></template>
          </el-table-column>
          <el-table-column label="平均技能" width="100"><template #default="{ row }"><b class="num">L{{ row.avgLevel }}</b></template></el-table-column>
          <el-table-column label="代表技能" min-width="260">
            <template #default="{ row }">
              <el-tag v-for="s in row.topSkills" :key="s.skill.id" size="small" effect="plain" class="tag-gap">{{ s.skill.name }} L{{ s.level }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="证据" width="80" prop="evidenceCount" />
          <el-table-column label="任务收入" width="110"><template #default="{ row }">{{ money(row.totalIncome) }}</template></el-table-column>
        </el-table>
      </div>
    </div>

    <el-drawer v-model="drawer" :title="profile ? `人才画像 · ${profile.user.nickname}` : '人才画像'" size="620px">
      <template v-if="profile">
        <div class="kv">
          <div class="kv__item"><span class="kv__label">职业目标</span><span class="kv__value">{{ profile.goalRole ? `${profile.goalCareer?.name} / ${profile.goalRole.name}` : '未设定' }}</span></div>
          <div class="kv__item"><span class="kv__label">岗位达成度</span><span class="kv__value">{{ profile.readiness }}%</span></div>
          <div class="kv__item"><span class="kv__label">知识掌握</span><span class="kv__value">{{ profile.knowledge.mastered }} 已掌握 / {{ profile.knowledge.total }} 已学习</span></div>
          <div class="kv__item"><span class="kv__label">累计收入</span><span class="kv__value">{{ money(profile.totalIncome) }}</span></div>
          <div class="kv__item"><span class="kv__label">学习事件</span><span class="kv__value">{{ profile.flywheel.learningEvents }} 次 / {{ profile.flywheel.learningMinutes }} 分钟</span></div>
          <div class="kv__item"><span class="kv__label">社区贡献</span><span class="kv__value">{{ profile.contributions.posts }} 帖 · {{ profile.contributions.answers }} 答 · {{ profile.contributions.resources }} 资产</span></div>
        </div>
        <h4>技能水平（知识 / 项目 / 任务）</h4>
        <el-table :data="profile.skills.filter((s) => s.level > 0 || s.knowledgeScore > 0)" size="small">
          <el-table-column label="技能" prop="skill.name" min-width="140" />
          <el-table-column label="等级" width="160"><template #default="{ row }">L{{ row.level }} {{ row.levelName }}</template></el-table-column>
          <el-table-column label="知识" width="70"><template #default="{ row }">{{ pct(row.knowledgeScore) }}%</template></el-table-column>
          <el-table-column label="项目" width="70"><template #default="{ row }">{{ pct(row.projectScore) }}%</template></el-table-column>
          <el-table-column label="任务" width="70"><template #default="{ row }">{{ pct(row.taskScore) }}%</template></el-table-column>
        </el-table>
        <h4>能力证据</h4>
        <el-timeline>
          <el-timeline-item v-for="e in profile.evidence.slice(0, 10)" :key="e.id" :timestamp="dateTime(e.createdAt)" :type="e.sourceType === 'opportunity' ? 'warning' : 'primary'">
            {{ e.title }} · {{ e.score }} 分 <span class="muted">（{{ e.sourceType === 'opportunity' ? '真实任务' : '项目实践' }}）</span>
          </el-timeline-item>
        </el-timeline>
      </template>
    </el-drawer>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import PageHeader from '../../components/common/PageHeader.vue'
import { ecoAdminApi } from '../../api'
import { dateTime, money, pct } from '../../utils/eco'

const list = ref([])
const keyword = ref('')
const loading = ref(false)
const drawer = ref(false)
const profile = ref(null)

async function load() {
  loading.value = true
  try {
    list.value = await ecoAdminApi.talents(keyword.value)
  } finally {
    loading.value = false
  }
}

async function open(row) {
  profile.value = null
  drawer.value = true
  profile.value = await ecoAdminApi.talent(row.user.id)
}

onMounted(load)
</script>

<style scoped>
.page { display: flex; flex-direction: column; gap: 18px; }
.small { font-size: 12px; }
.tag-gap { margin: 2px 4px 2px 0; }
.avatar { display: inline-grid; place-items: center; width: 34px; height: 34px; border-radius: 10px; color: #fff; font-weight: 700; }
h4 { margin: 20px 0 10px; }
:deep(.el-table__row) { cursor: pointer; }
</style>
