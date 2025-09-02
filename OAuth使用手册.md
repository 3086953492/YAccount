# YAccount OAuth2.0 使用手册

## 概述

你已经成功将YAccount改造为OAuth2.0授权服务器！现在它可以为其他应用提供统一的用户认证和授权服务。

## OAuth2.0基本概念

### 角色说明
- **授权服务器（Authorization Server）**：你的YAccount系统
- **资源服务器（Resource Server）**：提供用户数据的服务器（也是YAccount）
- **客户端（Client）**：想要访问用户数据的第三方应用
- **资源所有者（Resource Owner）**：用户本人

### 支持的授权模式
1. **授权码模式（Authorization Code）**：最安全，适用于有后端的Web应用
2. **客户端凭证模式（Client Credentials）**：用于应用间调用，无用户参与
3. **刷新令牌（Refresh Token）**：延长访问令牌的有效期

## 使用步骤详解

### 第一步：注册OAuth客户端

在使用OAuth之前，第三方应用需要先在你的系统中注册。

**API端点**：`POST /api/oauth/v1/clients`

**请求示例**：
```bash
POST /api/oauth/v1/clients
Authorization: Bearer {你的登录JWT令牌}
Content-Type: application/json

{
  "name": "我的第三方应用",
  "description": "这是一个测试应用",
  "redirect_uris": [
    "http://localhost:3000/callback",
    "https://myapp.com/auth/callback"
  ],
  "grant_types": [
    "authorization_code",
    "refresh_token",
    "client_credentials"
  ],
  "scopes": [
    "read",
    "profile"
  ],
  "client_type": "confidential"
}
```

**响应示例**：
```json
{
  "code": 200,
  "message": "OAuth客户端创建成功",
  "data": {
    "id": 1,
    "client_id": "client_AbCdEfGhIjKlMnOpQr",
    "client_secret": "abc123...xyz789",
    "name": "我的第三方应用",
    "redirect_uris": "[\"http://localhost:3000/callback\"]",
    "grant_types": "authorization_code,refresh_token,client_credentials",
    "scopes": "read,profile",
    "client_type": "confidential"
  }
}
```

**⚠️ 重要提示**：`client_secret` 只会在创建时返回一次，请务必保存好！

### 第二步：授权码模式完整流程

这是最常用的OAuth流程，适用于有后端的Web应用。

#### 2.1 引导用户授权

**步骤1**：构建授权URL
```
GET /oauth/authorize?response_type=code&client_id={CLIENT_ID}&redirect_uri={REDIRECT_URI}&scope=read profile&state={RANDOM_STATE}
```

**参数说明**：
- `response_type`: 固定值 `code`
- `client_id`: 你的客户端ID
- `redirect_uri`: 授权后的回调地址（必须在注册时声明过）
- `scope`: 请求的权限范围，空格分隔
- `state`: 随机字符串，防止CSRF攻击

**示例URL**：
```
http://localhost:8081/oauth/authorize?response_type=code&client_id=client_AbCdEfGhIjKlMnOpQr&redirect_uri=http://localhost:3000/callback&scope=read profile&state=xyz123
```

#### 2.2 用户授权流程

1. **用户访问授权URL**
2. **如果用户未登录**：
   - 后端返回401错误
   - 前端跳转到登录页面
   - 用户登录后重新访问授权URL

3. **如果用户已登录**：
   - 后端返回授权成功响应：
   ```json
   {
     "code": 200,
     "message": "授权成功",
     "data": {
       "redirect_uri": "http://localhost:3000/callback?code=auth_code_here&state=xyz123",
       "state": "xyz123"
     }
   }
   ```
   - 前端根据`redirect_uri`重定向到客户端应用

#### 2.3 交换访问令牌

客户端应用的后端收到授权码后，交换访问令牌：

**API端点**：`POST /oauth/token`

**请求示例**：
```bash
POST /oauth/token
Content-Type: application/x-www-form-urlencoded

grant_type=authorization_code&
client_id=client_AbCdEfGhIjKlMnOpQr&
client_secret=abc123...xyz789&
code=auth_code_here&
redirect_uri=http://localhost:3000/callback
```

**响应示例**：
```json
{
  "code": 200,
  "message": "令牌获取成功",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "scope": "read profile"
  }
}
```

### 第三步：使用访问令牌

#### 3.1 访问用户信息

**API端点**：`GET /oauth/userinfo`

**请求示例**：
```bash
GET /oauth/userinfo
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例**：
```json
{
  "code": 200,
  "message": "获取用户信息成功",
  "data": {
    "sub": 1,
    "username": "testuser",
    "nickname": "测试用户",
    "avatar": "avatar_url",
    "role": "user"
  }
}
```

#### 3.2 访问其他受保护的API

你可以使用同样的Bearer令牌访问其他需要认证的API：

```bash
GET /api/account/v1/user/profile
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### 第四步：刷新访问令牌

当访问令牌过期时，使用刷新令牌获取新的访问令牌：

**API端点**：`POST /oauth/token`

**请求示例**：
```bash
POST /oauth/token
Content-Type: application/x-www-form-urlencoded

grant_type=refresh_token&
client_id=client_AbCdEfGhIjKlMnOpQr&
client_secret=abc123...xyz789&
refresh_token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例**：
```json
{
  "code": 200,
  "message": "令牌获取成功",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "scope": "read profile"
  }
}
```

### 第五步：客户端凭证模式

这种模式用于应用间调用，不涉及用户：

**API端点**：`POST /oauth/token`

**请求示例**：
```bash
POST /oauth/token
Content-Type: application/x-www-form-urlencoded

grant_type=client_credentials&
client_id=client_AbCdEfGhIjKlMnOpQr&
client_secret=abc123...xyz789&
scope=read
```

**响应示例**：
```json
{
  "code": 200,
  "message": "令牌获取成功",
  "data": {
    "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "token_type": "Bearer",
    "expires_in": 3600,
    "scope": "read"
  }
}
```

## 权限范围（Scope）说明

你的系统支持以下权限范围：

- **read**：读取基本用户信息（默认）
- **profile**：访问完整用户资料（包括头像、角色等）
- **write**：修改用户信息
- **admin**：管理员权限

## 客户端管理API

### 获取客户端列表（管理员）
```bash
GET /api/oauth/v1/clients
Authorization: Bearer {管理员JWT令牌}
```

### 获取客户端详情
```bash
GET /api/oauth/v1/clients/{client_id}
Authorization: Bearer {JWT令牌}
```

### 更新客户端
```bash
PUT /api/oauth/v1/clients/{client_id}
Authorization: Bearer {JWT令牌}
Content-Type: application/json

{
  "name": "更新后的应用名称",
  "description": "更新后的描述"
}
```

### 删除客户端
```bash
DELETE /api/oauth/v1/clients/{client_id}
Authorization: Bearer {JWT令牌}
```

## 令牌内省

检查令牌是否有效：

**API端点**：`POST /oauth/introspect`

**请求示例**：
```bash
POST /oauth/introspect
Content-Type: application/x-www-form-urlencoded

token=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

**响应示例**：
```json
{
  "active": true,
  "sub": 1,
  "client_id": "client_AbCdEfGhIjKlMnOpQr",
  "scope": "read profile",
  "exp": 1699999999,
  "iat": 1699996399,
  "token_type": "access_token"
}
```

## 实际应用场景

### 场景1：第三方登录
```javascript
// 前端JavaScript示例
function loginWithYAccount() {
  const authUrl = 'http://localhost:8081/oauth/authorize?' +
    'response_type=code&' +
    'client_id=client_AbCdEfGhIjKlMnOpQr&' +
    'redirect_uri=http://localhost:3000/callback&' +
    'scope=read profile&' +
    'state=' + generateRandomState();
  
  window.location.href = authUrl;
}

// 回调页面处理
function handleCallback() {
  const urlParams = new URLSearchParams(window.location.search);
  const code = urlParams.get('code');
  const state = urlParams.get('state');
  
  // 将授权码发送到你的后端服务器
  fetch('/api/auth/callback', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ code, state })
  });
}
```

### 场景2：API调用
```javascript
// 使用访问令牌调用API
async function getUserInfo(accessToken) {
  const response = await fetch('http://localhost:8081/oauth/userinfo', {
    headers: {
      'Authorization': `Bearer ${accessToken}`
    }
  });
  
  return await response.json();
}
```

### 场景3：服务间调用
```javascript
// 客户端凭证模式获取令牌
async function getServiceToken() {
  const response = await fetch('http://localhost:8081/oauth/token', {
    method: 'POST',
    headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
    body: new URLSearchParams({
      grant_type: 'client_credentials',
      client_id: 'your_client_id',
      client_secret: 'your_client_secret',
      scope: 'read'
    })
  });
  
  const data = await response.json();
  return data.data.access_token;
}
```

## 错误处理

### 常见错误码

- `INVALID_CLIENT`：客户端认证失败
- `INVALID_GRANT`：授权码无效或已过期
- `INVALID_SCOPE`：请求的权限范围无效
- `INVALID_TOKEN`：访问令牌无效
- `ACCESS_DENIED`：用户拒绝授权
- `UNSUPPORTED_GRANT_TYPE`：不支持的授权类型

### 错误响应示例
```json
{
  "code": 400,
  "message": "授权码无效或已过期",
  "data": null
}
```

## 安全建议

1. **保护客户端密钥**：永远不要在前端代码中暴露`client_secret`
2. **使用HTTPS**：生产环境必须使用HTTPS
3. **验证state参数**：防止CSRF攻击
4. **定期轮换密钥**：定期更新客户端密钥
5. **限制重定向URI**：只允许可信的回调地址
6. **令牌有效期**：设置合理的令牌过期时间

## 配置说明

在你的配置文件中，相关的OAuth配置：

```json
{
  "oauth": {
    "issuer": "http://localhost:8081",
    "authorize_ui": "http://localhost:5173/oauth/authorize",
    "authorization_code_ttl": "10m",
    "access_token_ttl": "1h",
    "refresh_token_ttl": "168h",
    "supported_grant_types": [
      "authorization_code",
      "refresh_token",
      "client_credentials"
    ],
    "default_scopes": ["read"]
  }
}
```

## 总结

现在你的YAccount系统已经是一个完整的OAuth2.0授权服务器了！其他应用可以通过标准的OAuth2.0流程接入你的认证系统，实现统一登录和用户信息共享。

这大大简化了多服务架构下的用户认证复杂度，新的服务只需要实现OAuth2.0客户端逻辑即可快速接入。
