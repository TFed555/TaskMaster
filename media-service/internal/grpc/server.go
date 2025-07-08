package grpc_media

import (
	"context"
	"log"
	"media-service/internal/services"
	"net"
	pb "shared/pkg/media-service"

	"google.golang.org/grpc"
)

type MediaServer struct {
	pb.UnimplementedMediaServiceServer
	service services.MediaService
}


func (m MediaServer) GetAvatarPic(ctx context.Context, req *pb.AvatarRequest) (*pb.AvatarResponse, error) {
	imgURL, err := m.service.GetAvatarPic(req.ImgName)
	if err != nil {
		return nil, err
	}
	return &pb.AvatarResponse{
		ImgURL: imgURL,
	}, nil
}

func StartGRPCServer(mediaService services.MediaService, port string) error {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	s := grpc.NewServer()
	pb.RegisterMediaServiceServer(s, &MediaServer{service: mediaService})
	
	log.Printf("gRPC server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		return err
	}

	return nil
}