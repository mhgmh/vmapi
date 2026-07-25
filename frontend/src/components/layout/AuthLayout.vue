<template>
  <div class="flex min-h-screen bg-white dark:bg-dark-950">
    <!-- Brand panel (desktop only) -->
    <aside
      class="relative hidden flex-[1.1] flex-col justify-between overflow-hidden border-r border-gray-100 bg-gray-50 p-12 dark:border-dark-700 dark:bg-dark-900 lg:flex"
    >
      <!-- Grid texture -->
      <div
        class="pointer-events-none absolute inset-0 bg-[linear-gradient(rgba(120,120,128,0.05)_1px,transparent_1px),linear-gradient(90deg,rgba(120,120,128,0.05)_1px,transparent_1px)] bg-[size:56px_56px] [mask-image:radial-gradient(ellipse_90%_80%_at_30%_20%,black_30%,transparent_75%)]"
      ></div>

      <!-- Logo -->
      <div class="relative">
        <div class="flex items-center gap-3">
          <div class="flex h-8 w-8 items-center justify-center overflow-hidden rounded-lg">
            <img :src="siteLogo || '/logo.svg'" alt="Logo" class="h-full w-full object-contain" />
          </div>
          <span class="text-[15px] font-semibold text-gray-900 dark:text-dark-50">{{ siteName }}</span>
        </div>
      </div>

      <!-- Pitch + sample request -->
      <div class="relative max-w-md">
        <p class="font-mono text-xs font-medium tracking-wide text-primary-600 dark:text-primary-400">
          {{ t('auth.brandEyebrow') }}
        </p>
        <h2 class="mt-3.5 text-3xl font-semibold leading-[1.2] tracking-tight text-gray-900 dark:text-dark-50">
          {{ t('auth.brandHeadline') }}
        </h2>
        <p class="mt-3 max-w-sm text-sm leading-relaxed text-gray-500 dark:text-dark-400">
          {{ t('auth.brandSubline') }}
        </p>

        <div class="mt-6 max-w-md overflow-hidden rounded-lg border border-gray-200 bg-white dark:border-dark-700 dark:bg-dark-950">
          <div class="flex items-center gap-2 border-b border-gray-100 px-3 py-2 dark:border-dark-700">
            <span class="flex gap-1.5">
              <i class="block h-2 w-2 rounded-full bg-gray-200 dark:bg-dark-600"></i>
              <i class="block h-2 w-2 rounded-full bg-gray-200 dark:bg-dark-600"></i>
              <i class="block h-2 w-2 rounded-full bg-gray-200 dark:bg-dark-600"></i>
            </span>
            <span class="font-mono text-[11px] text-gray-400 dark:text-dark-500">curl — first request</span>
          </div>
          <pre class="overflow-x-auto p-4 font-mono text-xs leading-relaxed text-gray-600 dark:text-dark-300"><span class="text-gray-400 dark:text-dark-500"># Point your client at our base URL — that's it.</span>
curl {{ apiBaseUrl }}/chat/completions \
  -H <span class="text-primary-600 dark:text-primary-400">"Authorization: Bearer sk-…"</span> \
  -d <span class="text-primary-600 dark:text-primary-400">'{
    "model": "claude-sonnet-4.5",
    "messages": [{"role":"user","content":"Hello"}]
  }'</span></pre>
        </div>
      </div>

      <!-- Footer meta -->
      <div class="relative flex items-center gap-3 text-xs text-gray-400 dark:text-dark-500">
        <span>&copy; {{ currentYear }} {{ siteName }}</span>
      </div>
    </aside>

    <!-- Form column -->
    <main class="relative flex flex-1 items-center justify-center px-4 py-10 sm:px-8">
      <!-- Mobile brand -->
      <div class="absolute left-0 right-0 top-0 flex items-center justify-between p-5 lg:hidden">
        <div class="flex items-center gap-2.5">
          <div class="flex h-7 w-7 items-center justify-center overflow-hidden rounded-md">
            <img :src="siteLogo || '/logo.svg'" alt="Logo" class="h-full w-full object-contain" />
          </div>
          <span class="text-sm font-semibold text-gray-900 dark:text-dark-50">{{ siteName }}</span>
        </div>
        <LocaleSwitcher />
      </div>

      <div class="w-full max-w-[360px]">
        <slot />

        <!-- Footer Links -->
        <div class="mt-6 text-center text-sm text-gray-500 dark:text-dark-400">
          <slot name="footer" />
        </div>
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores'
import { sanitizeUrl } from '@/utils/url'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'

const { t } = useI18n()
const appStore = useAppStore()

const siteName = computed(() => appStore.siteName || 'Sub2API')
const siteLogo = computed(() => sanitizeUrl(appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))

const apiBaseUrl = computed(() => {
  if (typeof window !== 'undefined' && window.location?.origin) {
    return `${window.location.origin}/v1`
  }
  return 'https://api.example.com/v1'
})

const currentYear = computed(() => new Date().getFullYear())

onMounted(() => {
  appStore.fetchPublicSettings()
})
</script>
