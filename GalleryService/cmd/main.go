package main

import (
	"bytes"
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/codes"


	"GalleryService/internal/db"
	"GalleryService/internal/jwt"
	"GalleryService/internal/middleware"
	"GalleryService/internal/models"
	proto "GalleryService/internal/proto"
	"GalleryService/internal/services"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

type galleryServer struct {
	proto.UnimplementedAlbumServiceServer
	proto.UnimplementedMediaServiceServer
	proto.UnimplementedUserServiceServer
	albumService *services.AlbumService
	mediaService *services.MediaService
	userService  *services.UserService
}

// Album Service methods
func (s *galleryServer) CreateAlbum(ctx context.Context, req *proto.CreateAlbumRequest) (*proto.CreateAlbumResponse, error) {
	album := models.Album{
		Name:        req.Name,
		Description: req.Description,
		UserID:      uint(req.UserId),
	}

	if err := s.albumService.CreateAlbum(album); err != nil {
		log.Printf("Error creating album: %v", err)
		return nil, err
	}

	return &proto.CreateAlbumResponse{
		Message: "Album created successfully",
	}, nil
}

func (s *galleryServer) GetAlbumsByUser(ctx context.Context, req *proto.GetAlbumsByUserRequest) (*proto.GetAlbumsByUserResponse, error) {
	albums, err := s.albumService.GetAlbumsByUser(uint(req.UserId))
	if err != nil {
		log.Printf("Error getting albums by user: %v", err)
		return nil, err
	}

	var protoAlbums []*proto.AlbumWithMedia
	for _, album := range albums {
		var protoMedia []*proto.Media
		for _, media := range album.Media {
			protoMedia = append(protoMedia, &proto.Media{
				Id:       uint32(media.ID),
				Name:     media.Name,
				Path:     media.Path,
				FileSize: uint32(media.FileSize),
				AlbumId:  uint32(media.AlbumID),
			})
		}

		protoAlbums = append(protoAlbums, &proto.AlbumWithMedia{
			Id:          uint32(album.ID),
			Name:        album.Name,
			Description: album.Description,
			UserId:      uint32(album.UserID),
			Media:       protoMedia,
		})
	}

	return &proto.GetAlbumsByUserResponse{
		Albums: protoAlbums,
	}, nil
}



func (s *galleryServer) UpdateAlbum(ctx context.Context, req *proto.UpdateAlbumRequest) (*proto.UpdateAlbumResponse, error) {
	if err := s.albumService.UpdateAlbum(uint(req.AlbumId), req.Name, req.Description); err != nil {
		log.Printf("Error updating album: %v", err)
		return nil, err
	}

	return &proto.UpdateAlbumResponse{}, nil
}

func (s *galleryServer) DeleteAlbum(ctx context.Context, req *proto.DeleteAlbumRequest) (*proto.DeleteAlbumResponse, error) {
	if err := s.albumService.DeleteAlbum(uint(req.AlbumId)); err != nil {
		log.Printf("Error deleting album: %v", err)
		return nil, err
	}

	return &proto.DeleteAlbumResponse{}, nil
}

func (s *galleryServer) GetPrivateAlbum(ctx context.Context, req *proto.GetPrivateAlbumRequest) (*proto.GetPrivateAlbumResponse, error) {
	album, err := s.albumService.GetPrivateAlbum(uint(req.UserId), req.Type)
	if err != nil {
		log.Printf("Error getting private album: %v", err)
		return nil, err
	}

	return &proto.GetPrivateAlbumResponse{
		Album: &proto.Album{
			Id:          uint32(album.ID),
			Name:        album.Name,
			Description: album.Description,
			UserId:      uint32(album.UserID),
		},
	}, nil
}

// Media Service methods
func (s *galleryServer) AddMedia(ctx context.Context, req *proto.AddMediaRequest) (*proto.AddMediaResponse, error) {
	media := &models.Media{
		Name:    req.Name,
		AlbumID: uint(req.AlbumId),
	}

	reader := bytes.NewReader(req.FileData)
	if err := s.mediaService.AddMedia(media, reader, int64(len(req.FileData))); err != nil {
		log.Printf("Error adding media: %v", err)
		return nil, err
	}

	return &proto.AddMediaResponse{
		Message: "Media added successfully",
	}, nil
}

func (s *galleryServer) GetMediaByUser(ctx context.Context, req *proto.GetMediaByUserRequest) (*proto.GetMediaByUserResponse, error) {
	media, err := s.mediaService.GetMediaByUser(uint(req.UserId))
	if err != nil {
		log.Printf("Error getting media by user: %v", err)
		return nil, err
	}

	var protoMedia []*proto.Media
	for _, m := range media {
		protoMedia = append(protoMedia, &proto.Media{
			Id:       uint32(m.ID),
			Name:     m.Name,
			AlbumId:  uint32(m.AlbumID),
			FileSize: uint32(m.FileSize),
		})
	}

	return &proto.GetMediaByUserResponse{
		MediaList: protoMedia,
	}, nil
}

func (s *galleryServer) MarkAsPrivate(ctx context.Context, req *proto.MarkAsPrivateRequest) (*proto.MarkAsPrivateResponse, error) {
	// Extraire le userID depuis le contexte (via le JWT transmis dans le header Authorization)
	userID, err := jwt.ExtractUserIDFromContext(ctx)
	if err != nil {
		log.Printf(" Impossible d'extraire le userID : %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "Token invalide ou manquant")
	}

	// Appeler le service métier avec l'ID utilisateur correct
	if err := s.mediaService.MarkAsPrivate(uint(req.MediaId), userID); err != nil {
		log.Printf(" Erreur lors du passage en privé : %v", err)
		return nil, status.Errorf(codes.Unknown, "Erreur lors du passage en privé : %v", err)
	}

	log.Printf("Média %d marqué comme privé par userID=%d", req.MediaId, userID)
	return &proto.MarkAsPrivateResponse{}, nil
}

func (s *galleryServer) GetPrivateMedia(ctx context.Context, req *proto.GetPrivateMediaRequest) (*proto.GetPrivateMediaResponse, error) {

	err := s.userService.VerifyPrivateAlbumPin(uint(req.UserId), req.Pin)
	if err != nil {
		log.Printf("Error verifying private album pin: %v", err)
		return nil, err
	}

	media, err := s.mediaService.GetPrivateMedia(uint(req.UserId))
	if err != nil {
		log.Printf("Error getting private media: %v", err)
		return nil, err
	}

	var protoMedia []*proto.Media
	for _, m := range media {
		protoMedia = append(protoMedia, &proto.Media{
			Id:       uint32(m.ID),
			Name:     m.Name,
			AlbumId:  uint32(m.AlbumID),
			FileSize: uint32(m.FileSize),
		})
	}

	return &proto.GetPrivateMediaResponse{
		Media: protoMedia,
	}, nil
}

func (s *galleryServer) DownloadMedia(ctx context.Context, req *proto.DownloadMediaRequest) (*proto.DownloadMediaResponse, error) {
	userID, err := jwt.ExtractUserIDFromContext(ctx)
	if err != nil {
		log.Printf("Erreur d'extraction du userID : %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "token invalide : %v", err)
	}

	log.Printf("Téléchargement demandé pour mediaID=%d par userID=%d", req.MediaId, userID)

	var buf bytes.Buffer
	if err := s.mediaService.DownloadMedia(uint(req.MediaId), userID, &buf); err != nil {
		log.Printf("Erreur lors du téléchargement du média : %v", err)
		return nil, status.Errorf(codes.Internal, "échec du téléchargement du média : %v", err)
	}

	return &proto.DownloadMediaResponse{
		FileData: buf.Bytes(),
	}, nil
}

func (s *galleryServer) DeleteMedia(ctx context.Context, req *proto.DeleteMediaRequest) (*proto.DeleteMediaResponse, error) {

	userID, err := jwt.ExtractUserIDFromContext(ctx)
	if err != nil {
		log.Printf("Erreur d'extraction du userID : %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "token invalide : %v", err)
	}

	if err := s.mediaService.DeleteMedia(uint(req.MediaId), userID); err != nil {
		log.Printf("Erreur lors de la suppression du média : %v", err)
		return nil, status.Errorf(codes.Internal, "Erreur lors de la suppression du média : %v", err)
	}

	return &proto.DeleteMediaResponse{
		Message: "Média supprimé avec succès",
	}, nil
}

func (s *galleryServer) DetectSimilarMedia(ctx context.Context, req *proto.DetectSimilarMediaRequest) (*proto.DetectSimilarMediaResponse, error) {
	userID, err := jwt.ExtractUserIDFromContext(ctx)
	if err != nil {
		log.Printf("Token invalide : %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "token invalide : %v", err)
	}

	log.Printf("🔍 Détection de similarité sur albumID=%d pour userID=%d", req.AlbumId, userID)

	// Appel à la méthode modifiée
	similarGroups, err := s.mediaService.DetectSimilarMedia(userID, uint(req.AlbumId))
	if err != nil {
		log.Printf("Erreur détection similarité : %v", err)
		return nil, status.Errorf(codes.Unknown, err.Error())
	}

	// Convertit les groupes en format gRPC
	var protoGroups []*proto.MediaGroup
	for _, group := range similarGroups {
		var protoMedia []*proto.Media
		for _, m := range group {
			protoMedia = append(protoMedia, &proto.Media{
				Id:       uint32(m.ID),
				Name:     m.Name,
				AlbumId:  uint32(m.AlbumID),
				FileSize: uint32(m.FileSize),
			})
		}
		protoGroups = append(protoGroups, &proto.MediaGroup{Media: protoMedia})
	}

	return &proto.DetectSimilarMediaResponse{
		Groups: protoGroups,
	}, nil
}

// User Service methods
func (s *galleryServer) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	if err := s.userService.CreateUser(req.Username, req.Email); err != nil {
		log.Printf("Error creating user: %v", err)
		return nil, err
	}

	return &proto.CreateUserResponse{}, nil
}

// func (s *galleryServer) AddMediaToFavorite(ctx context.Context, req *proto.AddMediaToFavoriteRequest) (*proto.AddMediaToFavoriteResponse, error) {
// 	if err := s.mediaService.AddMediaToFavorite(uint(req.MediaId)); err != nil {
// 		log.Printf("Error adding media to favorite: %v", err)
// 		return nil, err
// 	}

// 	return &proto.AddMediaToFavoriteResponse{}, nil
// }

func (s *galleryServer) GetMediaByAlbum(ctx context.Context, req *proto.GetMediaByAlbumRequest) (*proto.GetMediaByAlbumResponse, error) {
    medias, err := s.mediaService.GetMediaByAlbum(uint(req.AlbumId))
    if err != nil {
        return nil, status.Errorf(codes.Internal, "failed to retrieve media: %v", err)
    }

    var protoMedias []*proto.Media
    for _, m := range medias {
        protoMedias = append(protoMedias, &proto.Media{
            Id:        uint32(m.ID),
            Name:      m.Name,
            Path:      m.Path,
            AlbumId:   uint32(m.AlbumID),
			IsFavorite: m.IsFavorite,
        })
    }

    return &proto.GetMediaByAlbumResponse{Media: protoMedias}, nil
}


func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Avertissement : Impossible de charger le fichier .env, utilisation des variables système.")
	}

	// Initialiser le gestionnaire de base de données
	dbManager, err := db.NewDBManagerService()
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation de la base de données : %v", err)
	}
	defer func() {
		log.Println("Fermeture de la connexion à la base de données...")
		dbManager.CloseConnection()
	}()

	// Effectuer la migration des modèles
	if err := dbManager.AutoMigrate(); err != nil {
		log.Fatalf("Erreur lors de la migration des modèles : %v", err)
	}

	// Initialiser le S3Service
	s3Service := services.NewS3Service("http://my-s3-clone:9090")

	// Initialize services
	albumService := services.NewAlbumService(dbManager, s3Service)
	mediaService := services.NewMediaService(dbManager, s3Service)
	userService := services.NewUserService(dbManager, s3Service)

	// Initialiser le service JWT
	jwtService, err := jwt.NewJWTService()
	if err != nil {
		log.Fatalf("Erreur lors de l'initialisation de JWTService : %v", err)
	}

	// Définir les méthodes protégées (authentification requise)
	methodsToIntercept := map[string]bool{
		"/proto.AlbumService/CreateAlbum":     true,
		"/proto.AlbumService/UpdateAlbum":     true,
		"/proto.AlbumService/DeleteAlbum":     true,
		"/proto.AlbumService/GetPrivateAlbum": true,
		"/proto.MediaService/AddMedia":        true,
		"/proto.MediaService/MarkAsPrivate":   true,
		"/proto.MediaService/GetPrivateMedia": true,
		"/proto.MediaService/DownloadMedia":   true,
		"/proto.MediaService/DeleteMedia":     true,
		"/proto.MediaService/GetMediaByAlbum": true,
	}

	// Créer le serveur gRPC avec intercepteur JWT
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.AuthInterceptor(jwtService, methodsToIntercept)),
	)

	galleryServer := &galleryServer{
		albumService: albumService,
		mediaService: mediaService,
		userService:  userService,
	}

	// Enregistrer les services gRPC
	proto.RegisterAlbumServiceServer(grpcServer, galleryServer)
	proto.RegisterMediaServiceServer(grpcServer, galleryServer)
	proto.RegisterUserServiceServer(grpcServer, galleryServer)

	// Démarrer le serveur gRPC
	grpcListener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Erreur lors de l'écoute du serveur gRPC : %v", err)
	}

	// Canal pour gérer les signaux système (interruption ou arrêt)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Démarrer le serveur gRPC dans une goroutine
	go func() {
		log.Println("gRPC server started on port 50052...")
		if err := grpcServer.Serve(grpcListener); err != nil {
			log.Fatalf("Erreur lors de l'exécution du serveur gRPC : %v", err)
		}
	}()

	// Attendre un signal d'arrêt
	<-stop
	log.Println("Signal reçu, arrêt des services...")

	// Arrêter gracieusement les serveurs
	grpcServer.GracefulStop()
	log.Println("Server stopped successfully.")
}
