<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { getBalance, type CreditBalance } from '@/api/credits'
import { apiClient } from '@/api/client'

const auth = useAuthStore()
const router = useRouter()
const isLoggedIn = computed(() => auth.isAuthenticated)
const isAdmin = computed(() => auth.isAdmin)

// Data
const creditBal = ref<CreditBalance | null>(null)
const subscription = ref<any>(null)
const recentUsage = ref<any[]>([])
const recentOrders = ref<any[]>([])
const loading = ref(true)

async function load() {
  if (!isLoggedIn.value) return
  loading.value = true
  try {
    const [bal, sub, usage, orders] = await Promise.allSettled([
      getBalance(),
      apiClient.get('/subscriptions/active').then(r => r.data).catch(() => null),
      apiClient.get('/usage?page=1&page_size=5').then(r => r.data).catch(() => ({ items: [] })),
      apiClient.get('/orders?page=1&page_size=5').then(r => r.data).catch(() => ({ items: [] })),
    ])
    creditBal.value = bal.status === 'fulfilled' ? bal.value : null
    subscription.value = sub.status === 'fulfilled' ? sub.value : null
    recentUsage.value = usage.status === 'fulfilled' ? (usage.value?.items ?? []) : []
    recentOrders.value = orders.status === 'fulfilled' ? (orders.value?.items ?? []) : []
  } finally {
    loading.value = false
  }
}

onMounted(load)

function formatDate(ts: string | number) {
  const d = typeof ts === 'number' ? new Date(ts * 1000) : new Date(ts)
  return d.toLocaleDateString('zh-CN', { month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' })
}

function formatExpiry(ts: string | null) {
  if (!ts) return '—'
  return new Date(ts).toLocaleDateString('zh-CN', { year: 'numeric', month: 'long', day: 'numeric' })
}

const isExpired = computed(() => {
  if (!subscription.value?.expires_at) return false
  return new Date(subscription.value.expires_at) < new Date()
})
</script>

<template>
  <div class="min-h-screen bg-gray-950 text-white flex flex-col">

    <!-- ── Header ── -->
    <header class="border-b border-gray-800/60 bg-gray-950/80 backdrop-blur sticky top-0 z-50">
      <div class="mx-auto flex h-14 max-w-4xl items-center justify-between px-5">
        <RouterLink to="/" class="flex items-center gap-2 font-bold">
          <span class="flex h-7 w-7 items-center justify-center rounded-lg bg-gradient-to-br from-violet-600 to-fuchsia-600 text-xs font-black">S</span>
          <span>StoryClaw</span>
        </RouterLink>

        <div v-if="isLoggedIn" class="flex items-center gap-3">
          <RouterLink v-if="isAdmin" to="/admin/dashboard"
            class="text-xs text-violet-400 border border-violet-500/30 rounded px-2.5 py-1 hover:bg-violet-500/10 transition-colors">
            管理后台
          </RouterLink>
          <RouterLink to="/profile" class="text-sm text-gray-400 hover:text-white transition-colors">个人设置</RouterLink>
          <RouterLink to="/purchase"
            class="rounded-lg bg-violet-600 px-3.5 py-1.5 text-sm font-medium hover:bg-violet-500 transition-colors">
            充值积分
          </RouterLink>
        </div>
        <div v-else class="flex items-center gap-3">
          <RouterLink to="/login" class="text-sm text-gray-400 hover:text-white transition-colors">登录</RouterLink>
          <RouterLink to="/register"
            class="rounded-lg bg-violet-600 px-3.5 py-1.5 text-sm font-medium hover:bg-violet-500 transition-colors">
            注册
          </RouterLink>
        </div>
      </div>
    </header>

    <!-- ── Main ── -->
    <main class="flex-1 mx-auto w-full max-w-4xl px-5 py-10">

      <!-- 未登录 -->
      <template v-if="!isLoggedIn">
        <div class="text-center space-y-6 pt-16">
          <div class="flex h-16 w-16 mx-auto items-center justify-center rounded-2xl bg-gradient-to-br from-violet-600 to-fuchsia-600 text-3xl font-black">S</div>
          <h1 class="text-3xl font-bold">欢迎使用 StoryClaw</h1>
          <p class="text-gray-400 max-w-sm mx-auto leading-relaxed">
            登录后查看积分余额、购买套餐，以及通过 StoryClaw 桌面应用一键授权配置 AI 模型。
          </p>
          <div class="flex justify-center gap-4 pt-2">
            <RouterLink to="/login"
              class="rounded-xl bg-violet-600 px-7 py-3 font-semibold hover:bg-violet-500 transition-colors">
              登录账号
            </RouterLink>
            <RouterLink to="/register"
              class="rounded-xl border border-gray-700 px-7 py-3 font-medium text-gray-300 hover:border-gray-500 transition-colors">
              创建账号
            </RouterLink>
          </div>
          <div class="pt-4">
            <RouterLink to="/#pricing" class="text-sm text-violet-400 hover:text-violet-300 transition-colors">
              查看套餐定价 →
            </RouterLink>
          </div>
        </div>
      </template>

      <!-- 已登录 -->
      <template v-else>
        <div class="space-y-8">

          <!-- 页标题 -->
          <div class="flex items-center justify-between">
            <h1 class="text-2xl font-bold">我的账户</h1>
            <span class="text-sm text-gray-500">{{ auth.user?.email }}</span>
          </div>

          <!-- 积分 + 订阅状态 -->
          <div class="grid gap-4 sm:grid-cols-2">

            <!-- 积分余额卡片 -->
            <div class="rounded-2xl border border-gray-800 bg-gray-900/60 p-6 space-y-4">
              <div class="flex items-center justify-between">
                <span class="text-sm text-gray-400">积分余额</span>
                <RouterLink to="/purchase"
                  class="text-xs text-violet-400 hover:text-violet-300 transition-colors">
                  充值 →
                </RouterLink>
              </div>
              <div v-if="loading" class="h-9 w-28 rounded-lg bg-gray-800 animate-pulse" />
              <div v-else>
                <div class="text-3xl font-bold tabular-nums"
                  :class="(creditBal?.balance ?? 0) <= 0 ? 'text-red-400' : 'text-white'">
                  {{ (creditBal?.balance ?? 0).toLocaleString() }}
                  <span class="text-base font-normal text-gray-400 ml-1">积分</span>
                </div>
                <p v-if="creditBal?.expires_at" class="text-xs text-gray-500 mt-1">
                  到期：{{ formatExpiry(String(creditBal.expires_at)) }}
                </p>
                <p v-if="(creditBal?.balance ?? 0) <= 0" class="text-xs text-red-400 mt-1">
                  积分已耗尽，请充值后继续使用
                </p>
              </div>
            </div>

            <!-- 订阅状态卡片 -->
            <div class="rounded-2xl border border-gray-800 bg-gray-900/60 p-6 space-y-4">
              <div class="flex items-center justify-between">
                <span class="text-sm text-gray-400">订阅状态</span>
                <RouterLink to="/subscriptions"
                  class="text-xs text-violet-400 hover:text-violet-300 transition-colors">
                  详情 →
                </RouterLink>
              </div>
              <div v-if="loading" class="h-9 w-32 rounded-lg bg-gray-800 animate-pulse" />
              <div v-else-if="!subscription">
                <div class="text-lg font-medium text-gray-400">暂无订阅</div>
                <RouterLink to="/purchase"
                  class="mt-3 inline-flex items-center gap-1 text-sm text-violet-400 hover:text-violet-300 transition-colors">
                  立即订阅 →
                </RouterLink>
              </div>
              <div v-else>
                <div class="flex items-center gap-2">
                  <span class="text-lg font-semibold text-white">{{ subscription.group?.name ?? '标准订阅' }}</span>
                  <span class="rounded-full px-2 py-0.5 text-xs font-medium"
                    :class="isExpired ? 'bg-red-900/40 text-red-400' : 'bg-green-900/40 text-green-400'">
                    {{ isExpired ? '已过期' : '有效中' }}
                  </span>
                </div>
                <p class="text-xs text-gray-500 mt-1">
                  到期：{{ formatExpiry(subscription.expires_at) }}
                </p>
              </div>
            </div>
          </div>

          <!-- 快捷入口 -->
          <div class="grid gap-3 grid-cols-2 sm:grid-cols-4">
            <RouterLink
              v-for="item in [
                { icon: '💳', label: '购买套餐', path: '/purchase' },
                { icon: '📊', label: '使用记录', path: '/usage' },
                { icon: '📋', label: '我的订单', path: '/orders' },
                { icon: '👤', label: '个人设置', path: '/profile' },
              ]" :key="item.path"
              :to="item.path"
              class="rounded-xl border border-gray-800 bg-gray-900/40 px-4 py-4 text-center hover:border-gray-600 hover:bg-gray-900/80 transition-all"
            >
              <div class="text-2xl mb-2">{{ item.icon }}</div>
              <div class="text-sm font-medium text-gray-300">{{ item.label }}</div>
            </RouterLink>
          </div>

          <!-- 最近使用记录 -->
          <div class="space-y-3">
            <div class="flex items-center justify-between">
              <h2 class="font-semibold text-gray-300">最近使用</h2>
              <RouterLink to="/usage" class="text-xs text-violet-400 hover:text-violet-300 transition-colors">
                全部记录 →
              </RouterLink>
            </div>
            <div class="rounded-xl border border-gray-800 bg-gray-900/40 overflow-hidden">
              <div v-if="loading" class="py-8 text-center text-gray-600 text-sm">加载中…</div>
              <div v-else-if="recentUsage.length === 0" class="py-8 text-center text-gray-600 text-sm">
                暂无使用记录
              </div>
              <table v-else class="w-full text-sm">
                <thead class="border-b border-gray-800">
                  <tr class="text-xs text-gray-500">
                    <th class="px-4 py-2.5 text-left">时间</th>
                    <th class="px-4 py-2.5 text-left">模型</th>
                    <th class="px-4 py-2.5 text-right">输入 tokens</th>
                    <th class="px-4 py-2.5 text-right">输出 tokens</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-800/60">
                  <tr v-for="item in recentUsage" :key="item.id" class="hover:bg-gray-800/30">
                    <td class="px-4 py-2.5 text-gray-400 text-xs">{{ formatDate(item.created_at) }}</td>
                    <td class="px-4 py-2.5 text-gray-300 font-mono text-xs">{{ item.model }}</td>
                    <td class="px-4 py-2.5 text-right tabular-nums text-gray-400">{{ (item.input_tokens ?? 0).toLocaleString() }}</td>
                    <td class="px-4 py-2.5 text-right tabular-nums text-gray-400">{{ (item.output_tokens ?? 0).toLocaleString() }}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <!-- 最近付款记录 -->
          <div class="space-y-3">
            <div class="flex items-center justify-between">
              <h2 class="font-semibold text-gray-300">付款记录</h2>
              <RouterLink to="/orders" class="text-xs text-violet-400 hover:text-violet-300 transition-colors">
                全部订单 →
              </RouterLink>
            </div>
            <div class="rounded-xl border border-gray-800 bg-gray-900/40 overflow-hidden">
              <div v-if="loading" class="py-8 text-center text-gray-600 text-sm">加载中…</div>
              <div v-else-if="recentOrders.length === 0" class="py-8 text-center text-gray-600 text-sm">
                暂无付款记录
              </div>
              <table v-else class="w-full text-sm">
                <thead class="border-b border-gray-800">
                  <tr class="text-xs text-gray-500">
                    <th class="px-4 py-2.5 text-left">时间</th>
                    <th class="px-4 py-2.5 text-left">套餐</th>
                    <th class="px-4 py-2.5 text-right">金额</th>
                    <th class="px-4 py-2.5 text-right">状态</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-800/60">
                  <tr v-for="order in recentOrders" :key="order.id" class="hover:bg-gray-800/30">
                    <td class="px-4 py-2.5 text-gray-400 text-xs">{{ formatDate(order.created_at) }}</td>
                    <td class="px-4 py-2.5 text-gray-300 text-xs">{{ order.description ?? '订阅套餐' }}</td>
                    <td class="px-4 py-2.5 text-right font-medium text-white">¥{{ order.amount?.toFixed(2) }}</td>
                    <td class="px-4 py-2.5 text-right">
                      <span class="rounded-full px-2 py-0.5 text-xs font-medium"
                        :class="{
                          'bg-green-900/40 text-green-400': order.status === 'completed',
                          'bg-yellow-900/40 text-yellow-400': order.status === 'pending',
                          'bg-red-900/40 text-red-400': order.status === 'failed' || order.status === 'expired',
                          'bg-gray-800 text-gray-400': !['completed','pending','failed','expired'].includes(order.status)
                        }">
                        {{ { completed: '已完成', pending: '待付款', failed: '失败', expired: '已过期' }[order.status] ?? order.status }}
                      </span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

        </div>
      </template>
    </main>
  </div>
</template>
