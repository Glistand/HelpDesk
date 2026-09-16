package grpcserver

import (
	"context"

	authv1 "github.com/Glistand/HelpDesk/api/gen/go/helpdesk/auth/v1"
	"github.com/Glistand/HelpDesk/libs/grpckit/statuserr"
	"github.com/Glistand/HelpDesk/services/auth-service/internal/store"
	"github.com/Glistand/HelpDesk/services/auth-service/internal/tokens"
)

type Server struct {
	authv1.UnimplementedAuthServiceServer
	users  *store.Store
	issuer *tokens.Issuer
}

func New(users *store.Store, issuer *tokens.Issuer) *Server {
	return &Server{users: users, issuer: issuer}
}

func (s *Server) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	u, err := s.users.Authenticate(req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, statuserr.InvalidArgument("invalid email or password")
	}
	token, exp, err := s.issuer.Issue(u)
	if err != nil {
		return nil, statuserr.FromError(err)
	}
	return &authv1.LoginResponse{
		AccessToken:   token,
		User:          u.Proto(),
		ExpiresAtUnix: exp.Unix(),
	}, nil
}

func (s *Server) ValidateToken(ctx context.Context, req *authv1.ValidateTokenRequest) (*authv1.ValidateTokenResponse, error) {
	claims, err := s.issuer.Parse(req.GetAccessToken())
	if err != nil {
		return &authv1.ValidateTokenResponse{Valid: false}, nil
	}
	return &authv1.ValidateTokenResponse{
		Valid: true,
		User: &authv1.User{
			Id:    claims.UserID,
			Email: claims.Email,
			Name:  claims.Name,
			Role:  claims.Role,
		},
	}, nil
}

func (s *Server) GetUser(ctx context.Context, req *authv1.GetUserRequest) (*authv1.GetUserResponse, error) {
	u, err := s.users.Get(req.GetId())
	if err != nil {
		if err == store.ErrNotFound {
			return nil, statuserr.NotFound("user not found")
		}
		return nil, statuserr.FromError(err)
	}
	return &authv1.GetUserResponse{User: u.Proto()}, nil
}
