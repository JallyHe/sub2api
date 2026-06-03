<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { RouterLink } from 'vue-router'

const scrolled = ref(false)

function onScroll() {
  scrolled.value = window.scrollY > 24
}

onMounted(() => window.addEventListener('scroll', onScroll, { passive: true }))
onUnmounted(() => window.removeEventListener('scroll', onScroll))
</script>

<template>
  <div class="min-h-screen bg-white text-gray-900">
    <!-- Sticky header -->
    <header
      class="fixed inset-x-0 top-0 z-50 h-16 transition-all duration-200"
      :class="scrolled ? 'bg-white/95 backdrop-blur-md shadow-sm' : 'bg-transparent'"
    >
      <nav class="mx-auto flex h-full max-w-6xl items-center justify-between px-5">
        <!-- Logo -->
        <RouterLink to="/" class="flex items-center gap-2 font-bold text-lg">
          <span class="text-indigo-600">◆</span>
          <span class="text-gray-900">Sub2API</span>
        </RouterLink>

        <!-- Nav links (hidden on mobile) -->
        <div class="hidden items-center gap-7 text-sm text-gray-600 md:flex">
          <a href="#features" class="hover:text-gray-900 transition-colors">功能</a>
          <a href="#pricing" class="hover:text-gray-900 transition-colors">定价</a>
          <a href="#how" class="hover:text-gray-900 transition-colors">使用方法</a>
        </div>

        <!-- CTA buttons -->
        <div class="flex items-center gap-3">
          <RouterLink
            to="/login"
            class="text-sm text-gray-600 hover:text-gray-900 px-3 py-1.5 transition-colors"
          >
            登录
          </RouterLink>
          <RouterLink
            to="/register"
            class="rounded-lg bg-indigo-600 px-4 py-1.5 text-sm font-medium text-white hover:bg-indigo-700 transition-colors shadow-sm"
          >
            免费注册
          </RouterLink>
        </div>
      </nav>
    </header>

    <!-- Main content -->
    <main>
      <slot />
    </main>

    <!-- Footer -->
    <footer class="border-t border-gray-100 py-10 mt-20">
      <div class="mx-auto max-w-6xl px-5 text-center space-y-3">
        <p class="text-sm text-gray-400">
          © {{ new Date().getFullYear() }} Sub2API &nbsp;·&nbsp; 开源 AI API 中转订阅平台
        </p>
        <div class="flex justify-center gap-5 text-xs text-gray-400">
          <RouterLink to="/legal/terms" class="hover:text-gray-600 transition-colors">服务条款</RouterLink>
          <RouterLink to="/legal/privacy" class="hover:text-gray-600 transition-colors">隐私政策</RouterLink>
        </div>
      </div>
    </footer>
  </div>
</template>
