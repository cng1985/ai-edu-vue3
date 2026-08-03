<template>
  <div>
    <div class="page-header">
      <h2>AI 大模型配置</h2>
      <p>分层管理 Provider（厂商）、CanonicalModel（统一模型）、Capability（能力标签）与 VirtualModel（虚拟模型路由）。</p>
    </div>

    <el-row :gutter="16" class="stats-row">
      <el-col :span="4">
        <el-card shadow="never">
          <div class="stat-label">厂商</div>
          <div class="stat-value">{{ overview.providerCount }}</div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="never">
          <div class="stat-label">统一模型</div>
          <div class="stat-value">{{ overview.canonicalModelCount }}</div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="never">
          <div class="stat-label">能力标签</div>
          <div class="stat-value">{{ overview.capabilityCount }}</div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="never">
          <div class="stat-label">虚拟模型</div>
          <div class="stat-value">{{ overview.virtualModelCount }}</div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="never">
          <div class="stat-label">厂商实现</div>
          <div class="stat-value">{{ overview.providerModelCount }}</div>
        </el-card>
      </el-col>
      <el-col :span="4">
        <el-card shadow="never">
          <div class="stat-label">默认虚拟模型</div>
          <div class="stat-value small">{{ overview.defaultVirtualModel || '-' }}</div>
        </el-card>
      </el-col>
    </el-row>

    <el-card shadow="never" class="resolve-card">
      <template #header>
        <span>路由解析测试</span>
      </template>
      <div class="resolve-bar">
        <el-input v-model="resolveCode" placeholder="虚拟模型编码，如 chat-default" clearable />
        <el-button type="primary" :loading="resolving" @click="handleResolve">解析</el-button>
      </div>
      <el-descriptions v-if="resolved" :column="2" border style="margin-top: 12px">
        <el-descriptions-item label="虚拟模型">{{ resolved.virtualModelCode }}</el-descriptions-item>
        <el-descriptions-item label="统一模型">{{ resolved.canonicalModelCode }}</el-descriptions-item>
        <el-descriptions-item label="厂商">{{ resolved.providerCode }}</el-descriptions-item>
        <el-descriptions-item label="调用模型">{{ resolved.modelCode }}</el-descriptions-item>
        <el-descriptions-item label="Base URL">{{ resolved.baseUrl }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="resolved.enabled ? 'success' : 'info'" size="small">
            {{ resolved.enabled ? '可用' : '未配置密钥' }}
          </el-tag>
        </el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-card shadow="never">
      <el-tabs v-model="activeTab">
        <el-tab-pane label="厂商" name="providers">
          <div class="tab-toolbar">
            <el-button v-if="canManage" type="primary" @click="openProviderDialog()">新增厂商</el-button>
          </div>
          <el-table :data="providers" stripe v-loading="loadingProviders">
            <el-table-column prop="code" label="编码" width="120" />
            <el-table-column prop="name" label="名称" width="140" />
            <el-table-column prop="baseUrl" label="Base URL" show-overflow-tooltip />
            <el-table-column prop="authType" label="认证" width="80" />
            <el-table-column prop="apiKeyMasked" label="API Key" width="140" />
            <el-table-column prop="modelCount" label="模型数" width="80" />
            <el-table-column label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
                  {{ row.status === 1 ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column v-if="canManage" label="操作" width="160" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="openProviderDialog(row)">编辑</el-button>
                <el-button link type="danger" @click="deleteProvider(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="统一模型" name="canonical">
          <div class="tab-toolbar">
            <el-button v-if="canManage" type="primary" @click="openCanonicalDialog()">新增统一模型</el-button>
          </div>
          <el-table :data="canonicalModels" stripe v-loading="loadingCanonical">
            <el-table-column prop="code" label="编码" width="140" />
            <el-table-column prop="name" label="名称" width="160" />
            <el-table-column prop="contextWindow" label="上下文窗口" width="120" />
            <el-table-column label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
                  {{ row.status === 1 ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column v-if="canManage" label="操作" width="160" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="openCanonicalDialog(row)">编辑</el-button>
                <el-button link type="danger" @click="deleteCanonical(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="厂商模型" name="providerModels">
          <div class="tab-toolbar">
            <el-button v-if="canManage" type="primary" @click="openProviderModelDialog()">新增厂商模型</el-button>
          </div>
          <el-table :data="providerModels" stripe v-loading="loadingProviderModels">
            <el-table-column prop="providerCode" label="厂商" width="100" />
            <el-table-column prop="canonicalModelCode" label="统一模型" width="140" />
            <el-table-column prop="modelCode" label="模型标识" width="140" />
            <el-table-column prop="deploymentName" label="部署名" width="120" />
            <el-table-column prop="priority" label="优先级" width="80" />
            <el-table-column prop="weight" label="权重" width="80" />
            <el-table-column label="推理" width="80">
              <template #default="{ row }">
                <el-tag :type="row.reasoningSupported ? 'warning' : 'info'" size="small">
                  {{ row.reasoningSupported ? '支持' : '否' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
                  {{ row.status === 1 ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column v-if="canManage" label="操作" width="160" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="openProviderModelDialog(row)">编辑</el-button>
                <el-button link type="danger" @click="deleteProviderModel(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="虚拟模型" name="virtual">
          <div class="tab-toolbar">
            <el-button v-if="canManage" type="primary" @click="openVirtualDialog()">新增虚拟模型</el-button>
          </div>
          <el-table :data="virtualModels" stripe v-loading="loadingVirtual">
            <el-table-column prop="code" label="编码" width="140" />
            <el-table-column prop="name" label="名称" width="160" />
            <el-table-column prop="mappingCount" label="映射数" width="80" />
            <el-table-column label="默认" width="80">
              <template #default="{ row }">
                <el-tag v-if="row.code === overview.defaultVirtualModel" type="warning" size="small">默认</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
                  {{ row.status === 1 ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column v-if="canManage" label="操作" width="220" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="setDefault(row)">设为默认</el-button>
                <el-button link type="primary" @click="openVirtualDialog(row)">编辑</el-button>
                <el-button link type="danger" @click="deleteVirtual(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="虚拟映射" name="mappings">
          <div class="tab-toolbar">
            <el-button v-if="canManage" type="primary" @click="openMappingDialog()">新增映射</el-button>
          </div>
          <el-table :data="mappings" stripe v-loading="loadingMappings">
            <el-table-column prop="virtualModelCode" label="虚拟模型" width="140" />
            <el-table-column prop="canonicalModelCode" label="统一模型" width="140" />
            <el-table-column prop="priority" label="优先级" width="80" />
            <el-table-column v-if="canManage" label="操作" width="160" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="openMappingDialog(row)">编辑</el-button>
                <el-button link type="danger" @click="deleteMapping(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="能力关联" name="capabilityModels">
          <div class="tab-toolbar">
            <el-button v-if="canManage" type="primary" @click="openCapModelDialog()">新增关联</el-button>
          </div>
          <el-table :data="capabilityModels" stripe v-loading="loadingCapModels">
            <el-table-column prop="canonicalModelCode" label="统一模型" width="160" />
            <el-table-column prop="canonicalModelName" label="模型名称" width="160" />
            <el-table-column prop="capabilityCode" label="能力编码" width="140" />
            <el-table-column prop="capabilityName" label="能力名称" width="160" />
            <el-table-column v-if="canManage" label="操作" width="100" fixed="right">
              <template #default="{ row }">
                <el-button link type="danger" @click="deleteCapModel(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="能力标签" name="capabilities">
          <div class="tab-toolbar">
            <el-button v-if="canManage" type="primary" @click="openCapabilityDialog()">新增能力</el-button>
          </div>
          <el-table :data="capabilities" stripe v-loading="loadingCapabilities">
            <el-table-column prop="code" label="编码" width="140" />
            <el-table-column prop="name" label="名称" width="160" />
            <el-table-column label="状态" width="80">
              <template #default="{ row }">
                <el-tag :type="row.status === 1 ? 'success' : 'info'" size="small">
                  {{ row.status === 1 ? '启用' : '禁用' }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column v-if="canManage" label="操作" width="160" fixed="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="openCapabilityDialog(row)">编辑</el-button>
                <el-button link type="danger" @click="deleteCapability(row)">删除</el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </el-card>

    <!-- 厂商对话框 -->
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
          <el-input v-model="providerDialog.form.apiKey" type="password" show-password placeholder="留空不修改" />
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

    <!-- 统一模型对话框 -->
    <el-dialog v-model="canonicalDialog.visible" :title="canonicalDialog.isEdit ? '编辑统一模型' : '新增统一模型'" width="480px">
      <el-form label-width="110px">
        <el-form-item label="编码" required>
          <el-input v-model="canonicalDialog.form.code" :disabled="canonicalDialog.isEdit" />
        </el-form-item>
        <el-form-item label="名称" required>
          <el-input v-model="canonicalDialog.form.name" />
        </el-form-item>
        <el-form-item label="上下文窗口">
          <el-input-number v-model="canonicalDialog.form.contextWindow" :min="0" :step="1000" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="canonicalDialog.form.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="canonicalDialog.visible = false">取消</el-button>
        <el-button type="primary" :loading="canonicalDialog.saving" @click="saveCanonical">保存</el-button>
      </template>
    </el-dialog>

    <!-- 厂商模型对话框 -->
    <el-dialog v-model="pmDialog.visible" :title="pmDialog.isEdit ? '编辑厂商模型' : '新增厂商模型'" width="520px">
      <el-form label-width="100px">
        <el-form-item label="厂商" required>
          <el-select v-model="pmDialog.form.providerId" :disabled="pmDialog.isEdit" filterable>
            <el-option v-for="p in providers" :key="p.id" :label="`${p.code} - ${p.name}`" :value="p.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="统一模型" required>
          <el-select v-model="pmDialog.form.canonicalModelId" :disabled="pmDialog.isEdit" filterable>
            <el-option v-for="m in canonicalModels" :key="m.id" :label="`${m.code} - ${m.name}`" :value="m.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="模型标识" required>
          <el-input v-model="pmDialog.form.modelCode" />
        </el-form-item>
        <el-form-item label="部署名">
          <el-input v-model="pmDialog.form.deploymentName" placeholder="Azure 部署名（可选）" />
        </el-form-item>
        <el-form-item label="优先级">
          <el-input-number v-model="pmDialog.form.priority" :min="1" />
        </el-form-item>
        <el-form-item label="权重">
          <el-input-number v-model="pmDialog.form.weight" :min="1" />
        </el-form-item>
        <el-form-item label="推理模型">
          <el-switch v-model="pmDialog.form.reasoningSupported" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="pmDialog.form.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="pmDialog.visible = false">取消</el-button>
        <el-button type="primary" :loading="pmDialog.saving" @click="saveProviderModel">保存</el-button>
      </template>
    </el-dialog>

    <!-- 虚拟模型对话框 -->
    <el-dialog v-model="virtualDialog.visible" :title="virtualDialog.isEdit ? '编辑虚拟模型' : '新增虚拟模型'" width="480px">
      <el-form label-width="80px">
        <el-form-item label="编码" required>
          <el-input v-model="virtualDialog.form.code" :disabled="virtualDialog.isEdit" />
        </el-form-item>
        <el-form-item label="名称" required>
          <el-input v-model="virtualDialog.form.name" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="virtualDialog.form.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="virtualDialog.visible = false">取消</el-button>
        <el-button type="primary" :loading="virtualDialog.saving" @click="saveVirtual">保存</el-button>
      </template>
    </el-dialog>

    <!-- 映射对话框 -->
    <el-dialog v-model="mappingDialog.visible" :title="mappingDialog.isEdit ? '编辑映射' : '新增映射'" width="480px">
      <el-form label-width="100px">
        <el-form-item label="虚拟模型" required>
          <el-select v-model="mappingDialog.form.virtualModelId" :disabled="mappingDialog.isEdit" filterable>
            <el-option v-for="v in virtualModels" :key="v.id" :label="`${v.code} - ${v.name}`" :value="v.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="统一模型" required>
          <el-select v-model="mappingDialog.form.canonicalModelId" :disabled="mappingDialog.isEdit" filterable>
            <el-option v-for="m in canonicalModels" :key="m.id" :label="`${m.code} - ${m.name}`" :value="m.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="优先级">
          <el-input-number v-model="mappingDialog.form.priority" :min="1" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="mappingDialog.visible = false">取消</el-button>
        <el-button type="primary" :loading="mappingDialog.saving" @click="saveMapping">保存</el-button>
      </template>
    </el-dialog>

    <!-- 能力关联对话框 -->
    <el-dialog v-model="capModelDialog.visible" title="新增能力关联" width="480px">
      <el-form label-width="100px">
        <el-form-item label="统一模型" required>
          <el-select v-model="capModelDialog.form.canonicalModelId" filterable>
            <el-option v-for="m in canonicalModels" :key="m.id" :label="`${m.code} - ${m.name}`" :value="m.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="能力标签" required>
          <el-select v-model="capModelDialog.form.capabilityId" filterable>
            <el-option v-for="c in capabilities" :key="c.id" :label="`${c.code} - ${c.name}`" :value="c.id" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="capModelDialog.visible = false">取消</el-button>
        <el-button type="primary" :loading="capModelDialog.saving" @click="saveCapModel">保存</el-button>
      </template>
    </el-dialog>

    <!-- 能力对话框 -->
    <el-dialog v-model="capDialog.visible" :title="capDialog.isEdit ? '编辑能力' : '新增能力'" width="480px">
      <el-form label-width="80px">
        <el-form-item label="编码" required>
          <el-input v-model="capDialog.form.code" :disabled="capDialog.isEdit" />
        </el-form-item>
        <el-form-item label="名称" required>
          <el-input v-model="capDialog.form.name" />
        </el-form-item>
        <el-form-item label="状态">
          <el-switch v-model="capDialog.form.status" :active-value="1" :inactive-value="0" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="capDialog.visible = false">取消</el-button>
        <el-button type="primary" :loading="capDialog.saving" @click="saveCapability">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { aiModelsApi } from '../api/aiModels.js'
import { useAuthStore } from '../stores/auth.js'
import { PERM } from '../constants/permissions.js'

const auth = useAuthStore()
const canManage = computed(() => auth.hasPermission(PERM.AI_MODEL_MANAGE))

const activeTab = ref('providers')
const overview = reactive({
  providerCount: 0,
  canonicalModelCount: 0,
  capabilityCount: 0,
  virtualModelCount: 0,
  providerModelCount: 0,
  defaultVirtualModel: ''
})

const providers = ref([])
const canonicalModels = ref([])
const providerModels = ref([])
const virtualModels = ref([])
const mappings = ref([])
const capabilities = ref([])
const capabilityModels = ref([])

const loadingProviders = ref(false)
const loadingCanonical = ref(false)
const loadingProviderModels = ref(false)
const loadingVirtual = ref(false)
const loadingMappings = ref(false)
const loadingCapabilities = ref(false)
const loadingCapModels = ref(false)

const resolveCode = ref('chat-default')
const resolved = ref(null)
const resolving = ref(false)

const providerDialog = reactive({ visible: false, isEdit: false, saving: false, id: '', form: emptyProvider() })
const canonicalDialog = reactive({ visible: false, isEdit: false, saving: false, id: '', form: emptyCanonical() })
const pmDialog = reactive({ visible: false, isEdit: false, saving: false, id: '', form: emptyProviderModel() })
const virtualDialog = reactive({ visible: false, isEdit: false, saving: false, id: '', form: emptyVirtual() })
const mappingDialog = reactive({ visible: false, isEdit: false, saving: false, id: '', form: emptyMapping() })
const capDialog = reactive({ visible: false, isEdit: false, saving: false, id: '', form: emptyCapability() })
const capModelDialog = reactive({ visible: false, saving: false, form: emptyCapModel() })

function emptyProvider() {
  return { code: '', name: '', baseUrl: '', authType: 'Bearer', apiKey: '', status: 1 }
}
function emptyCanonical() {
  return { code: '', name: '', contextWindow: 128000, status: 1 }
}
function emptyProviderModel() {
  return { providerId: '', canonicalModelId: '', modelCode: '', deploymentName: '', priority: 10, weight: 100, reasoningSupported: false, status: 1 }
}
function emptyCapModel() {
  return { canonicalModelId: '', capabilityId: '' }
}
function emptyVirtual() {
  return { code: '', name: '', status: 1 }
}
function emptyMapping() {
  return { virtualModelId: '', canonicalModelId: '', priority: 10 }
}
function emptyCapability() {
  return { code: '', name: '', status: 1 }
}

async function loadOverview() {
  const data = await aiModelsApi.overview()
  Object.assign(overview, data)
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

async function loadCanonical() {
  loadingCanonical.value = true
  try {
    const res = await aiModelsApi.listCanonicalModels({ page: 1, pageSize: 100 })
    canonicalModels.value = res.list || []
  } finally {
    loadingCanonical.value = false
  }
}

async function loadProviderModels() {
  loadingProviderModels.value = true
  try {
    providerModels.value = await aiModelsApi.listProviderModels()
  } finally {
    loadingProviderModels.value = false
  }
}

async function loadVirtual() {
  loadingVirtual.value = true
  try {
    const res = await aiModelsApi.listVirtualModels({ page: 1, pageSize: 100 })
    virtualModels.value = res.list || []
  } finally {
    loadingVirtual.value = false
  }
}

async function loadMappings() {
  loadingMappings.value = true
  try {
    mappings.value = await aiModelsApi.listVirtualModelMappings()
  } finally {
    loadingMappings.value = false
  }
}

async function loadCapabilities() {
  loadingCapabilities.value = true
  try {
    const res = await aiModelsApi.listCapabilities({ page: 1, pageSize: 100 })
    capabilities.value = res.list || []
  } finally {
    loadingCapabilities.value = false
  }
}

async function loadCapModels() {
  loadingCapModels.value = true
  try {
    capabilityModels.value = await aiModelsApi.listCapabilityModels()
  } finally {
    loadingCapModels.value = false
  }
}

async function loadAll() {
  await loadOverview()
  await Promise.all([
    loadProviders(),
    loadCanonical(),
    loadProviderModels(),
    loadVirtual(),
    loadMappings(),
    loadCapabilities(),
    loadCapModels()
  ])
}

async function handleResolve() {
  resolving.value = true
  try {
    resolved.value = await aiModelsApi.resolve(resolveCode.value)
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
  } finally {
    providerDialog.saving = false
  }
}

async function deleteProvider(row) {
  await ElMessageBox.confirm(`确定删除厂商「${row.name}」？`, '确认')
  await aiModelsApi.deleteProvider(row.id)
  ElMessage.success('已删除')
  await loadProviders()
  await loadOverview()
}

function openCanonicalDialog(row) {
  canonicalDialog.isEdit = !!row
  canonicalDialog.id = row?.id || ''
  canonicalDialog.form = row ? { ...row } : emptyCanonical()
  canonicalDialog.visible = true
}

async function saveCanonical() {
  canonicalDialog.saving = true
  try {
    if (canonicalDialog.isEdit) {
      await aiModelsApi.updateCanonicalModel(canonicalDialog.id, canonicalDialog.form)
    } else {
      await aiModelsApi.createCanonicalModel(canonicalDialog.form)
    }
    canonicalDialog.visible = false
    ElMessage.success('保存成功')
    await loadCanonical()
    await loadOverview()
  } finally {
    canonicalDialog.saving = false
  }
}

async function deleteCanonical(row) {
  await ElMessageBox.confirm(`确定删除统一模型「${row.name}」？`, '确认')
  await aiModelsApi.deleteCanonicalModel(row.id)
  ElMessage.success('已删除')
  await loadCanonical()
  await loadOverview()
}

function openProviderModelDialog(row) {
  pmDialog.isEdit = !!row
  pmDialog.id = row?.id || ''
  pmDialog.form = row ? { ...row } : emptyProviderModel()
  pmDialog.visible = true
}

async function saveProviderModel() {
  pmDialog.saving = true
  try {
    if (pmDialog.isEdit) {
      await aiModelsApi.updateProviderModel(pmDialog.id, pmDialog.form)
    } else {
      await aiModelsApi.createProviderModel(pmDialog.form)
    }
    pmDialog.visible = false
    ElMessage.success('保存成功')
    await loadProviderModels()
    await loadOverview()
  } finally {
    pmDialog.saving = false
  }
}

async function deleteProviderModel(row) {
  await ElMessageBox.confirm('确定删除该厂商模型？', '确认')
  await aiModelsApi.deleteProviderModel(row.id)
  ElMessage.success('已删除')
  await loadProviderModels()
  await loadOverview()
}

function openVirtualDialog(row) {
  virtualDialog.isEdit = !!row
  virtualDialog.id = row?.id || ''
  virtualDialog.form = row ? { ...row } : emptyVirtual()
  virtualDialog.visible = true
}

async function saveVirtual() {
  virtualDialog.saving = true
  try {
    if (virtualDialog.isEdit) {
      await aiModelsApi.updateVirtualModel(virtualDialog.id, virtualDialog.form)
    } else {
      await aiModelsApi.createVirtualModel(virtualDialog.form)
    }
    virtualDialog.visible = false
    ElMessage.success('保存成功')
    await loadVirtual()
    await loadOverview()
  } finally {
    virtualDialog.saving = false
  }
}

async function deleteVirtual(row) {
  await ElMessageBox.confirm(`确定删除虚拟模型「${row.name}」？`, '确认')
  await aiModelsApi.deleteVirtualModel(row.id)
  ElMessage.success('已删除')
  await loadVirtual()
  await loadOverview()
}

async function setDefault(row) {
  await aiModelsApi.setDefault(row.code)
  ElMessage.success('默认虚拟模型已更新')
  await loadOverview()
}

function openMappingDialog(row) {
  mappingDialog.isEdit = !!row
  mappingDialog.id = row?.id || ''
  mappingDialog.form = row ? { ...row } : emptyMapping()
  mappingDialog.visible = true
}

async function saveMapping() {
  mappingDialog.saving = true
  try {
    if (mappingDialog.isEdit) {
      await aiModelsApi.updateVirtualModelMapping(mappingDialog.id, mappingDialog.form)
    } else {
      await aiModelsApi.createVirtualModelMapping(mappingDialog.form)
    }
    mappingDialog.visible = false
    ElMessage.success('保存成功')
    await loadMappings()
  } finally {
    mappingDialog.saving = false
  }
}

async function deleteMapping(row) {
  await ElMessageBox.confirm('确定删除该映射？', '确认')
  await aiModelsApi.deleteVirtualModelMapping(row.id)
  ElMessage.success('已删除')
  await loadMappings()
}

function openCapabilityDialog(row) {
  capDialog.isEdit = !!row
  capDialog.id = row?.id || ''
  capDialog.form = row ? { ...row } : emptyCapability()
  capDialog.visible = true
}

async function saveCapability() {
  capDialog.saving = true
  try {
    if (capDialog.isEdit) {
      await aiModelsApi.updateCapability(capDialog.id, capDialog.form)
    } else {
      await aiModelsApi.createCapability(capDialog.form)
    }
    capDialog.visible = false
    ElMessage.success('保存成功')
    await loadCapabilities()
    await loadOverview()
  } finally {
    capDialog.saving = false
  }
}

async function deleteCapability(row) {
  await ElMessageBox.confirm(`确定删除能力「${row.name}」？`, '确认')
  await aiModelsApi.deleteCapability(row.id)
  ElMessage.success('已删除')
  await loadCapabilities()
  await loadOverview()
}

function openCapModelDialog() {
  capModelDialog.form = emptyCapModel()
  capModelDialog.visible = true
}

async function saveCapModel() {
  capModelDialog.saving = true
  try {
    await aiModelsApi.createCapabilityModel(capModelDialog.form)
    capModelDialog.visible = false
    ElMessage.success('保存成功')
    await loadCapModels()
  } finally {
    capModelDialog.saving = false
  }
}

async function deleteCapModel(row) {
  await ElMessageBox.confirm('确定删除该能力关联？', '确认')
  await aiModelsApi.deleteCapabilityModel(row.id)
  ElMessage.success('已删除')
  await loadCapModels()
}

onMounted(async () => {
  await loadAll()
  await handleResolve()
})
</script>

<style scoped>
.page-header p {
  margin: 6px 0 0;
  color: #6b7280;
  font-size: 14px;
}
.stats-row {
  margin-bottom: 16px;
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
  font-size: 14px;
  font-weight: 500;
}
.resolve-card {
  margin-bottom: 16px;
}
.resolve-bar {
  display: flex;
  gap: 12px;
}
.tab-toolbar {
  margin-bottom: 12px;
}
</style>
