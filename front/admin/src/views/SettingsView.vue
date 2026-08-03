<template>
  <div>
    <div class="page-header">
      <h2>系统设置</h2>
      <p>查看系统运行状态。大模型通过服务端环境变量配置，不在管理端提供配置入口。</p>
    </div>

    <el-card v-loading="loading" shadow="never">
      <template #header>
        <div class="card-head">
          <span>AI 大模型状态</span>
          <el-tag :type="form.llm.enabled ? 'success' : 'danger'" size="small">
            {{ form.llm.enabled ? '已启用' : '未配置' }}
          </el-tag>
        </div>
      </template>

      <el-descriptions :column="1" border style="max-width: 640px">
        <el-descriptions-item label="配置来源">
          {{ form.llm.source === 'database' ? '数据库' : '环境变量' }}
        </el-descriptions-item>
        <el-descriptions-item label="API Base URL">
          {{ form.llm.baseUrl || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="模型名称">
          {{ form.llm.model || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="API Key">
          {{ form.llm.apiKeyConfigured ? form.llm.apiKeyMasked : '未配置' }}
        </el-descriptions-item>
      </el-descriptions>

      <el-alert
        title="大模型配置说明"
        type="info"
        :closable="false"
        show-icon
        style="margin-top: 16px; max-width: 640px"
        description="请在服务端设置环境变量 LLM_API_KEY、LLM_BASE_URL（可选，默认 OpenAI）、LLM_MODEL（可选，默认 gpt-4o-mini）后重启服务。未配置时 AI 对话等功能将不可用，系统不会使用本地模拟模型。"
      />
    </el-card>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { settingsApi } from '../api/settings.js'

const loading = ref(false)
const form = reactive({
  llm: {
    baseUrl: '',
    model: '',
    enabled: false,
    apiKeyConfigured: false,
    apiKeyMasked: '',
    source: 'environment'
  }
})

async function loadData() {
  loading.value = true
  try {
    const data = await settingsApi.get()
    form.llm.baseUrl = data.llm.baseUrl
    form.llm.model = data.llm.model
    form.llm.enabled = data.llm.enabled
    form.llm.apiKeyConfigured = data.llm.apiKeyConfigured
    form.llm.apiKeyMasked = data.llm.apiKeyMasked
    form.llm.source = data.llm.source
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
</style>
