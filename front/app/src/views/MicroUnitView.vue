<script setup>
import { computed, ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { getMicroUnit, microUnits } from '../data/careerPath'
import { useGrowthStore } from '../stores/growth'
import Icon from '../components/Icon.vue'

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
      <router-link to="/path" class="back-link"><Icon name="arrowLeft" :size="15" /> 返回学习路径</router-link>
      <div class="micro-meta">
        <span class="tag">{{ unit.competency }}</span>
        <span class="micro-meta__item"><Icon name="clock" :size="13" /> {{ unit.duration }} 分钟</span>
        <span class="micro-meta__item"><Icon name="layers" :size="13" /> {{ unit.difficulty }}</span>
      </div>
      <h1>{{ unit.title }}</h1>
      <div class="unit-progress">
        <div class="progress progress--thin"><div :style="{ width: ((unitIndex + 1) / microUnits.length * 100) + '%' }"></div></div>
        <small>路径进度 <b class="num">{{ unitIndex + 1 }}</b> / {{ microUnits.length }}</small>
      </div>
    </header>

    <main class="lesson-card card">
      <section class="intro-block">
        <span class="intro-block__label">01 · 引入</span>
        <p>{{ unit.intro }}</p>
      </section>

      <section class="block">
        <div class="section-title"><span>02</span><h2>核心讲解</h2><small><Icon name="clock" :size="12" /> 约 2 分钟</small></div>
        <article v-for="block in unit.content" :key="block.title" class="content-item">
          <h3>{{ block.title }}</h3>
          <p>{{ block.text }}</p>
        </article>
        <div class="code-block">
          <div class="code-block__header">示例代码</div>
          <pre><code>{{ unit.example }}</code></pre>
        </div>
      </section>

      <section class="block">
        <div class="section-title"><span>03</span><h2>快速练习</h2><small>即学即测</small></div>
        <article v-for="(question, qIndex) in unit.questions" :key="question.text" class="question">
          <h3><b class="num">{{ qIndex + 1 }}</b>{{ question.text }}</h3>
          <div class="options">
            <button
              v-for="(option, oIndex) in question.options"
              :key="option"
              class="option"
              :class="{
                'option--selected': answers[qIndex] === oIndex,
                'option--correct': submitted && oIndex === question.answer,
                'option--wrong': submitted && answers[qIndex] === oIndex && oIndex !== question.answer
              }"
              :disabled="submitted"
              @click="answers[qIndex] = oIndex"
            >
              <span class="option__letter">{{ String.fromCharCode(65 + oIndex) }}</span>
              <span class="option__text">{{ option }}</span>
              <Icon v-if="submitted && oIndex === question.answer" name="check" :size="16" :stroke="3" class="option__mark" />
              <Icon v-else-if="submitted && answers[qIndex] === oIndex" name="x" :size="16" :stroke="3" class="option__mark option__mark--wrong" />
            </button>
          </div>
          <p v-if="submitted" class="explanation"><Icon name="sparkles" :size="14" /> {{ question.explanation }}</p>
        </article>

        <button v-if="!submitted" class="btn btn--primary btn--lg btn--block" :disabled="answeredCount !== unit.questions.length" @click="submit">
          提交快测（{{ answeredCount }}/{{ unit.questions.length }}）
        </button>
        <div v-else class="result" :class="{ 'result--passed': score >= 75 }">
          <div class="result__score"><strong class="num">{{ score }}</strong><span>分</span></div>
          <div class="result__text">
            <strong>{{ score >= 75 ? '掌握得不错！' : '还差一点点' }}</strong>
            <span>{{ score >= 75 ? '达成度与积分已更新，继续保持节奏。' : '建议复习讲解后再测一次。' }}</span>
          </div>
        </div>
      </section>

      <section v-if="submitted" class="block fade-up">
        <div class="section-title"><span>04</span><h2>要点回顾</h2></div>
        <div class="summary-grid">
          <div v-for="(item, index) in unit.summary" :key="item"><b class="num">{{ index + 1 }}</b>{{ item }}</div>
        </div>
        <button class="btn btn--primary btn--lg btn--block" @click="goNext">
          {{ nextUnit ? `下一步：${nextUnit.title}` : '查看完整学习路径' }} <Icon name="arrowRight" :size="16" />
        </button>
      </section>
    </main>
  </div>
  <div v-else class="page">
    <div class="empty-state card">
      <div class="empty-state__icon"><Icon name="alert" :size="30" /></div>
      <h2>微单元不存在</h2>
      <router-link to="/path" class="btn btn--primary">返回学习路径</router-link>
    </div>
  </div>
</template>

<style scoped>
.micro-page { max-width: 860px; }
.micro-header { margin-bottom: 22px; }
.micro-header h1 { margin: 12px 0 16px; font-size: 30px; }
.micro-meta { display: flex; align-items: center; gap: 12px; margin-top: 6px; color: var(--text-3); font-size: 13px; }
.micro-meta__item { display: inline-flex; align-items: center; gap: 4px; }
.unit-progress { display: flex; align-items: center; gap: 14px; }
.unit-progress .progress { flex: 1; }
.unit-progress small { color: var(--text-3); font-size: 12.5px; white-space: nowrap; }
.unit-progress b { color: var(--text); }

.lesson-card { overflow: hidden; border-radius: var(--radius-lg); }
.intro-block {
  position: relative;
  padding: 32px 38px;
  color: white;
  background:
    radial-gradient(90% 100% at 100% 0%, rgba(160, 113, 255, 0.5), transparent 60%),
    linear-gradient(120deg, #1b1d3e, #3e30c4);
}
.intro-block__label { color: #c7bfff; font-size: 12px; font-weight: 800; letter-spacing: 0.1em; }
.intro-block p { margin: 10px 0 0; font-size: 18px; line-height: 1.8; }
.block { padding: 30px 38px; border-top: 1px solid var(--border); }
.section-title { display: flex; align-items: center; gap: 12px; margin-bottom: 18px; }
.section-title > span { display: grid; place-items: center; width: 32px; height: 32px; border-radius: 10px; background: var(--primary-soft); color: var(--primary-strong); font-family: 'Sora', sans-serif; font-size: 12px; font-weight: 800; }
.section-title h2 { margin: 0; font-size: 19px; }
.section-title small { display: inline-flex; align-items: center; gap: 4px; margin-left: auto; color: var(--text-3); font-size: 12.5px; }
.content-item { margin: 18px 0; padding-left: 16px; border-left: 3px solid var(--primary-soft-2); }
.content-item h3 { margin: 0 0 6px; font-size: 15.5px; }
.content-item p { margin: 0; color: var(--text-2); }
.code-block pre { color: #e2e8f0; font: 13px/1.7 'JetBrains Mono', 'SFMono-Regular', Consolas, monospace; }
.code-block code { font-family: inherit; }

.question { margin: 0 0 26px; }
.question h3 { display: flex; align-items: flex-start; gap: 12px; margin: 0 0 12px; font-size: 15.5px; line-height: 1.6; }
.question h3 b { display: grid; flex: 0 0 26px; place-items: center; width: 26px; height: 26px; border-radius: 8px; background: var(--primary-soft); color: var(--primary-strong); font-size: 13px; }
.options { display: flex; flex-direction: column; gap: 9px; }
.option {
  display: flex;
  align-items: center;
  gap: 12px;
  width: 100%;
  padding: 12px 14px;
  border: 1.5px solid var(--border);
  border-radius: 13px;
  background: var(--surface);
  color: var(--text);
  font: inherit;
  font-size: 14.5px;
  cursor: pointer;
  text-align: left;
  transition: all var(--t-fast);
}
.option__letter { display: grid; place-items: center; width: 28px; height: 28px; flex: 0 0 28px; border-radius: 8px; background: var(--surface-3); color: var(--text-2); font-size: 12.5px; font-weight: 700; }
.option__text { flex: 1; }
.option__mark { color: var(--success); }
.option__mark--wrong { color: var(--danger); }
.option:hover:not(:disabled) { border-color: var(--primary); background: var(--primary-soft); }
.option--selected { border-color: var(--primary); background: var(--primary-soft); font-weight: 600; }
.option--selected .option__letter { background: var(--primary); color: #fff; }
.option--correct { border-color: var(--success); background: var(--success-soft); font-weight: 600; }
.option--correct .option__letter { background: var(--success); color: #fff; }
.option--wrong { border-color: var(--danger); background: var(--danger-soft); }
.option--wrong .option__letter { background: var(--danger); color: #fff; }
.option:disabled { cursor: default; }
.explanation { display: flex; gap: 8px; margin: 12px 0 0; padding: 12px 14px; border-radius: 12px; background: var(--primary-soft); color: var(--primary-deep); font-size: 13.5px; line-height: 1.65; }
.explanation .icon { flex-shrink: 0; margin-top: 3px; }

.result { display: flex; align-items: center; gap: 18px; padding: 20px 22px; border-radius: 16px; background: var(--warning-soft); color: var(--warning-strong); }
.result--passed { background: var(--success-soft); color: var(--success-strong); }
.result__score { display: flex; align-items: baseline; gap: 4px; }
.result__score strong { font-size: 40px; font-weight: 800; line-height: 1; }
.result__text { display: flex; flex-direction: column; }
.result__text strong { font-size: 16px; }
.result__text span { font-size: 13.5px; opacity: 0.85; }

.summary-grid { display: grid; grid-template-columns: repeat(3, 1fr); gap: 12px; margin-bottom: 22px; }
.summary-grid div { display: flex; gap: 10px; padding: 14px; border-radius: 13px; background: var(--surface-2); border: 1px solid var(--border); font-size: 13.5px; line-height: 1.6; }
.summary-grid b { display: grid; flex: 0 0 24px; place-items: center; width: 24px; height: 24px; border-radius: 7px; background: var(--primary-soft); color: var(--primary-strong); font-size: 12px; }

@media (max-width: 640px) {
  .block, .intro-block { padding: 24px 20px; }
  .summary-grid { grid-template-columns: 1fr; }
  .micro-header h1 { font-size: 24px; }
}
</style>
