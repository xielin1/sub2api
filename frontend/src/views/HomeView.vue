<template>
  <!-- Custom Home Content: Full Page Mode -->
  <div v-if="hasHomeContent" class="min-h-screen">
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

  <!-- Compact Home Page -->
  <div
    v-else-if="compactHomeEnabled"
    data-testid="compact-home"
    class="flex min-h-screen flex-col bg-gray-50 text-gray-900 dark:bg-dark-950 dark:text-white"
  >
    <header class="border-b border-gray-200 px-4 py-4 sm:px-6 dark:border-dark-800">
      <nav class="mx-auto flex max-w-5xl flex-wrap items-center justify-between gap-3 sm:gap-4">
        <div class="flex min-w-0 flex-1 items-center gap-3">
          <img
            :src="siteLogo || '/logo.svg'"
            alt="Logo"
            class="h-9 w-9 shrink-0 rounded-lg object-contain"
          />
          <span class="min-w-0 truncate text-base font-semibold">{{ siteName }}</span>
        </div>
        <div class="flex max-w-full shrink-0 flex-wrap items-center justify-end gap-2">
          <LocaleSwitcher />
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg text-gray-500 hover:bg-gray-100 dark:text-dark-400 dark:hover:bg-dark-800"
            :title="t('home.viewDocs')"
          >
            <Icon name="book" size="md" />
          </a>
          <router-link
            v-if="showModelPlazaEntry"
            to="/model-plaza"
            class="flex h-10 shrink-0 items-center gap-1.5 rounded-lg px-2.5 text-sm font-medium text-gray-500 hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="t('nav.modelPlaza')"
          >
            <Icon name="grid" size="md" />
            <span class="hidden sm:inline">{{ t('nav.modelPlaza') }}</span>
          </router-link>
          <button
            class="flex h-10 w-10 shrink-0 items-center justify-center rounded-lg text-gray-500 hover:bg-gray-100 dark:text-dark-400 dark:hover:bg-dark-800"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
            @click="toggleTheme"
          >
            <Icon v-if="isDark" name="sun" size="md" />
            <Icon v-else name="moon" size="md" />
          </button>
          <router-link
            :to="isAuthenticated ? dashboardPath : '/login'"
            class="inline-flex min-h-10 shrink-0 items-center justify-center rounded-lg bg-gray-900 px-4 py-2 text-sm font-medium text-white hover:bg-gray-800 dark:bg-white dark:text-gray-900 dark:hover:bg-gray-200"
          >
            {{ isAuthenticated ? t('home.dashboard') : t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <main class="flex min-w-0 flex-1 items-center justify-center px-4 py-16 sm:px-6">
      <div class="min-w-0 max-w-2xl text-center">
        <img
          :src="siteLogo || '/logo.svg'"
          alt="Logo"
          class="mx-auto mb-6 h-20 w-20 rounded-2xl object-contain"
        />
        <h1 class="[overflow-wrap:anywhere] text-3xl font-bold md:text-4xl">{{ siteName }}</h1>
        <p class="mt-4 whitespace-pre-wrap [overflow-wrap:anywhere] text-base text-gray-600 dark:text-dark-300">{{ siteSubtitle }}</p>
        <router-link
          :to="isAuthenticated ? dashboardPath : '/login'"
          class="mt-8 inline-flex min-h-10 items-center justify-center rounded-lg bg-primary-600 px-5 py-2.5 text-sm font-medium text-white hover:bg-primary-700"
        >
          {{ isAuthenticated ? t('home.goToDashboard') : t('home.login') }}
        </router-link>
      </div>
    </main>

    <footer class="min-w-0 border-t border-gray-200 px-4 py-5 text-center text-sm text-gray-500 [overflow-wrap:anywhere] sm:px-6 dark:border-dark-800 dark:text-dark-400">
      &copy; {{ currentYear }} {{ siteName }}
    </footer>
  </div>

  <!-- Default Home Page - Aphelion Theme -->
  <div
    v-else
    class="aphelion relative flex min-h-screen flex-col overflow-hidden bg-gradient-to-b from-[#f7f4ed] via-[#f4f1e9] to-[#efe9dd] dark:from-[#1b1c19] dark:via-[#171815] dark:to-[#121310]"
  >
    <!-- Deep Space Backdrop -->
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      <!-- Starfield -->
      <div class="starfield starfield--far"></div>
      <div class="starfield starfield--near"></div>
      <!-- Distant sun: small, cold, far away -->
      <div class="distant-sun"></div>
      <div
        class="absolute -bottom-52 -right-40 h-[32rem] w-[32rem] rounded-full bg-[#d66b4d]/10 blur-3xl dark:bg-[#d66b4d]/15"
      ></div>
      <div
        class="absolute left-1/3 top-1/3 h-72 w-72 rounded-full bg-[#c9a227]/10 blur-3xl dark:bg-[#8a6a3b]/15"
      ></div>
      <!-- Faint orbital arcs sweeping across the sky -->
      <div class="sky-orbit sky-orbit--one"></div>
      <div class="sky-orbit sky-orbit--two"></div>
      <div
        class="absolute inset-0 bg-[linear-gradient(rgba(140,132,110,0.06)_1px,transparent_1px),linear-gradient(90deg,rgba(140,132,110,0.06)_1px,transparent_1px)] bg-[size:72px_72px]"
      ></div>
    </div>

    <!-- Header -->
    <header class="relative z-20 px-6 py-4">
      <nav class="mx-auto flex max-w-6xl items-center justify-between">
        <!-- Logo -->
        <div class="flex items-center">
          <div class="h-12 w-12 overflow-hidden">
            <img :src="siteLogo || '/logo.svg'" alt="Logo" class="h-full w-full object-contain" />
          </div>
        </div>

        <!-- Nav Actions -->
        <div class="flex items-center gap-3">
          <!-- Language Switcher -->
          <LocaleSwitcher />

          <!-- Doc Link -->
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="t('home.viewDocs')"
          >
            <Icon name="book" size="md" />
          </a>

          <!-- Model Plaza Link -->
          <router-link
            v-if="showModelPlazaEntry"
            to="/model-plaza"
            class="inline-flex items-center gap-1.5 rounded-lg p-2 text-sm text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="t('nav.modelPlaza')"
          >
            <Icon name="grid" size="md" />
            <span class="hidden sm:inline">{{ t('nav.modelPlaza') }}</span>
          </router-link>

          <!-- Theme Toggle -->
          <button
            @click="toggleTheme"
            class="rounded-lg p-2 text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="isDark ? t('home.switchToLight') : t('home.switchToDark')"
          >
            <Icon v-if="isDark" name="sun" size="md" />
            <Icon v-else name="moon" size="md" />
          </button>

          <!-- Login / Dashboard Button -->
          <router-link
            v-if="isAuthenticated"
            :to="dashboardPath"
            class="inline-flex items-center gap-1.5 rounded-full bg-gray-900 py-1 pl-1 pr-2.5 transition-colors hover:bg-gray-800 dark:bg-gray-800 dark:hover:bg-gray-700"
          >
            <span
              class="flex h-5 w-5 items-center justify-center rounded-full bg-gradient-to-br from-primary-400 to-primary-600 text-[10px] font-semibold text-white"
            >
              {{ userInitial }}
            </span>
            <span class="text-xs font-medium text-white">{{ t('home.dashboard') }}</span>
            <svg
              class="h-3 w-3 text-gray-400"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                d="M4.5 19.5l15-15m0 0H8.25m11.25 0v11.25"
              />
            </svg>
          </router-link>
          <router-link
            v-else
            to="/login"
            class="inline-flex items-center rounded-full bg-gray-900 px-3 py-1 text-xs font-medium text-white transition-colors hover:bg-gray-800 dark:bg-gray-800 dark:hover:bg-gray-700"
          >
            {{ t('home.login') }}
          </router-link>
        </div>
      </nav>
    </header>

    <!-- Main Content - centered, deliberately minimal -->
    <main class="relative z-10 flex flex-1 items-center justify-center px-6 py-12">
      <div class="mx-auto flex w-full max-w-3xl flex-col items-center text-center">
        <!-- Aphelion badge -->
        <div
          class="mb-5 inline-flex items-center gap-2 rounded-full border border-[#d66b4d]/40 bg-white/70 px-3.5 py-1.5 text-xs font-medium tracking-wide text-[#b4523a] shadow-sm backdrop-blur-sm dark:border-[#d66b4d]/30 dark:bg-white/5 dark:text-[#e79878]"
        >
          <span class="aphelion-pip"></span>
          {{ t('home.aphelion.badge') }}
        </div>

        <h1
          class="mb-4 [overflow-wrap:anywhere] text-4xl font-bold text-gray-900 md:text-5xl lg:text-6xl dark:text-[#f3f1ea]"
        >
          {{ siteName }}
        </h1>
        <p class="mb-3 text-lg text-gray-600 dark:text-[#c9c4b6] md:text-xl">
          {{ siteSubtitle }}
        </p>
        <p class="mb-8 max-w-xl text-sm leading-relaxed text-gray-500 dark:text-[#8f897b]">
          {{ t('home.aphelion.tagline') }}
        </p>

        <!-- Orbit Scene: the terminal sits at the far point of the orbit -->
        <div class="orbit-stage">
          <svg class="orbit-svg" viewBox="0 0 560 420" aria-hidden="true" focusable="false">
            <defs>
              <radialGradient id="aphelionSunGlow">
                <stop offset="0%" stop-color="#f3f1ea" stop-opacity="0.7" />
                <stop offset="35%" stop-color="#e2b08a" stop-opacity="0.22" />
                <stop offset="100%" stop-color="#d66b4d" stop-opacity="0" />
              </radialGradient>
            </defs>
            <g transform="rotate(-22 280 210)">
              <ellipse class="orbit-path" cx="280" cy="210" rx="285" ry="150" />
              <circle cx="38" cy="210" r="40" fill="url(#aphelionSunGlow)" />
              <circle class="orbit-sun" cx="38" cy="210" r="8" />
              <circle class="orbit-apex" cx="565" cy="210" r="5" />
              <circle class="orbit-planet" r="7">
                <animateMotion
                  dur="24s"
                  repeatCount="indefinite"
                  path="M 565 210 A 285 150 0 1 1 -5 210 A 285 150 0 1 1 565 210"
                  calcMode="spline"
                  keyPoints="0;1"
                  keyTimes="0;1"
                  keySplines="0.45 0 0.55 1"
                />
              </circle>
            </g>
          </svg>

          <!-- Card + caption stack, sized by the stage so the orbit never reflows it -->
          <div class="orbit-inner">
            <div class="terminal-container">
              <div class="terminal-window">
                <!-- Window header -->
                <div class="terminal-header">
                  <div class="terminal-buttons">
                    <span class="btn-close"></span>
                    <span class="btn-minimize"></span>
                    <span class="btn-maximize"></span>
                  </div>
                  <span class="terminal-title">aphelion ~ terminal</span>
                </div>
                <!-- Terminal content -->
                <div class="terminal-body">
                  <div class="code-line line-1">
                    <span class="code-prompt">$</span>
                    <span class="code-cmd">curl</span>
                    <span class="code-flag">-X POST</span>
                    <span class="code-url">/v1/messages</span>
                  </div>
                  <div class="code-line line-2">
                    <span class="code-comment"># Routing to upstream...</span>
                  </div>
                  <div class="code-line line-3">
                    <span class="code-success">200 OK</span>
                    <span class="code-response">{ "content": "Hello!" }</span>
                  </div>
                  <div class="code-line line-4">
                    <span class="code-prompt">$</span>
                    <span class="cursor"></span>
                  </div>
                </div>
              </div>
            </div>
            <p class="orbit-caption">{{ t('home.aphelion.orbitCaption') }}</p>
          </div>
        </div>

        <!-- Actions -->
        <div class="mt-8 flex flex-wrap items-center justify-center gap-4">
          <router-link
            :to="isAuthenticated ? dashboardPath : '/login'"
            class="inline-flex min-h-11 items-center justify-center rounded-lg bg-[#d66b4d] px-8 py-3 text-base font-medium text-[#fdf8f2] shadow-lg shadow-[#d66b4d]/30 transition-colors hover:bg-[#c25f43]"
          >
            {{ isAuthenticated ? t('home.goToDashboard') : t('home.getStarted') }}
            <Icon name="arrowRight" size="md" class="ml-2" :stroke-width="2" />
          </router-link>

          <!-- Contact support: a link when the setting is a URL, otherwise reveal the details -->
          <a
            v-if="contactUrl"
            :href="contactUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="inline-flex min-h-11 items-center justify-center rounded-lg border border-gray-300/70 bg-white/70 px-8 py-3 text-base font-medium text-gray-700 backdrop-blur-sm transition-colors hover:bg-white dark:border-[#3a382f] dark:bg-white/5 dark:text-[#d9d4c6] dark:hover:bg-white/10"
          >
            {{ t('common.contactSupport') }}
          </a>
          <button
            v-else-if="contactInfo"
            type="button"
            class="inline-flex min-h-11 items-center justify-center rounded-lg border border-gray-300/70 bg-white/70 px-8 py-3 text-base font-medium text-gray-700 backdrop-blur-sm transition-colors hover:bg-white dark:border-[#3a382f] dark:bg-white/5 dark:text-[#d9d4c6] dark:hover:bg-white/10"
            @click="contactRevealed = !contactRevealed"
          >
            {{ t('common.contactSupport') }}
          </button>
        </div>

        <p
          v-if="contactRevealed && contactInfo && !contactUrl"
          class="mt-4 [overflow-wrap:anywhere] text-sm text-gray-600 dark:text-[#c9c4b6]"
        >
          {{ contactInfo }}
        </p>
      </div>
    </main>

    <!-- Footer -->
    <footer class="relative z-10 border-t border-gray-200/50 px-6 py-8 dark:border-dark-800/50">
      <div
        class="mx-auto flex max-w-6xl flex-col items-center justify-center gap-4 text-center sm:flex-row sm:text-left"
      >
        <p class="text-sm text-gray-500 dark:text-dark-400">
          &copy; {{ currentYear }} {{ siteName }}. {{ t('home.footer.allRightsReserved') }}
        </p>
        <div class="flex items-center gap-4">
          <a
            v-if="docUrl"
            :href="docUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="text-sm text-gray-500 transition-colors hover:text-gray-700 dark:text-dark-400 dark:hover:text-white"
          >
            {{ t('home.docs') }}
          </a>
          <a
            :href="githubUrl"
            target="_blank"
            rel="noopener noreferrer"
            class="text-sm text-gray-500 transition-colors hover:text-gray-700 dark:text-dark-400 dark:hover:text-white"
          >
            GitHub
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
import { FeatureFlags, isFeatureFlagEnabled } from '@/utils/featureFlags'

const { t } = useI18n()

const authStore = useAuthStore()
const appStore = useAppStore()

// Site settings - directly from appStore (already initialized from injected config)
const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Aphelion')
const siteLogo = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const siteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle || 'AI API Gateway Platform')
const docUrl = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.doc_url || appStore.docUrl || ''))
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const hasHomeContent = computed(() => homeContent.value.trim().length > 0)
const compactHomeEnabled = computed(() => appStore.cachedPublicSettings?.compact_home_enabled === true)

// Customer support contact, set by admins in site settings. It is free text
// (e.g. "QQ: 123456789"), so only treat it as a link when it is a URL.
const contactInfo = computed(() =>
  (appStore.cachedPublicSettings?.contact_info || appStore.contactInfo || '').trim(),
)
const contactUrl = computed(() =>
  /^https?:\/\//i.test(contactInfo.value) ? sanitizeUrl(contactInfo.value) : '',
)
const contactRevealed = ref(false)
const modelPlazaEnabled = computed(() => isFeatureFlagEnabled(FeatureFlags.modelPlaza))

// Check if homeContent is a URL (for iframe display)
const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

// Theme
const isDark = ref(document.documentElement.classList.contains('dark'))

// GitHub URL
const githubUrl = 'https://github.com/Wei-Shaw/sub2api'

// Auth state
const isAuthenticated = computed(() => authStore.isAuthenticated)
const modelPlazaRequiresAuth = computed(
  () => appStore.cachedPublicSettings?.model_plaza_require_auth === true,
)
const showModelPlazaEntry = computed(
  () => modelPlazaEnabled.value && (isAuthenticated.value || !modelPlazaRequiresAuth.value),
)
const isAdmin = computed(() => authStore.isAdmin)
const dashboardPath = computed(() => isAdmin.value ? '/admin/dashboard' : '/dashboard')
const userInitial = computed(() => {
  const user = authStore.user
  if (!user || !user.email) return ''
  return user.email.charAt(0).toUpperCase()
})

// Current year for footer
const currentYear = computed(() => new Date().getFullYear())

// Toggle theme
function toggleTheme() {
  isDark.value = !isDark.value
  document.documentElement.classList.toggle('dark', isDark.value)
  localStorage.setItem('theme', isDark.value ? 'dark' : 'light')
}

// Initialize theme
function initTheme() {
  const savedTheme = localStorage.getItem('theme')
  if (
    savedTheme === 'dark' ||
    (!savedTheme && window.matchMedia('(prefers-color-scheme: dark)').matches)
  ) {
    isDark.value = true
    document.documentElement.classList.add('dark')
  }
}

onMounted(() => {
  initTheme()

  // Check auth state
  authStore.checkAuth()

  // Ensure public settings are loaded (will use cache if already loaded from injected config)
  if (!appStore.publicSettingsLoaded) {
    appStore.fetchPublicSettings()
  }
})
</script>

<style scoped>
/* ==========================================================================
   Aphelion theme: the far point of the orbit — deep space, a distant sun,
   and a station that keeps running out here in the dark.
   ========================================================================== */

/* Starfield */
.starfield {
  position: absolute;
  inset: -20%;
  background-repeat: repeat;
  opacity: 0.35;
}

.starfield--far {
  background-image:
    radial-gradient(1px 1px at 20% 30%, rgba(243, 241, 234, 0.45) 50%, transparent 50%),
    radial-gradient(1px 1px at 70% 15%, rgba(214, 107, 77, 0.4) 50%, transparent 50%),
    radial-gradient(1px 1px at 45% 70%, rgba(243, 241, 234, 0.35) 50%, transparent 50%),
    radial-gradient(1px 1px at 85% 60%, rgba(201, 162, 39, 0.3) 50%, transparent 50%);
  background-size: 260px 260px;
  animation: drift 180s linear infinite;
}

.starfield--near {
  background-image:
    radial-gradient(1.6px 1.6px at 12% 55%, rgba(243, 241, 234, 0.5) 50%, transparent 50%),
    radial-gradient(1.6px 1.6px at 62% 40%, rgba(226, 216, 196, 0.45) 50%, transparent 50%),
    radial-gradient(1.4px 1.4px at 88% 85%, rgba(214, 107, 77, 0.45) 50%, transparent 50%);
  background-size: 420px 420px;
  animation: drift 110s linear infinite reverse;
}

:global(.dark) .starfield {
  opacity: 0.9;
}

@keyframes drift {
  from {
    transform: translate3d(0, 0, 0);
  }
  to {
    transform: translate3d(-260px, 130px, 0);
  }
}

/* The sun, seen from aphelion: small and cold */
.distant-sun {
  position: absolute;
  top: -8rem;
  left: -6rem;
  width: 26rem;
  height: 26rem;
  border-radius: 9999px;
  background: radial-gradient(
    circle at 55% 55%,
    rgba(243, 241, 234, 0.34) 0%,
    rgba(214, 107, 77, 0.2) 28%,
    rgba(214, 107, 77, 0.07) 55%,
    transparent 72%
  );
  filter: blur(6px);
  animation: sun-breathe 12s ease-in-out infinite;
}

:global(.dark) .distant-sun {
  background: radial-gradient(
    circle at 55% 55%,
    rgba(243, 241, 234, 0.3) 0%,
    rgba(214, 107, 77, 0.22) 26%,
    rgba(120, 52, 33, 0.15) 52%,
    transparent 72%
  );
}

@keyframes sun-breathe {
  0%,
  100% {
    opacity: 0.85;
    transform: scale(1);
  }
  50% {
    opacity: 1;
    transform: scale(1.04);
  }
}

/* Faint orbital arcs across the whole sky */
.sky-orbit {
  position: absolute;
  border-radius: 9999px;
  border: 1px solid rgba(140, 132, 110, 0.16);
}

:global(.dark) .sky-orbit {
  border-color: rgba(243, 241, 234, 0.1);
}

.sky-orbit--one {
  top: -40%;
  left: -25%;
  width: 150%;
  height: 130%;
  transform: rotate(-12deg);
}

.sky-orbit--two {
  top: -10%;
  left: -5%;
  width: 120%;
  height: 160%;
  transform: rotate(8deg);
}

/* Badge pip - a tiny body glinting in the dark */
.aphelion-pip {
  width: 7px;
  height: 7px;
  border-radius: 9999px;
  background: radial-gradient(circle at 35% 35%, #f3f1ea, #d66b4d);
  box-shadow: 0 0 8px rgba(214, 107, 77, 0.75);
  animation: pip-pulse 3s ease-in-out infinite;
}

@keyframes pip-pulse {
  0%,
  100% {
    opacity: 0.7;
    transform: scale(0.9);
  }
  50% {
    opacity: 1;
    transform: scale(1.15);
  }
}

/* Orbit stage: the terminal sits at the far end of the ellipse */
.orbit-stage {
  position: relative;
  width: min(100%, 560px);
  aspect-ratio: 560 / 420;
}

/* Card + caption stack, sized by the stage so the orbit never reflows it */
.orbit-inner {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 16px;
}

.orbit-svg {
  position: absolute;
  inset: 0;
  width: 100%;
  height: 100%;
  overflow: visible;
}

.orbit-path {
  fill: none;
  stroke: rgba(120, 112, 94, 0.45);
  stroke-width: 1.3;
  stroke-dasharray: 5 7;
}

:global(.dark) .orbit-path {
  stroke: rgba(243, 241, 234, 0.35);
}

.orbit-sun {
  fill: #c96a4c;
  filter: drop-shadow(0 0 10px rgba(214, 107, 77, 0.5));
}

:global(.dark) .orbit-sun {
  fill: #f0ece0;
  filter: drop-shadow(0 0 10px rgba(243, 241, 234, 0.85));
}

.orbit-planet {
  fill: #d66b4d;
  filter: drop-shadow(0 0 9px rgba(214, 107, 77, 0.85));
}

.orbit-apex {
  fill: none;
  stroke: rgba(214, 107, 77, 0.75);
  stroke-width: 1.5;
}

.orbit-caption {
  position: relative;
  z-index: 2;
  margin: 0;
  white-space: nowrap;
  font-family: ui-monospace, 'Fira Code', monospace;
  font-size: 11px;
  letter-spacing: 0.08em;
  color: rgba(122, 114, 98, 0.95);
}

:global(.dark) .orbit-caption {
  color: rgba(214, 107, 77, 0.85);
}

@media (max-width: 480px) {
  .orbit-caption {
    font-size: 10px;
  }
}

@media (prefers-reduced-motion: reduce) {
  .starfield,
  .distant-sun,
  .aphelion-pip {
    animation: none;
  }
}

/* Terminal Container */
.terminal-container {
  position: relative;
  z-index: 2;
  display: flex;
  width: min(100%, 380px);
  justify-content: center;
}

/* Terminal Window */
.terminal-window {
  width: min(100%, 380px);
  background: linear-gradient(145deg, #26261f 0%, #131410 100%);
  border-radius: 14px;
  box-shadow:
    0 25px 50px -12px rgba(10, 10, 6, 0.5),
    0 0 0 1px rgba(243, 241, 234, 0.14),
    0 0 45px rgba(214, 107, 77, 0.14),
    inset 0 1px 0 rgba(255, 255, 255, 0.08);
  overflow: hidden;
  transform: perspective(1000px) rotateX(2deg) rotateY(-2deg);
  transition: transform 0.3s ease;
}

.terminal-window:hover {
  transform: perspective(1000px) rotateX(0deg) rotateY(0deg) translateY(-4px);
}

/* Terminal Header */
.terminal-header {
  display: flex;
  align-items: center;
  padding: 12px 16px;
  background: rgba(40, 39, 32, 0.85);
  border-bottom: 1px solid rgba(243, 241, 234, 0.1);
}

.terminal-buttons {
  display: flex;
  gap: 8px;
}

.terminal-buttons span {
  width: 12px;
  height: 12px;
  border-radius: 50%;
}

.btn-close {
  background: #ef4444;
}
.btn-minimize {
  background: #eab308;
}
.btn-maximize {
  background: #22c55e;
}

.terminal-title {
  flex: 1;
  text-align: center;
  font-size: 12px;
  font-family: ui-monospace, monospace;
  color: #8f897b;
  margin-right: 52px;
}

/* Terminal Body */
.terminal-body {
  padding: 20px 24px;
  font-family: ui-monospace, 'Fira Code', monospace;
  font-size: 14px;
  line-height: 2;
}

.code-line {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
  opacity: 0;
  animation: line-appear 0.5s ease forwards;
}

.line-1 {
  animation-delay: 0.3s;
}
.line-2 {
  animation-delay: 1s;
}
.line-3 {
  animation-delay: 1.8s;
}
.line-4 {
  animation-delay: 2.5s;
}

@keyframes line-appear {
  from {
    opacity: 0;
    transform: translateY(5px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.code-prompt {
  color: #22c55e;
  font-weight: bold;
}
.code-cmd {
  color: #38bdf8;
}
.code-flag {
  color: #a78bfa;
}
.code-url {
  color: #e2b08a;
}
.code-comment {
  color: #837d6d;
  font-style: italic;
}
.code-success {
  color: #22c55e;
  background: rgba(34, 197, 94, 0.15);
  padding: 2px 8px;
  border-radius: 4px;
  font-weight: 600;
}
.code-response {
  color: #e6c07b;
}

/* Blinking Cursor */
.cursor {
  display: inline-block;
  width: 8px;
  height: 16px;
  background: #22c55e;
  animation: blink 1s step-end infinite;
}

@keyframes blink {
  0%,
  50% {
    opacity: 1;
  }
  51%,
  100% {
    opacity: 0;
  }
}

/* Dark mode adjustments */
:global(.dark) .terminal-window {
  box-shadow:
    0 25px 50px -12px rgba(0, 0, 0, 0.65),
    0 0 0 1px rgba(243, 241, 234, 0.16),
    0 0 55px rgba(214, 107, 77, 0.18),
    inset 0 1px 0 rgba(255, 255, 255, 0.08);
}
</style>
