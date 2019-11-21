package storage

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/ncw/swift"
)

// StoreClient struct containing a conn with swift library
type StorageClient struct {
	conn swift.Connection
}

// NewStorageClient Create a client connect with Cloud
func NewStorageClient() *StorageClient {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	return &StorageClient{conn: swift.Connection{UserName: os.Getenv("USERNAME"), ApiKey: os.Getenv("APIKEY"), AuthUrl: os.Getenv("AUTHURL"), Domain: os.Getenv("DOMAIN")}}
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
	f := bytes.NewReader(content)
	hasher := md5.New()
	if _, err := io.Copy(hasher, f); err != nil {
		return "", fmt.Errorf("error copying file content to hash.Hash")
	}
	value := hex.EncodeToString(hasher.Sum(nil)[:16])
	return value, nil
}

func getFileName(path string) string {
	split := strings.SplitAfter(path, "/")
	return split[len(split)-1]
}

//UploadFile Store a file in cloud container and return a Backup file containing a URL and a Hash for that file.
func (sc *StorageClient) UploadFile(path string) (*Backup, error) {
	if !sc.conn.Authenticated() {
		return nil, fmt.Errorf("Not Authenticated")
	}

	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading file at %s: %q", path, err)
	}

	hash, err := md5Hash(content)
	if err != nil {
		return nil, fmt.Errorf("error generating file hash at %s: %q", path, err)
	}
	fileName := getFileName(path)
	headers, err := sc.conn.ObjectPut("DadosJusBr", fileName, bytes.NewReader(content), true, hash, "", nil)
	if err != nil {
		return nil, fmt.Errorf("error trying to upload file at %s to storage: %q\nHeaders: %v", path, err, headers)
	}
	return &Backup{URL: os.Getenv("ENDURL") + fileName, Hash: hash}, nil
}

//DeleteFile delete a file from cloud container.
func (sc *StorageClient) DeleteFile(path string) error {
	if !sc.conn.Authenticated() {
		return fmt.Errorf("Not Authenticated")
	}
	fileName := getFileName(path)
	err := sc.conn.ObjectDelete("DadosJusBr", fileName)
	if err != nil {
		return fmt.Errorf("error trying to delete file at %s to storage: %q\nHeaders", path, err)
	}
	return nil
}
