/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

// API Response types
interface ApiResponse<T> {
  code: number
  data: T
  msg: string
}

interface PageResponse<T> {
  list: T[]
  total: number
  page: number
  pageSize: number
}

// Entity types
interface SMTPConfig {
  id: number
  userId: number
  userName: string
  name: string
  host: string
  port: number
  username: string
  password: string
  mail: string
  createTime: string
  updateTime: string
}

interface MailConfig {
  id: number
  userId: number
  userName: string
  name: string
  mail: string
  createTime: string
  updateTime: string
}

interface TokenConfig {
  id: number
  userId: number
  userName: string
  token: string
  name: string
  SMTPId: number
  SMTP: SMTPConfig | null
  SMTPStr: string
  receiveEmailIds: string
  receiveEmails: MailConfig[]
  receiveEmailsStr: string
  subject: string
  state: number // 0: 使用中 1: 用户禁用 2: 管理员禁用
  createTime: string
  updateTime: string
}

interface ScheduleConfig {
  id: number
  userId: number
  userName: string
  name: string
  tokenId: number
  cron: string
  url: string
  method: string
  headerStr: string
  headers: Array<{name: string; value: string}>
  body: string
  useTemplate: boolean
  template: string
  isUse: boolean
  isSendOnlyByError: boolean
  createTime: string
  updateTime: string
}

interface LogEntry {
  id: number
  userId: number
  /** Owner's name, filled in for admins (the list shows every user's entries there) */
  userName?: string
  name: string
  mail: string
  token: string
  smtp: string
  subject: string
  content: string
  /** JSON string: [{ name, filePath }] */
  files: string
  /** Numeric status code from the backend (0 success, 1 warning, 2 error) */
  status: number
  /** Human-readable status, resolved server-side to match the log.* i18n keys */
  statusStr: string
  result: string
  createTime: string
  updateTime: string
}

interface UserConfig {
  id: number
  name: string
  password: string
  isAdmin: boolean
  isUse: boolean
  createTime: string
  updateTime: string
}

interface SetInfo {
  webPort: number
  apiPort: number
  dbType: 'sqlite' | 'mysql'
  dbHost: string
  dbPort: number
  dbName: string
  dbUser: string
  dbPass: string
  dbCharset: string
  dbFile: string
  adminUser: string
  adminPass: string
}

interface SystemInfo {
  initialized: boolean
  dbInitialized: boolean
  hasAdmin: boolean
  hasLogin?: boolean
  isDocker: boolean
  /** Only set when the session cookie is valid */
  username?: string
  isAdmin?: boolean
  version: string
}
