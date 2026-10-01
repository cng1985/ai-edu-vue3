<script setup>
import { nextTick, onBeforeUnmount, ref, watch } from 'vue'
import { kernelApi } from '../api'
import MarkdownRenderer from './MarkdownRenderer.vue'
import Icon from './Icon.vue'

const props = defineProps({
  agent: { type: Object, required: true },
  compact: { type: Boolean, default: false },
  initialMessage: { type: String, default: '' }
})

const messages = ref([])
const input = ref(props.initialMessage)
const streaming = ref(false)
const error = ref('')
const scroller = ref(null)
const contextText = ref('')
const showContext = ref(false)
let cancel = null

watch(
  () => props.agent.code,
  () => {
    cancel?.()
    messages.value = []
    error.value = ''
    contextText.value = ''
    showContext.value = false
  }
)

async function scrollDown() {
  await nextTick()
  if (scroller.value) scroller.value.scrollTop = scroller.value.scrollHeight
}

function send(text) {
  const content = (text ?? input.value).trim()
  if (!content || streaming.value) return
  error.value = ''
  const history = messages.value.map(({ role, content }) => ({ role, content }))
  messages.value.push({ role: 'user', content })
  const reply = { role: 'assistant', content: '', sources: [] }
  messages.value.push(reply)
  input.value = ''
  streaming.value = true
  scrollDown()
  cancel = kernelApi.agentChat(props.agent.code, content, history, {
    token: ({ content }) => {
      reply.content += content
      messages.value = [...messages.value]
      scrollDown()
    },
    done: ({ sources }) => {
      reply.sources = sources || []
      streaming.value = false
      messages.value = [...messages.value]
    },
    error: (e) => {
      streaming.value = false
      error.value = e.message
      if (!reply.content) messages.value = messages.value.filter((m) => m !== reply)
    }
  })
}

function stop() {
  cancel?.()
  streaming.value = false
}

async function toggleContext() {
  showContext.value = !showContext.value
  if (showContext.value && !contextText.value) {
    const last = [...messages.value].reverse().find((m) => m.role === 'user')
    try {
      const res = await kernelApi.agentContext(props.agent.code, last?.content || '')
      contextText.value = res.context
    } catch (e) {
      contextText.value = e.message
    }
  }
}

onBeforeUnmount(() => cancel?.())
defineExpose({ send })
</script>

<template>
  <div class="agent-chat" :class="{ 'agent-chat--compact': compact }">
    <div ref="scroller" class="agent-chat__log">
      <div v-if="!messages.length" class="agent-chat__intro">
        <span class="agent-chat__badge" :style="{ background: agent.color }"><Icon :name="agent.icon" :size="20" /></span>
        <strong>{{ agent.name }} · {{ agent.title }}</strong>
        <p>{{ agent.description }}</p>
        <div class="agent-chat__starters">
          <button v-for="s in agent.starters" :key="s" class="chip" @click="send(s)">{{ s }}</button>
        </div>
      </div>
      <div v-for="(m, i) in messages" :key="i" class="msg" :class="`msg--${m.role}`">
        <span v-if="m.role === 'assistant'" class="msg__avatar" :style="{ background: agent.color }"><Icon :name="agent.icon" :size="14" /></span>
        <div class="msg__bubble">
          <MarkdownRenderer v-if="m.role === 'assistant'" :source="m.content || '…'" :live="streaming && i === messages.length - 1" />
          <template v-else>{{ m.content }}</template>
          <div v-if="m.sources?.length" class="msg__sources">
            <router-link
              v-for="s in m.sources"
              :key="s.courseId + s.chapterId"
              :to="`/courses/${s.courseId}/${s.chapterId}`"
              class="chip"
            >
              <Icon name="book" :size="11" /> {{ s.chapterTitle }}
            </router-link>
          </div>
        </div>
      </div>
    </div>

    <p v-if="error" class="agent-chat__error"><Icon name="alert" :size="14" /> {{ error }}</p>

    <div v-if="showContext" class="agent-chat__context">
      <div class="row row--between"><strong>Agent 使用的个人成长上下文</strong><button class="btn btn--sm btn--ghost" @click="showContext = false">收起</button></div>
      <pre>{{ contextText || '加载中…' }}</pre>
    </div>

    <form class="agent-chat__input" @submit.prevent="send()">
      <textarea
        v-model="input"
        class="textarea"
        rows="2"
        :placeholder="`向 ${agent.name} 提问…（Enter 发送，Shift+Enter 换行）`"
        @keydown.enter.exact.prevent="send()"
      ></textarea>
      <div class="agent-chat__actions">
        <button type="button" class="btn btn--sm btn--ghost" title="查看 Agent 依据的数据" @click="toggleContext">
          <Icon name="eye" :size="14" /> 上下文
        </button>
        <button v-if="streaming" type="button" class="btn btn--sm btn--danger" @click="stop"><Icon name="stop" :size="14" /> 停止</button>
        <button v-else type="submit" class="btn btn--sm btn--primary" :disabled="!input.trim()"><Icon name="send" :size="14" /> 发送</button>
      </div>
    </form>
  </div>
</template>

<style scoped>
.agent-chat { display: flex; flex-direction: column; height: 100%; min-height: 0; }
.agent-chat__log { flex: 1; min-height: 240px; overflow-y: auto; padding: 4px 2px 12px; }
.agent-chat__intro { display: flex; flex-direction: column; align-items: center; gap: 8px; padding: 36px 16px; text-align: center; }
.agent-chat__badge { display: grid; place-items: center; width: 48px; height: 48px; border-radius: 16px; color: #fff; }
.agent-chat__intro p { max-width: 420px; margin: 0; color: var(--text-2); font-size: 13.5px; }
.agent-chat__starters { display: flex; flex-wrap: wrap; justify-content: center; gap: 8px; margin-top: 8px; }
.msg { display: flex; gap: 10px; margin: 14px 0; }
.msg--user { justify-content: flex-end; }
.msg__avatar { display: grid; place-items: center; width: 28px; height: 28px; flex: 0 0 28px; border-radius: 9px; color: #fff; }
.msg__bubble { max-width: 82%; padding: 10px 14px; border-radius: 14px; background: var(--surface-2); border: 1px solid var(--border); font-size: 14px; }
.msg--user .msg__bubble { background: linear-gradient(135deg, var(--primary), var(--primary-strong)); border: none; color: #fff; white-space: pre-wrap; }
.msg__bubble :deep(.markdown-body) { font-size: 14px; line-height: 1.75; }
.msg__bubble :deep(.markdown-body p:first-child) { margin-top: 0; }
.msg__bubble :deep(.markdown-body p:last-child) { margin-bottom: 0; }
.msg__sources { display: flex; flex-wrap: wrap; gap: 6px; margin-top: 10px; }
.agent-chat__error { display: flex; align-items: center; gap: 6px; margin: 0 0 10px; padding: 10px 12px; border-radius: 12px; background: var(--danger-soft); color: var(--danger-strong); font-size: 13px; }
.agent-chat__context { margin-bottom: 10px; padding: 12px 14px; border-radius: 12px; background: var(--surface-2); border: 1px dashed var(--border-strong); font-size: 12.5px; }
.agent-chat__context pre { max-height: 220px; overflow: auto; margin: 8px 0 0; white-space: pre-wrap; color: var(--text-2); font-family: inherit; }
.agent-chat__input { display: flex; flex-direction: column; gap: 8px; }
.agent-chat__input .textarea { min-height: 64px; }
.agent-chat__actions { display: flex; justify-content: flex-end; gap: 8px; }
.agent-chat--compact .agent-chat__log { min-height: 180px; max-height: 360px; }
</style>
