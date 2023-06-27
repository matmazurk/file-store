package filestore

import (
	"context"

	"google.golang.org/grpc"

	pb "github.com/matmazurk/file-store/proto"
)

type service struct {
	pb.UnsafeFileServer
}

func NewService() service {
	return service{}
}

func (s service) ListFiles(context.Context, *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	return nil, nil
}

func (s service) StoreFile(stream pb.File_StoreFileServer) error {
	return nil
}

func (s service) RegisterIn(serv *grpc.Server) {
	pb.RegisterFileServer(serv, s)
}
