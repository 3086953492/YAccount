# 可复用组件库

本项目提取了登录和注册页面的共同元素，创建了以下可复用组件：

## 组件列表

### 1. FormCard (表单卡片)
用于创建统一的表单容器，包含标题、副标题和内容区域。

**Props:**
- `title: string` - 表单标题
- `subtitle?: string` - 表单副标题（可选）

**使用示例:**
```vue
<FormCard title="用户登录" subtitle="欢迎使用YAccount系统">
  <!-- 表单内容 -->
</FormCard>
```

### 2. FormInput (表单输入框)
统一的输入框组件，支持标签、验证状态和错误提示。

**Props:**
- `modelValue: string` - 输入值（支持v-model）
- `id: string` - 输入框ID
- `label: string` - 标签文本
- `type?: string` - 输入类型（默认：text）
- `placeholder?: string` - 占位符文本
- `required?: boolean` - 是否必填
- `disabled?: boolean` - 是否禁用
- `errorMessage?: string` - 错误提示信息

**Events:**
- `update:modelValue` - 值更新事件
- `input` - 输入事件
- `blur` - 失焦事件

**使用示例:**
```vue
<FormInput
  id="username"
  v-model="form.username"
  label="用户名"
  placeholder="请输入用户名"
  required
  :error-message="errors.username"
/>
```

### 3. FormButton (表单按钮)
统一的按钮组件，支持加载状态和多种样式变体。

**Props:**
- `text: string` - 按钮文本
- `loadingText?: string` - 加载时的文本
- `type?: 'button' | 'submit' | 'reset'` - 按钮类型
- `variant?: 'primary' | 'secondary'` - 按钮样式变体
- `loading?: boolean` - 是否显示加载状态
- `disabled?: boolean` - 是否禁用

**Events:**
- `click` - 点击事件

**使用示例:**
```vue
<FormButton
  type="submit"
  text="登录"
  :loading="loading"
  loading-text="登录中..."
/>
```

### 4. FormFooter (表单页脚)
用于显示表单底部的链接和提示信息。

**Props:**
- `text: string` - 提示文本
- `linkText: string` - 链接文本
- `linkTo: string` - 链接地址

**使用示例:**
```vue
<FormFooter
  text="还没有账号？"
  link-text="立即注册"
  link-to="/register"
/>
```

### 5. FormValidator (表单验证器)
提供表单验证功能的容器组件，支持多种验证规则。

**验证规则类型:**
- `required: boolean` - 必填验证
- `minLength: number` - 最小长度验证
- `maxLength: number` - 最大长度验证
- `pattern: RegExp` - 正则表达式验证
- `custom: function` - 自定义验证函数

**使用示例:**
```vue
<FormValidator ref="validator">
  <FormInput
    v-model="form.username"
    :rules="[
      { required: true },
      { minLength: 3 },
      { maxLength: 15 }
    ]"
  />
</FormValidator>
```

## 样式特点

所有组件都采用统一的设计风格：
- 渐变背景和阴影效果
- 响应式设计
- 平滑的过渡动画
- 一致的间距和字体大小
- 支持错误状态的视觉反馈

## 使用建议

1. **保持一致性**: 在项目中统一使用这些组件，确保UI的一致性
2. **合理组合**: 根据具体需求组合使用不同的组件
3. **扩展性**: 组件设计支持通过props和slots进行定制
4. **可访问性**: 组件内置了适当的ARIA属性和键盘导航支持

## 注意事项

- 确保在使用组件前正确导入
- 注意组件的依赖关系（如FormButton依赖LoadingSpinner）
- 可以根据项目需求进一步定制组件的样式和行为
