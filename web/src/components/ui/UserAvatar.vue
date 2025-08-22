<template>
    <el-avatar 
        :size="size" 
        :src="avatarUrl" 
        :class="avatarClass"
        @error="handleAvatarError"
        :style="fallbackStyle">
        {{ displayText }}
    </el-avatar>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

interface Props {
    size?: number | string
    avatar?: string
    username?: string
    nickname?: string
    className?: string
}

const props = withDefaults(defineProps<Props>(), {
    size: 40,
    avatar: '',
    username: '',
    nickname: '',
    className: ''
})

const hasError = ref(false)

// 计算头像URL，如果加载失败则返回空字符串
const avatarUrl = computed(() => {
    if (hasError.value || !props.avatar) {
        return ''
    }
    return props.avatar
})

// 计算显示文本（昵称或用户名的首字母）
const displayText = computed(() => {
    if (props.nickname) {
        return props.nickname.charAt(0).toUpperCase()
    }
    if (props.username) {
        return props.username.charAt(0).toUpperCase()
    }
    return '?'
})

// 计算头像样式类
const avatarClass = computed(() => {
    return `user-avatar ${props.className}`.trim()
})

// 计算降级样式（渐变背景）
const fallbackStyle = computed(() => {
    if (hasError.value || !props.avatar) {
        return {
            background: 'linear-gradient(135deg, #667eea 0%, #764ba2 100%)',
            color: 'white',
            fontWeight: '600'
        }
    }
    return {}
})

// 处理头像加载失败
const handleAvatarError = () => {
    hasError.value = true
}
</script>

<style scoped>
.user-avatar {
    transition: all 0.3s ease;
}

.user-avatar:hover {
    transform: scale(1.05);
}
</style>
