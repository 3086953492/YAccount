package services

import (
	"YAccount/global"
	"YAccount/models"
	"YAccount/pkg/apperrors"
	"YAccount/repositories"
	"YAccount/utils/logger"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/go-redis/cache/v9"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

// 创建 OAuth 客户端
func CreateOAuthClient(req *models.CreateOAuthClientRequest, ownerID uint) (*models.OAuthClient, error) {
	// 生成客户端ID和密钥
	clientID, err := generateClientID()
	if err != nil {
		return nil, apperrors.ErrServerInternal
	}

	clientSecret, err := generateClientSecret()
	if err != nil {
		return nil, apperrors.ErrServerInternal
	}

	// 加密客户端密钥
	hashedSecret, err := bcrypt.GenerateFromPassword([]byte(clientSecret), bcrypt.DefaultCost)
	if err != nil {
		return nil, apperrors.ErrServerInternal
	}

	// 序列化数组字段
	redirectURIsJSON, _ := json.Marshal(req.RedirectURIs)
	grantTypesStr := strings.Join(req.GrantTypes, ",")
	scopesStr := strings.Join(req.Scopes, ",")

	client := &models.OAuthClient{
		ClientID:     clientID,
		ClientSecret: string(hashedSecret),
		Name:         req.Name,
		Description:  req.Description,
		RedirectURIs: string(redirectURIsJSON),
		GrantTypes:   grantTypesStr,
		Scopes:       scopesStr,
		ClientType:   req.ClientType,
		OwnerID:      ownerID,
		OwnerType:    "user",
	}

	// 保存到数据库
	if err := repositories.CreateOAuthClient(client); err != nil {
		logger.LogError("CreateOAuthClient", "database", "创建OAuth客户端失败", err, zap.Any("client", client))
		return nil, apperrors.ErrServerInternal
	}

	// 在响应中包含未加密的客户端密钥（仅此一次）
	client.ClientSecret = clientSecret

	return client, nil
}

// 验证客户端凭证
func ValidateClientCredentials(clientID, clientSecret string) (*models.OAuthClient, error) {

	var client models.OAuthClient

	client, err := repositories.GetOAuthClientByID(clientID)
	if err != nil {
		if !apperrors.IsNotFoundError(err) {
			logger.LogError("ValidateClientCredentials", "database", "获取OAuth客户端失败", err, zap.String("clientID", clientID))
		}
		return nil, err
	}

	// 验证客户端密钥
	if err := bcrypt.CompareHashAndPassword([]byte(client.ClientSecret), []byte(clientSecret)); err != nil {
		return nil, apperrors.ErrInvalidClient
	}

	return &client, nil
}

// 验证重定向URI
func ValidateRedirectURI(client *models.OAuthClient, redirectURI string) bool {
	var redirectURIs []string
	if err := json.Unmarshal([]byte(client.RedirectURIs), &redirectURIs); err != nil {
		return false
	}

	return slices.Contains(redirectURIs, redirectURI)
}

// 生成客户端ID
func generateClientID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return fmt.Sprintf("client_%s", base64.URLEncoding.EncodeToString(bytes)[:22]), nil
}

// 生成客户端密钥
func generateClientSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// 获取客户端信息
func GetOAuthClientByID(clientID string) (*models.OAuthClient, error) {
	var client models.OAuthClient
	if err := global.Cache.Once(&cache.Item{
		Key:   fmt.Sprintf("oauth_client:%s", clientID),
		Value: &client,
		Do: func(*cache.Item) (any, error) {
			client, err := repositories.GetOAuthClientByID(clientID)
			if err != nil {
				if !apperrors.IsNotFoundError(err) {
					logger.LogError("GetOAuthClientByID", "database", "获取OAuth客户端失败", err, zap.String("clientID", clientID))
				}
				return nil, err
			}
			return client, nil
		},
	}); err != nil {
		return nil, apperrors.ErrInvalidClient
	}
	return &client, nil
}

// 获取客户端列表
func ListOAuthClients() ([]models.OAuthClient, error) {
	clients, err := repositories.ListOAuthClients()
	if err != nil {
		if !apperrors.IsNotFoundError(err) {
			logger.LogError("ListOAuthClients", "database", "获取OAuth客户端列表失败", err)
		}
		return nil, apperrors.ErrInvalidClient
	}
	return clients, nil
}

// 获取客户端详情
func GetOAuthClientDetail(clientID string, userID uint, role string) (*models.OAuthClient, error) {
	client, err := repositories.GetOAuthClientDetail(clientID, userID, role)
	if err != nil {
		if !apperrors.IsNotFoundError(err) {
			logger.LogError("GetOAuthClientDetail", "database", "获取OAuth客户端详情失败", err, zap.String("clientID", clientID))
		}
		return nil, apperrors.ErrInvalidClient
	}
	return &client, nil
}

func UpdateOAuthClient(clientID string, request models.UpdateOAuthClientRequest) error {
	// 转换切片字段为数据库格式，就像创建时一样
	redirectURIsJSON, _ := json.Marshal(request.RedirectURIs)
	grantTypesStr := strings.Join(request.GrantTypes, ",")
	scopesStr := strings.Join(request.Scopes, ",")

	// 创建用于数据库更新的结构体，使用与OAuthClient模型相同的字段类型
	updateData := map[string]any{
		"name":          request.Name,
		"description":   request.Description,
		"redirect_uris": string(redirectURIsJSON),
		"grant_types":   grantTypesStr,
		"scopes":        scopesStr,
		"client_type":   request.ClientType,
	}

	if err := repositories.UpdateOAuthClient(clientID, updateData); err != nil {
		logger.LogError("UpdateOAuthClient", "database", "更新OAuth客户端失败", err, zap.String("clientID", clientID))
		return apperrors.ErrUpdateClientFailed
	}
	return nil
}

func DeleteOAuthClient(clientID string) error {
	if err := repositories.DeleteOAuthClient(clientID); err != nil {
		logger.LogError("DeleteOAuthClient", "database", "删除OAuth客户端失败", err, zap.String("clientID", clientID))
		return apperrors.ErrDeleteClientFailed
	}
	return nil
}
