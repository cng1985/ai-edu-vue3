/** 后台导航布局模式配置 */
export const NAV_MODES = [
  {
    key: 'classic',
    label: '经典风格',
    description: '深色侧边栏 + 顶部分组标签',
    icon: 'Menu'
  },
  {
    key: 'win11',
    label: 'Win11 风格',
    description: '图标导航栏 + 悬浮展开面板',
    icon: 'Grid'
  },
  {
    key: 'modern',
    label: '现代风格',
    description: '顶部水平导航，无侧边栏',
    icon: 'Operation'
  }
]

export const DEFAULT_NAV_MODE = 'classic'
export const DEFAULT_COLLAPSED = true
