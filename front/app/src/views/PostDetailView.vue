<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { ecoApi } from '../api'
import { useAuthStore } from '../stores/auth'
import { POST_TYPES, ROLE_NAMES, timeAgo } from '../utils/eco'
import MarkdownRenderer from '../components/MarkdownRenderer.vue'
import UserAvatar from '../components/UserAvatar.vue'
import Icon from '../components/Icon.vue'

const route = useRoute()
const auth = useAuthStore()
const d = ref(null)
const answer = ref('')
const summarizing = ref(false)
const error = ref('')

const canPost = computed(() => auth.hasPermission('community:post'))
const isAuthor = computed(() => d.value?.post.authorId === auth.user?.id)

async function load() {
  d.value = await ecoApi.post(route.params.id)
}

async function submitAnswer() {
  error.value = ''
  try {
    await ecoApi.answer(d.value.post.id, answer.value)
    answer.value = ''
    await load()
  } catch (e) {
    error.value = e.message
  }
}

async function likePost() {
  await ecoApi.likePost(d.value.post.id)
  d.value.post.likes++
}

async function likeAnswer(a) {
  const res = await ecoApi.likeAnswer(a.id)
  a.likes = res.likes
}

async function accept(a) {
  await ecoApi.acceptAnswer(a.id)
  await load()
}

async function summarize() {
  summarizing.value = true
  error.value = ''
  try {
    d.value = await ecoApi.summarize(d.value.post.id)
  } catch (e) {
    error.value = e.message
  } finally {
    summarizing.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="page pdx">
    <router-link to="/community" class="back-link"><Icon name="arrowLeft" :size="14" /> 知识社区</router-link>
    <template v-if="d">
      <div class="pdx__body">
        <div class="stack">
          <article class="card panel">
            <div class="row row--wrap meta">
              <UserAvatar :user="d.post.author" />
              <b>{{ d.post.author?.nickname }}</b>
              <span class="chip">{{ ROLE_NAMES[d.post.author?.role] || '成员' }}</span>
              <small class="muted">{{ timeAgo(d.post.createdAt) }} · {{ d.post.views }} 浏览</small>
              <span class="tag tag--info">{{ POST_TYPES[d.post.type]?.label }}</span>
            </div>
            <h1>{{ d.post.title }}</h1>
            <MarkdownRenderer :source="d.post.content" />
            <div class="row row--wrap foot">
              <span v-for="t in d.post.tags" :key="t" class="chip">#{{ t }}</span>
              <button v-if="canPost" class="btn btn--sm btn--ghost foot__like" @click="likePost"><Icon name="heart" :size="13" /> {{ d.post.likes }}</button>
            </div>
          </article>

          <section v-if="d.post.aiSummary" class="card panel summary">
            <div class="panel-head">
              <div><h2><span class="tag tag--ai">AI</span> 知识总结</h2><p>Knowledge Agent 将问答整理为结构化知识资产</p></div>
              <router-link :to="`/resources?id=${d.post.resourceId}`" class="btn btn--sm btn--soft">查看知识资产</router-link>
            </div>
            <MarkdownRenderer :source="d.post.aiSummary" />
          </section>

          <section class="card panel">
            <div class="panel-head">
              <div><h2>{{ d.answers.length }} 个回答</h2><p>采纳的回答与高赞回答会优先进入 AI 总结</p></div>
              <button v-if="canPost && d.answers.length" class="btn btn--sm btn--primary" :disabled="summarizing" @click="summarize">
                <Icon name="sparkles" :size="13" /> {{ summarizing ? '总结中…' : d.post.aiSummary ? '重新总结' : 'AI 总结为知识资产' }}
              </button>
            </div>
            <div v-for="a in d.answers" :key="a.id" class="answer" :class="{ 'answer--accepted': a.accepted }">
              <div class="row meta">
                <UserAvatar :user="a.author" />
                <b>{{ a.author?.nickname }}</b>
                <small class="muted">{{ timeAgo(a.createdAt) }}</small>
                <span v-if="a.accepted" class="tag tag--success"><Icon name="check" :size="11" /> 已采纳</span>
              </div>
              <MarkdownRenderer :source="a.content" />
              <div class="row">
                <button v-if="canPost" class="btn btn--sm btn--ghost" @click="likeAnswer(a)"><Icon name="heart" :size="12" /> {{ a.likes }}</button>
                <button v-if="isAuthor && !a.accepted && d.post.type === 'question'" class="btn btn--sm btn--soft" @click="accept(a)">采纳</button>
              </div>
            </div>
            <form v-if="canPost" class="stack reply" @submit.prevent="submitAnswer">
              <textarea v-model="answer" class="textarea" rows="4" placeholder="分享你的经验或解答（支持 Markdown）"></textarea>
              <p v-if="error" class="error-text">{{ error }}</p>
              <div class="row row--between">
                <span class="muted small">社区贡献会计入你的人才画像</span>
                <button class="btn btn--primary" :disabled="!answer.trim()">发布回答</button>
              </div>
            </form>
          </section>
        </div>

        <aside class="stack">
          <section class="card panel">
            <h3 class="side-title">关联知识点</h3>
            <p v-if="!d.knowledge.length" class="muted small">未关联知识点</p>
            <router-link v-for="k in d.knowledge" :key="k.id" :to="`/knowledge/${k.id}`" class="list-row">
              <Icon name="brain" :size="15" />
              <div class="list-row__body"><strong>{{ k.name }}</strong><small>{{ k.domain }}</small></div>
            </router-link>
          </section>
          <section class="card panel">
            <h3 class="side-title">知识生产链路</h3>
            <ol class="flow">
              <li class="done">提出问题</li>
              <li :class="{ done: d.answers.length }">社区回答</li>
              <li :class="{ done: d.post.aiSummary }">AI 总结</li>
              <li :class="{ done: d.post.resourceId }">知识资产</li>
            </ol>
          </section>
        </aside>
      </div>
    </template>
  </div>
</template>

<style scoped>
.pdx .back-link { margin-bottom: 14px; }
.pdx__body { display: grid; grid-template-columns: minmax(0, 1fr) 280px; gap: 18px; align-items: start; }
.pdx h1 { margin: 14px 0 6px; font-size: 24px; }
.meta { font-size: 13px; }
.foot { margin-top: 14px; }
.foot__like { margin-left: auto; }
.summary { border-color: var(--primary-soft-2); background: linear-gradient(160deg, #f6f3ff, #fff 50%); }
.answer { display: flex; flex-direction: column; gap: 8px; padding: 16px 0; border-bottom: 1px dashed var(--border); }
.answer--accepted { margin: 0 -12px; padding: 16px 12px; border-radius: 14px; background: var(--success-soft); }
.answer :deep(.markdown-body) { font-size: 14.5px; }
.reply { margin-top: 16px; gap: 10px; }
.side-title { margin: 0 0 10px; font-size: 14px; }
.list-row { margin-bottom: 8px; }
.flow { margin: 0; padding-left: 18px; color: var(--text-3); font-size: 13.5px; line-height: 2; }
.flow li.done { color: var(--success-strong); font-weight: 600; }
@media (max-width: 960px) { .pdx__body { grid-template-columns: 1fr; } }
</style>
