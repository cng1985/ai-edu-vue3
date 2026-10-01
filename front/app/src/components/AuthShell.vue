<script setup>
import ParticleField from './ParticleField.vue'
import Icon from './Icon.vue'

defineProps({
  mode: { type: String, default: 'login' }
})

const highlights = [
  { icon: 'compass', title: '职业 → 能力 → 技能', desc: '选择目标岗位，AI 分析你与岗位要求的技能差距' },
  { icon: 'brain', title: '知识图谱与 AI 伙伴', desc: '按前置关系学习，六个 Agent 陪伴规划、讲解与评审' },
  { icon: 'target', title: '项目实践与真实任务', desc: '用项目证据与企业任务证明能力，获得真实收益' }
]

const stats = [
  { value: '45', label: '知识点图谱' },
  { value: '19', label: '技能模型' },
  { value: '6', label: 'AI Agent' }
]
</script>

<template>
  <div class="auth">
    <section class="auth__brand">
      <ParticleField />
      <div class="auth__brand-inner">
        <router-link to="/login" class="auth__logo">
          <span class="auth__mark"><Icon name="compass" :size="20" :stroke="2.2" /></span>
          <span class="auth__brand-name">知航</span>
        </router-link>

        <div class="auth__copy">
          <span class="auth__eyebrow">AI 时代个人成长操作系统</span>
          <h1>
            {{ mode === 'register' ? '加入知识工作生态' : '学习、实践、工作' }}<br />
            <em>{{ mode === 'register' ? '让能力持续增值' : '持续提升个人价值' }}</em>
          </h1>
          <p>
            {{
              mode === 'register'
                ? '注册后形成你的人才画像：学习数据、项目证据、任务表现与收入，共同证明你的能力。'
                : '从职业目标出发，AI 规划学习，社区交流与项目实践把知识变成能力，再通过真实任务获得收益。'
            }}
          </p>
        </div>

        <ul class="auth__features">
          <li v-for="item in highlights" :key="item.title">
            <span class="auth__feature-icon"><Icon :name="item.icon" :size="18" /></span>
            <div>
              <strong>{{ item.title }}</strong>
              <span>{{ item.desc }}</span>
            </div>
          </li>
        </ul>

        <div class="auth__stats">
          <div v-for="item in stats" :key="item.label">
            <strong class="num">{{ item.value }}</strong>
            <span>{{ item.label }}</span>
          </div>
        </div>
      </div>
    </section>

    <section class="auth__panel">
      <div class="auth__frame">
        <header class="auth__head">
          <h2>{{ mode === 'register' ? '创建账号' : '欢迎回来' }}</h2>
          <p>
            {{ mode === 'register' ? '填写以下信息即可开始学习之旅。' : '登录以继续你的学习路径。' }}
          </p>
        </header>

        <nav class="auth__tabs segmented" aria-label="账号入口">
          <router-link
            class="auth__tab"
            :class="{ active: mode === 'login' }"
            :to="{ path: '/login', query: $route.query }"
          >
            登录
          </router-link>
          <router-link
            class="auth__tab"
            :class="{ active: mode === 'register' }"
            :to="{ path: '/register', query: $route.query }"
          >
            注册
          </router-link>
        </nav>

        <div class="auth__body">
          <slot />
        </div>

        <p class="auth__foot">演示环境 · 请勿使用真实密码</p>
      </div>
    </section>
  </div>
</template>

<style scoped>
.auth {
  display: grid;
  grid-template-columns: minmax(0, 1.05fr) minmax(0, 1fr);
  min-height: 100vh;
  background: var(--bg);
}

/* 左侧品牌区 */
.auth__brand {
  position: relative;
  overflow: hidden;
  color: #fff;
  background:
    radial-gradient(80% 60% at 100% 0%, rgba(107, 92, 255, 0.5), transparent 60%),
    radial-gradient(60% 50% at 0% 100%, rgba(15, 185, 129, 0.3), transparent 60%),
    linear-gradient(160deg, #1b1d3e 0%, #121430 100%);
}

.auth__brand :deep(.particles) {
  position: absolute;
  opacity: 0.55;
}

.auth__brand-inner {
  position: relative;
  z-index: 1;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  gap: 40px;
  height: 100%;
  min-height: 100vh;
  padding: 44px 56px;
  max-width: 640px;
}

.auth__logo {
  display: inline-flex;
  align-items: center;
  gap: 12px;
  color: #fff;
  width: fit-content;
}

.auth__mark {
  display: grid;
  place-items: center;
  width: 42px;
  height: 42px;
  border-radius: 13px;
  background: linear-gradient(135deg, var(--primary), #a071ff);
  box-shadow: 0 10px 24px var(--primary-glow), inset 0 1px 0 rgba(255, 255, 255, 0.35);
}

.auth__brand-name {
  font-family: 'Sora', 'Noto Sans SC', sans-serif;
  font-size: 18px;
  font-weight: 700;
  letter-spacing: -0.01em;
}

.auth__eyebrow {
  display: inline-block;
  margin-bottom: 18px;
  padding: 5px 12px;
  border-radius: 999px;
  border: 1px solid rgba(255, 255, 255, 0.16);
  background: rgba(255, 255, 255, 0.06);
  color: #c7bfff;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.08em;
}

.auth__copy h1 {
  margin: 0 0 16px;
  font-size: 44px;
  line-height: 1.15;
  color: #fff;
}

.auth__copy h1 em {
  font-style: normal;
  background: linear-gradient(90deg, #c7bfff, #6ee7b7);
  -webkit-background-clip: text;
  background-clip: text;
  color: transparent;
}

.auth__copy p {
  margin: 0;
  max-width: 460px;
  color: rgba(201, 203, 230, 0.85);
  font-size: 16px;
  line-height: 1.75;
}

.auth__features {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.auth__features li {
  display: flex;
  gap: 14px;
  align-items: flex-start;
  padding: 14px 16px;
  border-radius: 16px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  background: rgba(255, 255, 255, 0.04);
  backdrop-filter: blur(8px);
}

.auth__feature-icon {
  display: grid;
  place-items: center;
  width: 38px;
  height: 38px;
  flex: 0 0 38px;
  border-radius: 12px;
  background: rgba(107, 92, 255, 0.25);
  color: #c7bfff;
}

.auth__features div {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.auth__features strong {
  font-size: 14.5px;
}

.auth__features span:last-child {
  color: rgba(201, 203, 230, 0.7);
  font-size: 13px;
}

.auth__stats {
  display: flex;
  gap: 36px;
  padding-top: 24px;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.auth__stats div {
  display: flex;
  flex-direction: column;
}

.auth__stats strong {
  font-size: 28px;
  font-weight: 700;
  line-height: 1.1;
}

.auth__stats span {
  color: rgba(201, 203, 230, 0.65);
  font-size: 12.5px;
}

/* 右侧表单区 */
.auth__panel {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 32px;
  background:
    radial-gradient(60% 40% at 80% 0%, rgba(107, 92, 255, 0.1), transparent 60%),
    var(--bg);
}

.auth__frame {
  width: 100%;
  max-width: 420px;
  animation: fade-up 0.45s var(--ease) both;
}

.auth__head {
  margin-bottom: 22px;
}

.auth__head h2 {
  margin: 0 0 6px;
  font-size: 28px;
}

.auth__head p {
  margin: 0;
  color: var(--text-2);
  font-size: 14.5px;
}

.auth__tabs {
  display: grid;
  grid-template-columns: 1fr 1fr;
  width: 100%;
  margin-bottom: 22px;
}

.auth__tab {
  padding: 9px 0;
  text-align: center;
  border-radius: 9px;
  font-size: 14px;
  font-weight: 600;
  color: var(--text-2);
  transition: all var(--t-fast);
}

.auth__tab.active {
  background: var(--surface);
  color: var(--primary-strong);
  box-shadow: var(--shadow-xs);
}

.auth__body {
  padding: 28px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow);
}

.auth__foot {
  margin: 18px 0 0;
  text-align: center;
  font-size: 12px;
  color: var(--text-3);
}

@media (max-width: 960px) {
  .auth {
    grid-template-columns: 1fr;
  }

  .auth__brand-inner {
    min-height: 0;
    padding: 32px 24px;
    gap: 28px;
  }

  .auth__copy h1 {
    font-size: 30px;
  }

  .auth__features,
  .auth__stats {
    display: none;
  }

  .auth__panel {
    padding: 28px 18px 40px;
  }
}
</style>
