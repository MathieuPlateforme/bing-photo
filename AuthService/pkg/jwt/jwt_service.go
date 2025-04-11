package jwt

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"google.golang.org/grpc/metadata"
)

type JWTService struct {
	Token      string
	Expiration int64
	IssuedAt   int64
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

func (j *JWTService) VerifyToken(tokenString string) (map[string]interface{}, error) {
	fmt.Println("Vérification du token brut :", tokenString)

	// Nettoyer "Bearer " si jamais il est encore présent (par précaution)
	tokenString = strings.TrimPrefix(tokenString, "Bearer ")
	tokenString = strings.TrimSpace(tokenString)

	if tokenString == "" {
		return nil, errors.New("token vide ou mal formaté")
	}

	// Parse le token et vérifie la signature
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// S'assurer que l'algo est bien HMAC (HS256)
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("méthode de signature inattendue : %v", token.Header["alg"])
		}
		return []byte(j.SecretKey), nil
	})
	if err != nil || !token.Valid {
		return nil, fmt.Errorf("token invalide ou signature incorrecte")
	}

	// Récupérer les claims
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("claims introuvables ou invalides")
	}

	// Vérification manuelle de l'expiration
	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return nil, errors.New("le token est expiré")
		}
	}

	return claims, nil
}

func (j *JWTService) GenerateToken(userID uint, username string) (string, error) {
	claims := jwt.MapClaims{
		"userID":   userID,
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(j.SecretKey))
}

// VerifyTokenFromContext extrait le token du contexte et le vérifie
func (j *JWTService) VerifyTokenFromContext(ctx context.Context) (map[string]interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("métadonnées manquantes dans le contexte")
	}

	authHeaders := md["authorization"]
	if len(authHeaders) == 0 {
		return nil, errors.New("en-tête Authorization manquant")
	}

	token := strings.TrimPrefix(authHeaders[0], "Bearer ")
	if token == "" {
		return nil, errors.New("token vide ou mal formaté")
	}

	// Appelle ta méthode de vérification déjà existante
	return j.VerifyToken(token)
}

func (j *JWTService) ExtractTokenFromContext(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", errors.New("aucune métadonnée dans le contexte")
	}

	authHeaders := md["authorization"]
	if len(authHeaders) == 0 {
		return "", errors.New("en-tête Authorization manquant")
	}

	token := strings.TrimSpace(strings.TrimPrefix(authHeaders[0], "Bearer "))
	if token == "" {
		return "", errors.New("token vide après extraction")
	}

	return token, nil
}
