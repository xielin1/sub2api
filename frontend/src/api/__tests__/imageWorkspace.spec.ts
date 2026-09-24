import { afterEach, describe, expect, it, vi } from 'vitest'
import { getImageBlob, getImageTask, listImageModels, submitImage, validateImageResult } from '../imageWorkspace'

describe('图片工作台网关协议', () => {
  afterEach(() => vi.unstubAllGlobals())

  it('生成和改图沿用网关协议，原图使用 multipart 且不覆盖 boundary', async () => {
    // 1. 返回同步位图，校验真实提交的地址、密钥、数量和请求体。
    const fetcher = vi.fn().mockResolvedValue(new Response(JSON.stringify({ data: [{ b64_json: 'aW1hZ2U=' }] })))
    vi.stubGlobal('fetch', fetcher)
    const signal = new AbortController().signal
    await submitImage('test-key', { model: 'image-model', prompt: '画一只猫', n: 2, platform: 'openai' }, false, signal)
    expect(fetcher.mock.calls[0][0]).toMatch(/\/v1\/images\/generations$/)
    expect(JSON.parse(fetcher.mock.calls[0][1].body)).toEqual({ model: 'image-model', prompt: '画一只猫', n: 2, response_format: 'b64_json' })
    expect(fetcher.mock.calls[0][1].headers.Authorization).toBe('Bearer test-key')
    // 2. 异步改图返回任务编号，不自动二次提交。
    fetcher.mockImplementation(() => Promise.resolve(new Response(JSON.stringify({ task_id: 'task-1', status: 'processing', created_at: 1, expires_at: 999 }))))
    const original = new File(['image'], 'cat.png', { type: 'image/png' })
    await submitImage('edit-key', { model: 'image-model', prompt: '换成蓝色', n: 1, image: original, platform: 'grok' }, true, signal)
    const [url, request] = fetcher.mock.calls[1]
    expect(url).toMatch(/\/v1\/images\/edits\/async$/)
    expect(request.headers).toEqual({ Authorization: 'Bearer edit-key' })
    expect(request.body.get('image')).toBe(original)
    expect(request.body.get('prompt')).toBe('换成蓝色')
    expect(fetcher).toHaveBeenCalledTimes(2)
    // 3. OpenAI/CPR 改图发送 JSON data URL，防止 multipart 被 CPR 拒绝。
    await submitImage('edit-key', { model: 'image-model', prompt: '换成蓝色', n: 1, image: original, platform: 'openai' }, true, signal)
    const requestBody = JSON.parse(fetcher.mock.calls[2][1].body)
    expect(requestBody.images).toEqual([{ image_url: 'data:image/png;base64,aW1hZ2U=' }])
    expect(fetcher.mock.calls[2][1].headers['Content-Type']).toBe('application/json')
  })

  it('拒绝空图片和未知任务状态，保留网关失败原因', async () => {
    // 1. 空成功响应、可执行 URL 与未知状态不能显示为生成成功。
    expect(() => validateImageResult({ data: [] })).toThrow('no images')
    expect(() => validateImageResult({ data: [{ url: 'javascript:alert(1)' }] })).toThrow('Invalid image URL')
    const fetcher = vi.fn().mockResolvedValue(new Response(JSON.stringify({ task_id: 'task-1', status: 'unknown', created_at: 1, expires_at: 2 })))
    vi.stubGlobal('fetch', fetcher)
    await expect(getImageTask('key', 'task-1', new AbortController().signal)).rejects.toThrow('Invalid image task')
    // 2. 403 必须原样报出，不自动切换接口或重新扣费。
    fetcher.mockResolvedValue(new Response(JSON.stringify({ error: { message: 'Images are disabled' } }), { status: 403, headers: { 'Content-Type': 'application/json' } }))
    await expect(submitImage('key', { model: 'image', prompt: 'cat', n: 1, platform: 'openai' }, false, new AbortController().signal)).rejects.toThrow('Images are disabled')
    expect(fetcher).toHaveBeenCalledTimes(2)
  })

  it('模型按当前密钥加载，图片下载不向外部 URL 泄露密钥', async () => {
    // 1. 模型建议过滤文本模型，同时保持网关返回的别名字符串。
    const fetcher = vi.fn().mockResolvedValue(new Response(JSON.stringify({ data: [{ id: 'gpt-image-2' }, { id: 'text-model' }] })))
    vi.stubGlobal('fetch', fetcher)
    expect(await listImageModels('key', new AbortController().signal)).toEqual(['gpt-image-2'])
    // 2. 异步下载走受鉴权的任务接口；同步 URL 下载不携带认证头或 Cookie。
    fetcher.mockImplementation(() => Promise.resolve(new Response('image', { headers: { 'Content-Type': 'image/png' } })))
    await getImageBlob({ url: 'https://images.example/result.png' }, 'secret', 'task-1', 0)
    expect(fetcher.mock.calls[1][0]).toMatch(/\/v1\/images\/tasks\/task-1\/images\/0$/)
    expect(fetcher.mock.calls[1][1].headers.Authorization).toBe('Bearer secret')
    await getImageBlob({ url: 'https://images.example/result.png' }, 'secret', undefined, 0)
    expect(fetcher.mock.calls[2][1].headers).toBeUndefined()
    expect(fetcher.mock.calls[2][1].credentials).toBe('omit')
  })
})
