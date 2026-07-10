import { createApp } from 'vue'
import './styles.css'

const courses = [
  { id: 1, title: 'AI 通识与素养', level: '入门', duration: 12, tag: '通识', desc: '理解人工智能发展、伦理边界与典型教育应用场景。', skills: ['AI 认知', '应用案例', '安全伦理'] },
  { id: 2, title: 'Python 数据分析基础', level: '入门', duration: 24, tag: '编程', desc: '用 Python 完成数据清洗、可视化和学习行为分析。', skills: ['Python', 'Pandas', '可视化'] },
  { id: 3, title: '机器学习项目实战', level: '进阶', duration: 36, tag: '算法', desc: '围绕分类、回归与聚类完成真实项目训练。', skills: ['建模', '评估', '调参'] },
  { id: 4, title: '大模型提示词工程', level: '进阶', duration: 18, tag: '大模型', desc: '掌握教学、办公和业务场景中的提示词设计方法。', skills: ['Prompt', 'RAG', 'Agent'] },
  { id: 5, title: 'AI 应用产品实战营', level: '高级', duration: 48, tag: '产品', desc: '从需求分析到原型、开发协作和上线复盘，完成 AI 产品闭环。', skills: ['产品设计', '低代码', '项目路演'] },
  { id: 6, title: '教师 AI 教学工作坊', level: '入门', duration: 10, tag: '教师', desc: '帮助教师快速构建 AI 备课、出题、批改和学情分析流程。', skills: ['备课', '题库', '评价'] },
]

const questions = [
  { q: '人工智能在课堂中最适合优先落地的场景是？', options: ['完全替代教师授课', '个性化练习与反馈', '取消学习评价'], answer: 1 },
  { q: '提示词工程的核心目标是？', options: ['让模型理解任务、约束和输出格式', '提高电脑硬件配置', '减少所有课程内容'], answer: 0 },
  { q: '教学数据看板主要帮助管理者完成什么？', options: ['隐藏学习风险', '识别学情趋势与干预机会', '替代课程内容'], answer: 1 },
]

const assistantReplies = [
  { keys: ['课程', '学习'], text: '建议先选择目标岗位或教学场景，再从入门课程开始，配合项目实训和阶段测评形成闭环。' },
  { keys: ['教师', '老师', '备课'], text: '教师可优先使用 AI 进行备课大纲、分层练习、课堂互动问题和课后讲评建议生成。' },
  { keys: ['数据', '看板', '学情'], text: '学情看板会聚合出勤、进度、练习、测评和项目表现，帮助定位班级共性问题与个人风险。' },
]

createApp({
  data() {
    return {
      courses,
      questions,
      selectedLevel: '全部',
      keyword: '',
      goal: '提升 AI 通识素养',
      base: '零基础',
      weeklyHours: 4,
      assistantInput: '',
      chat: [{ role: 'assistant', text: '你好，我是 AI 教育助教。你可以咨询课程选择、学习路径、教师备课或数据看板。' }],
      answers: Array(questions.length).fill(null),
      submittedQuiz: false,
      lead: { name: '', phone: '', organization: '', need: '建设 AI 课程体系' },
      leadMessage: '',
    }
  },
  computed: {
    levels() {
      return ['全部', ...new Set(this.courses.map((course) => course.level))]
    },
    filteredCourses() {
      const keyword = this.keyword.trim().toLowerCase()
      return this.courses.filter((course) => {
        const matchLevel = this.selectedLevel === '全部' || course.level === this.selectedLevel
        const matchKeyword = !keyword || [course.title, course.tag, course.desc, ...course.skills].join(' ').toLowerCase().includes(keyword)
        return matchLevel && matchKeyword
      })
    },
    totalHours() {
      return this.filteredCourses.reduce((sum, course) => sum + course.duration, 0)
    },
    planWeeks() {
      return Math.max(1, Math.ceil(this.totalHours / Number(this.weeklyHours || 1)))
    },
    recommendedPath() {
      const order = { 零基础: ['入门', '进阶', '高级'], 有基础: ['进阶', '高级', '入门'], 项目提升: ['高级', '进阶', '入门'] }
      const levelOrder = order[this.base] || order.零基础
      return [...this.courses]
        .sort((a, b) => levelOrder.indexOf(a.level) - levelOrder.indexOf(b.level))
        .slice(0, 4)
    },
    quizScore() {
      return this.answers.reduce((score, item, index) => score + (item === this.questions[index].answer ? 1 : 0), 0)
    },
  },
  methods: {
    askAssistant() {
      const text = this.assistantInput.trim()
      if (!text) return
      this.chat.push({ role: 'user', text })
      const matched = assistantReplies.find((reply) => reply.keys.some((key) => text.includes(key)))
      this.chat.push({ role: 'assistant', text: matched?.text || '可以从“目标—基础—时间”三项信息开始拆解，我会推荐课程组合、练习任务和阶段测评方式。' })
      this.assistantInput = ''
    },
    submitLead() {
      if (!this.lead.name || !this.lead.phone || !this.lead.organization) {
        this.leadMessage = '请完善姓名、联系方式和机构名称。'
        return
      }
      this.leadMessage = `已收到 ${this.lead.organization} 的需求，我们将围绕“${this.lead.need}”准备演示方案。`
    },
  },
  template: `
    <main class="page-shell">
      <section class="hero">
        <nav class="nav" aria-label="主导航">
          <div class="brand">AI Edu<span>智慧教育平台</span></div>
          <div class="nav-links">
            <a href="#courses">课程</a><a href="#path">路径</a><a href="#assistant">助教</a><a href="#contact">咨询</a>
          </div>
        </nav>
        <div class="hero-grid">
          <div>
            <p class="eyebrow">Vue3 + JavaScript 全功能演示</p>
            <h1>面向学校与机构的 AI 教育运营平台</h1>
            <p class="lead">集课程筛选、个性化学习路径、智能助教、在线测评、学情数据和需求提交于一体，覆盖教、学、练、评、管全流程。</p>
            <div class="actions"><a class="primary" href="#courses">开始体验</a><a class="secondary" href="#contact">预约演示</a></div>
          </div>
          <div class="dashboard-card">
            <p class="card-title">实时运营概览</p>
            <div class="metric-row"><b>{{ courses.length }}</b><span>课程模块</span></div>
            <div class="metric-row"><b>{{ totalHours }}</b><span>筛选课时</span></div>
            <div class="metric-row"><b>{{ planWeeks }}</b><span>预计学习周</span></div>
            <div class="progress"><i :style="{ width: Math.min(100, quizScore / questions.length * 100) + '%' }"></i></div>
            <small>测评正确率：{{ quizScore }}/{{ questions.length }}</small>
          </div>
        </div>
      </section>

      <section id="courses" class="section">
        <div class="section-head"><p class="section-kicker">课程中心</p><h2>可搜索、可筛选的课程目录</h2></div>
        <div class="toolbar">
          <input v-model="keyword" placeholder="搜索课程、技能或场景" />
          <select v-model="selectedLevel"><option v-for="level in levels" :key="level">{{ level }}</option></select>
        </div>
        <div class="course-grid">
          <article v-for="course in filteredCourses" :key="course.id" class="course-card">
            <span>{{ course.level }} · {{ course.duration }}课时</span><h3>{{ course.title }}</h3><p>{{ course.desc }}</p>
            <div class="tags"><em v-for="skill in course.skills" :key="skill">{{ skill }}</em></div>
          </article>
        </div>
      </section>

      <section id="path" class="section split-card">
        <div><p class="section-kicker">学习路径</p><h2>根据目标自动生成推荐路径</h2><p class="muted">选择学习目标、基础和每周投入时间，平台会给出课程顺序与周期预估。</p></div>
        <div class="planner">
          <label>学习目标<input v-model="goal" /></label>
          <label>当前基础<select v-model="base"><option>零基础</option><option>有基础</option><option>项目提升</option></select></label>
          <label>每周学习小时<input v-model.number="weeklyHours" type="number" min="1" max="40" /></label>
          <ol><li v-for="course in recommendedPath" :key="course.id">{{ course.title }} <small>{{ course.level }}</small></li></ol>
          <strong>预计 {{ Math.max(1, Math.ceil(recommendedPath.reduce((sum, item) => sum + item.duration, 0) / weeklyHours)) }} 周完成：{{ goal }}</strong>
        </div>
      </section>

      <section id="assistant" class="section two-columns">
        <div class="panel"><p class="section-kicker">AI 助教</p><h2>场景化问答</h2><div class="chat"><p v-for="(item, index) in chat" :key="index" :class="item.role">{{ item.text }}</p></div><form @submit.prevent="askAssistant" class="ask"><input v-model="assistantInput" placeholder="例如：教师如何用 AI 备课？" /><button>发送</button></form></div>
        <div class="panel"><p class="section-kicker">在线测评</p><h2>即时反馈</h2><div v-for="(item, index) in questions" :key="item.q" class="question"><p>{{ index + 1 }}. {{ item.q }}</p><label v-for="(option, optionIndex) in item.options" :key="option"><input v-model="answers[index]" :value="optionIndex" type="radio" />{{ option }}</label></div><button class="primary button" @click="submittedQuiz = true">提交测评</button><strong v-if="submittedQuiz">得分：{{ quizScore }}/{{ questions.length }}</strong></div>
      </section>

      <section id="contact" class="section contact-card">
        <div><p class="section-kicker">预约演示</p><h2>提交建设需求</h2><p class="muted">留下机构信息，获取课程体系、平台部署和运营方案建议。</p></div>
        <form @submit.prevent="submitLead" class="lead-form"><input v-model="lead.name" placeholder="姓名" /><input v-model="lead.phone" placeholder="联系方式" /><input v-model="lead.organization" placeholder="学校/机构名称" /><select v-model="lead.need"><option>建设 AI 课程体系</option><option>教师 AI 培训</option><option>学生项目实训</option><option>教学数据看板</option></select><button>提交需求</button><p>{{ leadMessage }}</p></form>
      </section>
    </main>
  `,
}).mount('#app')
