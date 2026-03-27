<template>
  <el-card class="subscription-card" :class="{ active: subscription.active }">
    <div v-if="subscription.active" class="active-badge"></div>
    <div class="card-content">
      <div class="card-main" @click="toggleExpand">
        <div class="card-header">
          <div class="header-row">
            <div class="name-row">
              <span class="name">{{ subscription.name }}</span>
              <el-tag v-if="subscription.active" type="success" size="small">激活中</el-tag>
            </div>
            <div class="header-actions" @click.stop>
              <el-button
                v-if="!subscription.active"
                type="primary"
                :size="isMobile ? 'small' : 'default'"
                @click="$emit('activate')"
              >
                激活
              </el-button>
              <el-dropdown trigger="click">
                <el-button :size="isMobile ? 'small' : 'default'">
                  操作 <el-icon class="el-icon--right"><ArrowDown /></el-icon>
                </el-button>
                <template #dropdown>
                  <el-dropdown-menu>
                    <el-dropdown-item @click.stop="$emit('edit')">
                      <el-icon><Edit /></el-icon>编辑
                    </el-dropdown-item>
                    <el-dropdown-item v-if="subscription.source_type === 'remote'" @click.stop="$emit('refresh')">
                      <el-icon><Refresh /></el-icon>刷新
                    </el-dropdown-item>
                    <el-dropdown-item @click.stop="$emit('rules')">
                      <el-icon><Document /></el-icon>规则 ({{ ruleCount }})
                    </el-dropdown-item>
                    <el-dropdown-item @click.stop="$emit('scripts')">
                      <el-icon><Notebook /></el-icon>脚本 ({{ scriptCount }})
                    </el-dropdown-item>
                    <el-dropdown-item @click.stop="$emit('content')">
                      <el-icon><Tickets /></el-icon>查看源文件
                    </el-dropdown-item>
                    <el-dropdown-item v-if="subscription.active" @click.stop="$emit('mergedConfig')">
                      <el-icon><View /></el-icon>查看合并配置
                    </el-dropdown-item>
                    <el-dropdown-item divided @click.stop="$emit('delete')">
                      <el-icon color="#ef4444"><Delete /></el-icon>
                      <span style="color: #ef4444;">删除</span>
                    </el-dropdown-item>
                  </el-dropdown-menu>
                </template>
              </el-dropdown>
            </div>
          </div>
          <div class="meta-row">
            <el-tag :type="subscription.source_type === 'remote' ? 'info' : 'warning'" size="small" effect="plain">
              {{ subscription.source_type === 'remote' ? '远程订阅' : '本地配置' }}
            </el-tag>
            <span v-if="subscription.node_count" class="meta-item">{{ subscription.node_count }} 节点</span>
            <span v-if="subscription.source_type === 'remote' && subscription.last_refresh" class="meta-item">
              {{ formatRelativeTime(subscription.last_refresh) }}更新
            </span>
            <span v-else-if="subscription.source_type === 'local' && subscription.updated_at" class="meta-item">
              {{ formatRelativeTime(subscription.updated_at) }}编辑
            </span>
          </div>
        </div>

        <div v-if="subscription.total_transfer > 0" class="traffic-section">
          <TrafficInfo
            :upload-used="subscription.upload_used"
            :download-used="subscription.download_used"
            :total-transfer="subscription.total_transfer"
            :expire-at="subscription.expire_at"
          />
        </div>

        <el-collapse-transition>
          <div v-show="expanded" class="detail-section">
            <div class="detail-grid">
              <template v-if="subscription.source_type === 'remote'">
                <div class="detail-item">
                  <span class="detail-label">订阅地址</span>
                  <span class="detail-value url">{{ truncateUrl(subscription.url) }}</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">刷新间隔</span>
                  <span class="detail-value">每 {{ subscription.interval }} 分钟</span>
                </div>
                <div class="detail-item">
                  <span class="detail-label">上次刷新</span>
                  <span class="detail-value">{{ formatRelativeTime(subscription.last_refresh) }}</span>
                </div>
              </template>
              <template v-else>
                <div class="detail-item">
                  <span class="detail-label">上次编辑</span>
                  <span class="detail-value">{{ subscription.updated_at ? formatRelativeTime(subscription.updated_at) : '从未编辑' }}</span>
                </div>
              </template>
              <div class="detail-item">
                <span class="detail-label">创建时间</span>
                <span class="detail-value">{{ formatDateTime(subscription.created_at) }}</span>
              </div>
            </div>

            <div v-if="subscription.use_proxy || subscription.user_agent || subscription.skip_cert" class="advanced-section">
              <el-divider content-position="left">请求设置</el-divider>
              <div class="advanced-tags">
                <el-tag v-if="subscription.use_proxy" type="info" size="small">代理更新</el-tag>
                <el-tag v-if="subscription.user_agent" size="small">自定义 UA</el-tag>
                <el-tag v-if="subscription.skip_cert" type="warning" size="small">跳过证书校验</el-tag>
              </div>
            </div>
          </div>
        </el-collapse-transition>

        <div class="expand-hint">
          <el-icon :class="{ rotated: expanded }"><ArrowDown /></el-icon>
          <span>{{ expanded ? '收起详情' : '展开详情' }}</span>
        </div>
      </div>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { ArrowDown, Edit, Refresh, Document, Notebook, View, Delete, Tickets } from '@element-plus/icons-vue'
import type { Subscription } from '@/types/api'
import { formatDateTime, formatRelativeTime } from '@/utils/format'
import { useBreakpoint } from '@/composables/useBreakpoint'
import TrafficInfo from './TrafficInfo.vue'

const { isMobile } = useBreakpoint()

defineProps<{
  subscription: Subscription
  ruleCount: number
  scriptCount: number
  loading?: boolean
}>()

defineEmits<{
  activate: []
  refresh: []
  edit: []
  delete: []
  rules: []
  scripts: []
  mergedConfig: []
  content: []
}>()

const expanded = ref(false)

const toggleExpand = () => {
  expanded.value = !expanded.value
}

const truncateUrl = (url: string) => {
  if (url.length > 60) {
    return url.slice(0, 60) + '...'
  }
  return url
}
</script>

<style scoped>
.subscription-card {
  margin-bottom: 16px;
  transition: all 0.2s ease;
  position: relative;
  overflow: visible;
}

.active-badge {
  position: absolute;
  top: -1px;
  left: -1px;
  width: 0;
  height: 0;
  border-style: solid;
  border-width: 45px 45px 0 0;
  border-color: #60a5fa transparent transparent transparent;
  z-index: 1;
  filter: drop-shadow(2px 2px 3px rgba(0, 0, 0, 0.2));
}

.card-content {
  padding: 4px 0;
}

.card-main {
  cursor: pointer;
}

.card-header {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.header-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 12px;
}

.name-row {
  display: flex;
  align-items: center;
  gap: 10px;
  flex: 1;
  min-width: 0;
}

.name {
  font-size: 18px;
  font-weight: 600;
  color: #1e293b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.meta-row {
  display: flex;
  align-items: center;
  gap: 12px;
}

.meta-item {
  font-size: 13px;
  color: #64748b;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.traffic-section {
  margin-top: 16px;
}

.detail-section {
  margin-top: 16px;
  padding-top: 16px;
  border-top: 1px solid #e2e8f0;
}

.detail-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 12px 24px;
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.detail-label {
  font-size: 12px;
  color: #64748b;
}

.detail-value {
  font-size: 14px;
  color: #1e293b;
}

.detail-value.url {
  font-size: 12px;
  color: #64748b;
  word-break: break-all;
  font-family: monospace;
}

.advanced-section {
  margin-top: 16px;
}

.advanced-tags {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.expand-hint {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  margin-top: 12px;
  color: #64748b;
  font-size: 13px;
}

.expand-hint .el-icon {
  transition: transform 0.3s ease;
}

.expand-hint .el-icon.rotated {
  transform: rotate(180deg);
}

@media (max-width: 640px) {
  .name {
    font-size: 16px;
  }

  .detail-grid {
    grid-template-columns: 1fr;
  }
}
</style>