package middleware

import (
	"context"

	"google.golang.org/grpc"
	pb "auth-service/internal/grpc/pkg"
)

type GRPCAuthService struct {
	client pb.AuthServiceClient
}

func NewGRPCAuthService(conn *grpc.ClientConn) *GRPCAuthService {
	return &GRPCAuthService{
		client: pb.NewAuthServiceClient(conn),
	}
}

func (s *GRPCAuthService) ValidateToken(token string) (bool, string) {
	resp, err := s.client.ValidateToken(context.Background(), &pb.TokenRequest{Token: token})
	if err != nil {
		return false, err.Error()
	}
	return resp.Valid, resp.Error
}

func (s *GRPCAuthService) UpdateAccessToken(refreshToken string) (bool, string, error) {
	resp, err := s.client.UpdateAccessToken(context.Background(), &pb.RefreshTokenRequest{RefreshToken: refreshToken})
	if err != nil {
		return false, "", err
	}
	return resp.Success, resp.NewToken, nil
}