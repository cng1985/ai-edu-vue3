import {
  DataAnalysis,
  User,
  Reading,
  Service,
  Ticket,
  Collection,
  Setting
} from '@element-plus/icons-vue'

export const GROUP_ICONS = {
  operations: DataAnalysis,
  users: User,
  content: Reading,
  service: Service,
  documents: Ticket,
  knowledge: Collection,
  system: Setting
}

export function getGroupIcon(key) {
  return GROUP_ICONS[key] || DataAnalysis
}
