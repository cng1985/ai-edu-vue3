<template>
  <div class="page">
    <PageHeader eyebrow="成长生态" title="技能与知识图谱" subtitle="技能是可执行的能力，知识是人知道什么。课程结构与知识结构分离，通过章节映射复用知识点。">
      <el-button v-if="tab === 'skills'" type="primary" :icon="Plus" @click="openSkill()">新增技能</el-button>
      <el-button v-if="tab === 'knowledge'" type="primary" :icon="Plus" @click="openKnowledge()">新增知识点</el-button>
    </PageHeader>

    <div class="stat-grid">
      <StatCard label="技能" :value="skills.length" icon="Medal" tone="primary" />
      <StatCard label="知识点" :value="knowledge.length" icon="Collection" tone="info" />
      <StatCard label="图谱关系" :value="relations.length" icon="Share" tone="sky" />
      <StatCard label="练习题" :value="knowledge.reduce((s, k) => s + k.questionCount, 0)" icon="EditPen" tone="success" />
    </div>

    <div class="panel">
      <el-tabs v-model="tab" class="tabs">
        <el-tab-pane label="技能" name="skills">
          <el-table :data="skills">
            <el-table-column label="技能" min-width="200">
              <template #default="{ row }"><strong>{{ row.name }}</strong> <span class="muted mono">{{ row.id }}</span></template>
            </el-table-column>
            <el-table-column label="分类" prop="category" width="120" />
            <el-table-column label="描述" prop="description" min-width="260" show-overflow-tooltip />
            <el-table-column label="知识点" width="90">
              <template #default="{ row }"><span class="num">{{ knowledge.filter((k) => k.skillIds.includes(row.id)).length }}</span></template>
            </el-table-column>
            <el-table-column label="操作" width="140" align="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="openSkill(row)">编辑</el-button>
                <el-popconfirm title="删除技能会移除相关的能力要求与学员技能状态，确定？" @confirm="removeSkill(row.id)">
                  <template #reference><el-button link type="danger">删除</el-button></template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="知识点" name="knowledge">
          <div class="toolbar tab-toolbar">
            <el-input v-model="keyword" placeholder="搜索知识点" clearable style="width: 220px" />
            <el-select v-model="domain" placeholder="全部领域" clearable style="width: 160px">
              <el-option v-for="d in domains" :key="d" :label="d" :value="d" />
            </el-select>
          </div>
          <el-table :data="filteredKnowledge">
            <el-table-column label="知识点" min-width="220">
              <template #default="{ row }">
                <div><strong>{{ row.name }}</strong> <span class="muted mono">{{ row.id }}</span></div>
                <div class="muted small">{{ row.summary }}</div>
              </template>
            </el-table-column>
            <el-table-column label="领域" prop="domain" width="110" />
            <el-table-column label="难度" width="80"><template #default="{ row }">{{ '★'.repeat(row.difficulty) }}</template></el-table-column>
            <el-table-column label="时长" width="80"><template #default="{ row }">{{ row.estimatedMinutes }}′</template></el-table-column>
            <el-table-column label="关联技能" min-width="180">
              <template #default="{ row }">
                <el-tag v-for="sid in row.skillIds" :key="sid" size="small" effect="plain" class="tag-gap">{{ skillName(sid) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="练习题" width="80"><template #default="{ row }"><span class="num">{{ row.questionCount }}</span></template></el-table-column>
            <el-table-column label="操作" width="140" align="right">
              <template #default="{ row }">
                <el-button link type="primary" @click="openKnowledge(row)">编辑</el-button>
                <el-popconfirm title="确定删除该知识点及其图谱关系？" @confirm="removeKnowledge(row.id)">
                  <template #reference><el-button link type="danger">删除</el-button></template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="知识关系" name="relations">
          <div class="toolbar tab-toolbar">
            <el-select v-model="rel.fromId" filterable placeholder="起点知识" style="width: 220px">
              <el-option v-for="k in knowledge" :key="k.id" :label="k.name" :value="k.id" />
            </el-select>
            <el-select v-model="rel.type" style="width: 140px">
              <el-option v-for="(label, k) in RELATION_TYPES" :key="k" :label="label" :value="k" />
            </el-select>
            <el-select v-model="rel.toId" filterable placeholder="终点知识" style="width: 220px">
              <el-option v-for="k in knowledge" :key="k.id" :label="k.name" :value="k.id" />
            </el-select>
            <el-button type="primary" :icon="Plus" @click="addRelation">添加关系</el-button>
            <span class="muted small">前置：起点是终点的前置；高级：终点是起点的进阶。系统会阻止形成环路。</span>
          </div>
          <el-table :data="relations">
            <el-table-column label="起点" min-width="200"><template #default="{ row }">{{ knowledgeName(row.fromId) }}</template></el-table-column>
            <el-table-column label="关系" width="140">
              <template #default="{ row }"><el-tag size="small" :type="row.type === 'prerequisite' ? 'primary' : row.type === 'advanced' ? 'warning' : 'info'">{{ RELATION_TYPES[row.type] }}</el-tag></template>
            </el-table-column>
            <el-table-column label="终点" min-width="200"><template #default="{ row }">{{ knowledgeName(row.toId) }}</template></el-table-column>
            <el-table-column label="操作" width="100" align="right">
              <template #default="{ row }">
                <el-popconfirm title="确定删除该关系？" @confirm="removeRelation(row.id)">
                  <template #reference><el-button link type="danger">删除</el-button></template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="课程章节映射" name="mapping">
          <div class="toolbar tab-toolbar">
            <el-select v-model="courseId" placeholder="选择课程" style="width: 260px" @change="loadMappings">
              <el-option v-for="c in courses" :key="c.id" :label="c.title" :value="c.id" />
            </el-select>
            <span class="muted small">学员完成章节或测验后，会按映射更新对应知识点的学习状态。</span>
          </div>
          <el-table v-if="course" :data="course.chapters">
            <el-table-column label="章节" min-width="220"><template #default="{ row }">{{ row.title }} <span class="muted mono">{{ row.id }}</span></template></el-table-column>
            <el-table-column label="映射知识点" min-width="360">
              <template #default="{ row }">
                <el-select v-model="mapping[row.id]" multiple filterable placeholder="选择知识点" style="width: 100%">
                  <el-option v-for="k in knowledge" :key="k.id" :label="`${k.domain} · ${k.name}`" :value="k.id" />
                </el-select>
              </template>
            </el-table-column>
            <el-table-column label="操作" width="100" align="right">
              <template #default="{ row }"><el-button link type="primary" @click="saveMapping(row.id)">保存</el-button></template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </div>

    <el-dialog v-model="skillDialog.visible" :title="skillDialog.id ? '编辑技能' : '新增技能'" width="520px">
      <el-form :model="skillDialog.form" label-width="80px">
        <el-form-item v-if="!skillDialog.id" label="标识"><el-input v-model="skillDialog.form.id" placeholder="如 redis，留空自动生成" /></el-form-item>
        <el-form-item label="名称" required><el-input v-model="skillDialog.form.name" /></el-form-item>
        <el-form-item label="分类"><el-input v-model="skillDialog.form.category" placeholder="后端开发 / AI 应用 …" /></el-form-item>
        <el-form-item label="描述"><el-input v-model="skillDialog.form.description" type="textarea" :rows="3" /></el-form-item>
        <el-form-item label="关键词">
          <el-select v-model="skillDialog.form.tags" multiple filterable allow-create default-first-option placeholder="用于规则评审识别技能" style="width: 100%" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="skillDialog.visible = false">取消</el-button>
        <el-button type="primary" @click="saveSkill">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="kDialog.visible" :title="kDialog.id ? '编辑知识点' : '新增知识点'" width="780px" top="5vh">
      <el-form :model="kDialog.form" label-width="88px">
        <div class="form-row">
          <el-form-item v-if="!kDialog.id" label="标识"><el-input v-model="kDialog.form.id" placeholder="留空自动生成" /></el-form-item>
          <el-form-item label="名称" required><el-input v-model="kDialog.form.name" /></el-form-item>
          <el-form-item label="领域"><el-input v-model="kDialog.form.domain" placeholder="Java 并发" /></el-form-item>
        </div>
        <div class="form-row">
          <el-form-item label="难度"><el-rate v-model="kDialog.form.difficulty" :max="5" /></el-form-item>
          <el-form-item label="时长(分)"><el-input-number v-model="kDialog.form.estimatedMinutes" :min="5" :max="240" /></el-form-item>
        </div>
        <el-form-item label="关联技能">
          <el-select v-model="kDialog.form.skillIds" multiple filterable style="width: 100%">
            <el-option v-for="s in skills" :key="s.id" :label="`${s.category} · ${s.name}`" :value="s.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="摘要"><el-input v-model="kDialog.form.summary" /></el-form-item>
        <el-form-item label="正文"><el-input v-model="kDialog.form.content" type="textarea" :rows="8" placeholder="支持 Markdown" /></el-form-item>
        <el-form-item label="练习题">
          <div class="questions">
            <div v-for="(q, i) in kDialog.form.questions" :key="i" class="question">
              <div class="question__head">
                <el-input v-model="q.text" placeholder="题干" />
                <el-button link type="danger" :icon="Delete" @click="kDialog.form.questions.splice(i, 1)" />
              </div>
              <el-radio-group v-model="q.answer" class="question__opts">
                <div v-for="(opt, j) in q.options" :key="j" class="question__opt">
                  <el-radio :value="j">{{ String.fromCharCode(65 + j) }}</el-radio>
                  <el-input v-model="q.options[j]" size="small" />
                </div>
              </el-radio-group>
              <el-input v-model="q.explanation" size="small" placeholder="答案解析" />
            </div>
            <el-button :icon="Plus" size="small" @click="kDialog.form.questions.push({ text: '', options: ['', '', '', ''], answer: 0, explanation: '' })">添加练习题</el-button>
          </div>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="kDialog.visible = false">取消</el-button>
        <el-button type="primary" @click="saveKnowledge">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
import { computed, onMounted, reactive, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { Delete, Plus } from '@element-plus/icons-vue'
import PageHeader from '../../components/common/PageHeader.vue'
import StatCard from '../../components/common/StatCard.vue'
import { coursesApi, ecoAdminApi } from '../../api'
import { RELATION_TYPES } from '../../utils/eco'

const tab = ref('skills')
const skills = ref([])
const knowledge = ref([])
const relations = ref([])
const courses = ref([])
const keyword = ref('')
const domain = ref('')
const rel = reactive({ fromId: '', toId: '', type: 'prerequisite' })
const courseId = ref('')
const mapping = ref({})
const skillDialog = reactive({ visible: false, id: '', form: {} })
const kDialog = reactive({ visible: false, id: '', form: {} })

const domains = computed(() => [...new Set(knowledge.value.map((k) => k.domain))])
const filteredKnowledge = computed(() =>
  knowledge.value.filter(
    (k) => (!domain.value || k.domain === domain.value) && (!keyword.value || k.name.includes(keyword.value) || k.id.includes(keyword.value))
  )
)
const course = computed(() => courses.value.find((c) => c.id === courseId.value))

const skillName = (id) => skills.value.find((s) => s.id === id)?.name || id
const knowledgeName = (id) => knowledge.value.find((k) => k.id === id)?.name || id

async function loadAll() {
  const [s, k, r] = await Promise.all([ecoAdminApi.skills(), ecoAdminApi.knowledge(), ecoAdminApi.relations()])
  skills.value = s
  knowledge.value = k
  relations.value = r
}

function openSkill(s) {
  skillDialog.id = s?.id || ''
  skillDialog.form = s ? { ...s, tags: [...(s.tags || [])] } : { id: '', name: '', category: '', description: '', tags: [] }
  skillDialog.visible = true
}

async function saveSkill() {
  await ecoAdminApi.saveSkill(skillDialog.form, skillDialog.id)
  ElMessage.success('技能已保存')
  skillDialog.visible = false
  await loadAll()
}

async function removeSkill(id) {
  await ecoAdminApi.removeSkill(id)
  ElMessage.success('已删除')
  await loadAll()
}

function openKnowledge(k) {
  kDialog.id = k?.id || ''
  kDialog.form = k
    ? JSON.parse(JSON.stringify({ ...k, questions: k.questions || [], skillIds: k.skillIds || [] }))
    : { id: '', name: '', domain: '', summary: '', content: '', difficulty: 2, estimatedMinutes: 20, skillIds: [], questions: [], tags: [] }
  kDialog.visible = true
}

async function saveKnowledge() {
  const form = { ...kDialog.form, questions: kDialog.form.questions.filter((q) => q.text.trim()) }
  await ecoAdminApi.saveKnowledge(form, kDialog.id)
  ElMessage.success('知识点已保存')
  kDialog.visible = false
  await loadAll()
}

async function removeKnowledge(id) {
  await ecoAdminApi.removeKnowledge(id)
  ElMessage.success('已删除')
  await loadAll()
}

async function addRelation() {
  if (!rel.fromId || !rel.toId) return ElMessage.warning('请选择起点与终点知识')
  await ecoAdminApi.addRelation({ ...rel })
  ElMessage.success('关系已添加')
  rel.fromId = ''
  rel.toId = ''
  relations.value = await ecoAdminApi.relations()
}

async function removeRelation(id) {
  await ecoAdminApi.removeRelation(id)
  relations.value = await ecoAdminApi.relations()
}

async function loadMappings() {
  const list = await ecoAdminApi.chapterMappings(courseId.value)
  const m = {}
  for (const ch of course.value?.chapters || []) m[ch.id] = []
  for (const it of list) (m[it.chapterId] ||= []).push(it.knowledgeId)
  mapping.value = m
}

async function saveMapping(chapterId) {
  await ecoAdminApi.setChapterKnowledge(courseId.value, chapterId, mapping.value[chapterId] || [])
  ElMessage.success('章节映射已保存')
}

onMounted(async () => {
  await loadAll()
  const res = await coursesApi.list({ pageSize: 100 })
  courses.value = res.list || []
  if (courses.value[0]) {
    courseId.value = courses.value[0].id
    await loadMappings()
  }
})
</script>

<style scoped>
.page { display: flex; flex-direction: column; gap: 18px; }
.tabs { padding: 6px 22px 22px; }
.tab-toolbar { margin-bottom: 14px; }
.small { font-size: 12.5px; }
.tag-gap { margin: 2px 4px 2px 0; }
.form-row { display: grid; grid-template-columns: repeat(3, 1fr); gap: 0 12px; }
.questions { display: flex; flex-direction: column; gap: 12px; width: 100%; }
.question { display: flex; flex-direction: column; gap: 8px; padding: 12px; border: 1px solid var(--border); border-radius: 10px; background: var(--surface-2); }
.question__head { display: flex; gap: 8px; }
.question__opts { display: grid; grid-template-columns: 1fr 1fr; gap: 6px 12px; width: 100%; }
.question__opt { display: flex; align-items: center; gap: 4px; }
.question__opt .el-radio { margin-right: 0; }
</style>
