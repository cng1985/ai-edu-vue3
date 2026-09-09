<template>
  <div class="page">
    <PageHeader eyebrow="系统管理" title="权限管理" subtitle="按角色分配后台功能权限，修改后即时生效。" />

    <div v-loading="loading" class="roles">
      <div v-for="role in roles" :key="role.role" class="role" :class="`role--${role.role}`">
        <div class="role__head">
          <span class="role__icon">
            <el-icon :size="20"><component :is="roleIcon(role.role)" /></el-icon>
          </span>
          <div class="role__title">
            <h3>{{ role.name }}</h3>
            <span class="mono">{{ role.role }}</span>
          </div>
          <el-button v-permission="PERM.ROLE_MANAGE" size="small" @click="openEdit(role)">
            <el-icon class="el-icon--left"><EditPen /></el-icon>编辑
          </el-button>
        </div>

        <div class="role__meter">
          <div class="role__meter-label">
            <span>已授权</span>
            <span class="num">{{ role.permissions.length }} / {{ permissions.length }}</span>
          </div>
          <div class="role__meter-track">
            <span :style="{ width: ratio(role) }" />
          </div>
        </div>

        <div class="role__perms">
          <span v-for="p in role.permissions" :key="p" class="perm">{{ permLabel(p) }}</span>
          <span v-if="!role.permissions.length" class="muted">尚未分配任何权限</span>
        </div>
      </div>
    </div>

    <el-dialog v-model="dialogVisible" :title="`编辑权限 · ${editing?.name}`" width="640px">
      <div class="perm-summary">
        已选择 <strong class="num">{{ selectedPerms.length }}</strong> 项权限
        <el-button link type="primary" size="small" @click="selectedPerms = permissions.map((p) => p.code)">全选</el-button>
        <el-button link size="small" @click="selectedPerms = []">清空</el-button>
      </div>
      <el-checkbox-group v-model="selectedPerms" class="perm-groups">
        <div v-for="group in permGroups" :key="group.name" class="perm-group">
          <div class="perm-group__head">
            <strong>{{ group.name }}</strong>
            <span class="muted num">{{ groupSelected(group) }} / {{ group.items.length }}</span>
          </div>
          <div class="perm-group__items">
            <el-checkbox v-for="p in group.items" :key="p.code" :value="p.code" :label="p.code" class="perm-check">
              {{ p.name }}
            </el-checkbox>
          </div>
        </div>
      </el-checkbox-group>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { EditPen, Avatar, DocumentChecked, Promotion, User } from '@element-plus/icons-vue'
import { rolesApi } from '../api/roles.js'
import { useAuthStore } from '../stores/auth'
import { PERM } from '../constants/permissions'
import PageHeader from '../components/common/PageHeader.vue'

const auth = useAuthStore()
const loading = ref(false)
const saving = ref(false)
const roles = ref([])
const permissions = ref([])
const dialogVisible = ref(false)
const editing = ref(null)
const selectedPerms = ref([])

const permMap = computed(() => Object.fromEntries(permissions.value.map((p) => [p.code, p.name])))
const permGroups = computed(() => {
  const groups = {}
  for (const p of permissions.value) {
    if (!groups[p.group]) groups[p.group] = { name: p.group, items: [] }
    groups[p.group].items.push(p)
  }
  return Object.values(groups)
})

function permLabel(code) {
  return permMap.value[code] || code
}

function roleIcon(role) {
  return { admin: Avatar, reviewer: DocumentChecked, operator: Promotion }[role] || User
}

function ratio(role) {
  if (!permissions.value.length) return '0%'
  return `${Math.round((role.permissions.length / permissions.value.length) * 100)}%`
}

function groupSelected(group) {
  return group.items.filter((p) => selectedPerms.value.includes(p.code)).length
}

async function loadData() {
  loading.value = true
  try {
    const [r, p] = await Promise.all([rolesApi.list(), rolesApi.listPermissions()])
    roles.value = r
    permissions.value = p
  } finally {
    loading.value = false
  }
}

function openEdit(role) {
  editing.value = role
  selectedPerms.value = [...role.permissions]
  dialogVisible.value = true
}

async function handleSave() {
  saving.value = true
  try {
    await rolesApi.update(editing.value.role, selectedPerms.value)
    ElMessage.success('权限已更新')
    if (editing.value.role === auth.user?.role) {
      await auth.refreshPermissions()
    }
    dialogVisible.value = false
    loadData()
  } finally {
    saving.value = false
  }
}

onMounted(loadData)
</script>

<style scoped>
.roles {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 18px;
  min-height: 200px;
}

.role {
  --tone: var(--primary);
  --tone-soft: var(--primary-soft);
  display: flex;
  flex-direction: column;
  gap: 18px;
  padding: 22px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius);
  box-shadow: var(--shadow-sm);
  transition: transform var(--t), box-shadow var(--t);
}

.role:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow);
}

.role--admin { --tone: var(--danger); --tone-soft: var(--danger-soft); }
.role--reviewer { --tone: var(--warning); --tone-soft: var(--warning-soft); }
.role--operator { --tone: var(--info); --tone-soft: var(--info-soft); }

.role__head {
  display: flex;
  align-items: center;
  gap: 12px;
}

.role__icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  border-radius: 13px;
  background: var(--tone-soft);
  color: var(--tone);
  flex-shrink: 0;
}

.role__title {
  flex: 1;
  min-width: 0;
  line-height: 1.2;
}

.role__title h3 {
  margin: 0;
  font-size: 16px;
}

.role__title span {
  color: var(--text-3);
}

.role__meter-label {
  display: flex;
  justify-content: space-between;
  font-size: 12.5px;
  color: var(--text-3);
  font-weight: 600;
  margin-bottom: 8px;
}

.role__meter-track {
  height: 8px;
  border-radius: 999px;
  background: var(--surface-3);
  overflow: hidden;
}

.role__meter-track span {
  display: block;
  height: 100%;
  border-radius: 999px;
  background: linear-gradient(90deg, var(--tone), color-mix(in srgb, var(--tone) 60%, #fff));
  transition: width 0.6s var(--ease);
}

.role__perms {
  display: flex;
  flex-wrap: wrap;
  gap: 6px;
}

.perm {
  display: inline-flex;
  align-items: center;
  padding: 3px 10px;
  border-radius: 999px;
  background: var(--surface-2);
  border: 1px solid var(--border);
  font-size: 12px;
  font-weight: 600;
  color: var(--text-2);
}

.perm-summary {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 14px;
  margin-bottom: 16px;
  border-radius: var(--radius-sm);
  background: var(--primary-soft);
  color: var(--primary-deep);
  font-size: 13px;
  font-weight: 600;
}

.perm-summary strong {
  font-size: 16px;
}

.perm-groups {
  display: flex;
  flex-direction: column;
  gap: 14px;
  max-height: 460px;
  overflow-y: auto;
  padding-right: 4px;
}

.perm-group {
  border: 1px solid var(--border);
  border-radius: var(--radius-sm);
  overflow: hidden;
}

.perm-group__head {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 10px 14px;
  background: var(--surface-2);
  font-size: 13.5px;
  color: var(--text);
}

.perm-group__items {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 4px 12px;
  padding: 10px 14px;
}

.perm-check {
  margin-right: 0;
  height: 32px;
}

@media (max-width: 1100px) {
  .roles {
    grid-template-columns: 1fr;
  }
}
</style>
