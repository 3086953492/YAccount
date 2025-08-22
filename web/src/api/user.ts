import api from './auth'
import { useAuthStore } from '@/stores/auth'

// 获取用户信息接口
export const getUserInfo = () => {
    const authStore = useAuthStore()
    const userId = authStore.user?.id
    if (!userId) {
        throw new Error('用户未登录')
    }
    return api.get(`/users/${userId}`)
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
