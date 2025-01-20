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
