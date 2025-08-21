<template>
  <div class="form-group">
    <label :for="id">{{ label }}</label>
    <div class="input-wrapper">
      <input :id="id" v-model="inputValue" :type="type" :placeholder="placeholder" :required="required"
        :disabled="disabled" @input="handleInput" @blur="handleBlur" @focus="handleFocus" 
        :class="{ 'error': hasError, 'focused': isFocused }" />
      <div class="input-border"></div>
    </div>
    <span v-if="errorMessage" class="error-message">{{ errorMessage }}</span>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'

interface Props {
  modelValue: string
  id: string
  label: string
  type?: string
  placeholder?: string
  required?: boolean
  disabled?: boolean
  errorMessage?: string
}

const props = withDefaults(defineProps<Props>(), {
  type: 'text',
  placeholder: '',
  required: false,
  disabled: false,
  errorMessage: ''
})

const emit = defineEmits<{
  'update:modelValue': [value: string]
  'input': [value: string]
  'blur': [value: string]
  'focus': [value: string]
}>()

const inputValue = computed({
  get: () => props.modelValue,
  set: (value: string) => emit('update:modelValue', value)
})

const hasError = computed(() => !!props.errorMessage)
const isFocused = ref(false)

const handleInput = (event: Event) => {
  const target = event.target as HTMLInputElement
  emit('input', target.value)
}

const handleBlur = (event: Event) => {
  const target = event.target as HTMLInputElement
  isFocused.value = false
  emit('blur', target.value)
}

const handleFocus = (event: Event) => {
  const target = event.target as HTMLInputElement
  isFocused.value = true
  emit('focus', target.value)
}
</script>

<style scoped>
.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-group label {
  color: #374151;
  font-weight: 600;
  font-size: 14px;
  margin-left: 2px;
}

.input-wrapper {
  position: relative;
}

.input-wrapper input {
  width: 100%;
  padding: 16px 20px;
  border: 2px solid transparent;
  border-radius: 12px;
  font-size: 16px;
  background: #f9fafb;
  color: #111827;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-sizing: border-box;
}

.input-wrapper input::placeholder {
  color: #9ca3af;
  font-weight: 400;
}

.input-wrapper input:focus {
  outline: none;
  background: #ffffff;
  box-shadow: 
    0 0 0 4px rgba(99, 102, 241, 0.1),
    0 4px 6px -1px rgba(0, 0, 0, 0.1),
    0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

.input-wrapper input:hover:not(:focus):not(:disabled) {
  background: #f3f4f6;
}

.input-wrapper input:disabled {
  background-color: #f3f4f6;
  color: #9ca3af;
  cursor: not-allowed;
}

.input-wrapper input.error {
  background: #fef2f2;
  border-color: #ef4444;
}

.input-wrapper input.error:focus {
  box-shadow: 
    0 0 0 4px rgba(239, 68, 68, 0.1),
    0 4px 6px -1px rgba(0, 0, 0, 0.1),
    0 2px 4px -1px rgba(0, 0, 0, 0.06);
}

.input-wrapper input.focused {
  background: #ffffff;
}

.input-border {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  border-radius: 12px;
  pointer-events: none;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.error-message {
  color: #dc2626;
  font-size: 13px;
  margin-top: 6px;
  margin-left: 2px;
  font-weight: 500;
}

@media (max-width: 640px) {
  .input-wrapper input {
    padding: 14px 16px;
    font-size: 16px;
  }
}
</style>
