<template>
  <div class="page">
    <PageHeader eyebrow="用户管理" title="用户管理" subtitle="管理学习者、创作者、企业与后台账号，分配角色并控制启用状态。">
      <el-button v-permission="PERM.USER_CREATE" type="primary" :icon="Plus" @click="openDialog()">新增用户</el-button>
    </PageHeader>

    <div class="panel">
      <div class="panel__head">
        <div class="toolbar">
          <el-input v-model="filters.keyword" placeholder="搜索用户名 / 昵称" clearable :prefix-icon="Search" style="width: 240px" @clear="loadData" @keyup.enter="loadData" />
          <el-select v-model="filters.role" placeholder="全部角色" clearable style="width: 140px" @change="loadData">
            <el-option v-for="(label, value) in roleMap" :key="value" :label="label" :value="value" />
          </el-select>
          <el-select v-model="filters.status" placeholder="全部状态" clearable style="width: 130px" @change="loadData">
            <el-option label="正常" value="active" />
            <el-option label="禁用" value="disabled" />
          </el-select>
          <el-button @click="loadData">查询</el-button>
        </div>
        <span class="muted">共 {{ total }} 位用户</span>
      </div>

      <div class="panel__body panel__body--flush">
        <el-table :data="list" v-loading="loading">
          <el-table-column label="用户" min-width="240">
            <template #default="{ row }">
              <div class="cell">
                <span class="avatar" :style="{ background: row.avatarColor || avatarColor(row.username) }">
                  {{ row.avatar || row.nickname?.slice(0, 1) || row.username?.slice(0, 1) }}
                </span>
                <div class="cell__main">
                  <div class="cell__title">{{ row.nickname || row.username }}</div>
                  <div class="cell__sub mono">@{{ row.username }}</div>
                </div>
              </div>
            </template>
          </el-table-column>
          <el-table-column label="角色" width="120">
            <template #default="{ row }">
              <el-tag :type="roleTagType(row.role)" size="small">{{ roleLabel(row.role) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="状态" width="110">
            <template #default="{ row }">
              <span class="status" :class="row.status === 'active' ? 'status--on' : 'status--off'">
                <i />{{ row.status === 'active' ? '正常' : '禁用' }}
              </span>
            </template>
          </el-table-column>
          <el-table-column label="注册时间" min-width="180">
            <template #default="{ row }"><span class="muted">{{ formatDate(row.joinedAt) }}</span></template>
          </el-table-column>
          <el-table-column label="操作" width="150" fixed="right" align="right">
            <template #default="{ row }">
              <el-button link type="primary" @click="openDialog(row)">编辑</el-button>
              <el-popconfirm v-if="auth.hasPermission(PERM.USER_DELETE)" title="确定删除该用户？" @confirm="handleDelete(row.id)">
                <template #reference>
                  <el-button link type="danger">删除</el-button>
                </template>
              </el-popconfirm>
            </template>
          </el-table-column>
          <template #empty>
            <el-empty description="没有符合条件的用户" :image-size="90" />
          </template>
        </el-table>
      </div>

      <div v-if="total > pageSize" class="panel__foot pager">
        <el-pagination
          background
          layout="total, prev, pager, next"
          :total="total"
          :page-size="pageSize"
          v-model:current-page="page"
          @current-change="loadData"
        />
      </div>
    </div>

    <el-dialog v-model="dialogVisible" :title="editing ? '编辑用户' : '新增用户'" width="480px">
      <el-form ref="formRef" :model="form" :rules="formRules" label-width="80px">
        <el-form-item label="用户名" prop="username">
          <el-input v-model="form.username" :disabled="editing" />
        </el-form-item>
        <el-form-item label="昵称" prop="nickname">
          <el-input v-model="form.nickname" />
        </el-form-item>
        <el-form-item :label="editing ? '新密码' : '密码'" :prop="editing ? '' : 'password'">
          <el-input v-model="form.password" type="password" show-password :placeholder="editing ? '留空则不修改' : ''" />
        </el-form-item>
        <el-form-item label="角色" prop="role">
          <el-select v-model="form.role" style="width: 100%">
            <el-option v-for="(label, value) in roleMap" :key="value" :label="label" :value="value" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="form.status" style="width: 100%">
            <el-option label="正常" value="active" />
            <el-option label="禁用" value="disabled" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="handleSave">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { Plus, Search } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { usersApi } from '../api'
import { useAuthStore } from '../stores/auth'
import { PERM } from '../constants/permissions'
import PageHeader from '../components/common/PageHeader.vue'

const AVATAR_COLORS = ['#6b5cff', '#0fb981', '#f59e0b', '#2a8cf4', '#f0589a', '#0ea5e9']
function avatarColor(seed = '') {
  let h = 0
  for (const ch of seed) h = (h * 31 + ch.charCodeAt(0)) >>> 0
  return AVATAR_COLORS[h % AVATAR_COLORS.length]
}

const auth = useAuthStore()
const loading = ref(false)
const saving = ref(false)
const list = ref([])
const total = ref(0)
const page = ref(1)
const pageSize = 20
const filters = reactive({ keyword: '', role: '', status: '' })

const dialogVisible = ref(false)
const editing = ref(false)
const editingId = ref('')
const formRef = ref()
const form = reactive({ username: '', nickname: '', password: '', role: 'learner', status: 'active' })
const formRules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

const roleMap = { learner: '学习者', creator: '创作者', enterprise: '企业', admin: '管理员', reviewer: '审核员', operator: '运营' }
function roleLabel(r) { return roleMap[r] || r }
function roleTagType(r) {
  if (r === 'admin') return 'danger'
  if (r === 'reviewer') return 'warning'
  if (r === 'operator') return 'info'
  return ''
}

function formatDate(ts) {
  return new Date(ts).toLocaleString('zh-CN')
}

async function loadData() {
  loading.value = true
  try {
    const data = await usersApi.list({ ...filters, page: page.value, pageSize })
    list.value = data.list
    total.value = data.total
  } finally {
    loading.value = false
  }
}

function openDialog(row) {
  editing.value = Boolean(row)
  editingId.value = row?.id || ''
  Object.assign(form, {
    username: row?.username || '',
    nickname: row?.nickname || '',
    password: '',
    role: row?.role || 'learner',
    status: row?.status || 'active'
  })
  dialogVisible.value = true
}

async function handleSave() {
  await formRef.value.validate()
  saving.value = true
  try {
    if (editing.value) {
      const payload = { nickname: form.nickname, role: form.role, status: form.status }
      if (form.password) payload.password = form.password
      await usersApi.update(editingId.value, payload)
    } else {
      await usersApi.create(form)
    }
    ElMessage.success('保存成功')
    dialogVisible.value = false
    loadData()
  } finally {
    saving.value = false
  }
}

async function handleDelete(id) {
  await usersApi.remove(id)
  ElMessage.success('删除成功')
  loadData()
}

onMounted(loadData)
</script>

<style scoped>
.avatar {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 38px;
  height: 38px;
  border-radius: 12px;
  color: #fff;
  font-weight: 700;
  font-size: 14px;
  flex-shrink: 0;
}

.status {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  font-size: 13px;
  font-weight: 600;
}

.status i {
  width: 8px;
  height: 8px;
  border-radius: 50%;
}

.status--on { color: var(--success-strong); }
.status--on i { background: var(--success); box-shadow: 0 0 0 3px var(--success-soft); }
.status--off { color: var(--danger-strong); }
.status--off i { background: var(--danger); box-shadow: 0 0 0 3px var(--danger-soft); }

.pager {
  display: flex;
  justify-content: flex-end;
}
</style>
