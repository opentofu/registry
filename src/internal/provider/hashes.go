package provider

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"fmt"
	"io"

	"golang.org/x/mod/sumdb/dirhash"
)

func (p Provider) CalculateHash1(releaseDownloadUrl string, shaExpected string) (string, error) {
	contents, assetErr := p.Github.DownloadAssetContents(releaseDownloadUrl)
	if assetErr != nil {
		return "", fmt.Errorf("failed to download release for calcualting hashes: %w", assetErr)
	}

	shaHasher := sha256.New()
	_, err := io.Copy(shaHasher, bytes.NewBuffer(contents))
	if err != nil {
		panic(err)
	}

	shaCalculated := fmt.Sprintf("%x", shaHasher.Sum(nil))
	if shaCalculated != shaExpected {
		return "", fmt.Errorf("expected SHA256 %q, got %q", shaExpected, shaCalculated)
	}

	h1, err := HashZip(contents, dirhash.Hash1)
	if err != nil {
		return "", fmt.Errorf("unable to hash provider release: %w", err)
	}
	return h1, nil
}

// Inspired heavily from golang.org/x/mod/sumdb/dirhash, but re-written to use a in-memory zip stream
// Ideally this would be contributed upstream
func HashZip(zipdata []byte, hash dirhash.Hash) (string, error) {
	z, err := zip.NewReader(bytes.NewReader(zipdata), int64(len(zipdata)))
	if err != nil {
		return "", err
	}
	var files []string
	zfiles := make(map[string]*zip.File)
	for _, file := range z.File {
		files = append(files, file.Name)
		zfiles[file.Name] = file
	}
	zipOpen := func(name string) (io.ReadCloser, error) {
		f := zfiles[name]
		if f == nil {
			return nil, fmt.Errorf("file %q not found in zip", name) // should never happen
		}
		return f.Open()
	}
	return hash(files, zipOpen)
}
