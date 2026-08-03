<template>
  <div>
    <div class="page-header">
      <h2>系统设置</h2>
      <p>管理大模型默认路由、厂商密钥与知识库索引等运行时可调项。</p>
    </div>

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
        v-if="!allReady"
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
          <el-button
            v-if="canManageModel"
            type="primary"
            :loading="savingDefault"
            :disabled="!defaultModel"
            @click="saveDefaultModel"
          >
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
          <el-button v-if="canManageModel" type="primary" size="small" @click="openProviderDialog()">
            新增厂商
          </el-button>
        </div>
      </template>

      <el-table :data="providers" stripe v-loading="loadingProviders">
        <el-table-column prop="code" label="编码" width="120" />
        <el-table-column prop="name" label="名称" width="140" />
        <el-table-column prop="baseUrl" label="Base URL" show-overflow-tooltip />
        <el-table-column prop="apiKeyMasked" label="API Key" width="140">
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
        <el-table-column v-if="canManageModel" label="操作" width="100" fixed="right">
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
        <el-button
          v-if="canManageKb"
          type="primary"
          :loading="reindexing"
          @click="handleReindex"
        >
          重建索引
        </el-button>
        <el-button @click="loadKbStatus">刷新</el-button>
      </div>
    </el-card>

    <!-- 厂商编辑对话框 -->
    <el-dialog
      v-model="providerDialog.visible"
      :title="providerDialog.isEdit ? '编辑厂商' : '新增厂商'"
      width="520px"
    >
      <el-form label-width="100px">
        <el-form-item label="编码" required>
          <el-input v-model="providerDialog.form.code" :disabled="providerDialog.isEdit" placeholder="如 openai" />
        </el-form-item>
        <el-form-item label="名称" required>
          <el-input v-model="providerDialog.form.name" placeholder="如 OpenAI" />
        </el-form-item>
        <el-form-item label="Base URL">
          <el-input v-model="providerDialog.form.baseUrl" placeholder="https://api.openai.com/v1" />
        </el-form-item>
        <el-form-item label="认证类型">
          <el-input v-model="providerDialog.form.authType" placeholder="Bearer" />
        </el-form-item>
        <el-form-item label="API Key">
          <el-input
            v-model="providerDialog.form.apiKey"
            type="password"
            show-password
            :placeholder="providerDialog.isEdit ? '留空不修改' : '输入 API Key'"
          />
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
import { aiModelsApi } from '../api/aiModels.js'
import { knowledgeApi } from '../api/knowledge.js'
import { useAuthStore } from '../stores/auth.js'
import { PERM } from '../constants/permissions.js'

const auth = useAuthStore()
const canManageModel = computed(() => auth.hasPermission(PERM.AI_MODEL_MANAGE))
const canManageKb = computed(() => auth.hasPermission(PERM.KNOWLEDGE_MANAGE))

const loading = ref(false)
const loadingProviders = ref(false)
const savingDefault = ref(false)
const resolving = ref(false)
const reindexing = ref(false)

const overview = reactive({
  providerCount: 0,
  virtualModelCount: 0,
  defaultVirtualModel: ''
})
const providers = ref([])
const virtualModels = ref([])
const defaultModel = ref('')
const resolved = ref(null)
const kbStatus = ref(null)

const providersWithKey = computed(() => providers.value.filter((p) => p.apiKeyMasked).length)

const allReady = computed(() =>
  overview.providerCount > 0 &&
  providersWithKey.value > 0 &&
  overview.virtualModelCount > 0 &&
  resolved.value?.enabled
)

const readyStep = computed(() => {
  if (resolved.value?.enabled) return 4
  if (overview.virtualModelCount > 0) return 3
  if (providersWithKey.value > 0) return 2
  if (overview.providerCount > 0) return 1
  return 0
})

const checklistHint = computed(() => {
  if (overview.providerCount === 0) return '请先添加至少一个厂商（如 OpenAI、DeepSeek）'
  if (providersWithKey.value === 0) return '请为厂商配置 API Key'
  if (overview.virtualModelCount === 0) return '请在大模型配置中创建虚拟模型'
  if (!resolved.value?.enabled) return '路由未连通，请检查厂商模型映射与 API Key 是否正确'
  return ''
})

const providerDialog = reactive({
  visible: false,
  isEdit: false,
  saving: false,
  id: '',
  form: emptyProvider()
})

function emptyProvider() {
  return { code: '', name: '', baseUrl: '', authType: 'Bearer', apiKey: '', status: 1 }
}

async function loadOverview() {
  const data = await aiModelsApi.overview()
  Object.assign(overview, data)
  defaultModel.value = data.defaultVirtualModel || ''
}

async function loadProviders() {
  loadingProviders.value = true
  try {
    const res = await aiModelsApi.listProviders({ page: 1, pageSize: 100 })
    providers.value = res.list || []
  } finally {
    loadingProviders.value = false
  }
}

async function loadVirtualModels() {
  virtualModels.value = await aiModelsApi.listVirtualModelOptions()
}

async function loadResolve() {
  const code = defaultModel.value || overview.defaultVirtualModel || 'chat-default'
  if (!code) return
  resolved.value = await aiModelsApi.resolve(code)
}

async function loadKbStatus() {
  kbStatus.value = await knowledgeApi.status()
}

async function loadData() {
  loading.value = true
  try {
    await Promise.all([loadOverview(), loadProviders(), loadVirtualModels(), loadKbStatus()])
    await loadResolve()
  } finally {
    loading.value = false
  }
}

async function saveDefaultModel() {
  savingDefault.value = true
  try {
    await aiModelsApi.setDefault(defaultModel.value)
    overview.defaultVirtualModel = defaultModel.value
    ElMessage.success('默认虚拟模型已更新')
    await loadResolve()
  } finally {
    savingDefault.value = false
  }
}

async function testResolve() {
  resolving.value = true
  try {
    await loadResolve()
    ElMessage.success(resolved.value?.enabled ? '路由解析成功，大模型可用' : '路由已解析，但厂商 API Key 未配置')
  } finally {
    resolving.value = false
  }
}

function openProviderDialog(row) {
  providerDialog.isEdit = !!row
  providerDialog.id = row?.id || ''
  providerDialog.form = row ? { ...row, apiKey: '' } : emptyProvider()
  providerDialog.visible = true
}

async function saveProvider() {
  providerDialog.saving = true
  try {
    if (providerDialog.isEdit) {
      await aiModelsApi.updateProvider(providerDialog.id, providerDialog.form)
    } else {
      await aiModelsApi.createProvider(providerDialog.form)
    }
    providerDialog.visible = false
    ElMessage.success('保存成功')
    await loadProviders()
    await loadOverview()
    await loadResolve()
  } finally {
    providerDialog.saving = false
  }
}

async function handleReindex() {
  reindexing.value = true
  try {
    kbStatus.value = await knowledgeApi.reindex()
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
.card-footer {
  margin-top: 12px;
}
.actions {
  display: flex;
  gap: 12px;
}
</style>
