package handlers

import (
	"GalleryService/internal/models"
	"GalleryService/internal/services"
	"encoding/json"
	"net/http"
)

type AlbumHandler struct {
	AlbumService *services.AlbumService
}

// NewAlbumHandler initialise un gestionnaire AlbumHandler
func NewAlbumHandler(albumService *services.AlbumService) *AlbumHandler {
	return &AlbumHandler{AlbumService: albumService}
}

// CreateAlbum gère la création d'un album
func (h *AlbumHandler) CreateAlbum(w http.ResponseWriter, r *http.Request) {
	var album models.Album
	if err := json.NewDecoder(r.Body).Decode(&album); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	if err := h.AlbumService.CreateAlbum(album); err != nil {
		http.Error(w, "Failed to create album: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "Album created successfully"})
}

func (h *AlbumHandler) GetAlbums(w http.ResponseWriter, r *http.Request) {
	// Appel du service pour récupérer les albums
	albums, err := h.AlbumService.ListAlbums()
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des albums : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Encoder la réponse en JSON
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(albums); err != nil {
		http.Error(w, "Erreur lors de l'encodage de la réponse : "+err.Error(), http.StatusInternalServerError)
	}
}