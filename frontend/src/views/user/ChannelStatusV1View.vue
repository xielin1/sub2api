<template>
  <!-- 1. 已登录：套后台布局；2. 未登录访客：使用模型广场同款独立导航条 -->
  <component :is="isAuthenticated ? AppLayout : 'div'" :class="isAuthenticated ? '' : 'min-h-screen bg-gray-50 dark:bg-dark-950'">
    <PlazaNavBar v-if="!isAuthenticated" />
    <main :class="isAuthenticated ? '' : 'mx-auto max-w-7xl px-4 py-6 sm:px-6 lg:px-8 lg:py-8'">
    <MonitorHero
      :overall-status="overallStatus"
      :interval-seconds="DEFAULT_INTERVAL_SECONDS"
      :window="currentWindow"
      :loading="loading"
      :auto-refresh="autoRefresh"
      @update:window="handleWindowChange"
      @refresh="manualReload"
    />

    <!-- 1. 上方平台分类 tab，点击后下方只展示该平台的监控卡片 -->
    <nav
      v-if="platformTabs.length > 0"
      class="mb-5 flex gap-2 overflow-x-auto"
      role="tablist"
    >
      <button
        v-for="tab in platformTabs"
        :key="tab.value"
        type="button"
        role="tab"
        class="flex items-center gap-2 whitespace-nowrap rounded-xl px-4 py-2.5 text-sm font-medium transition-colors"
        :class="activePlatform === tab.value
          ? 'bg-white text-gray-900 shadow-sm ring-1 ring-gray-900/5 dark:bg-dark-700 dark:text-white dark:ring-dark-600'
          : 'text-gray-600 hover:bg-white/60 hover:text-gray-900 dark:text-dark-400 dark:hover:bg-dark-800 dark:hover:text-white'"
        :aria-selected="activePlatform === tab.value"
        @click="activePlatform = tab.value"
      >
        <ProviderIcon :provider="tab.value" :size="18" />
        <span>{{ tab.label }}</span>
        <span class="text-xs tabular-nums text-gray-400">{{ tab.count }}</span>
      </button>
    </nav>

    <!-- 2. 当前平台的监控卡片 -->
    <MonitorCardGrid
      :items="filteredItems"
      :window="currentWindow"
      :countdown-seconds="countdown"
      :loading="loading"
      :detail-cache="detailCache"
      @card-click="openDetail"
    />

    <MonitorDetailDialog
      :show="showDetail"
      :monitor-id="detailTarget?.id ?? null"
      :title="detailTitle"
      @close="closeDetail"
    />
    </main>
  </component>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onBeforeUnmount, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { extractApiErrorMessage } from '@/utils/apiError'
import {
  list as listChannelMonitorViews,
  status as fetchChannelMonitorDetail,
  type UserMonitorView,
  type UserMonitorDetail,
} from '@/api/channelMonitor'
import AppLayout from '@/components/layout/AppLayout.vue'
import PlazaNavBar from '@/components/modelPlaza/PlazaNavBar.vue'
import { useAuthStore } from '@/stores/auth'
import MonitorHero, {
  type MonitorWindow,
  type OverallStatus,
} from '@/components/user/monitor/MonitorHero.vue'
import MonitorCardGrid from '@/components/user/monitor/MonitorCardGrid.vue'
import MonitorDetailDialog from '@/components/user/MonitorDetailDialog.vue'
import ProviderIcon from '@/components/user/monitor/ProviderIcon.vue'
import { platformLabel } from '@/utils/platformColors'
import { DEFAULT_INTERVAL_SECONDS, STATUS_OPERATIONAL } from '@/constants/channelMonitor'
import { useAutoRefresh } from '@/composables/useAutoRefresh'

const { t } = useI18n()
const appStore = useAppStore()
const authStore = useAuthStore()
// 访客访问时不渲染后台侧边栏
const isAuthenticated = computed(() => authStore.isAuthenticated)

// ── State ──
const items = ref<UserMonitorView[]>([])
const loading = ref(false)
const currentWindow = ref<MonitorWindow>('7d')
const detailCache = reactive<Record<number, UserMonitorDetail>>({})
const showDetail = ref(false)
const detailTarget = ref<UserMonitorView | null>(null)
// 当前选中的平台 tab；为空时默认取第一个平台
const selectedPlatform = ref('')

let abortController: AbortController | null = null

const autoRefresh = useAutoRefresh({
  storageKey: 'channel-status-auto-refresh',
  intervals: [30, 60, 120] as const,
  defaultInterval: DEFAULT_INTERVAL_SECONDS,
  defaultEnabled: true,
  onRefresh: () => reload(true),
  shouldPause: () => document.hidden || loading.value,
})
const countdown = autoRefresh.countdown

// ── Computed ──
const overallStatus = computed<OverallStatus>(() => {
  if (items.value.length === 0) return 'operational'
  for (const it of items.value) {
    if (it.primary_status === 'failed' || it.primary_status === 'error') return 'degraded'
    if (it.primary_status !== STATUS_OPERATIONAL) return 'degraded'
  }
  return 'operational'
})

// 平台 tab 列表：监控的 provider 即其探测的平台（openai / anthropic / grok ...）
const platformTabs = computed(() => {
  // 1. 统计每个平台的监控数
  const counts = new Map<string, number>()
  for (const it of items.value) {
    counts.set(it.provider, (counts.get(it.provider) ?? 0) + 1)
  }
  // 2. 生成 tab，anthropic 平台对用户展示为 Claude
  const tabs = Array.from(counts, ([platform, count]) => ({
    value: platform,
    label: platform === 'anthropic' ? 'Claude' : platformLabel(platform),
    count,
  }))
  // 3. 固定顺序 OpenAI、Claude、Grok，其余平台排在后面
  const order = ['openai', 'anthropic', 'grok']
  const rank = (p: string) => (order.indexOf(p) === -1 ? order.length : order.indexOf(p))
  return tabs.sort((a, b) => rank(a.value) - rank(b.value))
})

// 当前生效的平台：用户选中的平台已不存在时回落到第一个平台
const activePlatform = computed({
  get: () => {
    if (platformTabs.value.some(tab => tab.value === selectedPlatform.value)) return selectedPlatform.value
    return platformTabs.value[0]?.value ?? ''
  },
  set: (value: string) => {
    selectedPlatform.value = value
  },
})

// 当前平台下展示的监控卡片，按倍率从低到高排序，未配置倍率的排在最后
const filteredItems = computed(() => items.value
  .filter(it => it.provider === activePlatform.value)
  .sort((a, b) => (a.rate_multiplier ?? Infinity) - (b.rate_multiplier ?? Infinity)))

const detailTitle = computed(() => {
  return detailTarget.value?.name || t('channelStatus.detailTitle')
})

// ── Loaders ──
async function reload(silent = false) {
  if (abortController) abortController.abort()
  const ctrl = new AbortController()
  abortController = ctrl
  if (!silent) loading.value = true
  try {
    const res = await listChannelMonitorViews({ signal: ctrl.signal })
    if (ctrl.signal.aborted || abortController !== ctrl) return
    items.value = res.items || []
  } catch (err: unknown) {
    const e = err as { name?: string; code?: string }
    if (e?.name === 'AbortError' || e?.code === 'ERR_CANCELED') return
    appStore.showError(extractApiErrorMessage(err, t('channelStatus.loadError')))
  } finally {
    if (abortController === ctrl) {
      if (!silent) loading.value = false
      autoRefresh.resetCountdown()
      abortController = null
    }
  }
}

async function manualReload() {
  await reload(false)
  // After base reload, refresh any cached detail records so non-7d availability
  // values stay in sync without forcing the user to switch tabs again.
  if (currentWindow.value !== '7d') {
    await Promise.all(items.value.map(it => loadDetail(it.id, true)))
  }
}

async function loadDetail(id: number, force = false) {
  if (!force && detailCache[id]) return
  try {
    detailCache[id] = await fetchChannelMonitorDetail(id)
  } catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t('channelStatus.detailLoadError')))
  }
}

async function ensureDetailsForWindow() {
  if (currentWindow.value === '7d') return
  await Promise.all(items.value.map(it => loadDetail(it.id)))
}

// ── Handlers ──
async function handleWindowChange(value: MonitorWindow) {
  currentWindow.value = value
  await ensureDetailsForWindow()
}

function openDetail(row: UserMonitorView) {
  detailTarget.value = row
  showDetail.value = true
}

function closeDetail() {
  showDetail.value = false
  detailTarget.value = null
}

watch(items, () => {
  void ensureDetailsForWindow()
})

watch(
  () => appStore.cachedPublicSettings?.channel_monitor_enabled,
  (enabled) => {
    if (enabled === false) autoRefresh.stop()
    else if (autoRefresh.enabled.value) autoRefresh.start()
  },
)

onMounted(() => {
  // 访客直接打开状态页时，导航条需要站点名/Logo
  void appStore.fetchPublicSettings()
  void reload(false)
  if (appStore.cachedPublicSettings?.channel_monitor_enabled !== false) {
    autoRefresh.setEnabled(autoRefresh.enabled.value)
  }
})

onBeforeUnmount(() => {
  if (abortController) abortController.abort()
})
</script>
