package seed

import "github.com/cng1985/ai-learning-server/internal/model"

type skillDef struct {
	ID, Name, Category, Desc string
	Tags                     []string
}

var skillDefs = []skillDef{
	{"java-core", "Java 基础", "后端开发", "面向对象、集合框架与常用 API，能编写规范的 Java 代码。", []string{"java", "oop", "集合"}},
	{"java-concurrency", "Java 并发编程", "后端开发", "线程模型、JMM、锁与 JUC 工具，能编写正确高效的并发代码。", []string{"并发", "线程", "juc", "锁", "volatile", "synchronized"}},
	{"spring-boot", "Spring Boot", "后端开发", "基于 Spring Boot 构建 RESTful 服务、管理事务与配置。", []string{"spring", "springboot", "接口", "rest", "controller"}},
	{"mysql", "MySQL 与数据库优化", "数据存储", "表结构设计、索引、事务与慢查询优化。", []string{"mysql", "sql", "索引", "事务", "表"}},
	{"redis", "Redis 应用", "数据存储", "数据结构、缓存模式、缓存一致性与分布式锁。", []string{"redis", "缓存", "lua"}},
	{"microservice", "微服务设计", "架构设计", "服务拆分、服务治理与接口契约。", []string{"微服务", "服务治理", "网关", "拆分"}},
	{"distributed", "分布式系统", "架构设计", "消息队列、分布式事务与一致性方案。", []string{"分布式", "一致性", "消息队列", "mq", "幂等"}},
	{"performance", "性能优化", "架构设计", "JVM 调优、压测与全链路性能分析。", []string{"性能", "压测", "jvm", "调优", "qps"}},
	{"system-design", "系统设计", "架构设计", "需求抽象、容量评估、高可用与技术方案设计。", []string{"架构", "高可用", "方案", "容量"}},
	{"prompt-eng", "提示词工程", "AI 应用", "结构化提示词、Few-shot、思维链与注入防御。", []string{"prompt", "提示词", "few-shot"}},
	{"rag", "RAG 检索增强", "AI 应用", "文档切分、向量化、混合检索与重排序。", []string{"rag", "检索", "向量", "embedding", "切分", "rerank"}},
	{"ai-app-dev", "AI 应用开发", "AI 应用", "Function Calling、Agent 设计、评估与生产级容错。", []string{"agent", "function calling", "llm", "大模型", "评估"}},
	{"python", "Python 编程", "数据与 AI", "Python 语法、数据处理与脚本开发。", []string{"python", "pandas"}},
	{"data-pipeline", "数据管道", "数据与 AI", "ETL、数据建模与数据仓库。", []string{"etl", "数据仓库", "spark", "调度"}},
	{"html-css", "HTML / CSS", "前端开发", "语义化结构、布局与响应式。", []string{"html", "css", "flex", "布局"}},
	{"javascript", "JavaScript", "前端开发", "语言核心、异步编程与模块化。", []string{"javascript", "js", "es6", "promise", "async"}},
	{"vue", "Vue 应用开发", "前端开发", "响应式原理、组件设计与工程化。", []string{"vue", "组件", "响应式", "pinia"}},
	{"requirement", "需求分析", "产品设计", "用户研究、需求拆解与 PRD 撰写。", []string{"需求", "prd", "用户", "场景"}},
	{"ai-product", "AI 产品设计", "产品设计", "AI 场景识别、能力边界与效果评估。", []string{"ai产品", "场景", "评估", "体验"}},
}

type kpDef struct {
	ID, Name, Domain string
	Difficulty       int
	Minutes          int
	Summary, Content string
	Skills           []string
	Questions        []model.Question
}

func q(text string, options []string, answer int, explanation string) model.Question {
	return model.Question{Text: text, Options: options, Answer: answer, Explanation: explanation}
}

var kpDefs = []kpDef{
	// ---------- Java 基础 ----------
	{"java-oop", "面向对象与接口设计", "Java", 1, 25, "封装、继承、多态与面向接口编程。",
		"## 核心要点\n\n- **封装**：隐藏实现细节，只暴露必要的行为\n- **多态**：同一接口不同实现，调用方依赖抽象\n- **面向接口编程**：依赖倒置，便于替换与测试\n\n```java\nList<String> list = new ArrayList<>(); // 依赖接口而非实现\n```",
		[]string{"java-core"}, []model.Question{
			q("\"面向接口编程\"最主要的好处是？", []string{"代码更短", "调用方与实现解耦，便于替换与测试", "运行更快", "减少内存占用"}, 1, "依赖抽象让实现可替换，是可测试性与扩展性的基础。"),
			q("下列哪项体现了多态？", []string{"private 字段", "父类引用指向子类对象并调用重写方法", "static 方法", "final 类"}, 1, "运行时根据实际类型分派方法调用即多态。"),
		}},
	{"java-collections", "集合框架", "Java", 2, 30, "List/Set/Map 的实现差异与选型。",
		"## 选型速查\n\n| 场景 | 推荐 |\n| --- | --- |\n| 随机访问 | ArrayList |\n| 去重 | HashSet |\n| 键值查找 | HashMap |\n| 有序 | TreeMap |\n\nHashMap 在 JDK8 中采用 **数组 + 链表 + 红黑树**，链表长度超过 8 时树化。",
		[]string{"java-core"}, []model.Question{
			q("JDK8 HashMap 链表长度超过多少会尝试树化？", []string{"4", "8", "16", "64"}, 1, "链表长度 ≥ 8 且数组长度 ≥ 64 时转为红黑树。"),
			q("需要按键有序遍历时应选择？", []string{"HashMap", "TreeMap", "ArrayList", "HashSet"}, 1, "TreeMap 基于红黑树，按键排序。"),
		}},
	// ---------- Java 并发 ----------
	{"thread", "Thread 线程基础", "Java 并发", 1, 20, "线程的创建、状态流转与中断机制。",
		"## 线程状态\n\nNEW → RUNNABLE → BLOCKED / WAITING / TIMED_WAITING → TERMINATED\n\n- `start()` 才会创建新线程，`run()` 只是普通方法调用\n- 中断是**协作式**的：`interrupt()` 只设置标志位",
		[]string{"java-concurrency"}, []model.Question{
			q("直接调用 thread.run() 会发生什么？", []string{"启动新线程", "在当前线程中执行 run 方法", "抛出异常", "线程进入 WAITING"}, 1, "只有 start() 会创建新线程。"),
			q("Java 线程中断机制的特点是？", []string{"强制终止线程", "协作式，设置中断标志由线程自行响应", "只能中断守护线程", "会释放所有锁"}, 1, "interrupt() 设置标志，阻塞方法会抛出 InterruptedException。"),
		}},
	{"jmm", "Java 内存模型 JMM", "Java 并发", 3, 35, "主内存与工作内存、happens-before 规则。",
		"## 三大问题\n\n- **可见性**：一个线程的修改其他线程能否看到\n- **有序性**：指令重排序\n- **原子性**：操作是否可分割\n\nJMM 通过 **happens-before** 规则定义可见性保证，例如：锁的释放 happens-before 后续加锁。",
		[]string{"java-concurrency"}, []model.Question{
			q("JMM 主要解决哪三类问题？", []string{"内存泄漏、GC、OOM", "可见性、有序性、原子性", "加载、链接、初始化", "编译、解释、执行"}, 1, "并发三大问题即可见性、有序性、原子性。"),
			q("以下哪条属于 happens-before 规则？", []string{"线程启动早于 main 方法", "对一个锁的解锁 happens-before 随后对这个锁的加锁", "所有写操作立即对其他线程可见", "volatile 变量写操作可以被重排序到读之后"}, 1, "监视器锁规则是 happens-before 的核心规则之一。"),
		}},
	{"volatile", "volatile 关键字", "Java 并发", 3, 25, "保证可见性与禁止重排序，但不保证复合操作原子性。",
		"## volatile 能做什么\n\n- ✅ 写操作立即刷新到主内存，读操作从主内存读取（**可见性**）\n- ✅ 通过内存屏障**禁止指令重排序**（DCL 单例必须加 volatile）\n- ❌ **不保证原子性**：`count++` 仍然线程不安全\n\n```java\nprivate volatile boolean running = true;\n```",
		[]string{"java-concurrency"}, []model.Question{
			q("volatile 能保证 count++ 的线程安全吗？", []string{"能", "不能，count++ 是读-改-写复合操作", "只在单核 CPU 上能", "只在 JDK17 之后能"}, 1, "volatile 不保证复合操作的原子性，应使用 AtomicInteger 或加锁。"),
			q("DCL 双重检查单例中 instance 为什么要加 volatile？", []string{"提高性能", "防止对象初始化指令重排序导致拿到未初始化对象", "保证 instance 不被 GC", "让构造方法同步执行"}, 1, "new 对象可能被重排为先赋值引用再初始化，volatile 禁止该重排序。"),
		}},
	{"thread-safety", "线程安全", "Java 并发", 2, 25, "竞态条件、临界区与线程安全的实现方式。",
		"## 实现线程安全的思路\n\n1. **不共享**：线程封闭、ThreadLocal\n2. **不可变**：final + 不可变对象\n3. **同步**：synchronized、Lock\n4. **无锁**：CAS 原子类\n\n竞态条件的典型形式：**检查再执行（check-then-act）**。",
		[]string{"java-concurrency"}, []model.Question{
			q("\"if (map.get(k) == null) map.put(k, v)\" 在并发下的问题属于？", []string{"死锁", "检查再执行的竞态条件", "内存泄漏", "活锁"}, 1, "检查与执行之间可能被其他线程插入，应使用 putIfAbsent。"),
			q("以下哪种方式不需要加锁也能保证线程安全？", []string{"共享可变 HashMap", "不可变对象", "普通 ArrayList", "static 可变字段"}, 1, "不可变对象天然线程安全。"),
		}},
	{"synchronized", "synchronized 与锁升级", "Java 并发", 3, 30, "对象监视器、可重入与锁升级过程。",
		"## 要点\n\n- 修饰实例方法锁 `this`，修饰静态方法锁 `Class` 对象\n- **可重入**：同一线程可重复获取\n- 锁升级：**无锁 → 偏向锁 → 轻量级锁 → 重量级锁**（JDK15 后偏向锁默认禁用）\n- 同时保证原子性、可见性与有序性",
		[]string{"java-concurrency"}, []model.Question{
			q("synchronized 修饰静态方法时锁的是？", []string{"当前实例 this", "类的 Class 对象", "方法本身", "所有实例"}, 1, "静态同步方法锁的是 Class 对象。"),
			q("关于 synchronized 正确的是？", []string{"不可重入", "可重入，且能保证可见性", "只能修饰方法", "只保证可见性不保证原子性"}, 1, "synchronized 可重入，并同时保证原子性与可见性。"),
		}},
	{"cas", "CAS 与原子类", "Java 并发", 3, 25, "比较并交换的无锁算法及 ABA 问题。",
		"## CAS\n\n`compareAndSwap(内存值, 期望值, 新值)`：只有内存值等于期望值时才更新。\n\n- 优点：无锁，避免线程切换\n- 问题：**ABA**（用 `AtomicStampedReference` 加版本号）、自旋开销、只能保证单变量",
		[]string{"java-concurrency"}, []model.Question{
			q("解决 CAS 的 ABA 问题常用？", []string{"synchronized", "AtomicStampedReference 版本号", "volatile", "ThreadLocal"}, 1, "为值附加版本号，值相同但版本不同即可识别 ABA。"),
			q("CAS 适合的场景是？", []string{"长时间持有的复杂临界区", "竞争不激烈的单变量原子更新", "跨 JVM 的分布式锁", "IO 密集操作"}, 1, "高竞争下自旋开销大，CAS 适合低竞争单变量。"),
		}},
	{"aqs", "AQS 抽象队列同步器", "Java 并发", 4, 40, "JUC 锁与同步器的基础框架。",
		"## AQS 结构\n\n- `volatile int state`：同步状态\n- **CLH 双向队列**：获取失败的线程排队阻塞\n- 独占模式（ReentrantLock）与共享模式（Semaphore、CountDownLatch）\n\n通过 **CAS 修改 state** + 队列 + `LockSupport.park/unpark` 实现。",
		[]string{"java-concurrency"}, []model.Question{
			q("AQS 用什么记录同步状态？", []string{"synchronized 块", "volatile int state", "ThreadLocal", "HashMap"}, 1, "state 由 CAS 修改，volatile 保证可见性。"),
			q("CountDownLatch 使用 AQS 的哪种模式？", []string{"独占模式", "共享模式", "条件模式", "公平模式"}, 1, "多个线程可同时被唤醒，属于共享模式。"),
		}},
	{"thread-pool", "ThreadPoolExecutor 线程池", "Java 并发", 4, 40, "核心参数、任务调度流程与生产配置。",
		"## 七大参数\n\ncorePoolSize、maximumPoolSize、keepAliveTime、unit、workQueue、threadFactory、handler\n\n## 提交流程\n\n核心线程未满 → 创建核心线程；否则入队；队列满 → 创建非核心线程；达到最大线程 → **拒绝策略**。\n\n> 生产环境禁止使用无界队列的 `Executors.newFixedThreadPool`，以免 OOM。",
		[]string{"java-concurrency", "performance"}, []model.Question{
			q("核心线程已满时，新任务会先？", []string{"直接创建非核心线程", "进入工作队列", "执行拒绝策略", "阻塞提交线程"}, 1, "核心线程满后先入队，队列满才创建非核心线程。"),
			q("为什么不推荐 Executors.newFixedThreadPool？", []string{"线程数不可配", "使用无界队列，任务堆积可能 OOM", "不支持 Callable", "无法关闭"}, 1, "LinkedBlockingQueue 默认无界，可能导致内存溢出。"),
		}},
	// ---------- Spring Boot ----------
	{"spring-ioc", "IoC 与依赖注入", "Spring", 2, 25, "容器管理 Bean 生命周期与依赖注入。",
		"## 要点\n\n- IoC：对象创建与依赖关系交给容器\n- 推荐**构造器注入**：依赖不可变、便于测试\n- Bean 默认单例，注意不要在单例中保存请求状态",
		[]string{"spring-boot"}, []model.Question{
			q("Spring 官方推荐的注入方式是？", []string{"字段注入", "构造器注入", "静态注入", "setter 注入"}, 1, "构造器注入保证依赖不可变且便于单元测试。"),
			q("Spring Bean 默认作用域是？", []string{"prototype", "singleton", "request", "session"}, 1, "默认 singleton，需注意线程安全。"),
		}},
	{"spring-rest", "RESTful 接口开发", "Spring", 2, 30, "Controller、参数校验、统一响应与异常处理。",
		"## 规范\n\n- 资源用名词：`GET /orders/{id}`\n- 参数校验：`@Valid` + Bean Validation\n- 统一异常处理：`@RestControllerAdvice`\n- 写接口要考虑**幂等**（下单使用幂等号）",
		[]string{"spring-boot"}, []model.Question{
			q("统一处理 Controller 异常常用？", []string{"@Service", "@RestControllerAdvice", "@Bean", "@Async"}, 1, "@RestControllerAdvice + @ExceptionHandler 统一转换错误响应。"),
			q("创建订单接口防止重复提交的关键是？", []string{"加大线程池", "接口幂等（如幂等号/唯一约束）", "返回 200", "使用 GET 请求"}, 1, "幂等设计保证重复请求只生效一次。"),
		}},
	{"spring-transaction", "Spring 事务管理", "Spring", 3, 30, "@Transactional 传播行为与失效场景。",
		"## 常见失效场景\n\n- 同类内部方法调用（绕过代理）\n- 方法非 public\n- 异常被 catch 吞掉，或抛出受检异常未配置 rollbackFor\n\n传播行为默认 `REQUIRED`。",
		[]string{"spring-boot", "mysql"}, []model.Question{
			q("同一个类中 A 方法调用带 @Transactional 的 B 方法，事务会？", []string{"正常生效", "失效，因为绕过了代理对象", "抛出异常", "变成只读事务"}, 1, "Spring 事务基于代理，内部调用不经过代理。"),
			q("@Transactional 默认对哪类异常回滚？", []string{"所有异常", "RuntimeException 与 Error", "仅 IOException", "不回滚"}, 1, "默认只回滚非受检异常，受检异常需配置 rollbackFor。"),
		}},
	// ---------- MySQL ----------
	{"sql-basics", "SQL 基础与表设计", "MySQL", 1, 25, "范式、主键选择与常用 SQL。",
		"## 表设计原则\n\n- 主键使用自增或趋势递增 ID，避免页分裂\n- 金额使用 `DECIMAL`\n- 适度反范式：订单快照冗余商品名称与价格",
		[]string{"mysql"}, []model.Question{
			q("订单金额字段推荐使用？", []string{"FLOAT", "DOUBLE", "DECIMAL", "VARCHAR"}, 2, "浮点数有精度误差，金额必须使用 DECIMAL。"),
			q("订单表冗余商品名称与下单价格的原因是？", []string{"节省空间", "保留下单时快照，避免商品变更影响历史订单", "提高写入速度", "满足第三范式"}, 1, "订单是历史事实，需要快照。"),
		}},
	{"mysql-index", "索引原理与 B+ 树", "MySQL", 3, 35, "聚簇索引、联合索引与最左前缀。",
		"## 要点\n\n- InnoDB 使用 **B+ 树**，叶子节点有序链表便于范围查询\n- 聚簇索引存整行，二级索引存主键 → **回表**\n- 联合索引遵循**最左前缀**；**覆盖索引**可避免回表",
		[]string{"mysql"}, []model.Question{
			q("联合索引 (a,b,c) 能用于以下哪个查询？", []string{"WHERE b=1 AND c=2", "WHERE a=1 AND b=2", "WHERE c=3", "WHERE b>1"}, 1, "最左前缀原则，必须从 a 开始。"),
			q("覆盖索引的好处是？", []string{"减少磁盘空间", "查询列都在索引中，避免回表", "加快写入", "自动分区"}, 1, "二级索引包含所需列时无需回表。"),
		}},
	{"mysql-transaction", "事务与隔离级别", "MySQL", 3, 30, "ACID、MVCC 与四种隔离级别。",
		"## 隔离级别\n\n| 级别 | 脏读 | 不可重复读 | 幻读 |\n| --- | --- | --- | --- |\n| RU | ✔ | ✔ | ✔ |\n| RC | ✘ | ✔ | ✔ |\n| RR（默认） | ✘ | ✘ | 部分避免 |\n| Serializable | ✘ | ✘ | ✘ |\n\nInnoDB 通过 **MVCC + Next-Key Lock** 实现 RR。",
		[]string{"mysql"}, []model.Question{
			q("InnoDB 默认隔离级别是？", []string{"READ UNCOMMITTED", "READ COMMITTED", "REPEATABLE READ", "SERIALIZABLE"}, 2, "MySQL InnoDB 默认 RR。"),
			q("MVCC 主要依赖什么实现？", []string{"表锁", "undo log 版本链 + ReadView", "binlog", "redo log"}, 1, "通过 undo 版本链和 ReadView 判断可见性。"),
		}},
	{"slow-query", "慢查询分析与优化", "MySQL", 4, 40, "EXPLAIN、索引失效与 SQL 改写。",
		"## 排查步骤\n\n1. 开启慢查询日志定位 SQL\n2. `EXPLAIN` 看 type / key / rows / Extra\n3. 常见失效：函数作用于索引列、隐式类型转换、`LIKE '%x'`\n4. 深分页改为**游标分页**：`WHERE id > ? LIMIT 20`",
		[]string{"mysql", "performance"}, []model.Question{
			q("EXPLAIN 中 type=ALL 表示？", []string{"使用了覆盖索引", "全表扫描", "唯一索引等值查询", "范围扫描"}, 1, "ALL 为全表扫描，通常需要优化。"),
			q("LIMIT 1000000, 20 深分页的推荐优化是？", []string{"加大 buffer pool", "基于上一页最大 id 的游标分页", "改用 MyISAM", "增加 LIMIT"}, 1, "游标分页避免扫描并丢弃大量行。"),
		}},
	// ---------- Redis ----------
	{"redis-datatypes", "Redis 数据结构", "Redis", 2, 25, "String/Hash/List/Set/ZSet 的典型场景。",
		"## 场景\n\n- String：计数、缓存对象\n- Hash：对象字段\n- ZSet：排行榜、延时队列\n- 原子性：单命令原子，多命令用 **Lua 脚本**",
		[]string{"redis"}, []model.Question{
			q("实现排行榜最合适的数据结构是？", []string{"List", "ZSet", "String", "Hash"}, 1, "ZSet 按分数有序。"),
			q("需要多个 Redis 命令原子执行时用？", []string{"多次调用", "Lua 脚本", "pipeline 一定原子", "KEYS 命令"}, 1, "Lua 脚本在 Redis 中原子执行。"),
		}},
	{"cache-patterns", "缓存模式与一致性", "Redis", 3, 35, "Cache Aside、穿透/击穿/雪崩与一致性方案。",
		"## Cache Aside\n\n读：先缓存，未命中查库回填；写：**先更新数据库，再删除缓存**。\n\n| 问题 | 方案 |\n| --- | --- |\n| 穿透 | 缓存空值 / 布隆过滤器 |\n| 击穿 | 互斥锁 / 逻辑过期 |\n| 雪崩 | 过期时间加随机值 |",
		[]string{"redis"}, []model.Question{
			q("Cache Aside 写操作推荐顺序？", []string{"先删缓存再更新库", "先更新数据库再删除缓存", "只更新缓存", "先更新缓存再更新库"}, 1, "先更新库再删缓存，不一致窗口最小。"),
			q("大量 key 同时过期导致数据库压力激增称为？", []string{"缓存穿透", "缓存雪崩", "缓存击穿", "缓存预热"}, 1, "雪崩可通过过期时间随机化缓解。"),
		}},
	{"distributed-lock", "分布式锁", "Redis", 4, 35, "SET NX PX、锁续期与 Redisson。",
		"## 正确姿势\n\n```\nSET lock:order:1 <uuid> NX PX 30000\n```\n\n- 释放锁要校验持有者（Lua 脚本 compare-and-delete）\n- 业务超时用**看门狗续期**（Redisson）\n- 库存扣减也可直接用 Lua 原子扣减替代锁",
		[]string{"redis", "distributed"}, []model.Question{
			q("释放 Redis 分布式锁时为什么要校验 value？", []string{"提高速度", "防止误删其他线程持有的锁", "节省内存", "Redis 要求"}, 1, "锁可能已过期被他人获取，需校验持有者。"),
			q("Redisson 看门狗的作用是？", []string{"监控 Redis 内存", "在业务未完成时自动续期锁", "自动重试命令", "主从切换"}, 1, "看门狗定期延长锁过期时间。"),
		}},
	// ---------- 分布式与微服务 ----------
	{"service-split", "服务拆分与治理", "分布式", 3, 35, "按业务边界拆分服务，注册发现、熔断限流。",
		"## 拆分原则\n\n- 按**领域边界**（DDD 限界上下文）拆分\n- 服务自治：独立数据库\n- 治理：注册发现、负载均衡、**熔断降级**、限流",
		[]string{"microservice"}, []model.Question{
			q("微服务拆分的首要依据是？", []string{"代码行数", "业务领域边界", "团队人数", "数据库表数量"}, 1, "按限界上下文拆分，保证高内聚低耦合。"),
			q("下游服务故障时防止级联雪崩的手段是？", []string{"增加重试次数", "熔断降级", "关闭日志", "扩大超时"}, 1, "熔断快速失败，保护调用方。"),
		}},
	{"message-queue", "消息队列", "分布式", 3, 30, "削峰解耦、可靠投递与幂等消费。",
		"## 关键点\n\n- 场景：异步解耦、削峰填谷\n- 可靠性：生产者确认 + 持久化 + 消费者手动 ACK\n- 消息可能重复投递 → **消费端必须幂等**",
		[]string{"distributed"}, []model.Question{
			q("消息队列中为什么消费端必须幂等？", []string{"提高吞吐", "消息可能重复投递", "节省存储", "保证顺序"}, 1, "至少一次投递语义下可能重复。"),
			q("秒杀场景引入 MQ 的主要目的是？", []string{"增加功能", "削峰填谷", "替代数据库", "加密数据"}, 1, "MQ 缓冲瞬时流量。"),
		}},
	{"distributed-transaction", "分布式事务", "分布式", 4, 40, "2PC、TCC、本地消息表与最终一致性。",
		"## 方案对比\n\n- **2PC**：强一致，阻塞、性能差\n- **TCC**：Try/Confirm/Cancel，侵入业务\n- **本地消息表 / 事务消息**：最终一致，互联网主流\n- Saga：长事务补偿",
		[]string{"distributed", "microservice"}, []model.Question{
			q("互联网下单扣库存最常用的一致性方案是？", []string{"2PC 强一致", "基于事务消息的最终一致性", "不处理", "单库事务"}, 1, "最终一致兼顾性能与可靠性。"),
			q("TCC 中 Cancel 阶段的作用是？", []string{"提交资源", "释放 Try 阶段预留的资源", "重试 Try", "记录日志"}, 1, "Cancel 用于补偿回滚预留资源。"),
		}},
	// ---------- 性能与系统设计 ----------
	{"jvm-gc", "JVM 内存与 GC", "性能", 4, 40, "内存区域、GC 算法与 G1 调优。",
		"## 要点\n\n- 堆：新生代（Eden/S0/S1）+ 老年代\n- G1 以 Region 为单位，目标停顿时间可配（`-XX:MaxGCPauseMillis`）\n- 排查：GC 日志、`jstat`、堆转储分析内存泄漏",
		[]string{"performance"}, []model.Question{
			q("G1 收集器的主要特点是？", []string{"只有单线程", "基于 Region、可预测停顿时间", "不分代", "只回收老年代"}, 1, "G1 将堆划分为 Region，按收益优先回收。"),
			q("频繁 Full GC 首先应该？", []string{"重启服务", "分析 GC 日志与堆转储定位对象", "增加 CPU", "关闭 GC"}, 1, "先定位是否存在内存泄漏或大对象。"),
		}},
	{"load-testing", "压测与性能分析", "性能", 3, 35, "QPS、RT、瓶颈定位与容量评估。",
		"## 方法\n\n1. 明确目标：QPS、P99 RT、错误率\n2. 阶梯加压，观察拐点\n3. 定位瓶颈：CPU / IO / 锁 / 连接池\n4. 优化后**回归压测**对比",
		[]string{"performance"}, []model.Question{
			q("衡量接口尾延迟最常用的指标是？", []string{"平均 RT", "P99 RT", "最小 RT", "QPS"}, 1, "P99 反映长尾用户体验。"),
			q("性能优化完成后应该？", []string{"直接上线", "回归压测对比优化前后指标", "删除监控", "降低日志级别即可"}, 1, "用数据证明优化效果。"),
		}},
	{"high-availability", "高可用设计", "系统设计", 4, 40, "冗余、隔离、降级与故障演练。",
		"## 高可用手段\n\n- 冗余：多副本、多机房\n- 隔离：线程池隔离、舱壁\n- 降级与限流：保核心链路\n- 可观测：监控、告警、链路追踪\n- 演练：混沌工程",
		[]string{"system-design"}, []model.Question{
			q("\"保证核心链路可用，暂时关闭非核心功能\"属于？", []string{"扩容", "降级", "分库", "缓存"}, 1, "降级牺牲次要功能保障核心。"),
			q("以下哪项不是提升可用性的手段？", []string{"多副本", "限流", "单点部署", "故障演练"}, 2, "单点部署恰恰降低可用性。"),
		}},
	{"system-design-method", "系统设计方法论", "系统设计", 4, 45, "需求澄清 → 容量估算 → 架构 → 深挖 → 权衡。",
		"## 五步法\n\n1. **澄清需求**：功能与非功能（QPS、延迟、一致性）\n2. **容量估算**：存储、带宽、QPS\n3. **高层架构**：接入、服务、存储、缓存、MQ\n4. **深挖瓶颈**：热点、一致性、扩展性\n5. **权衡取舍**：说明为什么这样设计",
		[]string{"system-design", "microservice"}, []model.Question{
			q("系统设计的第一步应该是？", []string{"画架构图", "澄清功能与非功能需求", "选数据库", "写代码"}, 1, "需求不清，设计无从谈起。"),
			q("技术方案中\"权衡取舍\"的意义是？", []string{"展示知识面", "说明选择的理由与代价，便于决策", "增加篇幅", "可以省略"}, 1, "没有银弹，方案需要说明代价。"),
		}},
	// ---------- AI 应用 ----------
	{"llm-basics", "大模型基础", "AI 应用", 1, 20, "Token、上下文窗口、temperature 与幻觉。",
		"## 关键概念\n\n- **Token**：模型处理的最小单位，影响成本与上下文\n- **上下文窗口**：一次可处理的 token 上限\n- **temperature**：越高越随机\n- **幻觉**：模型生成看似合理但错误的内容 → 用 RAG 与约束输出缓解",
		[]string{"prompt-eng", "ai-app-dev"}, []model.Question{
			q("需要稳定、可复现的输出时 temperature 应该？", []string{"调高", "调低", "设为随机", "无影响"}, 1, "低 temperature 输出更确定。"),
			q("缓解大模型幻觉的有效手段是？", []string{"增加 temperature", "引入检索增强（RAG）提供事实依据", "缩短提示词", "使用更多 emoji"}, 1, "RAG 让模型基于可靠资料作答。"),
		}},
	{"prompt-structure", "结构化提示词", "AI 应用", 2, 25, "角色、任务、上下文、输出约束四要素。",
		"## 四要素\n\n1. 角色设定\n2. 任务描述\n3. 上下文信息（用 XML 标签隔离）\n4. 输出约束（格式、长度、JSON Schema）\n\n关键指令放在**开头和结尾**，对抗中间迷失。",
		[]string{"prompt-eng"}, []model.Question{
			q("用 XML 标签包裹用户输入的主要目的是？", []string{"美观", "隔离不可信内容，降低注入与歧义", "减少 token", "强制输出 XML"}, 1, "标签划分语义区域，防御提示词注入。"),
			q("长提示词中关键指令建议放在？", []string{"中间", "开头和结尾", "随机位置", "单独一条消息"}, 1, "模型对首尾更敏感。"),
		}},
	{"few-shot", "Few-shot 与思维链", "AI 应用", 3, 30, "示例驱动输出格式，CoT 提升推理。",
		"## 要点\n\n- Few-shot：给出 \"输入 → 期望输出\" 样例，**示例优于描述**\n- 思维链（CoT）：让模型分步推理，适合数学与逻辑任务\n- 示例要覆盖边界情况，避免偏差",
		[]string{"prompt-eng"}, []model.Question{
			q("Few-shot 示例的核心价值是？", []string{"增加字数", "让模型模仿示例的格式与风格", "降低成本", "替代系统提示词"}, 1, "示例比冗长描述更有效。"),
			q("思维链提示最适合？", []string{"简单分类", "多步推理问题", "图片生成", "翻译单词"}, 1, "CoT 通过显式推理步骤提高复杂问题准确率。"),
		}},
	{"embedding", "Embedding 向量化", "AI 应用", 2, 25, "语义向量、相似度与模型选择。",
		"## 要点\n\n- 文本 → 高维向量，语义相近则距离近\n- 相似度：余弦相似度\n- **更换 Embedding 模型必须全量重建索引**（向量空间不兼容）",
		[]string{"rag"}, []model.Question{
			q("更换 Embedding 模型后需要？", []string{"无需处理", "全量重建向量索引", "只重建新增数据", "调整维度即可"}, 1, "不同模型的向量空间互不兼容。"),
			q("衡量两个向量语义相似度常用？", []string{"欧氏距离的平方根之和", "余弦相似度", "字符串长度", "哈希值"}, 1, "余弦相似度与向量长度无关。"),
		}},
	{"chunking", "文档切分策略", "AI 应用", 3, 30, "固定长度、结构感知与重叠窗口。",
		"## 策略\n\n- 固定长度 + 重叠：简单但可能截断语义\n- **结构感知切分**：按标题、段落，推荐默认\n- Chunk 过大：噪声多；过小：上下文缺失",
		[]string{"rag"}, []model.Question{
			q("生产环境推荐的默认切分策略是？", []string{"不切分", "结构感知切分", "每句一块", "随机切分"}, 1, "结构感知切分保留语义完整性。"),
			q("Chunk 过小的主要问题是？", []string{"存储太大", "上下文缺失，回答不完整", "检索太慢", "无法向量化"}, 1, "过小的块缺乏足够上下文。"),
		}},
	{"hybrid-retrieval", "混合检索与重排序", "AI 应用", 4, 40, "向量 + BM25 召回，交叉编码器精排。",
		"## 两阶段检索\n\n1. 召回：稠密向量（语义）+ 稀疏 BM25（关键词），Top-20\n2. 精排：**交叉编码器 Rerank**，Top-3\n\n多轮对话需要**查询重写**补全指代。",
		[]string{"rag"}, []model.Question{
			q("混合检索结合了？", []string{"两个大模型", "向量检索与关键词检索", "图片与文本", "缓存与数据库"}, 1, "语义与精确匹配互补。"),
			q("Rerank 阶段通常使用？", []string{"双编码器", "交叉编码器", "决策树", "正则表达式"}, 1, "交叉编码器逐对打分精度更高。"),
		}},
	{"function-calling", "Function Calling", "AI 应用", 3, 30, "模型输出调用意图，应用层执行工具。",
		"## 流程\n\n1. 声明工具（名称、描述、JSON Schema 参数）\n2. 模型返回 `tool_call`（函数名 + 参数）\n3. **应用层执行**并把结果回传\n4. 模型基于结果继续回答\n\n参数必须校验，高风险操作加人工确认。",
		[]string{"ai-app-dev"}, []model.Question{
			q("Function Calling 中真正执行函数的是？", []string{"大模型内部", "应用层代码", "向量数据库", "浏览器"}, 1, "模型只输出意图，执行在应用层。"),
			q("对模型返回的工具参数应该？", []string{"直接执行", "进行校验并对高风险操作加确认", "忽略", "转发给用户"}, 1, "模型输出不可完全信任。"),
		}},
	{"agent-react", "Agent 与 ReAct", "AI 应用", 4, 40, "思考-行动-观察循环与死循环防护。",
		"## ReAct\n\n思考 → 调用工具 → 观察结果 → 再思考，直到完成目标。\n\n## 防护\n\n- 最大迭代次数\n- 重复调用熔断\n- 清晰的错误反馈",
		[]string{"ai-app-dev"}, []model.Question{
			q("ReAct 范式的循环是？", []string{"编码-测试-部署", "思考-行动-观察", "训练-验证-上线", "提问-回答-结束"}, 1, "Reason + Act 交替进行。"),
			q("防止 Agent 死循环的手段不包括？", []string{"最大迭代次数", "重复调用熔断", "清晰错误反馈", "无限增加额度"}, 3, "增加额度只会放大损失。"),
		}},
	{"llm-eval", "评估驱动开发", "AI 应用", 4, 35, "黄金数据集、自动评估与回归。",
		"## EDD\n\n- 构建**黄金数据集**（问题 + 期望答案）\n- 每次迭代自动评估：准确率、相关性、忠实度\n- 指标退化即阻断发布",
		[]string{"ai-app-dev", "rag"}, []model.Question{
			q("LLM 应用为什么需要评估驱动开发？", []string{"输出非确定，传统断言失效", "为了写文档", "模型会自动变好", "评估可以替代测试"}, 0, "需要用数据集与指标衡量质量。"),
			q("RAG 回答\"忠实度\"指的是？", []string{"回答长度", "回答是否基于检索内容而非编造", "响应速度", "用户满意度"}, 1, "忠实度衡量是否幻觉。"),
		}},
	// ---------- Python & 数据 ----------
	{"python-basics", "Python 基础", "Python", 1, 25, "数据类型、函数、推导式与虚拟环境。",
		"## 要点\n\n- 列表/字典推导式\n- 函数与装饰器\n- 使用 `venv` 隔离依赖\n\n```python\nsquares = [x * x for x in range(10)]\n```",
		[]string{"python"}, []model.Question{
			q("隔离项目依赖推荐使用？", []string{"全局 pip install", "venv 虚拟环境", "复制 site-packages", "sudo 安装"}, 1, "虚拟环境避免依赖冲突。"),
			q("[x*x for x in range(3)] 的结果是？", []string{"[1,4,9]", "[0,1,4]", "[0,1,2]", "9"}, 1, "range(3) 为 0,1,2。"),
		}},
	{"etl", "ETL 与数据管道", "数据", 3, 35, "抽取、转换、加载与调度。",
		"## 要点\n\n- E：从业务库/日志抽取（增量优先）\n- T：清洗、去重、口径统一\n- L：写入数仓分层（ODS → DWD → DWS → ADS）\n- 调度与重跑需要**幂等**",
		[]string{"data-pipeline", "python"}, []model.Question{
			q("数仓分层中最贴近原始数据的是？", []string{"ADS", "ODS", "DWS", "DIM"}, 1, "ODS 为贴源层。"),
			q("ETL 任务支持重跑的关键是？", []string{"多开线程", "幂等写入", "关闭日志", "全量覆盖所有表"}, 1, "幂等保证重跑结果一致。"),
		}},
	// ---------- 前端 ----------
	{"html-semantic", "HTML 语义化", "前端", 1, 15, "用合适的标签表达内容结构。",
		"## 要点\n\n- `header` / `nav` / `main` / `article` / `footer`\n- 语义化提升可访问性与 SEO\n- 按钮用 `button`，不要用 `div` 模拟",
		[]string{"html-css"}, []model.Question{
			q("页面主要内容区域应使用？", []string{"div", "main", "span", "section"}, 1, "main 表示文档主体内容。"),
			q("语义化 HTML 的好处不包括？", []string{"可访问性", "SEO", "代码可读性", "让 JS 执行更快"}, 3, "语义化与 JS 执行速度无关。"),
		}},
	{"flex-layout", "Flex 布局", "前端", 2, 20, "主轴、交叉轴与常用布局模式。",
		"## 常用\n\n```css\n.row { display: flex; justify-content: space-between; align-items: center; gap: 12px; }\n```\n\n- `justify-content` 控制主轴\n- `align-items` 控制交叉轴\n- `flex: 1` 占满剩余空间",
		[]string{"html-css"}, []model.Question{
			q("控制 flex 主轴对齐的属性是？", []string{"align-items", "justify-content", "flex-wrap", "order"}, 1, "justify-content 作用于主轴。"),
			q("让一个子元素占满剩余空间用？", []string{"flex: 1", "width: auto", "display: block", "float: left"}, 0, "flex: 1 分配剩余空间。"),
		}},
	{"js-async", "Promise 与异步编程", "前端", 3, 30, "事件循环、Promise 与 async/await。",
		"## 要点\n\n- 事件循环：同步 → 微任务（Promise）→ 宏任务（setTimeout）\n- `Promise.all` 并发、`Promise.allSettled` 不因单个失败中断\n- `async/await` 中用 `try/catch` 处理错误",
		[]string{"javascript"}, []model.Question{
			q("Promise.then 回调属于？", []string{"宏任务", "微任务", "同步任务", "渲染任务"}, 1, "Promise 回调进入微任务队列。"),
			q("多个请求并发且需要全部结果（即使部分失败）用？", []string{"Promise.all", "Promise.allSettled", "Promise.race", "for 循环 await"}, 1, "allSettled 返回每个结果的状态。"),
		}},
	{"vue-reactivity", "Vue 响应式原理", "前端", 3, 30, "Proxy、依赖收集与 computed。",
		"## 要点\n\n- Vue3 基于 **Proxy** 拦截读写\n- 读取时收集依赖，写入时触发更新\n- `computed` 有缓存；`watch` 处理副作用\n- 解构 reactive 会丢失响应式 → 使用 `toRefs`",
		[]string{"vue"}, []model.Question{
			q("Vue3 响应式基于？", []string{"Object.defineProperty", "Proxy", "setInterval", "MutationObserver"}, 1, "Vue3 使用 Proxy。"),
			q("直接解构 reactive 对象会？", []string{"性能提升", "丢失响应式", "报错", "自动转为 ref"}, 1, "需使用 toRefs 保持响应式。"),
		}},
	{"vue-component", "组件设计", "前端", 3, 30, "Props/Emits、插槽与组合式函数复用。",
		"## 原则\n\n- 单向数据流：Props 向下，事件向上\n- 用插槽提升灵活性\n- 逻辑复用使用 **composables**\n- 组件职责单一",
		[]string{"vue"}, []model.Question{
			q("子组件修改父组件数据的正确方式是？", []string{"直接修改 props", "emit 事件由父组件修改", "使用全局变量", "修改 $parent"}, 1, "遵循单向数据流。"),
			q("Vue3 推荐的逻辑复用方式是？", []string{"mixins", "composables 组合式函数", "继承", "全局方法"}, 1, "composables 来源清晰、无命名冲突。"),
		}},
	// ---------- 产品 ----------
	{"requirement-analysis", "需求分析方法", "产品", 2, 30, "用户故事、场景拆解与优先级。",
		"## 方法\n\n- 用户故事：作为 <角色>，我想要 <功能>，以便 <价值>\n- 场景拆解：主流程 + 异常流程\n- 优先级：价值 × 紧急度，MVP 先行",
		[]string{"requirement"}, []model.Question{
			q("用户故事的标准结构包含？", []string{"技术方案", "角色、功能、价值", "数据库表", "开发排期"}, 1, "以用户价值为中心描述需求。"),
			q("MVP 的核心思想是？", []string{"功能越多越好", "用最小可用版本验证核心价值", "先做完美设计", "跳过用户研究"}, 1, "快速验证，降低风险。"),
		}},
	{"ai-scenario", "AI 场景设计与评估", "产品", 3, 30, "识别适合 AI 的场景，定义效果指标。",
		"## 判断是否适合 AI\n\n- 任务有模式但规则难以穷举\n- 能容忍一定错误率，或有人机协同兜底\n- 有可衡量的效果指标（准确率、采纳率、节省时长）",
		[]string{"ai-product", "prompt-eng"}, []model.Question{
			q("以下哪个场景最适合引入大模型？", []string{"精确的金额计算", "客服工单自动分类与回复草稿", "银行转账执行", "固定格式的数据校验"}, 1, "模式多样、可人机协同的文本任务适合大模型。"),
			q("衡量 AI 回复草稿功能效果的关键指标是？", []string{"模型参数量", "客服采纳率与节省时长", "服务器数量", "页面 PV"}, 1, "以业务价值指标衡量。"),
		}},
}

type relDef struct{ From, To, Type string }

var relationDefs = []relDef{
	{"java-oop", "java-collections", model.RelationPrerequisite},
	{"java-oop", "thread", model.RelationPrerequisite},
	{"thread", "thread-safety", model.RelationPrerequisite},
	{"thread", "jmm", model.RelationPrerequisite},
	{"jmm", "volatile", model.RelationPrerequisite},
	{"thread-safety", "synchronized", model.RelationPrerequisite},
	{"jmm", "cas", model.RelationPrerequisite},
	{"cas", "aqs", model.RelationPrerequisite},
	{"synchronized", "aqs", model.RelationAdvanced},
	{"aqs", "thread-pool", model.RelationPrerequisite},
	{"volatile", "synchronized", model.RelationRelated},
	{"cas", "synchronized", model.RelationSimilar},
	{"java-oop", "spring-ioc", model.RelationPrerequisite},
	{"spring-ioc", "spring-rest", model.RelationPrerequisite},
	{"spring-ioc", "spring-transaction", model.RelationPrerequisite},
	{"mysql-transaction", "spring-transaction", model.RelationPrerequisite},
	{"sql-basics", "mysql-index", model.RelationPrerequisite},
	{"sql-basics", "mysql-transaction", model.RelationPrerequisite},
	{"mysql-index", "slow-query", model.RelationPrerequisite},
	{"redis-datatypes", "cache-patterns", model.RelationPrerequisite},
	{"redis-datatypes", "distributed-lock", model.RelationPrerequisite},
	{"synchronized", "distributed-lock", model.RelationRelated},
	{"spring-rest", "service-split", model.RelationPrerequisite},
	{"service-split", "message-queue", model.RelationRelated},
	{"message-queue", "distributed-transaction", model.RelationPrerequisite},
	{"mysql-transaction", "distributed-transaction", model.RelationPrerequisite},
	{"thread-pool", "load-testing", model.RelationRelated},
	{"load-testing", "jvm-gc", model.RelationRelated},
	{"slow-query", "load-testing", model.RelationRelated},
	{"service-split", "high-availability", model.RelationPrerequisite},
	{"cache-patterns", "high-availability", model.RelationRelated},
	{"high-availability", "system-design-method", model.RelationPrerequisite},
	{"llm-basics", "prompt-structure", model.RelationPrerequisite},
	{"prompt-structure", "few-shot", model.RelationPrerequisite},
	{"llm-basics", "embedding", model.RelationPrerequisite},
	{"embedding", "chunking", model.RelationPrerequisite},
	{"chunking", "hybrid-retrieval", model.RelationPrerequisite},
	{"prompt-structure", "function-calling", model.RelationPrerequisite},
	{"function-calling", "agent-react", model.RelationPrerequisite},
	{"hybrid-retrieval", "llm-eval", model.RelationRelated},
	{"agent-react", "llm-eval", model.RelationAdvanced},
	{"python-basics", "etl", model.RelationPrerequisite},
	{"sql-basics", "etl", model.RelationPrerequisite},
	{"html-semantic", "flex-layout", model.RelationPrerequisite},
	{"js-async", "vue-reactivity", model.RelationPrerequisite},
	{"vue-reactivity", "vue-component", model.RelationPrerequisite},
	{"flex-layout", "vue-component", model.RelationRelated},
	{"requirement-analysis", "ai-scenario", model.RelationPrerequisite},
	{"llm-basics", "ai-scenario", model.RelationPrerequisite},
}

type capDef struct {
	ID, Name, Desc string
	Weight         int
	Skills         [][2]interface{} // {skillID, requiredLevel}
}

type roleDef struct {
	ID, Name, Level, Desc string
	Caps                  []capDef
}

type careerDef struct {
	ID, Name, Category, Desc, Icon, Color, Demand, Salary string
	Roles                                                 []roleDef
}

func req(skill string, level int) [2]interface{} { return [2]interface{}{skill, level} }

var careerDefs = []careerDef{
	{"software-engineer", "软件工程师", "研发", "设计、开发与维护企业级软件系统，是数字化业务的基石。", "code", "#6b5cff", "高", "12–45K", []roleDef{
		{"java-engineer", "Java 工程师", "初级", "独立完成业务模块的开发与维护。", []capDef{
			{"je-dev", "企业应用开发能力", "能基于 Spring Boot 完成业务接口开发", 3, [][2]interface{}{req("java-core", 3), req("spring-boot", 2), req("mysql", 2)}},
			{"je-data", "数据存储能力", "能设计表结构并正确使用缓存", 2, [][2]interface{}{req("mysql", 2), req("redis", 1)}},
		}},
		{"senior-java", "高级 Java 工程师", "高级", "负责核心系统设计与性能、稳定性，能指导初级工程师。", []capDef{
			{"sj-dev", "企业应用开发能力", "高质量交付复杂业务系统", 2, [][2]interface{}{req("java-core", 3), req("spring-boot", 3)}},
			{"sj-dist", "分布式系统能力", "解决分布式环境下的一致性与可靠性问题", 3, [][2]interface{}{req("distributed", 3), req("microservice", 2), req("redis", 3)}},
			{"sj-perf", "性能优化能力", "定位并解决性能瓶颈", 3, [][2]interface{}{req("performance", 3), req("mysql", 3), req("java-concurrency", 3)}},
			{"sj-design", "系统设计能力", "完成模块级架构设计", 2, [][2]interface{}{req("system-design", 3)}},
			{"sj-solution", "技术方案设计能力", "输出可评审、可落地的技术方案", 2, [][2]interface{}{req("system-design", 3), req("microservice", 3)}},
		}},
		{"java-architect", "Java 架构师", "专家", "主导系统架构演进与技术决策。", []capDef{
			{"ja-arch", "架构设计能力", "设计高可用、可扩展的大型系统", 4, [][2]interface{}{req("system-design", 5), req("distributed", 4), req("microservice", 4)}},
			{"ja-stab", "性能与稳定性", "保障系统在高并发下稳定运行", 3, [][2]interface{}{req("performance", 4), req("java-concurrency", 4)}},
		}},
	}},
	{"ai-engineer", "AI 工程师", "研发", "把大模型能力落地为可靠的业务应用，是 AI 时代增长最快的岗位。", "sparkles", "#0fb981", "极高", "18–60K", []roleDef{
		{"ai-app-engineer", "AI 应用工程师", "中级", "基于大模型构建知识库问答、智能体等 AI 应用。", []capDef{
			{"ai-llm", "大模型交互能力", "设计稳定可控的提示词与输出约束", 2, [][2]interface{}{req("prompt-eng", 3), req("ai-app-dev", 2)}},
			{"ai-kb", "知识库构建能力", "构建高召回、低幻觉的 RAG 系统", 3, [][2]interface{}{req("rag", 3), req("python", 2)}},
			{"ai-eng", "AI 应用工程化能力", "完成 Agent、工具调用与评估体系", 3, [][2]interface{}{req("ai-app-dev", 3), req("spring-boot", 2)}},
		}},
		{"llm-architect", "大模型应用架构师", "专家", "设计企业级 AI 平台与智能体架构。", []capDef{
			{"la-arch", "AI 架构设计能力", "设计可扩展的 AI 应用平台", 3, [][2]interface{}{req("ai-app-dev", 4), req("rag", 4), req("system-design", 4)}},
			{"la-quality", "效果与质量保障", "建立评估与持续优化机制", 2, [][2]interface{}{req("prompt-eng", 4), req("ai-app-dev", 4)}},
		}},
	}},
	{"frontend-engineer", "前端工程师", "研发", "构建体验优秀的 Web 应用，连接用户与产品。", "layers", "#2a8cf4", "高", "10–35K", []roleDef{
		{"web-frontend", "Web 前端工程师", "初级", "独立完成页面与组件开发。", []capDef{
			{"wf-page", "页面开发能力", "还原设计稿并适配多端", 2, [][2]interface{}{req("html-css", 3), req("javascript", 3)}},
			{"wf-app", "应用开发能力", "使用 Vue 构建交互应用", 2, [][2]interface{}{req("vue", 2), req("javascript", 3)}},
		}},
		{"senior-frontend", "高级前端工程师", "高级", "负责前端架构、性能与工程化。", []capDef{
			{"sf-arch", "前端架构能力", "设计可维护的组件体系与状态管理", 3, [][2]interface{}{req("vue", 4), req("javascript", 4), req("system-design", 2)}},
			{"sf-perf", "体验与性能优化", "优化加载与交互性能", 2, [][2]interface{}{req("performance", 2), req("html-css", 3)}},
		}},
	}},
	{"data-engineer", "数据工程师", "数据", "构建稳定的数据管道与数仓，为分析与 AI 提供高质量数据。", "chart", "#f59e0b", "高", "15–40K", []roleDef{
		{"data-engineer-role", "数据工程师", "中级", "负责数据采集、建模与管道开发。", []capDef{
			{"de-pipe", "数据管道能力", "开发可重跑、可监控的 ETL 任务", 3, [][2]interface{}{req("data-pipeline", 3), req("python", 3)}},
			{"de-model", "数据建模能力", "设计数仓分层与指标口径", 2, [][2]interface{}{req("mysql", 3), req("data-pipeline", 2)}},
		}},
	}},
	{"product-manager", "产品经理", "产品", "发现用户价值，定义并推动产品落地。", "target", "#ef4d63", "中", "12–40K", []roleDef{
		{"ai-pm", "AI 产品经理", "中级", "识别 AI 场景，定义 AI 产品与效果指标。", []capDef{
			{"pm-req", "需求分析能力", "把模糊问题拆解为可交付的需求", 2, [][2]interface{}{req("requirement", 3)}},
			{"pm-ai", "AI 产品设计能力", "判断 AI 能力边界并设计人机协同", 3, [][2]interface{}{req("ai-product", 3), req("prompt-eng", 2)}},
		}},
	}},
}

// chapterKnowledgeDefs 现有课程章节与知识点的映射。
var chapterKnowledgeDefs = map[string]map[string][]string{
	"prompt-engineering": {
		"basics":     {"llm-basics", "prompt-structure"},
		"structured": {"prompt-structure"},
		"advanced":   {"few-shot"},
	},
	"rag-in-action": {
		"overview":  {"llm-basics", "embedding"},
		"chunking":  {"chunking", "embedding"},
		"retrieval": {"hybrid-retrieval"},
	},
	"ai-native-dev": {
		"core-concepts":  {"llm-basics", "function-calling"},
		"architecture":   {"agent-react"},
		"implementation": {"function-calling"},
		"code-practice":  {"agent-react"},
		"best-practices": {"llm-eval"},
		"pitfalls":       {"agent-react"},
		"performance":    {"llm-eval"},
	},
}
