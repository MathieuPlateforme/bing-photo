package google_test

import (
	"AuthService/pkg/google"
	"golang.org/x/oauth2"
	"testing"
	"os"
)

// Helper function to initialize GoogleAuthService
func initGoogleAuthService(t *testing.T) *google.GoogleAuthService {
	googleAuthService, err := google.NewGoogleAuthService()
	if err != nil {
		t.Fatalf("Erreur lors de l'initialisation de GoogleAuthService : %v", err)
	}
	return googleAuthService
}

// Test pour s'assurer que GoogleAuthService est initialisé correctement
func TestNewGoogleAuthService(t *testing.T) {
	googleAuthService := initGoogleAuthService(t)
	if googleAuthService == nil {
		t.Error("GoogleAuthService n'est pas initialisé")
	}
	t.Log("GoogleAuthService initialisé avec succès")
}

// Test pour vérifier que AuthenticateWithGoogle renvoie une URL d'authentification valide
func TestAuthenticateWithGoogle(t *testing.T) {
	googleAuthService := initGoogleAuthService(t)
	authURL,err := googleAuthService.AuthenticateWithGoogle()
	if(err != nil){	
		t.Errorf("Erreur lors de la génération de l'URL d'authentification Google : %v", err)
	}
	if authURL == "" {
		t.Error("URL d'authentification Google vide")
	}
	t.Logf("URL d'authentification Google : %v", authURL)
}

// Helper function to simulate OAuth2 token
func getOAuth2Token(accessToken string) *oauth2.Token {
	return &oauth2.Token{
		AccessToken: accessToken,
	}
}

// Test pour vérifier que GetGoogleUserProfile renvoie des informations utilisateur
func TestGetGoogleUserProfile(t *testing.T) {
	googleAuthService := initGoogleAuthService(t)
	token := getOAuth2Token(os.Getenv("GOOGLE_ACCESS_TOKEN")) // Assurez-vous que "GOOGLE_ACCESS_TOKEN" est défini dans votre .env

	_, err := googleAuthService.GetGoogleUserProfile(token)
	if err != nil {
		t.Errorf("Erreur lors de la récupération des informations utilisateur : %v", err)
	} else {
		t.Log("Informations utilisateur récupérées avec succès")
	}
}

// Test pour vérifier le comportement avec un jeton OAuth2 vide
func TestGetGoogleUserProfileEmptyToken(t *testing.T) {
	googleAuthService := initGoogleAuthService(t)
	token := &oauth2.Token{} // Jeton OAuth2 vide

	_, err := googleAuthService.GetGoogleUserProfile(token)
	if err == nil {
		t.Error("Aucune erreur retournée pour un jeton OAuth2 vide")
	} else {
		t.Logf("Erreur correcte pour un jeton OAuth2 vide : %v", err)
	}
}

// Test pour vérifier le comportement avec un jeton OAuth2 invalide
func TestGetGoogleUserProfileInvalidToken(t *testing.T) {
	googleAuthService := initGoogleAuthService(t)
	token := getOAuth2Token("invalid_token") 

	_, err := googleAuthService.GetGoogleUserProfile(token)
	if err == nil {
		t.Error("Aucune erreur retournée pour un jeton OAuth2 invalide")
	} else {
		t.Logf("Erreur correcte pour un jeton OAuth2 invalide : %v", err)
	}
}

// Test pour vérifier le comportement avec une URL d'API Google invalide
func TestGetGoogleUserProfileInvalidURL(t *testing.T) {
	googleAuthService := initGoogleAuthService(t)

	// Changer l'URL de l'API pour une URL invalide
	googleAuthService.Config.Endpoint = oauth2.Endpoint{
		AuthURL:  "https://invalid-url.com/auth",
		TokenURL: "https://invalid-url.com/token",
	}

	token := getOAuth2Token("valid_token")
	_, err := googleAuthService.GetGoogleUserProfile(token)
	if err == nil {
		t.Error("Aucune erreur retournée pour une URL d'API Google invalide")
	} else {
		t.Logf("Erreur correcte pour une URL d'API Google invalide : %v", err)
	}
}
