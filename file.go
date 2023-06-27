package filestore

import (
	"context"

	"github.com/pkg/errors"

	"crypto"
	"os"
	"path"
)

func saveFile(
	ctx context.Context,
	rootDir, destPath string,
	dataCh <-chan []byte,
) ([]byte, error) {
	absolutePath := path.Join(rootDir, destPath)
	err := os.MkdirAll(absolutePath, os.ModePerm)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(absolutePath)
	if err != nil {
		return nil, err
	}

	md5, err := writeChunks(ctx, f, absolutePath, dataCh)
	f.Close()
	if err != nil {
		os.Remove(absolutePath)
	}
	return md5, err
}

func writeChunks(
	ctx context.Context,
	f *os.File,
	absolutePath string,
	dataCh <-chan []byte,
) ([]byte, error) {
	md5 := crypto.MD5.New()
	for {
		select {
		case <-ctx.Done():
			return nil, errors.Errorf("context cancelled for %s", absolutePath)
		case chunk, ok := <-dataCh:
			if !ok {
				return md5.Sum(nil), nil
			}

			_, err := f.Write(chunk)
			if err != nil {
				return nil, errors.Errorf("couldn't write chunk to file %s", absolutePath)
			}
			_, err = md5.Write(chunk)
			if err != nil {
				return nil, errors.Errorf("couldn't write md5 chunk for file %s", absolutePath)
			}
		}
	}
}
