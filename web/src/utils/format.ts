export function formatBytes(bytes: number): string {
  if (!bytes || !isFinite(bytes) || bytes <= 0) return '0 B'
  const units = ['B', 'KB', 'MB', 'GB', 'TB']
  const k = 1024
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  const safeI = Math.max(0, Math.min(i, units.length - 1))
  const value = bytes / Math.pow(k, safeI)
  return parseFloat(value.toFixed(safeI === 0 ? 0 : 2)) + ' ' + units[safeI]
}

export function formatPercent(used: number, total: number): number {
  if (total === 0) return 0
  return Math.round((used / total) * 100)
}

export function formatDaysLeft(expireAt: string | null): number | null {
  if (!expireAt) return null
  const expire = new Date(expireAt)
  const now = new Date()
  const diff = expire.getTime() - now.getTime()
  const days = Math.ceil(diff / (1000 * 60 * 60 * 24))
  return days > 0 ? days : 0
}

export function isExpired(expireAt: string | null): boolean {
  if (!expireAt) return false
  return new Date(expireAt) < new Date()
}

export function isExpiringSoon(expireAt: string | null, days: number = 7): boolean {
  if (!expireAt) return false
  const daysLeft = formatDaysLeft(expireAt)
  return daysLeft !== null && daysLeft <= days && daysLeft > 0
}

export function formatDateTime(dateStr: string | null): string {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  return date.toLocaleString('zh-CN', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit',
  })
}

export function formatRelativeTime(dateStr: string | null): string {
  if (!dateStr) return '-'
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / (1000 * 60))
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 30) return `${days}天前`
  return formatDateTime(dateStr)
}

export function formatTraffic(bytes: number, simplify = false): string {
  if (bytes === 0) return '0 B/s'
  const units = ['B/s', 'KB/s', 'MB/s', 'GB/s']
  const k = 1024
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  const value = parseFloat((bytes / Math.pow(k, i)).toFixed(2))

  if (simplify) {
    const simplifiedUnits = ['B', 'K', 'M', 'G']
    return `${value}${simplifiedUnits[i]}`
  }

  return `${value} ${units[i]}`
}