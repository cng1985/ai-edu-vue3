# 常见坑与排错指南

在将AI原生应用推向生产环境的过程中，开发者往往会遭遇传统软件工程中不存在的“暗礁”。大模型的非确定性特征，使得系统在具备高度智能的同时，也引入了不可忽视的脆弱性。本章将深入剖析模型输出格式不稳定、无限递归调用、Token超限及API超时四大高频陷阱，并提供针对性的调试与降级容错策略。

### 1. 模型输出格式不稳定

**陷阱表现**：在提示词中要求模型输出JSON，但模型偶尔会“自作主张”地加上Markdown代码块标记（如 ```json ... ```），或者在JSON前后附带冗长的解释性文本，导致下游代码执行 `JSON.parse()` 时抛出异常。

**调试与容错策略**：
*   **正则截取提取**：不要假定模型输出100%是纯JSON。在解析前，使用正则表达式提取第一个 `{` 和最后一个 `}` 之间的内容。
*   **JsonRepair 库兜底**：面对缺失引号、尾随逗号等常见LLM语法错误，引入第三方修复库（如 Python 的 `json-repair` 或 JS 的 `jsonrepair`）进行自动修复。
*   **结构化输出 API**：优先使用 OpenAI 等平台提供的 Structured Outputs / Function Calling 特性，在API层面强制模型输出符合特定 JSON Schema 的结果。

以下是一个包含正则截取与容错修复的 Python 代码示例：

```python
import re
import json
from json_repair import repair_json

def parse_llm_json(raw_output: str) -> dict:
    """
    容错解析 LLM 的 JSON 输出
    """
    # 策略1：直接尝试解析（针对表现良好的输出）
    try:
        return json.loads(raw_output)
    except json.JSONDecodeError:
        pass
    
    # 策略2：正则提取首个 { 到最后一个 } 之间的内容
    match = re.search(r'\{.*\}', raw_output, re.DOTALL)
    if match:
        extracted_str = match.group(0)
        try:
            return json.loads(extracted_str)
        except json.JSONDecodeError:
            # 策略3：使用 json-repair 修复残缺的 JSON 语法
            repaired_str = repair_json(extracted_str)
            try:
                return json.loads(repaired_str)
            except Exception as e:
                raise ValueError(f"JSON 解析与修复均失败: {e}")
    
    raise ValueError("未能在输出中找到有效的 JSON 结构")

# 测试用例
test_cases = [
    '{"name": "Alice", "age": 30}',                     # 正常 JSON
    '```json\n{"name": "Bob", "age": 25}\n```',        # 包含 Markdown 标记
    '好的，这是结果：{"name": "Charlie", "age": 20}',  # 包含前置文本
    '{"name": "David", "age": 20,}'                     # 尾随逗号错误
]

for case in test_cases:
    print(f"原始输出: {case}")
    print(f"解析结果: {parse_llm_json(case)}\n")
```

### 2. 无限递归调用与工具调用死循环

**陷阱表现**：在赋予LLM Agent工具调用能力后，模型可能陷入“调用工具 -> 解析结果失败 -> 再次调用相同工具”的死循环。这会迅速耗尽API额度，甚至导致系统资源耗尽。

**调试与容错策略**：
*   **强制最大迭代次数**：在 Agent 循环逻辑中硬编码最大执行步数（如 `max_iterations = 5`），超过阈值则强制中断并返回降级提示。
*   **错误反馈显式化**：当工具执行报错时，不要只返回简单的错误码，应将清晰的、指导性的文本反馈给模型（如：“数据库查询失败，原因是参数格式不正确，请修正后重试”）。
*   **重复调用熔断**：在上下文记忆中记录最近的工具调用，若连续两次调用相同工具且参数一致，则触发熔断机制。

### 3. Token 超限与上下文截断

**陷阱表现**：在长对话或多轮RAG检索中，输入Token数不知不觉超过了模型的上下文窗口限制（如超过8K或128K），导致API报错，或者更糟糕的是，系统静默截断了早期的关键指令，引发“失忆”或指令偏离。

**调试与容错策略**：
*   **Token 预计算**：在发送请求前，使用 `tiktoken` 等库计算当前 Prompt 的 Token 数。
*   **滑动窗口与摘要压缩**：保留最近的 N 轮对话原文，对更早的历史调用小模型生成摘要，用摘要替代原文注入上下文。
*   **RAG 重排截断**：在RAG检索阶段召回大量文档后，引入 Reranker 模型进行相关性重排，只截取 Top-K（如 Top-3）最相关的片段，严格控制注入上下文的体积。

### 4. API 超时与限流

**陷阱表现**：大模型推理耗时较长，特别是在要求输出长文本时。遇到网络波动或服务商负载高峰时，应用极易出现网关超时（504 Gateway Timeout）或限流错误（429 Too Many Requests）。

**调试与容错策略**：
*   **指数退避重试**：针对 429 或 5xx 错误，实现带有随机抖动的指数退避重试机制，避免重试风暴。
*   **流式输出 首字优化**：强制使用 SSE (Server-Sent Events) 流式接口。流式输出不仅能让用户更快看到响应（降低前端体感超时），还能在底层保持长连接，避免传统 HTTP 请求的网关超时切断。
*   **多级降级路由**：配置多模型路由策略。当主模型（如 GPT-4）响应超时或限流时，自动降级到备用模型（如 GPT-3.5 或开源 Llama 系列），保障核心业务可用性。

下表总结了各类陷阱的核心降级策略：

| 陷阱类型 | 核心原因 | 降级容错策略 | 兜底方案 |
| :--- | :--- | :--- | :--- |
| **格式不稳定** | 非确定性生成 | 正则截取、JsonRepair、结构化API | 返回默认模板，引导用户重试 |
| **无限递归** | 工具反馈缺失/幻觉 | 最大步数限制、重复调用熔断 | 返回“当前无法处理该请求” |
| **Token超限** | 上下文无限膨胀 | Token预计算、历史摘要、RAG重排 | 截断最早的历史记录 |
| **API超时** | 网络/服务商负载 | 指数退避重试、流式输出 | 降级到轻量级备用模型 |

> **架构提示**：在AI原生架构中，“容错”不再是异常处理的附属品，而是核心设计准则。开发者必须在每一层（提示词层、Agent调度层、API通信层）都预设“模型会出错”的假设，并构建相应的防护栏。
