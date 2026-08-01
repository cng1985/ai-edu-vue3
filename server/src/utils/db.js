import fs from 'fs'
import path from 'path'
import { fileURLToPath } from 'url'

const __dirname = path.dirname(fileURLToPath(import.meta.url))
const DB_PATH = path.join(__dirname, '../../data/store.json')

const defaultData = {
  users: [],
  courses: [],
  quizzes: [],
  reviews: [],
  settings: {
    siteName: 'AI 学习系统',
    maintenanceMode: false
  }
}

let cache = null

export function loadDb() {
  if (cache) return cache
  try {
    if (fs.existsSync(DB_PATH)) {
      cache = JSON.parse(fs.readFileSync(DB_PATH, 'utf-8'))
      return cache
    }
  } catch {
    /* 损坏时重建 */
  }
  cache = structuredClone(defaultData)
  saveDb()
  return cache
}

export function saveDb() {
  const dir = path.dirname(DB_PATH)
  if (!fs.existsSync(dir)) fs.mkdirSync(dir, { recursive: true })
  fs.writeFileSync(DB_PATH, JSON.stringify(cache, null, 2), 'utf-8')
}

export function resetCache() {
  cache = null
}

export function genId(prefix = 'id') {
  return `${prefix}_${Date.now().toString(36)}_${Math.random().toString(36).slice(2, 7)}`
}
