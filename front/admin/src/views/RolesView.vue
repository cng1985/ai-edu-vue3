<template>
  <div>
    <div class="page-header">
      <h2>权限管理</h2>
    </div>

    <el-row :gutter="16" v-loading="loading">
      <el-col :span="8" v-for="role in roles" :key="role.role">
        <el-card shadow="hover" style="margin-bottom: 16px">
          <template #header>
            <div style="display: flex; justify-content: space-between; align-items: center">
              <span>{{ role.name }} ({{ role.role }})</span>
              <el-button v-if="auth.isAdmin" size="small" type="primary" @click="openEdit(role)">编辑</el-button>
            </div>
          </template>
          <el-tag v-for="p in role.permissions" :key="p" size="small" style="margin: 2px">{{ permLabel(p) }}</el-tag>
        </el-card>
      </el-col>
    </el-row>

    <el-dialog v-model="dialogVisible" :title="`编辑权限 - ${editing?.name}`" width="560px">
      <el-checkbox-group v-model="selectedPerms">
        <div v-for="group in permGroups" :key="group.name" style="margin-bottom: 16px">
          <div style="font-weight: 600; margin-bottom: 8px">{{ group.name }}</div>
          <el-checkbox v-for="p in group.items" :key="p.code" :value="p.code" :label="p.code">{{ p.name }}</el-checkbox>
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
import { rolesApi } from '../api/roles.js'
import { useAuthStore } from '../stores/auth'

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
    dialogVisible.value = false
    loadData()
  } finally {
    saving.value = false
  }
}

onMounted(loadData)
</script>
