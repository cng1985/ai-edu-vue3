<template>
  <div class="login">
    <section class="login__brand">
      <div class="login__brand-glow login__brand-glow--a" aria-hidden="true" />
      <div class="login__brand-glow login__brand-glow--b" aria-hidden="true" />
      <div class="login__brand-grid" aria-hidden="true" />

      <div class="login__brand-top">
        <span class="login__logo">AI</span>
        <div>
          <strong>知航 · 运营管理</strong>
          <span>Admin Console</span>
        </div>
      </div>

      <div class="login__brand-body">
        <p class="login__eyebrow">运营 · 内容 · AI 服务</p>
        <h1>一个后台，<br />掌控整条学习链路。</h1>
        <p class="login__lede">
          从课程与题库到 AI 大模型路由、知识库索引与客户咨询，所有运营动作在这里汇聚成一张清晰的看板。
        </p>

        <ul class="login__features">
          <li v-for="f in features" :key="f.title">
            <span class="login__feature-icon">
              <el-icon :size="18"><component :is="f.icon" /></el-icon>
            </span>
            <div>
              <strong>{{ f.title }}</strong>
              <span>{{ f.desc }}</span>
            </div>
          </li>
        </ul>
      </div>

      <div class="login__brand-foot">
        <span>© {{ year }} AI Learning System</span>
        <span class="login__brand-dot" />
        <span>Powered by Vue 3 + Go</span>
      </div>
    </section>

    <section class="login__panel">
      <div class="login__card">
        <div class="login__card-head">
          <p class="login__eyebrow login__eyebrow--dark">欢迎回来</p>
          <h2>登录管理后台</h2>
          <p>使用后台账号登录，权限将根据角色自动分配。</p>
        </div>

        <el-form ref="formRef" :model="form" :rules="rules" class="login__form" @submit.prevent="handleLogin">
          <el-form-item prop="username">
            <label class="login__label">用户名</label>
            <el-input v-model="form.username" placeholder="请输入用户名" size="large" :prefix-icon="User" />
          </el-form-item>
          <el-form-item prop="password">
            <label class="login__label">密码</label>
            <el-input
              v-model="form.password"
              type="password"
              placeholder="请输入密码"
              size="large"
              :prefix-icon="Lock"
              show-password
              @keyup.enter="handleLogin"
            />
          </el-form-item>
          <el-form-item class="login__submit">
            <el-button type="primary" size="large" :loading="loading" class="login__btn" @click="handleLogin">
              登录
              <el-icon class="el-icon--right"><Right /></el-icon>
            </el-button>
          </el-form-item>
        </el-form>

        <div class="login__demo">
          <div class="login__demo-title">演示账号 · 点击快速填入</div>
          <div class="login__demo-list">
            <button
              v-for="acc in demoAccounts"
              :key="acc.username"
              type="button"
              class="login__demo-item"
              @click="fill(acc)"
            >
              <span class="login__demo-role" :class="`login__demo-role--${acc.role}`">{{ acc.label }}</span>
              <span class="login__demo-cred mono">{{ acc.username }} / {{ acc.password }}</span>
            </button>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { User, Lock, Right, DataAnalysis, Reading, Cpu, Service } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { useAuthStore } from '../stores/auth'

const router = useRouter()
const route = useRoute()
const auth = useAuthStore()
const formRef = ref()
const loading = ref(false)
const year = new Date().getFullYear()

const features = [
  { icon: DataAnalysis, title: '实时运营看板', desc: '用户、课程、审核与题库指标一屏可见' },
  { icon: Reading, title: '内容全生命周期', desc: '课程、章节、题库与审核流程一体化管理' },
  { icon: Cpu, title: '大模型路由中心', desc: '厂商 → 统一模型 → 虚拟模型分层配置' },
  { icon: Service, title: '客户咨询实时响应', desc: 'WebSocket 驱动的工单与会话协作' }
]

const demoAccounts = [
  { role: 'admin', label: '管理员', username: 'admin', password: 'admin123' },
  { role: 'reviewer', label: '审核员', username: 'reviewer', password: 'review123' },
  { role: 'operator', label: '运营', username: 'operator', password: 'oper123' }
]

const form = reactive({ username: '', password: '' })
const rules = {
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
}

function fill(acc) {
  form.username = acc.username
  form.password = acc.password
}

async function handleLogin() {
  await formRef.value.validate()
  loading.value = true
  try {
    await auth.login(form.username, form.password)
    ElMessage.success('登录成功')
    router.push(route.query.redirect || '/dashboard')
  } catch {
    /* 错误已由拦截器处理 */
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.login {
  display: grid;
  grid-template-columns: minmax(0, 1.1fr) minmax(0, 1fr);
  min-height: 100vh;
  background: var(--bg);
}

/* 左侧品牌区 */
.login__brand {
  position: relative;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 44px 56px;
  background: linear-gradient(160deg, var(--ink-2) 0%, var(--ink) 60%, #0d0f26 100%);
  color: var(--ink-text);
  overflow: hidden;
}

.login__brand-glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
  pointer-events: none;
}

.login__brand-glow--a {
  width: 520px;
  height: 520px;
  top: -180px;
  left: -140px;
  background: rgba(107, 92, 255, 0.5);
}

.login__brand-glow--b {
  width: 420px;
  height: 420px;
  bottom: -160px;
  right: -100px;
  background: rgba(15, 185, 129, 0.28);
}

.login__brand-grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(rgba(255, 255, 255, 0.04) 1px, transparent 1px),
    linear-gradient(90deg, rgba(255, 255, 255, 0.04) 1px, transparent 1px);
  background-size: 44px 44px;
  mask-image: radial-gradient(ellipse at 30% 40%, #000 30%, transparent 75%);
  pointer-events: none;
}

.login__brand-top,
.login__brand-body,
.login__brand-foot {
  position: relative;
}

.login__brand-top {
  display: flex;
  align-items: center;
  gap: 14px;
}

.login__brand-top div {
  display: flex;
  flex-direction: column;
  line-height: 1.2;
}

.login__brand-top strong {
  font-family: 'Sora', 'Noto Sans SC', sans-serif;
  font-size: 16px;
  color: #fff;
}

.login__brand-top span {
  font-size: 11px;
  letter-spacing: 0.1em;
  text-transform: uppercase;
  color: var(--ink-text-2);
}

.login__logo {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  border-radius: 14px;
  background: linear-gradient(135deg, var(--primary) 0%, #a396ff 100%);
  color: #fff;
  font-family: 'Sora', sans-serif;
  font-weight: 800;
  font-size: 16px;
  box-shadow: 0 12px 28px rgba(107, 92, 255, 0.45), inset 0 1px 0 rgba(255, 255, 255, 0.35);
}

.login__brand-body {
  max-width: 520px;
  padding: 40px 0;
}

.login__eyebrow {
  margin: 0 0 14px;
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.14em;
  text-transform: uppercase;
  color: #a396ff;
}

.login__eyebrow--dark {
  color: var(--primary);
}

.login__brand-body h1 {
  margin: 0;
  font-size: clamp(30px, 3.2vw, 42px);
  line-height: 1.2;
  color: #fff;
  letter-spacing: -0.02em;
}

.login__lede {
  margin: 20px 0 0;
  font-size: 15px;
  line-height: 1.8;
  color: var(--ink-text);
  opacity: 0.85;
}

.login__features {
  list-style: none;
  margin: 36px 0 0;
  padding: 0;
  display: grid;
  gap: 14px;
}

.login__features li {
  display: flex;
  align-items: flex-start;
  gap: 14px;
  padding: 14px 16px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.07);
  backdrop-filter: blur(8px);
}

.login__feature-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 38px;
  height: 38px;
  flex-shrink: 0;
  border-radius: 12px;
  background: rgba(107, 92, 255, 0.22);
  color: #c8c0ff;
}

.login__features li div {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.login__features strong {
  font-size: 14px;
  color: #fff;
}

.login__features span:not(.login__feature-icon) {
  font-size: 13px;
  color: var(--ink-text-2);
}

.login__brand-foot {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 12px;
  color: var(--ink-text-2);
}

.login__brand-dot {
  width: 3px;
  height: 3px;
  border-radius: 50%;
  background: var(--ink-text-2);
}

/* 右侧表单区 */
.login__panel {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 40px 32px;
  background:
    radial-gradient(600px 400px at 100% 0%, rgba(107, 92, 255, 0.08), transparent 60%),
    var(--bg);
}

.login__card {
  width: 100%;
  max-width: 420px;
}

.login__card-head h2 {
  margin: 0;
  font-size: 28px;
}

.login__card-head p:last-child {
  margin: 10px 0 0;
  color: var(--text-2);
  font-size: 14px;
}

.login__form {
  margin-top: 32px;
}

.login__form :deep(.el-form-item) {
  margin-bottom: 20px;
}

.login__form :deep(.el-form-item__content) {
  flex-direction: column;
  align-items: stretch;
}

.login__label {
  display: block;
  margin-bottom: 8px;
  font-size: 13px;
  font-weight: 600;
  color: var(--text-2);
}

.login__form :deep(.el-input--large .el-input__wrapper) {
  padding: 4px 14px;
  border-radius: 12px !important;
}

.login__submit {
  margin-top: 8px;
}

.login__btn {
  width: 100%;
  height: 48px;
  font-size: 15px;
  border-radius: 13px;
}

.login__demo {
  margin-top: 28px;
  padding-top: 22px;
  border-top: 1px dashed var(--border-strong);
}

.login__demo-title {
  font-size: 12px;
  font-weight: 700;
  letter-spacing: 0.06em;
  color: var(--text-3);
  margin-bottom: 10px;
}

.login__demo-list {
  display: grid;
  gap: 8px;
}

.login__demo-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 10px 14px;
  border: 1px solid var(--border);
  border-radius: 12px;
  background: var(--surface);
  cursor: pointer;
  text-align: left;
  transition: all var(--t-fast);
}

.login__demo-item:hover {
  border-color: var(--primary);
  background: var(--primary-soft);
  transform: translateY(-1px);
}

.login__demo-role {
  font-size: 12.5px;
  font-weight: 700;
  padding: 2px 10px;
  border-radius: 999px;
}

.login__demo-role--admin { background: var(--danger-soft); color: var(--danger-strong); }
.login__demo-role--reviewer { background: var(--warning-soft); color: var(--warning-strong); }
.login__demo-role--operator { background: var(--info-soft); color: #1c5fb0; }

.login__demo-cred {
  color: var(--text-2);
}

@media (max-width: 960px) {
  .login {
    grid-template-columns: 1fr;
  }

  .login__brand {
    padding: 28px 24px;
  }

  .login__brand-body {
    padding: 24px 0;
  }

  .login__features {
    display: none;
  }

  .login__brand-foot {
    display: none;
  }

  .login__panel {
    padding: 32px 20px 48px;
  }
}
</style>
