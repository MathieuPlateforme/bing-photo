package google

import (
	"fmt"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"os"
	"errors"
	"encoding/json"
	"net/http"
)

// GoogleAuthService structure
type GoogleAuthService struct {
	Config *oauth2.Config
}

// NewGoogleAuthService initialise et retourne une nouvelle instance de GoogleAuthService
func NewGoogleAuthService() (*GoogleAuthService, error) {
	fmt.Println("Initializing GoogleAuthService...")

	// Initialiser la configuration OAuth2 pour Google
	clientID := os.Getenv("GOOGLE_CLIENT_ID")
	clientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	redirectURL := os.Getenv("GOOGLE_REDIRECT_URL")

	if clientID == "" || clientSecret == "" || redirectURL == "" {
		return nil, errors.New("Google OAuth2 configuration is missing")
	}

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email","https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	return &GoogleAuthService{Config: config}, nil
}

func (s *GoogleAuthService) AuthenticateWithGoogle() string {
	// Génère et renvoie l'URL d'authentification Google
	return s.Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
}

func (s *GoogleAuthService) GetGoogleUserProfile(token *oauth2.Token) (map[string]interface{}, error) {
	client := s.Config.Client(oauth2.NoContext, token)
	userInfoResp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, fmt.Errorf("Échec de la récupération des informations utilisateur : %v", err)
	}
	defer userInfoResp.Body.Close()

	// Vérifiez le code de statut HTTP de la réponse
	if userInfoResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Requête échouée avec le code de statut : %d", userInfoResp.StatusCode)
	}

	userInfo := map[string]interface{}{}
	if err := json.NewDecoder(userInfoResp.Body).Decode(&userInfo); err != nil {
		return nil, fmt.Errorf("Échec de la lecture des informations utilisateur : %v", err)
	}

	return userInfo, nil
}
