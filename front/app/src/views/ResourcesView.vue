<script setup>
import { onMounted, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ecoApi } from '../api'
import { useAuthStore } from '../stores/auth'
import { RESOURCE_TYPES, timeAgo } from '../utils/eco'
import MarkdownRenderer from '../components/MarkdownRenderer.vue'
import UserAvatar from '../components/UserAvatar.vue'
import Icon from '../components/Icon.vue'

const route = useRoute()
const router = useRouter()
const auth = useAuthStore()
const type = ref('')
const keyword = ref('')
const page = ref({ list: [], total: 0 })
const current = ref(null)
const showForm = ref(false)
const form = ref({ type: 'article', title: '', summary: '', content: '', url: '', tags: '' })
const error = ref('')

async function load() {
  page.value = await ecoApi.resources({ type: type.value, keyword: keyword.value, pageSize: 60 })
}

async function open(id) {
  current.value = await ecoApi.resource(id)
  router.replace({ query: { ...route.query, id } })
}

function close() {
  current.value = null
  const { id, ...rest } = route.query
  router.replace({ query: rest })
}

function isInternal(url) {
  return url?.startsWith('#/')
}

async function create() {
  error.value = ''
  try {
    const r = await ecoApi.createResource({ ...form.value, tags: form.value.tags.split(/[,，\s]+/).filter(Boolean) })
    showForm.value = false
    form.value = { type: 'article', title: '', summary: '', content: '', url: '', tags: '' }
    await load()
    open(r.id)
  } catch (e) {
    error.value = e.message
  }
}

watch(type, load)
onMounted(async () => {
  await load()
  if (route.query.id) open(route.query.id)
})
</script>

<template>
  <div class="page res">
    <header class="page-header page-header--split">
      <div>
        <span class="eyebrow">Knowledge Resource</span>
        <h1>知识资产</h1>
        <p>统一的知识资源：文章、课程、视频、代码、案例、项目、任务、Prompt 与 SOP。来自创作者发布与社区问答的 AI 沉淀。</p>
      </div>
      <button v-if="auth.hasPermission('resource:publish')" class="btn btn--primary" @click="showForm = true"><Icon name="plus" :size="15" /> 发布资产</button>
    </header>

    <div class="row row--between row--wrap">
      <div class="row row--wrap">
        <button class="chip" :class="{ active: type === '' }" @click="type = ''">全部</button>
        <button v-for="(t, k) in RESOURCE_TYPES" :key="k" class="chip" :class="{ active: type === k }" @click="type = k"><Icon :name="t.icon" :size="11" /> {{ t.label }}</button>
      </div>
      <form class="row" @submit.prevent="load"><input v-model="keyword" class="input" placeholder="搜索知识资产" /></form>
    </div>

    <div class="grid-3">
      <button v-for="r in page.list" :key="r.id" class="card card--hover resource" @click="open(r.id)">
        <div class="row">
          <span class="resource__icon"><Icon :name="RESOURCE_TYPES[r.type]?.icon || 'note'" :size="16" /></span>
          <span class="tag tag--neutral">{{ RESOURCE_TYPES[r.type]?.label }}</span>
          <span v-if="r.sourceType === 'community'" class="tag tag--success">社区沉淀</span>
        </div>
        <h3>{{ r.title }}</h3>
        <p>{{ r.summary }}</p>
        <div class="row small muted resource__foot">
          <UserAvatar :user="r.author" />
          <span>{{ r.author?.nickname }}</span>
          <span class="resource__stats"><Icon name="eye" :size="12" /> {{ r.views }} · <Icon name="heart" :size="12" /> {{ r.likes }}</span>
        </div>
      </button>
    </div>
    <div v-if="!page.list.length" class="card empty-state"><h3>暂无资源</h3></div>

    <div v-if="current" class="modal-mask" @click.self="close">
      <div class="modal modal--wide">
        <div class="row row--between">
          <span class="tag tag--neutral">{{ RESOURCE_TYPES[current.type]?.label }}</span>
          <button class="btn btn--sm btn--ghost" @click="close"><Icon name="x" :size="14" /></button>
        </div>
        <h2>{{ current.title }}</h2>
        <p class="muted small">{{ current.author?.nickname }} · {{ timeAgo(current.createdAt) }} · {{ current.views }} 浏览</p>
        <p>{{ current.summary }}</p>
        <MarkdownRenderer v-if="current.content" :source="current.content" />
        <div v-if="current.url" class="modal__foot">
          <router-link v-if="isInternal(current.url)" :to="current.url.slice(1)" class="btn btn--primary" @click="close">前往学习</router-link>
          <a v-else :href="current.url" target="_blank" rel="noopener" class="btn btn--primary">打开链接</a>
        </div>
      </div>
    </div>

    <div v-if="showForm" class="modal-mask" @click.self="showForm = false">
      <form class="modal" @submit.prevent="create">
        <h2>发布知识资产</h2>
        <div class="form-grid">
          <label class="field"><span>类型</span>
            <select v-model="form.type" class="select"><option v-for="(t, k) in RESOURCE_TYPES" :key="k" :value="k">{{ t.label }}</option></select>
          </label>
          <label class="field"><span>标签</span><input v-model="form.tags" class="input" placeholder="逗号分隔" /></label>
          <label class="field field--full"><span>标题</span><input v-model="form.title" class="input" /></label>
          <label class="field field--full"><span>摘要</span><input v-model="form.summary" class="input" /></label>
          <label class="field field--full"><span>正文（Markdown）</span><textarea v-model="form.content" class="textarea" rows="6"></textarea></label>
          <label class="field field--full"><span>链接（视频/外部资源可选）</span><input v-model="form.url" class="input" placeholder="https://" /></label>
        </div>
        <p v-if="error" class="error-text">{{ error }}</p>
        <div class="modal__foot">
          <button type="button" class="btn btn--ghost" @click="showForm = false">取消</button>
          <button type="submit" class="btn btn--primary">发布</button>
        </div>
      </form>
    </div>
  </div>
</template>

<style scoped>
.res { display: flex; flex-direction: column; gap: 16px; }
.res .page-header { margin-bottom: 4px; }
.resource { display: flex; flex-direction: column; gap: 8px; padding: 18px; font: inherit; text-align: left; cursor: pointer; color: inherit; }
.resource__icon { display: grid; place-items: center; width: 32px; height: 32px; border-radius: 10px; background: var(--primary-soft); color: var(--primary); }
.resource h3 { margin: 4px 0 0; font-size: 15.5px; }
.resource p { flex: 1; margin: 0; color: var(--text-2); font-size: 13px; }
.resource__foot .avatar { width: 22px; height: 22px; border-radius: 7px; font-size: 10px; }
.resource__stats { margin-left: auto; display: inline-flex; align-items: center; gap: 3px; }
.modal--wide { width: min(760px, 100%); }
</style>
