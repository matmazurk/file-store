package filestore

import (
	"path/filepath"

	"os"
	"path"

	"github.com/pkg/errors"
)

func saveFile(
	rootDir, destPath string,
	dataCh <-chan []byte,
) error {
	absolutePath := path.Join(rootDir, destPath)
	err := os.MkdirAll(absolutePath, os.ModePerm)
	if err != nil {
		return err
	}

	f, err := os.Open(absolutePath)
	if err != nil {
		return err
	}

	err = writeChunks(f, dataCh)
	f.Close()
	if err != nil {
		os.Remove(absolutePath)
		return errors.Wrapf(err, "couldn't write chunks to %s", absolutePath)
	}

	return nil
}

func writeChunks(
	f *os.File,
	dataCh <-chan []byte,
) error {
	for chunk := range dataCh {
		_, err := f.Write(chunk)
		if err != nil {
			return err
		}
	}
	return nil
}

func listAllFiles(rootDir string) ([]string, error) {
	var ret []string
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			ret = append(ret, path)
		}
		return nil
	})
	return nil, err
}
