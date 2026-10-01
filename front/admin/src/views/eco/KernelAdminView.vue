<template>
  <div class="page">
    <PageHeader eyebrow="AI 服务" title="AI 学习内核" subtitle="AI Learning Kernel → AI Execution Kernel → Pipeline Runtime → Resource / Tool Layer。查看 Pipeline 结构、Agent 体系与每次运行的 Stage 轨迹。">
      <el-button :icon="Refresh" @click="load">刷新</el-button>
    </PageHeader>

    <div class="grid">
      <section class="panel">
        <div class="panel__head"><h3 class="panel__title">学习 Pipeline</h3></div>
        <div class="panel__body">
          <el-steps direction="vertical" :active="stages.length">
            <el-step v-for="s in stages" :key="s.name" :title="s.title" :description="s.name" />
          </el-steps>
        </div>
      </section>

      <section class="panel">
        <div class="panel__head"><h3 class="panel__title">Agent 体系</h3></div>
        <div class="panel__body agents">
          <div v-for="a in agents" :key="a.code" class="agent">
            <span class="agent__dot" :style="{ background: a.color }" />
            <div>
              <strong>{{ a.name }}</strong> <span class="muted">{{ a.title }}</span>
              <p class="muted">{{ a.description }}</p>
              <el-tag v-for="r in a.responsibilities" :key="r" size="small" effect="plain" class="tag-gap">{{ r }}</el-tag>
            </div>
          </div>
        </div>
      </section>
    </div>

    <div class="panel">
      <div class="panel__head">
        <div><h3 class="panel__title">运行记录</h3><p class="panel__sub">最近 {{ runs.length }} 次学习 Pipeline 运行</p></div>
      </div>
      <div class="panel__body panel__body--flush">
        <el-table :data="runs" v-loading="loading" row-key="id">
          <el-table-column type="expand">
            <template #default="{ row }">
              <div class="expand">
                <el-table :data="row.trace" size="small">
                  <el-table-column label="Stage" width="220"><template #default="{ row: t }"><strong>{{ t.title }}</strong> <span class="muted mono">{{ t.name }}</span></template></el-table-column>
                  <el-table-column label="状态" width="90">
                    <template #default="{ row: t }"><el-tag size="small" :type="t.status === 'ok' ? 'success' : t.status === 'skipped' ? 'warning' : 'danger'">{{ t.status }}</el-tag></template>
                  </el-table-column>
                  <el-table-column label="摘要" prop="summary" min-width="320" />
                  <el-table-column label="耗时" width="80"><template #default="{ row: t }">{{ t.durationMs }}ms</template></el-table-column>
                </el-table>
                <p v-if="row.output?.response" class="response"><b>AI 回答（{{ row.output.responseSource === 'ai' ? '大模型' : '规则模板' }}）：</b>{{ row.output.response }}</p>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="运行时间" width="190"><template #default="{ row }">{{ dateTime(row.createdAt) }}</template></el-table-column>
          <el-table-column label="用户" prop="userId" width="160" />
          <el-table-column label="Pipeline" prop="pipeline" width="110" />
          <el-table-column label="状态" width="100"><template #default="{ row }"><el-tag size="small" :type="row.status === 'ok' ? 'success' : 'danger'">{{ row.status }}</el-tag></template></el-table-column>
          <el-table-column label="岗位达成度" width="110"><template #default="{ row }">{{ row.output?.gap ? row.output.gap.readiness + '%' : '—' }}</template></el-table-column>
          <el-table-column label="推荐知识" width="100"><template #default="{ row }">{{ row.output?.recommendations?.length || 0 }}</template></el-table-column>
          <el-table-column label="耗时" width="100"><template #default="{ row }">{{ row.durationMs }}ms</template></el-table-column>
        </el-table>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { Refresh } from '@element-plus/icons-vue'
import PageHeader from '../../components/common/PageHeader.vue'
import { ecoAdminApi } from '../../api'
import { dateTime } from '../../utils/eco'

const stages = ref([])
const agents = ref([])
const runs = ref([])
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    const [s, a, r] = await Promise.all([ecoAdminApi.kernelStages(), ecoAdminApi.agents(), ecoAdminApi.kernelRuns({ limit: 50 })])
    stages.value = s
    agents.value = a
    runs.value = r
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<style scoped>
.page { display: flex; flex-direction: column; gap: 18px; }
.grid { display: grid; grid-template-columns: 360px minmax(0, 1fr); gap: 18px; }
.agents { display: grid; grid-template-columns: 1fr 1fr; gap: 14px; }
.agent { display: flex; gap: 10px; }
.agent__dot { width: 10px; height: 10px; flex: 0 0 10px; margin-top: 6px; border-radius: 50%; }
.agent p { margin: 4px 0 6px; font-size: 13px; }
.tag-gap { margin: 2px 4px 2px 0; }
.expand { padding: 4px 24px 12px 48px; }
.response { margin: 12px 0 0; padding: 10px 12px; border-radius: 8px; background: var(--surface-2); font-size: 13px; white-space: pre-wrap; }
@media (max-width: 1100px) {
  .grid, .agents { grid-template-columns: 1fr; }
}
</style>
