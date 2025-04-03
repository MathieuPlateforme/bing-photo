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
	// Parse et valide le token
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(j.SecretKey), nil
	})
	if err != nil || !parsedToken.Valid {
		return nil, fmt.Errorf("token invalide ou expiré 5")
	}

	// Extraire les claims
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("claims introuvables")
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

func ParseToken(tokenStr string) (jwt.MapClaims, error) {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		return nil, errors.New("la variable d'environnement JWT_SECRET est manquante")
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("méthode de signature inattendue : %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("échec du parsing du token : %v", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("token invalide ou claims manquants")
	}

	if exp, ok := claims["exp"].(float64); ok {
		if int64(exp) < time.Now().Unix() {
			return nil, errors.New("le token est expiré")
		}
	}
	return claims, nil
}


func ExtractUserIDFromContext(ctx context.Context) (uint, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return 0, fmt.Errorf("aucune metadata dans le contexte")
	}

	// Supporte les deux formes : "authorization" ou "Authorization"
	var tokenStr string
	if authHeaders := md.Get("authorization"); len(authHeaders) > 0 {
		tokenStr = strings.TrimSpace(authHeaders[0])
	} else if authHeaders := md.Get("Authorization"); len(authHeaders) > 0 {
		tokenStr = strings.TrimSpace(authHeaders[0])
	} else {
		return 0, fmt.Errorf("authorization header manquant")
	}

	// Supprime "Bearer " s'il est présent
	tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")
	tokenStr = strings.TrimSpace(tokenStr)

	// Parse le token
	claims, err := ParseToken(tokenStr)
	if err != nil {
		return 0, fmt.Errorf("token invalide : %v", err)
	}

	userIDFloat, ok := claims["userID"].(float64)
	if !ok {
		return 0, fmt.Errorf("userID non présent ou invalide dans le token")
	}

	return uint(userIDFloat), nil
}