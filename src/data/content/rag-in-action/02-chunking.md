# 文档切分与向量化策略

"垃圾进，垃圾出"在 RAG 系统中体现得淋漓尽致。文档切分（Chunking）的质量直接决定了检索召回的上限。

## 为什么切分策略如此关键

- **太大的块**：包含多个主题，向量表征被稀释，且占用宝贵的上下文窗口
- **太小的块**：语义不完整，模型拿到碎片信息无法回答
- **切断语义**：在句子或段落中间硬切，破坏语义连贯性

## 常见切分策略

### 1. 固定长度切分

最简单的方案：按固定字符数（如 500 字）切分，配合重叠窗口（overlap，如 50 字）缓解语义切断。

```javascript
function chunkText(text, chunkSize = 500, overlap = 50) {
  const chunks = []
  let start = 0
  while (start < text.length) {
    const end = Math.min(start + chunkSize, text.length)
    chunks.push(text.slice(start, end))
    if (end === text.length) break
    start = end - overlap
  }
  return chunks
}
```

### 2. 结构感知切分

利用文档的天然结构（Markdown 标题、HTML 标签、代码函数边界）作为切分点，保证每个块是语义完整的单元。这是**生产环境推荐的默认策略**。

### 3. 语义切分

用 Embedding 模型计算相邻句子的语义相似度，在语义突变处切分。效果最好但计算成本最高。

## 元数据的重要性

每个 Chunk 除了正文，还应携带元数据（metadata），用于检索时过滤：

```json
{
  "content": "报销单需要在费用发生后 30 天内提交……",
  "metadata": {
    "source": "员工手册.pdf",
    "chapter": "财务制度",
    "tenant_id": "org_42",
    "updated_at": "2026-05-01"
  }
}
```

检索时结合元数据过滤（如只查当前租户、只查最新版本文档），能显著提升精度并保障数据隔离。

## Embedding 模型的选择与陷阱

1. **维度不是越高越好**：高维向量存储与检索成本更高，需在效果与成本间平衡
2. **查询与文档的不对称性**：查询通常很短、文档块较长，部分模型（如 BGE）要求对查询添加特定前缀指令
3. **切换模型必须全量重建索引**：不同模型的向量空间不兼容，混用会导致检索完全失效

> **实践建议**：先用结构感知切分 + 主流 Embedding 模型跑通基线，再基于评估集数据驱动地调整块大小与策略，不要过早优化。
