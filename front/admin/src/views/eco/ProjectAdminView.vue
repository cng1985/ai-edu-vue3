<template>
  <div class="page">
    <PageHeader eyebrow="成长生态" title="项目实践" subtitle="Project → Task → Execution → Evidence：维护实践项目与任务技能要求，查看学员提交与 Review Agent 评审结果。">
      <el-button type="primary" :icon="Plus" @click="openProject()">新增项目</el-button>
    </PageHeader>

    <div class="panel">
      <el-tabs v-model="tab" class="tabs">
        <el-tab-pane label="项目" name="projects">
          <el-table :data="projects" v-loading="loading">
            <el-table-column label="项目" min-width="260">
              <template #default="{ row }">
                <div class="cell">
                  <span class="dot" :style="{ background: row.color || '#6b5cff' }" />
                  <div><strong>{{ row.title }}</strong><div class="muted small">{{ row.summary }}</div></div>
                </div>
              </template>
            </el-table-column>
            <el-table-column label="领域" prop="domain" width="110" />
            <el-table-column label="难度" width="90"><template #default="{ row }">{{ '★'.repeat(row.difficulty) }}</template></el-table-column>
            <el-table-column label="任务" min-width="260">
              <template #default="{ row }">
                <el-tag v-for="t in row.tasks" :key="t.id" size="small" effect="plain" class="tag-gap">{{ t.title }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="90">
              <template #default="{ row }"><el-tag size="small" :type="row.status === 'published' ? 'success' : 'info'">{{ row.status === 'published' ? '已发布' : '草稿' }}</el-tag></template>
            </el-table-column>
            <el-table-column label="操作" width="150" align="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="openProject(row)">编辑</el-button>
                <el-popconfirm title="确定删除该项目？" @confirm="remove(row.id)">
                  <template #reference><el-button link type="danger">删除</el-button></template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="学员提交" name="submissions">
          <el-table :data="submissions">
            <el-table-column label="项目 / 任务" min-width="240">
              <template #default="{ row }"><strong>{{ row.taskTitle }}</strong><div class="muted small">{{ row.projectTitle }}</div></template>
            </el-table-column>
            <el-table-column label="学员" prop="userId" width="150" />
            <el-table-column label="得分" width="90">
              <template #default="{ row }"><span class="num score" :class="{ ok: row.score >= 60 }">{{ row.score }}</span></template>
            </el-table-column>
            <el-table-column label="评审" width="110">
              <template #default="{ row }"><el-tag size="small" :type="row.evaluator === 'ai' ? 'primary' : 'info'">{{ row.evaluator === 'ai' ? 'AI 评审' : '规则评审' }}</el-tag></template>
            </el-table-column>
            <el-table-column label="提交时间" width="180"><template #default="{ row }">{{ dateTime(row.createdAt) }}</template></el-table-column>
            <el-table-column type="expand">
              <template #default="{ row }">
                <div class="expand">
                  <h4>提交内容</h4><pre>{{ row.content }}</pre>
                  <h4>评审反馈</h4><pre>{{ row.feedback }}</pre>
                </div>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </div>

    <el-dialog v-model="dialog.visible" :title="dialog.id ? '编辑项目' : '新增项目'" width="860px" top="4vh">
      <el-form :model="dialog.form" label-width="88px">
        <div class="form-row">
          <el-form-item label="项目名称" required><el-input v-model="dialog.form.title" /></el-form-item>
          <el-form-item label="领域"><el-input v-model="dialog.form.domain" placeholder="后端开发 / AI 应用" /></el-form-item>
        </div>
        <el-form-item label="一句话简介"><el-input v-model="dialog.form.summary" /></el-form-item>
        <el-form-item label="项目说明"><el-input v-model="dialog.form.description" type="textarea" :rows="5" placeholder="Markdown：项目背景、你将实践什么" /></el-form-item>
        <div class="form-row">
          <el-form-item label="难度"><el-rate v-model="dialog.form.difficulty" :max="5" /></el-form-item>
          <el-form-item label="状态">
            <el-radio-group v-model="dialog.form.status"><el-radio value="published">发布</el-radio><el-radio value="draft">草稿</el-radio></el-radio-group>
          </el-form-item>
          <el-form-item label="主题色"><el-color-picker v-model="dialog.form.color" /></el-form-item>
        </div>
        <el-form-item label="涉及技能">
          <el-select v-model="dialog.form.skillIds" multiple filterable style="width: 100%">
            <el-option v-for="s in skills" :key="s.id" :label="s.name" :value="s.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="任务链">
          <div class="tasks">
            <div v-for="(t, i) in dialog.form.tasks" :key="i" class="task">
              <div class="task__head">
                <span class="task__idx">{{ i + 1 }}</span>
                <el-input v-model="t.title" placeholder="任务名称，如：库存扣减" />
                <el-input-number v-model="t.estimatedHours" :min="1" :max="200" style="width: 130px" />
                <span class="muted small">小时</span>
                <el-button link type="danger" :icon="Delete" @click="dialog.form.tasks.splice(i, 1)" />
              </div>
              <el-input v-model="t.description" type="textarea" :rows="2" placeholder="任务说明" />
              <el-input v-model="t.deliverable" placeholder="交付物" />
              <div class="reqs">
                <div v-for="(r, j) in t.requirements" :key="j" class="req">
                  <el-select v-model="r.skillId" filterable placeholder="技能" style="width: 200px">
                    <el-option v-for="s in skills" :key="s.id" :label="s.name" :value="s.id" />
                  </el-select>
                  <el-select v-model="r.level" style="width: 90px"><el-option v-for="l in 5" :key="l" :label="`L${l}`" :value="l" /></el-select>
                  <el-button link :icon="Close" @click="t.requirements.splice(j, 1)" />
                </div>
                <el-button size="small" text :icon="Plus" @click="t.requirements.push({ skillId: '', level: 2 })">技能要求</el-button>
              </div>
            </div>
            <el-button :icon="Plus" size="small" @click="dialog.form.tasks.push({ title: '', description: '', deliverable: '', estimatedHours: 4, requirements: [] })">添加任务</el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialog.visible = false">取消</el-button>
        <el-button type="primary" @click="save">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Close, Delete, Plus } from '@element-plus/icons-vue'
import PageHeader from '../../components/common/PageHeader.vue'
import { ecoAdminApi } from '../../api'
import { dateTime } from '../../utils/eco'

const tab = ref('projects')
const projects = ref([])
const submissions = ref([])
const skills = ref([])
const loading = ref(false)
const dialog = reactive({ visible: false, id: '', form: {} })

async function load() {
  loading.value = true
  try {
    projects.value = await ecoAdminApi.projects()
    submissions.value = await ecoAdminApi.submissions({ limit: 200 })
  } finally {
    loading.value = false
  }
}

async function openProject(p) {
  dialog.id = p?.id || ''
  if (p) {
    const full = await ecoAdminApi.project(p.id)
    dialog.form = {
      title: full.title, summary: full.summary, description: full.description, domain: full.domain,
      difficulty: full.difficulty, icon: full.icon, color: full.color, status: full.status, skillIds: [...full.skillIds],
      tasks: full.tasks.map((t) => ({ id: t.id, title: t.title, description: t.description, deliverable: t.deliverable, estimatedHours: t.estimatedHours, requirements: t.requirements.map((r) => ({ ...r })) }))
    }
  } else {
    dialog.form = { title: '', summary: '', description: '', domain: '', difficulty: 3, icon: 'layers', color: '#6b5cff', status: 'published', skillIds: [], tasks: [] }
  }
  dialog.visible = true
}

async function save() {
  const form = { ...dialog.form, tasks: dialog.form.tasks.map((t) => ({ ...t, requirements: t.requirements.filter((r) => r.skillId) })) }
  await ecoAdminApi.saveProject(form, dialog.id)
  ElMessage.success('项目已保存')
  dialog.visible = false
  await load()
}

async function remove(id) {
  await ecoAdminApi.removeProject(id)
  ElMessage.success('已删除')
  await load()
}

onMounted(async () => {
  skills.value = await ecoAdminApi.skills()
  await load()
})
</script>

<style scoped>
.page { display: flex; flex-direction: column; gap: 18px; }
.tabs { padding: 6px 22px 22px; }
.small { font-size: 12.5px; }
.dot { width: 10px; height: 10px; flex: 0 0 10px; border-radius: 50%; }
.tag-gap { margin: 2px 4px 2px 0; }
.score { font-weight: 700; color: var(--warning-strong); }
.score.ok { color: var(--success-strong); }
.expand { padding: 4px 24px; }
.expand h4 { margin: 8px 0 4px; }
.expand pre { max-height: 220px; overflow: auto; margin: 0; padding: 10px; border-radius: 8px; background: var(--surface-2); white-space: pre-wrap; font-family: inherit; font-size: 13px; }
.form-row { display: grid; grid-template-columns: repeat(3, 1fr); gap: 0 12px; }
.tasks { display: flex; flex-direction: column; gap: 12px; width: 100%; }
.task { display: flex; flex-direction: column; gap: 8px; padding: 12px; border: 1px solid var(--border); border-radius: 10px; background: var(--surface-2); }
.task__head { display: flex; align-items: center; gap: 8px; }
.task__idx { display: inline-grid; place-items: center; width: 24px; height: 24px; flex: 0 0 24px; border-radius: 50%; background: var(--primary); color: #fff; font-size: 12px; font-weight: 700; }
.reqs { display: flex; flex-wrap: wrap; align-items: center; gap: 8px; }
.req { display: flex; align-items: center; gap: 4px; }
</style>
