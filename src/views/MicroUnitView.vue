<script setup>
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getMicroUnit, microUnits } from '../data/careerPath'
import { useGrowthStore } from '../stores/growth'

const route = useRoute()
const router = useRouter()
const growth = useGrowthStore()
const answers = ref({})
const submitted = ref(false)

const unit = computed(() => getMicroUnit(route.params.unitId))
const unitIndex = computed(() => microUnits.findIndex((item) => item.id === unit.value?.id))
const nextUnit = computed(() => microUnits[unitIndex.value + 1])
const answeredCount = computed(() => Object.keys(answers.value).length)
const score = computed(() => {
  if (!unit.value) return 0
  const correct = unit.value.questions.filter((question, index) => answers.value[index] === question.answer).length
  return Math.round(correct / unit.value.questions.length * 100)
})

watch(() => route.params.unitId, () => {
  answers.value = {}
  submitted.value = false
})

function submit() {
  if (answeredCount.value !== unit.value.questions.length) return
  submitted.value = true
  growth.completeUnit(unit.value.id, score.value)
}

function goNext() {
  if (nextUnit.value) router.push(`/micro/${nextUnit.value.id}`)
  else router.push('/path')
}
</script>

<template>
  <div v-if="unit" class="page micro-page">
    <header class="micro-header">
      <router-link to="/path">← 返回学习路径</router-link>
      <div class="micro-meta">
        <span class="tag">{{ unit.competency }}</span>
        <span>⏱ {{ unit.duration }} 分钟</span>
        <span>{{ unit.difficulty }}</span>
      </div>
      <h1>{{ unit.title }}</h1>
      <div class="unit-progress"><div :style="{ width: ((unitIndex + 1) / microUnits.length * 100) + '%' }"></div></div>
      <small>路径进度 {{ unitIndex + 1 }} / {{ microUnits.length }}</small>
    </header>

    <main class="lesson-card card">
      <section class="intro-block">
        <span>01 · 引入</span>
        <p>{{ unit.intro }}</p>
      </section>

      <section class="content-section">
        <div class="section-title"><span>02</span><h2>核心讲解</h2><small>约 2 分钟</small></div>
        <article v-for="block in unit.content" :key="block.title">
          <h3>{{ block.title }}</h3>
          <p>{{ block.text }}</p>
        </article>
        <pre><code>{{ unit.example }}</code></pre>
      </section>

      <section class="quiz-section">
        <div class="section-title"><span>03</span><h2>快速练习</h2><small>即学即测</small></div>
        <article v-for="(question, qIndex) in unit.questions" :key="question.text" class="question">
          <h3>{{ qIndex + 1 }}. {{ question.text }}</h3>
          <button
            v-for="(option, oIndex) in question.options"
            :key="option"
            class="option"
            :class="{
              selected: answers[qIndex] === oIndex,
              correct: submitted && oIndex === question.answer,
              wrong: submitted && answers[qIndex] === oIndex && oIndex !== question.answer
            }"
            :disabled="submitted"
            @click="answers[qIndex] = oIndex"
          >
            <span>{{ String.fromCharCode(65 + oIndex) }}</span>{{ option }}
          </button>
          <p v-if="submitted" class="explanation">{{ question.explanation }}</p>
        </article>
        <button v-if="!submitted" class="btn btn--primary submit-btn" :disabled="answeredCount !== unit.questions.length" @click="submit">
          提交快测（{{ answeredCount }}/{{ unit.questions.length }}）
        </button>
        <div v-else class="result" :class="{ passed: score >= 75 }">
          <strong>{{ score }} 分</strong>
          <span>{{ score >= 75 ? '掌握得不错！达成度与积分已更新。' : '建议复习讲解后再测一次。' }}</span>
        </div>
      </section>

      <section v-if="submitted" class="summary-section">
        <div class="section-title"><span>04</span><h2>要点回顾</h2></div>
        <div class="summary-grid">
          <div v-for="(item, index) in unit.summary" :key="item"><b>{{ index + 1 }}</b>{{ item }}</div>
        </div>
        <button class="btn btn--primary" @click="goNext">
          {{ nextUnit ? `下一步：${nextUnit.title}` : '查看完整学习路径' }} →
        </button>
      </section>
    </main>
  </div>
  <div v-else class="page"><div class="card not-found">微单元不存在。<router-link to="/path">返回学习路径</router-link></div></div>
</template>

<style scoped>
.micro-page { max-width: 850px; }
.micro-header { margin-bottom: 20px; }
.micro-header > a { font-size: 13px; color: var(--text-2); }
.micro-header h1 { margin: 12px 0 10px; font-size: 28px; }
.micro-meta { display: flex; align-items: center; gap: 12px; margin-top: 18px; color: var(--text-3); font-size: 12px; }
.unit-progress { height: 5px; background: var(--border); border-radius: 8px; overflow: hidden; }
.unit-progress div { height: 100%; background: linear-gradient(90deg, var(--primary), #8b5cf6); }
.micro-header small { display: block; margin-top: 5px; color: var(--text-3); text-align: right; }
.lesson-card { overflow: hidden; }
.intro-block { padding: 28px 34px; background: linear-gradient(110deg, #312e81, #6366f1); color: white; }
.intro-block span { color: #c7d2fe; font-size: 12px; font-weight: 700; letter-spacing: .08em; }
.intro-block p { margin: 8px 0 0; font-size: 18px; line-height: 1.8; }
.content-section, .quiz-section, .summary-section { padding: 28px 34px; border-top: 1px solid var(--border); }
.section-title { display: flex; align-items: center; gap: 10px; margin-bottom: 16px; }
.section-title > span { display: grid; place-items: center; width: 30px; height: 30px; border-radius: 8px; background: var(--primary-soft); color: var(--primary); font-size: 12px; font-weight: 800; }
.section-title h2 { margin: 0; font-size: 18px; }
.section-title small { margin-left: auto; color: var(--text-3); }
.content-section article { margin: 20px 0; }
.content-section h3, .question h3 { margin: 0 0 6px; font-size: 15px; }
.content-section p { color: var(--text-2); }
pre { margin: 18px 0 0; padding: 18px; overflow-x: auto; border-radius: 10px; background: #1e293b; color: #e2e8f0; font: 13px/1.7 'SFMono-Regular', Consolas, monospace; }
.question { margin: 22px 0; }
.option { display: flex; align-items: center; gap: 10px; width: 100%; margin: 8px 0; padding: 10px 12px; border: 1px solid var(--border); border-radius: 9px; background: white; color: var(--text); cursor: pointer; text-align: left; }
.option span { display: grid; place-items: center; width: 24px; height: 24px; border-radius: 6px; background: var(--surface-2); font-weight: 700; }
.option:hover, .option.selected { border-color: var(--primary); background: var(--primary-soft); }
.option.correct { border-color: var(--success); background: var(--success-soft); }
.option.wrong { border-color: var(--danger); background: #fef2f2; }
.explanation { padding: 8px 12px; border-left: 3px solid var(--success); background: var(--success-soft); color: #047857; font-size: 13px; }
.submit-btn { width: 100%; }
.result { display: flex; align-items: center; gap: 14px; padding: 16px; border-radius: 10px; background: #fff7ed; color: #c2410c; }
.result.passed { background: var(--success-soft); color: #047857; }
.result strong { font-size: 24px; }
.summary-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 10px; margin-bottom: 20px; }
.summary-grid div { display: flex; gap: 8px; padding: 12px; border-radius: 9px; background: var(--surface-2); font-size: 13px; }
.summary-grid b { color: var(--primary); }
.summary-section .btn { width: 100%; }
.not-found { margin-top: 50px; padding: 40px; text-align: center; }
@media (max-width: 600px) {
  .content-section, .quiz-section, .summary-section, .intro-block { padding: 22px 20px; }
  .summary-grid { grid-template-columns: 1fr; }
}
</style>
