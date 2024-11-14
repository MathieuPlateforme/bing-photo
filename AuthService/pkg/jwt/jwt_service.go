package jwt

type JWTService struct {
	Token     string
	Expiration int64
	IssuedAt  int64
}

func NewJWTService() (*JWTService, error) {
	// Initialiser un nouveau service JWT
	return &JWTService{}, nil
}

func GenerateToken() {
	// Logique pour générer un token
}

func VerifyToken() {
	// Logique pour vérifier un token
}
