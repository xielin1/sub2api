<template>
  <AppLayout>
    <!-- 1. 页头：说明检测方式，提供手动刷新 -->
    <section class="mb-5 flex flex-col gap-3 rounded-2xl border border-gray-200/80 bg-white/70 p-5 dark:border-dark-700/70 dark:bg-dark-800/60 sm:flex-row sm:items-center sm:justify-between">
      <div>
        <h1 class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ t('iqRadar.title') }}</h1>
        <p class="mt-1 max-w-3xl text-sm leading-6 text-gray-500 dark:text-gray-400">{{ t('iqRadar.subtitle') }}</p>
      </div>
      <button class="btn btn-secondary shrink-0" :disabled="loading" @click="load()">
        {{ t('common.refresh') }}
      </button>
    </section>

    <!-- 2. 首次加载骨架 -->
    <div v-if="loading && items.length === 0" class="grid grid-cols-1 gap-5 md:grid-cols-2 xl:grid-cols-3">
      <div v-for="i in 3" :key="i" class="h-72 animate-pulse rounded-2xl bg-white/70 dark:bg-dark-800/60"></div>
    </div>

    <EmptyState
      v-else-if="items.length === 0"
      :title="t('iqRadar.emptyTitle')"
      :description="t('iqRadar.emptyDescription')"
    />

    <!-- 3. 每个检测分组一张卡片 -->
    <div v-else class="grid grid-cols-1 gap-5 md:grid-cols-2 xl:grid-cols-3">
      <article
        v-for="item in items"
        :key="item.id"
        class="flex flex-col rounded-2xl border border-gray-200/80 bg-white/70 p-5 shadow-card dark:border-dark-700/70 dark:bg-dark-800/60"
      >
        <div class="flex items-start gap-3">
          <span class="grid h-9 w-9 flex-shrink-0 place-items-center rounded-xl ring-1 ring-black/5 dark:ring-white/10">
            <ProviderIcon :provider="item.provider" :size="20" />
          </span>
          <div class="min-w-0 flex-1">
            <div class="truncate text-base font-semibold text-gray-900 dark:text-gray-100">{{ item.name }}</div>
            <div class="mt-0.5 flex min-w-0 items-center gap-1.5">
              <span
                v-if="item.group_name"
                class="inline-flex min-w-0 items-center truncate rounded-md bg-gray-100 px-1.5 py-0.5 text-[10px] font-medium text-gray-600 dark:bg-dark-700 dark:text-gray-300"
              >
                {{ item.group_name }}
              </span>
            </div>
            <div class="mt-1 truncate font-mono text-xs text-gray-500 dark:text-gray-400">{{ item.model }}</div>
          </div>
          <span class="flex-shrink-0 rounded-full px-2.5 py-1 text-xs font-semibold" :class="badgeClass(item.history[0]?.status)">
            {{ statusText(item.history[0]?.status) }}
          </span>
        </div>

        <!-- 4. 通过率与最近一次耗时 -->
        <div class="mt-5 grid grid-cols-2 gap-2">
          <div class="rounded-xl border border-gray-100 bg-gray-50/80 p-3 dark:border-dark-700/50 dark:bg-dark-900/40">
            <div class="text-[10px] font-semibold uppercase tracking-wider text-gray-400">
              {{ t('iqRadar.passRate') }} · {{ t('iqRadar.passRateHint', { n: item.history.length }) }}
            </div>
            <div class="mt-1.5 font-mono text-lg font-bold tabular-nums text-gray-900 dark:text-gray-100">
              {{ passRate(item) }}<span v-if="item.history.length" class="ml-0.5 text-xs font-normal text-gray-400">%</span>
            </div>
          </div>
          <div class="rounded-xl border border-gray-100 bg-gray-50/80 p-3 dark:border-dark-700/50 dark:bg-dark-900/40">
            <div class="text-[10px] font-semibold uppercase tracking-wider text-gray-400">{{ t('iqRadar.lastElapsed') }}</div>
            <div class="mt-1.5 font-mono text-lg font-bold tabular-nums text-gray-900 dark:text-gray-100">
              {{ formatSeconds(item.history[0]?.latency_ms) }}<span v-if="item.history[0]?.latency_ms != null" class="ml-0.5 text-xs font-normal text-gray-400">s</span>
            </div>
          </div>
        </div>

        <!-- 5. 时间线：左旧右新，颜色对应通过 / 降智 / 请求失败 -->
        <div class="mt-4 border-t border-gray-100 pt-3 dark:border-dark-700/60">
          <div class="mb-2 flex justify-between text-[10px] font-semibold uppercase tracking-widest text-gray-400">
            <span>{{ t('iqRadar.timeline', { n: TIMELINE_LENGTH }) }}</span>
            <span>{{ t('iqRadar.every', { n: Math.round(item.interval_seconds / 60) }) }}</span>
          </div>
          <div class="flex h-5 w-full items-end gap-[2px]">
            <div
              v-for="(bar, idx) in timelineBars(item)"
              :key="idx"
              class="min-w-0 flex-1 rounded-sm"
              :class="bar.colorClass"
              :style="{ height: bar.heightPct + '%' }"
              :title="bar.title"
            ></div>
          </div>
          <div class="mt-1 flex justify-between text-[9px] uppercase tracking-widest text-gray-400">
            <span>{{ t('iqRadar.past') }}</span>
            <span>{{ t('iqRadar.now') }}</span>
          </div>
        </div>

        <!-- 6. 最近几次记录：未通过时展示回答结尾或错误信息 -->
        <div class="mt-4 border-t border-gray-100 pt-3 dark:border-dark-700/60">
          <div class="mb-1.5 text-[10px] font-semibold uppercase tracking-widest text-gray-400">{{ t('iqRadar.recent') }}</div>
          <p v-if="item.history.length === 0" class="text-xs text-gray-500 dark:text-gray-400">{{ t('iqRadar.noRecords') }}</p>
          <ul v-else class="space-y-1.5">
            <li v-for="record in item.history.slice(0, RECENT_LENGTH)" :key="record.checked_at" class="text-xs text-gray-600 dark:text-gray-300">
              <div class="flex items-center gap-1.5">
                <span class="h-1.5 w-1.5 flex-shrink-0 rounded-full" :class="STATUS_COLOR[record.status] ?? STATUS_COLOR.empty"></span>
                <span class="tabular-nums">{{ formatTime(record.checked_at) }}</span>
                <span class="font-medium">{{ statusText(record.status) }}</span>
                <span v-if="record.latency_ms != null" class="tabular-nums text-gray-400">{{ formatSeconds(record.latency_ms) }}s</span>
              </div>
              <p v-if="record.message" class="mt-0.5 break-all pl-3 text-[11px] text-gray-400">{{ record.message }}</p>
            </li>
          </ul>
        </div>
      </article>
    </div>

    <p class="mt-5 text-xs leading-5 text-gray-400">{{ t('iqRadar.legend') }}</p>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, onBeforeUnmount } from 'vue'
import { useI18n } from 'vue-i18n'
import AppLayout from '@/components/layout/AppLayout.vue'
import EmptyState from '@/components/common/EmptyState.vue'
import ProviderIcon from '@/components/user/monitor/ProviderIcon.vue'
import { useAppStore } from '@/stores/app'
import { extractApiErrorMessage } from '@/utils/apiError'
import { iqRadar, type IQRadarItem } from '@/api/channelMonitor'

const { t } = useI18n()
const appStore = useAppStore()

// 时间线固定展示的格数（与后端返回的最多 60 条一致），最近记录列表条数
const TIMELINE_LENGTH = 60
const RECENT_LENGTH = 5
// 智商雷达检测间隔以分钟计，页面每分钟同步一次即可
const REFRESH_INTERVAL_MS = 60_000

const STATUS_COLOR: Record<string, string> = {
  operational: 'bg-emerald-500',
  failed: 'bg-red-500',
  error: 'bg-gray-400 dark:bg-dark-500',
  empty: 'bg-gray-200 dark:bg-dark-700',
}
const STATUS_HEIGHT: Record<string, number> = {
  operational: 100,
  failed: 45,
  error: 30,
  empty: 15,
}

const items = ref<IQRadarItem[]>([])
const loading = ref(false)
let timer: ReturnType<typeof setInterval> | null = null

// 拉取智商雷达数据；失败只提示，不清空已展示的数据
async function load() {
  loading.value = true
  try {
    const res = await iqRadar()
    items.value = res.items || []
  } catch (err: unknown) {
    appStore.showError(extractApiErrorMessage(err, t('iqRadar.loadError')))
  } finally {
    loading.value = false
  }
}

function statusText(status?: string): string {
  if (status === 'operational' || status === 'failed' || status === 'error') return t(`iqRadar.status.${status}`)
  return t('iqRadar.status.unknown')
}

function badgeClass(status?: string): string {
  switch (status) {
    case 'operational':
      return 'bg-emerald-100 text-emerald-700 dark:bg-emerald-500/15 dark:text-emerald-300'
    case 'failed':
      return 'bg-red-100 text-red-700 dark:bg-red-500/15 dark:text-red-300'
    default:
      return 'bg-gray-100 text-gray-600 dark:bg-dark-700 dark:text-gray-300'
  }
}

// 通过率 = 答对次数 / 有效判分次数；请求失败不计入，避免把上游故障算成降智
function passRate(item: IQRadarItem): string {
  const judged = item.history.filter((r) => r.status === 'operational' || r.status === 'failed')
  if (judged.length === 0) return '—'
  const passed = judged.filter((r) => r.status === 'operational').length
  return ((passed / judged.length) * 100).toFixed(0)
}

function formatSeconds(ms?: number | null): string {
  if (ms == null) return '—'
  return (ms / 1000).toFixed(1)
}

function formatTime(iso: string): string {
  return new Date(iso).toLocaleString(undefined, { month: '2-digit', day: '2-digit', hour: '2-digit', minute: '2-digit', hour12: false })
}

// 后端记录最新在前，这里翻转为左旧右新，并在左侧补空格保持格数稳定
function timelineBars(item: IQRadarItem) {
  const real = item.history.slice(0, TIMELINE_LENGTH).reverse()
  const bars = Array.from({ length: TIMELINE_LENGTH - real.length }, () => ({
    colorClass: STATUS_COLOR.empty,
    heightPct: STATUS_HEIGHT.empty,
    title: '',
  }))
  for (const record of real) {
    bars.push({
      colorClass: STATUS_COLOR[record.status] ?? STATUS_COLOR.empty,
      heightPct: STATUS_HEIGHT[record.status] ?? STATUS_HEIGHT.empty,
      title: `${formatTime(record.checked_at)} · ${statusText(record.status)} · ${formatSeconds(record.latency_ms)}s`,
    })
  }
  return bars
}

onMounted(() => {
  void load()
  timer = setInterval(() => {
    if (!document.hidden) void load()
  }, REFRESH_INTERVAL_MS)
})

onBeforeUnmount(() => {
  if (timer) clearInterval(timer)
})
</script>
