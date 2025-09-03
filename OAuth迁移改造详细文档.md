# YAccount系统OAuth完全迁移改造文档

## 概述

基于对现有代码的深入分析，您的系统完全可以从JWT+OAuth混合架构迁移到纯OAuth架构。当前系统已经具备完整的OAuth 2.0实现，包括：

- ✅ 授权码模式 (Authorization Code)
- ✅ 刷新令牌模式 (Refresh Token)  
- ✅ 客户端凭证模式 (Client Credentials)
- ✅ 完整的Token管理机制
- ✅ 权限范围(Scopes)控制

## 改造优势

### 1. 标准化优势
- 使用标准OAuth 2.0协议，更好的第三方集成
- 统一的认证授权体系，降低维护复杂度
- 更好的安全性和可扩展性

### 2. 架构简化
- 消除JWT和OAuth的重复逻辑
- 统一的Token管理机制
- 减少配置复杂度

### 3. 功能增强
- 更细粒度的权限控制(Scopes)
- 更好的Token撤销机制
- 支持多种授权模式

## 改造计划

## 第一阶段：OAuth登录模式实现

### 1.1 创建内部OAuth客户端

**目标**：为Web应用创建专用的OAuth客户端来替代JWT登录

**步骤**：

1. **创建系统客户端注册脚本**

```go
// scripts/create_system_client.go
package main

import (
    "YAccount/models"
    "YAccount/services"
    "log"
)

func main() {
    // 创建系统内部客户端
    req := &models.CreateOAuthClientRequest{
        Name:         "YAccount Web App",
        ClientType:   "confidential", // 机密客户端
        RedirectURIs: []string{
            "http://localhost:3000/auth/callback",
            "https://yourdomain.com/auth/callback",
        },
        Scopes: "read,write,admin", // 根据需要调整
        GrantTypes: []string{
            "authorization_code",
            "refresh_token",
        },
    }
    
    client, err := services.CreateOAuthClient(req, 1) // 假设系统用户ID为1
    if err != nil {
        log.Fatal("创建系统客户端失败:", err)
    }
    
    log.Printf("系统客户端创建成功:")
    log.Printf("Client ID: %s", client.ClientID)
    log.Printf("Client Secret: %s", client.ClientSecret)
}
```

2. **执行客户端创建**

```bash
cd server
go run scripts/create_system_client.go
```

### 1.2 实现OAuth登录流程

**目标**：使用OAuth授权码模式替代JWT登录

**步骤**：

1. **修改登录服务** - 文件：`server/services/user.go`

```go
// 添加OAuth登录方法
func OAuthLoginService(req *models.LoginRequest, clientID string) (*models.UserResponse, map[string]any, error) {
    // 验证用户凭据
    user, err := repositories.Login(req)
    if err != nil {
        if !apperrors.IsNotFoundError(err) {
            logger.LogError("OAuthLoginService", "database query", "从数据库中获取用户失败", err, zap.String("username", req.Username))
        }
        return nil, nil, apperrors.ErrUsernameOrPasswordError
    }

    // 获取客户端信息
    client, err := services.GetOAuthClientByID(clientID)
    if err != nil {
        return nil, nil, err
    }

    // 根据用户角色确定授权范围
    scopes := []string{"read"}
    if user.Role == "admin" {
        scopes = append(scopes, "write", "admin")
    }

    // 生成访问令牌和刷新令牌
    accessToken, err := oauth.GenerateAccessToken(user.ID, clientID, scopes)
    if err != nil {
        return nil, nil, err
    }

    refreshToken, err := oauth.GenerateRefreshToken(user.ID, clientID, scopes)
    if err != nil {
        return nil, nil, err
    }

    // 保存令牌记录
    tokenRecord := &models.OAuthAccessToken{
        AccessToken:      accessToken,
        RefreshToken:     refreshToken,
        ClientID:         clientID,
        UserID:           user.ID,
        Scopes:           strings.Join(scopes, " "),
        ExpiresAt:        time.Now().Add(global.Cfg.OAuth.AccessTokenTTL),
        RefreshExpiresAt: time.Now().Add(global.Cfg.OAuth.RefreshTokenTTL),
    }

    if err := repositories.CreateOAuthAccessToken(tokenRecord); err != nil {
        logger.LogError("OAuthLoginService", "database", "保存令牌失败", err)
        return nil, nil, apperrors.ErrServerInternal
    }

    logger.Info("用户OAuth登录成功", zap.String("username", user.Username))

    userResponse := &models.UserResponse{
        ID:        user.ID,
        Username:  user.Username,
        Role:      user.Role,
        Nickname:  user.Nickname,
        Avatar:    user.Avatar,
        Status:    user.Status,
        CreatedAt: user.CreatedAt,
        UpdatedAt: user.UpdatedAt,
    }

    tokenResponse := map[string]any{
        "access_token":  accessToken,
        "token_type":    "Bearer",
        "expires_in":    int(global.Cfg.OAuth.AccessTokenTTL.Seconds()),
        "refresh_token": refreshToken,
        "scope":         strings.Join(scopes, " "),
    }

    return userResponse, tokenResponse, nil
}
```

2. **修改登录控制器** - 文件：`server/controllers/user.go`

```go
// 修改LoginHandler以支持OAuth
func LoginHandler(c *gin.Context) {
    var req models.LoginRequest

    if !validator.ValidateStruct(c, &req) {
        response.Error(c, apperrors.ErrInvalidInput)
        return
    }

    // 从请求头或参数中获取客户端ID
    clientID := c.GetHeader("X-Client-ID")
    if clientID == "" {
        clientID = c.Query("client_id")
    }
    if clientID == "" {
        // 如果没有指定客户端ID，使用默认系统客户端
        clientID = "your-system-client-id" // 替换为实际的系统客户端ID
    }

    user, tokenResponse, err := services.OAuthLoginService(&req, clientID)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.Success(c, "登录成功", gin.H{
        "user":  user,
        "token": tokenResponse,
    })
}
```

### 1.3 配置OAuth作为主要认证方式

**目标**：将OAuth中间件设置为主要认证方式

**步骤**：

1. **修改中间件管理器** - 文件：`server/middleware/manager.go`

```go
// 添加统一认证中间件，优先使用OAuth
func (m *Manager) UnifiedAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            response.Error(c, apperrors.ErrUnauthorized)
            c.Abort()
            return
        }

        tokenParts := strings.Split(authHeader, " ")
        if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
            response.Error(c, apperrors.ErrTokenInvalid)
            c.Abort()
            return
        }

        token := tokenParts[1]

        // 首先尝试OAuth token解析
        if claims, err := oauth.ParseOAuthToken(token); err == nil {
            // OAuth token 验证成功
            c.Set("user_id", claims.UserID)
            c.Set("client_id", claims.ClientID)
            c.Set("scopes", claims.Scopes)
            c.Set("token_type", "oauth")
            c.Next()
            return
        }

        // 如果OAuth解析失败，尝试JWT（兼容性支持）
        if claims, err := auth.ParseToken(token); err == nil {
            // JWT token 验证成功
            c.Set("user_id", claims.UserID)
            c.Set("username", claims.Username)
            c.Set("role", claims.Role)
            c.Set("token_type", "jwt")
            c.Next()
            return
        }

        // 两种token都解析失败
        response.Error(c, apperrors.ErrTokenInvalid)
        c.Abort()
    }
}
```

2. **更新路由配置** - 文件：`server/routers/user.go`

```go
func LoadUserRouters(router *gin.Engine) {
    m := middleware.NewManager()
    userRouters := router.Group("/api/account/v1/users")
    {
        // 使用统一认证中间件
        userRouters.PUT("/:user_id", m.UnifiedAuth(), controllers.UpdateHandler)
        userRouters.GET("", m.UnifiedAuth(), m.AdminPermission(), controllers.UserListHandler)
        userRouters.GET("/:user_id", m.UnifiedAuth(), controllers.UserProfileHandler)
    }
}
```

## 第二阶段：权限系统OAuth化

### 2.1 实现基于Scopes的权限控制

**目标**：使用OAuth scopes替代基于角色的权限控制

**步骤**：

1. **创建Scopes权限中间件** - 文件：`server/middleware/permission/scopes.go`

```go
package permission

import (
    "YAccount/pkg/apperrors"
    "YAccount/pkg/response"
    "github.com/gin-gonic/gin"
)

// Scopes权限中间件
func RequiredScopes(requiredScopes ...string) gin.HandlerFunc {
    return func(c *gin.Context) {
        tokenType := c.GetString("token_type")
        
        if tokenType == "oauth" {
            // OAuth token 使用 scopes 验证
            userScopes, exists := c.Get("scopes")
            if !exists {
                response.Error(c, apperrors.ErrInsufficientScope)
                c.Abort()
                return
            }

            scopes, ok := userScopes.([]string)
            if !ok {
                response.Error(c, apperrors.ErrInsufficientScope)
                c.Abort()
                return
            }

            // 检查是否有所需的权限范围
            for _, required := range requiredScopes {
                found := false
                for _, scope := range scopes {
                    if scope == required {
                        found = true
                        break
                    }
                }
                if !found {
                    response.Error(c, apperrors.ErrInsufficientScope)
                    c.Abort()
                    return
                }
            }
        } else if tokenType == "jwt" {
            // JWT token 兼容性处理，基于角色映射到scopes
            role := c.GetString("role")
            userScopes := rolesToScopes(role)

            for _, required := range requiredScopes {
                found := false
                for _, scope := range userScopes {
                    if scope == required {
                        found = true
                        break
                    }
                }
                if !found {
                    response.Error(c, apperrors.ErrInsufficientScope)
                    c.Abort()
                    return
                }
            }
        } else {
            response.Error(c, apperrors.ErrUnauthorized)
            c.Abort()
            return
        }

        c.Next()
    }
}

// 角色到权限范围的映射（兼容性支持）
func rolesToScopes(role string) []string {
    switch role {
    case "admin":
        return []string{"read", "write", "admin"}
    case "user":
        return []string{"read"}
    default:
        return []string{}
    }
}
```

2. **更新中间件管理器** - 文件：`server/middleware/manager.go`

```go
// 添加权限范围中间件
func (m *Manager) RequiredScopes(scopes ...string) gin.HandlerFunc {
    return permission.RequiredScopes(scopes...)
}
```

3. **更新路由配置** - 文件：`server/routers/user.go`

```go
func LoadUserRouters(router *gin.Engine) {
    m := middleware.NewManager()
    userRouters := router.Group("/api/account/v1/users")
    {
        userRouters.PUT("/:user_id", m.UnifiedAuth(), controllers.UpdateHandler)
        userRouters.GET("", m.UnifiedAuth(), m.RequiredScopes("admin"), controllers.UserListHandler)
        userRouters.GET("/:user_id", m.UnifiedAuth(), controllers.UserProfileHandler)
    }
}
```

### 2.2 更新OAuth路由权限

**目标**：使用scopes替代管理员权限检查

**步骤**：

**修改OAuth路由** - 文件：`server/routers/oauth.go`

```go
func LoadOAuthRouters(router *gin.Engine) {
    m := middleware.NewManager()

    // OAuth 标准端点
    oauthGroup := router.Group("/api/account/v1/oauth")
    {
        // 授权端点（需要用户登录）
        oauthGroup.GET("/authorize", m.UnifiedAuth(), controllers.OAuthAuthorizeHandler)

        // 令牌端点（公开）
        oauthGroup.POST("/token", controllers.OAuthTokenHandler)

        // 令牌内省端点（公开）
        oauthGroup.POST("/introspect", controllers.OAuthIntrospectHandler)
    }

    clientGroup := router.Group("/api/account/v1/oauth/clients")
    {
        // 使用scopes替代AdminPermission
        clientGroup.POST("", m.UnifiedAuth(), m.RequiredScopes("admin"), controllers.OAuthClientRegisterHandler)
        clientGroup.GET("", m.UnifiedAuth(), m.RequiredScopes("admin"), controllers.ListOAuthClientsHandler)
        clientGroup.GET("/:client_id", m.UnifiedAuth(), m.RequiredScopes("admin"), controllers.GetOAuthClientHandler)
        clientGroup.PUT("/:client_id", m.UnifiedAuth(), m.RequiredScopes("admin"), controllers.UpdateOAuthClientHandler)
        clientGroup.DELETE("/:client_id", m.UnifiedAuth(), m.RequiredScopes("admin"), controllers.DeleteOAuthClientHandler)
    }
}
```

## 第三阶段：Token刷新机制OAuth化

### 3.1 实现OAuth刷新令牌端点

**目标**：使用OAuth标准的刷新令牌机制替代JWT刷新

**步骤**：

1. **修改认证路由** - 文件：`server/routers/auth.go`

```go
func LoadAuthRouters(router *gin.Engine) {
    m := middleware.NewManager()
    authRouters := router.Group("/api/account/v1/auth")
    {
        authRouters.POST("/login", controllers.LoginHandler)
        authRouters.POST("/register", controllers.RegisterHandler)
        // 使用OAuth的token端点替代JWT refresh
        authRouters.POST("/refresh", controllers.OAuthRefreshHandler)
    }
}
```

2. **实现OAuth刷新控制器** - 文件：`server/controllers/auth.go`

```go
// OAuth刷新令牌处理器
func OAuthRefreshHandler(c *gin.Context) {
    var req struct {
        RefreshToken string `json:"refresh_token" binding:"required"`
        ClientID     string `json:"client_id"`
    }

    if err := c.ShouldBindJSON(&req); err != nil {
        response.Error(c, apperrors.ErrInvalidInput)
        return
    }

    // 如果没有提供客户端ID，使用默认系统客户端
    if req.ClientID == "" {
        req.ClientID = "your-system-client-id" // 替换为实际的系统客户端ID
    }

    // 构建标准OAuth刷新请求
    tokenReq := services.TokenRequest{
        GrantType:    "refresh_token",
        RefreshToken: req.RefreshToken,
        ClientID:     req.ClientID,
        ClientSecret: "your-system-client-secret", // 从配置获取
    }

    tokenResponse, err := services.HandleTokenRequest(&tokenReq)
    if err != nil {
        response.Error(c, err)
        return
    }

    response.Success(c, "刷新令牌成功", tokenResponse)
}
```

## 第四阶段：移除JWT依赖

### 4.1 修改配置文件

**目标**：移除JWT相关配置，保留OAuth配置

**步骤**：

1. **更新配置结构** - 文件：`server/configs/config.go`

```go
type Config struct {
    Database   DatabaseConfig   `mapstructure:"database"`
    Redis      RedisConfig      `mapstructure:"redis"`
    Log        LogConfig        `mapstructure:"log"`
    Server     ServerConfig     `mapstructure:"server"`
    OAuth      OAuthConfig      `mapstructure:"oauth"`      // 保留
    // JWT        JWTConfig        `mapstructure:"jwt"`      // 移除
    Middleware MiddlewareConfig `mapstructure:"middleware"`
}
```

2. **更新配置文件示例** - 文件：`server/configs/config.example.yaml`

```yaml
# 移除jwt配置块
# jwt:
#   secret: "your-jwt-secret"
#   expire: "24h"

# 保留并完善oauth配置
oauth:
  issuer: "https://yourdomain.com"
  authorize_ui: "https://yourdomain.com/oauth/authorize"
  authorization_code_ttl: "10m"
  access_token_ttl: "1h"
  refresh_token_ttl: "168h" # 7天
  supported_grant_types:
    - "authorization_code"
    - "refresh_token" 
    - "client_credentials"
  supported_response_types:
    - "code"
  default_scopes:
    - "read"
  require_pkce: false
  supported_challenge_methods:
    - "S256"
```

### 4.2 更新中间件配置

**目标**：移除JWT中间件配置

**步骤**：

**修改中间件配置** - 文件：`server/configs/middleware.go`

```go
type MiddlewareConfig struct {
    CORS CORSConfig `mapstructure:"cors"`
    // JWT  JWTMiddlewareConfig `mapstructure:"jwt"` // 移除
}

// 移除JWTMiddlewareConfig结构体
// type JWTMiddlewareConfig struct {
//     SkipPaths []string `mapstructure:"skip_paths"`
// }
```

### 4.3 清理JWT相关文件

**目标**：移除不再需要的JWT文件

**执行以下清理操作**：

```bash
# 备份重要文件
mkdir -p backup/jwt
cp server/configs/jwt.go backup/jwt/
cp server/pkg/auth/jwt.go backup/jwt/
cp server/middleware/auth/jwt.go backup/jwt/

# 删除JWT相关文件
rm server/configs/jwt.go
rm server/pkg/auth/jwt.go  
rm server/middleware/auth/jwt.go
rm -rf server/middleware/auth/  # 如果目录为空
```

### 4.4 更新导入和引用

**目标**：移除所有JWT相关的导入和引用

**步骤**：

1. **更新全局配置初始化** - 文件：`server/global/global.go`

```go
// 移除JWT相关的全局变量和初始化
```

2. **更新服务文件** - 检查以下文件并移除JWT导入：
   - `server/services/user.go` - 移除 `auth.GenerateToken` 调用
   - `server/controllers/auth.go` - 移除JWT刷新逻辑

## 第五阶段：前端适配

### 5.1 更新前端API调用

**目标**：修改前端以使用OAuth令牌格式

**步骤**：

1. **更新认证API** - 文件：`web/src/api/auth.ts`

```typescript
// 更新登录响应接口
interface LoginResponse {
  user: User;
  token: {
    access_token: string;
    token_type: string;
    expires_in: number;
    refresh_token: string;
    scope: string;
  };
}

// 更新刷新令牌方法
export const refreshToken = async (refreshToken: string) => {
  const response = await request.post<{
    access_token: string;
    token_type: string; 
    expires_in: number;
    refresh_token: string;
    scope: string;
  }>('/auth/refresh', {
    refresh_token: refreshToken,
    client_id: 'your-system-client-id' // 配置中获取
  });
  return response.data;
};
```

2. **更新token存储** - 文件：`web/src/stores/auth.ts`

```typescript
interface TokenInfo {
  access_token: string;
  token_type: string;
  expires_in: number;
  refresh_token: string;
  scope: string;
  expires_at: number; // 计算得出的过期时间
}

export const useAuthStore = defineStore('auth', () => {
  // 更新token相关的存储和处理逻辑
  const tokenInfo = ref<TokenInfo | null>(null);

  const setToken = (token: TokenInfo) => {
    token.expires_at = Date.now() + (token.expires_in * 1000);
    tokenInfo.value = token;
    localStorage.setItem('oauth_token', JSON.stringify(token));
  };

  const getAccessToken = () => {
    return tokenInfo.value?.access_token || localStorage.getItem('oauth_token') 
      ? JSON.parse(localStorage.getItem('oauth_token')!).access_token 
      : null;
  };

  // 其他方法更新...
});
```

### 5.2 更新HTTP拦截器

**目标**：适配OAuth令牌格式

**步骤**：

**修改请求拦截器** - 更新axios配置以使用OAuth token：

```typescript
// 请求拦截器
request.interceptors.request.use((config) => {
  const authStore = useAuthStore();
  const token = authStore.getAccessToken();
  
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  
  return config;
});
```

## 第六阶段：测试和验证

### 6.1 功能测试清单

**测试项目**：

- [ ] OAuth登录流程
- [ ] Token刷新机制  
- [ ] 权限控制(scopes)
- [ ] 客户端管理
- [ ] API访问控制
- [ ] 前端集成

### 6.2 性能测试

**测试项目**：

- [ ] Token生成性能
- [ ] 权限验证性能
- [ ] 数据库查询优化

### 6.3 安全测试

**测试项目**：

- [ ] Token安全性
- [ ] 权限控制有效性
- [ ] 客户端认证

## 迁移时间规划

### 建议分阶段实施：

1. **第1-2周**：第一、二阶段（OAuth登录和权限系统）
2. **第3周**：第三阶段（Token刷新机制）  
3. **第4周**：第四阶段（移除JWT依赖）
4. **第5周**：第五阶段（前端适配）
5. **第6周**：第六阶段（测试和验证）

## 风险控制

### 回滚计划：

1. **保留备份**：所有删除的文件先备份
2. **分支管理**：在新分支进行改造
3. **灰度发布**：逐步切换用户流量
4. **监控告警**：设置关键指标监控

## 注意事项

1. **配置管理**：确保正确配置OAuth客户端信息
2. **权限映射**：确保角色到scopes的映射正确
3. **Token管理**：注意OAuth token的生命周期管理
4. **兼容性**：保持API接口的向后兼容性
5. **安全性**：确保client_secret的安全存储

## 支持文档

- [OAuth 2.0 RFC 6749](https://tools.ietf.org/html/rfc6749)
- [OAuth 2.0 Security Best Practices](https://tools.ietf.org/html/draft-ietf-oauth-security-topics)
- 当前系统的OAuth使用手册.md

---

**免责声明**：建议在生产环境实施前，先在测试环境进行完整的验证测试。
