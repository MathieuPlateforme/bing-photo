package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	proto "ApiGateway/proto"
	"google.golang.org/grpc/metadata"
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
// @Failure 500 {string} string "Erreur serveur"
// @Router /albums/{id} [put]
// @Security BearerAuth
func (g *GalleryGateway) UpdateAlbumHandler(w http.ResponseWriter, r *http.Request) {
	var req proto.UpdateAlbumRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Failed to parse request: %v\n", err)
		return
	}

	res, err := g.GalleryClient.UpdateAlbum(context.Background(), &req)
	if err != nil {
		http.Error(w, "Failed to update album: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Update album error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// @Summary Supprimer un album
// @Description Supprime un album existant par son ID
// @Tags Albums
// @Produce json
// @Param id path int true "ID de l'album"
// @Success 200 {object} proto.DeleteAlbumResponse
// @Failure 400 {string} string "ID invalide"
// @Failure 500 {string} string "Erreur serveur"
// @Router /albums/{id} [delete]
// @Security BearerAuth
func (g *GalleryGateway) DeleteAlbumHandler(w http.ResponseWriter, r *http.Request) {
	albumID, err := strconv.ParseUint(r.URL.Query().Get("album_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid album ID", http.StatusBadRequest)
		return
	}

	req := &proto.DeleteAlbumRequest{
		AlbumId: uint32(albumID),
	}

	res, err := g.GalleryClient.DeleteAlbum(context.Background(), req)
	if err != nil {
		http.Error(w, "Failed to delete album: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Delete album error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// @Summary Obtenir un album privé
// @Description Récupère un album privé à partir de son ID
// @Tags Albums
// @Produce json
// @Param album_id query int true "ID de l'album"
// @Success 200 {object} proto.GetPrivateAlbumResponse
// @Failure 400 {string} string "ID invalide"
// @Failure 500 {string} string "Erreur serveur"
// @Router /albums/private [get]
// @Security BearerAuth
func (g *GalleryGateway) GetPrivateAlbumHandler(w http.ResponseWriter, r *http.Request) {
	albumID, err := strconv.ParseUint(r.URL.Query().Get("album_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid album ID", http.StatusBadRequest)
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

	// Construire la requête
	req := &proto.GetPrivateAlbumRequest{
		AlbumId: uint32(albumID),
	}

	// Appeler le service
	res, err := g.GalleryClient.GetPrivateAlbum(ctx, req)
	if err != nil {
		http.Error(w, "Failed to get private album: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Get private album error: %v\n", err)
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

	res, err := g.MediaClient.AddMedia(context.Background(), req)
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

// @Summary Rendre un média privé
// @Description Marque un média comme privé
// @Tags Media
// @Accept json
// @Produce json
// @Param request body proto.MarkAsPrivateRequest true "Données de la requête"
// @Success 200 {object} proto.MarkAsPrivateResponse
// @Failure 400 {string} string "Requête invalide"
// @Failure 500 {string} string "Erreur serveur"
// @Router /media/{id}/private [post]
// @Security BearerAuth
func (g *GalleryGateway) MarkAsPrivateHandler(w http.ResponseWriter, r *http.Request) {
	var req proto.MarkAsPrivateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Failed to parse request: %v\n", err)
		return
	}

	res, err := g.MediaClient.MarkAsPrivate(context.Background(), &req)
	if err != nil {
		http.Error(w, "Failed to mark as private: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Mark as private error: %v\n", err)
		return
	}

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
	mediaID, err := strconv.ParseUint(r.URL.Query().Get("media_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid media ID", http.StatusBadRequest)
		return
	}

	req := &proto.DownloadMediaRequest{
		MediaId: uint32(mediaID),
	}

	res, err := g.MediaClient.DownloadMedia(context.Background(), req)
	if err != nil {
		http.Error(w, "Failed to download media: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Download media error: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=media")
	w.WriteHeader(http.StatusOK)
	w.Write(res.FileData)
}

// @Summary Supprimer un média
// @Description Supprime un média par son ID
// @Tags Media
// @Produce json
// @Param id path int true "ID du média"
// @Success 200 {object} proto.DeleteMediaResponse
// @Failure 400 {string} string "ID invalide"
// @Failure 500 {string} string "Erreur serveur"
// @Router /media/{id} [delete]
// @Security BearerAuth
func (g *GalleryGateway) DeleteMediaHandler(w http.ResponseWriter, r *http.Request) {
	mediaID, err := strconv.ParseUint(r.URL.Query().Get("media_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid media ID", http.StatusBadRequest)
		return
	}

	req := &proto.DeleteMediaRequest{
		MediaId: uint32(mediaID),
	}

	res, err := g.MediaClient.DeleteMedia(context.Background(), req)
	if err != nil {
		http.Error(w, "Failed to delete media: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Delete media error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// @Summary Détecter les médias similaires
// @Description Détecte les fichiers médias similaires à celui donné
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
	mediaID, err := strconv.ParseUint(r.URL.Query().Get("media_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid media ID", http.StatusBadRequest)
		return
	}

	req := &proto.DetectSimilarMediaRequest{
		MediaId: uint32(mediaID),
	}

	res, err := g.MediaClient.DetectSimilarMedia(context.Background(), req)
	if err != nil {
		http.Error(w, "Failed to detect similar media: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Detect similar media error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
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
