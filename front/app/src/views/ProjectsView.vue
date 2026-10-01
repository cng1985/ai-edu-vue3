<script setup>
import { onMounted, ref } from 'vue'
import { ecoApi } from '../api'
import Icon from '../components/Icon.vue'

const projects = ref([])
const loading = ref(true)

onMounted(async () => {
  try {
    projects.value = await ecoApi.projects()
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="page">
    <header class="page-header">
      <span class="eyebrow">Project · Task · Execution · Evidence</span>
      <h1>项目实践</h1>
      <p>项目用于将知识转化为能力。每个任务都标注了技能要求，提交后由 Review Agent 评审打分，生成可验证的能力证据，直接提升你的技能等级。</p>
    </header>

    <div v-if="loading" class="muted">加载中…</div>
    <div class="grid-3 stagger">
      <router-link v-for="p in projects" :key="p.id" :to="`/projects/${p.id}`" class="card card--hover project">
        <div class="project__top" :style="{ background: `linear-gradient(135deg, ${p.color}, ${p.color}cc)` }">
          <Icon :name="p.icon || 'layers'" :size="22" />
          <span class="tag tag--glass">{{ p.domain }}</span>
        </div>
        <div class="project__body">
          <h3>{{ p.title }}</h3>
          <p>{{ p.summary }}</p>
          <div class="row small muted">
            <span><Icon name="layers" :size="12" /> {{ p.tasks.length }} 个任务</span>
            <span>难度 {{ '★'.repeat(p.difficulty) }}</span>
          </div>
          <div class="project__progress">
            <div class="row row--between small"><span>完成 {{ p.completed }}/{{ p.tasks.length }}</span><b class="num">{{ p.progress }}%</b></div>
            <div class="progress progress--thin"><i :style="{ width: p.progress + '%', background: p.color }"></i></div>
          </div>
          <div class="project__tasks">
            <span v-for="t in p.tasks" :key="t.id" class="chip" :class="{ 'chip--ok': t.bestScore >= 60 }">
              <Icon v-if="t.bestScore >= 60" name="check" :size="11" /> {{ t.title }}
            </span>
          </div>
        </div>
      </router-link>
    </div>
  </div>
</template>

<style scoped>
.project { display: flex; flex-direction: column; overflow: hidden; color: inherit; }
.project__top { display: flex; align-items: center; justify-content: space-between; padding: 18px 20px; color: #fff; }
.project__body { display: flex; flex: 1; flex-direction: column; gap: 10px; padding: 18px 20px 20px; }
.project__body h3 { margin: 0; font-size: 17px; }
.project__body p { margin: 0; color: var(--text-2); font-size: 13.5px; }
.project__progress { margin-top: 4px; }
.project__progress .progress { margin-top: 5px; }
.project__tasks { display: flex; flex-wrap: wrap; gap: 6px; }
</style>
