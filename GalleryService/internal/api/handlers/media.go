package handlers

import (
	"GalleryService/internal/models"
	"GalleryService/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
)

type MediaHandler struct {
	MediaService *services.MediaService
}

// NewMediaHandler initialise un gestionnaire MediaHandler
func NewMediaHandler(mediaService *services.MediaService) *MediaHandler {
	return &MediaHandler{MediaService: mediaService}
}

// AddMedia gère l'ajout d'un fichier
func (h *MediaHandler) AddMedia(w http.ResponseWriter, r *http.Request) {
	// Analyse des données multipart/form-data
	err := r.ParseMultipartForm(10 << 20) // Limite de 10 MB
	if err != nil {
		http.Error(w, "Erreur lors de l'analyse du formulaire : "+err.Error(), http.StatusBadRequest)
		return
	}

	// Récupérer le fichier et ses métadonnées
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Fichier manquant dans la requête : "+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	albumIDStr := r.FormValue("albumID")
	albumID, err := strconv.Atoi(albumIDStr)
	if err != nil {
		http.Error(w, "ID d'album invalide", http.StatusBadRequest)
		return
	}

	// Créer une instance de Media
	media := models.Media{
		Name:     fileHeader.Filename,
		AlbumID:  uint(albumID),
		FileSize: uint(fileHeader.Size),
	}

	// Appeler le service pour ajouter le fichier
	err = h.MediaService.AddMedia(&media, file, fileHeader.Size)
	if err != nil {
		http.Error(w, "Erreur lors de l'ajout du fichier : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Répondre avec succès
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Fichier ajouté avec succès",
	})
}

func (h *MediaHandler) GetMediaByUser(w http.ResponseWriter, r *http.Request) {
	// Récupérer l'ID de l'utilisateur depuis les paramètres de la route
	vars := mux.Vars(r)
	userIDStr := vars["id"]

	// Convertir l'ID en uint
	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "ID d'utilisateur invalide", http.StatusBadRequest)
		return
	}

	// Appeler le service pour récupérer les médias
	mediaList, err := h.MediaService.GetMediaByUser(uint(userID))
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des fichiers : "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Répondre avec les médias sous forme JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(mediaList)
}

