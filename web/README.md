# YAccount Web 前端

这是一个基于 Vue 3 + TypeScript + Vite 构建的现代化前端应用，实现了完整的用户认证系统。

## 功能特性

- 🔐 用户登录/注册
- 🛡️ JWT Token 认证
- 📱 响应式设计
- 🎨 现代化 UI 界面
- 🔄 自动 Token 刷新
- 🚫 路由守卫保护

## 技术栈

- **前端框架**: Vue 3
- **构建工具**: Vite
- **语言**: TypeScript
- **状态管理**: Pinia
- **路由**: Vue Router 4
- **HTTP 客户端**: Axios
- **样式**: CSS3 + 响应式设计

## 项目结构

```
src/
├── api/          # API 接口
├── components/   # 公共组件
├── config/       # 配置文件
├── router/       # 路由配置
├── stores/       # 状态管理
├── views/        # 页面组件
└── utils/        # 工具函数
```

## 快速开始

### 安装依赖

```bash
npm install
```

### 开发环境运行

```bash
npm run dev
```

### 生产环境构建

```bash
npm run build
```

## 使用说明

### 1. 登录功能

访问 `/login` 页面，输入用户名和密码进行登录。登录成功后会自动跳转到首页。

### 2. 注册功能

访问 `/register` 页面，填写完整的注册信息创建新账户。

### 3. 用户管理

登录后可以在首页查看个人信息，包括：
- 用户基本信息
- 账户状态
- 注册时间
- 最后更新时间

### 4. 安全特性

- JWT Token 自动管理
- 路由守卫保护
- 自动 Token 刷新
- 401 错误自动跳转登录

## 配置说明

### API 配置

在 `src/config/index.ts` 中可以修改 API 基础 URL：

```typescript
export const config = {
  apiBaseURL: 'http://localhost:8080/api/account/v1',
  // ... 其他配置
}
```

### 环境变量

可以通过环境变量 `VITE_API_BASE_URL` 来配置 API 地址。

## 开发指南

### 添加新页面

1. 在 `src/views/` 目录下创建新的 Vue 组件
2. 在 `src/router/index.ts` 中添加路由配置
3. 设置适当的 `meta.requiresAuth` 值

### 添加新 API

1. 在 `src/api/` 目录下创建或修改相应的 API 文件
2. 使用统一的 axios 实例和错误处理

### 状态管理

使用 Pinia 进行状态管理，认证相关的状态在 `src/stores/auth.ts` 中管理。

## 注意事项

- 确保后端服务正在运行
- 默认 API 地址为 `http://localhost:8080`
- 支持自动 Token 刷新
- 响应式设计，支持移动端访问

## 浏览器支持

- Chrome >= 87
- Firefox >= 78
- Safari >= 14
- Edge >= 88
