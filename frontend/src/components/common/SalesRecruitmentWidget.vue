<template>
  <template v-if="config?.enabled && visible">
    <!-- 1. 入口可拖动到任意位置；点击在入口旁展开面板，不遮罩页面，再次点击收起。 -->
    <div ref="widgetRef" class="recruitment-widget" :style="widgetStyle">
      <button
        type="button"
        class="recruitment-launcher"
        :class="{ 'recruitment-launcher-dragging': dragging }"
        :aria-label="config.title"
        :aria-expanded="open"
        @pointerdown="onPointerDown"
        @pointermove="onPointerMove"
        @pointerup="onPointerUp"
        @pointercancel="onPointerCancel"
        @keydown.enter.prevent="toggle"
        @keydown.space.prevent="toggle"
      >
        <Icon name="users" size="md" />
        <span><strong>{{ config.title }}</strong><small>{{ config.badge || config.subtitle }}</small></span>
        <Icon :name="open ? 'chevronDown' : 'chevronUp'" size="sm" />
      </button>
    </div>
    <div v-if="open" class="recruitment-panel" :style="panelStyle" role="dialog" :aria-label="config.title">
      <button type="button" class="recruitment-panel-close" :aria-label="t('salesRecruitment.collapse')" @click="open = false"><Icon name="x" size="sm" /></button>
      <SalesRecruitmentContent :config="config" :page="page" @navigate="page = $event" />
    </div>
  </template>
</template>

<script setup lang="ts">
import { computed, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { apiClient } from '@/api/client'
import type { SalesRecruitment } from '@/types'
import Icon from '@/components/icons/Icon.vue'
import SalesRecruitmentContent from '@/components/common/SalesRecruitmentContent.vue'

// 入口位置保存在本机浏览器，只影响当前用户的界面布局。
const POSITION_KEY = 'sales-recruitment-launcher-position'
// 入口与视口边缘、面板与入口之间的间距。
const EDGE_GAP = 12
const PANEL_GAP = 10
// 指针移动超过该距离才视为拖动，避免点击时轻微抖动被当成拖动。
const DRAG_THRESHOLD = 4

const route = useRoute()
const appStore = useAppStore()
const authStore = useAuthStore()
const { t } = useI18n()
const open = ref(false)
const page = ref<'intro' | 'rules' | 'contact'>('intro')
const config = ref<SalesRecruitment | null>(null)
const widgetRef = ref<HTMLElement | null>(null)
// position 为入口左上角坐标；null 表示使用默认的右下角位置。
const position = ref<{ x: number; y: number } | null>(readPosition())
const launcherRect = ref<DOMRect | null>(null)
const viewport = ref({ width: window.innerWidth, height: window.innerHeight })
const dragging = ref(false)
let dragStart: { pointerX: number; pointerY: number; x: number; y: number } | null = null

// 1. 登录后所有页面（用户工作台和管理后台）展示；匿名页面读不到鉴权配置，不展示。
const visible = computed(() => authStore.isAuthenticated && config.value?.floating_enabled)

watch([() => authStore.isAuthenticated, () => appStore.cachedPublicSettings], async ([authenticated], _, onCleanup) => {
  // 2. 登录或管理员保存设置后读取鉴权接口；退出立即清除内容并取消未完成请求。
  config.value = null
  open.value = false
  if (!authenticated) return
  const controller = new AbortController()
  onCleanup(() => controller.abort())
  try {
    const { data } = await apiClient.get<SalesRecruitment | null>('/settings/sales-recruitment', { signal: controller.signal })
    config.value = data
  } catch (error) {
    // 3. 退出导致的取消不是错误，真实接口故障交给现有提示组件显示。
    if (!controller.signal.aborted) appStore.showError(t('salesRecruitment.loadFailed'))
  }
}, { immediate: true })

watch(() => route.path, () => {
  // 4. 导航后收起面板，避免旧内容停留在新页面上。
  open.value = false
})

watch(open, (value) => {
  // 5. 每次展开回到介绍页，并按入口当前位置计算面板位置。
  if (!value) return
  page.value = 'intro'
  measureLauncher()
})

function readPosition(): { x: number; y: number } | null {
  // 1. 本地存储可能被清空或手动改坏，无法解析时回到默认位置。
  try {
    const saved = JSON.parse(localStorage.getItem(POSITION_KEY) || 'null')
    if (saved && Number.isFinite(saved.x) && Number.isFinite(saved.y)) return { x: saved.x, y: saved.y }
  } catch {
    // 2. 解析失败按未保存处理。
  }
  return null
}

function clampPosition(x: number, y: number): { x: number; y: number } {
  // 1. 入口始终完整留在视口内，窗口缩小后也不会被挤出屏幕。
  const el = widgetRef.value
  const width = el?.offsetWidth ?? 0
  const height = el?.offsetHeight ?? 0
  return {
    x: Math.min(Math.max(x, EDGE_GAP), Math.max(EDGE_GAP, viewport.value.width - width - EDGE_GAP)),
    y: Math.min(Math.max(y, EDGE_GAP), Math.max(EDGE_GAP, viewport.value.height - height - EDGE_GAP))
  }
}

function measureLauncher() {
  // 1. 面板定位依赖入口的实时位置和尺寸。
  launcherRect.value = widgetRef.value?.getBoundingClientRect() ?? null
}

function toggle() {
  // 1. 点击入口在展开和收起之间切换。
  open.value = !open.value
}

function onPointerDown(event: PointerEvent) {
  // 1. 只响应主按钮；记录起点，等移动超过阈值后才进入拖动。
  if (event.button !== 0 || !widgetRef.value) return
  const rect = widgetRef.value.getBoundingClientRect()
  dragStart = { pointerX: event.clientX, pointerY: event.clientY, x: rect.left, y: rect.top }
  ;(event.currentTarget as HTMLElement).setPointerCapture(event.pointerId)
}

function onPointerMove(event: PointerEvent) {
  // 1. 超过阈值后跟随指针移动，并限制在视口内。
  if (!dragStart) return
  const dx = event.clientX - dragStart.pointerX
  const dy = event.clientY - dragStart.pointerY
  if (!dragging.value && Math.hypot(dx, dy) < DRAG_THRESHOLD) return
  dragging.value = true
  position.value = clampPosition(dragStart.x + dx, dragStart.y + dy)
  // 2. 面板展开时跟随入口一起移动。
  if (open.value) measureLauncher()
}

function onPointerUp() {
  // 1. 没有发生拖动时按点击处理；拖动结束则保存位置。
  if (!dragStart) return
  if (dragging.value) {
    localStorage.setItem(POSITION_KEY, JSON.stringify(position.value))
  } else {
    toggle()
  }
  dragStart = null
  dragging.value = false
}

function onPointerCancel() {
  // 1. 系统中断手势时放弃本次操作，保留已移动到的位置。
  dragStart = null
  dragging.value = false
}

function onResize() {
  // 1. 窗口尺寸变化后把入口拉回视口内，并重新计算面板位置。
  viewport.value = { width: window.innerWidth, height: window.innerHeight }
  if (position.value) position.value = clampPosition(position.value.x, position.value.y)
  if (open.value) measureLauncher()
}

function onKeydown(event: KeyboardEvent) {
  // 1. 展开时按 Escape 收起。
  if (event.key === 'Escape' && open.value) open.value = false
}

onMounted(() => {
  window.addEventListener('resize', onResize)
  window.addEventListener('keydown', onKeydown)
})
onBeforeUnmount(() => {
  window.removeEventListener('resize', onResize)
  window.removeEventListener('keydown', onKeydown)
})

const widgetStyle = computed(() => {
  // 1. 用户拖动过则使用保存的位置，否则使用 CSS 中的右下角默认位置。
  if (!position.value) return {}
  return { left: `${position.value.x}px`, top: `${position.value.y}px`, right: 'auto', bottom: 'auto' }
})

const panelStyle = computed(() => {
  // 1. 面板宽度固定、不超过视口；高度取入口上方或下方较大的可用空间。
  const rect = launcherRect.value
  const { width: vw, height: vh } = viewport.value
  const width = Math.min(380, vw - EDGE_GAP * 2)
  if (!rect) return { width: `${width}px`, right: `${EDGE_GAP}px`, bottom: `${EDGE_GAP}px` }
  const spaceAbove = rect.top - PANEL_GAP - EDGE_GAP
  const spaceBelow = vh - rect.bottom - PANEL_GAP - EDGE_GAP
  const openUp = spaceAbove >= spaceBelow
  // 2. 入口在屏幕左半边时面板左对齐入口，右半边时右对齐入口，再夹回视口内。
  const alignedLeft = rect.left + rect.width / 2 < vw / 2 ? rect.left : rect.right - width
  const left = Math.min(Math.max(alignedLeft, EDGE_GAP), vw - width - EDGE_GAP)
  const style: Record<string, string> = {
    width: `${width}px`,
    left: `${left}px`,
    maxHeight: `${Math.min(560, openUp ? spaceAbove : spaceBelow)}px`
  }
  if (openUp) style.bottom = `${vh - rect.top + PANEL_GAP}px`
  else style.top = `${rect.bottom + PANEL_GAP}px`
  return style
})
</script>

<style scoped>
/* 1. 默认停在右下角；拖动后改用保存的坐标。 */
.recruitment-widget { position: fixed; z-index: 35; right: max(20px, env(safe-area-inset-right)); bottom: max(20px, env(safe-area-inset-bottom)); max-width: calc(100vw - 24px); }
/* 2. 入口使用站点主色（陶土暖橙），浅底深字；禁用触摸滚动以支持手机拖动。 */
.recruitment-launcher { display: flex; align-items: center; gap: 10px; max-width: 240px; padding: 9px 14px; color: #a24a33; background: #fcf5f1; border: 1px solid #f2cdbd; border-radius: 999px; box-shadow: 0 4px 16px #6c342714; cursor: grab; touch-action: none; user-select: none; transition: box-shadow .15s; }
.recruitment-launcher:hover { box-shadow: 0 6px 20px #6c342724; }
.recruitment-launcher-dragging { cursor: grabbing; box-shadow: 0 10px 28px #6c342733; }
.recruitment-launcher span { display: grid; text-align: left; min-width: 0; }
.recruitment-launcher strong { font-size: 13px; line-height: 1.3; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.recruitment-launcher small { font-size: 11px; line-height: 1.3; max-width: 160px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; opacity: .75; }
.recruitment-widget button:focus-visible, .recruitment-panel button:focus-visible { outline: 3px solid #d66b4d; outline-offset: 3px; }
/* 3. 展开面板：贴着入口出现，不遮罩页面，内容超出时面板内部滚动。 */
.recruitment-panel { position: fixed; z-index: 36; overflow-y: auto; overscroll-behavior: contain; border-radius: 18px; border: 1px solid #f2cdbd; box-shadow: 0 16px 48px #2c2a2433; background: white; animation: recruitment-pop .16s ease-out; }
.recruitment-panel-close { position: sticky; top: 8px; float: right; margin: 8px 8px -36px 0; z-index: 2; display: grid; place-items: center; width: 28px; height: 28px; border-radius: 50%; background: #ffffffd9; color: #a24a33; border: 1px solid #f2cdbd; }
@keyframes recruitment-pop { from { opacity: 0; transform: translateY(6px) scale(.98); } to { opacity: 1; transform: none; } }
/* 4. 手机端只保留图标和标题，减少对页面内容的遮挡。 */
@media (max-width: 640px) {
  .recruitment-widget { right: max(12px, env(safe-area-inset-right)); bottom: max(12px, env(safe-area-inset-bottom)); }
  .recruitment-launcher { padding: 8px 12px; gap: 8px; }
  .recruitment-launcher small { display: none; }
}
/* 5. 深色模式使用首页暖黑底色。 */
:global(.dark .recruitment-launcher) { color: #e9ad96; background: #26261f; border-color: #3a382f; box-shadow: 0 4px 16px #00000040; }
:global(.dark .recruitment-panel) { border-color: #3a382f; box-shadow: 0 16px 48px #00000066; }
:global(.dark .recruitment-panel-close) { color: #e9ad96; background: #1b1c19d9; border-color: #3a382f; }
</style>
