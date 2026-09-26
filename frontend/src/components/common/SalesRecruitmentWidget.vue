<template>
  <template v-if="config?.enabled && visible">
    <!-- 1. 入口不自动弹出；同一会话可收起，避免持续遮挡页面操作。 -->
    <div class="recruitment-widget">
      <button v-if="!collapsed" type="button" class="recruitment-collapse" :aria-label="t('salesRecruitment.collapse')" @click="collapse"><Icon name="x" size="xs" /></button>
      <button type="button" class="recruitment-launcher" :class="{ 'recruitment-launcher-small': collapsed }" :aria-label="config.title" @click="show = true">
        <Icon name="users" size="md" />
        <span v-if="!collapsed"><strong>{{ config.title }}</strong><small>{{ config.badge || config.subtitle }}</small></span>
        <Icon v-if="!collapsed" name="chevronUp" size="sm" />
      </button>
    </div>
    <SalesRecruitmentDialog :show="show" :config="config" @close="show = false" />
  </template>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAppStore } from '@/stores/app'
import { useAuthStore } from '@/stores/auth'
import { apiClient } from '@/api/client'
import type { SalesRecruitment } from '@/types'
import Icon from '@/components/icons/Icon.vue'
import SalesRecruitmentDialog from '@/components/common/SalesRecruitmentDialog.vue'

const route = useRoute()
const appStore = useAppStore()
const authStore = useAuthStore()
const { t } = useI18n()
const show = ref(false)
const collapsed = ref(false)
const config = ref<SalesRecruitment | null>(null)
// 1. 登录后所有页面（用户工作台和管理后台）固定展示；匿名页面读不到鉴权配置，不展示。
const visible = computed(() => authStore.isAuthenticated && config.value?.floating_enabled)
watch([() => authStore.isAuthenticated, () => appStore.cachedPublicSettings], async ([authenticated], _, onCleanup) => {
  // 2. 登录或管理员保存设置后读取鉴权接口；退出立即清除内容并取消未完成请求。
  config.value = null
  show.value = false
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
function collapse() {
  // 2. 收起只影响当前会话中的组件状态，不写入业务配置。
  collapsed.value = true
}
watch(() => route.path, () => {
  // 3. 导航后关闭旧弹窗，避免覆盖新页面。
  show.value = false
})
</script>

<style scoped>
/* 1. 全局常驻入口：缩小尺寸，减少对页面右下角操作区的遮挡。 */
.recruitment-widget { position: fixed; z-index: 35; right: max(20px, env(safe-area-inset-right)); bottom: max(20px, env(safe-area-inset-bottom)); max-width: calc(100vw - 40px); }
/* 2. 配色沿用站点主色（陶土暖橙），浅底深字，不抢页面主按钮的视觉。 */
.recruitment-launcher { display: flex; align-items: center; gap: 10px; max-width: 240px; padding: 9px 14px; color: #a24a33; background: #fcf5f1; border: 1px solid #f2cdbd; border-radius: 999px; box-shadow: 0 4px 16px #6c342714; transition: transform .15s, box-shadow .15s; }
.recruitment-launcher:hover { transform: translateY(-2px); box-shadow: 0 6px 20px #6c342724; }
.recruitment-launcher span { display: grid; text-align: left; min-width: 0; }
.recruitment-launcher strong { font-size: 13px; line-height: 1.3; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.recruitment-launcher small { font-size: 11px; line-height: 1.3; max-width: 160px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; opacity: .75; }
.recruitment-launcher-small { padding: 11px; }
.recruitment-collapse { position: absolute; top: -8px; right: -4px; display: grid; place-items: center; width: 20px; height: 20px; background: white; border: 1px solid #f2cdbd; border-radius: 50%; color: #a24a33; }
.recruitment-widget button:focus-visible { outline: 3px solid #d66b4d; outline-offset: 3px; }
/* 3. 手机端只保留图标和标题，减少对页面底部内容的遮挡。 */
@media (max-width: 640px) {
  .recruitment-widget { right: max(12px, env(safe-area-inset-right)); bottom: max(12px, env(safe-area-inset-bottom)); }
  .recruitment-launcher { padding: 8px 12px; gap: 8px; }
  .recruitment-launcher small { display: none; }
}
/* 4. 深色模式使用首页暖黑底色。 */
:global(.dark .recruitment-launcher) { color: #e9ad96; background: #26261f; border-color: #3a382f; box-shadow: 0 4px 16px #00000040; }
:global(.dark .recruitment-collapse) { color: #e9ad96; background: #1b1c19; border-color: #3a382f; }
</style>
