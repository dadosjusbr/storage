package storage

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ncw/swift"
)

// StoreClient struct containing a conn with swift library
type StorageClient struct {
	conn swift.Connection
}

const (
	jusContainer = "DadosJusBr"
)

// NewStorageClient Create a client connect with Cloud
func NewStorageClient(userName, apiKey, authURL, domain string) *StorageClient {
	return &StorageClient{conn: swift.Connection{UserName: userName, ApiKey: apiKey, AuthUrl: authURL, Domain: domain}}
}

//Authenticate Authenticate a client in Cloud Service
func (sc *StorageClient) Authenticate() error {
	err := sc.conn.Authenticate()
	if err != nil {
		return fmt.Errorf("error creating swift.Connection: %q", err)
	}
	return nil
}

//md5Hash calculate a md5Hasg for a file
func md5Hash(content []byte) (string, error) {
	hasher := md5.New()
	if _, err := io.Copy(hasher, bytes.NewReader(content)); err != nil {
		return "", fmt.Errorf("error hashing file contents: %q", err)
	}
	return hex.EncodeToString(hasher.Sum(nil)[:16]), nil
}

//UploadFile Store a file in cloud container and return a Backup file containing a URL and a Hash for that file.
func (sc *StorageClient) UploadFile(path string) (*Backup, error) {
	if !sc.conn.Authenticated() {
		return nil, fmt.Errorf("Not Authenticated")
	}
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("error Opening file at %s: %q", path, err)
	}

	headers, err := sc.conn.ObjectPut(jusContainer, filepath.Base(path), f, true, "", "", nil)
	if err != nil {
		return nil, fmt.Errorf("error trying to upload file at %s to storage: %q\nHeaders: %v", path, err, headers)
	}
	return &Backup{URL: fmt.Sprintf("%s/%s/%s", sc.conn.StorageUrl, jusContainer, filepath.Base(path)), Hash: headers["Etag"]}, nil
}

//DeleteFile delete a file from cloud container.
func (sc *StorageClient) deleteFile(path string) error {
	if !sc.conn.Authenticated() {
		return fmt.Errorf("Not Authenticated")
	}
	err := sc.conn.ObjectDelete("DadosJusBr", filepath.Base(path))
	if err != nil {
		return fmt.Errorf("error trying to delete file at %s to storage: %q\nHeaders", path, err)
	}
	return nil
}

func (sc *StorageClient) backup(Files []string) ([]Backup, error) {
	if len(Files) == 0 {
		return nil, fmt.Errorf("There is no file to upload")
	}
	var backups []Backup
	if err := sc.Authenticate(); err != nil {
		return nil, fmt.Errorf("Authentication error: %q", err)
	}

	for _, value := range Files {
		back, err := sc.UploadFile(value)
		if err != nil {
			return nil, fmt.Errorf("Error no upload do arquivo %v", err)
		}
		backups = append(backups, *back)

	}
	return backups, nil
}
