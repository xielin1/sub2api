<template>
  <AppLayout>
    <div class="mx-auto max-w-[1600px] space-y-6" @paste="pasteImage">
      <!-- 1. 明确显示保存能力，未启用对象存储时不能承诺刷新恢复。 -->
      <section class="flex flex-wrap items-center justify-between gap-4 rounded-2xl border border-primary-100 bg-gradient-to-r from-primary-50 to-white p-5 dark:border-primary-900 dark:from-primary-950 dark:to-dark-900">
        <div class="flex items-center gap-3">
          <div class="rounded-xl bg-primary-600 p-3 text-white"><Icon name="sparkles" size="lg" /></div>
          <div><h2 class="font-semibold text-primary-800 dark:text-primary-200">{{ t('imageWorkspace.workspace') }}</h2><p class="mt-1 text-sm text-gray-500 dark:text-gray-400">{{ t('imageWorkspace.automatic') }}</p></div>
        </div>
        <div class="max-w-md rounded-xl bg-white/70 px-4 py-3 text-sm dark:bg-dark-800">
          <p class="font-medium text-primary-700 dark:text-primary-300">{{ t(asynchronous ? 'imageWorkspace.saved' : 'imageWorkspace.sessionOnly') }}</p>
          <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ t(asynchronous ? 'imageWorkspace.savedHint' : 'imageWorkspace.sessionHint') }}</p>
        </div>
      </section>

      <div v-if="pageError" role="alert" class="rounded-xl border border-red-200 bg-red-50 p-4 text-sm text-red-700 dark:border-red-900 dark:bg-red-950 dark:text-red-300">
        {{ pageError }} <button class="ml-3 underline" type="button" @click="initialize">{{ t('imageWorkspace.retry') }}</button>
      </div>
      <div v-if="storageError" role="alert" class="rounded-xl bg-amber-50 p-4 text-sm text-amber-800 dark:bg-amber-950 dark:text-amber-200">{{ t('imageWorkspace.storageError') }}</div>

      <div class="grid items-start gap-6 xl:grid-cols-[360px_minmax(0,1fr)] 2xl:grid-cols-[400px_minmax(0,1fr)]">
        <!-- 2. 表单复用站点样式；上传、拖拽和粘贴共用同一图片校验。 -->
        <form class="card overflow-hidden" @submit.prevent="generate">
          <fieldset :disabled="loading || submitting" class="min-w-0">
            <div class="border-b border-gray-100 p-5 dark:border-dark-700">
              <div class="grid grid-cols-2 gap-1 rounded-xl bg-gray-100 p-1 dark:bg-dark-900">
                <button v-for="item in modes" :key="item" type="button" class="rounded-lg px-3 py-2.5 text-sm font-medium transition-colors" :class="mode === item ? 'bg-white text-primary-700 shadow-sm dark:bg-dark-700 dark:text-primary-300' : 'text-gray-500 dark:text-gray-400'" :aria-pressed="mode === item" @click="mode = item">{{ t(`imageWorkspace.${item}`) }}</button>
              </div>
            </div>
            <div class="space-y-5 p-5">
              <div>
                <label for="image-prompt" class="input-label">{{ t('imageWorkspace.prompt') }}</label>
                <textarea id="image-prompt" v-model="prompt" class="input min-h-[160px] resize-y" maxlength="8000" required :placeholder="t('imageWorkspace.placeholder')" />
                <div class="mt-2 flex gap-3 text-xs text-gray-500"><span class="flex-1">{{ t('imageWorkspace.promptHint') }}</span><span>{{ prompt.length }}/8000</span></div>
              </div>
              <div>
                <label for="image-key" class="input-label">{{ t('imageWorkspace.apiKey') }}</label>
                <select id="image-key" v-model="keyId" required class="input">
                  <option :value="0" disabled>{{ t('imageWorkspace.selectKey') }}</option>
                  <option v-for="key in eligibleKeys" :key="key.id" :value="key.id">{{ key.name }} · {{ key.group?.name }} · {{ key.group?.platform.toUpperCase() }}</option>
                </select>
                <p v-if="selectedKey" class="mt-2 text-xs text-gray-500">{{ selectedKey.quota > 0 ? t('imageWorkspace.quota', { amount: Math.max(0, selectedKey.quota - selectedKey.quota_used).toFixed(2) }) : t('imageWorkspace.unlimited') }}</p>
                <p v-if="!loading && !eligibleKeys.length" class="mt-2 text-sm text-amber-700 dark:text-amber-300">{{ t('imageWorkspace.noKeys') }} <RouterLink to="/keys" class="underline">{{ t('imageWorkspace.manageKeys') }}</RouterLink></p>
              </div>
              <div>
                <label for="image-model" class="input-label">{{ t('imageWorkspace.model') }}</label>
                <select id="image-model" v-model="model" class="input" required :disabled="modelsLoading || !models.length">
                  <option value="" disabled>{{ t(modelsLoading ? 'imageWorkspace.modelsLoading' : 'imageWorkspace.selectModel') }}</option>
                  <option v-for="id in models" :key="id" :value="id">{{ id }}</option>
                </select>
                <p class="mt-2 text-xs text-gray-500">{{ t('imageWorkspace.modelHint') }}</p>
                <p v-if="modelError" role="alert" class="mt-2 text-xs text-amber-700 dark:text-amber-300">{{ modelError }} <button type="button" class="ml-2 underline" @click="loadModels">{{ t('imageWorkspace.retry') }}</button></p>
                <p v-else-if="selectedKey && !modelsLoading && !models.length" class="mt-2 text-xs text-amber-700 dark:text-amber-300">{{ t('imageWorkspace.noModels') }}</p>
              </div>
              <div v-if="mode === 'generate'">
                <label for="image-count" class="input-label">{{ t('imageWorkspace.count') }}</label>
                <select id="image-count" v-model="count" class="input"><option v-for="n in 10" :key="n" :value="n">{{ t('imageWorkspace.imageCount', { n }) }}</option></select>
                <p class="mt-2 text-xs text-gray-500">{{ t('imageWorkspace.billingHint') }}</p>
              </div>
              <div v-else>
                <label for="image-file" class="input-label">{{ t('imageWorkspace.original') }}</label>
                <label class="flex min-h-[160px] cursor-pointer flex-col items-center justify-center gap-2 rounded-xl border-2 border-dashed border-gray-200 p-4 text-center hover:border-primary-400 dark:border-dark-600" :class="{ 'border-primary-400 bg-primary-50 dark:bg-primary-950': dragging }" @dragover.prevent="dragging = true" @dragleave.prevent="dragging = false" @drop.prevent="dropImage">
                  <img v-if="referenceUrl" :src="referenceUrl" :alt="t('imageWorkspace.original')" class="max-h-40 rounded-lg object-contain" />
                  <Icon v-else name="upload" size="xl" class="text-primary-500" />
                  <span class="text-sm text-gray-600 dark:text-gray-300">{{ reference?.name || t('imageWorkspace.upload') }}</span>
                  <span class="text-xs text-gray-500">{{ t('imageWorkspace.uploadHint') }}</span>
                  <input id="image-file" type="file" class="sr-only" accept="image/png,image/jpeg,image/webp" @change="fileChanged" />
                </label>
                <button v-if="reference" type="button" class="mt-2 text-sm text-red-600" @click="clearReference">{{ t('imageWorkspace.removeImage') }}</button>
              </div>
              <button type="submit" class="btn btn-primary w-full gap-2 py-3" :disabled="!ready || !selectedKey || !prompt.trim() || modelsLoading || !models.includes(model) || (mode === 'edit' && !reference) || submitting || processing">
                <Icon :name="submitting || processing ? 'refresh' : 'sparkles'" :class="{ 'animate-spin': submitting || processing }" size="md" />
                {{ t(submitting ? 'imageWorkspace.submitting' : processing ? 'imageWorkspace.processing' : mode === 'edit' ? 'imageWorkspace.startEdit' : 'imageWorkspace.start') }}
              </button>
              <p class="text-center text-xs leading-relaxed text-gray-400">{{ t('imageWorkspace.auditHint') }}</p>
            </div>
          </fieldset>
        </form>

        <!-- 3. 画布和临时作品共享任务记录，失败显示原因且不自动重提计费请求。 -->
        <section class="card min-w-0 overflow-hidden">
          <div class="flex items-center justify-between gap-3 border-b border-gray-100 p-4 dark:border-dark-700">
            <div class="flex gap-1 rounded-xl bg-gray-100 p-1 dark:bg-dark-900">
              <button v-for="tab in tabs" :key="tab" type="button" class="rounded-lg px-4 py-2 text-sm font-medium" :class="activeTab === tab ? 'bg-white text-primary-700 shadow-sm dark:bg-dark-700 dark:text-primary-300' : 'text-gray-500'" :aria-pressed="activeTab === tab" @click="activeTab = tab">{{ t(`imageWorkspace.${tab}`) }}<span v-if="tab === 'history'" class="ml-2 rounded-full bg-primary-50 px-2 py-0.5 text-xs text-primary-700 dark:bg-primary-950 dark:text-primary-300">{{ jobs.length }}</span></button>
            </div>
          </div>
          <div v-if="!shownJobs.length" class="flex min-h-[520px] flex-col items-center justify-center p-8 text-center">
            <div class="mb-6 rounded-3xl border border-dashed border-primary-300 bg-primary-50 p-6 dark:border-primary-800 dark:bg-primary-950"><Icon name="sparkles" size="xl" class="text-primary-500" /></div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-gray-100">{{ t('imageWorkspace.empty') }}</h3>
            <p class="mt-3 max-w-sm text-sm leading-relaxed text-gray-500">{{ t('imageWorkspace.emptyHint') }}</p>
          </div>
          <div v-else class="space-y-6 p-5">
            <article v-for="job in shownJobs" :key="job.id" class="space-y-4">
              <div class="flex items-start justify-between gap-3">
                <div class="min-w-0"><p class="line-clamp-3 whitespace-pre-wrap break-words text-sm text-gray-800 dark:text-gray-200">{{ job.prompt }}</p><p class="mt-2 text-xs text-gray-500">{{ job.model }} · {{ new Date(job.createdAt).toLocaleString() }}</p></div>
                <button v-if="job.status !== 'processing' || job.pollError" type="button" class="rounded-lg p-2 text-gray-400 hover:text-red-500" :aria-label="t('imageWorkspace.delete')" @click="removeJob(job)"><Icon name="trash" size="sm" /></button>
              </div>
              <div v-if="job.status === 'processing'" role="status" class="flex min-h-[240px] flex-col items-center justify-center gap-4 rounded-xl bg-gray-50 p-5 text-center dark:bg-dark-900">
                <Icon name="refresh" size="xl" class="animate-spin text-primary-500" /><p class="font-medium text-gray-700 dark:text-gray-200">{{ t('imageWorkspace.processing') }}</p><p class="max-w-md text-sm text-gray-500">{{ t(job.taskId ? 'imageWorkspace.processingAsync' : 'imageWorkspace.processingSync') }}</p>
              </div>
              <div v-if="job.error || job.pollError" role="alert" class="rounded-xl bg-red-50 p-4 text-sm text-red-700 dark:bg-red-950 dark:text-red-300">
                {{ job.error || job.pollError }}
                <button v-if="job.pollError" type="button" class="ml-3 underline" :disabled="polling" @click="pollJobs(true)">{{ t('imageWorkspace.refreshTask') }}</button>
              </div>
              <div v-if="job.result" class="grid gap-4" :class="job.result.data.length > 1 || activeTab === 'history' ? 'sm:grid-cols-2' : ''">
                <figure v-for="(picture, index) in job.result.data" :key="index" class="overflow-hidden rounded-xl border border-gray-100 dark:border-dark-700">
                  <button type="button" class="block w-full bg-gray-50 dark:bg-dark-900" :aria-label="t('imageWorkspace.preview')" @click="preview = { src: imageSource(picture), prompt: job.prompt }"><img :src="imageSource(picture)" :alt="job.prompt" class="mx-auto max-h-[560px] w-full object-contain" loading="lazy" @error="brokenImage" /></button>
                  <figcaption class="flex flex-wrap items-center justify-between gap-2 p-3">
                    <span class="text-xs text-gray-400">{{ index + 1 }} / {{ job.result.data.length }}</span>
                    <div class="flex gap-2"><button type="button" class="btn btn-secondary btn-sm gap-1" :disabled="imageBusy" @click="download(job, index)"><Icon name="download" size="sm" />{{ t('imageWorkspace.download') }}</button><button type="button" class="btn btn-secondary btn-sm" :disabled="imageBusy || submitting" @click="editResult(job, index)">{{ t('imageWorkspace.continueEdit') }}</button></div>
                  </figcaption>
                  <p v-if="picture.revised_prompt" class="border-t border-gray-100 p-3 text-xs text-gray-500 dark:border-dark-700">{{ picture.revised_prompt }}</p>
                </figure>
              </div>
              <div v-if="job.status === 'failed'" class="flex flex-wrap items-center gap-3"><button type="button" class="btn btn-secondary btn-sm" @click="reusePrompt(job)">{{ t('imageWorkspace.reuse') }}</button><span class="text-xs text-gray-500">{{ t('imageWorkspace.retryHint') }}</span></div>
            </article>
          </div>
        </section>
      </div>
    </div>
    <BaseDialog :show="!!preview" :title="t('imageWorkspace.preview')" width="extra-wide" :close-on-click-outside="true" @close="preview = null"><img v-if="preview" :src="preview.src" :alt="preview.prompt" class="mx-auto max-h-[75vh] object-contain" /></BaseDialog>
  </AppLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { onBeforeRouteLeave } from 'vue-router'
import { saveAs } from 'file-saver'
import AppLayout from '@/components/layout/AppLayout.vue'
import BaseDialog from '@/components/common/BaseDialog.vue'
import Icon from '@/components/icons/Icon.vue'
import { keysAPI } from '@/api/keys'
import { getWorkspaceConfig, listImageModels, submitImage, getImageTask, getImageBlob, type ImageResult, type ImageTask, type WorkspaceImage } from '@/api/imageWorkspace'
import { useAuthStore } from '@/stores/auth'
import { useAppStore } from '@/stores/app'
import type { ApiKey } from '@/types'

interface WorkspaceJob {
  id: string
  taskId?: string
  keyId: number
  prompt: string
  model: string
  mode: 'generate' | 'edit'
  createdAt: number
  expiresAt: number
  status: 'processing' | 'completed' | 'failed'
  result?: ImageResult
  error?: string
  pollError?: string
}

const { t } = useI18n()
const authStore = useAuthStore()
const appStore = useAppStore()
const modes = ['generate', 'edit'] as const
const tabs = ['canvas', 'history'] as const
const mode = ref<'generate' | 'edit'>('generate')
const activeTab = ref<'canvas' | 'history'>('canvas')
const prompt = ref('')
const model = ref('')
const count = ref(1)
const keyId = ref(0)
const keys = ref<ApiKey[]>([])
const models = ref<string[]>([])
const reference = ref<File | null>(null)
const referenceUrl = ref('')
const dragging = ref(false)
const jobs = ref<WorkspaceJob[]>([])
const currentId = ref('')
const preview = ref<{ src: string; prompt: string } | null>(null)
const loading = ref(true)
const ready = ref(false)
const submitting = ref(false)
const polling = ref(false)
const imageBusy = ref(false)
const asynchronous = ref(false)
const pageError = ref('')
const modelError = ref('')
const modelsLoading = ref(false)
const storageError = ref(false)
const lifetime = new AbortController()
let historyRestored = false
let modelController: AbortController | undefined
let pollTimer: ReturnType<typeof setTimeout> | undefined
let expiryTimer: ReturnType<typeof setInterval> | undefined
const storageKey = `image-workspace:${authStore.user!.id}`
const eligibleKeys = computed(() => keys.value.filter(key => key.status === 'active' && key.group?.status === 'active' && key.group.allow_image_generation && ['openai', 'grok'].includes(key.group.platform) && (!key.expires_at || Date.parse(key.expires_at) > Date.now()) && (key.quota === 0 || key.quota_used < key.quota)))
const selectedKey = computed(() => eligibleKeys.value.find(key => key.id === keyId.value))
const processing = computed(() => jobs.value.some(job => job.status === 'processing'))
const shownJobs = computed(() => activeTab.value === 'history' ? jobs.value : jobs.value.filter(job => job.id === currentId.value))

function errorMessage(error: unknown): string {
  // 1. 同时兼容网关 Error 和管理面 API 返回的错误对象。
  return typeof (error as { message?: unknown })?.message === 'string' ? (error as { message: string }).message : t('imageWorkspace.failed')
}

function persistJobs() {
  // 1. 初始化失败或历史损坏时不覆盖原记录，避免定时清理误删尚未恢复的任务。
  if (!historyRestored) return
  // 2. 只保存任务标识和提示词；API 密钥、原图和大体积结果不写入 localStorage。
  try {
    const records = jobs.value.filter(job => job.taskId).map(({ result: _result, error: _error, pollError: _pollError, ...job }) => job)
    localStorage.setItem(storageKey, JSON.stringify(records))
    storageError.value = false
  } catch {
    // 3. 浏览器禁用存储或配额耗尽时保留当前画布，并明确提示无法恢复。
    storageError.value = true
  }
}

function restoreJobs() {
  // 1. 本地存储也是输入边界，拒绝损坏记录，不把任意字段传进网关。
  const records = JSON.parse(localStorage.getItem(storageKey) || '[]')
  if (!Array.isArray(records) || records.some(job => !job || typeof job.id !== 'string' || typeof job.taskId !== 'string' || !Number.isInteger(job.keyId) || typeof job.prompt !== 'string' || typeof job.model !== 'string' || !modes.includes(job.mode) || !Number.isFinite(job.createdAt) || !Number.isFinite(job.expiresAt) || !['processing', 'completed', 'failed'].includes(job.status))) {
    throw new Error(t('imageWorkspace.historyInvalid'))
  }
  // 2. 状态和图片以服务器为准；恢复后查询一次完成和失败任务。
  jobs.value = records.filter(job => job.expiresAt > Date.now()).map(job => ({ id: job.id, taskId: job.taskId, keyId: job.keyId, prompt: job.prompt, model: job.model, mode: job.mode, createdAt: job.createdAt, expiresAt: job.expiresAt, status: job.status }))
  currentId.value = jobs.value[0]?.id || ''
  historyRestored = true
}

async function initialize() {
  // 1. 分页获取所有密钥，保留耗尽密钥用于取回已经付费的结果。
  loading.value = true
  ready.value = false
  pageError.value = ''
  try {
    const config = await getWorkspaceConfig()
    asynchronous.value = config.async_enabled
    const all: ApiKey[] = []
    let page = 1
    let pages = 1
    do {
      const response = await keysAPI.list(page, 100, undefined, { signal: lifetime.signal })
      all.push(...response.items)
      pages = response.pages
      page++
    } while (page <= pages)
    keys.value = all
    keyId.value = eligibleKeys.value[0]?.id || 0
    ready.value = true
    // 2. 存储损坏不阻断新创作，但不能静默伪装成恢复成功。
    if (!jobs.value.length) {
      try { restoreJobs() } catch (error) { storageError.value = true; appStore.showError(errorMessage(error)) }
    }
    await pollJobs(true)
  } catch (error) {
    if (!lifetime.signal.aborted) pageError.value = errorMessage(error)
  } finally {
    loading.value = false
  }
}

async function loadModels() {
  // 1. 切换密钥取消旧模型请求，避免慢响应覆盖新分组的模型。
  modelController?.abort()
  modelController = new AbortController()
  const controller = modelController
  models.value = []
  model.value = ''
  modelError.value = ''
  modelsLoading.value = false
  const key = selectedKey.value
  if (!key) return
  modelsLoading.value = true
  try {
    const available = await listImageModels(key.key, controller.signal)
    // 2. 已切换密钥时丢弃旧响应，当前下拉框只接受最新密钥的模型。
    if (controller.signal.aborted) return
    models.value = available
    model.value = models.value[0] || ''
  } catch (error) {
    if (!controller.signal.aborted) modelError.value = `${t('imageWorkspace.modelsFailed')} ${errorMessage(error)}`
  } finally {
    if (!controller.signal.aborted) modelsLoading.value = false
  }
}

// 1. 选择密钥后自动加载模型；失败时按钮复用同一加载入口。
watch(keyId, loadModels)

async function chooseFile(file: File) {
  // 1. 文件首次进入系统时统一校验 MIME、20 MB 上限和实际可解码性。
  if (submitting.value) return
  if (!['image/png', 'image/jpeg', 'image/webp'].includes(file.type) || file.size === 0 || file.size > 20 * 1024 * 1024) {
    appStore.showError(t('imageWorkspace.invalidFile'))
    return
  }
  const url = URL.createObjectURL(file)
  try {
    const picture = new Image()
    picture.src = url
    await picture.decode()
    clearReference()
    reference.value = file
    referenceUrl.value = url
    mode.value = 'edit'
  } catch {
    URL.revokeObjectURL(url)
    appStore.showError(t('imageWorkspace.invalidFile'))
  }
}

function clearReference() {
  // 1. 替换或移除图片时释放浏览器对象 URL。
  URL.revokeObjectURL(referenceUrl.value)
  reference.value = null
  referenceUrl.value = ''
}

function fileChanged(event: Event) {
  // 1. 文件选择器允许取消，选择后清空控件以便再次选择同一文件。
  const input = event.target as HTMLInputElement
  if (input.files?.[0]) void chooseFile(input.files[0])
  input.value = ''
}

function dropImage(event: DragEvent) {
  // 1. 拖拽和粘贴都委托同一个上传入口。
  dragging.value = false
  if (event.dataTransfer?.files[0]) void chooseFile(event.dataTransfer.files[0])
}

function pasteImage(event: ClipboardEvent) {
  // 1. 只处理图片粘贴，普通提示词文本仍交给浏览器。
  if (event.clipboardData?.files[0]) { event.preventDefault(); void chooseFile(event.clipboardData.files[0]) }
}

async function generate() {
  // 1. 提交边界验证表单；后端仍独立执行权限、模型和余额校验。
  const key = selectedKey.value
  if (!ready.value || submitting.value || processing.value || !key || !prompt.value.trim() || prompt.value.length > 8000 || modelsLoading.value || !models.value.includes(model.value) || !Number.isInteger(count.value) || count.value < 1 || count.value > 10 || (mode.value === 'edit' && !reference.value)) return
  const job: WorkspaceJob = { id: crypto.randomUUID(), keyId: key.id, prompt: prompt.value.trim(), model: model.value.trim(), mode: mode.value, createdAt: Date.now(), expiresAt: Date.now() + 86400000, status: 'processing' }
  jobs.value.unshift(job)
  currentId.value = job.id
  activeTab.value = 'canvas'
  submitting.value = true
  try {
    // 2. 每次点击只提交一次；文生图支持数量，改图固定一张。
    const result = await submitImage(key.key, { model: job.model, prompt: job.prompt, n: mode.value === 'edit' ? 1 : count.value, image: mode.value === 'edit' ? reference.value! : undefined, platform: key.group!.platform as 'openai' | 'grok' }, asynchronous.value, lifetime.signal)
    const current = jobs.value.find(item => item.id === job.id)!
    if ('task_id' in result) {
      current.taskId = result.task_id
      applyTask(current, result)
      persistJobs()
      void pollJobs()
    } else {
      current.status = 'completed'
      current.result = result
    }
  } catch (error) {
    if (!lifetime.signal.aborted) {
      const current = jobs.value.find(item => item.id === job.id)!
      current.status = 'failed'
      current.error = `${errorMessage(error)} ${t('imageWorkspace.unknownOutcome')}`
    }
  } finally {
    submitting.value = false
  }
}

function applyTask(job: WorkspaceJob, task: ImageTask) {
  // 1. 只更新服务端负责的状态，不改变提交时的密钥、提示词和模型。
  job.status = task.status
  job.expiresAt = Math.min(task.expires_at * 1000, job.createdAt + 86400000)
  job.result = task.result
  job.error = task.status === 'failed' ? task.error?.message || t('imageWorkspace.failed') : undefined
  job.pollError = undefined
}

async function pollJobs(restore = false) {
  // 1. 单个轮询循环串行读取任务；页面离开只停止查询，不取消已受理的生图。
  if (polling.value || lifetime.signal.aborted) return
  clearTimeout(pollTimer)
  polling.value = true
  try {
    for (const job of jobs.value.filter(item => item.taskId && (restore || (item.status === 'processing' && !item.pollError)))) {
      const key = keys.value.find(item => item.id === job.keyId)
      if (!key) { job.pollError = t('imageWorkspace.missingKey'); continue }
      try {
        applyTask(job, await getImageTask(key.key, job.taskId!, lifetime.signal))
      } catch (error) {
        if (!lifetime.signal.aborted) job.pollError = errorMessage(error)
      }
    }
    persistJobs()
  } finally {
    polling.value = false
    // 2. 只自动查询健康的进行中任务；查询失败可手动恢复，绝不重提原请求。
    if (!lifetime.signal.aborted && jobs.value.some(job => job.taskId && job.status === 'processing' && !job.pollError)) pollTimer = setTimeout(() => void pollJobs(), 3000)
  }
}

function imageSource(picture: WorkspaceImage): string {
  // 1. 展示响应中的位图，不插入上游 HTML。
  return picture.b64_json ? `data:image/png;base64,${picture.b64_json}` : picture.url!
}

function brokenImage(event: Event) {
  // 1. 临时链接过期或网络失败时在图片占位中给出可理解的说明。
  const image = event.target as HTMLImageElement
  image.alt = t('imageWorkspace.imageUnavailable')
}

async function download(job: WorkspaceJob, index: number) {
  // 1. 异步下载通过任务密钥鉴权，同步下载使用当前结果。
  imageBusy.value = true
  try {
    const key = keys.value.find(item => item.id === job.keyId)
    if (job.taskId && !key) throw new Error(t('imageWorkspace.missingKey'))
    const blob = await getImageBlob(job.result!.data[index], key?.key || '', job.taskId, index)
    const extension = blob.type === 'image/jpeg' ? 'jpg' : blob.type === 'image/webp' ? 'webp' : 'png'
    saveAs(blob, `${job.id}-${index + 1}.${extension}`)
  } catch (error) {
    appStore.showError(`${t('imageWorkspace.downloadFailed')} ${errorMessage(error)}`)
  } finally { imageBusy.value = false }
}

async function editResult(job: WorkspaceJob, index: number) {
  // 1. 继续修改复用下载入口与上传校验，不重新生成原图。
  imageBusy.value = true
  try {
    const key = keys.value.find(item => item.id === job.keyId)
    if (job.taskId && !key) throw new Error(t('imageWorkspace.missingKey'))
    const blob = await getImageBlob(job.result!.data[index], key?.key || '', job.taskId, index)
    const extension = blob.type === 'image/jpeg' ? 'jpg' : blob.type === 'image/webp' ? 'webp' : 'png'
    await chooseFile(new File([blob], `image.${extension}`, { type: blob.type }))
    prompt.value = ''
    window.scrollTo({ top: 0, behavior: 'smooth' })
  } catch (error) { appStore.showError(errorMessage(error)) }
  finally { imageBusy.value = false }
}

function reusePrompt(job: WorkspaceJob) {
  // 1. 只填回表单，不自动触发新的计费请求。
  prompt.value = job.prompt
  // 2. 历史模型只有仍在当前密钥列表中才回填，不能绕过可选范围。
  if (models.value.includes(job.model)) model.value = job.model
  mode.value = job.mode
  window.scrollTo({ top: 0, behavior: 'smooth' })
}

function removeJob(job: WorkspaceJob) {
  // 1. 删除浏览器临时记录不代表删除平台用量或审核记录。
  jobs.value = jobs.value.filter(item => item.id !== job.id)
  currentId.value = jobs.value[0]?.id || ''
  persistJobs()
}

function beforeUnload(event: BeforeUnloadEvent) {
  // 1. 同步请求或尚未获得任务编号时离页会丢失结果，使用浏览器原生提示。
  if (submitting.value || jobs.value.some(job => !job.taskId && job.status === 'processing')) {
    event.preventDefault()
    event.returnValue = ''
  }
}

onBeforeRouteLeave(() => {
  // 1. 站内导航不会触发 beforeunload，同步生成离页同样需要明确确认。
  if (submitting.value || jobs.value.some(job => !job.taskId && job.status === 'processing')) {
    return window.confirm(t('imageWorkspace.leaveConfirm'))
  }
  return true
})

onMounted(() => {
  // 1. 初始化完成后恢复任务；到期清理仅作用于浏览器中的临时作品。
  void initialize()
  window.addEventListener('beforeunload', beforeUnload)
  expiryTimer = setInterval(() => {
    jobs.value = jobs.value.filter(job => job.expiresAt > Date.now())
    if (!jobs.value.some(job => job.id === currentId.value)) currentId.value = jobs.value[0]?.id || ''
    persistJobs()
  }, 60000)
})

onUnmounted(() => {
  // 1. 释放请求、定时器、原图预览和离页监听。
  lifetime.abort()
  modelController?.abort()
  clearTimeout(pollTimer)
  clearInterval(expiryTimer)
  clearReference()
  window.removeEventListener('beforeunload', beforeUnload)
})
</script>
