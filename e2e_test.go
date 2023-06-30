package test

import (
	"context"
	"log"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/matmazurk/file-store/client"
	server "github.com/matmazurk/file-store/server"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

var lis *bufconn.Listener

const (
	testFilesDir = "./_test-files"
	bufSize      = 1024 * 1024
)

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestService(t *testing.T) {
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	service := server.NewService(testFilesDir)
	service.RegisterIn(s)
	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("Server exited with error: %v", err)
		}
	}()
	defer s.GracefulStop()

	ctx := context.Background()
	conn, err := grpc.DialContext(
		ctx,
		"bufnet",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}

	c := client.New(conn)

	t.Run("test", func(t *testing.T) {
		const smallFilePath = "./testfiles/small-file.txt"
		err = c.SendFile(smallFilePath)
		require.NoError(t, err)
		require.Eventually(
			t,
			func() bool {
				_, err := os.Stat(filepath.Join(testFilesDir, smallFilePath))
				return err == nil
			},
			time.Second,
			time.Millisecond*100)
	})

}
