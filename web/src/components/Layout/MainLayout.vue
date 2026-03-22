<template>
  <el-container class="layout">
    <el-aside :width="isCollapsed ? '64px' : '200px'" class="sidebar" v-if="!isMobile">
      <div class="logo">
        <img src="/favicon.svg" class="logo-icon" alt="Logo" />
        <span v-if="!isCollapsed" class="logo-text">ClashServer</span>
      </div>
      <el-menu
        :default-active="currentRoute"
        :collapse="isCollapsed"
        router
        class="sidebar-menu"
      >
        <el-menu-item index="/">
          <el-icon><Odometer /></el-icon>
          <template #title>仪表盘</template>
        </el-menu-item>
        <el-menu-item index="/proxies">
          <el-icon><Connection /></el-icon>
          <template #title>代理</template>
        </el-menu-item>
        <el-menu-item index="/rules">
          <el-icon><Document /></el-icon>
          <template #title>规则</template>
        </el-menu-item>
        <el-menu-item index="/connections">
          <el-icon><Link /></el-icon>
          <template #title>连接</template>
        </el-menu-item>
        <el-menu-item index="/logs">
          <el-icon><Notebook /></el-icon>
          <template #title>日志</template>
        </el-menu-item>
        <el-menu-item index="/subscriptions">
          <el-icon><FolderOpened /></el-icon>
          <template #title>订阅</template>
        </el-menu-item>
        <el-menu-item index="/settings">
          <el-icon><Setting /></el-icon>
          <template #title>设置</template>
        </el-menu-item>
      </el-menu>
      <div class="sidebar-footer">
        <SidebarTraffic v-if="!isCollapsed" />
        <div class="user-section" v-if="!isCollapsed">
          <span class="username">{{ authStore.user?.username }}</span>
          <el-button type="danger" size="small" text @click="logout">
            <el-icon><SwitchButton /></el-icon>
            退出
          </el-button>
        </div>
        <el-button
          :icon="isCollapsed ? Expand : Fold"
          @click="toggleCollapse"
          text
        />
      </div>
    </el-aside>

    <el-drawer
      v-model="drawerVisible"
      direction="ltr"
      :show-close="false"
      :with-header="false"
      size="200px"
      v-if="isMobile"
    >
      <div class="logo">
        <img src="/favicon.svg" class="logo-icon" alt="Logo" />
        <span class="logo-text">ClashServer</span>
      </div>
      <el-menu
        :default-active="currentRoute"
        router
        class="drawer-menu"
        @select="drawerVisible = false"
      >
        <el-menu-item index="/">
          <el-icon><Odometer /></el-icon>
          <template #title>仪表盘</template>
        </el-menu-item>
        <el-menu-item index="/proxies">
          <el-icon><Connection /></el-icon>
          <template #title>代理</template>
        </el-menu-item>
        <el-menu-item index="/rules">
          <el-icon><Document /></el-icon>
          <template #title>规则</template>
        </el-menu-item>
        <el-menu-item index="/connections">
          <el-icon><Link /></el-icon>
          <template #title>连接</template>
        </el-menu-item>
        <el-menu-item index="/logs">
          <el-icon><Notebook /></el-icon>
          <template #title>日志</template>
        </el-menu-item>
        <el-menu-item index="/subscriptions">
          <el-icon><FolderOpened /></el-icon>
          <template #title>订阅</template>
        </el-menu-item>
        <el-menu-item index="/settings">
          <el-icon><Setting /></el-icon>
          <template #title>设置</template>
        </el-menu-item>
      </el-menu>
    </el-drawer>

    <el-container>
      <el-header class="header" v-if="isMobile">
        <el-button :icon="Expand" @click="drawerVisible = true" text />
        <div class="header-center">
          <TrafficIndicator />
        </div>
        <div class="header-right">
          <el-dropdown>
            <span class="user-dropdown">
              {{ authStore.user?.username }}
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="logout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>
      <el-main class="main">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useBreakpoint } from '@/composables/useBreakpoint'
import { useWebSocket } from '@/composables/useWebSocket'
import {
  Odometer,
  Connection,
  Document,
  Link,
  Notebook,
  FolderOpened,
  Setting,
  Fold,
  Expand,
  ArrowDown,
  SwitchButton,
} from '@element-plus/icons-vue'
import TrafficIndicator from '@/components/dashboard/TrafficIndicator.vue'
import SidebarTraffic from '@/components/Layout/SidebarTraffic.vue'

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const { isMobile } = useBreakpoint()
const { connect, disconnect } = useWebSocket()

const isCollapsed = ref(localStorage.getItem('sidebar-collapsed') === 'true')
const drawerVisible = ref(false)

const currentRoute = computed(() => route.path)

const toggleCollapse = () => {
  isCollapsed.value = !isCollapsed.value
  localStorage.setItem('sidebar-collapsed', String(isCollapsed.value))
}

const logout = async () => {
  await authStore.logout()
  router.push('/login')
}

watch(isMobile, (mobile) => {
  if (mobile) {
    drawerVisible.value = false
  }
})

onMounted(() => {
  connect()
})

onUnmounted(() => {
  disconnect()
})
</script>

<style scoped>
.layout {
  height: 100vh;
  background: #f1f5f9;
}

.layout > .el-container {
  height: 100%;
}

.sidebar {
  background: #ffffff;
  display: flex;
  flex-direction: column;
  transition: width 0.3s ease;
  overflow: hidden;
  border-right: 1px solid #e2e8f0;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: #3b82f6;
  font-size: 20px;
  font-weight: bold;
  border-bottom: 1px solid #e2e8f0;
}

.logo-icon {
  width: 36px;
  height: 36px;
  flex-shrink: 0;
}

.logo-text {
  white-space: nowrap;
}

.sidebar-menu {
  flex: 1;
  border-right: none;
}

.sidebar-menu:not(.el-menu--collapse) {
  width: 200px;
}

.drawer-menu {
  border-right: none;
}

.sidebar-footer {
  border-top: 1px solid #e2e8f0;
  display: flex;
  flex-direction: column;
}

.user-section {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 10px;
  justify-content: space-between;
}

.username {
  font-size: 14px;
  color: #1e293b;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.header {
  background: #ffffff;
  display: flex;
  align-items: center;
  padding: 0 16px;
  border-bottom: 1px solid #e2e8f0;
}

.header-center {
  flex: 1;
  display: flex;
  justify-content: center;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;
}

.user-dropdown {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #64748b;
  cursor: pointer;
}

.main {
  padding: 16px;
  overflow-y: auto;
  flex: 1;
  min-height: 0;
}

@media (max-width: 768px) {
  .main {
    padding: 12px;
  }
}
</style>