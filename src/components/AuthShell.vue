<script setup>
import { courses, totalChapterCount } from '../data/courses'

defineProps({
  mode: { type: String, default: 'login' } // login | register
})

const stats = [
  { value: '1.2万+', label: '学习者' },
  { value: String(courses.length), label: '精品课程' },
  { value: String(totalChapterCount), label: '实战章节' }
]

const quote = {
  avatar: '🦊',
  name: '前端小狐',
  role: 'RAG 课程学员',
  text: '照着课程把混合检索与重排序跑通之后,召回质量肉眼可见地提升。这是我用过最系统的 AI 工程路线。'
}
</script>

<template>
  <div class="auth" :class="`auth--${mode}`">
    <aside class="panel">
      <div class="panel__glow panel__glow--a" aria-hidden="true"></div>
      <div class="panel__glow panel__glow--b" aria-hidden="true"></div>
      <div class="panel__grid" aria-hidden="true"></div>

      <div class="panel__inner">
        <div class="panel__brand font-display">
          <span class="panel__mark">AI</span>
          <span class="panel__name">学习系统</span>
        </div>

        <div class="panel__hero">
          <p class="panel__eyebrow">Developer Community</p>
          <h1 class="font-display">
            {{ mode === 'register' ? '开启你的 AI 学习之旅' : '回到你的学习社区' }}
          </h1>
          <p class="panel__lead">
            从提示词工程到 AI 原生架构,与上万开发者一起完成从「写代码」到「调度智能」的跃迁。
          </p>
        </div>

        <div class="panel__stats">
          <div v-for="s in stats" :key="s.label" class="panel__stat">
            <strong class="font-display">{{ s.value }}</strong>
            <span>{{ s.label }}</span>
          </div>
        </div>

        <figure class="panel__quote">
          <blockquote>{{ quote.text }}</blockquote>
          <figcaption>
            <span class="panel__quote-avatar">{{ quote.avatar }}</span>
            <span>
              <strong>{{ quote.name }}</strong>
              <em>{{ quote.role }}</em>
            </span>
          </figcaption>
        </figure>
      </div>
    </aside>

    <main class="stage">
      <div class="stage__card">
        <slot />
      </div>
      <p class="stage__note">演示版账号仅保存在当前浏览器本地,请勿使用真实密码。</p>
    </main>
  </div>
</template>

<style scoped>
.auth {
  display: grid;
  grid-template-columns: minmax(0, 1.05fr) minmax(0, 0.95fr);
  min-height: 100vh;
  background: #0b1020;
}

.panel {
  position: relative;
  overflow: hidden;
  color: #e8edf8;
  background:
    linear-gradient(165deg, #101935 0%, #162455 48%, #0f172a 100%);
}

.panel__glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(2px);
  pointer-events: none;
}

.panel__glow--a {
  width: 420px;
  height: 420px;
  top: -120px;
  right: -80px;
  background: radial-gradient(circle, rgba(99, 102, 241, 0.45), transparent 68%);
  animation: float-a 12s ease-in-out infinite;
}

.panel__glow--b {
  width: 340px;
  height: 340px;
  bottom: -80px;
  left: -60px;
  background: radial-gradient(circle, rgba(14, 165, 233, 0.32), transparent 70%);
  animation: float-b 14s ease-in-out infinite;
}

.panel__grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.04) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.04) 1px, transparent 1px);
  background-size: 48px 48px;
  mask-image: radial-gradient(ellipse at 40% 30%, black 20%, transparent 75%);
  opacity: 0.7;
  pointer-events: none;
}

.panel__inner {
  position: relative;
  z-index: 1;
  height: 100%;
  display: flex;
  flex-direction: column;
  gap: 36px;
  padding: 48px 56px;
}

.panel__brand {
  display: flex;
  align-items: center;
  gap: 12px;
}

.panel__mark {
  width: 42px;
  height: 42px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 12px;
  background: linear-gradient(135deg, #6366f1, #38bdf8);
  color: #fff;
  font-size: 15px;
  font-weight: 800;
  letter-spacing: 0.02em;
  box-shadow: 0 8px 24px rgba(99, 102, 241, 0.35);
}

.panel__name {
  font-size: 20px;
  font-weight: 700;
  color: #fff;
  letter-spacing: -0.02em;
}

.panel__hero {
  max-width: 460px;
  animation: rise 0.7s ease both;
}

.panel__eyebrow {
  margin: 0 0 14px;
  font-size: 12px;
  font-weight: 600;
  letter-spacing: 0.16em;
  text-transform: uppercase;
  color: #93c5fd;
}

.panel__hero h1 {
  margin: 0 0 16px;
  font-size: clamp(28px, 3.2vw, 40px);
  line-height: 1.25;
  font-weight: 700;
  color: #fff;
  letter-spacing: -0.03em;
}

.panel__lead {
  margin: 0;
  font-size: 15px;
  line-height: 1.85;
  color: rgba(226, 232, 240, 0.82);
}

.panel__stats {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: 12px;
  max-width: 420px;
  animation: rise 0.7s ease 0.1s both;
}

.panel__stat {
  padding: 16px 14px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.06);
  border: 1px solid rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
}

.panel__stat strong {
  display: block;
  margin-bottom: 4px;
  font-size: 22px;
  font-weight: 700;
  color: #fff;
  letter-spacing: -0.02em;
}

.panel__stat span {
  font-size: 12px;
  color: rgba(203, 213, 225, 0.85);
}

.panel__quote {
  margin: auto 0 0;
  max-width: 460px;
  padding: 22px 24px;
  border-radius: 18px;
  background: rgba(15, 23, 42, 0.35);
  border: 1px solid rgba(148, 163, 184, 0.18);
  backdrop-filter: blur(12px);
  animation: rise 0.7s ease 0.18s both;
}

.panel__quote blockquote {
  margin: 0 0 16px;
  font-size: 14.5px;
  line-height: 1.8;
  color: rgba(241, 245, 249, 0.92);
}

.panel__quote figcaption {
  display: flex;
  align-items: center;
  gap: 12px;
}

.panel__quote-avatar {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background: rgba(99, 102, 241, 0.25);
  border: 1px solid rgba(165, 180, 252, 0.35);
  font-size: 20px;
}

.panel__quote figcaption strong {
  display: block;
  font-size: 13.5px;
  color: #fff;
}

.panel__quote figcaption em {
  font-style: normal;
  font-size: 12px;
  color: #94a3b8;
}

.stage {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 40px 28px;
  background:
    radial-gradient(700px 420px at 90% 0%, rgba(99, 102, 241, 0.08), transparent 60%),
    radial-gradient(520px 360px at 10% 100%, rgba(14, 165, 233, 0.07), transparent 55%),
    #f5f7fc;
}

.stage__card {
  width: 100%;
  max-width: 440px;
  padding: 40px 40px 34px;
  border-radius: 24px;
  background: #fff;
  border: 1px solid rgba(226, 232, 240, 0.9);
  box-shadow:
    0 1px 2px rgba(15, 23, 42, 0.04),
    0 18px 48px rgba(15, 23, 42, 0.08);
  animation: rise 0.55s ease both;
}

.stage__note {
  margin: 18px 0 0;
  max-width: 440px;
  text-align: center;
  font-size: 12px;
  color: #94a3b8;
}

@keyframes float-a {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(-24px, 18px); }
}

@keyframes float-b {
  0%, 100% { transform: translate(0, 0); }
  50% { transform: translate(18px, -22px); }
}

@keyframes rise {
  from {
    opacity: 0;
    transform: translateY(14px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 980px) {
  .auth {
    grid-template-columns: 1fr;
  }

  .panel {
    min-height: auto;
  }

  .panel__inner {
    padding: 28px 22px 32px;
    gap: 22px;
  }

  .panel__hero h1 {
    font-size: 26px;
  }

  .panel__quote {
    display: none;
  }

  .stage {
    padding: 28px 16px 40px;
  }

  .stage__card {
    padding: 28px 22px;
    border-radius: 20px;
  }
}
</style>
