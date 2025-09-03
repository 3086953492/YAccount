import api from './auth'

export interface OAuthClient {
  id: number
  client_id: string
  client_secret?: string // 只在创建时返回
  name: string
  description: string
  redirect_uris: string
  grant_types: string
  scopes: string
  client_type: string
  status: string
  access_token_ttl: number
  refresh_token_ttl: number
  owner_id: number
  owner_type: string
  created_at: string
  updated_at: string
}

export interface ListOAuthClientsResponse {
  data: OAuthClient[]
  pagination: {
    page: number
    page_size: number
    total: number
    total_pages: number
  }
}

export interface CreateOAuthClientRequest {
  name: string
  description?: string
  redirect_uris: string[]
  grant_types: string[]
  scopes: string[]
  client_type: 'public' | 'confidential'
}

export interface UpdateOAuthClientRequest {
  name: string
  description?: string
  redirect_uris: string[]
  grant_types: string[]
  scopes: string[]
  client_type: 'public' | 'confidential'
}

// 获取OAuth客户端列表
export const listOAuthClients = (params?: {
  page?: number
  page_size?: number
  name?: string
}) => {
  return api.get('/oauth/clients', { params })
}

// 获取OAuth客户端详情
export const getOAuthClient = (clientId: string) => {
  return api.get(`/oauth/clients/${clientId}`)
}

// 创建OAuth客户端
export const createOAuthClient = (data: CreateOAuthClientRequest) => {
  return api.post('/oauth/clients', data)
}

// 更新OAuth客户端
export const updateOAuthClient = (clientId: string, data: UpdateOAuthClientRequest) => {
  return api.put(`/oauth/clients/${clientId}`, data)
}

// 删除OAuth客户端
export const deleteOAuthClient = (clientId: string) => {
  return api.delete(`/oauth/clients/${clientId}`)
}
