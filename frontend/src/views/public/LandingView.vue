<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { RouterLink } from 'vue-router'
import LandingLayout from '@/layouts/LandingLayout.vue'
import { getPublicPlans, type PlanItem } from '@/api/credits'

const plans = ref<PlanItem[]>([])

onMounted(async () => {
  try {
    const res = await getPublicPlans()
    plans.value = (res.items ?? []).filter((p) => p.for_sale).sort((a, b) => a.sort_order - b.sort_order)
  } catch {
    // non-critical
  }
})

const features = [
  {
    icon: '🤖',
    title: '多模型统一接入',
    desc: 'Claude、GPT-4o、Gemini 一个平台全覆盖。一个账号、一个 API Key，搞定所有模型。',
  },
  {
    icon: '💎',
    title: '积分透明计费',
    desc: '订阅即获积分，按实际调用消耗，余额实时可见。不再为账单焦虑，清晰掌控用量。',
  },
  {
    icon: '⚡',
    title: 'StoryClaw 一键配置',
    desc: '在 StoryClaw 中输入账号一键登录，模型列表自动加载，无需手动填写任何 API 参数。',
  },
]

const models = [
  'Claude Opus 4',
  'Claude Sonnet 4',
  'Claude Haiku 4',
  'GPT-4o',
  'GPT-4o Mini',
  'Gemini 2.5 Pro',
  'Gemini 2.0 Flash',
]

const steps = [
  { num: '01', title: '注册账号', desc: '邮箱注册，30 秒完成' },
  { num: '02', title: '选择套餐', desc: '微信 / 支付宝一键购买积分' },
  { num: '03', title: 'StoryClaw 连接', desc: '输入服务地址，一键登录' },
  { num: '04', title: '开始使用', desc: '模型自动配置，立即可用' },
]

// Recommend the middle plan (index 1)
const recommendedPlanIndex = 1
</script>

<template>
  <LandingLayout>
    <!-- ========== Hero ========== -->
    <section class="relative overflow-hidden pt-32 pb-24 px-5 text-center">
      <!-- Background gradient -->
      <div class="pointer-events-none absolute inset-0 bg-gradient-to-b from-indigo-50/70 via-white to-white" />

      <div class="relative mx-auto max-w-3xl space-y-7">
        <div class="inline-flex items-center gap-2 rounded-full bg-indigo-50 px-4 py-1.5 text-xs font-semibold text-indigo-700">
          <span class="h-1.5 w-1.5 animate-pulse rounded-full bg-indigo-500" />
          开源 · 自托管 · 积分透明计费
        </div>

        <h1 class="text-5xl font-extrabold leading-tight tracking-tight text-gray-900 md:text-6xl">
          一个订阅，<br />
          <span class="bg-gradient-to-r from-indigo-600 to-violet-600 bg-clip-text text-transparent">
            接入所有 AI 模型
          </span>
        </h1>

        <p class="mx-auto max-w-xl text-xl leading-relaxed text-gray-500">
          基于 sub2api 构建的 API 中转平台。积分透明计费，StoryClaw 一键配置，
          拒绝繁琐的 API Key 管理。
        </p>

        <div class="flex flex-wrap items-center justify-center gap-4 pt-2">
          <RouterLink
            to="/register"
            class="inline-flex items-center gap-2 rounded-xl bg-indigo-600 px-7 py-3 font-semibold text-white shadow-lg shadow-indigo-200 transition-colors hover:bg-indigo-700"
          >
            免费开始使用 →
          </RouterLink>
          <a
            href="#pricing"
            class="inline-flex items-center gap-2 rounded-xl px-5 py-3 font-medium text-gray-600 transition-colors hover:text-gray-900"
          >
            查看定价
          </a>
        </div>
      </div>
    </section>

    <!-- ========== Features ========== -->
    <section id="features" class="py-20 px-5">
      <div class="mx-auto max-w-5xl">
        <h2 class="mb-12 text-center text-3xl font-bold text-gray-900">为什么选择我们</h2>
        <div class="grid gap-8 md:grid-cols-3">
          <div
            v-for="f in features" :key="f.title"
            class="rounded-2xl border border-gray-100 bg-gray-50 p-7 space-y-3 transition-shadow hover:shadow-md"
          >
            <div class="text-4xl leading-none">{{ f.icon }}</div>
            <h3 class="font-semibold text-gray-900">{{ f.title }}</h3>
            <p class="text-sm leading-relaxed text-gray-500">{{ f.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- ========== Supported Models ========== -->
    <section class="overflow-hidden bg-gray-50 py-14 px-5">
      <div class="mx-auto max-w-5xl">
        <p class="mb-6 text-center text-xs font-semibold uppercase tracking-widest text-gray-400">
          支持的 AI 模型
        </p>
        <div class="flex flex-wrap justify-center gap-3">
          <span
            v-for="m in models" :key="m"
            class="rounded-full border border-gray-200 bg-white px-4 py-1.5 text-sm font-medium text-gray-700 shadow-sm"
          >
            {{ m }}
          </span>
        </div>
      </div>
    </section>

    <!-- ========== Pricing ========== -->
    <section id="pricing" class="py-20 px-5">
      <div class="mx-auto max-w-5xl">
        <h2 class="text-center text-3xl font-bold text-gray-900 mb-3">透明定价</h2>
        <p class="text-center text-gray-500 mb-14">购买即得积分，按用量消耗，到期前不浪费</p>

        <!-- Placeholder when no plans loaded -->
        <div v-if="plans.length === 0" class="grid gap-6 md:grid-cols-3">
          <div v-for="i in 3" :key="i" class="rounded-2xl border-2 border-gray-100 p-6 space-y-4 animate-pulse">
            <div class="h-5 w-20 rounded bg-gray-100" />
            <div class="h-9 w-28 rounded bg-gray-100" />
            <div class="h-4 w-full rounded bg-gray-100" />
            <div class="h-4 w-3/4 rounded bg-gray-100" />
            <div class="h-10 rounded-xl bg-gray-100" />
          </div>
        </div>

        <div v-else class="grid gap-6 md:grid-cols-3">
          <div
            v-for="(plan, i) in plans.slice(0, 3)" :key="plan.id"
            class="relative rounded-2xl border-2 p-7 space-y-5 transition-shadow"
            :class="i === recommendedPlanIndex
              ? 'border-indigo-500 shadow-xl shadow-indigo-100'
              : 'border-gray-200 hover:shadow-md'"
          >
            <div
              v-if="i === recommendedPlanIndex"
              class="absolute -top-3.5 left-1/2 -translate-x-1/2 rounded-full bg-indigo-600 px-4 py-1 text-xs font-semibold text-white"
            >
              最受欢迎
            </div>

            <div>
              <h3 class="text-xl font-bold text-gray-900">{{ plan.name }}</h3>
              <div class="mt-2 flex items-end gap-1">
                <span class="text-4xl font-extrabold text-gray-900">¥{{ Math.round(plan.price) }}</span>
                <span class="mb-1 text-sm text-gray-400">/ {{ plan.validity_days }} 天</span>
              </div>
              <div v-if="plan.original_price" class="mt-0.5 text-xs text-gray-400 line-through">
                原价 ¥{{ Math.round(plan.original_price) }}
              </div>
            </div>

            <div class="flex items-baseline gap-1.5">
              <span class="text-2xl font-bold text-indigo-700">{{ plan.credits.toLocaleString() }}</span>
              <span class="text-sm text-gray-400">积分</span>
            </div>

            <ul class="space-y-2 text-sm text-gray-600">
              <li
                v-for="feat in plan.features.split('\n').filter(Boolean)" :key="feat"
                class="flex items-start gap-2"
              >
                <span class="mt-0.5 text-green-500 shrink-0">✓</span>{{ feat }}
              </li>
              <li v-if="!plan.features" class="text-gray-400">OpenAI 兼容接口，多模型统一接入</li>
            </ul>

            <RouterLink
              to="/register"
              class="block rounded-xl py-2.5 text-center font-medium transition-colors"
              :class="i === recommendedPlanIndex
                ? 'bg-indigo-600 text-white hover:bg-indigo-700'
                : 'bg-gray-100 text-gray-700 hover:bg-gray-200'"
            >
              立即购买
            </RouterLink>
          </div>
        </div>
      </div>
    </section>

    <!-- ========== How it works ========== -->
    <section id="how" class="bg-gray-50 py-20 px-5">
      <div class="mx-auto max-w-4xl">
        <h2 class="mb-14 text-center text-3xl font-bold text-gray-900">4 步开始使用</h2>
        <div class="grid gap-8 md:grid-cols-4">
          <div v-for="step in steps" :key="step.num" class="text-center space-y-3">
            <div
              class="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-indigo-100 text-xl font-bold text-indigo-700"
            >
              {{ step.num }}
            </div>
            <h4 class="font-semibold text-gray-900">{{ step.title }}</h4>
            <p class="text-sm text-gray-500">{{ step.desc }}</p>
          </div>
        </div>

        <div class="mt-14 text-center">
          <RouterLink
            to="/register"
            class="inline-flex items-center gap-2 rounded-xl bg-indigo-600 px-8 py-3.5 font-semibold text-white shadow-lg shadow-indigo-200 transition-colors hover:bg-indigo-700"
          >
            立即开始 →
          </RouterLink>
        </div>
      </div>
    </section>
  </LandingLayout>
</template>
