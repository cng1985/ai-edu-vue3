import {
  DataAnalysis,
  User,
  Reading,
  Service,
  Ticket,
  Collection,
  Setting,
  Cpu
} from '@element-plus/icons-vue'

export const GROUP_ICONS = {
  operations: DataAnalysis,
  users: User,
  content: Reading,
  service: Service,
  documents: Ticket,
  knowledge: Collection,
  ai: Cpu,
  system: Setting
}

export function getGroupIcon(key) {
  return GROUP_ICONS[key] || DataAnalysis
}
