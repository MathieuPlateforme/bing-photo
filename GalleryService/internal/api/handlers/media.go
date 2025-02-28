package handlers

import (
	"GalleryService/internal/models"
	"GalleryService/internal/services"
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"GalleryService/internal/utils"
	"fmt"
	"log"
)

type MediaHandler struct {
	MediaService *services.MediaService
    UserService  *services.UserService 
}

// NewMediaHandler initialise un gestionnaire MediaHandler
func NewMediaHandler(mediaService *services.MediaService, userService *services.UserService) *MediaHandler {
	return &MediaHandler{MediaService: mediaService, UserService: userService}
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

func (h *MediaHandler) MarkAsPrivate(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
		return
	}
		
	mediaID, _ := strconv.Atoi(mux.Vars(r)["id"])

	var request struct {
		Pin string `json:"pin"`
	}
	json.NewDecoder(r.Body).Decode(&request)

	err = h.UserService.SetPrivateAlbumPin(uint(userID), request.Pin)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = h.MediaService.MarkAsPrivate(uint(userID), uint(mediaID))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Photo marquée comme privée"})
}

func (h *MediaHandler) GetPrivateMedia(w http.ResponseWriter, r *http.Request) {
	userID, err := utils.GetUserIDFromContext(r.Context())
	if err != nil {
		http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
		return
	}
	
	pin := r.URL.Query().Get("pin")

	err = h.UserService.VerifyPrivateAlbumPin(uint(userID), pin)
	if err != nil {
		http.Error(w, "PIN incorrect", http.StatusUnauthorized)
		return
	}

	media, err := h.MediaService.GetPrivateMedia(uint(userID))
	if err != nil {
		http.Error(w, "Erreur lors de la récupération des photos privées", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(media)
}

func (h *MediaHandler) MarkMediaAsPrivate(w http.ResponseWriter, r *http.Request) {
    userID, err := utils.GetUserIDFromContext(r.Context())
    if err != nil {
        http.Error(w, "Utilisateur non authentifié", http.StatusUnauthorized)
        return
    }

    var request struct {
        MediaID uint `json:"mediaID"`
    }

    if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
        http.Error(w, "Données de requête invalides", http.StatusBadRequest)
        return
    }

    err = h.MediaService.MarkAsPrivate(request.MediaID, userID)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    json.NewEncoder(w).Encode(map[string]string{
        "message": "Média marqué comme privé avec succès",
    })
}

func (h *MediaHandler) DownloadMedia(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mediaID, err := strconv.ParseUint(vars["id"], 10, 64)
	if err != nil {
		http.Error(w, "Invalid media ID", http.StatusBadRequest)
		log.Printf("Invalid media ID: %v", err)
		return
	}

	// Télécharger le média
	if err := h.MediaService.DownloadMedia(uint(mediaID), w); err != nil {
		http.Error(w, "Error downloading media", http.StatusInternalServerError)
		log.Printf("Error downloading media: %v", err)
		return
	}

	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", vars["id"]))
	w.WriteHeader(http.StatusOK)
	log.Printf("Media %d downloaded successfully", mediaID)
}

func (h *MediaHandler) DeleteMedia(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    mediaID, err := strconv.ParseUint(vars["id"], 10, 64)
    if err != nil {
        http.Error(w, "Invalid media ID", http.StatusBadRequest)
        return
    }

    // Récupérer le userID à partir du contexte
    userID, err := utils.GetUserIDFromContext(r.Context())
    if err != nil {
        http.Error(w, "Unauthorized: userID missing in context", http.StatusUnauthorized)
        return
    }

    // Supprimer le média
    if err := h.MediaService.DeleteMedia(uint(mediaID), userID); err != nil {
        log.Printf("Error deleting media: %v", err)
        http.Error(w, fmt.Sprintf("Failed to delete media: %v", err), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusNoContent)
}

func (h *MediaHandler) DetectSimilarMedia(w http.ResponseWriter, r *http.Request) {
    userID, err := utils.GetUserIDFromContext(r.Context())
    if err != nil {
        log.Printf("Error retrieving userID from context: %v", err)
        http.Error(w, "User not authenticated", http.StatusUnauthorized)
        return
    }

    vars := mux.Vars(r)
    albumIDStr, exists := vars["albumID"]
    if !exists {
        http.Error(w, "albumID required in the URL", http.StatusBadRequest)
        return
    }

    albumID, err := strconv.ParseUint(albumIDStr, 10, 32)
    if err != nil {
        http.Error(w, "Invalid albumID", http.StatusBadRequest)
        return
    }

    log.Printf("Detecting similar media for userID=%d in albumID=%d", userID, albumID)

    similarMedia, err := h.MediaService.DetectSimilarMedia(uint(userID), uint(albumID))
    if err != nil {
        log.Printf("Error detecting similar media: %v", err)
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    response, err := json.Marshal(similarMedia)
    if err != nil {
        log.Printf("Error generating JSON response: %v", err)
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)
    w.Write(response)

    log.Printf("Successfully completed similar media detection for userID=%d and albumID=%d", userID, albumID)
}


