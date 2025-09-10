package repositories

import (
	"YAccount/models"
	"github.com/3086953492/YaBase/global"

	"gorm.io/gorm"
)

func oauthDB() *gorm.DB {
	return global.GetGlobalDB()
}

func CreateOAuthAuthorizationCode(authCode *models.OAuthAuthorizationCode) error {
	if err := oauthDB().Create(authCode).Error; err != nil {
		return err
	}
	return nil
}

func GetOAuthAuthorizationCodeByCode(code, clientID string, used bool) error {
	if err := oauthDB().Where("code = ? AND client_id = ? AND used = ?", code, clientID, used).First(&models.OAuthAuthorizationCode{}).Error; err != nil {
		return err
	}
	return nil
}

func UpdateOAuthAuthorizationCodeUsed(id uint, used bool) error {
	if err := oauthDB().Model(&models.OAuthAuthorizationCode{}).Where("id = ?", id).Update("used", used).Error; err != nil {
		return err
	}
	return nil
}

func CreateOAuthAccessToken(tokenRecord *models.OAuthAccessToken) error {
	if err := oauthDB().Create(tokenRecord).Error; err != nil {
		return err
	}
	return nil
}

func GetOAuthAccessTokenByRefreshToken(refreshToken, clientID string, revoked bool) error {
	if err := oauthDB().Where("refresh_token = ? AND client_id = ? AND revoked = ?", refreshToken, clientID, revoked).First(&models.OAuthAccessToken{}).Error; err != nil {
		return err
	}
	return nil
}

func UpdateOAuthAccessTokenRevoked(id uint, revoked bool) error {
	if err := oauthDB().Model(&models.OAuthAccessToken{}).Where("id = ?", id).Update("revoked", revoked).Error; err != nil {
		return err
	}
	return nil
}
