package grpc_client

import (
	"context"
	pb "shared/pkg/media-service"

	"google.golang.org/grpc"
)

type GRPCMediaService interface {
	GetAvatarPic(img_name string) (string, error)
}

type GRPCMediaServiceImpl struct {
	client pb.MediaServiceClient
}

func NewGRPCMediaService(conn *grpc.ClientConn) *GRPCMediaServiceImpl {
	return &GRPCMediaServiceImpl{
		client: pb.NewMediaServiceClient(conn),
	}
}

func (s GRPCMediaServiceImpl) GetAvatarPic(img_name string) (string, error) {
	resp, err := s.client.GetAvatarPic(context.Background(), &pb.AvatarRequest{ImgName: img_name})
	if err != nil {
		return "", err
	}
	return resp.ImgURL, nil
}