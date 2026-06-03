<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { getBalance, type CreditBalance } from '@/api/credits'

const balance = ref<CreditBalance | null>(null)
const loading = ref(true)
const error = ref(false)

onMounted(async () => {
  try {
    balance.value = await getBalance()
  } catch {
    error.value = true
  } finally {
    loading.value = false
  }
})

const expiresText = computed(() => {
  if (!balance.value?.expires_at) return '永不过期'
  const d = new Date(balance.value.expires_at * 1000)
  return d.toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
})

const isLow = computed(() => (balance.value?.balance ?? 0) <= 0)
</script>

<template>
  <div class="rounded-xl border border-gray-200 bg-white p-5 shadow-sm">
    <div class="flex items-center justify-between mb-3">
      <h3 class="text-sm font-medium text-gray-500">积分余额</h3>
      <span v-if="!loading && balance" class="text-xs text-gray-400">到期：{{ expiresText }}</span>
    </div>

    <!-- Loading skeleton -->
    <div v-if="loading" class="space-y-2">
      <div class="h-9 w-36 rounded-lg bg-gray-100 animate-pulse" />
      <div class="h-2 w-full rounded-full bg-gray-100 animate-pulse" />
    </div>

    <!-- Error state -->
    <div v-else-if="error" class="text-sm text-gray-400">积分信息加载失败</div>

    <!-- Balance display -->
    <div v-else-if="balance">
      <div class="flex items-baseline gap-1.5">
        <span class="text-3xl font-bold tabular-nums" :class="isLow ? 'text-red-600' : 'text-gray-900'">
          {{ balance.balance.toLocaleString() }}
        </span>
        <span class="text-sm text-gray-400">积分</span>
      </div>

      <!-- Progress bar (visual only, shows relative fullness up to 100k) -->
      <div class="mt-3 h-1.5 rounded-full bg-gray-100 overflow-hidden">
        <div
          class="h-full rounded-full transition-all duration-700 ease-out"
          :class="isLow ? 'bg-red-500' : balance.balance < 5000 ? 'bg-amber-500' : 'bg-blue-500'"
          :style="{ width: Math.min(100, (balance.balance / 100000) * 100) + '%' }"
        />
      </div>

      <p v-if="isLow" class="mt-2 text-xs text-red-500 font-medium">
        积分已耗尽，请前往「购买套餐」充值后继续使用
      </p>
    </div>
  </div>
</template>
