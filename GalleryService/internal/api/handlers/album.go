package handlers

import (
	"GalleryService/internal/models"
	"GalleryService/internal/services"
	"encoding/json"
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"strconv"
	"GalleryService/internal/utils"
)

type AlbumHandler struct {
	AlbumService *services.AlbumService
	UserService  *services.UserService 
}

// NewAlbumHandler initialise un gestionnaire AlbumHandler
func NewAlbumHandler(albumService *services.AlbumService, userService *services.UserService) *AlbumHandler {
	return &AlbumHandler{AlbumService: albumService, UserService: userService}
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

func (h *AlbumHandler) GetAlbumsByUser(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID de l'utilisateur depuis les paramètres de la route
	vars := mux.Vars(r)
	userIDStr := vars["id"]

	// Convertir l'ID en uint
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "ID d'utilisateur invalide", http.StatusBadRequest)
		return
	}

	// Appel du service pour récupérer les albums par user
	albums, err := h.AlbumService.GetAlbumsByUser(uint(userID))
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

// UpdateAlbum gère la mise à jour d'un album
func (h *AlbumHandler) UpdateAlbum(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID depuis les paramètres de l'URL
	vars := mux.Vars(r)
	albumID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	// Lire et décoder le corps de la requête
	var updateData struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		http.Error(w, "Requête invalide : "+err.Error(), http.StatusBadRequest)
		return
	}

	// Appeler le service pour mettre à jour l'album
	err = h.AlbumService.UpdateAlbum(uint(albumID), updateData.Name, updateData.Description)
	if err != nil {
		http.Error(w, "Erreur lors de la mise à jour de l'album : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Réponse en cas de succès
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Album %d mis à jour avec succès", albumID)
}

// DeleteAlbum gère la suppression d'un album
func (h *AlbumHandler) DeleteAlbum(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID de l'album depuis l'URL
	vars := mux.Vars(r)
	albumID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID invalide", http.StatusBadRequest)
		return
	}

	// Appeler le service pour supprimer l'album
	err = h.AlbumService.DeleteAlbum(uint(albumID))
	if err != nil {
		http.Error(w, "Erreur lors de la suppression de l'album : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Répondre avec un statut de succès
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Album supprimé avec succès"))
}

func (h *AlbumHandler) GetPrivateAlbum(w http.ResponseWriter, r *http.Request) {
    userID, err := utils.GetUserIDFromContext(r.Context())
    if err != nil {
        http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
        return
    }

    album, err := h.AlbumService.GetPrivateAlbum(userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(album)
}
