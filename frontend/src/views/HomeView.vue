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
    class="aphelion relative flex min-h-screen flex-col overflow-hidden bg-gradient-to-b from-[#eef1fa] via-[#f4f3fb] to-[#fdf5ec] dark:from-[#04060f] dark:via-[#080b1c] dark:to-[#100a18]"
  >
    <!-- Deep Space Backdrop -->
    <div class="pointer-events-none absolute inset-0 overflow-hidden">
      <!-- Starfield -->
      <div class="starfield starfield--far"></div>
      <div class="starfield starfield--near"></div>
      <!-- Distant sun: small, cold, far away -->
      <div class="distant-sun"></div>
      <div
        class="absolute -bottom-52 -right-40 h-[32rem] w-[32rem] rounded-full bg-indigo-500/10 blur-3xl dark:bg-indigo-500/20"
      ></div>
      <div
        class="absolute left-1/3 top-1/3 h-72 w-72 rounded-full bg-violet-400/10 blur-3xl dark:bg-violet-600/15"
      ></div>
      <!-- Faint orbital arcs sweeping across the sky -->
      <div class="sky-orbit sky-orbit--one"></div>
      <div class="sky-orbit sky-orbit--two"></div>
      <div
        class="absolute inset-0 bg-[linear-gradient(rgba(99,102,241,0.035)_1px,transparent_1px),linear-gradient(90deg,rgba(99,102,241,0.035)_1px,transparent_1px)] bg-[size:72px_72px]"
      ></div>
    </div>

    <!-- Header -->
    <header class="relative z-20 px-6 py-4">
      <nav class="mx-auto flex max-w-6xl items-center justify-between">
        <!-- Logo -->
        <div class="flex items-center">
          <div class="h-10 w-10 overflow-hidden rounded-xl shadow-md">
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

    <!-- Main Content -->
    <main class="relative z-10 flex-1 px-6 py-16">
      <div class="mx-auto max-w-6xl">
        <!-- Hero Section - Left/Right Layout -->
        <div class="mb-12 flex flex-col items-center justify-between gap-12 lg:flex-row lg:gap-16">
          <!-- Left: Text Content -->
          <div class="flex-1 text-center lg:text-left">
            <!-- Aphelion badge -->
            <div
              class="mb-5 inline-flex items-center gap-2 rounded-full border border-amber-300/50 bg-white/70 px-3.5 py-1.5 text-xs font-medium tracking-wide text-amber-700 shadow-sm backdrop-blur-sm dark:border-amber-400/25 dark:bg-white/5 dark:text-amber-200"
            >
              <span class="aphelion-pip"></span>
              {{ t('home.aphelion.badge') }}
            </div>
            <h1
              class="mb-4 text-4xl font-bold text-gray-900 dark:text-white md:text-5xl lg:text-6xl"
            >
              {{ siteName }}
            </h1>
            <p class="mb-4 text-lg text-gray-600 dark:text-dark-300 md:text-xl">
              {{ siteSubtitle }}
            </p>
            <p class="mb-8 text-sm leading-relaxed text-gray-500 dark:text-dark-400">
              {{ t('home.aphelion.tagline') }}
            </p>

            <!-- CTA Button -->
            <div>
              <router-link
                :to="isAuthenticated ? dashboardPath : '/login'"
                class="btn btn-primary px-8 py-3 text-base shadow-lg shadow-indigo-500/30"
              >
                {{ isAuthenticated ? t('home.goToDashboard') : t('home.getStarted') }}
                <Icon name="arrowRight" size="md" class="ml-2" :stroke-width="2" />
              </router-link>
            </div>
          </div>

          <!-- Right: Orbit Scene with Terminal at Aphelion -->
          <div class="flex flex-1 justify-center lg:justify-end">
            <div class="orbit-stage">
              <!-- Elliptical orbit: the sun sits at one focus, we sit at the far end -->
              <svg class="orbit-svg" viewBox="0 0 560 460" aria-hidden="true" focusable="false">
                <defs>
                  <radialGradient id="aphelionSunGlow">
                    <stop offset="0%" stop-color="#fde68a" stop-opacity="0.95" />
                    <stop offset="45%" stop-color="#fb923c" stop-opacity="0.35" />
                    <stop offset="100%" stop-color="#fb923c" stop-opacity="0" />
                  </radialGradient>
                </defs>
                <g transform="rotate(-14 280 230)">
                  <ellipse class="orbit-path" cx="280" cy="230" rx="250" ry="155" />
                  <circle cx="86" cy="230" r="58" fill="url(#aphelionSunGlow)" />
                  <circle class="orbit-sun" cx="86" cy="230" r="9" />
                  <circle class="orbit-apex" cx="530" cy="230" r="4" />
                  <circle class="orbit-planet" r="6">
                    <animateMotion
                      dur="24s"
                      repeatCount="indefinite"
                      path="M 530 230 A 250 155 0 1 1 30 230 A 250 155 0 1 1 530 230"
                      keyPoints="0;1"
                      keyTimes="0;1"
                      calcMode="spline"
                      keySplines="0.45 0 0.55 1"
                    />
                  </circle>
                </g>
              </svg>
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
                <p class="orbit-caption">{{ t('home.aphelion.orbitCaption') }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Feature Tags - Centered -->
        <div class="mb-12 flex flex-wrap items-center justify-center gap-4 md:gap-6">
          <div
            class="inline-flex items-center gap-2.5 rounded-full border border-gray-200/50 bg-white/80 px-5 py-2.5 shadow-sm backdrop-blur-sm dark:border-dark-700/50 dark:bg-dark-800/80"
          >
            <Icon name="swap" size="sm" class="text-amber-500" />
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{
              t('home.tags.subscriptionToApi')
            }}</span>
          </div>
          <div
            class="inline-flex items-center gap-2.5 rounded-full border border-gray-200/50 bg-white/80 px-5 py-2.5 shadow-sm backdrop-blur-sm dark:border-dark-700/50 dark:bg-dark-800/80"
          >
            <Icon name="shield" size="sm" class="text-indigo-500 dark:text-indigo-400" />
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{
              t('home.tags.stickySession')
            }}</span>
          </div>
          <div
            class="inline-flex items-center gap-2.5 rounded-full border border-gray-200/50 bg-white/80 px-5 py-2.5 shadow-sm backdrop-blur-sm dark:border-dark-700/50 dark:bg-dark-800/80"
          >
            <Icon name="chart" size="sm" class="text-violet-500 dark:text-violet-400" />
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{
              t('home.tags.realtimeBilling')
            }}</span>
          </div>
        </div>

        <!-- Features Grid -->
        <div class="mb-12 grid gap-6 md:grid-cols-3">
          <!-- Feature 1: Unified Gateway -->
          <div
            class="group rounded-2xl border border-gray-200/50 bg-white/60 p-6 backdrop-blur-sm transition-all duration-300 hover:shadow-xl hover:shadow-indigo-500/20 dark:border-dark-700/50 dark:bg-dark-800/60"
          >
            <div
              class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-amber-400 to-orange-500 shadow-lg shadow-amber-500/30 transition-transform group-hover:scale-110"
            >
              <Icon name="server" size="lg" class="text-white" />
            </div>
            <h3 class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('home.features.unifiedGateway') }}
            </h3>
            <p class="text-sm leading-relaxed text-gray-600 dark:text-dark-400">
              {{ t('home.features.unifiedGatewayDesc') }}
            </p>
          </div>

          <!-- Feature 2: Account Pool -->
          <div
            class="group rounded-2xl border border-gray-200/50 bg-white/60 p-6 backdrop-blur-sm transition-all duration-300 hover:shadow-xl hover:shadow-indigo-500/20 dark:border-dark-700/50 dark:bg-dark-800/60"
          >
            <div
              class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-indigo-500 to-indigo-600 shadow-lg shadow-indigo-500/30 transition-transform group-hover:scale-110"
            >
              <svg
                class="h-6 w-6 text-white"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                stroke-width="1.5"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M18 18.72a9.094 9.094 0 003.741-.479 3 3 0 00-4.682-2.72m.94 3.198l.001.031c0 .225-.012.447-.037.666A11.944 11.944 0 0112 21c-2.17 0-4.207-.576-5.963-1.584A6.062 6.062 0 016 18.719m12 0a5.971 5.971 0 00-.941-3.197m0 0A5.995 5.995 0 0012 12.75a5.995 5.995 0 00-5.058 2.772m0 0a3 3 0 00-4.681 2.72 8.986 8.986 0 003.74.477m.94-3.197a5.971 5.971 0 00-.94 3.197M15 6.75a3 3 0 11-6 0 3 3 0 016 0zm6 3a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0zm-13.5 0a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z"
                />
              </svg>
            </div>
            <h3 class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('home.features.multiAccount') }}
            </h3>
            <p class="text-sm leading-relaxed text-gray-600 dark:text-dark-400">
              {{ t('home.features.multiAccountDesc') }}
            </p>
          </div>

          <!-- Feature 3: Billing & Quota -->
          <div
            class="group rounded-2xl border border-gray-200/50 bg-white/60 p-6 backdrop-blur-sm transition-all duration-300 hover:shadow-xl hover:shadow-indigo-500/20 dark:border-dark-700/50 dark:bg-dark-800/60"
          >
            <div
              class="mb-4 flex h-12 w-12 items-center justify-center rounded-xl bg-gradient-to-br from-violet-500 to-fuchsia-600 shadow-lg shadow-violet-500/30 transition-transform group-hover:scale-110"
            >
              <svg
                class="h-6 w-6 text-white"
                fill="none"
                viewBox="0 0 24 24"
                stroke="currentColor"
                stroke-width="1.5"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M2.25 18.75a60.07 60.07 0 0115.797 2.101c.727.198 1.453-.342 1.453-1.096V18.75M3.75 4.5v.75A.75.75 0 013 6h-.75m0 0v-.375c0-.621.504-1.125 1.125-1.125H20.25M2.25 6v9m18-10.5v.75c0 .414.336.75.75.75h.75m-1.5-1.5h.375c.621 0 1.125.504 1.125 1.125v9.75c0 .621-.504 1.125-1.125 1.125h-.375m1.5-1.5H21a.75.75 0 00-.75.75v.75m0 0H3.75m0 0h-.375a1.125 1.125 0 01-1.125-1.125V15m1.5 1.5v-.75A.75.75 0 003 15h-.75M15 10.5a3 3 0 11-6 0 3 3 0 016 0zm3 0h.008v.008H18V10.5zm-12 0h.008v.008H6V10.5z"
                />
              </svg>
            </div>
            <h3 class="mb-2 text-lg font-semibold text-gray-900 dark:text-white">
              {{ t('home.features.balanceQuota') }}
            </h3>
            <p class="text-sm leading-relaxed text-gray-600 dark:text-dark-400">
              {{ t('home.features.balanceQuotaDesc') }}
            </p>
          </div>
        </div>

        <!-- Supported Providers -->
        <div class="mb-8 text-center">
          <h2 class="mb-3 text-2xl font-bold text-gray-900 dark:text-white">
            {{ t('home.providers.title') }}
          </h2>
          <p class="text-sm text-gray-600 dark:text-dark-400">
            {{ t('home.providers.description') }}
          </p>
        </div>

        <div class="mb-16 flex flex-wrap items-center justify-center gap-4">
          <!-- Claude - Supported -->
          <div
            class="flex items-center gap-2 rounded-xl border border-primary-200 bg-white/60 px-5 py-3 ring-1 ring-primary-500/20 backdrop-blur-sm dark:border-primary-800 dark:bg-dark-800/60"
          >
            <div
              class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-orange-400 to-orange-500"
            >
              <span class="text-xs font-bold text-white">C</span>
            </div>
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{ t('home.providers.claude') }}</span>
            <span
              class="rounded bg-primary-100 px-1.5 py-0.5 text-[10px] font-medium text-primary-600 dark:bg-primary-900/30 dark:text-primary-400"
              >{{ t('home.providers.supported') }}</span
            >
          </div>
          <!-- GPT - Supported -->
          <div
            class="flex items-center gap-2 rounded-xl border border-primary-200 bg-white/60 px-5 py-3 ring-1 ring-primary-500/20 backdrop-blur-sm dark:border-primary-800 dark:bg-dark-800/60"
          >
            <div
              class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-green-500 to-green-600"
            >
              <span class="text-xs font-bold text-white">G</span>
            </div>
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">GPT</span>
            <span
              class="rounded bg-primary-100 px-1.5 py-0.5 text-[10px] font-medium text-primary-600 dark:bg-primary-900/30 dark:text-primary-400"
              >{{ t('home.providers.supported') }}</span
            >
          </div>
          <!-- Gemini - Supported -->
          <div
            class="flex items-center gap-2 rounded-xl border border-primary-200 bg-white/60 px-5 py-3 ring-1 ring-primary-500/20 backdrop-blur-sm dark:border-primary-800 dark:bg-dark-800/60"
          >
            <div
              class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-blue-500 to-blue-600"
            >
              <span class="text-xs font-bold text-white">G</span>
            </div>
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{ t('home.providers.gemini') }}</span>
            <span
              class="rounded bg-primary-100 px-1.5 py-0.5 text-[10px] font-medium text-primary-600 dark:bg-primary-900/30 dark:text-primary-400"
              >{{ t('home.providers.supported') }}</span
            >
          </div>
          <!-- Antigravity - Supported -->
          <div
            class="flex items-center gap-2 rounded-xl border border-primary-200 bg-white/60 px-5 py-3 ring-1 ring-primary-500/20 backdrop-blur-sm dark:border-primary-800 dark:bg-dark-800/60"
          >
            <div
              class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-rose-500 to-pink-600"
            >
              <span class="text-xs font-bold text-white">A</span>
            </div>
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{ t('home.providers.antigravity') }}</span>
            <span
              class="rounded bg-primary-100 px-1.5 py-0.5 text-[10px] font-medium text-primary-600 dark:bg-primary-900/30 dark:text-primary-400"
              >{{ t('home.providers.supported') }}</span
            >
          </div>
          <!-- More - Coming Soon -->
          <div
            class="flex items-center gap-2 rounded-xl border border-gray-200/50 bg-white/40 px-5 py-3 opacity-60 backdrop-blur-sm dark:border-dark-700/50 dark:bg-dark-800/40"
          >
            <div
              class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-gray-500 to-gray-600"
            >
              <span class="text-xs font-bold text-white">+</span>
            </div>
            <span class="text-sm font-medium text-gray-700 dark:text-dark-200">{{ t('home.providers.more') }}</span>
            <span
              class="rounded bg-gray-100 px-1.5 py-0.5 text-[10px] font-medium text-gray-500 dark:bg-dark-700 dark:text-dark-400"
              >{{ t('home.providers.soon') }}</span
            >
          </div>
        </div>
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
const siteName = computed(() => appStore.cachedPublicSettings?.site_name || appStore.siteName || 'Sub2API')
const siteLogo = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.site_logo || appStore.siteLogo || '', { allowRelative: true, allowDataUrl: true }))
const siteSubtitle = computed(() => appStore.cachedPublicSettings?.site_subtitle || 'AI API Gateway Platform')
const docUrl = computed(() => sanitizeUrl(appStore.cachedPublicSettings?.doc_url || appStore.docUrl || ''))
const homeContent = computed(() => appStore.cachedPublicSettings?.home_content || '')
const hasHomeContent = computed(() => homeContent.value.trim().length > 0)
const compactHomeEnabled = computed(() => appStore.cachedPublicSettings?.compact_home_enabled === true)
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
    radial-gradient(1px 1px at 20% 30%, rgba(99, 102, 241, 0.55) 50%, transparent 50%),
    radial-gradient(1px 1px at 70% 15%, rgba(148, 163, 184, 0.5) 50%, transparent 50%),
    radial-gradient(1px 1px at 45% 70%, rgba(129, 140, 248, 0.45) 50%, transparent 50%),
    radial-gradient(1px 1px at 85% 60%, rgba(148, 163, 184, 0.4) 50%, transparent 50%);
  background-size: 260px 260px;
  animation: drift 180s linear infinite;
}

.starfield--near {
  background-image:
    radial-gradient(1.6px 1.6px at 12% 55%, rgba(253, 230, 138, 0.5) 50%, transparent 50%),
    radial-gradient(1.6px 1.6px at 62% 40%, rgba(226, 232, 240, 0.5) 50%, transparent 50%),
    radial-gradient(1.4px 1.4px at 88% 85%, rgba(196, 181, 253, 0.5) 50%, transparent 50%);
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
  width: 30rem;
  height: 30rem;
  border-radius: 9999px;
  background: radial-gradient(
    circle at 55% 55%,
    rgba(253, 224, 71, 0.5) 0%,
    rgba(251, 146, 60, 0.28) 28%,
    rgba(251, 146, 60, 0.08) 55%,
    transparent 72%
  );
  filter: blur(6px);
  animation: sun-breathe 12s ease-in-out infinite;
}

:global(.dark) .distant-sun {
  background: radial-gradient(
    circle at 55% 55%,
    rgba(253, 224, 71, 0.4) 0%,
    rgba(249, 115, 22, 0.22) 26%,
    rgba(124, 45, 18, 0.15) 52%,
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
  border: 1px solid rgba(99, 102, 241, 0.12);
}

:global(.dark) .sky-orbit {
  border-color: rgba(165, 180, 252, 0.12);
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
  background: radial-gradient(circle at 35% 35%, #fef3c7, #f59e0b);
  box-shadow: 0 0 8px rgba(245, 158, 11, 0.8);
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
  display: flex;
  width: min(100%, 560px);
  align-items: center;
  justify-content: center;
  aspect-ratio: 560 / 460;
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
  stroke: rgba(99, 102, 241, 0.3);
  stroke-width: 1.2;
  stroke-dasharray: 5 7;
}

:global(.dark) .orbit-path {
  stroke: rgba(165, 180, 252, 0.35);
}

.orbit-sun {
  fill: #fbbf24;
  filter: drop-shadow(0 0 10px rgba(251, 191, 36, 0.9));
}

.orbit-planet {
  fill: #c4b5fd;
  filter: drop-shadow(0 0 8px rgba(167, 139, 250, 0.9));
}

.orbit-apex {
  fill: none;
  stroke: rgba(245, 158, 11, 0.8);
  stroke-width: 1.5;
}

.orbit-caption {
  position: absolute;
  bottom: 4%;
  left: 50%;
  transform: translateX(-50%);
  margin: 0;
  white-space: nowrap;
  font-family: ui-monospace, 'Fira Code', monospace;
  font-size: 11px;
  letter-spacing: 0.08em;
  color: rgba(120, 113, 108, 0.9);
}

:global(.dark) .orbit-caption {
  color: rgba(196, 181, 253, 0.75);
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
  background: linear-gradient(145deg, #1b1a38 0%, #0b0a1c 100%);
  border-radius: 14px;
  box-shadow:
    0 25px 50px -12px rgba(2, 4, 16, 0.55),
    0 0 0 1px rgba(196, 181, 253, 0.18),
    0 0 45px rgba(99, 102, 241, 0.18),
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
  background: rgba(30, 27, 62, 0.8);
  border-bottom: 1px solid rgba(196, 181, 253, 0.12);
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
  color: #8b83b8;
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
  color: #a78bfa;
}
.code-comment {
  color: #7c7699;
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
  color: #fcd34d;
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
    0 25px 50px -12px rgba(0, 0, 0, 0.7),
    0 0 0 1px rgba(196, 181, 253, 0.22),
    0 0 55px rgba(99, 102, 241, 0.22),
    inset 0 1px 0 rgba(255, 255, 255, 0.08);
}
</style>
