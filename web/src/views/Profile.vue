<template>
    <div class="profile-container">
        <Navigation />

        <div class="main-content">
            <el-card class="profile-card" shadow="hover">
                <template #header>
                    <div class="card-header">
                        <el-icon>
                            <UserFilled />
                        </el-icon>
                        <span>个人信息</span>
                    </div>
                </template>

                <!-- 成功/错误消息 -->
                <el-alert v-if="successMessage" :title="successMessage" type="success" :closable="true"
                    @close="successMessage = ''" show-icon class="message-alert" />

                <el-alert v-if="errorMessage" :title="errorMessage" type="error" :closable="true"
                    @close="errorMessage = ''" show-icon class="message-alert" />

                <!-- 头像设置区域 -->
                <div class="avatar-section">
                    <div class="avatar-wrapper">
                        <UserAvatar :size="120" :avatar="userInfo.avatar" :username="userInfo.username"
                            :nickname="userInfo.nickname" class="profile-avatar" />
                    </div>
                    <div class="avatar-input">
                        <el-input v-model="avatarUrl" placeholder="请输入头像链接" :prefix-icon="Link" clearable
                            class="avatar-url-input" />
                        <p class="avatar-hint">支持 JPG、PNG、GIF 格式的图片链接</p>
                    </div>
                </div>

                <!-- 基本信息表单 -->
                <el-form ref="profileFormRef" :model="profileForm" :rules="profileRules"
                    @submit.prevent="handleUpdateProfile" class="profile-form" label-position="top" size="large">
                    <div class="form-row">
                        <el-form-item label="用户名" prop="username">
                            <el-input v-model="profileForm.username" placeholder="用户名" :prefix-icon="User"
                                :disabled="true" clearable />
                            <template #label>
                                <span>用户名</span>
                                <el-tooltip content="用户名不可修改" placement="top">
                                    <el-icon class="hint-icon">
                                        <QuestionFilled />
                                    </el-icon>
                                </el-tooltip>
                            </template>
                        </el-form-item>

                        <el-form-item label="昵称" prop="nickname">
                            <el-input v-model="profileForm.nickname" placeholder="请输入昵称" :prefix-icon="UserFilled"
                                :disabled="loading" clearable />
                            <template #label>
                                <span>昵称</span>
                                <el-tooltip content="将显示在个人资料中" placement="top">
                                    <el-icon class="hint-icon">
                                        <QuestionFilled />
                                    </el-icon>
                                </el-tooltip>
                            </template>
                        </el-form-item>
                    </div>

                    <div class="form-row">
                        <el-form-item label="角色" prop="role">
                            <el-input v-model="profileForm.role" placeholder="用户角色" :prefix-icon="UserFilled"
                                :disabled="true" clearable />
                            <template #label>
                                <span>角色</span>
                                <el-tooltip content="用户角色不可修改" placement="top">
                                    <el-icon class="hint-icon">
                                        <QuestionFilled />
                                    </el-icon>
                                </el-tooltip>
                            </template>
                        </el-form-item>

                        <el-form-item label="状态" prop="status">
                            <el-tag :type="profileForm.status === 1 ? 'success' : 'danger'" size="large" effect="light">
                                {{ profileForm.status === 1 ? '正常' : '禁用' }}
                            </el-tag>
                            <template #label>
                                <span>状态</span>
                                <el-tooltip content="账户状态" placement="top">
                                    <el-icon class="hint-icon">
                                        <QuestionFilled />
                                    </el-icon>
                                </el-tooltip>
                            </template>
                        </el-form-item>
                    </div>

                    <!-- 密码修改区域 -->
                    <el-divider content-position="left">
                        <span class="divider-title">修改密码（可选）</span>
                    </el-divider>

                    <div class="form-row">
                        <el-form-item label="新密码" prop="password">
                            <el-input v-model="passwordForm.password" type="password" placeholder="留空则不修改密码"
                                :prefix-icon="Lock" :disabled="loading" show-password clearable />
                            <template #label>
                                <span>新密码</span>
                                <el-tooltip content="留空则不修改密码" placement="top">
                                    <el-icon class="hint-icon">
                                        <QuestionFilled />
                                    </el-icon>
                                </el-tooltip>
                            </template>
                        </el-form-item>

                        <el-form-item label="确认新密码" prop="confirm_password">
                            <el-input v-model="passwordForm.confirm_password" type="password" placeholder="请再次输入新密码"
                                :prefix-icon="Lock" :disabled="loading" show-password clearable />
                            <template #label>
                                <span>确认新密码</span>
                                <el-tooltip content="确保两次输入一致" placement="top">
                                    <el-icon class="hint-icon">
                                        <QuestionFilled />
                                    </el-icon>
                                </el-tooltip>
                            </template>
                        </el-form-item>
                    </div>

                    <el-form-item>
                        <el-button type="primary" native-type="submit" :loading="loading" class="update-button"
                            size="large">
                            {{ loading ? '更新中...' : '保存修改' }}
                        </el-button>
                    </el-form-item>
                </el-form>

                <!-- 账户信息 -->
                <el-divider content-position="left">
                    <span class="divider-title">账户信息</span>
                </el-divider>

                <div class="account-info">
                    <div class="info-item">
                        <span class="info-label">注册时间：</span>
                        <span class="info-value">{{ formatDate(userInfo.created_at) }}</span>
                    </div>
                    <div class="info-item">
                        <span class="info-label">最后更新：</span>
                        <span class="info-value">{{ formatDate(userInfo.updated_at) }}</span>
                    </div>
                </div>
            </el-card>
        </div>
    </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { getUserInfo, updateUserInfo } from '@/api/user'
import { User, UserFilled, Lock, QuestionFilled, Link } from '@element-plus/icons-vue'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import Navigation from '@/components/layout/Navigation.vue'
import UserAvatar from '@/components/ui/UserAvatar.vue'

// 响应式数据
const loading = ref(false)
const successMessage = ref('')
const errorMessage = ref('')
const profileFormRef = ref<FormInstance>()
const avatarUrl = ref('')



// 用户信息
const userInfo = reactive({
    id: 0,
    username: '',
    nickname: '',
    avatar: '',
    role: '',
    status: 1,
    created_at: new Date(),
    updated_at: new Date()
})

// 基本信息表单
const profileForm = reactive({
    username: '',
    nickname: '',
    role: '',
    status: 1
})

// 密码表单
const passwordForm = reactive({
    password: '',
    confirm_password: ''
})

// 表单验证规则
const validateConfirmPassword = (rule: any, value: string, callback: any) => {
    // 如果密码为空，则确认密码也可以为空
    if (passwordForm.password === '') {
        callback()
        return
    }

    if (value === '') {
        callback(new Error('请再次输入密码'))
    } else if (value !== passwordForm.password) {
        callback(new Error('两次输入的密码不一致'))
    } else {
        callback()
    }
}

const profileRules: FormRules = {
    nickname: [
        { required: true, message: '请输入昵称', trigger: 'blur' },
        { min: 2, max: 20, message: '昵称长度必须在2-20位之间', trigger: 'blur' }
    ],
    password: [
        { min: 6, message: '密码长度至少6位', trigger: 'blur' }
    ],
    confirm_password: [
        { validator: validateConfirmPassword, trigger: 'blur' }
    ]
}

// 获取用户信息
const fetchUserInfo = async () => {
    try {
        const response = await getUserInfo()
        if (response.data.success) {
            const data = response.data.data
            Object.assign(userInfo, data)
            Object.assign(profileForm, {
                username: data.username,
                nickname: data.nickname,
                role: data.role,
                status: data.status
            })
            // 设置头像链接输入框的初始值
            avatarUrl.value = data.avatar || ''
        }
    } catch (error: any) {
        console.error('获取用户信息失败:', error)
        errorMessage.value = '获取用户信息失败，请稍后重试'
    }
}

// 统一更新所有信息
const handleUpdateProfile = async () => {
    if (!profileFormRef.value) return

    try {
        await profileFormRef.value.validate()
    } catch {
        return
    }

    loading.value = true
    errorMessage.value = ''

    try {
        // 构建更新数据，只包含有值的字段
        const updateData: any = {}

        // 昵称（必填）
        updateData.nickname = profileForm.nickname

        // 头像（如果有输入则更新）
        if (avatarUrl.value.trim()) {
            updateData.avatar = avatarUrl.value.trim()
        }

        // 密码（如果输入了则更新）
        if (passwordForm.password.trim()) {
            updateData.password = passwordForm.password.trim()
            updateData.confirm_password = passwordForm.confirm_password.trim()
        }

        const response = await updateUserInfo(updateData)

        if (response.data.success) {
            successMessage.value = '个人信息更新成功！'
            ElMessage.success('个人信息更新成功！')

            // 清空密码表单
            passwordForm.password = ''
            passwordForm.confirm_password = ''
            profileFormRef.value?.clearValidate()

            // 重新获取用户信息
            await fetchUserInfo()
        } else {
            errorMessage.value = response.data.message || '更新失败'
        }
    } catch (error: any) {
        console.error('更新失败:', error)
        errorMessage.value = error.response?.data?.message || '更新失败，请稍后重试'
    } finally {
        loading.value = false
    }
}



// 格式化日期
const formatDate = (date: Date | string) => {
    if (!date) return '未知'
    const d = new Date(date)
    return d.toLocaleString('zh-CN', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
    })
}

// 组件挂载时获取用户信息
onMounted(() => {
    fetchUserInfo()
})
</script>

<style scoped>
.profile-container {
    min-height: 100vh;
    background: var(--el-bg-color-page);
}

.main-content {
    max-width: 800px;
    margin: 0 auto;
    padding: 40px 20px;
}

.profile-card {
    border-radius: 16px;
}

.card-header {
    display: flex;
    align-items: center;
    gap: 8px;
    font-weight: 600;
    color: var(--el-text-color-primary);
    font-size: 18px;
}

.message-alert {
    margin-bottom: 20px;
}

.avatar-section {
    text-align: center;
    margin-bottom: 30px;
    padding: 20px;
    background: var(--el-bg-color-soft);
    border-radius: 12px;
}

.avatar-wrapper {
    margin-bottom: 20px;
}

.profile-avatar {
    border: 4px solid var(--el-border-color-lighter);
}

.avatar-input {
    max-width: 400px;
    margin: 0 auto;
}

.avatar-url-input {
    margin-bottom: 8px;
}

.avatar-hint {
    margin: 0;
    font-size: 12px;
    color: var(--el-text-color-secondary);
    text-align: center;
}

.profile-form {
    margin-bottom: 30px;
}

.form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 20px;
}

.update-button {
    width: 100%;
    height: 44px;
    font-size: 16px;
    font-weight: 500;
}

.divider-title {
    font-size: 16px;
    font-weight: 600;
    color: var(--el-text-color-primary);
}

.account-info {
    padding: 20px;
    background: var(--el-bg-color-soft);
    border-radius: 12px;
}

.info-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 0;
    border-bottom: 1px solid var(--el-border-color-lighter);
}

.info-item:last-child {
    border-bottom: none;
}

.info-label {
    font-weight: 500;
    color: var(--el-text-color-regular);
}

.info-value {
    color: var(--el-text-color-primary);
    font-family: monospace;
}

.hint-icon {
    margin-left: 4px;
    color: var(--el-text-color-secondary);
    cursor: help;
}

@media (max-width: 768px) {
    .main-content {
        padding: 20px 15px;
    }

    .profile-card {
        max-width: 100%;
    }

    .form-row {
        grid-template-columns: 1fr;
        gap: 20px;
    }

    .avatar-section {
        padding: 15px;
    }

    .profile-avatar {
        width: 100px !important;
        height: 100px !important;
    }
}
</style>
