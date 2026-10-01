<template>
  <div class="page">
    <PageHeader eyebrow="成长生态" title="职业能力体系" subtitle="维护 职业 → 岗位 → 能力 → 技能要求 的成长模型，学习者的差距分析与任务匹配都基于此计算。">
      <el-button type="primary" :icon="Plus" @click="openCareer()">新增职业</el-button>
    </PageHeader>

    <div class="layout">
      <aside class="panel tree">
        <div class="panel__head"><h3 class="panel__title">职业与岗位</h3></div>
        <div class="panel__body tree__body" v-loading="loading">
          <div v-for="c in careers" :key="c.id" class="tree__career">
            <div class="tree__career-head">
              <span class="tree__dot" :style="{ background: c.color }" />
              <strong>{{ c.name }}</strong>
              <span class="muted">{{ c.category }}</span>
              <span class="tree__ops">
                <el-button link :icon="Plus" title="新增岗位" @click="openRole(null, c.id)" />
                <el-button link :icon="Edit" title="编辑职业" @click="openCareer(c)" />
                <el-popconfirm title="删除职业将同时删除其岗位与能力，确定？" @confirm="removeCareer(c.id)">
                  <template #reference><el-button link type="danger" :icon="Delete" /></template>
                </el-popconfirm>
              </span>
            </div>
            <button
              v-for="r in c.roles"
              :key="r.id"
              class="tree__role"
              :class="{ active: r.id === roleId }"
              @click="selectRole(r.id)"
            >
              {{ r.name }} <el-tag size="small" effect="plain">{{ r.level }}</el-tag>
            </button>
            <div v-if="!c.roles?.length" class="muted tree__empty">暂无岗位</div>
          </div>
        </div>
      </aside>

      <section v-if="role" class="panel">
        <div class="panel__head">
          <div>
            <h3 class="panel__title">{{ role.name }} <el-tag size="small">{{ role.level }}</el-tag></h3>
            <p class="panel__sub">{{ role.description }}</p>
          </div>
          <div class="toolbar">
            <el-button :icon="Edit" @click="openRole(role)">编辑岗位</el-button>
            <el-button type="primary" :icon="Plus" @click="openCap()">新增能力</el-button>
            <el-popconfirm title="确定删除该岗位？" @confirm="removeRole(role.id)">
              <template #reference><el-button type="danger" plain :icon="Delete">删除</el-button></template>
            </el-popconfirm>
          </div>
        </div>
        <div class="panel__body caps">
          <el-empty v-if="!role.capabilities?.length" description="还没有能力，点击右上角新增" :image-size="80" />
          <div v-for="cap in role.capabilities" :key="cap.id" class="cap">
            <div class="cap__head">
              <div>
                <strong>{{ cap.name }}</strong>
                <el-tag size="small" type="info" effect="plain">权重 {{ cap.weight }}</el-tag>
                <p class="muted">{{ cap.description }}</p>
              </div>
              <div>
                <el-button link type="primary" @click="openCap(cap)">编辑</el-button>
                <el-popconfirm title="确定删除该能力？" @confirm="removeCap(cap.id)">
                  <template #reference><el-button link type="danger">删除</el-button></template>
                </el-popconfirm>
              </div>
            </div>
            <el-table :data="cap.skills" size="small">
              <el-table-column label="技能" min-width="180">
                <template #default="{ row }">{{ row.skill?.name || row.skillId }} <span class="muted mono">{{ row.skillId }}</span></template>
              </el-table-column>
              <el-table-column label="分类" width="120"><template #default="{ row }">{{ row.skill?.category }}</template></el-table-column>
              <el-table-column label="要求等级" width="220">
                <template #default="{ row }">
                  <span class="lvl" :style="{ background: LEVEL_COLORS[row.requiredLevel] }">L{{ row.requiredLevel }}</span>
                  {{ LEVELS[row.requiredLevel] }}
                </template>
              </el-table-column>
              <el-table-column label="权重" width="80" prop="weight" />
            </el-table>
          </div>
        </div>
      </section>
      <section v-else class="panel"><el-empty description="请选择左侧岗位" /></section>
    </div>

    <el-dialog v-model="careerDialog.visible" :title="careerDialog.id ? '编辑职业' : '新增职业'" width="520px">
      <el-form :model="careerDialog.form" label-width="88px">
        <el-form-item v-if="!careerDialog.id" label="标识"><el-input v-model="careerDialog.form.id" placeholder="如 ai-engineer，留空自动生成" /></el-form-item>
        <el-form-item label="名称" required><el-input v-model="careerDialog.form.name" /></el-form-item>
        <el-form-item label="分类"><el-input v-model="careerDialog.form.category" placeholder="研发 / 数据 / 产品" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="careerDialog.form.description" type="textarea" :rows="3" /></el-form-item>
        <el-form-item label="市场需求"><el-input v-model="careerDialog.form.demand" placeholder="高 / 中" /></el-form-item>
        <el-form-item label="薪资范围"><el-input v-model="careerDialog.form.salaryRange" placeholder="12–40K" /></el-form-item>
        <el-form-item label="图标"><el-input v-model="careerDialog.form.icon" placeholder="学习端图标名，如 code / sparkles" /></el-form-item>
        <el-form-item label="主题色"><el-color-picker v-model="careerDialog.form.color" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="careerDialog.form.sort" :min="0" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="careerDialog.visible = false">取消</el-button>
        <el-button type="primary" @click="saveCareer">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="roleDialog.visible" :title="roleDialog.id ? '编辑岗位' : '新增岗位'" width="520px">
      <el-form :model="roleDialog.form" label-width="88px">
        <el-form-item v-if="!roleDialog.id" label="标识"><el-input v-model="roleDialog.form.id" placeholder="如 senior-java，留空自动生成" /></el-form-item>
        <el-form-item label="所属职业" required>
          <el-select v-model="roleDialog.form.careerId" style="width: 100%"><el-option v-for="c in careers" :key="c.id" :label="c.name" :value="c.id" /></el-select>
        </el-form-item>
        <el-form-item label="名称" required><el-input v-model="roleDialog.form.name" /></el-form-item>
        <el-form-item label="级别">
          <el-select v-model="roleDialog.form.level" style="width: 100%"><el-option v-for="l in ROLE_LEVELS" :key="l" :label="l" :value="l" /></el-select>
        </el-form-item>
        <el-form-item label="描述"><el-input v-model="roleDialog.form.description" type="textarea" :rows="3" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="roleDialog.form.sort" :min="0" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="roleDialog.visible = false">取消</el-button>
        <el-button type="primary" @click="saveRole">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="capDialog.visible" :title="capDialog.id ? '编辑能力' : '新增能力'" width="640px">
      <el-form :model="capDialog.form" label-width="88px">
        <el-form-item label="能力名称" required><el-input v-model="capDialog.form.name" placeholder="如：分布式系统能力" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="capDialog.form.description" placeholder="能完成什么类型的问题" /></el-form-item>
        <el-form-item label="权重"><el-input-number v-model="capDialog.form.weight" :min="1" :max="10" /></el-form-item>
        <el-form-item label="排序"><el-input-number v-model="capDialog.form.sort" :min="0" /></el-form-item>
        <el-form-item label="技能要求">
          <div class="req-list">
            <div v-for="(s, i) in capDialog.form.skills" :key="i" class="req">
              <el-select v-model="s.skillId" filterable placeholder="选择技能" style="flex: 1">
                <el-option v-for="sk in skills" :key="sk.id" :label="`${sk.category} · ${sk.name}`" :value="sk.id" />
              </el-select>
              <el-select v-model="s.requiredLevel" style="width: 150px">
                <el-option v-for="(name, l) in LEVELS" v-show="l > 0" :key="l" :label="`L${l} ${name}`" :value="l" />
              </el-select>
              <el-input-number v-model="s.weight" :min="1" :max="5" style="width: 100px" />
              <el-button link type="danger" :icon="Delete" @click="capDialog.form.skills.splice(i, 1)" />
            </div>
            <el-button :icon="Plus" size="small" @click="capDialog.form.skills.push({ skillId: '', requiredLevel: 2, weight: 1 })">添加技能要求</el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="capDialog.visible = false">取消</el-button>
        <el-button type="primary" @click="saveCap">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Delete, Edit, Plus } from '@element-plus/icons-vue'
import PageHeader from '../../components/common/PageHeader.vue'
import { ecoAdminApi } from '../../api'
import { LEVELS, LEVEL_COLORS, ROLE_LEVELS } from '../../utils/eco'

const careers = ref([])
const skills = ref([])
const roleId = ref('')
const role = ref(null)
const loading = ref(false)

const careerDialog = reactive({ visible: false, id: '', form: {} })
const roleDialog = reactive({ visible: false, id: '', form: {} })
const capDialog = reactive({ visible: false, id: '', form: {} })

async function load() {
  loading.value = true
  try {
    careers.value = await ecoAdminApi.careers()
    if (!roleId.value) roleId.value = careers.value[0]?.roles?.[0]?.id || ''
    if (roleId.value) await selectRole(roleId.value)
  } finally {
    loading.value = false
  }
}

async function selectRole(id) {
  roleId.value = id
  const res = await ecoAdminApi.role(id)
  if (roleId.value === id) role.value = res.role
}

function openCareer(c) {
  careerDialog.id = c?.id || ''
  careerDialog.form = c ? { ...c } : { id: '', name: '', category: '研发', description: '', demand: '高', salaryRange: '', icon: 'code', color: '#6b5cff', sort: careers.value.length }
  careerDialog.visible = true
}

async function saveCareer() {
  const { roles, ...data } = careerDialog.form
  await ecoAdminApi.saveCareer(data, careerDialog.id)
  ElMessage.success('职业已保存')
  careerDialog.visible = false
  await load()
}

async function removeCareer(id) {
  await ecoAdminApi.removeCareer(id)
  ElMessage.success('已删除')
  if (careers.value.find((c) => c.id === id)?.roles?.some((r) => r.id === roleId.value)) {
    roleId.value = ''
    role.value = null
  }
  await load()
}

function openRole(r, careerId) {
  roleDialog.id = r?.id || ''
  roleDialog.form = r
    ? { careerId: r.careerId, name: r.name, level: r.level, description: r.description, sort: r.sort }
    : { id: '', careerId, name: '', level: '中级', description: '', sort: 0 }
  roleDialog.visible = true
}

async function saveRole() {
  const saved = await ecoAdminApi.saveRole(roleDialog.form, roleDialog.id)
  ElMessage.success('岗位已保存')
  roleDialog.visible = false
  roleId.value = saved.id
  await load()
}

async function removeRole(id) {
  await ecoAdminApi.removeRole(id)
  ElMessage.success('已删除')
  roleId.value = ''
  role.value = null
  await load()
}

function openCap(cap) {
  capDialog.id = cap?.id || ''
  capDialog.form = cap
    ? { roleId: role.value.id, name: cap.name, description: cap.description, weight: cap.weight, sort: cap.sort, skills: cap.skills.map((s) => ({ skillId: s.skillId, requiredLevel: s.requiredLevel, weight: s.weight })) }
    : { roleId: role.value.id, name: '', description: '', weight: 2, sort: role.value.capabilities?.length || 0, skills: [{ skillId: '', requiredLevel: 2, weight: 1 }] }
  capDialog.visible = true
}

async function saveCap() {
  await ecoAdminApi.saveCapability({ ...capDialog.form, skills: capDialog.form.skills.filter((s) => s.skillId) }, capDialog.id)
  ElMessage.success('能力已保存')
  capDialog.visible = false
  await selectRole(roleId.value)
}

async function removeCap(id) {
  await ecoAdminApi.removeCapability(id)
  ElMessage.success('已删除')
  await selectRole(roleId.value)
}

onMounted(async () => {
  skills.value = await ecoAdminApi.skills()
  await load()
})
</script>

<style scoped>
.page { display: flex; flex-direction: column; gap: 18px; }
.layout { display: grid; grid-template-columns: 320px minmax(0, 1fr); gap: 18px; align-items: start; }
.tree__body { padding: 12px; }
.tree__career { margin-bottom: 14px; }
.tree__career-head { display: flex; align-items: center; gap: 8px; padding: 4px 6px; font-size: 14px; }
.tree__career-head .muted { font-size: 12px; }
.tree__dot { width: 9px; height: 9px; border-radius: 50%; }
.tree__ops { margin-left: auto; display: inline-flex; }
.tree__role { display: flex; align-items: center; justify-content: space-between; width: 100%; margin: 2px 0; padding: 8px 10px 8px 24px; border: 1px solid transparent; border-radius: 10px; background: transparent; color: var(--text-2); font: inherit; font-size: 13.5px; text-align: left; cursor: pointer; }
.tree__role:hover { background: var(--surface-2); }
.tree__role.active { border-color: var(--primary-soft-2, #e2defc); background: var(--primary-soft); color: var(--primary); font-weight: 600; }
.tree__empty { padding: 4px 24px; font-size: 12.5px; }
.caps { display: flex; flex-direction: column; gap: 16px; }
.cap { padding: 14px; border: 1px solid var(--border); border-radius: var(--radius-sm); }
.cap__head { display: flex; justify-content: space-between; gap: 12px; margin-bottom: 10px; }
.cap__head strong { margin-right: 8px; font-size: 15px; }
.cap__head p { margin: 4px 0 0; font-size: 13px; }
.lvl { display: inline-block; margin-right: 6px; padding: 0 6px; border-radius: 6px; color: #fff; font-size: 11px; font-weight: 700; }
.req-list { display: flex; flex-direction: column; gap: 8px; width: 100%; }
.req { display: flex; align-items: center; gap: 8px; }
@media (max-width: 1100px) { .layout { grid-template-columns: 1fr; } }
</style>
