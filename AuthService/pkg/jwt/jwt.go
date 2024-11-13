package jwt

type JWT struct {
	Token     string
	Expiration int64
	IssuedAt  int64
}

func GenerateToken() {
	// Logique pour générer un token
}

func VerifyToken() {
	// Logique pour vérifier un token
}
