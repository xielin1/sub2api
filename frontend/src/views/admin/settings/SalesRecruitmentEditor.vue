<template>
  <section class="space-y-5 border-t border-gray-100 pt-6 dark:border-dark-700">
    <!-- 1. 配置归属系统设置，使用外层表单的保存动作，预览不写入服务器。 -->
    <div class="flex items-start justify-between gap-4">
      <div><h3 class="font-semibold text-gray-900 dark:text-white">{{ t('salesRecruitment.settingsTitle') }}</h3><p class="mt-1 text-xs leading-6 text-gray-500 dark:text-gray-400">{{ t('salesRecruitment.settingsHint') }}</p></div>
      <Toggle v-model="model.enabled" />
    </div>
    <div class="flex flex-wrap items-center gap-5 text-sm">
      <label class="flex items-center gap-2"><input v-model="model.floating_enabled" type="checkbox" class="rounded" />{{ t('salesRecruitment.floatingEnabled') }}</label>
      <button type="button" class="btn btn-secondary btn-sm" @click="preview = true">{{ t('salesRecruitment.preview') }}</button>
    </div>
    <p class="text-xs text-gray-500">{{ t('salesRecruitment.visibilityHint') }}</p>
    <div class="grid gap-4 sm:grid-cols-2">
      <label v-for="field in textFields" :key="field" class="block text-sm font-medium">
        {{ t(`salesRecruitment.fields.${field}`) }}
        <input v-model="model[field]" type="text" maxlength="500" class="input mt-2" />
      </label>
      <label class="block text-sm font-medium sm:col-span-2">{{ t('salesRecruitment.fields.intro') }}<textarea v-model="model.intro" rows="3" maxlength="500" class="input mt-2" /></label>
    </div>
    <div>
      <label class="mb-2 block text-sm font-medium">{{ t('salesRecruitment.heroImage') }}</label>
      <ImageUpload v-model="model.hero_image" :max-size="500 * 1024" :upload-label="t('salesRecruitment.upload')" :remove-label="t('salesRecruitment.remove')" :hint="t('salesRecruitment.heroHint')" />
      <input v-model="model.hero_image" type="text" class="input mt-2 text-xs" :placeholder="t('salesRecruitment.imageUrl')" :aria-label="t('salesRecruitment.heroImage')" />
    </div>
    <!-- 2. 档位与权益直接编辑数组，复用现有增删表单方式。 -->
    <div class="space-y-3">
      <h4 class="text-sm font-semibold">{{ t('salesRecruitment.commission') }}</h4>
      <div v-for="(tier, index) in model.tiers" :key="index" class="flex items-center gap-2">
        <input v-model="tier.label" class="input min-w-0 flex-1" maxlength="100" :aria-label="t('salesRecruitment.tierLabel')" :placeholder="t('salesRecruitment.tierLabel')" />
        <input v-model="tier.rate" class="input w-24 sm:w-32" maxlength="30" :aria-label="t('salesRecruitment.tierRate')" :placeholder="t('salesRecruitment.tierRate')" />
        <button type="button" class="btn btn-secondary px-2" :aria-label="t('salesRecruitment.remove')" @click="model.tiers.splice(index, 1)"><Icon name="x" size="sm" /></button>
      </div>
      <button type="button" class="btn btn-secondary btn-sm" :disabled="model.tiers.length >= 6" @click="model.tiers.push({ label: '', rate: '' })">{{ t('salesRecruitment.addTier') }}</button>
    </div>
    <div class="space-y-3">
      <h4 class="text-sm font-semibold">{{ t('salesRecruitment.benefits') }}</h4>
      <div v-for="(benefit, index) in model.benefits" :key="index" class="flex items-start gap-2">
        <div class="min-w-0 flex-1 space-y-2">
          <input v-model="benefit.title" class="input" maxlength="100" :aria-label="t('salesRecruitment.benefitTitle')" :placeholder="t('salesRecruitment.benefitTitle')" />
          <textarea v-model="benefit.description" class="input" rows="2" maxlength="500" :aria-label="t('salesRecruitment.benefitDescription')" :placeholder="t('salesRecruitment.benefitDescription')" />
        </div>
        <button type="button" class="btn btn-secondary px-2" :aria-label="t('salesRecruitment.remove')" @click="model.benefits.splice(index, 1)"><Icon name="x" size="sm" /></button>
      </div>
      <button type="button" class="btn btn-secondary btn-sm" :disabled="model.benefits.length >= 6" @click="model.benefits.push({ title: '', description: '' })">{{ t('salesRecruitment.addBenefit') }}</button>
    </div>
    <label class="block text-sm font-medium">{{ t('salesRecruitment.rulesTitle') }}<textarea v-model="model.rules" rows="10" maxlength="20000" class="input mt-2 leading-7" /><span class="mt-1 block text-xs font-normal text-gray-500">{{ t('salesRecruitment.rulesHint') }}</span></label>
    <!-- 3. 联系方式必须由站长配置真实值，不复制参考站二维码，也不生成虚假二维码。 -->
    <label class="block text-sm font-medium">{{ t('salesRecruitment.wechatId') }}<input v-model="model.wechat_id" class="input mt-2" maxlength="500" /></label>
    <div>
      <label class="mb-2 block text-sm font-medium">{{ t('salesRecruitment.qrAlt') }}</label>
      <ImageUpload v-model="model.wechat_qr_code" :max-size="500 * 1024" :upload-label="t('salesRecruitment.upload')" :remove-label="t('salesRecruitment.remove')" :hint="t('salesRecruitment.qrHint')" />
      <input v-model="model.wechat_qr_code" type="text" class="input mt-2 text-xs" :placeholder="t('salesRecruitment.imageUrl')" :aria-label="t('salesRecruitment.qrAlt')" />
    </div>
    <label class="block text-sm font-medium">{{ t('salesRecruitment.contactNote') }}<textarea v-model="model.contact_note" class="input mt-2" rows="2" maxlength="500" /></label>
    <SalesRecruitmentDialog :show="preview" :config="model" @close="preview = false" />
  </section>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import type { SalesRecruitment } from '@/types'
import Toggle from '@/components/common/Toggle.vue'
import ImageUpload from '@/components/common/ImageUpload.vue'
import Icon from '@/components/icons/Icon.vue'
import SalesRecruitmentDialog from '@/components/common/SalesRecruitmentDialog.vue'

// 1. 与外层系统设置共享当前草稿，统一提交，不引入单独保存接口。
const model = defineModel<SalesRecruitment>({ required: true })
const { t } = useI18n()
const preview = ref(false)
const textFields = ['title', 'subtitle', 'badge', 'button_text', 'footer'] as const
</script>
