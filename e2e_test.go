package test

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
	"time"

	fakegrpcserver "github.com/matmazurk/fake-grpc-server"
	"github.com/matmazurk/file-store/client"
	server "github.com/matmazurk/file-store/server"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

func TestService(t *testing.T) {
	const testfilesTargetDir = "./__test"
	service := server.NewService(testfilesTargetDir)
	fakeServer := fakegrpcserver.NewFakeServer(func(s *grpc.Server) {
		service.RegisterIn(s)
	})
	stop := fakeServer.Start()
	defer stop()

	conn, err := fakeServer.Conn()
	require.NoError(t, err)
	c := client.New(conn)

	t.Run("should properly store all test files", func(t *testing.T) {
		const testfilesDir = "./testfiles"
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
				requireFileEventuallyPresent(t, expectedLocationPath)
				requireFilesExact(t, toSendPath, expectedLocationPath)
			})
		}
	})
}

func requireFileEventuallyPresent(t *testing.T, path string) {
	t.Helper()

	require.Eventually(
		t,
		func() bool {
			_, err := os.Stat(path)
			return err == nil
		},
		500*time.Millisecond,
		50*time.Millisecond)
}

func requireFilesExact(t *testing.T, expectedPath, actualPath string) {
	t.Helper()

	expected, err := os.ReadFile(expectedPath)
	require.NoError(t, err)

	actual, err := os.ReadFile(actualPath)
	require.NoError(t, err)

	require.True(t, bytes.Equal(expected, actual))
}
