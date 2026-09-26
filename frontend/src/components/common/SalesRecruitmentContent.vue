<template>
  <div ref="contentRef" class="sales-recruitment">
    <!-- 1. 文案使用真实文字，图片只负责插画，管理员改字不需要重新生图。 -->
    <header class="recruitment-hero">
      <div class="recruitment-hero-copy">
        <p class="recruitment-eyebrow">{{ t(`salesRecruitment.eyebrow.${page}`) }}</p>
        <h2>{{ page === 'intro' ? config.title : t(`salesRecruitment.${page}Title`) }}</h2>
        <p class="recruitment-subtitle">{{ page === 'intro' ? config.subtitle : t(`salesRecruitment.${page}Subtitle`) }}</p>
        <span v-if="page === 'intro' && config.badge" class="recruitment-badge">{{ config.badge }}</span>
      </div>
      <img :src="heroImage" alt="" class="recruitment-art" />
    </header>

    <div class="recruitment-body">
      <!-- 2. 浮窗和管理预览复用同一套内容，规则没有平行实现。 -->
      <template v-if="page === 'intro'">
        <p v-if="config.intro" class="recruitment-intro">{{ config.intro }}</p>
        <template v-if="config.tiers?.length">
          <h3 class="recruitment-section-title">{{ t('salesRecruitment.commission') }}</h3>
          <div class="recruitment-tiers">
            <div v-for="(tier, index) in config.tiers" :key="index" class="recruitment-tier" :class="{ 'recruitment-tier-featured': index === 0 }">
              <span class="recruitment-tier-number">{{ String(config.tiers.length - index).padStart(2, '0') }}</span>
              <span class="recruitment-tier-label">{{ tier.label }}</span>
              <div class="recruitment-tier-value"><strong>{{ tier.rate }}</strong><span>{{ t('salesRecruitment.cashCommission') }}</span></div>
            </div>
          </div>
        </template>
        <template v-if="config.benefits?.length">
          <h3 class="recruitment-section-title">{{ t('salesRecruitment.benefits') }}</h3>
          <div class="recruitment-benefits">
            <div v-for="(benefit, index) in config.benefits" :key="index" class="recruitment-benefit">
              <span class="recruitment-benefit-icon"><Icon :name="benefitIcons[index % benefitIcons.length]" size="lg" /></span>
              <h4>{{ benefit.title }}</h4>
              <p>{{ benefit.description }}</p>
            </div>
          </div>
        </template>
        <button type="button" class="recruitment-primary" @click="emit('navigate', 'contact')">{{ config.button_text }}<Icon name="arrowRight" size="md" /></button>
        <button type="button" class="recruitment-link" @click="emit('navigate', 'rules')">{{ t('salesRecruitment.viewRules') }}<Icon name="chevronRight" size="sm" /></button>
        <p v-if="config.footer" class="recruitment-footer">{{ config.footer }}</p>
      </template>

      <template v-else-if="page === 'rules'">
        <button type="button" class="recruitment-back" @click="emit('navigate', 'intro')"><Icon name="arrowLeft" size="sm" />{{ t('salesRecruitment.back') }}</button>
        <!-- 3. 规则以纯文本保留换行，拒绝通过配置注入 HTML 或脚本。 -->
        <div class="recruitment-rules">{{ config.rules }}</div>
        <button type="button" class="recruitment-primary" @click="emit('navigate', 'contact')">{{ config.button_text }}<Icon name="arrowRight" size="md" /></button>
      </template>

      <template v-else>
        <button type="button" class="recruitment-back" @click="emit('navigate', 'intro')"><Icon name="arrowLeft" size="sm" />{{ t('salesRecruitment.back') }}</button>
        <div class="recruitment-contact">
          <p class="recruitment-wechat"><Icon name="chat" size="lg" />{{ t('salesRecruitment.wechatContact') }}</p>
          <img v-if="qrImage" :src="qrImage" :alt="t('salesRecruitment.qrAlt')" class="recruitment-qr" />
          <h3>{{ qrImage ? t('salesRecruitment.scan') : t('salesRecruitment.addWechat') }}</h3>
          <p class="recruitment-contact-note">{{ config.contact_note }}</p>
          <div v-if="config.wechat_id" class="recruitment-wechat-id">
            <div><span>{{ t('salesRecruitment.wechatId') }}</span><strong>{{ config.wechat_id }}</strong></div>
            <button type="button" class="btn btn-secondary" @click="copyToClipboard(config.wechat_id)"><Icon :name="copied ? 'check' : 'copy'" size="sm" />{{ copied ? t('salesRecruitment.copied') : t('salesRecruitment.copyWechat') }}</button>
          </div>
          <p v-if="!config.wechat_id && !qrImage" class="recruitment-contact-note">{{ t('salesRecruitment.contactMissing') }}</p>
          <p v-if="qrImage" class="recruitment-tip">{{ t('salesRecruitment.mobileHint') }}</p>
        </div>
      </template>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import Icon from '@/components/icons/Icon.vue'
import { useClipboard } from '@/composables/useClipboard'
import { sanitizeUrl } from '@/utils/url'
import type { SalesRecruitment } from '@/types'
import defaultHero from '@/assets/images/sales-partner.jpg'

const props = defineProps<{ config: SalesRecruitment; page: 'intro' | 'rules' | 'contact' }>()
const emit = defineEmits<{ navigate: [page: 'intro' | 'rules' | 'contact'] }>()
const { t } = useI18n()
const { copied, copyToClipboard } = useClipboard()
// 1. 未上传封面是正常状态，使用随应用发布的原创配图；二维码只能使用真实配置。
const heroImage = computed(() => props.config.hero_image ? sanitizeUrl(props.config.hero_image, { allowDataUrl: true }) : defaultHero)
const qrImage = computed(() => sanitizeUrl(props.config.wechat_qr_code || '', { allowDataUrl: true }))
const benefitIcons = ['shield', 'users', 'document'] as const
const contentRef = ref<HTMLElement | null>(null)
watch(() => props.page, async () => {
  // 2. 从长介绍切换到规则或联系方式时回到顶部，并将焦点移入当前页。
  await nextTick()
  contentRef.value?.closest('.modal-body, .recruitment-panel')?.scrollTo({ top: 0 })
  contentRef.value?.querySelector<HTMLButtonElement>('button')?.focus({ preventScroll: true })
})
</script>

<style scoped>
/* 1. 招募卡片局部使用暖黄与蓝色，不影响现有站点主题。 */
.sales-recruitment { --recruitment-ink: #15213c; --recruitment-muted: #66738a; color: var(--recruitment-ink); background: white; overflow: hidden; border-radius: 24px; }
.recruitment-hero { position: relative; display: flex; align-items: center; min-height: 226px; background: #fff1b9; isolation: isolate; }
.recruitment-hero-copy { position: relative; z-index: 1; width: 68%; padding: 30px; }
.recruitment-eyebrow { font-size: 11px; font-weight: 800; letter-spacing: .2em; color: #2859b7; margin-bottom: 14px; }
.recruitment-hero h2 { font-size: clamp(26px, 3vw, 40px); line-height: 1.2; font-weight: 850; letter-spacing: -.04em; overflow-wrap: anywhere; }
.recruitment-subtitle { margin-top: 12px; font-size: 15px; color: #685722; white-space: pre-wrap; overflow-wrap: anywhere; }
.recruitment-badge { display: inline-block; margin-top: 18px; padding: 5px 10px; border: 1px solid #c7a73a66; border-radius: 30px; font-size: 11px; color: #705808; background: #fffcdfb3; }
.recruitment-art { position: absolute; right: 0; height: 100%; width: 42%; object-fit: cover; object-position: 45% center; mask-image: linear-gradient(to right, transparent, black 24%); }
.recruitment-body { padding: 24px 30px; }
.recruitment-intro { color: var(--recruitment-muted); font-size: 14px; line-height: 1.8; margin-bottom: 22px; white-space: pre-wrap; overflow-wrap: anywhere; }
.recruitment-section-title { display: flex; align-items: center; gap: 14px; font-size: 16px; font-weight: 750; margin: 0 0 18px; }
.recruitment-section-title::before, .recruitment-section-title::after { content: ''; height: 1px; flex: 1; background: #e3e9f2; }
.recruitment-tiers { display: grid; gap: 10px; margin-bottom: 28px; }
.recruitment-tier { display: flex; align-items: center; gap: 16px; padding: 15px 18px; border: 1px solid #dce5f8; border-radius: 16px; background: #f6f8ff; }
.recruitment-tier-featured { background: #fffbef; border-color: #f5d083; }
.recruitment-tier-number { display: grid; place-items: center; width: 36px; height: 36px; flex-shrink: 0; border-radius: 50%; background: #e5edff; color: #2864dd; font-size: 13px; font-weight: 750; }
.recruitment-tier-featured .recruitment-tier-number { background: #ffdf87; color: #875514; }
.recruitment-tier-label { flex: 1; min-width: 0; font-size: 15px; font-weight: 650; overflow-wrap: anywhere; }
.recruitment-tier-value { text-align: right; max-width: 45%; overflow-wrap: anywhere; }
.recruitment-tier-value strong { display: block; font-size: 30px; line-height: 1.15; color: #2864dd; }
.recruitment-tier-featured .recruitment-tier-value strong { color: #e68720; }
.recruitment-tier-value span { font-size: 11px; color: var(--recruitment-muted); }
.recruitment-benefits { display: grid; grid-template-columns: repeat(auto-fit, minmax(120px, 1fr)); gap: 18px; margin-bottom: 26px; }
.recruitment-benefit { text-align: center; min-width: 0; overflow-wrap: anywhere; }
.recruitment-benefit-icon { display: inline-flex; padding: 12px; border-radius: 16px; background: #edf3ff; color: #2864dd; }
.recruitment-benefit:nth-child(2) .recruitment-benefit-icon { background: #fff3e6; color: #e68720; }
.recruitment-benefit:nth-child(3) .recruitment-benefit-icon { background: #eaf7f2; color: #14846b; }
.recruitment-benefit h4 { font-size: 14px; font-weight: 700; margin: 10px 0 5px; }
.recruitment-benefit p { font-size: 12px; line-height: 1.7; color: var(--recruitment-muted); white-space: pre-wrap; }
.recruitment-primary { display: flex; justify-content: center; align-items: center; gap: 12px; width: 100%; padding: 15px 20px; border-radius: 15px; background: linear-gradient(110deg, #ffdb43, #ffc511); color: #3e3000; font-size: 17px; font-weight: 750; box-shadow: 0 10px 24px #eab30822; overflow-wrap: anywhere; }
.recruitment-primary:hover { background: #f4bf06; }
.recruitment-link { display: flex; align-items: center; justify-content: center; gap: 4px; margin: 16px auto 0; font-size: 13px; color: #52617c; }
.recruitment-footer { margin-top: 20px; font-size: 11px; text-align: center; color: var(--recruitment-muted); white-space: pre-wrap; overflow-wrap: anywhere; }
.recruitment-back { display: inline-flex; align-items: center; gap: 8px; font-size: 13px; color: #52617c; margin-bottom: 22px; }
.recruitment-rules { white-space: pre-wrap; overflow-wrap: anywhere; line-height: 1.95; font-size: 15px; margin-bottom: 28px; }
.recruitment-contact { text-align: center; }
.recruitment-wechat { display: flex; align-items: center; justify-content: center; gap: 8px; color: #118466; font-size: 15px; font-weight: 650; }
.recruitment-qr { display: block; width: min(248px, 100%); aspect-ratio: 1; object-fit: contain; background: white; padding: 14px; margin: 22px auto; border: 1px solid #dde5f1; border-radius: 20px; box-shadow: 0 14px 40px #17294e0d; }
.recruitment-contact h3 { font-size: 20px; font-weight: 750; margin: 20px 0 10px; }
.recruitment-contact-note { font-size: 14px; line-height: 1.9; color: var(--recruitment-muted); white-space: pre-wrap; overflow-wrap: anywhere; }
.recruitment-wechat-id { display: flex; align-items: center; justify-content: space-between; gap: 12px; flex-wrap: wrap; margin-top: 24px; padding: 16px; border-radius: 14px; background: #f6f8fc; border: 1px solid #dfe6f0; text-align: left; }
.recruitment-wechat-id div { min-width: 0; overflow-wrap: anywhere; }
.recruitment-wechat-id span { display: block; font-size: 12px; color: var(--recruitment-muted); }
.recruitment-wechat-id strong { font-size: 18px; }
.recruitment-tip { margin-top: 16px; padding: 14px; border-radius: 12px; font-size: 12px; text-align: left; background: #fff8d9; color: #896719; }
.sales-recruitment button:focus-visible { outline: 3px solid #2864dd; outline-offset: 3px; }
/* 2. 移动端保持单列，允许较长文案换行。 */
@media (max-width: 540px) {
  .recruitment-hero { min-height: 196px; }
  .recruitment-hero-copy { padding: 24px 18px; width: 75%; }
  .recruitment-art { width: 40%; opacity: .6; }
  .recruitment-body { padding: 22px 16px; }
  .recruitment-tier { gap: 10px; padding: 14px 12px; }
  .recruitment-tier-label { font-size: 13px; }
  .recruitment-tier-value strong { font-size: 25px; }
  .recruitment-benefits { grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 10px; }
  .recruitment-benefit h4 { font-size: 13px; }
}
/* 3. 右下角浮动面板空间有限，统一缩小横幅、间距与字号；管理端预览弹窗不受影响。 */
:global(.recruitment-panel) .sales-recruitment { border-radius: 18px; }
:global(.recruitment-panel) .recruitment-hero { min-height: 132px; }
:global(.recruitment-panel) .recruitment-hero-copy { width: 72%; padding: 18px 18px 16px; }
:global(.recruitment-panel) .recruitment-eyebrow { font-size: 10px; letter-spacing: .14em; margin-bottom: 8px; }
:global(.recruitment-panel) .recruitment-hero h2 { font-size: 20px; letter-spacing: -.02em; }
:global(.recruitment-panel) .recruitment-subtitle { margin-top: 6px; font-size: 12px; }
:global(.recruitment-panel) .recruitment-badge { margin-top: 10px; padding: 3px 8px; font-size: 10px; }
:global(.recruitment-panel) .recruitment-art { width: 38%; }
:global(.recruitment-panel) .recruitment-body { padding: 16px 18px 18px; }
:global(.recruitment-panel) .recruitment-intro { font-size: 13px; line-height: 1.7; margin-bottom: 14px; }
:global(.recruitment-panel) .recruitment-section-title { font-size: 13px; gap: 10px; margin-bottom: 10px; }
:global(.recruitment-panel) .recruitment-tiers { gap: 8px; margin-bottom: 18px; }
:global(.recruitment-panel) .recruitment-tier { gap: 10px; padding: 10px 12px; border-radius: 12px; }
:global(.recruitment-panel) .recruitment-tier-number { width: 28px; height: 28px; font-size: 11px; }
:global(.recruitment-panel) .recruitment-tier-label { font-size: 13px; }
:global(.recruitment-panel) .recruitment-tier-value strong { font-size: 20px; }
:global(.recruitment-panel) .recruitment-tier-value span { font-size: 10px; }
:global(.recruitment-panel) .recruitment-benefits { grid-template-columns: repeat(3, minmax(0, 1fr)); gap: 10px; margin-bottom: 18px; }
:global(.recruitment-panel) .recruitment-benefit-icon { padding: 8px; border-radius: 12px; }
:global(.recruitment-panel) .recruitment-benefit h4 { font-size: 12px; margin: 6px 0 3px; }
:global(.recruitment-panel) .recruitment-benefit p { font-size: 11px; line-height: 1.5; }
:global(.recruitment-panel) .recruitment-primary { padding: 11px 16px; border-radius: 12px; font-size: 14px; gap: 8px; }
:global(.recruitment-panel) .recruitment-link { margin-top: 10px; font-size: 12px; }
:global(.recruitment-panel) .recruitment-footer { margin-top: 12px; }
:global(.recruitment-panel) .recruitment-back { margin-bottom: 14px; font-size: 12px; }
:global(.recruitment-panel) .recruitment-rules { font-size: 13px; line-height: 1.8; margin-bottom: 18px; }
:global(.recruitment-panel) .recruitment-wechat { font-size: 13px; }
:global(.recruitment-panel) .recruitment-qr { width: min(180px, 100%); padding: 10px; margin: 14px auto; border-radius: 14px; }
:global(.recruitment-panel) .recruitment-contact h3 { font-size: 16px; margin: 12px 0 6px; }
:global(.recruitment-panel) .recruitment-contact-note { font-size: 12px; line-height: 1.7; }
:global(.recruitment-panel) .recruitment-wechat-id { margin-top: 14px; padding: 12px; }
:global(.recruitment-panel) .recruitment-wechat-id strong { font-size: 15px; }
:global(.recruitment-panel) .recruitment-tip { margin-top: 10px; padding: 10px; font-size: 11px; }
/* 4. 暗色模式只调整内容区；插画和品牌横幅保持原色。 */
:global(.dark .sales-recruitment) { --recruitment-ink: #e9eef8; --recruitment-muted: #a3afc2; background: #172033; }
:global(.dark .recruitment-hero) { color: #15213c; }
:global(.dark .recruitment-tier) { background: #202d46; border-color: #344666; }
:global(.dark .recruitment-tier-featured) { background: #332e23; border-color: #665431; }
:global(.dark .recruitment-tier-value strong) { color: #8bb3ff; }
:global(.dark .recruitment-tier-featured .recruitment-tier-value strong) { color: #ffcb78; }
:global(.dark .recruitment-wechat-id) { background: #202d46; border-color: #344666; }
:global(.dark .recruitment-back), :global(.dark .recruitment-link) { color: #b9c7df; }
</style>
