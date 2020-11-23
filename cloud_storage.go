package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/ncw/swift"
)

// CloudClient takes care of files backup
type CloudClient struct {
	conn      swiftConnection
	container string
}

type swiftConnection interface {
	ObjectPut(container string, objectName string, contents io.Reader, checkHash bool, Hash string, contentType string, h swift.Headers) (headers swift.Headers, err error)
	ObjectDelete(container string, objectName string) error
}

// NewCloudClient Create a client connect with Cloud
func NewCloudClient(userName, apiKey, authURL, domain, containerName string) *CloudClient {
	return &CloudClient{conn: &swift.Connection{UserName: userName, ApiKey: apiKey, AuthUrl: authURL, Domain: domain}, container: containerName}
}

//UploadFile Store a file in cloud container and return a Backup file containing a URL and a Hash for that file.
func (cloud *CloudClient) UploadFile(srcPath string, dstFolder string) (*Backup, error) {
	f, err := os.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("error Opening file at %s: %v", f.Name, err)
	}
	defer f.Close()
	dstPath := filepath.Join(dstFolder, filepath.Base(srcPath))
	headers, err := cloud.conn.ObjectPut(cloud.container, dstPath, f, true, "", "", nil)
	if err != nil {
		return nil, fmt.Errorf("error trying to upload file at %s to storage: %v\nHeaders: %v", dstPath, err, headers)
	}
	return &Backup{URL: fmt.Sprintf("%s/%s/%s", cloud.storageURL(), cloud.container, dstPath), Hash: headers["Etag"]}, nil
}

//storageURL finds cloud repository url
func (cloud *CloudClient) storageURL() string {
	if v, ok := cloud.conn.(*swift.Connection); ok {
		return v.StorageUrl
	}
	return ""
}

//deleteFile delete a file from cloud container.
func (cloud *CloudClient) deleteFile(path string) error {
	err := cloud.conn.ObjectDelete(cloud.container, filepath.Base(path))
	if err != nil {
		return fmt.Errorf("delete file error: 'BackupClient:deleteFile' %s to storage: %q\nHeaders", path, err)
	}
	return nil
}

//Backup is the API to make URL and HASH files to be used in mgo store function
func (cloud *CloudClient) Backup(Files []string, dstFolder string) ([]Backup, error) {
	if len(Files) == 0 {
		return []Backup{}, nil
	}
	var backups []Backup
	for _, value := range Files {
		back, err := cloud.UploadFile(value, dstFolder)
		if err != nil {
			return nil, fmt.Errorf("Error in BackupClient:backup upload file %v", err)
		}
		backups = append(backups, *back)
	}
	return backups, nil
}
