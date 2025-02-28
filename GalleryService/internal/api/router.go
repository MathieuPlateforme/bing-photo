package api

import (
	"GalleryService/internal/api/handlers"
	"GalleryService/internal/db"
	"GalleryService/internal/services"
	"GalleryService/internal/middleware"

	"github.com/gorilla/mux"
)

func NewRouter(dbManager *db.DBManagerService, s3Service *services.S3Service, authClientService *services.AuthServiceClient) *mux.Router {

	// JWT_SECRET := os.Getenv("JWT_SECRET")
	router := mux.NewRouter()

	// Middleware d'authentification
	router.Use(middleware.AuthMiddleware(authClientService, "secretkey" ))

	// Initialiser le service UserService
	userService := services.NewUserService(dbManager, s3Service)

	// Initialiser le gestionnaire UserHandler
	userHandler := handlers.NewUserHandler(userService)

	// Initialiser le service AlbumService
	albumService := services.NewAlbumService(dbManager, s3Service)

	// Initialiser le gestionnaire AlbumHandler
	albumHandler := handlers.NewAlbumHandler(albumService, userService)

	// Initialiser le service MediaService
	mediaService := services.NewMediaService(dbManager, s3Service)

	// Initialiser le gestionnaire MediaHandler
	mediaHandler := handlers.NewMediaHandler(mediaService, userService)

	// Routes pour Albums
	router.HandleFunc("/albums", albumHandler.CreateAlbum).Methods("POST") 
	router.HandleFunc("/users/{id}/albums", albumHandler.GetAlbumsByUser).Methods("GET")
	router.HandleFunc("/albums/{id}", albumHandler.UpdateAlbum).Methods("PUT")
	router.HandleFunc("/albums/{id}", albumHandler.DeleteAlbum).Methods("DELETE")

	// Routes pour Medias
	router.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	router.HandleFunc("/media", mediaHandler.AddMedia).Methods("POST")
	router.HandleFunc("/users/{id}/media", mediaHandler.GetMediaByUser).Methods("GET")
	router.HandleFunc("/media/{id}/private", mediaHandler.MarkAsPrivate).Methods("PUT")
	router.HandleFunc("/media/private", mediaHandler.GetPrivateMedia).Methods("GET")
	router.HandleFunc("/media/{id}", mediaHandler.DownloadMedia).Methods("GET")
	router.HandleFunc("/media/{id}", mediaHandler.DeleteMedia).Methods("DELETE")
	router.HandleFunc("/{albumID}/media/similar", mediaHandler.DetectSimilarMedia).Methods("POST")

	return router
}
