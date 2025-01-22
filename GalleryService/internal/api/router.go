package api

import (
	"GalleryService/internal/api/handlers"
	"GalleryService/internal/db"
	"GalleryService/internal/services"

	"github.com/gorilla/mux"
)

func NewRouter(dbManager *db.DBManagerService, s3Service *services.S3Service) *mux.Router {
	router := mux.NewRouter()
	// Initialiser le service AlbumService
	albumService := services.NewAlbumService(dbManager, s3Service)

	// Initialiser le gestionnaire AlbumHandler
	albumHandler := handlers.NewAlbumHandler(albumService)

	// Initialiser le service MediaService
	mediaService := services.NewMediaService(dbManager, s3Service)

	// Initialiser le gestionnaire MediaHandler
	mediaHandler := handlers.NewMediaHandler(mediaService)

	// Routes pour Albums
	router.HandleFunc("/albums", albumHandler.CreateAlbum).Methods("POST") 
	router.HandleFunc("/users/{id}/albums", albumHandler.GetAlbumsByUser).Methods("GET")
	router.HandleFunc("/albums/{id}", albumHandler.UpdateAlbum).Methods("PUT")
	router.HandleFunc("/albums/{id}", albumHandler.DeleteAlbum).Methods("DELETE")

	// Route pour Medias
	router.HandleFunc("/media", mediaHandler.AddMedia).Methods("POST")

	return router
}
