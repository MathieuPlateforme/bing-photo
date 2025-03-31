package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"

	proto "ApiGateway/proto"
)

type GalleryGateway struct {
	GalleryClient proto.AlbumServiceClient
	MediaClient   proto.MediaServiceClient
	UserClient    proto.UserServiceClient
}

func NewGalleryGateway(albumClient proto.AlbumServiceClient, mediaClient proto.MediaServiceClient, userClient proto.UserServiceClient) *GalleryGateway {
	return &GalleryGateway{
		GalleryClient: albumClient,
		MediaClient:   mediaClient,
		UserClient:    userClient,
	}
}

// Album handlers
func (g *GalleryGateway) CreateAlbumHandler(w http.ResponseWriter, r *http.Request) {
	var req proto.CreateAlbumRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Failed to parse request: %v\n", err)
		return
	}

	res, err := g.GalleryClient.CreateAlbum(context.Background(), &req)
	if err != nil {
		http.Error(w, "Failed to create album: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Create album error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (g *GalleryGateway) GetAlbumsByUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(r.URL.Query().Get("user_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	req := &proto.GetAlbumsByUserRequest{
		UserId: uint32(userID),
	}

	res, err := g.GalleryClient.GetAlbumsByUser(context.Background(), req)
	if err != nil {
		http.Error(w, "Failed to get albums: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Get albums error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (g *GalleryGateway) UpdateAlbumHandler(w http.ResponseWriter, r *http.Request) {
	var req proto.UpdateAlbumRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Failed to parse request: %v\n", err)
		return
	}

	res, err := g.GalleryClient.UpdateAlbum(context.Background(), &req)
	if err != nil {
		http.Error(w, "Failed to update album: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Update album error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (g *GalleryGateway) DeleteAlbumHandler(w http.ResponseWriter, r *http.Request) {
	albumID, err := strconv.ParseUint(r.URL.Query().Get("album_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid album ID", http.StatusBadRequest)
		return
	}

	req := &proto.DeleteAlbumRequest{
		AlbumId: uint32(albumID),
	}

	res, err := g.GalleryClient.DeleteAlbum(context.Background(), req)
	if err != nil {
		http.Error(w, "Failed to delete album: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Delete album error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (g *GalleryGateway) GetPrivateAlbumHandler(w http.ResponseWriter, r *http.Request) {
	albumID, err := strconv.ParseUint(r.URL.Query().Get("album_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid album ID", http.StatusBadRequest)
		return
	}

	req := &proto.GetPrivateAlbumRequest{
		AlbumId: uint32(albumID),
	}

	res, err := g.GalleryClient.GetPrivateAlbum(context.Background(), req)
	if err != nil {
		http.Error(w, "Failed to get private album: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Get private album error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// Media handlers
func (g *GalleryGateway) AddMediaHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(10 << 20) // 10 MB max
	if err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Failed to read file", http.StatusInternalServerError)
		return
	}

	albumID, err := strconv.ParseUint(r.FormValue("album_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid album ID", http.StatusBadRequest)
		return
	}

	req := &proto.AddMediaRequest{
		Name:     header.Filename,
		AlbumId:  uint32(albumID),
		FileData: fileData,
	}

	res, err := g.MediaClient.AddMedia(context.Background(), req)
	if err != nil {
		http.Error(w, "Failed to add media: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Add media error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}

func (g *GalleryGateway) GetMediaByUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(r.URL.Query().Get("user_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	req := &proto.GetMediaByUserRequest{
		UserId: uint32(userID),
	}

	res, err := g.MediaClient.GetMediaByUser(context.Background(), req)
	if err != nil {
		http.Error(w, "Failed to get media: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Get media error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (g *GalleryGateway) MarkAsPrivateHandler(w http.ResponseWriter, r *http.Request) {
	var req proto.MarkAsPrivateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Failed to parse request: %v\n", err)
		return
	}

	res, err := g.MediaClient.MarkAsPrivate(context.Background(), &req)
	if err != nil {
		http.Error(w, "Failed to mark as private: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Mark as private error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (g *GalleryGateway) GetPrivateMediaHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseUint(r.URL.Query().Get("user_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	req := &proto.GetPrivateMediaRequest{
		UserId: uint32(userID),
	}

	res, err := g.MediaClient.GetPrivateMedia(context.Background(), req)
	if err != nil {
		http.Error(w, "Failed to get private media: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Get private media error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (g *GalleryGateway) DownloadMediaHandler(w http.ResponseWriter, r *http.Request) {
	mediaID, err := strconv.ParseUint(r.URL.Query().Get("media_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid media ID", http.StatusBadRequest)
		return
	}

	req := &proto.DownloadMediaRequest{
		MediaId: uint32(mediaID),
	}

	res, err := g.MediaClient.DownloadMedia(context.Background(), req)
	if err != nil {
		http.Error(w, "Failed to download media: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Download media error: %v\n", err)
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", "attachment; filename=media")
	w.WriteHeader(http.StatusOK)
	w.Write(res.FileData)
}

func (g *GalleryGateway) DeleteMediaHandler(w http.ResponseWriter, r *http.Request) {
	mediaID, err := strconv.ParseUint(r.URL.Query().Get("media_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid media ID", http.StatusBadRequest)
		return
	}

	req := &proto.DeleteMediaRequest{
		MediaId: uint32(mediaID),
	}

	res, err := g.MediaClient.DeleteMedia(context.Background(), req)
	if err != nil {
		http.Error(w, "Failed to delete media: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Delete media error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

func (g *GalleryGateway) DetectSimilarMediaHandler(w http.ResponseWriter, r *http.Request) {
	mediaID, err := strconv.ParseUint(r.URL.Query().Get("media_id"), 10, 32)
	if err != nil {
		http.Error(w, "Invalid media ID", http.StatusBadRequest)
		return
	}

	req := &proto.DetectSimilarMediaRequest{
		MediaId: uint32(mediaID),
	}

	res, err := g.MediaClient.DetectSimilarMedia(context.Background(), req)
	if err != nil {
		http.Error(w, "Failed to detect similar media: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Detect similar media error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(res)
}

// User handlers
func (g *GalleryGateway) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var req proto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		log.Printf("Failed to parse request: %v\n", err)
		return
	}

	res, err := g.UserClient.CreateUser(context.Background(), &req)
	if err != nil {
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusInternalServerError)
		log.Printf("Create user error: %v\n", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)
}
