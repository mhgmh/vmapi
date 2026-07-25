<template>
  <!-- Custom Home Content: Full Page Mode -->
  <div v-if="homeContent" class="min-h-screen">
    <!-- iframe mode -->
    <iframe
      v-if="isHomeContentUrl"
      :src="homeContent.trim()"
      class="h-screen w-full border-0"
      allowfullscreen
    ></iframe>
    <!-- HTML mode - SECURITY: homeContent is admin-only setting, XSS risk is acceptable -->
    <div v-else v-html="homeContent"></div>
  </div>

  <!-- Default Home Page -->
  <div v-else class="flex min-h-screen flex-col bg-white dark:bg-dark-950">
    <!-- Nav -->
    <header
      class="sticky top-0 z-20 border-b border-gray-100 bg-white/80 backdrop-blur-md dark:border-dark-700 dark:bg-dark-950/80"
    >
      <nav class="mx-auto flex h-14 max-w-6xl items-center justify-between px-6">
        <div class="flex items-center gap-2.5">
          <div class="h-7 w-7 overflow-hidden rounded-md">
            <img :src="siteLogo || '/logo.svg'" alt="Logo" class="h-full w-full object-contain" />
          </div>
          <span class="text-sm font-semibold text-gray-900 dark:text-dark-50">{{ siteName }}</span>
        </div>

        <div class="flex items-center gap-1.5">
          <LocaleSwitcher />
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="hidden items-center gap-1.5 rounded-md px-2.5 py-1.5 text-sm font-medium text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-dark-50 sm:flex"
          >
            <Icon name="book" size="sm" />
            {{ t('home.docs') }}
          </a>
          <button
            @click="toggleTheme"
            class="rounded-md p-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-dark-50"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
          >
            <Icon v-if="isDark" name="sun" size="md" />
            <Icon v-else name="moon" size="md" />
          </button>
          <router-link
            v-if="isAuthenticated"
            :to="dashboardPath"
            class="btn btn-primary btn-sm ml-1"
          >
            {{ t('home.dashboard') }}
          </router-link>
          <template v-else>
            <router-link
              to="/login"
              class="ml-1 rounded-md px-2.5 py-1.5 text-sm font-medium text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-dark-50"
            >
              {{ t('home.login') }}
            </router-link>
            <router-link to="/register" class="btn btn-primary btn-sm">
              {{ t('home.getStarted') }}
            </router-link>
          </template>
        </div>
      </nav>
    </header>

    <main class="flex-1">
      <!-- Hero -->
      <section class="mx-auto grid max-w-6xl items-center gap-12 px-6 pb-20 pt-16 lg:grid-cols-2 lg:pt-24">
        <div>
          <span
            class="inline-flex items-center gap-2 rounded-full border border-primary-200 bg-primary-50 px-3 py-1 text-xs font-medium text-primary-700 dark:border-primary-400/25 dark:bg-primary-400/10 dark:text-primary-400"
          >
            <span class="h-1.5 w-1.5 rounded-full bg-primary-500 dark:bg-primary-400"></span>
            {{ t('home.tags.subscriptionToApi') }} · {{ t('home.tags.realtimeBilling') }}
          </span>
          <h1
            class="mt-5 text-4xl font-semibold leading-[1.1] tracking-tight text-gray-900 dark:text-dark-50 md:text-5xl"
          >
            {{ siteSubtitle }}
          </h1>
          <p class="mt-4 max-w-lg text-base leading-relaxed text-gray-500 dark:text-dark-400">
            {{ t('home.heroDescription') }}
          </p>
          <div class="mt-8 flex flex-wrap gap-3">
            <router-link
              :to="isAuthenticated ? dashboardPath : '/register'"
              class="btn btn-primary h-11 px-5 text-sm"
            >
              {{ isAuthenticated ? t('home.goToDashboard') : t('home.getStarted') }}
              <Icon name="arrowRight" size="md" :stroke-width="2" />
            </router-link>
            <a
              v-if="docUrl"
              :href="docUrl"
              target="_blank"
              rel="noopener noreferrer"
              class="btn btn-secondary h-11 px-5 text-sm"
            >
              {{ t('home.viewDocs') }}
            </a>
          </div>
        </div>

        <!-- Code block: migrate in 3 lines -->
        <div>
          <div
            class="overflow-hidden rounded-lg border border-gray-200 bg-gray-50 dark:border-dark-700 dark:bg-dark-900"
          >
            <div
              class="flex items-center gap-2 border-b border-gray-100 px-3.5 py-2.5 dark:border-dark-700"
            >
              <span class="flex gap-1.5">
                <i class="block h-2.5 w-2.5 rounded-full bg-gray-200 dark:bg-dark-600"></i>
                <i class="block h-2.5 w-2.5 rounded-full bg-gray-200 dark:bg-dark-600"></i>
                <i class="block h-2.5 w-2.5 rounded-full bg-gray-200 dark:bg-dark-600"></i>
              </span>
              <span class="font-mono text-[11px] text-gray-400 dark:text-dark-500">
                Python — drop-in compatible
              </span>
            </div>
            <pre
              class="overflow-x-auto p-4 font-mono text-[12.5px] leading-relaxed text-gray-700 dark:text-dark-200"
            ><code><span class="text-blue-600 dark:text-blue-400">from</span> openai <span class="text-blue-600 dark:text-blue-400">import</span> OpenAI

client = OpenAI(
    base_url=<span class="text-primary-600 dark:text-primary-400">"{{ apiBaseUrl }}"</span>,
    api_key=<span class="text-primary-600 dark:text-primary-400">"sk-…"</span>,
)

resp = client.chat.completions.create(
    model=<span class="text-primary-600 dark:text-primary-400">"claude-sonnet-4.5"</span>,
    messages=[{<span class="text-primary-600 dark:text-primary-400">"role"</span>: <span class="text-primary-600 dark:text-primary-400">"user"</span>, <span class="text-primary-600 dark:text-primary-400">"content"</span>: <span class="text-primary-600 dark:text-primary-400">"Hello"</span>}],
)</code></pre>
          </div>
          <div class="mt-4 flex flex-wrap gap-2">
            <span
              v-for="tag in [t('home.tags.subscriptionToApi'), t('home.tags.stickySession'), t('home.tags.realtimeBilling')]"
              :key="tag"
              class="rounded-full border border-gray-200 px-3 py-1 text-xs font-medium text-gray-500 dark:border-dark-700 dark:text-dark-400"
            >
              {{ tag }}
            </span>
          </div>
        </div>
      </section>

      <!-- Features -->
      <section class="border-t border-gray-100 dark:border-dark-700">
        <div class="mx-auto max-w-6xl px-6 py-16">
          <p
            class="text-xs font-semibold uppercase tracking-widest text-primary-600 dark:text-primary-400"
          >
            {{ t('home.sections.howItWorks') }}
          </p>
          <h2 class="mt-3 text-2xl font-semibold tracking-tight text-gray-900 dark:text-dark-50 md:text-3xl">
            {{ t('home.solutions.subtitle') }}
          </h2>
          <div class="mt-10 grid gap-4 md:grid-cols-3">
            <div
              v-for="(f, i) in features"
              :key="f.title"
              class="rounded-lg border border-gray-200 bg-white p-6 transition-colors hover:border-gray-300 dark:border-dark-700 dark:bg-dark-800 dark:hover:border-dark-500"
            >
              <span class="font-mono text-xs font-medium text-primary-600 dark:text-primary-400">
                0{{ i + 1 }}
              </span>
              <h3 class="mt-3 text-[15px] font-semibold text-gray-900 dark:text-dark-50">
                {{ f.title }}
              </h3>
              <p class="mt-1.5 text-sm leading-relaxed text-gray-500 dark:text-dark-400">
                {{ f.desc }}
              </p>
            </div>
          </div>
        </div>
      </section>

      <!-- Comparison -->
      <section class="border-t border-gray-100 dark:border-dark-700">
        <div class="mx-auto max-w-6xl px-6 py-16">
          <h2 class="text-2xl font-semibold tracking-tight text-gray-900 dark:text-dark-50 md:text-3xl">
            {{ t('home.comparison.title') }}
          </h2>
          <div
            class="mt-8 overflow-x-auto rounded-lg border border-gray-200 dark:border-dark-700"
          >
            <table class="table">
              <thead>
                <tr>
                  <th>{{ t('home.comparison.headers.feature') }}</th>
                  <th>{{ t('home.comparison.headers.official') }}</th>
                  <th class="text-primary-600 dark:text-primary-400">
                    {{ siteName }}
                  </th>
                </tr>
              </thead>
              <tbody>
                <tr v-for="row in comparisonRows" :key="row.feature">
                  <td class="font-medium text-gray-900 dark:text-dark-100">{{ row.feature }}</td>
                  <td>{{ row.official }}</td>
                  <td class="text-gray-900 dark:text-dark-100">{{ row.us }}</td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </section>

      <!-- Providers -->
      <section class="border-t border-gray-100 dark:border-dark-700">
        <div class="mx-auto max-w-6xl px-6 py-16">
          <p
            class="text-xs font-semibold uppercase tracking-widest text-primary-600 dark:text-primary-400"
          >
            {{ t('home.providers.title') }}
          </p>
          <h2 class="mt-3 text-2xl font-semibold tracking-tight text-gray-900 dark:text-dark-50 md:text-3xl">
            {{ t('home.providers.description') }}
          </h2>
          <div class="mt-8 flex flex-wrap gap-2.5">
            <span
              v-for="p in providers"
              :key="p"
              class="inline-flex items-center gap-2 rounded-full border border-gray-200 bg-white px-4 py-2 text-sm font-medium text-gray-600 dark:border-dark-700 dark:bg-dark-800 dark:text-dark-300"
            >
              <span class="h-1.5 w-1.5 rounded-full bg-primary-500 dark:bg-primary-400"></span>
              {{ p }}
            </span>
          </div>
        </div>
      </section>

      <!-- CTA -->
      <section class="border-t border-gray-100 dark:border-dark-700">
        <div class="mx-auto max-w-6xl px-6 py-20 text-center">
          <h2 class="text-2xl font-semibold tracking-tight text-gray-900 dark:text-dark-50 md:text-3xl">
            {{ t('home.cta.title') }}
          </h2>
          <p class="mx-auto mt-3 max-w-md text-sm text-gray-500 dark:text-dark-400">
            {{ t('home.cta.description') }}
          </p>
          <router-link
            :to="isAuthenticated ? dashboardPath : '/register'"
            class="btn btn-primary mt-8 h-11 px-6 text-sm"
          >
            {{ isAuthenticated ? t('home.goToDashboard') : t('home.cta.button') }}
          </router-link>
        </div>
      </section>
    </main>

    <!-- Footer -->
    <footer class="border-t border-gray-100 dark:border-dark-700">
      <div
        class="mx-auto flex max-w-6xl flex-col items-center justify-between gap-3 px-6 py-8 text-sm text-gray-400 dark:text-dark-500 sm:flex-row"
      >
        <p>&copy; {{ currentYear }} {{ siteName }}. {{ t('home.footer.allRightsReserved') }}</p>
        <div class="flex items-center gap-5">
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="transition-colors hover:text-gray-700 dark:hover:text-dark-200"
          >
            {{ t('home.docs') }}
          </a>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'
import { sanitizeUrl } from '@/utils/url'

const { t } = useI18n()

const authStore = useAuthStore()
const appStore = useAppStore()

// Site settings - directly from appStore (already initialized from injected config)
const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const siteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle || 'One API key. Every model.')
const docUrl = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.doc_url || appStore.docUrl || ''))
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')

// Check if homeContent is a URL (for iframe display)
const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

const apiBaseUrl = computed(() => {
  if (typeof window !== 'undefined' && window.location?.origin) {
    return `${window.location.origin}/v1`
  }
  return 'https://api.example.com/v1'
})

const features = computed(() => [
  { title: t('home.features.unifiedGateway'), desc: t('home.features.unifiedGatewayDesc') },
  { title: t('home.features.multiAccount'), desc: t('home.features.multiAccountDesc') },
  { title: t('home.features.balanceQuota'), desc: t('home.features.balanceQuotaDesc') }
])

const comparisonRows = computed(() => {
  const keys = ['pricing', 'models', 'management', 'stability', 'control']
  return keys.map((k) => ({
    feature: t(`home.comparison.items.${k}.feature`),
    official: t(`home.comparison.items.${k}.official`),
    us: t(`home.comparison.items.${k}.us`)
  }))
})

const providers = computed(() => [
  'Claude',
  'GPT',
  'Gemini',
  t('home.providers.antigravity'),
  `${t('home.providers.more')}…`
])

// Theme
const isDark = ref(document.documentElement.classList.contains('dark'))

// Auth state
const isAuthenticated = computed(() => authStore.isAuthenticated)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')

// Current year for footer
const currentYear = computed(() => new Date().getFullYear())

// Toggle theme
function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

onMounted(() => {
  isDark.value = document.documentElement.classList.contains('dark')

  // Check auth state
  authStore.checkAuth()

  // Ensure public settings are loaded (will use cache if already loaded from injected config)
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>
