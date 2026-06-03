<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getModelRates, type ModelCreditRate } from '@/api/credits'

const rates = ref<ModelCreditRate[]>([])
const loading = ref(false)
const error = ref('')

async function load() {
  loading.value = true
  error.value = ''
  try {
    const res = await getModelRates()
    rates.value = res.items
  } catch (e: any) {
    error.value = e.message ?? '加载失败'
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>

<template>
  <div class="p-6 space-y-5 max-w-4xl">
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-xl font-semibold text-gray-900">模型积分费率</h1>
        <p class="text-sm text-gray-500 mt-1">
          配置每个模型每 1000 个 token 的积分消耗。<code class="bg-gray-100 px-1 rounded">*</code> 通配符匹配任意字符（如 <code class="bg-gray-100 px-1 rounded">claude-opus-4*</code>）。Priority 最高的规则生效。
        </p>
      </div>
      <button @click="load" class="text-sm text-blue-600 hover:text-blue-700">刷新</button>
    </div>

    <div v-if="error" class="rounded-md bg-red-50 border border-red-200 p-3 text-sm text-red-700">
      {{ error }}
    </div>

    <div v-if="loading" class="text-center py-12 text-gray-400">
      <div class="inline-block w-6 h-6 border-2 border-blue-500 border-t-transparent rounded-full animate-spin mb-2" />
      <p class="text-sm">加载中…</p>
    </div>

    <div v-else class="overflow-hidden rounded-xl border border-gray-200 shadow-sm">
      <table class="min-w-full divide-y divide-gray-200 text-sm">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-5 py-3 text-left text-xs font-semibold text-gray-500 uppercase tracking-wider">模型匹配模式</th>
            <th class="px-5 py-3 text-right text-xs font-semibold text-gray-500 uppercase tracking-wider">输入 /1K tokens</th>
            <th class="px-5 py-3 text-right text-xs font-semibold text-gray-500 uppercase tracking-wider">输出 /1K tokens</th>
            <th class="px-5 py-3 text-right text-xs font-semibold text-gray-500 uppercase tracking-wider">Priority</th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-100">
          <tr v-for="rate in rates" :key="rate.id" class="hover:bg-gray-50 transition-colors">
            <td class="px-5 py-3">
              <code class="font-mono text-gray-800 bg-gray-50 px-1.5 py-0.5 rounded text-xs">{{ rate.model_pattern }}</code>
            </td>
            <td class="px-5 py-3 text-right tabular-nums text-gray-700">{{ rate.credits_per_1k_tokens_input }}</td>
            <td class="px-5 py-3 text-right tabular-nums text-gray-700">{{ rate.credits_per_1k_tokens_output }}</td>
            <td class="px-5 py-3 text-right tabular-nums">
              <span class="inline-flex items-center justify-center w-8 h-5 text-xs font-medium rounded"
                :class="rate.priority >= 90 ? 'bg-blue-100 text-blue-700' : rate.priority >= 40 ? 'bg-gray-100 text-gray-600' : 'bg-gray-50 text-gray-400'">
                {{ rate.priority }}
              </span>
            </td>
          </tr>
          <tr v-if="rates.length === 0">
            <td colspan="4" class="px-5 py-10 text-center text-gray-400">暂无费率配置</td>
          </tr>
        </tbody>
      </table>
    </div>

    <p class="text-xs text-gray-400">
      如需修改费率，请在数据库 <code class="bg-gray-100 px-1 rounded">model_credit_rates</code> 表中直接操作，或通过 API 调用。
    </p>
  </div>
</template>
