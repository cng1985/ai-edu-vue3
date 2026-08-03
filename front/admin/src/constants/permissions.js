/** 与后端 server/pkg/rbac/rbac.go 保持一致 */
export const PERM = {
  USER_READ: 'user:read',
  USER_CREATE: 'user:create',
  USER_UPDATE: 'user:update',
  USER_DELETE: 'user:delete',
  COURSE_READ: 'course:read',
  COURSE_WRITE: 'course:write',
  COURSE_DELETE: 'course:delete',
  QUIZ_READ: 'quiz:read',
  QUIZ_WRITE: 'quiz:write',
  QUIZ_DELETE: 'quiz:delete',
  REVIEW_READ: 'review:read',
  REVIEW_APPROVE: 'review:approve',
  DASHBOARD: 'dashboard:read',
  ROLE_MANAGE: 'role:manage',
  AI_CHAT: 'ai:chat',
  CUSTOMER_READ: 'customer:read',
  CUSTOMER_REPLY: 'customer:reply'
}
