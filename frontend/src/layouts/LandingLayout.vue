<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { RouterLink } from 'vue-router'

const scrolled = ref(false)
function onScroll() { scrolled.value = window.scrollY > 24 }
onMounted(() => window.addEventListener('scroll', onScroll, { passive: true }))
onUnmounted(() => window.removeEventListener('scroll', onScroll))
</script>

<template>
  <div class="min-h-screen bg-gray-950 text-white">
    <!-- Sticky header -->
    <header
      class="fixed inset-x-0 top-0 z-50 h-16 transition-all duration-200"
      :class="scrolled ? 'bg-gray-950/95 backdrop-blur-md border-b border-gray-800/60' : 'bg-transparent'"
    >
      <nav class="mx-auto flex h-full max-w-6xl items-center justify-between px-5">
        <!-- Logo -->
        <RouterLink to="/" class="flex items-center gap-2.5 font-bold text-lg">
          <span class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-violet-600 to-fuchsia-600 text-sm font-black text-white">S</span>
          <span class="text-white">StoryClaw</span>
        </RouterLink>

        <!-- Nav links -->
        <div class="hidden items-center gap-7 text-sm text-gray-400 md:flex">
          <a href="#features" class="hover:text-white transition-colors">功能</a>
          <a href="#pricing" class="hover:text-white transition-colors">定价</a>
          <a href="#how" class="hover:text-white transition-colors">使用方法</a>
          <a href="#download" class="hover:text-white transition-colors">下载</a>
        </div>

        <!-- CTA buttons -->
        <div class="flex items-center gap-3">
          <RouterLink
            to="/portal"
            class="text-sm text-gray-400 hover:text-white px-3 py-1.5 transition-colors"
          >
            登录
          </RouterLink>
          <RouterLink
            to="/portal"
            class="rounded-lg bg-violet-600 px-4 py-1.5 text-sm font-medium text-white hover:bg-violet-500 transition-colors"
          >
            订阅套餐
          </RouterLink>
        </div>
      </nav>
    </header>

    <main><slot /></main>

    <!-- Footer -->
    <footer class="border-t border-gray-800/60 py-10 mt-10">
      <div class="mx-auto max-w-6xl px-5 text-center space-y-3">
        <div class="flex items-center justify-center gap-2 mb-4">
          <span class="flex h-7 w-7 items-center justify-center rounded-lg bg-gradient-to-br from-violet-600 to-fuchsia-600 text-xs font-black text-white">S</span>
          <span class="font-semibold text-white">StoryClaw</span>
        </div>
        <p class="text-sm text-gray-500">
          © {{ new Date().getFullYear() }} StoryClaw &nbsp;·&nbsp; AI 驱动的剧本创作工具
        </p>
        <div class="flex justify-center gap-5 text-xs text-gray-600">
          <RouterLink to="/legal/terms" class="hover:text-gray-400 transition-colors">服务条款</RouterLink>
          <RouterLink to="/legal/privacy" class="hover:text-gray-400 transition-colors">隐私政策</RouterLink>
        </div>
      </div>
    </footer>
  </div>
</template>
