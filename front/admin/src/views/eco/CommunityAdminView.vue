<template>
  <div class="page">
    <PageHeader eyebrow="任务与人才" title="社区与知识资产" subtitle="知识生产 → 知识传播：治理社区内容与统一知识资产（文章、课程、视频、代码、案例、项目、任务、Prompt、SOP）。" />

    <div class="panel">
      <el-tabs v-model="tab" class="tabs" @tab-change="load">
        <el-tab-pane label="社区内容" name="posts">
          <div class="toolbar bar">
            <el-select v-model="postType" placeholder="全部类型" clearable style="width: 130px" @change="load">
              <el-option v-for="(label, k) in POST_TYPES" :key="k" :label="label" :value="k" />
            </el-select>
            <el-input v-model="keyword" placeholder="搜索标题或正文" clearable style="width: 220px" @keyup.enter="load" @clear="load" />
          </div>
          <el-table :data="posts.list" v-loading="loading">
            <el-table-column label="内容" min-width="300">
              <template #default="{ row }"><strong>{{ row.title }}</strong><div class="muted small">{{ row.content }}</div></template>
            </el-table-column>
            <el-table-column label="类型" width="80"><template #default="{ row }"><el-tag size="small">{{ POST_TYPES[row.type] }}</el-tag></template></el-table-column>
            <el-table-column label="作者" width="110"><template #default="{ row }">{{ row.author?.nickname }}</template></el-table-column>
            <el-table-column label="互动" width="150"><template #default="{ row }">{{ row.answerCount }} 答 · {{ row.likes }} 赞 · {{ row.views }} 阅</template></el-table-column>
            <el-table-column label="知识资产" width="100"><template #default="{ row }"><el-tag v-if="row.resourceId" size="small" type="success">已沉淀</el-tag></template></el-table-column>
            <el-table-column label="发布时间" width="170"><template #default="{ row }">{{ dateTime(row.createdAt) }}</template></el-table-column>
            <el-table-column label="操作" width="90" align="right">
              <template #default="{ row }">
                <el-popconfirm title="确定删除该内容及其回答？" @confirm="removePost(row.id)">
                  <template #reference><el-button link type="danger">删除</el-button></template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="知识资产" name="resources">
          <div class="toolbar bar">
            <el-select v-model="resType" placeholder="全部类型" clearable style="width: 130px" @change="load">
              <el-option v-for="(label, k) in RESOURCE_TYPES" :key="k" :label="label" :value="k" />
            </el-select>
            <el-input v-model="keyword" placeholder="搜索资产" clearable style="width: 220px" @keyup.enter="load" @clear="load" />
          </div>
          <el-table :data="resources.list" v-loading="loading">
            <el-table-column label="资产" min-width="300">
              <template #default="{ row }"><strong>{{ row.title }}</strong><div class="muted small">{{ row.summary }}</div></template>
            </el-table-column>
            <el-table-column label="类型" width="90"><template #default="{ row }"><el-tag size="small" effect="plain">{{ RESOURCE_TYPES[row.type] }}</el-tag></template></el-table-column>
            <el-table-column label="来源" width="100"><template #default="{ row }">{{ row.sourceType === 'community' ? '社区沉淀' : '创作者' }}</template></el-table-column>
            <el-table-column label="作者" width="110"><template #default="{ row }">{{ row.author?.nickname }}</template></el-table-column>
            <el-table-column label="浏览 / 点赞" width="110"><template #default="{ row }">{{ row.views }} / {{ row.likes }}</template></el-table-column>
            <el-table-column label="操作" width="90" align="right">
              <template #default="{ row }">
                <el-popconfirm title="确定删除该知识资产？" @confirm="removeResource(row.id)">
                  <template #reference><el-button link type="danger">删除</el-button></template>
                </el-popconfirm>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import PageHeader from '../../components/common/PageHeader.vue'
import { ecoAdminApi } from '../../api'
import { POST_TYPES, RESOURCE_TYPES, dateTime } from '../../utils/eco'

const tab = ref('posts')
const posts = ref({ list: [] })
const resources = ref({ list: [] })
const postType = ref('')
const resType = ref('')
const keyword = ref('')
const loading = ref(false)

async function load() {
  loading.value = true
  try {
    if (tab.value === 'posts') posts.value = await ecoAdminApi.posts({ type: postType.value, keyword: keyword.value, pageSize: 100 })
    else resources.value = await ecoAdminApi.resources({ type: resType.value, keyword: keyword.value, pageSize: 100 })
  } finally {
    loading.value = false
  }
}

async function removePost(id) {
  await ecoAdminApi.removePost(id)
  ElMessage.success('已删除')
  await load()
}

async function removeResource(id) {
  await ecoAdminApi.removeResource(id)
  ElMessage.success('已删除')
  await load()
}

onMounted(load)
</script>

<style scoped>
.page { display: flex; flex-direction: column; gap: 18px; }
.tabs { padding: 6px 22px 22px; }
.bar { margin-bottom: 14px; }
.small { font-size: 12.5px; }
</style>
