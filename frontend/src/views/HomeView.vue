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

          <!-- 服务状态入口：渠道监控开启时对访客公开 -->
          <router-link
            v-if="channelMonitorEnabled"
            to="/monitor"
            class="inline-flex items-center gap-1.5 rounded-lg p-2 text-sm text-gray-500 transition-colors hover:bg-gray-100 hover:text-gray-700 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white"
            :title="t('home.trust.statusNav')"
          >
            <span class="h-2 w-2 rounded-full bg-emerald-500"></span>
            <span class="hidden sm:inline">{{ t('home.trust.statusNav') }}</span>
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
                  <!-- 1. 示例终端沿用站点配置名称，避免出现旧品牌。 -->
                  <span class="terminal-title">{{ siteName }} ~ terminal</span>
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

          <!-- 联系客服：点击弹窗展示联系方式；未配置 contact_info 时整块不渲染 -->
          <button
            v-if="contactInfo"
            type="button"
            class="inline-flex min-h-11 items-center justify-center gap-2 rounded-lg border border-gray-300/70 bg-white/70 px-8 py-3 text-base font-medium text-gray-700 backdrop-blur-sm transition-colors hover:bg-white dark:border-[#3a382f] dark:bg-white/5 dark:text-[#d9d4c6] dark:hover:bg-white/10"
            @click="contactDialogOpen = true"
          >
            <Icon name="chatBubble" size="md" :stroke-width="2" />
            {{ t('common.contactSupport') }}
          </button>
        </div>
      </div>
    </main>

    <!-- 信任区块：累计数据 / 模型价格 / 痛点 / 最近更新 -->
    <section class="relative z-10 px-6 pb-16">
      <div class="mx-auto flex w-full max-w-5xl flex-col gap-12">
        <!-- 1. 累计服务数据（接口失败时整块不渲染） -->
        <div v-if="publicStats" class="grid grid-cols-1 gap-4 sm:grid-cols-3">
          <div class="trust-card text-center">
            <p class="text-3xl font-bold text-gray-900 dark:text-[#f3f1ea]">{{ formatCompact(publicStats.total_requests) }}</p>
            <p class="mt-1 text-sm text-gray-500 dark:text-[#8f897b]">{{ t('home.trust.totalRequests') }}</p>
          </div>
          <div class="trust-card text-center">
            <p class="text-3xl font-bold text-gray-900 dark:text-[#f3f1ea]">{{ formatCompact(publicStats.total_tokens) }}</p>
            <p class="mt-1 text-sm text-gray-500 dark:text-[#8f897b]">{{ t('home.trust.totalTokens') }}</p>
          </div>
          <router-link v-if="channelMonitorEnabled" to="/monitor" class="trust-card text-center transition-colors hover:border-[#d66b4d]/50">
            <!-- 真实数据：各监控 7 天可用率均值；全部正常为绿点，否则为黄点 -->
            <p class="flex items-center justify-center gap-2 text-3xl font-bold text-gray-900 dark:text-[#f3f1ea]">
              <span class="h-3 w-3 rounded-full" :class="monitorSummary?.allOperational ? 'bg-emerald-500' : 'bg-amber-500'"></span>
              {{ monitorSummary ? `${monitorSummary.availability.toFixed(2)}%` : '-' }}
            </p>
            <p class="mt-1 text-sm text-gray-500 dark:text-[#8f897b]">{{ t('home.trust.statusDesc') }}</p>
          </router-link>
        </div>

        <!-- 2. 支持的模型与价格（取模型广场数据，最多展示 8 个） -->
        <div v-if="plazaModels.length > 0">
          <div class="mb-5 flex items-end justify-between gap-4">
            <div>
              <h2 class="text-2xl font-bold text-gray-900 dark:text-[#f3f1ea]">{{ t('home.trust.modelsTitle') }}</h2>
              <p class="mt-1 text-sm text-gray-500 dark:text-[#8f897b]">{{ t('home.trust.modelsDesc') }}</p>
            </div>
            <router-link to="/model-plaza" class="shrink-0 text-sm font-medium text-[#d66b4d] hover:underline">
              {{ t('home.trust.viewAllModels') }}
            </router-link>
          </div>
          <div class="grid grid-cols-1 gap-3 sm:grid-cols-2">
            <div v-for="m in plazaModels" :key="m.key" class="trust-card flex items-center justify-between gap-3">
              <div class="min-w-0">
                <p class="truncate font-mono text-sm font-semibold text-gray-900 dark:text-[#f3f1ea]">{{ m.name }}</p>
                <p class="text-xs text-gray-500 dark:text-[#8f897b]">{{ m.group }}</p>
              </div>
              <div class="shrink-0 text-right text-xs text-gray-600 dark:text-[#c9c4b6]">
                <p>{{ t('home.trust.input') }} ${{ m.input }} / {{ t('home.trust.output') }} ${{ m.output }}</p>
                <p class="text-gray-400 dark:text-[#8f897b]">{{ t('home.trust.perMillion') }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- 3. 用户痛点（复用已有文案） -->
        <div>
          <h2 class="mb-5 text-center text-2xl font-bold text-gray-900 dark:text-[#f3f1ea]">{{ t('home.painPoints.title') }}</h2>
          <div class="grid grid-cols-1 gap-4 sm:grid-cols-2 lg:grid-cols-4">
            <div v-for="key in painPointKeys" :key="key" class="trust-card">
              <p class="font-semibold text-gray-900 dark:text-[#f3f1ea]">{{ t(`home.painPoints.items.${key}.title`) }}</p>
              <p class="mt-2 text-sm leading-relaxed text-gray-500 dark:text-[#8f897b]">{{ t(`home.painPoints.items.${key}.desc`) }}</p>
            </div>
          </div>
          <div class="mt-6 flex flex-wrap items-center justify-center gap-2">
            <span v-for="tag in featureTagKeys" :key="tag" class="rounded-full border border-[#d66b4d]/40 px-3 py-1 text-xs font-medium text-[#b4523a] dark:text-[#e79878]">
              {{ t(`home.tags.${tag}`) }}
            </span>
          </div>
        </div>

        <!-- 4. 最近更新（后台公告，无公开公告时不渲染） -->
        <div v-if="announcements.length > 0">
          <h2 class="mb-5 text-2xl font-bold text-gray-900 dark:text-[#f3f1ea]">{{ t('home.trust.updatesTitle') }}</h2>
          <ul class="flex flex-col gap-3">
            <li v-for="a in announcements" :key="a.id" class="trust-card">
              <div class="flex items-baseline justify-between gap-3">
                <p class="font-semibold text-gray-900 dark:text-[#f3f1ea]">{{ a.title }}</p>
                <span class="shrink-0 text-xs text-gray-400 dark:text-[#8f897b]">{{ a.created_at.slice(0, 10) }}</span>
              </div>
              <p class="mt-1 line-clamp-2 whitespace-pre-line text-sm text-gray-500 dark:text-[#8f897b]">{{ a.content }}</p>
            </li>
          </ul>
        </div>
      </div>
    </section>

    <!-- 联系客服弹窗：沿用公告弹窗的暖色头部与卡片尺寸 -->
    <BaseDialog
      :show="contactDialogOpen"
      :title="t('common.contactSupport')"
      width="normal"
      close-on-click-outside
      @close="contactDialogOpen = false"
    >
      <div class="space-y-5">
        <div class="flex items-start gap-4">
          <div
            class="flex h-11 w-11 flex-none items-center justify-center rounded-xl bg-gradient-to-br from-amber-500 to-orange-600 text-white shadow-lg shadow-amber-500/30"
          >
            <Icon name="chatBubble" size="md" :stroke-width="2" />
          </div>
          <p class="pt-1 text-sm leading-relaxed text-gray-600 dark:text-dark-300">
            {{ t('home.aphelion.contactHint') }}
          </p>
        </div>

        <!-- 联系方式：链接可直接点，纯文本（微信号 / QQ 号）选中复制 -->
        <a
          v-if="contactUrl"
          :href="contactUrl"
          target="_blank"
          rel="noopener noreferrer"
          class="block rounded-xl border border-gray-200 bg-gray-50 px-4 py-3 text-sm font-medium text-amber-700 transition-colors hover:border-amber-400 hover:bg-amber-50 dark:border-dark-700 dark:bg-dark-900/40 dark:text-amber-400 dark:hover:border-amber-500/50 dark:hover:bg-amber-900/20"
        >
          <span class="mb-0.5 block text-xs font-normal text-gray-500 dark:text-dark-400">
            {{ t('home.aphelion.contactLinkLabel') }}
          </span>
          <span class="[overflow-wrap:anywhere]">{{ contactInfo }}</span>
        </a>
        <div
          v-else
          class="rounded-xl border border-gray-200 bg-gray-50 px-4 py-3 dark:border-dark-700 dark:bg-dark-900/40"
        >
          <span class="mb-0.5 block text-xs text-gray-500 dark:text-dark-400">
            {{ t('home.aphelion.contactTextLabel') }}
          </span>
          <span
            class="block select-all font-mono text-base font-semibold tracking-wide text-gray-900 [overflow-wrap:anywhere] dark:text-white"
          >
            {{ contactInfo }}
          </span>
        </div>
      </div>

      <template #footer>
        <div class="flex items-center justify-end gap-3">
          <button
            v-if="!contactUrl"
            type="button"
            class="btn btn-secondary"
            @click="copyContact()"
          >
            <Icon name="copy" size="sm" class="mr-1.5" />
            {{ copied ? t('common.copied') : t('common.copy') }}
          </button>
          <button type="button" class="btn btn-primary" @click="contactDialogOpen = false">
            {{ t('common.close') }}
          </button>
        </div>
      </template>
    </BaseDialog>

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
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore, useAppStore } from '@/stores'
import BaseDialog from '@/components/common/BaseDialog.vue'
import LocaleSwitcher from '@/components/common/LocaleSwitcher.vue'
import Icon from '@/components/icons/Icon.vue'
import { sanitizeUrl } from '@/utils/url'
import { useClipboard } from '@/composables/useClipboard'
import { FeatureFlags, isFeatureFlagEnabled, isChannelMonitorV1Mode } from '@/utils/featureFlags'
import { getModelPlaza } from '@/api/modelPlaza'
import { list as listChannelMonitors } from '@/api/channelMonitor'
import { getPublicAnnouncements, getPublicStats, type PublicAnnouncement, type PublicStats } from '@/api/publicInfo'

const { t } = useI18n()

const authStore = useAuthStore()
const appStore = useAppStore()

// 1. 优先使用已注入的站点配置，未配置时统一显示 mdai。
const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'mdai')
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

// 联繗客服弹窗开关
const contactDialogOpen = ref(false)

// 纯文本联系方式不支持点击跳转，改为提供一键复制
const { copied, copyToClipboard } = useClipboard()

function copyContact() {
  copyToClipboard(contactInfo.value)
}
const modelPlazaEnabled = computed(() => isFeatureFlagEnabled(FeatureFlags.modelPlaza))

// Check if homeContent is a URL (for iframe display)
const isHomeContentUrl = computed(() => {
  const content = homeContent.value.trim()
  return content.startsWith('http://') || content.startsWith('https://')
})

// Theme
const isDark = ref(document.documentElement.classList.contains('dark'))

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

// 服务状态入口：仅 V1 监控模式开启时展示（V2 为登录后被动视图）
const channelMonitorEnabled = computed(() => isChannelMonitorV1Mode())

// 首页信任区块数据：接口失败时对应区块不渲染，不影响首页主体
const publicStats = ref<PublicStats | null>(null)
const monitorSummary = ref<{ availability: number; allOperational: boolean } | null>(null)
const announcements = ref<PublicAnnouncement[]>([])
const plazaModels = ref<{ key: string; name: string; group: string; input: string; output: string }[]>([])
const painPointKeys = ['expensive', 'complex', 'unstable', 'noControl']
const featureTagKeys = ['subscriptionToApi', 'stickySession', 'realtimeBilling']

// 大数字压缩显示：1.2 万 / 3.4 亿（英文 1.2K / 3.4B）
function formatCompact(n: number) {
  return new Intl.NumberFormat(document.documentElement.lang || undefined, { notation: 'compact', maximumFractionDigits: 1 }).format(n)
}

// 单价（USD/token）换算为每百万 token 价格，乘分组倍率得到实付价
function perMillion(price: number | null | undefined, rate: number) {
  if (price == null) return '-'
  return (price * 1_000_000 * rate).toFixed(2)
}

async function loadTrustData() {
  // 1. 累计数据与公告：互不依赖，并行请求
  getPublicStats().then((v) => (publicStats.value = v)).catch(() => {})
  getPublicAnnouncements().then((v) => (announcements.value = v)).catch(() => {})
  // 2. 服务状态：取各监控 7 天可用率均值，无监控数据时卡片显示 "-"
  if (channelMonitorEnabled.value) {
    listChannelMonitors()
      .then((res) => {
        const items = res.items || []
        if (items.length === 0) return
        monitorSummary.value = {
          availability: items.reduce((sum, it) => sum + it.availability_7d, 0) / items.length,
          allOperational: items.every((it) => it.primary_status === 'operational'),
        }
      })
      .catch(() => {})
  }
  // 3. 模型价格：广场未公开（关闭或需登录）时跳过
  if (!showModelPlazaEntry.value) return
  try {
    const plaza = await getModelPlaza()
    const seen = new Set<string>()
    const rows: typeof plazaModels.value = []
    for (const g of plaza.groups) {
      for (const m of g.models) {
        // 只展示按 token 计费、且同名模型只出现一次
        if (seen.has(m.name) || !m.pricing || m.pricing.billing_mode !== 'token') continue
        seen.add(m.name)
        rows.push({
          key: `${g.id}-${m.name}`,
          name: m.name,
          group: g.name,
          input: perMillion(m.pricing.input_price, g.rate_multiplier),
          output: perMillion(m.pricing.output_price, g.rate_multiplier),
        })
      }
    }
    plazaModels.value = rows.slice(0, 8)
  } catch {
    plazaModels.value = []
  }
}

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
    appStore.fetchPublicSettings().then(() => loadTrustData())
  } else {
    void loadTrustData()
  }
})
</script>

<style scoped>
/* 信任区块卡片：与首页暖色主题一致 */
.trust-card {
  @apply rounded-xl border border-gray-200/70 bg-white/70 p-4 backdrop-blur-sm dark:border-[#3a382f] dark:bg-white/5;
}

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
