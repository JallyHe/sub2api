<script setup lang="ts">
import { ref } from 'vue'
import { getUserBalance, grantCredits, type CreditBalance } from '@/api/credits'

// --- Grant Credits ---
const grantForm = ref({ user_id: '', credits: '', notes: '' })
const grantLoading = ref(false)
const grantResult = ref('')
const grantError = ref('')

async function handleGrant() {
  grantError.value = ''
  grantResult.value = ''
  const uid = Number(grantForm.value.user_id)
  const credits = Number(grantForm.value.credits)
  if (!uid || !credits || credits < 1) {
    grantError.value = '请填写有效的用户 ID 和积分数量（≥1）'
    return
  }
  grantLoading.value = true
  try {
    await grantCredits({ user_id: uid, credits, notes: grantForm.value.notes })
    grantResult.value = `✓ 已向用户 ${uid} 充值 ${credits.toLocaleString()} 积分`
    grantForm.value = { user_id: '', credits: '', notes: '' }
  } catch (e: any) {
    grantError.value = e.response?.data?.error ?? e.message ?? '充值失败'
  } finally {
    grantLoading.value = false
  }
}

// --- Check Balance ---
const checkId = ref('')
const checkResult = ref<CreditBalance | null>(null)
const checkError = ref('')
const checkLoading = ref(false)

async function handleCheck() {
  checkError.value = ''
  checkResult.value = null
  const uid = Number(checkId.value)
  if (!uid) return
  checkLoading.value = true
  try {
    checkResult.value = await getUserBalance(uid)
  } catch (e: any) {
    checkError.value = e.response?.data?.error ?? e.message ?? '查询失败'
  } finally {
    checkLoading.value = false
  }
}

function formatExpiry(ts: number | null) {
  if (!ts) return '永不过期'
  return new Date(ts * 1000).toLocaleDateString('zh-CN')
}
</script>

<template>
  <div class="p-6 space-y-8 max-w-2xl">
    <h1 class="text-xl font-semibold text-gray-900">积分管理</h1>

    <!-- Grant Credits -->
    <section class="rounded-xl border border-gray-200 bg-white p-6 space-y-4 shadow-sm">
      <h2 class="font-medium text-gray-800">手动充值积分</h2>
      <div class="grid grid-cols-2 gap-4">
        <div>
          <label class="block text-xs text-gray-500 mb-1">用户 ID</label>
          <input
            v-model="grantForm.user_id" type="number" placeholder="123"
            class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>
        <div>
          <label class="block text-xs text-gray-500 mb-1">积分数量</label>
          <input
            v-model="grantForm.credits" type="number" placeholder="10000"
            class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          />
        </div>
      </div>
      <div>
        <label class="block text-xs text-gray-500 mb-1">备注（可选）</label>
        <input
          v-model="grantForm.notes" type="text" placeholder="管理员手动补偿"
          class="w-full rounded-lg border border-gray-200 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
        />
      </div>
      <div v-if="grantError" class="text-sm text-red-600 bg-red-50 rounded-lg px-3 py-2">{{ grantError }}</div>
      <div v-if="grantResult" class="text-sm text-green-700 bg-green-50 rounded-lg px-3 py-2">{{ grantResult }}</div>
      <button
        @click="handleGrant" :disabled="grantLoading"
        class="rounded-lg bg-blue-600 px-5 py-2 text-sm font-medium text-white hover:bg-blue-700 disabled:opacity-50 transition-colors"
      >
        {{ grantLoading ? '充值中…' : '充值' }}
      </button>
    </section>

    <!-- Check Balance -->
    <section class="rounded-xl border border-gray-200 bg-white p-6 space-y-4 shadow-sm">
      <h2 class="font-medium text-gray-800">查询用户积分余额</h2>
      <div class="flex gap-3">
        <input
          v-model="checkId" type="number" placeholder="用户 ID"
          class="flex-1 rounded-lg border border-gray-200 px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          @keydown.enter="handleCheck"
        />
        <button
          @click="handleCheck" :disabled="checkLoading"
          class="rounded-lg bg-gray-100 px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-200 disabled:opacity-50 transition-colors"
        >
          {{ checkLoading ? '查询中…' : '查询' }}
        </button>
      </div>
      <div v-if="checkError" class="text-sm text-red-600 bg-red-50 rounded-lg px-3 py-2">{{ checkError }}</div>
      <div v-if="checkResult" class="rounded-lg bg-gray-50 p-4 text-sm space-y-2">
        <div class="flex justify-between">
          <span class="text-gray-500">余额</span>
          <span class="font-semibold tabular-nums">{{ checkResult.balance.toLocaleString() }} 积分</span>
        </div>
        <div class="flex justify-between">
          <span class="text-gray-500">到期时间</span>
          <span>{{ formatExpiry(checkResult.expires_at) }}</span>
        </div>
      </div>
    </section>
  </div>
</template>
