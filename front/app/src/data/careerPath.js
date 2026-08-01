export const careers = [
  {
    id: 'frontend',
    category: '研发 · 前端',
    name: 'Web 前端工程师',
    icon: '🖥️',
    demand: '高',
    salary: '8–25K',
    difficulty: 3,
    duration: 16,
    match: 92,
    description: '用 HTML、CSS、JavaScript 和 Vue 构建高质量 Web 应用。',
    reasons: ['学习成果直观，适合快速建立正反馈', '岗位需求稳定，项目作品可直接展示', '已有课程可作为 AI 应用进阶方向']
  },
  {
    id: 'java',
    category: '研发 · 后端',
    name: 'Java 后端工程师',
    icon: '☕',
    demand: '高',
    salary: '10–30K',
    difficulty: 4,
    duration: 24,
    match: 84,
    description: '构建稳定、可扩展的服务端系统与企业级应用。',
    reasons: ['企业岗位覆盖广', '适合喜欢业务逻辑与系统设计的学习者']
  },
  {
    id: 'data',
    category: '数据 · 分析',
    name: '数据分析师',
    icon: '📈',
    demand: '高',
    salary: '8–20K',
    difficulty: 3,
    duration: 16,
    match: 81,
    description: '通过 SQL、Python 和可视化工具从数据中提炼业务洞察。',
    reasons: ['技术与业务结合紧密', '入门路径清晰，适用行业广']
  }
]

export const frontendPath = {
  careerId: 'frontend',
  name: '16 周成为初级 Web 前端工程师',
  description: '从网页基础到 Vue 工程化，通过微学习、练习与项目建立可验证的前端能力。',
  weeklyHours: 12,
  durationWeeks: 16,
  competencies: [
    {
      id: 'html-css',
      name: 'HTML / CSS',
      weight: 20,
      color: '#10b981',
      points: [
        { id: 'semantic-html', name: 'HTML 语义化', unitId: 'semantic-html', weight: 40 },
        { id: 'flex-layout', name: 'Flex 布局', unitId: 'flex-layout', weight: 60 }
      ]
    },
    {
      id: 'javascript',
      name: 'JavaScript 核心',
      weight: 30,
      color: '#f59e0b',
      points: [
        { id: 'js-variables', name: '变量与数据类型', unitId: 'js-variables', weight: 35 },
        { id: 'promise', name: 'Promise 异步编程', unitId: 'promise', weight: 65 }
      ]
    },
    {
      id: 'vue',
      name: 'Vue 应用开发',
      weight: 30,
      color: '#6366f1',
      points: [
        { id: 'vue-reactivity', name: 'Vue 响应式', unitId: 'vue-reactivity', weight: 45 },
        { id: 'vue-component', name: '组件设计', unitId: 'vue-component', weight: 55 }
      ]
    },
    {
      id: 'engineering',
      name: '工程化与实战',
      weight: 20,
      color: '#0ea5e9',
      points: [
        { id: 'git-basics', name: 'Git 协作', unitId: 'git-basics', weight: 40 },
        { id: 'vite-build', name: 'Vite 构建部署', unitId: 'vite-build', weight: 60 }
      ]
    }
  ],
  milestones: [
    { id: 'm1', name: '完成响应式静态页面', week: 4, standard: 'HTML/CSS 达成度 ≥ 75%' },
    { id: 'm2', name: '完成交互式 Todo 应用', week: 8, standard: 'JavaScript 达成度 ≥ 75%' },
    { id: 'm3', name: '完成 Vue 管理后台', week: 13, standard: 'Vue 达成度 ≥ 75%' },
    { id: 'm4', name: '作品部署并通过模拟面试', week: 16, standard: '总达成度 ≥ 75%' }
  ]
}

export const microUnits = [
  {
    id: 'semantic-html',
    title: '5 分钟理解 HTML 语义化',
    competency: 'HTML / CSS',
    duration: 4,
    difficulty: '入门',
    intro: '语义化标签就像给网页每个房间挂上门牌，让人和机器都能快速看懂内容结构。',
    content: [
      { title: '为什么不用全是 div？', text: 'div 只表示一个容器，而 header、nav、main、article、footer 会说明内容的角色。它们让代码更易维护，也帮助搜索引擎与屏幕阅读器理解页面。' },
      { title: '选择标签的方法', text: '先问“这段内容是什么”，再选标签：独立文章用 article，页面主内容用 main，导航链接用 nav。不要只为默认样式选择标签。' }
    ],
    example: '<main>\n  <article>\n    <h1>学习语义化</h1>\n    <p>内容正文</p>\n  </article>\n</main>',
    questions: [
      { text: '页面中主要的导航链接最适合使用哪个标签？', options: ['div', 'nav', 'span'], answer: 1, explanation: 'nav 明确表示一组导航链接。' },
      { text: '语义化标签的主要价值是什么？', options: ['自动美化页面', '提升结构可读性与可访问性', '减少所有 CSS'], answer: 1, explanation: '语义化关注含义与结构，不会自动完成视觉设计。' }
    ],
    summary: ['根据内容角色选择标签', '一个页面通常只有一个主要 main', '语义化改善维护、SEO 与无障碍体验']
  },
  {
    id: 'flex-layout',
    title: '5 分钟掌握 Flex 对齐',
    competency: 'HTML / CSS',
    duration: 5,
    difficulty: '入门',
    intro: 'Flex 像一排可自动伸缩的座位：你只需指定方向和对齐规则，浏览器负责排布。',
    content: [
      { title: '主轴与交叉轴', text: 'display: flex 开启弹性布局。默认主轴从左到右，justify-content 控制主轴，align-items 控制交叉轴。' },
      { title: '最常用的居中', text: '同时使用 justify-content: center 与 align-items: center，就能让子项水平、垂直居中。gap 可统一设置子项间距。' }
    ],
    example: '.container {\n  display: flex;\n  justify-content: center;\n  align-items: center;\n  gap: 16px;\n}',
    questions: [
      { text: '默认横向 Flex 中，水平居中使用？', options: ['align-items', 'justify-content', 'flex-wrap'], answer: 1, explanation: '默认主轴是水平方向，justify-content 控制主轴对齐。' },
      { text: '统一设置子项间距应优先使用？', options: ['gap', 'padding', 'position'], answer: 0, explanation: 'gap 能直接描述项目之间的间距。' }
    ],
    summary: ['Flex 默认主轴为水平方向', 'justify-content 管主轴', 'align-items 管交叉轴']
  },
  {
    id: 'js-variables',
    title: '4 分钟分清 const 与 let',
    competency: 'JavaScript 核心',
    duration: 4,
    difficulty: '入门',
    intro: '变量像贴了标签的盒子；const 表示标签不能换盒子，let 表示之后可以重新指向别的值。',
    content: [
      { title: '默认使用 const', text: '当变量不需要重新赋值时使用 const，可清晰表达意图并减少误改。需要计数或状态变化时才使用 let。' },
      { title: '对象仍可修改', text: 'const 限制的是变量绑定，不会冻结对象。const user = {} 后仍可执行 user.name = "小李"。' }
    ],
    example: 'const course = { progress: 0 }\ncourse.progress = 20 // 合法\n\nlet count = 0\ncount += 1 // 合法',
    questions: [
      { text: '不需要重新赋值的变量应优先使用？', options: ['var', 'let', 'const'], answer: 2, explanation: 'const 能明确表达不会重新绑定。' },
      { text: 'const 声明的对象能修改属性吗？', options: ['可以', '不可以'], answer: 0, explanation: 'const 不会自动冻结对象本身。' }
    ],
    summary: ['默认 const，需要重赋值才用 let', '避免使用 var', 'const 对象的属性仍可修改']
  },
  {
    id: 'promise',
    title: '5 分钟理解 Promise 三种状态',
    competency: 'JavaScript 核心',
    duration: 5,
    difficulty: '基础',
    intro: 'Promise 就像外卖订单：处理中、已送达或配送失败，最终结果一旦确定就不会改变。',
    content: [
      { title: '三种状态', text: 'pending 表示进行中，fulfilled 表示成功，rejected 表示失败。Promise 只能从 pending 变为成功或失败。' },
      { title: '消费结果', text: 'then 处理成功结果，catch 处理错误，finally 无论成败都会执行，适合关闭 loading。' }
    ],
    example: 'fetchUser()\n  .then(user => render(user))\n  .catch(error => showError(error))\n  .finally(() => hideLoading())',
    questions: [
      { text: 'Promise 成功后的状态是？', options: ['pending', 'fulfilled', 'rejected'], answer: 1, explanation: 'fulfilled 表示操作成功完成。' },
      { text: '无论成功失败都执行的回调是？', options: ['then', 'catch', 'finally'], answer: 2, explanation: 'finally 常用于清理工作。' }
    ],
    summary: ['Promise 有 pending / fulfilled / rejected 三态', '状态确定后不可逆', 'then、catch、finally 分别处理结果、错误与清理']
  },
  {
    id: 'vue-reactivity',
    title: '5 分钟理解 Vue 响应式',
    competency: 'Vue 应用开发',
    duration: 5,
    difficulty: '基础',
    intro: '响应式就像数据和界面之间连了一根自动更新的线：数据变化，相关界面随之刷新。',
    content: [
      { title: 'ref 与 reactive', text: 'ref 适合单个值，通过 .value 在脚本中读写；reactive 适合对象。模板会自动解包 ref。' },
      { title: '派生状态', text: 'computed 根据响应式数据计算新值，并自动缓存。不要把可计算出的结果再存一份，避免状态不一致。' }
    ],
    example: 'const count = ref(0)\nconst doubled = computed(() => count.value * 2)',
    questions: [
      { text: '脚本中读取 ref 值需要使用？', options: ['.value', '.current', '.data'], answer: 0, explanation: 'ref 的值保存在 value 属性。' },
      { text: '派生状态应优先使用？', options: ['watchEffect', 'computed', 'setTimeout'], answer: 1, explanation: 'computed 专门用于声明式派生状态。' }
    ],
    summary: ['ref 管理单值，reactive 管理对象', '模板自动解包 ref', 'computed 用于派生状态']
  },
  {
    id: 'vue-component',
    title: '5 分钟设计 Vue 组件接口',
    competency: 'Vue 应用开发',
    duration: 5,
    difficulty: '基础',
    intro: '组件像一台小机器：props 是输入，emit 是输出，slot 是可替换的内容区域。',
    content: [
      { title: '单向数据流', text: '父组件通过 props 传入数据，子组件不应直接修改 props，而应通过 emit 通知父组件发生了什么。' },
      { title: '保持职责单一', text: '一个组件解决一个清晰问题。接口命名表达业务含义，例如 emit("submit", formData)，而不是暴露内部点击细节。' }
    ],
    example: 'const props = defineProps({ title: String })\nconst emit = defineEmits(["submit"])\nemit("submit", formData)',
    questions: [
      { text: '子组件通知父组件应使用？', options: ['props', 'emit', 'slot'], answer: 1, explanation: 'emit 用于向父级发出事件。' },
      { text: '子组件可以直接修改 props 吗？', options: ['推荐', '不推荐'], answer: 1, explanation: 'props 遵循单向数据流。' }
    ],
    summary: ['props 输入，emit 输出', '不要直接修改 props', '组件保持单一职责']
  },
  {
    id: 'git-basics',
    title: '5 分钟学会 Git 提交流程',
    competency: '工程化与实战',
    duration: 5,
    difficulty: '入门',
    intro: 'Git 提交像给项目拍快照：先挑选要入镜的改动，再写清这张快照做了什么。',
    content: [
      { title: '三步流程', text: 'git status 查看改动，git add 选择暂存内容，git commit 创建有说明的提交。提交前应检查 diff。' },
      { title: '小而清晰的提交', text: '一次提交只完成一个逻辑变化，消息使用动词描述结果，例如“添加目标设置流程”。' }
    ],
    example: 'git status\ngit add src/\ngit commit -m "feat: add learning goal setup"',
    questions: [
      { text: '查看当前改动状态使用？', options: ['git status', 'git push', 'git clone'], answer: 0, explanation: 'git status 展示工作区和暂存区状态。' },
      { text: '创建提交前选择文件使用？', options: ['git add', 'git log', 'git fetch'], answer: 0, explanation: 'git add 将内容加入暂存区。' }
    ],
    summary: ['先 status，再 add，最后 commit', '提交前检查差异', '每次提交保持单一且清晰']
  },
  {
    id: 'vite-build',
    title: '5 分钟完成 Vite 构建',
    competency: '工程化与实战',
    duration: 5,
    difficulty: '基础',
    intro: '开发服务器是工作台，生产构建则是把源码整理、压缩成可以交付上线的成品。',
    content: [
      { title: '构建与预览', text: 'npm run build 生成 dist 目录，npm run preview 在本地模拟生产服务。上线前应确保构建成功。' },
      { title: '环境变量', text: 'Vite 暴露给客户端的变量必须以 VITE_ 开头。密钥不能放在前端环境变量中，因为构建后用户可见。' }
    ],
    example: 'npm run build\nnpm run preview\n\nconst apiUrl = import.meta.env.VITE_API_URL',
    questions: [
      { text: 'Vite 默认生产输出目录是？', options: ['src', 'dist', 'public'], answer: 1, explanation: '构建产物默认写入 dist。' },
      { text: '敏感密钥能放在 VITE_ 变量中吗？', options: ['能', '不能'], answer: 1, explanation: '前端构建变量会暴露给浏览器。' }
    ],
    summary: ['build 生成 dist', 'preview 用于本地验收产物', '客户端环境变量不能存密钥']
  }
]

export const getCareer = (id) => careers.find((item) => item.id === id)
export const getMicroUnit = (id) => microUnits.find((item) => item.id === id)
