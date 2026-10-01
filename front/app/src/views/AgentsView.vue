<script setup>
import { computed, onMounted, ref } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { aiApi, kernelApi } from '../api'
import AgentChat from '../components/AgentChat.vue'
import Icon from '../components/Icon.vue'

const route = useRoute()
const router = useRouter()
const agents = ref([])
const config = ref(null)

const code = computed(() => route.params.code || agents.value[0]?.code)
const agent = computed(() => agents.value.find((a) => a.code === code.value))

onMounted(async () => {
  agents.value = await kernelApi.agents()
  config.value = await aiApi.config().catch(() => null)
})
</script>

<template>
  <div class="page agents">
    <header class="page-header page-header--split">
      <div>
        <span class="eyebrow">AI Agents</span>
        <h1>AI 伙伴</h1>
        <p>AI 不是聊天机器人，而是平台的学习内核。六个 Agent 各司其职，回答时会读取你的职业目标、技能状态、学习记录与项目证据。</p>
      </div>
      <router-link to="/kernel" class="btn btn--soft"><Icon name="layers" :size="15" /> AI 学习内核</router-link>
    </header>

    <div v-if="config && !config.enabled" class="notice">
      <Icon name="alert" :size="16" />
      大模型尚未配置，Agent 对话暂不可用；你仍可点击「上下文」查看 Agent 会使用的个人数据。管理员可在管理端「大模型配置」中接入。
    </div>

    <div class="agents__body">
      <nav class="agents__list">
        <button
          v-for="a in agents"
          :key="a.code"
          class="card agent-card"
          :class="{ active: a.code === code }"
          @click="router.replace(`/agents/${a.code}`)"
        >
          <span class="agent-card__icon" :style="{ background: a.color }"><Icon :name="a.icon" :size="18" /></span>
          <span class="agent-card__body">
            <strong>{{ a.name }}</strong>
            <small>{{ a.title }} · {{ a.responsibilities.join(' / ') }}</small>
          </span>
        </button>
      </nav>
      <section v-if="agent" class="card panel agents__chat">
        <AgentChat :key="agent.code" :agent="agent" />
      </section>
    </div>
  </div>
</template>

<style scoped>
.agents { display: flex; flex-direction: column; gap: 16px; }
.agents .page-header { margin-bottom: 4px; }
.notice { display: flex; align-items: center; gap: 8px; padding: 12px 16px; border-radius: 14px; background: var(--warning-soft); color: var(--warning-strong); font-size: 13.5px; }
.agents__body { display: grid; grid-template-columns: 290px minmax(0, 1fr); gap: 18px; align-items: start; }
.agents__list { display: flex; flex-direction: column; gap: 10px; }
.agent-card { display: flex; align-items: center; gap: 12px; padding: 14px; font: inherit; text-align: left; cursor: pointer; transition: all var(--t-fast); }
.agent-card:hover { border-color: var(--border-strong); }
.agent-card.active { border-color: var(--primary); box-shadow: 0 0 0 3px var(--primary-soft-2); }
.agent-card__icon { display: grid; place-items: center; width: 38px; height: 38px; flex: 0 0 38px; border-radius: 12px; color: #fff; }
.agent-card__body { display: flex; flex-direction: column; line-height: 1.4; }
.agent-card__body strong { font-size: 14px; }
.agent-card__body small { color: var(--text-3); font-size: 11.5px; }
.agents__chat { height: calc(100vh - 230px); min-height: 520px; display: flex; flex-direction: column; }
@media (max-width: 960px) {
  .agents__body { grid-template-columns: 1fr; }
  .agents__list { flex-direction: row; overflow-x: auto; }
  .agent-card { min-width: 220px; }
}
</style>
