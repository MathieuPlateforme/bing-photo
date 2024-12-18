package handlers

import (
	"context"
	"net/http"

	proto "ApiGateway/proto"
)

type ApiGateway struct {
	AuthClient proto.AuthServiceClient
}

func NewApiGateway(authClient proto.AuthServiceClient) *ApiGateway {
	return &ApiGateway{AuthClient: authClient}
}

func (g *ApiGateway) LoginHandler(w http.ResponseWriter, r *http.Request) {
	req := &proto.LoginRequest{
		Email:    r.FormValue("email"),
		Password: r.FormValue("password"),
	}
	res, err := g.AuthClient.Login(context.Background(), req)
	if err != nil {
		http.Error(w, "Login failed", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Token: " + res.Token))
}

func (g *ApiGateway) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	req := &proto.RegisterRequest{
		Email:     r.FormValue("email"),
		Password:  r.FormValue("password"),
		Username:  r.FormValue("username"),
		FirstName: r.FormValue("firstName"),
		LastName:  r.FormValue("lastName"),
	}
	res, err := g.AuthClient.Register(context.Background(), req)
	if err != nil {
		http.Error(w, "Register failed", http.StatusInternalServerError)
		return
	}
	w.Write([]byte("Message: " + res.Message))
}
