// Markdown 渲染工具：marked + highlight.js，mermaid 代码块输出为占位容器，
// 由 MarkdownRenderer 组件在挂载后异步渲染为 SVG 图表。

import { Marked } from 'marked'
import hljs from 'highlight.js/lib/core'
import javascript from 'highlight.js/lib/languages/javascript'
import python from 'highlight.js/lib/languages/python'
import java from 'highlight.js/lib/languages/java'
import json from 'highlight.js/lib/languages/json'
import xml from 'highlight.js/lib/languages/xml'
import bash from 'highlight.js/lib/languages/bash'
import sql from 'highlight.js/lib/languages/sql'
import plaintext from 'highlight.js/lib/languages/plaintext'

hljs.registerLanguage('javascript', javascript)
hljs.registerLanguage('python', python)
hljs.registerLanguage('java', java)
hljs.registerLanguage('json', json)
hljs.registerLanguage('xml', xml)
hljs.registerLanguage('html', xml)
hljs.registerLanguage('bash', bash)
hljs.registerLanguage('sql', sql)
hljs.registerLanguage('text', plaintext)
hljs.registerLanguage('plaintext', plaintext)

function escapeHtml(text) {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
    .replace(/"/g, '&quot;')
}

const marked = new Marked({
  gfm: true,
  breaks: false,
  renderer: {
    code({ text, lang }) {
      const language = (lang || '').trim().toLowerCase()
      if (language === 'mermaid') {
        return `<div class="mermaid-block" data-mermaid="${escapeHtml(encodeURIComponent(text))}"><pre class="mermaid-source">${escapeHtml(text)}</pre></div>`
      }
      let highlighted
      try {
        highlighted = hljs.getLanguage(language)
          ? hljs.highlight(text, { language }).value
          : hljs.highlight(text, { language: 'plaintext' }).value
      } catch (e) {
        highlighted = escapeHtml(text)
      }
      const label = language || 'text'
      return `<div class="code-block"><div class="code-block__header"><span>${escapeHtml(label)}</span></div><pre><code class="hljs language-${escapeHtml(label)}">${highlighted}</code></pre></div>`
    }
  }
})

export function renderMarkdown(markdown) {
  return marked.parse(markdown || '')
}

let mermaidPromise = null

function loadMermaid() {
  if (!mermaidPromise) {
    mermaidPromise = import('mermaid').then(({ default: mermaid }) => {
      mermaid.initialize({
        startOnLoad: false,
        securityLevel: 'loose',
        theme: 'neutral',
        fontFamily: 'inherit'
      })
      return mermaid
    })
  }
  return mermaidPromise
}

let mermaidSeq = 0

/** 将容器内的 mermaid 占位块渲染为 SVG */
export async function renderMermaidIn(container) {
  const blocks = container.querySelectorAll('.mermaid-block[data-mermaid]')
  if (blocks.length === 0) return
  const mermaid = await loadMermaid()
  for (const block of blocks) {
    const source = decodeURIComponent(block.dataset.mermaid)
    block.removeAttribute('data-mermaid')
    try {
      const { svg } = await mermaid.render(`mmd-${Date.now()}-${mermaidSeq++}`, source)
      block.innerHTML = svg
      block.classList.add('mermaid-block--rendered')
    } catch (e) {
      // 渲染失败时保留源码展示
      block.classList.add('mermaid-block--error')
    }
  }
}
