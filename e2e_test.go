package test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	fakegrpcserver "github.com/matmazurk/fake-grpc-server"
	"github.com/matmazurk/file-store/client"
	server "github.com/matmazurk/file-store/server"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestService(t *testing.T) {
	const (
		testfilesTargetDir = "./__test"
		testfilesDir       = "./testfiles"
	)
	t.Cleanup(func() { os.RemoveAll(testfilesTargetDir) })

	service := server.NewService(testfilesTargetDir)
	fakeServer := fakegrpcserver.NewFakeServer(func(s *grpc.Server) {
		service.RegisterIn(s)
	})
	stop := fakeServer.Start()
	t.Cleanup(func() { stop() })

	conn, err := fakeServer.Conn()
	require.NoError(t, err)
	c := client.New(conn)

	t.Run("should properly store all test files", func(t *testing.T) {
		testfiles, err := os.ReadDir(testfilesDir)
		require.NoError(t, err)

		for _, testfile := range testfiles {
			testfile := testfile.Name()
			t.Run(testfile, func(t *testing.T) {
				t.Parallel()

				toSendPath := filepath.Join(testfilesDir, testfile)
				err := c.SendFile(toSendPath)
				require.NoError(t, err)

				expectedLocationPath := filepath.Join(testfilesTargetDir, testfilesDir, testfile)

				expected, err := os.ReadFile(expectedLocationPath)
				require.NoError(t, err)

				actual, err := os.ReadFile(toSendPath)
				require.NoError(t, err)
				require.True(t, bytes.Equal(expected, actual))
			})
		}
	})
}
