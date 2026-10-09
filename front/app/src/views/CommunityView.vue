<script setup>
import { onMounted, ref, watch } from 'vue'
import { useRouter } from 'vue-router'
import { ecoApi } from '../api'
import { useAuthStore } from '../stores/auth'
import { POST_TYPES, ROLE_NAMES, timeAgo } from '../utils/eco'
import UserAvatar from '../components/UserAvatar.vue'
import Icon from '../components/Icon.vue'

const router = useRouter()
const auth = useAuthStore()
const type = ref('')
const keyword = ref('')
const page = ref({ list: [], total: 0 })
const showForm = ref(false)
const form = ref({ type: 'question', title: '', content: '', tags: '' })
const error = ref('')

async function load() {
  page.value = await ecoApi.posts({ type: type.value, keyword: keyword.value, pageSize: 50 })
}

async function create() {
  error.value = ''
  try {
    const p = await ecoApi.createPost({ ...form.value, tags: form.value.tags.split(/[,，\s]+/).filter(Boolean) })
    showForm.value = false
    router.push(`/community/${p.id}`)
  } catch (e) {
    error.value = e.message
  }
}

watch(type, load)
onMounted(load)
</script>

<template>
  <div class="page community">
    <header class="page-header page-header--split">
      <div>
        <span class="eyebrow">Knowledge Community</span>
        <h1>知识社区</h1>
        <p>平台的核心不是课程，而是知识社区。问题 → 社区回答 → AI 总结 → 知识资产：每一次讨论都会沉淀为可复用的知识。</p>
      </div>
      <button v-if="auth.hasPermission('community:post')" class="btn btn--primary" @click="showForm = true"><Icon name="plus" :size="15" /> 发布内容</button>
    </header>

    <div class="row row--between row--wrap">
      <div class="segmented">
        <button :class="{ active: type === '' }" @click="type = ''">全部</button>
        <button v-for="(t, k) in POST_TYPES" :key="k" :class="{ active: type === k }" @click="type = k">{{ t.label }}</button>
      </div>
      <form class="search" @submit.prevent="load">
        <Icon name="search" :size="15" />
        <input v-model="keyword" class="input" placeholder="搜索社区内容" />
      </form>
    </div>

    <div class="stack">
      <div v-if="!page.list.length" class="card empty-state"><h3>暂无内容</h3><p>来发布第一个问题吧</p></div>
      <router-link v-for="p in page.list" :key="p.id" :to="`/community/${p.id}`" class="card card--hover post">
        <div class="post__meta">
          <UserAvatar :user="p.author" />
          <b>{{ p.author?.nickname }}</b>
          <span class="chip">{{ ROLE_NAMES[p.author?.role] || '成员' }}</span>
          <small>{{ timeAgo(p.createdAt) }}</small>
          <span class="tag" :class="p.type === 'question' ? 'tag--info' : 'tag--neutral'"><Icon :name="POST_TYPES[p.type]?.icon" :size="11" /> {{ POST_TYPES[p.type]?.label }}</span>
          <span v-if="p.resourceId" class="tag tag--success"><Icon name="sparkles" :size="11" /> 已沉淀知识资产</span>
        </div>
        <h3>{{ p.title }}</h3>
        <p>{{ p.content }}</p>
        <div class="post__foot">
          <span v-for="t in p.tags" :key="t" class="post__tag">#{{ t }}</span>
          <small><Icon name="heart" :size="12" /> {{ p.likes }}</small>
          <small><Icon name="message" :size="12" /> {{ p.answerCount }}</small>
          <small><Icon name="eye" :size="12" /> {{ p.views }}</small>
        </div>
      </router-link>
    </div>

    <div v-if="showForm" class="modal-mask" @click.self="showForm = false">
      <form class="modal" @submit.prevent="create">
        <h2>发布到知识社区</h2>
        <div class="form-grid">
          <label class="field field--full"><span>类型</span>
            <div class="segmented">
              <button v-for="(t, k) in POST_TYPES" :key="k" type="button" :class="{ active: form.type === k }" @click="form.type = k">{{ t.label }}</button>
            </div>
          </label>
          <label class="field field--full"><span>标题</span><input v-model="form.title" class="input" placeholder="一句话说清楚你的问题或主题" /></label>
          <label class="field field--full"><span>正文（支持 Markdown）</span><textarea v-model="form.content" class="textarea" rows="7"></textarea></label>
          <label class="field field--full"><span>标签（逗号分隔）</span><input v-model="form.tags" class="input" placeholder="Java并发, 面试" /></label>
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
.community { display: flex; flex-direction: column; gap: 16px; }
.community .page-header { margin-bottom: 4px; }
.search { position: relative; width: 280px; }
.search .icon { position: absolute; left: 12px; top: 50%; color: var(--text-3); transform: translateY(-50%); }
.search .input { padding-left: 36px; }
.post { display: flex; flex-direction: column; gap: 8px; padding: 20px 22px; color: inherit; }
.post__meta { display: flex; flex-wrap: wrap; align-items: center; gap: 8px; font-size: 13px; }
.post__meta small { color: var(--text-3); }
.post h3 { margin: 4px 0 0; font-size: 17px; }
.post p { margin: 0; color: var(--text-2); font-size: 13.5px; }
.post__foot { display: flex; flex-wrap: wrap; align-items: center; gap: 10px; color: var(--text-3); font-size: 12px; }
.post__tag { color: var(--primary); font-weight: 600; }
.post__foot small { display: inline-flex; align-items: center; gap: 3px; }
.post__foot small:first-of-type { margin-left: auto; }
@media (max-width: 560px) { .search { width: 100%; } }
</style>
