import { apiClient, buildGatewayUrl } from './client'

export interface WorkspaceImage {
  url?: string
  b64_json?: string
  revised_prompt?: string
}

export interface ImageResult {
  data: WorkspaceImage[]
}

export interface ImageTask {
  task_id: string
  status: 'processing' | 'completed' | 'failed'
  created_at: number
  expires_at: number
  result?: ImageResult
  error?: { message: string }
}

export interface ImageRequest {
  model: string
  prompt: string
  n: number
  image?: File
  platform: 'openai' | 'grok'
}

// 1. 管理面配置继续使用登录令牌，生图请求单独使用用户选择的 API 密钥。
export async function getWorkspaceConfig(): Promise<{ async_enabled: boolean }> {
  const { data } = await apiClient.get<{ async_enabled: boolean }>('/image-workspace/config')
  if (typeof data.async_enabled !== 'boolean') throw new Error('Invalid image workspace configuration')
  return data
}

async function imageRequest(apiKey: string, path: string, init: RequestInit = {}): Promise<Response> {
  // 1. 复用网关地址规则，不经过管理面 Axios 的登录令牌拦截器。
  const response = await fetch(buildGatewayUrl(path), {
    ...init,
    headers: { ...init.headers, Authorization: `Bearer ${apiKey}` },
  })
  // 2. 下游 HTTP 错误在边界解析，非 JSON 的代理错误保留状态码。
  if (!response.ok) {
    const body = await response.text()
    let message = `HTTP ${response.status}: ${response.statusText}`
    if (response.headers.get('content-type')?.includes('application/json')) {
      const error = JSON.parse(body)
      message = error.error?.message || error.message || message
    }
    throw new Error(message)
  }
  return response
}

export async function listImageModels(apiKey: string, signal: AbortSignal): Promise<string[]> {
  // 1. 用当前密钥查询网关模型目录，再筛选图片模型供下拉框选择。
  const response = await imageRequest(apiKey, '/v1/models', { signal })
  const body = await response.json()
  if (!Array.isArray(body.data) || body.data.some((item: { id?: unknown }) => typeof item?.id !== 'string')) {
    throw new Error('Invalid model list')
  }
  return body.data.map((item: { id: string }) => item.id).filter((id: string) => /image|dall-e|flux|imagen|seedream/i.test(id))
}

export function validateImageResult(value: ImageResult): ImageResult {
  // 1. 下游结果必须至少有一张图片，防止把空响应显示为生成成功。
  if (!Array.isArray(value?.data) || value.data.length === 0) throw new Error('Image response contains no images')
  for (const item of value.data) {
    if (typeof item?.b64_json === 'string' && item.b64_json.length > 0) continue
    if (typeof item?.url !== 'string' || !/^https?:\/\//i.test(item.url)) throw new Error('Invalid image URL')
  }
  return value
}

function validateTask(task: ImageTask): ImageTask {
  // 1. 任务协议只接受已知状态，缺少任务标识或过期时间不能进入恢复队列。
  if (!task || typeof task.task_id !== 'string' || !task.task_id ||
      !['processing', 'completed', 'failed'].includes(task.status) ||
      !Number.isFinite(task.expires_at) || !Number.isFinite(task.created_at)) {
    throw new Error('Invalid image task response')
  }
  if (task.status === 'completed') validateImageResult(task.result!)
  return task
}

export async function submitImage(apiKey: string, input: ImageRequest, asynchronous: boolean, signal: AbortSignal): Promise<ImageTask | ImageResult> {
  // 1. OpenAI 改图使用网关与 CPR 都支持的 JSON；Grok 沿用 multipart 上传。
  const path = `/v1/images/${input.image ? 'edits' : 'generations'}${asynchronous ? '/async' : ''}`
  let body: BodyInit
  let headers: Record<string, string> = {}
  if (input.image && input.platform === 'grok') {
    const form = new FormData()
    form.set('model', input.model)
    form.set('prompt', input.prompt)
    form.set('n', String(input.n))
    form.set('response_format', 'b64_json')
    form.set('image', input.image)
    body = form
  } else {
    headers = { 'Content-Type': 'application/json' }
    let images: { image_url: string }[] | undefined
    if (input.image) {
      const imageUrl = await new Promise<string>((resolve, reject) => {
        const reader = new FileReader()
        reader.onload = () => resolve(reader.result as string)
        reader.onerror = () => reject(reader.error)
        reader.readAsDataURL(input.image!)
      })
      images = [{ image_url: imageUrl }]
    }
    body = JSON.stringify({ model: input.model, prompt: input.prompt, n: input.n, response_format: 'b64_json', images })
  }
  // 2. 提交只执行一次；网络错误不能自动重提，以免重复扣费。
  const response = await imageRequest(apiKey, path, { method: 'POST', headers, body, signal })
  const result = await response.json()
  return asynchronous ? validateTask(result) : validateImageResult(result)
}

export async function getImageTask(apiKey: string, taskId: string, signal: AbortSignal): Promise<ImageTask> {
  // 1. 始终使用创建任务的同一密钥查询，后端负责验证任务归属。
  const response = await imageRequest(apiKey, `/v1/images/tasks/${encodeURIComponent(taskId)}`, { signal })
  return validateTask(await response.json())
}

export async function getImageBlob(image: WorkspaceImage, apiKey: string, taskId: string | undefined, index: number): Promise<Blob> {
  // 1. 异步结果通过受鉴权的同源下载接口，避免对象存储 CORS 阻断下载和继续改图。
  let blob: Blob
  if (taskId) {
    const response = await imageRequest(apiKey, `/v1/images/tasks/${encodeURIComponent(taskId)}/images/${index}`, { signal: AbortSignal.timeout(60000) })
    blob = await response.blob()
  } else if (image.b64_json) {
    // 2. 同步结果优先使用返回的 base64，避免把 API 密钥发往图片 URL。
    const bytes = Uint8Array.from(atob(image.b64_json), char => char.charCodeAt(0))
    const type = bytes[0] === 0xff && bytes[1] === 0xd8 ? 'image/jpeg'
      : bytes[0] === 0x52 && bytes[1] === 0x49 ? 'image/webp' : 'image/png'
    blob = new Blob([bytes], { type })
  } else {
    const response = await fetch(image.url!, { credentials: 'omit', signal: AbortSignal.timeout(60000) })
    if (!response.ok) throw new Error(`HTTP ${response.status}: ${response.statusText}`)
    blob = await response.blob()
  }
  // 3. 图片内容是外部输入，只允许受支持的位图，拒绝 HTML/SVG 等可执行内容。
  if (!['image/png', 'image/jpeg', 'image/webp'].includes(blob.type.split(';')[0])) throw new Error('Unsupported image content type')
  return blob
}
