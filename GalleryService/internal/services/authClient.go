package services

import (
	"context"
	"fmt"
	"net/http"
)

type AuthServiceClient struct {
	APIURL string
}

func NewAuthServiceClient(apiURL string) *AuthServiceClient {
	return &AuthServiceClient{APIURL: apiURL}
}


func (c *AuthServiceClient) ValidateToken(ctx context.Context, token string) (bool, error) {

	// Créer une requête POST pour valider le token
	req, err := http.NewRequest("POST", c.APIURL+"/validateToken", nil)
	if err != nil {
		return false, fmt.Errorf("échec de la création de la requête : %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	// Exécuter la requête
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("échec de l'exécution de la requête : %v", err)
	}
	defer resp.Body.Close()

	// Vérifier le code de statut
	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("statut de la réponse non OK : %s", resp.Status)
	}

	return true, nil
}
