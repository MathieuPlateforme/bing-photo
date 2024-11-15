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
	secretKey  []byte
}

func NewJWTService() (*JWTService, error) {
	// Initialiser un nouveau service JWT

	// Charger la clé secrète à partir des variables d'environnement
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	if len(secretKey) == 0 {
		return nil, fmt.Errorf("clé secrète JWT non configurée")
	}

	return &JWTService{
		secretKey: secretKey,
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
	tokenString, err := token.SignedString(s.secretKey)
	if err != nil {
		return "", err
	}

	// Mettre à jour le champ Token dans la structure
	s.Token = tokenString
	return tokenString, nil
}

func (s *JWTService) VerifyToken(tokenString string) error {
	// Analyser le token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Vérifier la méthode de signature
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("méthode de signature inattendue : %v", token.Header["alg"])
		}
		return s.secretKey, nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("token invalide")
	}

	return nil
}
