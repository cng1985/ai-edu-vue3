<template>
  <div>
    <div class="page-header">
      <h2>系统设置</h2>
      <p>查看系统运行状态。大模型通过分层配置模块管理（厂商 → 统一模型 → 虚拟模型）。</p>
    </div>

    <el-card v-loading="loading" shadow="never">
      <template #header>
        <div class="card-head">
          <span>AI 大模型状态</span>
          <el-tag :type="resolved?.enabled ? 'success' : 'warning'" size="small">
            {{ resolved?.enabled ? '路由可用' : '未就绪' }}
          </el-tag>
        </div>
      </template>

      <el-row :gutter="16" class="stats-row">
        <el-col :span="4">
          <div class="stat-label">厂商</div>
          <div class="stat-value">{{ overview.providerCount }}</div>
        </el-col>
        <el-col :span="4">
          <div class="stat-label">统一模型</div>
          <div class="stat-value">{{ overview.canonicalModelCount }}</div>
        </el-col>
        <el-col :span="4">
          <div class="stat-label">虚拟模型</div>
          <div class="stat-value">{{ overview.virtualModelCount }}</div>
        </el-col>
        <el-col :span="4">
          <div class="stat-label">厂商实现</div>
          <div class="stat-value">{{ overview.providerModelCount }}</div>
        </el-col>
        <el-col :span="8">
          <div class="stat-label">默认虚拟模型</div>
          <div class="stat-value small">{{ overview.defaultVirtualModel || '-' }}</div>
        </el-col>
      </el-row>

      <el-descriptions v-if="resolved" :column="2" border style="margin-top: 16px; max-width: 720px">
        <el-descriptions-item label="虚拟模型">{{ resolved.virtualModelCode || overview.defaultVirtualModel }}</el-descriptions-item>
        <el-descriptions-item label="统一模型">{{ resolved.canonicalModelCode }}</el-descriptions-item>
        <el-descriptions-item label="厂商">{{ resolved.providerCode }}</el-descriptions-item>
        <el-descriptions-item label="调用模型">{{ resolved.modelCode }}</el-descriptions-item>
        <el-descriptions-item label="Base URL" :span="2">{{ resolved.baseUrl }}</el-descriptions-item>
        <el-descriptions-item label="密钥状态">
          <el-tag :type="resolved.enabled ? 'success' : 'danger'" size="small">
            {{ resolved.enabled ? '已配置' : '未配置 API Key' }}
          </el-tag>
        </el-descriptions-item>
      </el-descriptions>

      <div class="actions">
        <el-button type="primary" @click="$router.push({ name: 'ai-models' })">前往大模型配置</el-button>
        <el-button @click="loadData">刷新状态</el-button>
      </div>

      <el-alert
        title="配置说明"
        type="info"
        :closable="false"
        show-icon
        style="margin-top: 16px; max-width: 720px"
        description="在「大模型配置」中管理 Provider（厂商 API Key）、CanonicalModel（统一模型）、VirtualModel（对外模型名）及映射关系。服务端口、JWT 密钥等仍通过环境变量配置。"
      />
    </el-card>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { aiModelsApi } from '../api/aiModels.js'

const loading = ref(false)
const overview = reactive({
  providerCount: 0,
  canonicalModelCount: 0,
  capabilityCount: 0,
  virtualModelCount: 0,
  providerModelCount: 0,
  defaultVirtualModel: ''
})
const resolved = ref(null)

async function loadData() {
  loading.value = true
  try {
    const data = await aiModelsApi.overview()
    Object.assign(overview, data)
    const code = data.defaultVirtualModel || 'chat-default'
    resolved.value = await aiModelsApi.resolve(code)
  } finally {
    loading.value = false
  }
}

onMounted(loadData)
</script>

<style scoped>
.page-header p {
  margin: 6px 0 0;
  color: #6b7280;
  font-size: 14px;
}
.card-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.stats-row {
  margin-top: 4px;
}
.stat-label {
  color: #6b7280;
  font-size: 13px;
}
.stat-value {
  font-size: 24px;
  font-weight: 600;
  margin-top: 4px;
}
.stat-value.small {
  font-size: 16px;
  font-weight: 500;
}
.actions {
  margin-top: 16px;
  display: flex;
  gap: 12px;
}
</style>
