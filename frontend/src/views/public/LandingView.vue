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
  } catch { /* non-critical */ }
})

const features = [
  {
    icon: '✍️',
    title: 'AI 驱动的剧本创作',
    desc: '从大纲到分集，AI 理解剧本结构，精准补全对白、场景描述和人物弧线。',
  },
  {
    icon: '⚡',
    title: '多模型无缝切换',
    desc: '订阅一次，按需选用 Claude、GPT-4o、Gemini 等顶尖模型，积分制计费透明可控。',
  },
  {
    icon: '🔒',
    title: '本地优先，数据安全',
    desc: '剧本文件存储在你的本地机器上。AI 服务按需调用，创意资产永远属于你。',
  },
]

const steps = [
  { num: '01', title: '下载 StoryClaw', desc: '桌面客户端，Windows / macOS' },
  { num: '02', title: '订阅积分', desc: '选择套餐，微信 / 支付宝一键购买' },
  { num: '03', title: '一键授权', desc: '在应用内点击「登录服务」，自动完成配置' },
  { num: '04', title: '开始创作', desc: 'AI 即刻就绪，无需手动填写任何 API 参数' },
]

const models = ['Claude Opus 4', 'Claude Sonnet 4', 'GPT-4o', 'GPT-4o Mini', 'Gemini 2.5 Pro', 'Gemini 2.0 Flash']
</script>

<template>
  <LandingLayout>
    <!-- ── Hero ── -->
    <section class="relative overflow-hidden pt-36 pb-28 px-5 text-center">
      <div class="pointer-events-none absolute inset-0 bg-gradient-to-b from-violet-950/30 via-gray-950 to-gray-950" />
      <!-- Decorative glow -->
      <div class="pointer-events-none absolute top-0 left-1/2 -translate-x-1/2 w-[600px] h-[300px] rounded-full bg-violet-600/20 blur-3xl" />

      <div class="relative mx-auto max-w-3xl space-y-7">
        <div class="inline-flex items-center gap-2 rounded-full border border-violet-500/30 bg-violet-500/10 px-4 py-1.5 text-xs font-semibold text-violet-300">
          <span class="h-1.5 w-1.5 animate-pulse rounded-full bg-violet-400" />
          AI 驱动 · 本地优先 · 开箱即用
        </div>

        <h1 class="text-5xl font-extrabold leading-tight tracking-tight text-white md:text-6xl">
          专为创作者打造的<br />
          <span class="bg-gradient-to-r from-violet-400 to-fuchsia-400 bg-clip-text text-transparent">
            AI 剧本工作台
          </span>
        </h1>

        <p class="mx-auto max-w-xl text-xl leading-relaxed text-gray-400">
          StoryClaw 把顶尖 AI 模型带进你的剧本工作流。从大纲到完稿，
          积分制透明计费，一次订阅，多模型自由切换。
        </p>

        <div class="flex flex-wrap items-center justify-center gap-4 pt-2">
          <a
            href="#download"
            class="inline-flex items-center gap-2 rounded-xl bg-violet-600 px-7 py-3.5 font-semibold text-white shadow-lg shadow-violet-900/40 transition-colors hover:bg-violet-500"
          >
            免费下载 →
          </a>
          <RouterLink
            to="/portal"
            class="inline-flex items-center gap-2 rounded-xl border border-gray-700 px-5 py-3.5 font-medium text-gray-300 transition-colors hover:border-gray-500 hover:text-white"
          >
            查看套餐 / 登录
          </RouterLink>
        </div>
      </div>
    </section>

    <!-- ── Features ── -->
    <section id="features" class="py-20 px-5">
      <div class="mx-auto max-w-5xl">
        <h2 class="mb-12 text-center text-3xl font-bold text-white">为什么选择 StoryClaw</h2>
        <div class="grid gap-8 md:grid-cols-3">
          <div
            v-for="f in features" :key="f.title"
            class="rounded-2xl border border-gray-800 bg-gray-900/60 p-7 space-y-3 transition-shadow hover:shadow-lg hover:shadow-violet-900/20"
          >
            <div class="text-4xl leading-none">{{ f.icon }}</div>
            <h3 class="font-semibold text-white">{{ f.title }}</h3>
            <p class="text-sm leading-relaxed text-gray-400">{{ f.desc }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- ── Supported Models ── -->
    <section class="overflow-hidden border-y border-gray-800/60 bg-gray-900/30 py-14 px-5">
      <div class="mx-auto max-w-5xl">
        <p class="mb-6 text-center text-xs font-semibold uppercase tracking-widest text-gray-500">
          支持的 AI 模型
        </p>
        <div class="flex flex-wrap justify-center gap-3">
          <span
            v-for="m in models" :key="m"
            class="rounded-full border border-gray-700 bg-gray-800/60 px-4 py-1.5 text-sm font-medium text-gray-300"
          >
            {{ m }}
          </span>
        </div>
      </div>
    </section>

    <!-- ── Pricing ── -->
    <section id="pricing" class="py-20 px-5">
      <div class="mx-auto max-w-5xl">
        <h2 class="text-center text-3xl font-bold text-white mb-3">订阅套餐</h2>
        <p class="text-center text-gray-400 mb-14">购买即得积分，按用量消耗，清晰透明</p>

        <div v-if="plans.length === 0" class="grid gap-6 md:grid-cols-3">
          <div v-for="i in 3" :key="i" class="rounded-2xl border border-gray-800 p-6 space-y-4 animate-pulse">
            <div class="h-5 w-20 rounded bg-gray-800" />
            <div class="h-9 w-28 rounded bg-gray-800" />
            <div class="h-4 w-full rounded bg-gray-800" />
            <div class="h-10 rounded-xl bg-gray-800" />
          </div>
        </div>

        <div v-else class="grid gap-6 md:grid-cols-3">
          <div
            v-for="(plan, i) in plans.slice(0, 3)" :key="plan.id"
            class="relative rounded-2xl border-2 p-7 space-y-5 transition-shadow"
            :class="i === 1
              ? 'border-violet-500 shadow-xl shadow-violet-900/30'
              : 'border-gray-800 hover:border-gray-700'"
          >
            <div
              v-if="i === 1"
              class="absolute -top-3.5 left-1/2 -translate-x-1/2 rounded-full bg-violet-600 px-4 py-1 text-xs font-semibold text-white"
            >
              最受欢迎
            </div>
            <div>
              <h3 class="text-xl font-bold text-white">{{ plan.name }}</h3>
              <div class="mt-2 flex items-end gap-1">
                <span class="text-4xl font-extrabold text-white">¥{{ Math.round(plan.price) }}</span>
                <span class="mb-1 text-sm text-gray-400">/ {{ plan.validity_days }} 天</span>
              </div>
            </div>
            <div class="flex items-baseline gap-1.5">
              <span class="text-2xl font-bold text-violet-400">{{ plan.credits.toLocaleString() }}</span>
              <span class="text-sm text-gray-400">积分</span>
            </div>
            <ul class="space-y-2 text-sm text-gray-400">
              <li
                v-for="feat in plan.features.split('\n').filter(Boolean)" :key="feat"
                class="flex items-start gap-2"
              >
                <span class="mt-0.5 text-violet-400 shrink-0">✓</span>{{ feat }}
              </li>
              <li v-if="!plan.features" class="flex items-start gap-2">
                <span class="mt-0.5 text-violet-400">✓</span>多模型统一接入，积分透明计费
              </li>
            </ul>
            <RouterLink
              to="/portal"
              class="block rounded-xl py-2.5 text-center font-medium transition-colors"
              :class="i === 1
                ? 'bg-violet-600 text-white hover:bg-violet-500'
                : 'bg-gray-800 text-gray-300 hover:bg-gray-700'"
            >
              立即订阅
            </RouterLink>
          </div>
        </div>
      </div>
    </section>

    <!-- ── How it works ── -->
    <section id="how" class="border-t border-gray-800/60 bg-gray-900/30 py-20 px-5">
      <div class="mx-auto max-w-4xl">
        <h2 class="mb-14 text-center text-3xl font-bold text-white">4 步开始创作</h2>
        <div class="grid gap-8 md:grid-cols-4">
          <div v-for="step in steps" :key="step.num" class="text-center space-y-3">
            <div
              class="mx-auto flex h-14 w-14 items-center justify-center rounded-2xl bg-violet-600/20 text-xl font-bold text-violet-400"
            >
              {{ step.num }}
            </div>
            <h4 class="font-semibold text-white">{{ step.title }}</h4>
            <p class="text-sm text-gray-400">{{ step.desc }}</p>
          </div>
        </div>
        <div class="mt-14 text-center">
          <a
            id="download"
            href="#"
            class="inline-flex items-center gap-2 rounded-xl bg-violet-600 px-8 py-3.5 font-semibold text-white shadow-lg shadow-violet-900/40 transition-colors hover:bg-violet-500"
          >
            下载 StoryClaw →
          </a>
        </div>
      </div>
    </section>
  </LandingLayout>
</template>
