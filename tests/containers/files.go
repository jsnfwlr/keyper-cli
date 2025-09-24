package containers

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
)

func TarFile(basePath string, fileContent []byte, fileMode int64) (*bytes.Buffer, error) {
	buffer := &bytes.Buffer{}

	zr := gzip.NewWriter(buffer)
	tw := tar.NewWriter(zr)

	hdr := &tar.Header{
		Name: basePath,
		Mode: fileMode,
		Size: int64(len(fileContent)),
	}
	if err := tw.WriteHeader(hdr); err != nil {
		return &bytes.Buffer{}, err
	}

	_, err := tw.Write(fileContent)
	if err != nil {
		return &bytes.Buffer{}, err
	}

	// produce tar
	if err := tw.Close(); err != nil {
		return &bytes.Buffer{}, fmt.Errorf("error closing tar file: %w", err)
	}
	// produce gzip
	if err := zr.Close(); err != nil {
		return &bytes.Buffer{}, fmt.Errorf("error closing gzip file: %w", err)
	}

	return buffer, nil
}
