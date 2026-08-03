<template>
  <div>
    <div class="page-header">
      <h2>系统设置</h2>
      <p>管理大模型分层配置（厂商 → 统一模型 → 虚拟模型）与知识库索引。</p>
    </div>

    <!-- 未配置时：一键初始化向导 -->
    <el-card v-if="!loading && overview.providerCount === 0" shadow="never" class="section-card">
      <template #header>
        <span>快速初始化大模型</span>
      </template>
      <p class="wizard-hint">检测到尚未配置厂商，可通过向导一键创建完整路由链路。</p>
      <el-form label-width="120px" style="max-width: 640px">
        <el-form-item label="预设模板">
          <el-radio-group v-model="preset" @change="applyPreset">
            <el-radio-button value="openai">OpenAI</el-radio-button>
            <el-radio-button value="deepseek">DeepSeek</el-radio-button>
            <el-radio-button value="custom">自定义</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="厂商编码" required>
          <el-input v-model="setupForm.providerCode" placeholder="openai" />
        </el-form-item>
        <el-form-item label="厂商名称" required>
          <el-input v-model="setupForm.providerName" placeholder="OpenAI" />
        </el-form-item>
        <el-form-item label="Base URL">
          <el-input v-model="setupForm.baseUrl" placeholder="https://api.openai.com/v1" />
        </el-form-item>
        <el-form-item label="API Key" required>
          <el-input v-model="setupForm.apiKey" type="password" show-password placeholder="sk-..." />
        </el-form-item>
        <el-form-item label="模型编码" required>
          <el-input v-model="setupForm.canonicalCode" placeholder="gpt-4o-mini" />
        </el-form-item>
        <el-form-item label="模型名称">
          <el-input v-model="setupForm.canonicalName" placeholder="GPT-4o Mini" />
        </el-form-item>
        <el-form-item label="虚拟模型">
          <el-input v-model="setupForm.virtualCode" placeholder="chat-default" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="settingUp" @click="handleQuickSetup">
            一键初始化
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <!-- 配置健康检查 -->
    <el-card v-loading="loading" shadow="never" class="section-card">
      <template #header>
        <div class="card-head">
          <span>大模型配置检查</span>
          <el-tag :type="allReady ? 'success' : 'warning'" size="small">
            {{ allReady ? '可正常使用' : '需要完善配置' }}
          </el-tag>
        </div>
      </template>

      <el-steps :active="readyStep" finish-status="success" align-center class="check-steps">
        <el-step title="厂商" :description="`${overview.providerCount} 个`" />
        <el-step title="密钥" :description="`${providersWithKey} 个已配置`" />
        <el-step title="虚拟模型" :description="`${overview.virtualModelCount} 个`" />
        <el-step title="路由可用" :description="resolved?.enabled ? '已连通' : '未连通'" />
      </el-steps>

      <el-alert
        v-if="!allReady && overview.providerCount > 0"
        :title="checklistHint"
        type="warning"
        :closable="false"
        show-icon
        style="margin-top: 16px"
      />
    </el-card>

    <!-- 默认虚拟模型 -->
    <el-card shadow="never" class="section-card">
      <template #header>
        <span>默认虚拟模型</span>
      </template>
      <el-form label-width="120px" style="max-width: 560px">
        <el-form-item label="虚拟模型">
          <el-select v-model="defaultModel" filterable placeholder="选择默认虚拟模型" style="width: 100%">
            <el-option
              v-for="vm in virtualModels"
              :key="vm.code"
              :label="`${vm.code} - ${vm.name}`"
              :value="vm.code"
            />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="savingDefault" :disabled="!defaultModel" @click="saveDefaultModel">
            保存为默认
          </el-button>
          <el-button :loading="resolving" @click="testResolve">测试路由解析</el-button>
        </el-form-item>
      </el-form>

      <el-descriptions v-if="resolved" :column="2" border style="max-width: 720px">
        <el-descriptions-item label="虚拟模型">{{ resolved.virtualModelCode }}</el-descriptions-item>
        <el-descriptions-item label="统一模型">{{ resolved.canonicalModelCode }}</el-descriptions-item>
        <el-descriptions-item label="厂商">{{ resolved.providerCode }}</el-descriptions-item>
        <el-descriptions-item label="调用模型">{{ resolved.modelCode }}</el-descriptions-item>
        <el-descriptions-item label="Base URL" :span="2">{{ resolved.baseUrl }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="resolved.enabled ? 'success' : 'danger'" size="small">
            {{ resolved.enabled ? '可用' : '未配置 API Key' }}
          </el-tag>
        </el-descriptions-item>
      </el-descriptions>
    </el-card>

    <!-- 厂商密钥 -->
    <el-card shadow="never" class="section-card">
      <template #header>
        <div class="card-head">
          <span>厂商 API Key</span>
          <el-button type="primary" size="small" @click="openProviderDialog()">新增厂商</el-button>
        </div>
      </template>

      <el-table :data="providers" stripe>
        <el-table-column prop="code" label="编码" width="120" />
        <el-table-column prop="name" label="名称" width="140" />
        <el-table-column prop="baseUrl" label="Base URL" show-overflow-tooltip />
        <el-table-column label="API Key" width="140">
          <template #default="{ row }">
            <el-tag v-if="row.apiKeyMasked" type="success" size="small">{{ row.apiKeyMasked }}</el-tag>
            <el-tag v-else type="danger" size="small">未配置</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
              {{ row.status === 1 ? '启用' : '禁用' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="openProviderDialog(row)">编辑</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="card-footer">
        <el-button link type="primary" @click="$router.push({ name: 'ai-models' })">
          前往完整大模型配置（统一模型、虚拟映射、能力标签）→
        </el-button>
      </div>
    </el-card>

    <!-- 知识库 -->
    <el-card shadow="never" class="section-card">
      <template #header>
        <div class="card-head">
          <span>知识库索引</span>
          <el-tag :type="kbStatus?.indexStatus === 'ready' ? 'success' : 'info'" size="small">
            {{ kbStatus?.indexStatus || '未知' }}
          </el-tag>
        </div>
      </template>

      <el-descriptions v-if="kbStatus" :column="3" border>
        <el-descriptions-item label="课程数">{{ kbStatus.courseCount }}</el-descriptions-item>
        <el-descriptions-item label="章节数">{{ kbStatus.chapterCount }}</el-descriptions-item>
        <el-descriptions-item label="切片数">{{ kbStatus.chunkCount }}</el-descriptions-item>
        <el-descriptions-item label="嵌入模型">{{ kbStatus.embedModel }}</el-descriptions-item>
        <el-descriptions-item label="嵌入来源">{{ kbStatus.embedSource === 'api' ? 'API' : '本地哈希' }}</el-descriptions-item>
        <el-descriptions-item label="向量维度">{{ kbStatus.dimensions }}</el-descriptions-item>
      </el-descriptions>

      <div class="actions" style="margin-top: 12px">
        <el-button type="primary" :loading="reindexing" @click="handleReindex">重建索引</el-button>
        <el-button @click="loadData">刷新</el-button>
      </div>
    </el-card>

    <!-- 厂商编辑对话框 -->
    <el-dialog v-model="providerDialog.visible" :title="providerDialog.isEdit ? '编辑厂商' : '新增厂商'" width="520px">
      <el-form label-width="100px">
        <el-form-item label="编码" required>
          <el-input v-model="providerDialog.form.code" :disabled="providerDialog.isEdit" />
        </el-form-item>
        <el-form-item label="名称" required>
          <el-input v-model="providerDialog.form.name" />
        </el-form-item>
        <el-form-item label="Base URL">
          <el-input v-model="providerDialog.form.baseUrl" />
        </el-form-item>
        <el-form-item label="认证类型">
          <el-input v-model="providerDialog.form.authType" placeholder="Bearer" />
        </el-form-item>
        <el-form-item label="API Key">
          <el-input v-model="providerDialog.form.apiKey" type="password" show-password :placeholder="providerDialog.isEdit ? '留空不修改' : '输入 API Key'" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="providerDialog.form.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="providerDialog.visible = false">取消</el-button>
        <el-button type="primary" :loading="providerDialog.saving" @click="saveProvider">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { settingsApi } from '../api/settings.js'

const PRESETS = {
  openai: {
    providerCode: 'openai',
    providerName: 'OpenAI',
    baseUrl: 'https://api.openai.com/v1',
    canonicalCode: 'gpt-4o-mini',
    canonicalName: 'GPT-4o Mini',
    virtualCode: 'chat-default',
    virtualName: '默认对话模型'
  },
  deepseek: {
    providerCode: 'deepseek',
    providerName: 'DeepSeek',
    baseUrl: 'https://api.deepseek.com/v1',
    canonicalCode: 'deepseek-chat',
    canonicalName: 'DeepSeek Chat',
    virtualCode: 'chat-default',
    virtualName: '默认对话模型'
  },
  custom: {
    providerCode: '',
    providerName: '',
    baseUrl: '',
    canonicalCode: '',
    canonicalName: '',
    virtualCode: 'chat-default',
    virtualName: '默认对话模型'
  }
}

const loading = ref(false)
const settingUp = ref(false)
const savingDefault = ref(false)
const resolving = ref(false)
const reindexing = ref(false)
const preset = ref('openai')

const overview = reactive({
  providerCount: 0,
  canonicalModelCount: 0,
  virtualModelCount: 0,
  defaultVirtualModel: ''
})
const providers = ref([])
const virtualModels = ref([])
const defaultModel = ref('')
const resolved = ref(null)
const kbStatus = ref(null)

const setupForm = reactive({
  providerCode: 'openai',
  providerName: 'OpenAI',
  baseUrl: 'https://api.openai.com/v1',
  apiKey: '',
  canonicalCode: 'gpt-4o-mini',
  canonicalName: 'GPT-4o Mini',
  modelCode: '',
  virtualCode: 'chat-default',
  virtualName: '默认对话模型'
})

const providersWithKey = computed(() => providers.value.filter((p) => p.apiKeyMasked).length)
const allReady = computed(() =>
  overview.providerCount > 0 && providersWithKey.value > 0 && overview.virtualModelCount > 0 && resolved.value?.enabled
)
const readyStep = computed(() => {
  if (resolved.value?.enabled) return 4
  if (overview.virtualModelCount > 0) return 3
  if (providersWithKey.value > 0) return 2
  if (overview.providerCount > 0) return 1
  return 0
})
const checklistHint = computed(() => {
  if (providersWithKey.value === 0) return '请为厂商配置 API Key'
  if (overview.virtualModelCount === 0) return '请创建虚拟模型，或使用「一键初始化」'
  if (!resolved.value?.enabled) return '路由未连通，请检查厂商模型映射与 API Key'
  return ''
})

const providerDialog = reactive({
  visible: false,
  isEdit: false,
  saving: false,
  form: emptyProvider()
})

function emptyProvider() {
  return { code: '', name: '', baseUrl: '', authType: 'Bearer', apiKey: '', status: 1 }
}

function applyPreset(val) {
  const p = PRESETS[val] || PRESETS.custom
  Object.assign(setupForm, { ...p, apiKey: setupForm.apiKey })
}

function applyView(data) {
  Object.assign(overview, data.aiModel)
  providers.value = data.providers || []
  virtualModels.value = data.virtualModels || []
  resolved.value = data.resolved || null
  kbStatus.value = data.knowledge || null
  defaultModel.value = data.aiModel?.defaultVirtualModel || ''
}

async function loadData() {
  loading.value = true
  try {
    const data = await settingsApi.get()
    applyView(data)
  } finally {
    loading.value = false
  }
}

async function handleQuickSetup() {
  settingUp.value = true
  try {
    const data = await settingsApi.quickSetup({
      ...setupForm,
      modelCode: setupForm.canonicalCode
    })
    applyView(data)
    ElMessage.success('大模型配置已初始化，路由链路已建立')
  } finally {
    settingUp.value = false
  }
}

async function saveDefaultModel() {
  savingDefault.value = true
  try {
    const data = await settingsApi.setDefaultVirtualModel(defaultModel.value)
    applyView(data)
    ElMessage.success('默认虚拟模型已更新')
  } finally {
    savingDefault.value = false
  }
}

async function testResolve() {
  resolving.value = true
  try {
    const code = defaultModel.value || overview.defaultVirtualModel || 'chat-default'
    resolved.value = await settingsApi.resolve(code)
    ElMessage.success(resolved.value?.enabled ? '路由解析成功' : '路由已解析，但 API Key 未配置')
  } finally {
    resolving.value = false
  }
}

function openProviderDialog(row) {
  providerDialog.isEdit = !!row
  providerDialog.form = row ? { ...row, apiKey: '' } : emptyProvider()
  providerDialog.visible = true
}

async function saveProvider() {
  providerDialog.saving = true
  try {
    if (providerDialog.isEdit) {
      await settingsApi.updateProvider(providerDialog.form.id, providerDialog.form)
    } else {
      await settingsApi.createProvider(providerDialog.form)
    }
    providerDialog.visible = false
    ElMessage.success('保存成功')
    await loadData()
  } finally {
    providerDialog.saving = false
  }
}

async function handleReindex() {
  reindexing.value = true
  try {
    kbStatus.value = await settingsApi.reindexKnowledge()
    ElMessage.success('知识库索引重建完成')
  } finally {
    reindexing.value = false
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
.section-card {
  margin-bottom: 16px;
}
.card-head {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.check-steps {
  margin-top: 8px;
}
.wizard-hint {
  color: #6b7280;
  margin: 0 0 16px;
  font-size: 14px;
}
.card-footer {
  margin-top: 12px;
}
.actions {
  display: flex;
  gap: 12px;
}
</style>
