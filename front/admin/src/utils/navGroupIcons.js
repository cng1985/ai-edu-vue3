import {
  DataAnalysis,
  User,
  Reading,
  Service,
  TrendCharts,
  Medal,
  Collection,
  Setting,
  Cpu
} from '@element-plus/icons-vue'

export const GROUP_ICONS = {
  operations: DataAnalysis,
  users: User,
  content: Reading,
  service: Service,
  ecosystem: TrendCharts,
  talent: Medal,
  knowledge: Collection,
  ai: Cpu,
  system: Setting
}

export function getGroupIcon(key) {
  return GROUP_ICONS[key] || DataAnalysis
}
