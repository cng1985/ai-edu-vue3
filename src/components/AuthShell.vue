<script setup>
import ParticleField from './ParticleField.vue'

defineProps({
  mode: { type: String, default: 'login' }
})
</script>

<template>
  <div class="auth">
    <ParticleField />

    <div class="auth__frame">
      <header class="auth__hero">
        <div class="auth__logo">
          <span class="auth__mark">AI</span>
          <span class="auth__brand">学习社区</span>
        </div>
        <h1>{{ mode === 'register' ? '加入社区，把知识学透' : '与同行一起，把 AI 做成产品' }}</h1>
        <p>
          {{
            mode === 'register'
              ? '注册后同步学习进度、笔记与测验成绩，和社区学员一起推进。'
              : '提示词 · RAG · Agent 架构 —— 面向开发者的系统化知识社区。'
          }}
        </p>
      </header>

      <section class="auth__card">
        <nav class="auth__tabs" aria-label="账号入口">
          <router-link
            class="auth__tab"
            :class="{ 'auth__tab--on': mode === 'login' }"
            :to="{ path: '/login', query: $route.query }"
          >
            登录
          </router-link>
          <router-link
            class="auth__tab"
            :class="{ 'auth__tab--on': mode === 'register' }"
            :to="{ path: '/register', query: $route.query }"
          >
            注册
          </router-link>
        </nav>

        <div class="auth__body">
          <slot />
        </div>
      </section>

      <p class="auth__foot">演示环境 · 账号仅保存在本机浏览器 · 请勿使用真实密码</p>
    </div>
  </div>
</template>

<style scoped>
.auth {
  position: relative;
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 20px;
  background: #f6f8fb;
  overflow: hidden;
}

.auth__frame {
  position: relative;
  z-index: 1;
  width: 100%;
  max-width: 420px;
  animation: rise 0.45s ease both;
}

.auth__hero {
  text-align: center;
  margin-bottom: 22px;
}

.auth__logo {
  display: inline-flex;
  align-items: center;
  gap: 10px;
  margin-bottom: 18px;
}

.auth__mark {
  width: 36px;
  height: 36px;
  border-radius: 8px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  background: #1772f6;
  color: #fff;
  font-size: 13px;
  font-weight: 800;
  letter-spacing: 0.02em;
  font-family: 'Sora', 'Noto Sans SC', sans-serif;
}

.auth__brand {
  font-size: 18px;
  font-weight: 700;
  color: #1a1a1a;
  letter-spacing: -0.02em;
  font-family: 'Sora', 'Noto Sans SC', sans-serif;
}

.auth__hero h1 {
  margin: 0 0 10px;
  font-size: 26px;
  line-height: 1.35;
  font-weight: 700;
  color: #121212;
  letter-spacing: -0.02em;
}

.auth__hero p {
  margin: 0 auto;
  max-width: 340px;
  font-size: 14px;
  line-height: 1.7;
  color: #6b7280;
}

.auth__card {
  background: #fff;
  border-radius: 4px;
  box-shadow: 0 1px 3px rgba(26, 26, 26, 0.06), 0 8px 28px rgba(26, 26, 26, 0.06);
  overflow: hidden;
}

.auth__tabs {
  display: grid;
  grid-template-columns: 1fr 1fr;
  border-bottom: 1px solid #ebebeb;
}

.auth__tab {
  padding: 15px 0;
  text-align: center;
  font-size: 15px;
  font-weight: 500;
  color: #8590a6;
  position: relative;
  transition: color 0.15s ease;
}

.auth__tab:hover {
  color: #121212;
}

.auth__tab--on {
  color: #121212;
  font-weight: 600;
}

.auth__tab--on::after {
  content: '';
  position: absolute;
  left: 28%;
  right: 28%;
  bottom: 0;
  height: 2px;
  background: #1772f6;
  border-radius: 1px;
}

.auth__body {
  padding: 28px 32px 30px;
}

.auth__foot {
  margin: 18px 0 0;
  text-align: center;
  font-size: 12px;
  color: #a0a6b0;
}

@keyframes rise {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 480px) {
  .auth__hero h1 {
    font-size: 22px;
  }

  .auth__body {
    padding: 24px 20px 26px;
  }
}
</style>
