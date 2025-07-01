package grpc

import (
	"auth-service/internal/services"
	"context"
	"log"
	"net"

	"google.golang.org/grpc"
	pb "shared/pkg"
)

type AuthServer struct {
	pb.UnimplementedAuthServiceServer
	service services.AuthService
}

func (s *AuthServer) ValidateToken(ctx context.Context, req *pb.TokenRequest) (*pb.TokenResponse, error) {
	valid, errMsg := s.service.ValidateToken(req.Token)
	return &pb.TokenResponse{Valid: valid, Error: errMsg}, nil
}

func (s *AuthServer) UpdateAccessToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.AccessTokenResponse, error) {
	success, newToken, err := s.service.UpdateAccessToken(req.RefreshToken)
	return &pb.AccessTokenResponse{
		Success:  success,
		AccessToken: newToken,
		Error:    err.Error(),
	}, nil
}

func StartGRPCServer(authService services.AuthService, port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, &AuthServer{service: authService})
	
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil
}