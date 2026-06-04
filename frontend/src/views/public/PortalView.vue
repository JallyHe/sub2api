<script setup lang="ts">
/**
 * 用户门户 /portal
 * 仅展示：登录入口 / 套餐购买 / 积分余额 / 使用记录
 * 不需要完整 sub2api 管理后台
 */
import { computed } from 'vue'
import { RouterLink, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const auth = useAuthStore()
const router = useRouter()
const isLoggedIn = computed(() => auth.isAuthenticated)

function goLogin() { void router.push('/login?redirect=/portal') }
function goRegister() { void router.push('/register') }
</script>

<template>
  <div class="min-h-screen bg-gray-950 text-white flex flex-col">
    <!-- Header -->
    <header class="border-b border-gray-800/60 bg-gray-950/80 backdrop-blur sticky top-0 z-50">
      <div class="mx-auto flex h-14 max-w-4xl items-center justify-between px-5">
        <RouterLink to="/" class="flex items-center gap-2 font-bold">
          <span class="flex h-7 w-7 items-center justify-center rounded-lg bg-gradient-to-br from-violet-600 to-fuchsia-600 text-xs font-black">S</span>
          <span>StoryClaw</span>
        </RouterLink>
        <div v-if="isLoggedIn" class="flex items-center gap-3 text-sm">
          <RouterLink to="/usage" class="text-gray-400 hover:text-white transition-colors">使用记录</RouterLink>
          <RouterLink to="/subscriptions" class="text-gray-400 hover:text-white transition-colors">我的订阅</RouterLink>
          <RouterLink to="/purchase" class="rounded-lg bg-violet-600 px-3.5 py-1.5 text-sm font-medium hover:bg-violet-500 transition-colors">购买套餐</RouterLink>
        </div>
      </div>
    </header>

    <!-- Main -->
    <main class="flex-1 mx-auto w-full max-w-4xl px-5 py-12">

      <!-- 未登录：登录/注册入口 -->
      <template v-if="!isLoggedIn">
        <div class="text-center space-y-6 pt-12">
          <div class="flex h-16 w-16 mx-auto items-center justify-center rounded-2xl bg-gradient-to-br from-violet-600 to-fuchsia-600 text-3xl font-black">S</div>
          <h1 class="text-3xl font-bold">欢迎使用 StoryClaw</h1>
          <p class="text-gray-400 max-w-sm mx-auto leading-relaxed">
            登录后即可查看积分余额、购买套餐，以及通过 StoryClaw 桌面应用一键完成授权。
          </p>
          <div class="flex justify-center gap-4 pt-2">
            <button
              @click="goLogin"
              class="rounded-xl bg-violet-600 px-7 py-3 font-semibold hover:bg-violet-500 transition-colors"
            >
              登录账号
            </button>
            <button
              @click="goRegister"
              class="rounded-xl border border-gray-700 px-7 py-3 font-medium text-gray-300 hover:border-gray-500 transition-colors"
            >
              创建账号
            </button>
          </div>
        </div>
      </template>

      <!-- 已登录：用户门户 -->
      <template v-else>
        <div class="space-y-8">
          <h1 class="text-2xl font-bold">我的账户</h1>

          <!-- 积分卡片 -->
          <div class="grid gap-4 sm:grid-cols-2">
            <RouterLink
              to="/usage"
              class="group rounded-2xl border border-gray-800 bg-gray-900/60 p-6 hover:border-violet-500/40 hover:shadow-lg hover:shadow-violet-900/10 transition-all"
            >
              <div class="flex items-center justify-between mb-4">
                <span class="text-sm text-gray-400">积分余额</span>
                <span class="text-xs text-gray-600 group-hover:text-violet-400 transition-colors">查看明细 →</span>
              </div>
              <CreditBalanceInline />
            </RouterLink>

            <RouterLink
              to="/purchase"
              class="group rounded-2xl border border-gray-800 bg-gray-900/60 p-6 hover:border-violet-500/40 hover:shadow-lg hover:shadow-violet-900/10 transition-all flex flex-col justify-between"
            >
              <div>
                <span class="text-sm text-gray-400">购买积分</span>
                <p class="text-gray-500 text-sm mt-2 leading-relaxed">选择套餐补充积分，微信 / 支付宝支付</p>
              </div>
              <span class="mt-4 inline-flex items-center text-violet-400 text-sm font-medium group-hover:gap-2 transition-all">
                选择套餐 →
              </span>
            </RouterLink>
          </div>

          <!-- 快捷导航 -->
          <div class="grid gap-3 sm:grid-cols-3">
            <RouterLink
              v-for="item in quickLinks" :key="item.path"
              :to="item.path"
              class="rounded-xl border border-gray-800 bg-gray-900/40 px-5 py-4 text-sm font-medium text-gray-300 hover:border-gray-600 hover:text-white transition-all flex items-center gap-3"
            >
              <span>{{ item.icon }}</span>
              {{ item.label }}
            </RouterLink>
          </div>
        </div>
      </template>
    </main>
  </div>
</template>

<script lang="ts">
// CreditBalanceInline: inline component defined in same file for simplicity
import { defineComponent, ref, onMounted } from 'vue'
import { getBalance } from '@/api/credits'

const CreditBalanceInline = defineComponent({
  name: 'CreditBalanceInline',
  setup() {
    const balance = ref<number | null>(null)
    onMounted(async () => {
      try { balance.value = (await getBalance()).balance } catch { balance.value = 0 }
    })
    return { balance }
  },
  template: `
    <div>
      <div v-if="balance === null" class="h-9 w-24 rounded bg-gray-800 animate-pulse" />
      <div v-else class="text-3xl font-bold tabular-nums" :class="balance <= 0 ? 'text-red-400' : 'text-white'">
        {{ balance.toLocaleString() }}
        <span class="text-base font-normal text-gray-400 ml-1">积分</span>
      </div>
    </div>
  `
})

const quickLinks = [
  { icon: '📋', label: '使用记录', path: '/usage' },
  { icon: '💳', label: '我的订阅', path: '/subscriptions' },
  { icon: '📦', label: '我的订单', path: '/orders' },
]

export { CreditBalanceInline, quickLinks }
</script>
