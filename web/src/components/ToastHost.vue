<script setup lang="ts">
import { useToastStore } from '@/stores/toast'

const toastStore = useToastStore()

const typeClass = (type: string) => `toast-${type}`
</script>

<template>
  <div class="toast-host">
    <TransitionGroup name="toast">
      <div
        v-for="toast in toastStore.toasts"
        :key="toast.id"
        class="toast-item"
        :class="typeClass(toast.type)"
        @click="toastStore.remove(toast.id)"
      >
        <span class="toast-dot"></span>
        <span class="toast-message">{{ toast.message }}</span>
      </div>
    </TransitionGroup>
  </div>
</template>

<style scoped>
.toast-host {
  position: fixed;
  top: 16px;
  right: 16px;
  z-index: 3000;
  display: flex;
  flex-direction: column;
  gap: 8px;
  pointer-events: none;
  max-width: min(420px, calc(100vw - 32px));
}

.toast-item {
  pointer-events: auto;
  display: flex;
  align-items: flex-start;
  gap: 8px;
  padding: 10px 14px;
  border-radius: 6px;
  background: #ffffff;
  border: 1px solid #e2e8f0;
  box-shadow: 0 4px 12px rgba(15, 23, 42, 0.12);
  font-size: 14px;
  line-height: 1.5;
  color: #1e293b;
  cursor: pointer;
}

.toast-dot {
  width: 8px;
  height: 8px;
  border-radius: 50%;
  margin-top: 6px;
  flex-shrink: 0;
}

.toast-success { border-left: 3px solid #10b981; }
.toast-success .toast-dot { background: #10b981; }

.toast-error { border-left: 3px solid #ef4444; }
.toast-error .toast-dot { background: #ef4444; }

.toast-warning { border-left: 3px solid #f59e0b; }
.toast-warning .toast-dot { background: #f59e0b; }

.toast-info { border-left: 3px solid #3b82f6; }
.toast-info .toast-dot { background: #3b82f6; }

.toast-enter-active,
.toast-leave-active {
  transition: all 0.25s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(20px);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(20px);
}
</style>
