package google

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// GoogleAuthService structure
type GoogleAuthService struct {
	Config *oauth2.Config
}
type GoogleUserProfile struct {
	ID      string `json:"sub"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
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
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	return &GoogleAuthService{Config: config}, nil
}

func (s *GoogleAuthService) AuthenticateWithGoogle() (string, error) {
	if s.Config == nil {
		return "", errors.New("OAuth2 configuration is not initialized")
	}

	authURL := s.Config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	if authURL == "" {
		return "", errors.New("Failed to generate Google authentication URL")
	}

	return authURL, nil
}

func (s *GoogleAuthService) GetGoogleUserProfile(token *oauth2.Token) (*GoogleUserProfile, error) {
	client := s.Config.Client(oauth2.NoContext, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v3/userinfo")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var userProfile GoogleUserProfile
	if err := json.NewDecoder(resp.Body).Decode(&userProfile); err != nil {
		return nil, err
	}

	return &userProfile, nil
}

