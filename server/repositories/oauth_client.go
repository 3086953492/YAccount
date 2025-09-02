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

func ListOAuthClients() ([]models.OAuthClient, error) {
	var clients []models.OAuthClient
	if err := global.DB.Where("status = ?", "active").Find(&clients).Error; err != nil {
		return nil, err
	}
	return clients, nil
}

func GetOAuthClientDetail(clientID string, userID uint, role string) (models.OAuthClient, error) {
	var client models.OAuthClient
	if err := global.DB.Where("client_id = ? AND status = ?", clientID, "active").First(&client).Error; err != nil {
		return models.OAuthClient{}, err
	}
	return client, nil
}

func UpdateOAuthClient(clientID string, updateData map[string]any) error {
	if err := global.DB.Model(&models.OAuthClient{}).Where("client_id = ?", clientID).Updates(updateData).Error; err != nil {
		return err
	}
	return nil
}

func DeleteOAuthClient(clientID string) error {
	if err := global.DB.Model(&models.OAuthClient{}).Where("client_id = ?", clientID).Delete(&models.OAuthClient{}).Error; err != nil {
		return err
	}
	return nil
}