package filestore

import (
	"bytes"
	"context"
	"crypto"
	"errors"
	"io"

	"google.golang.org/grpc"

	pb "github.com/matmazurk/file-store/proto"
)

type service struct {
	pb.UnimplementedFileServer

	rootDir string
}

func NewService(rootDir string) service {
	return service{
		rootDir: rootDir,
	}
}

func (s service) ListFiles(context.Context, *pb.ListFilesRequest) (*pb.ListFilesResponse, error) {
	paths, err := listAllFiles(s.rootDir)
	if err != nil {
		return nil, err
	}

	return &pb.ListFilesResponse{
		Paths: paths,
	}, nil
}

func (s service) StoreFile(stream pb.File_StoreFileServer) error {
	rec, err := stream.Recv()
	if err != nil {

	}

	path := rec.GetPath()
	if path == "" {

	}

	errCh := make(chan error)
	dataCh := make(chan []byte)
	go func() {
		err := saveFile(s.rootDir, path, dataCh)
		if err != nil {
			errCh <- err
		}
	}()
	md5 := crypto.MD5.New()
	md5.Write(rec.GetChunkData())
	dataCh <- rec.GetChunkData()

	receivedMD5 := rec.GetMd_5()
	for {
		rec, err := stream.Recv()
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
		}
		md5.Write(rec.GetChunkData())
		dataCh <- rec.GetChunkData()

		if len(rec.Md_5) > 0 {
			receivedMD5 = rec.GetMd_5()
		}
	}
	close(dataCh)

	if err := <-errCh; err != nil {
		return err
	}
	if !bytes.Equal(receivedMD5, md5.Sum(nil)) {
		return errors.New("invalid md5")
	}

	return stream.SendAndClose(&pb.StoreFileResponse{})
}

func (s service) RegisterIn(serv *grpc.Server) {
	pb.RegisterFileServer(serv, s)
}
