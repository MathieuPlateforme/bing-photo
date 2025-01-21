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

	// Routes pour Albums
	router.HandleFunc("/albums", albumHandler.CreateAlbum).Methods("POST") 
	router.HandleFunc("/albums", albumHandler.GetAlbums).Methods("GET")
	router.HandleFunc("/albums/{id}", albumHandler.UpdateAlbum).Methods("PUT")


	return router
}
