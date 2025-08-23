import api from './auth'
import { useAuthStore } from '@/stores/auth'

// 根据用户ID获取用户信息接口
export const getUserInfoById = (userId: number) => {
    return api.get(`/users/${userId}`)
}

// 获取当前登录用户信息接口
export const getCurrentUserInfo = () => {
    const authStore = useAuthStore()
    const userId = authStore.user?.id
    if (!userId) {
        throw new Error('用户未登录')
    }
    return getUserInfoById(userId)
}

// 更新用户信息接口
export const updateUserInfo = (data: {
    nickname?: string
    avatar?: string
    password?: string
    confirm_password?: string
}) => {
    const authStore = useAuthStore()
    const userId = authStore.user?.id
    if (!userId) {
        throw new Error('用户未登录')
    }

    // 转换字段名以匹配后端API
    const requestData: any = {}
    if (data.nickname !== undefined) requestData.nickname = data.nickname
    if (data.avatar !== undefined) requestData.avatar = data.avatar
    if (data.password !== undefined) requestData.password = data.password
    if (data.confirm_password !== undefined) requestData.ConfirmPassword = data.confirm_password

    return api.put(`/users/${userId}`, requestData)
}

// 根据用户ID更新用户信息接口（管理员权限）
export const updateUserInfoById = (userId: number, data: {
    nickname?: string
    avatar?: string
    password?: string
    confirm_password?: string
}) => {
    // 转换字段名以匹配后端API
    const requestData: any = {}
    if (data.nickname !== undefined) requestData.nickname = data.nickname
    if (data.avatar !== undefined) requestData.avatar = data.avatar
    if (data.password !== undefined) requestData.password = data.password
    if (data.confirm_password !== undefined) requestData.ConfirmPassword = data.confirm_password

    return api.put(`/users/${userId}`, requestData)
}

// 获取用户列表接口（管理员权限）
export const getUserList = (params: {
    page?: number
    pageSize?: number
    username?: string
    nickname?: string
}) => {
    return api.get('/users', { params })
}
