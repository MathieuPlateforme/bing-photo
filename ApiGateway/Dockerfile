# Étape de construction
FROM golang:1.23 AS builder

WORKDIR /app

# Copier les fichiers de dépendances
COPY go.mod go.sum ./
RUN go mod download

# Copier le reste des fichiers du projet
COPY . .

# Construire l'application avec CGO désactivé (nécessaire pour Alpine)
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./main.go

# Vérifier que l'exécutable est bien créé
RUN ls -la /app

# Étape finale pour exécuter l'application
FROM alpine:latest
WORKDIR /root/

# Copier l'exécutable depuis l'étape de construction
COPY --from=builder /app/main .

# Lister les fichiers pour vérifier la présence de `main`
RUN ls -la

# Exposer le port
EXPOSE 8080

# Commande par défaut pour exécuter l'application
CMD ["./main"]
