package api

import (
	"GalleryService/internal/api/handlers"
	"GalleryService/internal/db"
	"GalleryService/internal/services"

	"github.com/gorilla/mux"
)

func NewRouter(dbManager *db.DBManagerService) *mux.Router {
	router := mux.NewRouter()

	// Initialiser le client S3
	s3Service := services.NewS3Service("http://my-s3-clone:9090")

	// Initialiser le service AlbumService
	albumService := services.NewAlbumService(dbManager, s3Service)

	// Initialiser le gestionnaire AlbumHandler
	albumHandler := handlers.NewAlbumHandler(albumService)

	// Définir les routes
	router.HandleFunc("/albums", albumHandler.CreateAlbum).Methods("POST")

	return router
}
