package client

import (
	"context"
	"crypto"
	"errors"
	"io"
	"os"

	pb "github.com/matmazurk/file-store/proto"
	"google.golang.org/grpc"
)

const bufSize = 1 << 16

type client struct {
	c pb.FileClient
}

func New(
	cc grpc.ClientConnInterface,
) client {

	return client{
		c: pb.NewFileClient(cc),
	}
}

func (c client) Send(content io.Reader, destination string) error {
	stream, err := c.c.StoreFile(context.Background())
	if err != nil {
		return err
	}

	md5 := crypto.MD5.New()
	buf := make([]byte, bufSize)
	for {
		n, err := content.Read(buf)
		buf = buf[:n]
		if err != nil {
			if errors.Is(err, io.EOF) {
				err := stream.Send(&pb.StoreFileMsg{
					Path:      destination,
					ChunkData: buf,
					Md_5:      md5.Sum(buf),
				})
				if err != nil {
					return err
				}

				_, err = stream.CloseAndRecv()
				return err
			}

			return err
		}

		md5.Write(buf)
		err = stream.Send(&pb.StoreFileMsg{
			Path:      destination,
			ChunkData: buf,
		})
		if err != nil {
			return err
		}
	}
}

func (c client) SendFile(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	return c.Send(f, filePath)
}
