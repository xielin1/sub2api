/**
 * 首页公开信息 API（匿名可访问）
 * 1. 最近公告：面向所有用户的生效公告，用作「最近更新」
 * 2. 累计服务数据：累计请求数与 Token 数
 */

import { apiClient } from './client'

export interface PublicAnnouncement {
  id: number
  title: string
  content: string
  created_at: string
}

export interface PublicStats {
  total_requests: number
  total_tokens: number
}

export async function getPublicAnnouncements(): Promise<PublicAnnouncement[]> {
  const { data } = await apiClient.get<PublicAnnouncement[]>('/public/announcements')
  return data
}

export async function getPublicStats(): Promise<PublicStats> {
  const { data } = await apiClient.get<PublicStats>('/public/stats')
  return data
}
