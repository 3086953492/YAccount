<template>
  <el-dialog
    v-model="dialogVisible"
    title="客户端密钥"
    width="600px"
    :close-on-click-modal="false"
    :close-on-press-escape="false"
    :show-close="false"
    center
  >
    <div class="secret-display-container">
      <!-- 警告提示 -->
      <el-alert
        type="warning"
        show-icon
        :closable="false"
        class="security-warning"
      >
        <template #title>
          <div class="warning-title">
            <el-icon><Warning /></el-icon>
            <span>重要安全提示</span>
          </div>
        </template>
        <div class="warning-content">
          <p>• 客户端密钥只会显示一次，请立即保存到安全的地方</p>
          <p>• 密钥泄露可能导致安全风险，请妥善保管</p>
          <p>• 如果丢失密钥，您需要重新生成客户端</p>
        </div>
      </el-alert>

      <!-- 客户端信息 -->
      <div class="client-info">
        <h3 class="section-title">
          <el-icon><Key /></el-icon>
          客户端信息
        </h3>
        
        <div class="info-grid">
          <div class="info-item">
            <label class="info-label">客户端名称</label>
            <div class="info-value">{{ clientData?.name }}</div>
          </div>
          
          <div class="info-item">
            <label class="info-label">Client ID</label>
            <div class="info-value client-id">
              <code>{{ clientData?.client_id }}</code>
              <el-button
                type="primary"
                link
                size="small"
                @click="copyToClipboard(clientData?.client_id || '')"
                class="copy-btn"
              >
                <el-icon><CopyDocument /></el-icon>
              </el-button>
            </div>
          </div>
          
          <div class="info-item secret-item">
            <label class="info-label">Client Secret</label>
            <div class="info-value client-secret">
              <div class="secret-container">
                <code class="secret-code">{{ clientData?.client_secret }}</code>
                <div class="secret-actions">
                  <el-button
                    type="primary"
                    @click="copyToClipboard(clientData?.client_secret || '')"
                    class="copy-secret-btn"
                  >
                    <el-icon><CopyDocument /></el-icon>
                    复制密钥
                  </el-button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- 使用说明 -->
      <div class="usage-guide">
        <h3 class="section-title">
          <el-icon><Document /></el-icon>
          使用说明
        </h3>
        <div class="guide-content">
          <ol>
            <li>将 Client ID 和 Client Secret 保存到您的应用配置中</li>
            <li>在OAuth授权流程中使用这些凭证进行身份验证</li>
            <li>确保在生产环境中通过环境变量或安全配置文件管理密钥</li>
            <li>定期轮换客户端密钥以确保安全</li>
          </ol>
        </div>
      </div>
    </div>

    <template #footer>
      <div class="dialog-footer">
        <el-checkbox v-model="confirmSaved" class="confirm-checkbox">
          我已安全保存了客户端密钥
        </el-checkbox>
        <el-button
          type="primary"
          @click="handleConfirm"
          :disabled="!confirmSaved"
        >
          我已保存，继续
        </el-button>
      </div>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import { Warning, Key, CopyDocument, Document } from '@element-plus/icons-vue'
import type { OAuthClient } from '@/api/oauth'

interface Props {
  visible: boolean
  clientData?: OAuthClient | null
}

interface Emits {
  (e: 'update:visible', value: boolean): void
  (e: 'confirmed'): void
}

const props = withDefaults(defineProps<Props>(), {
  visible: false,
  clientData: null
})

const emit = defineEmits<Emits>()

// 对话框显示状态
const dialogVisible = computed({
  get: () => props.visible,
  set: (value) => emit('update:visible', value)
})

// 确认已保存状态
const confirmSaved = ref(false)

// 复制到剪贴板
const copyToClipboard = async (text: string) => {
  try {
    await navigator.clipboard.writeText(text)
    ElMessage.success('已复制到剪贴板')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

// 确认处理
const handleConfirm = () => {
  if (!confirmSaved.value) {
    ElMessage.warning('请确认您已安全保存了客户端密钥')
    return
  }
  
  emit('confirmed')
  dialogVisible.value = false
  confirmSaved.value = false
}
</script>

<style scoped>
.secret-display-container {
  max-height: 70vh;
  overflow-y: auto;
}

.security-warning {
  margin-bottom: 24px;
  border-radius: 8px;
}

.warning-title {
  display: flex;
  align-items: center;
  gap: 6px;
  font-weight: 600;
  font-size: 16px;
}

.warning-content {
  margin-top: 8px;
  font-size: 14px;
  line-height: 1.6;
}

.warning-content p {
  margin: 4px 0;
}

.client-info {
  margin-bottom: 24px;
}

.section-title {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 16px;
  font-weight: 600;
  color: var(--el-text-color-primary);
  margin-bottom: 16px;
  padding-bottom: 8px;
  border-bottom: 2px solid var(--el-color-primary);
}

.info-grid {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.info-item {
  padding: 16px;
  background: var(--el-fill-color-extra-light);
  border-radius: 8px;
  border: 1px solid var(--el-border-color-lighter);
}

.secret-item {
  border: 2px solid var(--el-color-warning);
  background: var(--el-color-warning-light-9);
}

.info-label {
  display: block;
  font-weight: 600;
  color: var(--el-text-color-regular);
  margin-bottom: 8px;
  font-size: 14px;
}

.info-value {
  font-size: 14px;
  color: var(--el-text-color-primary);
}

.client-id {
  display: flex;
  align-items: center;
  gap: 8px;
}

.client-id code {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  background: var(--el-fill-color-lighter);
  padding: 6px 12px;
  border-radius: 4px;
  font-size: 13px;
  color: var(--el-color-primary);
  flex: 1;
}

.copy-btn {
  padding: 4px;
  min-height: auto;
}

.client-secret {
  width: 100%;
}

.secret-container {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.secret-code {
  font-family: 'Monaco', 'Menlo', 'Ubuntu Mono', monospace;
  background: var(--el-color-warning-light-8);
  padding: 12px;
  border-radius: 6px;
  font-size: 13px;
  color: var(--el-text-color-primary);
  word-break: break-all;
  line-height: 1.5;
  border: 1px solid var(--el-color-warning-light-5);
  display: block;
  width: 100%;
  box-sizing: border-box;
}

.secret-actions {
  display: flex;
  justify-content: center;
}

.copy-secret-btn {
  background: var(--el-color-warning);
  border-color: var(--el-color-warning);
  color: white;
  font-weight: 500;
  padding: 8px 16px;
}

.copy-secret-btn:hover {
  background: var(--el-color-warning-dark-2);
  border-color: var(--el-color-warning-dark-2);
}

.usage-guide {
  margin-bottom: 16px;
}

.guide-content {
  background: var(--el-fill-color-extra-light);
  padding: 16px;
  border-radius: 8px;
  border-left: 4px solid var(--el-color-primary);
}

.guide-content ol {
  margin: 0;
  padding-left: 20px;
  line-height: 1.8;
}

.guide-content li {
  margin: 8px 0;
  color: var(--el-text-color-regular);
  font-size: 14px;
}

.dialog-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
}

.confirm-checkbox {
  font-size: 14px;
  color: var(--el-text-color-regular);
}

/* 响应式设计 */
@media (max-width: 768px) {
  .info-grid {
    gap: 12px;
  }
  
  .info-item {
    padding: 12px;
  }
  
  .client-id {
    flex-direction: column;
    align-items: stretch;
    gap: 8px;
  }
  
  .secret-container {
    gap: 8px;
  }
  
  .dialog-footer {
    flex-direction: column;
    align-items: stretch;
  }
}

/* Element Plus 样式覆盖 */
:deep(.el-dialog__header) {
  background: linear-gradient(135deg, var(--el-color-primary-light-9), var(--el-color-primary-light-8));
  border-bottom: 1px solid var(--el-border-color-lighter);
}

:deep(.el-dialog__title) {
  font-weight: 600;
  color: var(--el-color-primary);
}

:deep(.el-alert__content) {
  width: 100%;
}

:deep(.el-checkbox__label) {
  color: var(--el-text-color-regular);
}
</style>
