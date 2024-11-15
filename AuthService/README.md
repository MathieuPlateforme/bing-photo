# Fichiers

**cmd/main.go** : Point d'entrée principal de l'application.

**Services/auth_service.go** : Contient la logique du service d'authentification, y compris loginWithEmail, registerWithEmail, forgotPassword, resetPassword, loginWithGoogle, et validateToken.

**Models/user/user.go** : Définit la structure de la classe User et ses méthodes associées (hashPassword, validatePassword, generatePasswordResetToken).

**Services/jwt/jwt.go** : Gère la logique de génération et de vérification des tokens JWT.

**Services/email_service.go** : Contient les fonctions pour envoyer des emails de vérification et de réinitialisation de mot de passe.

**Services/db_manager_service.go** : Gère la connexion à la base de données et l'exécution des requêtes.



# Dépendences

**JWT pour la gestion des tokens**
go get github.com/golang-jwt/jwt/v4

**Bcrypt pour le hashage des mots de passe**
go get golang.org/x/crypto/bcrypt

**Gomail pour l'envoi d'emails**
go get github.com/go-gomail/gomail

**Sqlx pour l'interaction avec la base de données**
go get github.com/jmoiron/sqlx

**Driver PostgreSQL** 
go get github.com/lib/pq

**OAuth2 pour l'authentification Google**
go get golang.org/x/oauth2
go get google.golang.org/api/oauth2/v2

**Installer GORM et le Pilote PostgreSQL**
go get gorm.io/gorm
go get gorm.io/driver/postgres


**Gorilla Mux pour le routage HTTP**
go get github.com/gorilla/mux
