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

const chunkSize = 1 << 16

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

func (c client) SendFile(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer f.Close()

	stream, err := c.c.StoreFile(context.Background())
	if err != nil {
		return err
	}
	md5 := crypto.MD5.New()
	chunk := make([]byte, chunkSize)
	for {
		n, err := f.Read(chunk)
		if err != nil {
			if errors.Is(err, io.EOF) {
				chunk := chunk[:n]
				md5.Write(chunk)
				stream.Send(&pb.StoreFileMsg{
					ChunkData: chunk,
					Md_5:      md5.Sum(nil),
				})

				return stream.CloseSend()
			}
		}

		chunk := chunk[:n]
		md5.Write(chunk)
		err = stream.Send(&pb.StoreFileMsg{
			Path:      filePath,
			ChunkData: chunk,
		})
		if err != nil {
			return err
		}
	}
}
