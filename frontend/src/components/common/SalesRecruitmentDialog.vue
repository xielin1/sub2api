<template>
  <BaseDialog :show="show" :title="config.title" width="normal" :close-on-click-outside="true" @close="emit('close')">
    <SalesRecruitmentContent :config="config" :page="page" @navigate="page = $event" />
  </BaseDialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import type { SalesRecruitment } from '@/types'
import BaseDialog from '@/components/common/BaseDialog.vue'
import SalesRecruitmentContent from '@/components/common/SalesRecruitmentContent.vue'

const props = defineProps<{ show: boolean; config: SalesRecruitment }>()
const emit = defineEmits<{ close: [] }>()
const page = ref<'intro' | 'rules' | 'contact'>('intro')
// 1. 每次重新打开回到介绍页；复用已有弹窗的焦点、Escape 和滚动锁定。
watch(() => props.show, (show) => {
  if (show) page.value = 'intro'
})
</script>

<style scoped>
/* 1. 只调整包含招募内容的弹窗，保留通用弹窗其余行为。 */
:global(.modal-content:has(> .modal-body > .sales-recruitment)) { width: min(660px, 100%); max-width: 660px; border-radius: 24px; position: relative; overflow: hidden; }
:global(.modal-content:has(> .modal-body > .sales-recruitment) > .modal-body) { padding: 0; }
:global(.modal-content:has(> .modal-body > .sales-recruitment) > .modal-header) { position: absolute; top: 12px; right: 14px; padding: 0; border: 0; z-index: 2; }
:global(.modal-content:has(> .modal-body > .sales-recruitment) > .modal-header > h3) { position: absolute; width: 1px; height: 1px; overflow: hidden; clip-path: inset(50%); }
:global(.modal-content:has(> .modal-body > .sales-recruitment) > .modal-header > button) { margin: 0; background: #fff9df; color: #776024; border: 1px solid #d5bc7166; border-radius: 50%; }
</style>
