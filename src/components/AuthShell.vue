<script setup>
import { courses, totalChapterCount } from '../data/courses'

const stats = [
  { value: '12,800+', label: '注册学习者' },
  { value: courses.length, label: '精品课程' },
  { value: totalChapterCount, label: '实战章节' },
  { value: '96%', label: '好评率' }
]

const posts = [
  {
    avatar: '🦊',
    name: '前端小狐',
    tag: 'RAG 实战',
    text: '照着课程把混合检索 + 重排序跑通了,召回质量肉眼可见地提升,强推!'
  },
  {
    avatar: '🐼',
    name: 'Panda_Dev',
    tag: 'Agent',
    text: '第一次搞懂 ReAct 循环和 Function Calling 的边界,测验错题解析很到位。'
  },
  {
    avatar: '🚀',
    name: '架构老王',
    tag: 'AI 原生',
    text: '从 CRUD 思维切到"调度智能"确实需要一个过程,这套课程路线很系统。'
  }
]

const members = ['🧑‍💻', '👩‍💻', '🦉', '🐯', '🌟', '🤖']
</script>

<template>
  <div class="auth">
    <aside class="auth__side">
      <router-link to="/login" class="auth__brand">
        <span class="auth__logo">🧠</span>
        <span>AI 学习系统</span>
      </router-link>

      <div class="auth__slogan">
        <h1>加入 AI 学习者社区</h1>
        <p>与上万名开发者一起,从提示词工程到 AI 原生架构,完成从"写代码"到"调度智能"的跃迁。</p>
      </div>

      <div class="auth__stats">
        <div v-for="s in stats" :key="s.label" class="auth__stat">
          <strong>{{ s.value }}</strong>
          <span>{{ s.label }}</span>
        </div>
      </div>

      <div class="auth__posts">
        <div v-for="p in posts" :key="p.name" class="auth__post">
          <span class="auth__post-avatar">{{ p.avatar }}</span>
          <div class="auth__post-body">
            <div class="auth__post-head">
              <strong>{{ p.name }}</strong>
              <span class="auth__post-tag"># {{ p.tag }}</span>
            </div>
            <p>{{ p.text }}</p>
          </div>
        </div>
      </div>

      <div class="auth__members">
        <span v-for="(m, i) in members" :key="i" class="auth__member">{{ m }}</span>
        <span class="auth__members-text">今天已有 328 位同学在学习</span>
      </div>
    </aside>

    <main class="auth__main">
      <div class="auth__card fade-up">
        <slot />
      </div>
      <p class="auth__footnote">
        演示版账号数据仅保存在当前浏览器本地,请勿使用真实密码。
      </p>
    </main>
  </div>
</template>

<style scoped>
.auth {
  display: flex;
  min-height: 100vh;
}

.auth__side {
  width: 46%;
  min-width: 420px;
  padding: 44px 52px;
  display: flex;
  flex-direction: column;
  gap: 30px;
  color: #eef2ff;
  background:
    radial-gradient(900px 500px at -10% -10%, rgba(167, 139, 250, 0.35), transparent 60%),
    radial-gradient(700px 480px at 110% 110%, rgba(56, 189, 248, 0.28), transparent 60%),
    linear-gradient(155deg, #312e81 0%, #4338ca 55%, #1e1b4b 100%);
}

.auth__brand {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 18px;
  font-weight: 800;
  color: #fff;
  letter-spacing: -0.01em;
}

.auth__logo {
  font-size: 26px;
}

.auth__slogan h1 {
  margin: 0 0 10px;
  font-size: 30px;
  line-height: 1.35;
  color: #fff;
  letter-spacing: -0.02em;
}

.auth__slogan p {
  margin: 0;
  font-size: 14.5px;
  line-height: 1.8;
  color: rgba(238, 242, 255, 0.82);
  max-width: 440px;
}

.auth__stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 12px;
}

.auth__stat {
  background: rgba(255, 255, 255, 0.09);
  border: 1px solid rgba(255, 255, 255, 0.14);
  border-radius: 12px;
  padding: 13px 10px;
  text-align: center;
  backdrop-filter: blur(4px);
}

.auth__stat strong {
  display: block;
  font-size: 17px;
  color: #fff;
}

.auth__stat span {
  font-size: 11.5px;
  color: rgba(238, 242, 255, 0.72);
}

.auth__posts {
  display: flex;
  flex-direction: column;
  gap: 12px;
  flex: 1;
}

.auth__post {
  display: flex;
  gap: 12px;
  padding: 14px 16px;
  background: rgba(255, 255, 255, 0.08);
  border: 1px solid rgba(255, 255, 255, 0.13);
  border-radius: 14px;
  backdrop-filter: blur(4px);
}

.auth__post-avatar {
  width: 38px;
  height: 38px;
  min-width: 38px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 20px;
  background: rgba(255, 255, 255, 0.14);
  border-radius: 50%;
}

.auth__post-head {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 3px;
}

.auth__post-head strong {
  font-size: 13.5px;
  color: #fff;
}

.auth__post-tag {
  font-size: 11.5px;
  padding: 1px 9px;
  border-radius: 999px;
  background: rgba(167, 139, 250, 0.28);
  color: #ddd6fe;
}

.auth__post-body p {
  margin: 0;
  font-size: 13px;
  line-height: 1.7;
  color: rgba(238, 242, 255, 0.85);
}

.auth__members {
  display: flex;
  align-items: center;
}

.auth__member {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  background: rgba(255, 255, 255, 0.16);
  border: 2px solid rgba(49, 46, 129, 0.9);
  border-radius: 50%;
  margin-left: -8px;
}

.auth__member:first-child {
  margin-left: 0;
}

.auth__members-text {
  margin-left: 14px;
  font-size: 12.5px;
  color: rgba(238, 242, 255, 0.78);
}

.auth__main {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 24px;
  background: var(--bg);
}

.auth__card {
  width: 100%;
  max-width: 420px;
  background: var(--surface);
  border: 1px solid var(--border);
  border-radius: 18px;
  box-shadow: var(--shadow);
  padding: 36px 38px;
}

.auth__footnote {
  margin: 18px 0 0;
  font-size: 12px;
  color: var(--text-3);
  text-align: center;
  max-width: 420px;
}

@media (max-width: 920px) {
  .auth {
    flex-direction: column;
  }

  .auth__side {
    width: 100%;
    min-width: 0;
    padding: 28px 22px;
    gap: 20px;
  }

  .auth__posts,
  .auth__members {
    display: none;
  }

  .auth__slogan h1 {
    font-size: 23px;
  }

  .auth__card {
    padding: 26px 22px;
  }
}
</style>
