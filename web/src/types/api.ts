export interface ApiResponse<T = unknown> {
  code: number
  message: string
  data: T
  core_error?: string
}

export interface User {
  id: number
  username: string
  created_at: string
}

export type SourceType = 'remote' | 'local'

export interface Subscription {
  id: number
  name: string
  source_type: SourceType
  url: string
  interval: number
  content: string
  active: boolean
  use_proxy: boolean
  user_agent: string
  skip_cert: boolean
  last_refresh: string | null
  node_count: number
  upload_used: number
  download_used: number
  total_transfer: number
  expire_at: string | null
  created_at: string
  updated_at: string
  rule_count?: number
  script_count?: number
}

export interface SubscriptionWithCounts {
  subscription: Subscription
  rule_count: number
  script_count: number
}

export interface Rule {
  id: number
  subscription_id: number
  name: string
  type: string
  payload: string
  proxy: string
  enabled: boolean
  mode: 'insert' | 'append'
  priority: number
  created_at: string
}

export interface Script {
  id: number
  subscription_id: number
  name: string
  description: string
  content: string
  enabled: boolean
  created_at: string
}

export interface MergedConfig {
  config: Record<string, unknown>
  yaml: string
}

export interface SubscriptionDetail {
  subscription: Subscription
  rules: Rule[]
  scripts: Script[]
}

export interface RefreshResult {
  subscription: Subscription
  node_count: number
  refresh_error?: string
}

export interface Proxy {
  name: string
  type: string
  alive: boolean
  delay: number
}

export interface MihomoProxyGroup {
  name: string
  type: string
  all: string[]
  now: string
  alive?: boolean
}

export interface MihomoProxy {
  name: string
  type: string
  alive: boolean
  history: Array<{ time: string; delay: number }>
}

export interface Proxy {
  name: string
  type: string
  alive: boolean
  delay: number
}

export interface ProxyGroup {
  name: string
  type: string
  all: Proxy[]
  now: string
}

export interface Connection {
  id: string
  metadata?: {
    host?: string
    network?: string
    type?: string
    sourceIP?: string
    sourcePort?: string
    destinationIP?: string
    destinationPort?: string
    process?: string
  }
  upload?: number
  download?: number
  start?: string
  chains?: string[]
  rule?: string
  rulePayload?: string
  disconnected?: boolean
}

export interface ConnectionsData {
  connections: Connection[]
  history: Connection[]
  downloadTotal: number
  uploadTotal: number
}

export interface Traffic {
  up: number
  down: number
}

export interface ConnectionsData {
  connections: Connection[]
  downloadTotal: number
  uploadTotal: number
}

export interface LogData {
  type: string
  payload: string
}

export interface MemoryData {
  inuse: number
  oslimit: number
}

export interface CoreRule {
  index: number
  type: string
  payload: string
  proxy: string
  size: number
  extra?: {
    disabled: boolean
    hitCount: number
    hitAt: string
    missCount: number
    missAt: string
  }
}

export interface CoreRulesData {
  rules: CoreRule[]
}

export interface SystemInfo {
  core: CoreStatus
  subscription: {
    count: number
    proxy_count: number
  }
  traffic: Traffic | null
}

export interface Config {
  server: ServerConfig
  core: CoreConfig
}

export interface ServerConfig {
  host: string
  port: number
  database: string
}

export interface CoreConfig {
  api_host: string
  api_port: number
  api_secret: string
  mixed_port: number
  allow_lan: boolean
  mode: string
  log_level: string
  ipv6: boolean
}

export interface CoreStatus {
  running: boolean
  version?: string
  error?: string
}