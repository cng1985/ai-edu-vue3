<template>
  <div class="page">
    <PageHeader eyebrow="系统管理" title="系统设置" subtitle="管理大模型分层配置（厂商 → 统一模型 → 虚拟模型）与知识库索引。">
      <el-button :icon="Refresh" :loading="loading" @click="loadData">刷新</el-button>
    </PageHeader>

    <!-- 配置健康检查 -->
    <div v-loading="loading" class="panel health" :class="{ 'health--ok': allReady }">
      <div class="panel__head">
        <h3 class="panel__title">
          <span class="panel__title-icon"><el-icon><component :is="allReady ? CircleCheck : Warning" /></el-icon></span>
          大模型配置检查
        </h3>
        <el-tag :type="allReady ? 'success' : 'warning'" size="small">
          {{ allReady ? '可正常使用' : '需要完善配置' }}
        </el-tag>
      </div>
      <div class="panel__body">
        <div class="steps">
          <div
            v-for="(s, i) in checkSteps"
            :key="s.title"
            class="step"
            :class="{ 'step--done': i < readyStep, 'step--current': i === readyStep }"
          >
            <span class="step__index">
              <el-icon v-if="i < readyStep" :size="14"><Check /></el-icon>
              <span v-else class="num">{{ i + 1 }}</span>
            </span>
            <div class="step__text">
              <strong>{{ s.title }}</strong>
              <span>{{ s.desc }}</span>
            </div>
            <span v-if="i < checkSteps.length - 1" class="step__line" />
          </div>
        </div>
        <div v-if="!allReady && overview.providerCount > 0" class="hint hint--warning" style="margin-top: 18px">
          <el-icon :size="16" style="margin-top: 2px"><Warning /></el-icon>
          <span>{{ checklistHint }}</span>
        </div>
      </div>
    </div>

    <!-- 一键初始化/补全向导 -->
    <div v-if="!loading" class="panel">
      <div class="panel__head">
        <div>
          <h3 class="panel__title">
            <span class="panel__title-icon"><el-icon><MagicStick /></el-icon></span>
            快速初始化大模型
          </h3>
          <p class="panel__sub">
            {{ overview.providerCount === 0 ? '尚未配置厂商，可通过向导创建完整路由链路。' : '可新增或补全厂商 → 统一模型 → 厂商模型 → 虚拟模型的完整链路。' }}
          </p>
        </div>
      </div>
      <div class="panel__body">
      <el-form label-width="120px" style="max-width: 720px">
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
        <el-form-item label="厂商模型标识" required>
          <el-input v-model="setupForm.modelCode" placeholder="留空时使用模型编码" />
        </el-form-item>
        <el-form-item label="上下文窗口">
          <el-input-number v-model="setupForm.contextWindow" :min="1" :step="1000" />
        </el-form-item>
        <el-form-item label="推理模型">
          <el-switch v-model="setupForm.reasoningSupported" />
        </el-form-item>
        <el-form-item label="虚拟模型">
          <el-input v-model="setupForm.virtualCode" placeholder="chat-default" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="settingUp" @click="handleQuickSetup">
            <el-icon class="el-icon--left"><MagicStick /></el-icon>一键初始化
          </el-button>
        </el-form-item>
      </el-form>
      </div>
    </div>

    <div class="grid-2">
    <!-- 默认虚拟模型 -->
    <div class="panel">
      <div class="panel__head">
        <h3 class="panel__title">
          <span class="panel__title-icon"><el-icon><Connection /></el-icon></span>
          默认虚拟模型
        </h3>
      </div>
      <div class="panel__body">
      <el-form label-position="top">
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

      <div v-if="resolved" class="route">
        <div class="route__chain">
          <span class="route__node route__node--virtual"><small>虚拟模型</small><b class="mono">{{ resolved.virtualModelCode }}</b></span>
          <el-icon class="route__arrow"><Right /></el-icon>
          <span class="route__node"><small>统一模型</small><b class="mono">{{ resolved.canonicalModelCode }}</b></span>
          <el-icon class="route__arrow"><Right /></el-icon>
          <span class="route__node"><small>厂商</small><b class="mono">{{ resolved.providerCode }}</b></span>
          <el-icon class="route__arrow"><Right /></el-icon>
          <span class="route__node route__node--model"><small>调用模型</small><b class="mono">{{ resolved.modelCode }}</b></span>
        </div>
        <div class="kv" style="margin-top: 14px">
          <div class="kv__item"><span class="kv__label">部署名</span><span class="kv__value mono">{{ resolved.deploymentName || '—' }}</span></div>
          <div class="kv__item"><span class="kv__label">上下文窗口</span><span class="kv__value num">{{ resolved.contextWindow || '—' }}</span></div>
          <div class="kv__item"><span class="kv__label">推理能力</span><span class="kv__value">{{ resolved.reasoningSupported ? '支持' : '普通模型' }}</span></div>
          <div class="kv__item">
            <span class="kv__label">状态</span>
            <span class="kv__value">
              <el-tag :type="resolved.enabled ? 'success' : 'danger'" size="small">{{ resolved.enabled ? '可用' : '未配置 API Key' }}</el-tag>
            </span>
          </div>
          <div class="kv__item" style="grid-column: 1 / -1"><span class="kv__label">Base URL</span><span class="kv__value mono">{{ resolved.baseUrl }}</span></div>
        </div>
      </div>
      </div>
    </div>

    <!-- 知识库 -->
    <div class="panel">
      <div class="panel__head">
        <h3 class="panel__title">
          <span class="panel__title-icon"><el-icon><Collection /></el-icon></span>
          知识库索引
        </h3>
        <el-tag :type="kbStatus?.indexStatus === 'ready' ? 'success' : 'info'" size="small">
          {{ kbStatus?.indexStatus || '未知' }}
        </el-tag>
      </div>
      <div class="panel__body">
        <div v-if="kbStatus" class="kb">
          <div class="kb__stat"><span>课程数</span><strong class="num">{{ kbStatus.courseCount }}</strong></div>
          <div class="kb__stat"><span>章节数</span><strong class="num">{{ kbStatus.chapterCount }}</strong></div>
          <div class="kb__stat"><span>切片数</span><strong class="num">{{ kbStatus.chunkCount }}</strong></div>
        </div>
        <div v-if="kbStatus" class="kv" style="margin-top: 14px">
          <div class="kv__item"><span class="kv__label">嵌入模型</span><span class="kv__value mono">{{ kbStatus.embedModel }}</span></div>
          <div class="kv__item"><span class="kv__label">嵌入来源</span><span class="kv__value">{{ kbStatus.embedSource === 'api' ? 'API' : '本地哈希' }}</span></div>
          <div class="kv__item"><span class="kv__label">向量维度</span><span class="kv__value num">{{ kbStatus.dimensions }}</span></div>
        </div>
      </div>
      <div class="panel__foot actions">
        <el-button type="primary" :loading="reindexing" :icon="Refresh" @click="handleReindex">重建索引</el-button>
        <el-button link type="primary" @click="$router.push({ name: 'knowledge' })">前往知识库管理 →</el-button>
      </div>
    </div>
    </div>

    <!-- 厂商密钥 -->
    <div class="panel">
      <div class="panel__head">
        <h3 class="panel__title">
          <span class="panel__title-icon"><el-icon><Key /></el-icon></span>
          厂商 API Key
        </h3>
        <el-button type="primary" size="small" :icon="Plus" @click="openProviderDialog()">新增厂商</el-button>
      </div>
      <div class="panel__body panel__body--flush">
        <el-table :data="providers">
          <el-table-column label="厂商" min-width="200">
            <template #default="{ row }">
              <div class="cell">
                <span class="cell__icon provider-icon">{{ row.name?.slice(0, 1) }}</span>
                <div class="cell__main">
                  <div class="cell__title">{{ row.name }}</div>
                  <div class="cell__sub mono">{{ row.code }}</div>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="Base URL" show-overflow-tooltip>
            <template #default="{ row }"><span class="mono muted">{{ row.baseUrl || '—' }}</span></template>
          </el-table-column>
          <el-table-column label="API Key" width="170">
            <template #default="{ row }">
              <el-tag v-if="row.apiKeyMasked" type="success" size="small" class="mono">{{ row.apiKeyMasked }}</el-tag>
              <el-tag v-else type="danger" size="small">未配置</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="90">
            <template #default="{ row }">
              <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
                {{ row.status === 1 ? '启用' : '禁用' }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="90" fixed="right" align="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="openProviderDialog(row)">编辑</el-button>
            </template>
          </el-table-column>
          <template #empty>
            <el-empty description="尚未配置任何厂商" :image-size="80" />
          </template>
        </el-table>
      </div>
      <div class="panel__foot">
        <el-button link type="primary" @click="$router.push({ name: 'ai-models' })">
          前往完整大模型配置（统一模型、虚拟映射、能力标签）→
        </el-button>
      </div>
    </div>

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
          <el-select v-model="providerDialog.form.authType" style="width: 100%">
            <el-option label="Bearer Token" value="Bearer" />
            <el-option label="api-key（Azure OpenAI）" value="api-key" />
            <el-option label="x-api-key" value="x-api-key" />
          </el-select>
        </el-form-item>
        <el-form-item label="API Key">
          <el-input v-model="providerDialog.form.apiKey" type="password" show-password :placeholder="providerDialog.isEdit ? '留空不修改' : '输入 API Key'" />
        </el-form-item>
        <el-form-item v-if="providerDialog.isEdit && providerDialog.form.apiKeyMasked" label="清除密钥">
          <el-checkbox v-model="providerDialog.form.clearApiKey">删除已保存的 API Key</el-checkbox>
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
import { Refresh, CircleCheck, Warning, Check, MagicStick, Connection, Right, Collection, Key, Plus } from '@element-plus/icons-vue'
import { settingsApi } from '../api/settings.js'
import PageHeader from '../components/common/PageHeader.vue'

const PRESETS = {
  openai: {
    providerCode: 'openai',
    providerName: 'OpenAI',
    baseUrl: 'https://api.openai.com/v1',
    canonicalCode: 'gpt-4o-mini',
    canonicalName: 'GPT-4o Mini',
    virtualCode: 'chat-default',
    virtualName: '默认对话模型',
    modelCode: 'gpt-4o-mini',
    contextWindow: 128000,
    reasoningSupported: false
  },
  deepseek: {
    providerCode: 'deepseek',
    providerName: 'DeepSeek',
    baseUrl: 'https://api.deepseek.com/v1',
    canonicalCode: 'deepseek-chat',
    canonicalName: 'DeepSeek Chat',
    virtualCode: 'chat-default',
    virtualName: '默认对话模型',
    modelCode: 'deepseek-chat',
    contextWindow: 64000,
    reasoningSupported: false
  },
  custom: {
    providerCode: '',
    providerName: '',
    baseUrl: '',
    canonicalCode: '',
    canonicalName: '',
    virtualCode: 'chat-default',
    virtualName: '默认对话模型',
    modelCode: '',
    contextWindow: 128000,
    reasoningSupported: false
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
  virtualName: '默认对话模型',
  contextWindow: 128000,
  reasoningSupported: false
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
const checkSteps = computed(() => [
  { title: '厂商', desc: `${overview.providerCount} 个` },
  { title: '密钥', desc: `${providersWithKey.value} 个已配置` },
  { title: '虚拟模型', desc: `${overview.virtualModelCount} 个` },
  { title: '路由可用', desc: resolved.value?.enabled ? '已连通' : '未连通' }
])
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
  return { code: '', name: '', baseUrl: '', authType: 'Bearer', apiKey: '', clearApiKey: false, status: 1 }
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
      modelCode: setupForm.modelCode || setupForm.canonicalCode
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
.grid-2 {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 20px;
  align-items: start;
}

.health--ok .panel__title-icon {
  background: var(--success-soft);
  color: var(--success);
}

.steps {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 12px;
}

.step {
  position: relative;
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border-radius: var(--radius-sm);
  background: var(--surface-2);
  border: 1px solid var(--border);
}

.step__index {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  flex-shrink: 0;
  border-radius: 50%;
  background: var(--surface);
  border: 1.5px solid var(--border-strong);
  color: var(--text-3);
  font-size: 12.5px;
  font-weight: 700;
}

.step__text {
  display: flex;
  flex-direction: column;
  line-height: 1.3;
}

.step__text strong {
  font-size: 14px;
  color: var(--text);
}

.step__text span {
  font-size: 12px;
  color: var(--text-3);
}

.step--done {
  background: var(--success-soft);
  border-color: transparent;
}

.step--done .step__index {
  background: var(--success);
  border-color: var(--success);
  color: #fff;
}

.step--done .step__text span {
  color: var(--success-strong);
}

.step--current {
  border-color: var(--primary);
  box-shadow: 0 0 0 3px var(--primary-soft);
}

.step--current .step__index {
  border-color: var(--primary);
  color: var(--primary);
}

.route__chain {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  padding: 14px 16px;
  border-radius: var(--radius-sm);
  background: var(--surface-2);
}

.route__node {
  display: flex;
  flex-direction: column;
  padding: 8px 12px;
  border-radius: 10px;
  background: var(--surface);
  border: 1px solid var(--border);
  line-height: 1.3;
  min-width: 0;
}

.route__node small {
  font-size: 11px;
  color: var(--text-3);
}

.route__node b {
  font-size: 13px;
  color: var(--text);
}

.route__node--virtual {
  border-color: var(--primary-soft-2);
  background: var(--primary-soft);
}

.route__node--model {
  border-color: #d1f3e6;
  background: var(--success-soft);
}

.route__arrow {
  color: var(--text-3);
}

.kb {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
}

.kb__stat {
  display: flex;
  flex-direction: column;
  gap: 4px;
  padding: 14px 16px;
  border-radius: var(--radius-sm);
  background: linear-gradient(135deg, var(--primary-soft) 0%, #f6f4ff 100%);
}

.kb__stat span {
  font-size: 12px;
  color: var(--text-3);
}

.kb__stat strong {
  font-size: 24px;
  color: var(--primary-deep);
  letter-spacing: -0.02em;
}

.provider-icon {
  background: var(--primary-soft);
  color: var(--primary-deep);
  font-size: 15px;
  font-weight: 800;
}

.actions {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
}

@media (max-width: 1100px) {
  .grid-2 {
    grid-template-columns: 1fr;
  }

  .steps {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}
</style>
