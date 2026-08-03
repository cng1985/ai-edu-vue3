<template>
  <div>
    <div class="page-header">
      <h2>系统设置</h2>
      <p>管理 AI 大模型等运行时可调配置，保存后立即生效。</p>
    </div>

    <el-card v-loading="loading" shadow="never">
      <template #header>
        <div class="card-head">
          <span>AI 大模型配置</span>
          <el-tag :type="form.llm.enabled ? 'success' : 'info'" size="small">
            {{ form.llm.enabled ? '已启用' : '本地模式' }}
          </el-tag>
        </div>
      </template>

      <el-form label-width="120px" style="max-width: 640px">
        <el-form-item label="配置来源">
          <el-tag size="small">{{ form.llm.source === 'database' ? '数据库' : '环境变量' }}</el-tag>
          <span class="hint">数据库配置优先于环境变量</span>
        </el-form-item>
        <el-form-item label="API Base URL">
          <el-input v-model="form.llm.baseUrl" placeholder="https://api.openai.com/v1" />
        </el-form-item>
        <el-form-item label="模型名称">
          <el-input v-model="form.llm.model" placeholder="gpt-4o-mini" />
        </el-form-item>
        <el-form-item label="API Key">
          <el-input
            v-model="form.llm.apiKey"
            type="password"
            show-password
            :placeholder="form.llm.apiKeyConfigured ? `已配置 ${form.llm.apiKeyMasked}，留空不修改` : '输入 API Key 启用大模型'"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="saving" @click="handleSave">保存配置</el-button>
          <el-button @click="loadData">重置</el-button>
        </el-form-item>
      </el-form>

      <el-alert
        title="提示"
        type="info"
        :closable="false"
        show-icon
        description="服务端口、JWT 密钥、数据库路径仍通过环境变量配置。此处保存的 LLM 配置会写入 SQLite 数据库并立即生效。完整的多厂商模型路由请前往「AI 大模型配置」页面管理。"
      />
    </el-card>
  </div>
</template>

<script setup>
import { reactive, ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { settingsApi } from '../api/settings.js'

const loading = ref(false)
const saving = ref(false)
const form = reactive({
  llm: {
    apiKey: '',
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
    form.llm.apiKey = ''
  } finally {
    loading.value = false
  }
}

async function handleSave() {
  saving.value = true
  try {
    const data = await settingsApi.update({
      llm: {
        apiKey: form.llm.apiKey,
        baseUrl: form.llm.baseUrl,
        model: form.llm.model
      }
    })
    form.llm.baseUrl = data.llm.baseUrl
    form.llm.model = data.llm.model
    form.llm.enabled = data.llm.enabled
    form.llm.apiKeyConfigured = data.llm.apiKeyConfigured
    form.llm.apiKeyMasked = data.llm.apiKeyMasked
    form.llm.source = data.llm.source
    form.llm.apiKey = ''
    ElMessage.success('配置已保存并生效')
  } finally {
    saving.value = false
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
.hint {
  margin-left: 10px;
  color: #9ca3af;
  font-size: 12px;
}
</style>
