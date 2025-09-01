package repositories

import (
	"YAccount/global"
	"YAccount/models"
)

func CreateOAuthClient(client *models.OAuthClient) error {
	if err := global.DB.Create(client).Error; err != nil {
		return err
	}
	return nil
}

func GetOAuthClientByID(clientID string) (models.OAuthClient, error) {
	var client models.OAuthClient
	if err := global.DB.Where("client_id = ? AND status = ?", clientID, "active").First(&client).Error; err != nil {
		return models.OAuthClient{}, err
	}
	return client, nil
}
