package services

import (
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
	client, err := repositories.GetOAuthClientByID(clientID)
	if err != nil {
		if !apperrors.IsNotFoundError(err) {
			logger.LogError("GetOAuthClientByID", "database", "获取OAuth客户端失败", err, zap.String("clientID", clientID))
		}
		return nil, apperrors.ErrInvalidClient
	}
	return &client, nil
}
