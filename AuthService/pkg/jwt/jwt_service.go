package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type JWTService struct {
	Token     string
	Expiration int64
	IssuedAt  int64
	SecretKey  []byte
}

func NewJWTService() (*JWTService, error) {
	// Initialiser un nouveau service JWT

	// Charger la clé secrète à partir des variables d'environnement
	SecretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	if len(SecretKey) == 0 {
		return nil, fmt.Errorf("clé secrète JWT non configurée")
	}

	return &JWTService{
		SecretKey: SecretKey,
	}, nil
}

func (s *JWTService) GenerateToken(username string) (string, error) {
	// Définir l'heure actuelle et l'expiration
	s.IssuedAt = time.Now().Unix()
	s.Expiration = time.Now().Add(24 * time.Hour).Unix()

	// Créer un nouveau token JWT avec les claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      s.Expiration,
		"iat":      s.IssuedAt,
	})

	// Signer le token avec la clé secrète
	tokenString, err := token.SignedString(s.SecretKey)
	if err != nil {
		return "", err
	}

	// Mettre à jour le champ Token dans la structure
	s.Token = tokenString
	return tokenString, nil
}

func (j *JWTService) VerifyToken(token string) (map[string]interface{}, error) {
	// Parse et valide le token
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})
	if err != nil || !parsedToken.Valid {
		return nil, fmt.Errorf("token invalide ou expiré")
	}

	// Extraire les claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("claims introuvables")
	}

	return claims, nil
}

