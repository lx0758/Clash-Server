<template>
  <el-dialog
    :model-value="true"
    :title="'规则编辑 - ' + subscriptionName"
    :width="isMobile ? '100%' : '900px'"
    :fullscreen="isMobile"
    :class="{ 'mobile-dialog': isMobile }"
    @close="$emit('close')"
  >
    <div class="toolbar">
      <el-select v-model="ruleMode">
        <el-option value="insert" label="插入模式 (规则插入到订阅规则之前)" />
        <el-option value="append" label="追加模式 (规则追加到订阅规则之后)" />
      </el-select>
      <el-button type="primary" @click="addRule">+ 添加规则</el-button>
    </div>

    <div class="table-wrapper">
      <el-empty v-if="rules.length === 0" description="暂无规则" />

      <div v-else class="rules-container">
        <div v-for="(rule, index) in rules" :key="rule.id" class="rule-item">
          <el-button
            type="danger"
            size="small"
            circle
            class="rule-delete-mobile"
            @click="removeRule(index)"
          >
            <el-icon><Delete /></el-icon>
          </el-button>
          <div class="rule-row">
            <el-select v-model="rule.type" size="small" class="rule-type">
              <el-option value="DOMAIN" label="DOMAIN" />
              <el-option value="DOMAIN-SUFFIX" label="DOMAIN-SUFFIX" />
              <el-option value="DOMAIN-KEYWORD" label="DOMAIN-KEYWORD" />
              <el-option value="IP-CIDR" label="IP-CIDR" />
              <el-option value="GEOIP" label="GEOIP" />
              <el-option value="GEOSITE" label="GEOSITE" />
              <el-option value="PROCESS-NAME" label="PROCESS-NAME" />
              <el-option value="RULE-SET" label="RULE-SET" />
            </el-select>
            <el-input v-model="rule.payload" placeholder="匹配内容" size="small" class="rule-payload" />
            <el-input v-model="rule.proxy" placeholder="代理" size="small" class="rule-proxy" />
            <el-button type="danger" size="small" text @click="removeRule(index)" class="rule-delete">
              <el-icon><Delete /></el-icon>
            </el-button>
          </div>
        </div>
      </div>
    </div>

    <template #footer>
      <el-button @click="$emit('close')">取消</el-button>
      <el-button type="primary" :loading="loading" @click="saveRules">保存</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, onMounted, watch } from 'vue'
import { Delete } from '@element-plus/icons-vue'
import { ruleApi } from '@/api/rule'
import { useBreakpoint } from '@/composables/useBreakpoint'
import type { Rule } from '@/types/api'

const { isMobile } = useBreakpoint()

const props = defineProps<{
  subscriptionId: number
  subscriptionName: string
}>()

const emit = defineEmits<{
  close: []
}>()

const loading = ref(false)
const rules = ref<Array<{ id: string; type: string; payload: string; proxy: string; dbId?: number }>>([])
const ruleMode = ref<'insert' | 'append'>('append')

const fetchRules = async () => {
  const res = await ruleApi.list(props.subscriptionId)
  const data = res.data.data.rules || []
  rules.value = data.map((r: Rule) => ({
    id: `edit-${r.id}`,
    type: r.type,
    payload: r.payload,
    proxy: r.proxy,
    dbId: r.id,
  }))
  if (data.length > 0 && data[0]) {
    ruleMode.value = data[0].mode
  }
}

onMounted(fetchRules)
watch(() => props.subscriptionId, fetchRules)

const addRule = () => {
  rules.value.push({
    id: `new-${Date.now()}`,
    type: 'DOMAIN',
    payload: '',
    proxy: '',
  })
}

const removeRule = (index: number) => {
  rules.value.splice(index, 1)
}

const saveRules = async () => {
  loading.value = true
  try {
    const existing = await ruleApi.list(props.subscriptionId)
    const existingIds = new Set(existing.data.data.rules?.map((r: Rule) => r.id) || [])
    const editedIds = new Set(rules.value.filter(r => r.dbId).map(r => r.dbId!))

    for (const id of existingIds) {
      if (!editedIds.has(id)) {
        await ruleApi.delete(props.subscriptionId, id)
      }
    }

    for (let i = 0; i < rules.value.length; i++) {
      const r = rules.value[i]
      if (!r || !r.payload || !r.proxy) continue
      const data = {
        name: `${r.type}-${r.payload.slice(0, 20)}`,
        type: r.type,
        payload: r.payload,
        proxy: r.proxy,
        mode: ruleMode.value,
        priority: i,
        enabled: true,
      }
      if (r.dbId) {
        await ruleApi.update(props.subscriptionId, r.dbId, data)
      } else {
        await ruleApi.create(props.subscriptionId, data)
      }
    }
    emit('close')
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.toolbar {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  flex-wrap: wrap;
}

.toolbar .el-select {
  flex: 1;
  min-width: 200px;
}

.rules-container {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.rule-item {
  padding: 8px;
  background: #f8fafc;
  border-radius: 6px;
  position: relative;
}

.rule-delete-mobile {
  display: none;
  position: absolute;
  top: 4px;
  right: 4px;
  width: 24px;
  height: 24px;
}

.rule-row {
  display: flex;
  gap: 8px;
  align-items: center;
}

.rule-type {
  width: 140px;
  flex-shrink: 0;
}

.rule-payload {
  flex: 1;
  min-width: 100px;
}

.rule-proxy {
  width: 120px;
  flex-shrink: 0;
}

@media (max-width: 768px) {
  .toolbar {
    flex-direction: column;
  }

  .toolbar .el-select {
    width: 100%;
  }

  .rule-item {
    padding-top: 32px;
  }

  .rule-delete-mobile {
    display: flex;
  }

  .rule-row {
    flex-wrap: wrap;
  }

  .rule-type {
    width: 100%;
  }

  .rule-payload {
    width: 100%;
  }

  .rule-proxy {
    width: 100%;
  }

  .rule-delete {
    display: none;
  }
}
</style>

<style>
.mobile-dialog .el-dialog__body {
  display: flex;
  flex-direction: column;
  height: calc(100vh - 120px);
  overflow: hidden;
}

.mobile-dialog .el-dialog__body > .toolbar {
  flex-shrink: 0;
}

.mobile-dialog .el-dialog__body > .table-wrapper {
  flex: 1;
  overflow-y: auto;
}

.mobile-dialog .el-dialog__footer {
  position: fixed;
  bottom: 0;
  left: 0;
  right: 0;
  background: #fff;
  border-top: 1px solid #e2e8f0;
  padding: 16px 20px;
  z-index: 100;
}
</style>