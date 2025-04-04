package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"fmt"
	proto "ApiGateway/proto"

	"google.golang.org/grpc/metadata"
	"github.com/gorilla/mux"
)

type GalleryGateway struct {
	GalleryClient proto.AlbumServiceClient
	MediaClient   proto.MediaServiceClient
	UserClient    proto.UserServiceClient
}

func NewGalleryGateway(albumClient proto.AlbumServiceClient, mediaClient proto.MediaServiceClient, userClient proto.UserServiceClient) *GalleryGateway {
	return &GalleryGateway{
		GalleryClient: albumClient,
		MediaClient:   mediaClient,
		UserClient:    userClient,
	}
}

// @Summary Créer un album
// @Description Crée un nouvel album avec un titre et un identifiant utilisateur
// @Tags Albums
// @Accept json
// @Produce json
// @Param album body proto.CreateAlbumRequest true "Données de l'album"
// @Success 201 {object} proto.CreateAlbumResponse
// @Failure 400 {string} string "Requête invalide"
// @Failure 500 {string} string "Erreur serveur"
// @Router /albums [post]
// @Security BearerAuth
func (g *GalleryGateway) CreateAlbumHandler(w http.ResponseWriter, r *http.Request) {
	var req proto.CreateAlbumRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Failed to parse request: %v\n", err)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		log.Println("Authorization header missing")
		return
	}

	md := metadata.New(map[string]string{"authorization": authHeader})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Appel gRPC avec le contexte enrichi
	res, err := g.GalleryClient.CreateAlbum(ctx, &req)
	if err != nil {
		http.Error(w, "Failed to create album: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Create album error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Failed to encode response: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// @Summary Obtenir les albums d'un utilisateur
// @Description Récupère tous les albums appartenant à un utilisateur donné
// @Tags Albums
// @Produce json
// @Param user_id query int true "ID utilisateur"
// @Success 200 {object} proto.GetAlbumsByUserResponse
// @Failure 400 {string} string "ID utilisateur invalide"
// @Failure 500 {string} string "Erreur serveur"
// @Router /albums/user [get]
// @Security BearerAuth
func (g *GalleryGateway) GetAlbumsByUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(r.URL.Query().Get("user_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	req := &proto.GetAlbumsByUserRequest{
		UserId: uint32(userID),
	}

	res, err := g.GalleryClient.GetAlbumsByUser(context.Background(), req)
	if err != nil {
		http.Error(w, "Failed to get albums: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Get albums error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// @Summary Mettre à jour un album
// @Description Met à jour les informations d'un album
// @Tags Albums
// @Accept json
// @Produce json
// @Param id path int true "ID de l'album"
// @Param album body proto.UpdateAlbumRequest true "Mise à jour de l'album"
// @Success 200 {object} proto.UpdateAlbumResponse
// @Failure 400 {string} string "Requête invalide"
// @Failure 401 {string} string "Non autorisé"
// @Failure 500 {string} string "Erreur serveur"
// @Router /albums/{id} [put]
// @Security BearerAuth
func (g *GalleryGateway) UpdateAlbumHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifier la présence de l'en-tête Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		log.Println("Authorization header missing")
		return
	}

	// Décoder le corps de la requête
	var req proto.UpdateAlbumRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Failed to parse request: %v\n", err)
		return
	}

	// Injecter l'en-tête dans le contexte gRPC
	md := metadata.New(map[string]string{"authorization": authHeader})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Appel gRPC
	res, err := g.GalleryClient.UpdateAlbum(ctx, &req)
	if err != nil {
		http.Error(w, "Failed to update album: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Update album error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// @Summary Supprimer un album
// @Description Supprime un album par son ID
// @Tags Albums
// @Produce json
// @Param id path int true "ID de l'album"
// @Success 200 {object} proto.DeleteAlbumResponse
// @Failure 400 {string} string "Invalid album ID"
// @Failure 401 {string} string "Authorization header missing"
// @Failure 500 {string} string "Failed to delete album"
// @Router /albums/{id} [delete]
// @Security BearerAuth
func (g *GalleryGateway) DeleteAlbumHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		log.Println("Authorization header missing")
		return
	}

	vars := mux.Vars(r)
	albumIDStr := vars["id"]

	albumID, err := strconv.ParseUint(albumIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid album ID", http.StatusBadRequest)
		log.Printf("Invalid album ID: %v\n", err)
		return
	}

	req := &proto.DeleteAlbumRequest{
		AlbumId: uint32(albumID),
	}

	md := metadata.New(map[string]string{"authorization": authHeader})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	res, err := g.GalleryClient.DeleteAlbum(ctx, req)
	if err != nil {
		http.Error(w, "Failed to delete album: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Delete album error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// @Summary Obtenir un album par type
// @Description Récupère un album privé ou principal selon le type
// @Tags Albums
// @Produce json
// @Param user_id query int true "ID de l'utilisateur"
// @Param type query string true "Type d'album : 'private' ou 'main'"
// @Success 200 {object} proto.GetPrivateAlbumResponse
// @Failure 400 {string} string "ID ou type invalide"
// @Failure 401 {string} string "Authorization manquante"
// @Failure 500 {string} string "Erreur serveur"
// @Router /albums/type [get]
// @Security BearerAuth
func (g *GalleryGateway) GetPrivateAlbumHandler(w http.ResponseWriter, r *http.Request) {
	userId, err := strconv.ParseUint(r.URL.Query().Get("user_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	albumType := r.URL.Query().Get("type")
	if albumType != "private" && albumType != "main" {
		http.Error(w, "Invalid type, must be 'private' or 'main'", http.StatusBadRequest)
		return
	}

	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		log.Println("Authorization header missing")
		return
	}

	md := metadata.New(map[string]string{"authorization": authHeader})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Utilise toujours GetPrivateAlbum côté proto, mais passe le type dynamiquement
	req := &proto.GetPrivateAlbumRequest{
		UserId: uint32(userId),
		Type:   albumType,
	}

	res, err := g.GalleryClient.GetPrivateAlbum(ctx, req)
	if err != nil {
		http.Error(w, "Failed to get album: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Get album error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Failed to encode response: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// @Summary Ajouter un média
// @Description Ajoute un fichier média à un album
// @Tags Media
// @Accept multipart/form-data
// @Produce json
// @Param album_id formData int true "ID de l'album"
// @Param file formData file true "Fichier à uploader"
// @Success 201 {object} proto.AddMediaResponse
// @Failure 400 {string} string "Erreur de parsing du formulaire"
// @Failure 500 {string} string "Erreur serveur"
// @Router /media [post]
// @Security BearerAuth
func (g *GalleryGateway) AddMediaHandler(w http.ResponseWriter, r *http.Request) {
	// Vérifie l'en-tête Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		log.Println("Authorization header missing")
		return
	}

	// Parse le formulaire (fichier + album_id)
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	albumID, err := strconv.ParseUint(r.FormValue("album_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid album ID", http.StatusBadRequest)
		return
	}

	req := &proto.AddMediaRequest{
		Name:     header.Filename,
		AlbumId:  uint32(albumID),
		FileData: fileData,
	}

	// Injecter l'en-tête Authorization dans le contexte
	md := metadata.New(map[string]string{"authorization": authHeader})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Appel gRPC
	res, err := g.MediaClient.AddMedia(ctx, req)
	if err != nil {
		http.Error(w, "Failed to add media: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Add media error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}


// @Summary Médias d’un utilisateur
// @Description Récupère tous les médias appartenant à un utilisateur
// @Tags Media
// @Produce json
// @Param user_id query int true "ID utilisateur"
// @Success 200 {object} proto.GetMediaByUserResponse
// @Failure 400 {string} string "ID invalide"
// @Failure 500 {string} string "Erreur serveur"
// @Router /media/user [get]
// @Security BearerAuth
func (g *GalleryGateway) GetMediaByUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(r.URL.Query().Get("user_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	req := &proto.GetMediaByUserRequest{
		UserId: uint32(userID),
	}

	res, err := g.MediaClient.GetMediaByUser(context.Background(), req)
	if err != nil {
		http.Error(w, "Failed to get media: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Get media error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}


func (g *GalleryGateway) MarkAsPrivateHandler(w http.ResponseWriter, r *http.Request) {
	// Extraire le token du header Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		log.Println("Authorization header is missing")
		return
	}

	// Injecter le token dans un contexte gRPC
	ctx := metadata.AppendToOutgoingContext(context.Background(), "authorization", authHeader)

	// Lire et parser le corps de la requête
	var req proto.MarkAsPrivateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Failed to parse request: %v\n", err)
		return
	}

	// Appeler le service gRPC avec le contexte enrichi
	res, err := g.MediaClient.MarkAsPrivate(ctx, &req)
	if err != nil {
		http.Error(w, "Failed to mark as private: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Mark as private error: %v\n", err)
		return
	}

	// Réponse OK
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// @Summary Obtenir les médias privés
// @Description Récupère les médias privés d’un utilisateur
// @Tags Media
// @Produce json
// @Param user_id query int true "ID utilisateur"
// @Success 200 {object} proto.GetPrivateMediaResponse
// @Failure 400 {string} string "ID invalide"
// @Failure 500 {string} string "Erreur serveur"
// @Router /media/private [get]
// @Security BearerAuth
func (g *GalleryGateway) GetPrivateMediaHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(r.URL.Query().Get("user_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	req := &proto.GetPrivateMediaRequest{
		UserId: uint32(userID),
		Pin:    r.URL.Query().Get("pin"),
	}

	res, err := g.MediaClient.GetPrivateMedia(context.Background(), req)
	if err != nil {
		http.Error(w, "Failed to get private media: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Get private media error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// @Summary Télécharger un média
// @Description Télécharge le contenu d’un fichier média
// @Tags Media
// @Produce application/octet-stream
// @Param id path int true "ID du média"
// @Success 200 {file} file "Fichier binaire"
// @Failure 400 {string} string "ID invalide"
// @Failure 500 {string} string "Erreur serveur"
// @Router /media/{id}/download [get]
// @Security BearerAuth
func (g *GalleryGateway) DownloadMediaHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		log.Println("Authorization header missing")
		return
	}

	vars := mux.Vars(r)
	mediaIDStr := vars["id"]

	mediaID, err := strconv.ParseUint(mediaIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid media ID", http.StatusBadRequest)
		log.Printf("Invalid media ID: %v\n", err)
		return
	}

	req := &proto.DownloadMediaRequest{
		MediaId: uint32(mediaID),
	}

	md := metadata.New(map[string]string{"authorization": authHeader})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	res, err := g.MediaClient.DownloadMedia(ctx, req)
	if err != nil {
		http.Error(w, "Failed to download media: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Download media error: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", http.DetectContentType(res.FileData))
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=media_%d", mediaID))
	w.WriteHeader(http.StatusOK)
	w.Write(res.FileData)
}

// DeleteMediaHandler supprime un média spécifique
// @Summary Supprimer un média
// @Description Supprime un média si l'utilisateur est propriétaire
// @Tags Media
// @Produce json
// @Param id path int true "ID du média à supprimer"
// @Success 200 {object} proto.DeleteMediaResponse
// @Failure 400 {string} string "Invalid media ID"
// @Failure 401 {string} string "Authorization header missing"
// @Failure 500 {string} string "Failed to delete media"
// @Security BearerAuth
// @Router /media/{id} [delete]
func (g *GalleryGateway) DeleteMediaHandler(w http.ResponseWriter, r *http.Request) {
	// Récupération du media ID depuis l'URL
	vars := mux.Vars(r)
	mediaIDStr := vars["id"]

	mediaID, err := strconv.ParseUint(mediaIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid media ID", http.StatusBadRequest)
		return
	}

	// Récupération du token JWT
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		log.Println("Authorization header missing")
		return
	}
	log.Printf("Authorization header: %s", authHeader)

	// Ajout du header au contexte gRPC
	md := metadata.New(map[string]string{"authorization": authHeader})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Création de la requête gRPC
	req := &proto.DeleteMediaRequest{
		MediaId: uint32(mediaID),
	}

	res, err := g.MediaClient.DeleteMedia(ctx, req)
	if err != nil {
		http.Error(w, "Failed to delete media: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Delete media error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// @Summary Détecter les médias similaires dans un album
// @Description Détecte les fichiers médias similaires dans un album donné
// @Tags Media
// @Accept json
// @Produce json
// @Param request body proto.DetectSimilarMediaRequest true "Requête de détection"
// @Success 200 {object} proto.DetectSimilarMediaResponse
// @Failure 400 {string} string "Requête invalide"
// @Failure 500 {string} string "Erreur serveur"
// @Router /media/similar [post]
// @Security BearerAuth
func (g *GalleryGateway) DetectSimilarMediaHandler(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		log.Println("Authorization header missing")
		return
	}

	var req proto.DetectSimilarMediaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Failed to decode request: %v\n", err)
		return
	}

	if req.AlbumId == 0 {
		http.Error(w, "Missing album_id", http.StatusBadRequest)
		log.Println("album_id manquant")
		return
	}

	md := metadata.New(map[string]string{"authorization": authHeader})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	res, err := g.MediaClient.DetectSimilarMedia(ctx, &req)
	if err != nil {
		http.Error(w, "Failed to detect similar media: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Detect similar media error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		log.Printf("Failed to encode response: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}



func (g *GalleryGateway) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req proto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Failed to parse request: %v\n", err)
		return
	}

	res, err := g.UserClient.CreateUser(context.Background(), &req)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Create user error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (g *GalleryGateway) AddMediaToFavoriteHandler(w http.ResponseWriter, r *http.Request) {
	var req proto.AddMediaToFavoriteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Failed to parse request: %v\n", err)
		return
	}

	res, err := g.MediaClient.AddMediaToFavorite(context.Background(), &req)
	if err != nil {
		http.Error(w, "Failed to add media to favorite: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Add media to favorite error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// GetMediaByAlbumHandler godoc
// @Summary Récupérer les médias d’un album
// @Description Renvoie tous les médias appartenant à un album donné
// @Tags Media
// @Produce json
// @Param id path int true "ID de l’album"
// @Success 200 {object} proto.GetMediaByAlbumResponse
// @Failure 400 {string} string "Requête invalide"
// @Failure 500 {string} string "Erreur serveur"
// @Router /media/album/{id} [get]
// @Security BearerAuth
func (g *GalleryGateway) GetMediaByAlbumHandler(w http.ResponseWriter, r *http.Request) {
	// Extraire l'Authorization
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header missing", http.StatusUnauthorized)
		log.Println("Authorization header missing")
		return
	}

	// Extraire l'ID depuis l'URL
	vars := mux.Vars(r)
	albumIDStr := vars["id"]
	albumID, err := strconv.ParseUint(albumIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Invalid album ID", http.StatusBadRequest)
		return
	}

	// Contexte avec JWT
	md := metadata.New(map[string]string{"authorization": authHeader})
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	// Requête gRPC
	req := &proto.GetMediaByAlbumRequest{AlbumId: uint32(albumID)}
	res, err := g.MediaClient.GetMediaByAlbum(ctx, req)
	if err != nil {
		http.Error(w, "Failed to fetch media: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Get media by album error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}
