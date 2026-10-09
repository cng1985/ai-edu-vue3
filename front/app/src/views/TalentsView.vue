<script setup>
import { onMounted, ref } from 'vue'
import { talentApi } from '../api'
import { money } from '../utils/eco'
import SkillLevel from '../components/SkillLevel.vue'
import UserAvatar from '../components/UserAvatar.vue'
import Icon from '../components/Icon.vue'

const list = ref([])
const keyword = ref('')

async function load() {
  list.value = await talentApi.list(keyword.value)
}

onMounted(load)
</script>

<template>
  <div class="page">
    <header class="page-header page-header--split">
      <div>
        <span class="eyebrow">Talent Pool</span>
        <h1>人才库</h1>
        <p>基于个人成长数据的动态人才画像：技能等级、岗位达成度、能力证据与真实任务收入，帮助企业建立人才池。</p>
      </div>
      <form class="row" @submit.prevent="load">
        <input v-model="keyword" class="input" placeholder="搜索人才" />
        <button class="btn btn--primary"><Icon name="search" :size="14" /></button>
      </form>
    </header>

    <div class="grid-2">
      <router-link v-for="t in list" :key="t.user.id" :to="`/talents/${t.user.id}`" class="card card--hover talent">
        <div class="row">
          <UserAvatar :user="t.user" size="lg" />
          <div class="list-row__body">
            <strong>{{ t.user.nickname }}</strong>
            <small>{{ t.goal || '未设定职业目标' }}<template v-if="t.goal"> · 达成度 {{ t.readiness }}%</template></small>
          </div>
          <div class="talent__avg"><b class="num">L{{ t.avgLevel }}</b><small>平均技能</small></div>
        </div>
        <div class="talent__skills">
          <div v-for="s in t.topSkills" :key="s.skill.id" class="row row--between small">
            <span>{{ s.skill.name }}</span><SkillLevel :level="s.level" />
          </div>
        </div>
        <div class="row small muted">
          <span><Icon name="trophy" :size="12" /> {{ t.evidenceCount }} 项证据</span>
          <span><Icon name="gift" :size="12" /> 任务收入 {{ money(t.totalIncome) }}</span>
        </div>
      </router-link>
    </div>
  </div>
</template>

<style scoped>
.talent { display: flex; flex-direction: column; gap: 14px; padding: 20px; color: inherit; }
.talent__avg { display: flex; flex-direction: column; align-items: center; }
.talent__avg b { color: var(--primary-strong); font-size: 20px; }
.talent__avg small { color: var(--text-3); font-size: 11px; }
.talent__skills { display: flex; flex-direction: column; gap: 6px; padding: 12px; border-radius: 12px; background: var(--surface-2); }
</style>
