<template>
  <template v-if="config?.enabled && visible">
    <!-- 1. 入口不自动弹出；同一会话可收起，避免持续遮挡页面操作。 -->
    <div class="recruitment-widget">
      <button v-if="!collapsed" type="button" class="recruitment-collapse" :aria-label="t('salesRecruitment.collapse')" @click="collapse"><Icon name="x" size="xs" /></button>
      <button type="button" class="recruitment-launcher" :class="{ 'recruitment-launcher-small': collapsed }" :aria-label="config.title" @click="show = true">
        <Icon name="users" size="lg" />
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
// 1. 登录是必要条件；默认首页和所有匿名页面都没有招募入口。
const visible = computed(() => authStore.isAuthenticated && config.value?.floating_enabled && (
  route.path === '/admin/dashboard' ||
  (route.meta.requiresAuth && !route.meta.requiresAdmin && !route.path.startsWith('/admin'))
))
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
.recruitment-widget { position: fixed; z-index: 35; right: max(20px, env(safe-area-inset-right)); bottom: max(20px, env(safe-area-inset-bottom)); max-width: calc(100vw - 40px); }
.recruitment-launcher { display: flex; align-items: center; gap: 12px; max-width: 290px; padding: 14px 20px; color: #533e08; background: linear-gradient(120deg, #ffe581, #ffcd30); border: 1px solid #edbf39; border-radius: 999px; box-shadow: 0 6px 28px #8e690c26, 0 0 0 4px #fff8dbbb; transition: transform .15s; }
.recruitment-launcher:hover { transform: translateY(-2px); }
.recruitment-launcher span { display: grid; text-align: left; min-width: 0; }
.recruitment-launcher strong { font-size: 14px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
.recruitment-launcher small { font-size: 11px; max-width: 190px; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; opacity: .8; }
.recruitment-launcher-small { padding: 14px; }
.recruitment-collapse { position: absolute; top: -10px; right: -5px; display: grid; place-items: center; width: 24px; height: 24px; background: white; border: 1px solid #eadbae; border-radius: 50%; color: #776024; }
.recruitment-widget button:focus-visible { outline: 3px solid #2864dd; outline-offset: 3px; }
</style>
