# Étape 1 : Construction de l'application Go
FROM golang:1.21 as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o my-s3-clone

# Étape 2 : Image légère pour exécuter l'application
FROM debian:bookworm-slim

COPY --from=builder /app/my-s3-clone /my-s3-clone

EXPOSE 9090

CMD ["/my-s3-clone"]