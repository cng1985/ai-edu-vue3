# AI 学习系统

基于 **Vue 3 + JavaScript（Vite）** 开发的目标驱动 IT 学习平台，围绕“目标确认 → 路径分解 → 学习执行 → 达成度评估”构建学习闭环。

## 功能特性

- **登录 / 注册**:社区风格的认证页面(左侧社区氛围面板 + 右侧表单卡片),支持头像选择、表单校验、路由守卫(未登录自动跳转登录页并支持 redirect 回跳)。演示版账号数据保存在浏览器本地
- **AI 职业规划**：职业推荐、目标条件确认、能力域/知识点/里程碑自动分解，内置前端工程师完整 MVP 路径
- **学习驾驶舱**：登录即见目标达成度、能力域进度、下一里程碑、今日待办与 AI 学习建议
- **5 分钟微学习**：8 个结构化微单元，包含引入、核心讲解、代码示例、即时快测、解析与要点回顾
- **达成度评估**：按“微单元完成度 + 快测成绩”计算知识点掌握度，再按能力权重聚合总目标达成度
- **成长激励**：每日打卡、连续学习、积分、等级与勋章，数据本地持久化
- **课程体系**：内置 3 门课程（提示词工程入门 / RAG 实战 / AI 驱动的应用开发实战），共 18 个章节，内容以 Markdown 管理，支持代码高亮与 Mermaid 图表渲染
- **章节学习**：章节目录导航、上一章/下一章、完成打卡、本章笔记（本地持久化）
- **AI 学习助手**：内置本地知识库检索引擎（"检索 → 组装 → 流式生成"三段式，与真实 RAG 架构一致），支持流式输出、停止生成与知识来源溯源跳转
- **知识测验**：每门课程配套测验题，提交后逐题判分并展示解析
- **学习统计**：总体/分课程进度、近 7 天学习活跃度图表、测验成绩汇总与笔记一览
- **数据持久化**:学习进度、笔记、测验成绩与对话记录均保存在浏览器 localStorage 中

## 管理后台

项目包含独立的 **PC 管理端** 与 **API 服务**，用于平台运营与内容管理。

### 功能模块

| 模块 | 说明 |
| :--- | :--- |
| 数据看板 | 用户、课程、题库、审核队列概览统计 |
| 用户管理 | 学员/管理员/审核员/运营的增删改查与角色权限 |
| 课程管理 | 课程与章节的 CRUD，支持 Markdown 编辑与预览 |
| 题库管理 | 测验套题与单题的增删改查 |
| 内容审核 | AI 预审 + 人工终审工作流（通过/驳回） |

### 启动方式

```bash
# 安装依赖（应用端 + 管理端）
npm run install:all

# 终端 1：启动 Go API 服务（http://localhost:3001）
npm run dev:server

# 终端 2：启动管理端（http://localhost:5174）
npm run dev:admin

# 终端 3：启动学习者应用端（http://localhost:5173）
npm run dev
```

### 后端技术栈

| 类别 | 选型 |
| :--- | :--- |
| 语言 | Go 1.22 |
| 依赖注入 | Uber Fx |
| Web 框架 | Gin |
| ORM | GORM |
| 数据库 | SQLite（开发环境，可切换 MySQL） |
| 认证 | JWT + bcrypt |

### 演示账号

| 角色 | 用户名 | 密码 |
| :--- | :--- | :--- |
| 管理员 | admin | admin123 |
| 审核员 | reviewer | review123 |

### 目录结构

```
server/                Go API 服务（Fx + Gin + GORM）
├── cmd/server/        入口
├── internal/
│   ├── config/        配置
│   ├── database/      GORM 初始化与迁移
│   ├── model/         数据模型
│   ├── repository/    数据访问层
│   ├── service/       业务逻辑层
│   ├── handler/       HTTP Handler
│   ├── middleware/    JWT 认证与权限
│   ├── router/        路由注册
│   └── seed/          种子数据
├── pkg/               公共工具（JWT、响应封装）
└── data/              SQLite 数据库（运行时生成）

admin/                 Vue3 + Element Plus 管理端
├── src/views/         登录、看板、用户、课程、题库、审核
├── src/layouts/       后台布局
└── src/api/           API 封装
```

## 技术栈

| 类别 | 选型 |
| :--- | :--- |
| 框架 | Vue 3（Composition API + `<script setup>`） |
| 构建 | Vite 6 |
| 路由 | Vue Router 4（Hash 模式） |
| 状态管理 | Pinia 3 |
| Markdown 渲染 | marked |
| 代码高亮 | highlight.js |
| 图表渲染 | mermaid（按需异步加载） |

## 快速开始

```bash
# 安装依赖
npm install

# 启动开发服务器（默认 http://localhost:5173）
npm run dev

# 生产构建
npm run build

# 预览构建产物
npm run preview
```

## 目录结构

```
src/
├── assets/            全局样式
├── components/        通用组件（侧边栏、Markdown 渲染器、进度环、课程卡片）
├── data/
│   ├── content/       课程正文（Markdown，经 ?raw 导入）
│   ├── courses.js     课程目录数据
│   ├── careerPath.js  职业、能力图谱与微单元数据
│   └── quizzes.js     测验题库
├── router/            路由配置（含登录守卫）
├── stores/            Pinia 状态（账号认证 / 目标成长 / 学习进度 / 对话）
├── utils/
│   ├── aiEngine.js    AI 助手引擎（本地检索 + 流式输出）
│   └── markdown.js    Markdown/高亮/Mermaid 渲染工具
└── views/             页面（登录、注册、首页、课程、章节、AI 助手、测验、统计）
```

## 接入真实大模型

`src/utils/aiEngine.js` 中的 `ask(question, handlers)` 是 AI 助手的唯一入口，当前实现为本地知识库检索 + 模拟流式输出。若要接入真实 LLM API（如 OpenAI 兼容接口的 SSE 流式输出），只需替换该函数内部实现，保持 `onToken` / `onDone` 回调协议不变即可，界面层无需任何改动。
