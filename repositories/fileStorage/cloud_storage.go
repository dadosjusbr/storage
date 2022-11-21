package fileStorage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dadosjusbr/storage/models"
	"github.com/ncw/swift"
)

// CloudClient takes care of files backup
type CloudClient struct {
	conn      *swift.Connection
	container string
}

// NewCloudClient Create a client connect with Cloud
func NewCloudClient(userName, apiKey, authURL, domain, containerName string) *CloudClient {
	return &CloudClient{conn: &swift.Connection{UserName: userName, ApiKey: apiKey, AuthUrl: authURL, Domain: domain}, container: containerName}
}

//UploadFile Store a file in cloud container and return a Backup file containing a URL and a Hash for that file.
func (cloud *CloudClient) UploadFile(srcPath string, dstFolder string) (*models.Backup, error) {
	if !cloud.conn.Authenticated() {
		if err := cloud.conn.Authenticate(); err != nil {
			return nil, fmt.Errorf("error authenticating to swift:%q", err)
		}
		defer cloud.conn.UnAuthenticate()
	}

	f, err := os.Open(srcPath)
	if err != nil {
		return nil, fmt.Errorf("error Opening file at %s: %v", srcPath, err)
	}
	defer f.Close()
	d := filepath.Join(dstFolder, filepath.Base(srcPath))
	// this varaiable exists because we need to format the file path to something who swift can understand
	dstPath := strings.Replace(d, "\\", "/", -1)
	headers, err := cloud.conn.ObjectPut(cloud.container, dstPath, f, true, "", "", nil)
	if err != nil {
		return nil, fmt.Errorf("error trying to upload file at %s to storage: %v\nHeaders: %v", dstPath, err, headers)
	}
	return &models.Backup{URL: fmt.Sprintf("%s/%s/%s", cloud.conn.StorageUrl, cloud.container, dstPath), Hash: headers["Etag"]}, nil
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
func (cloud *CloudClient) Backup(Files []string, dstFolder string) ([]models.Backup, error) {
	if !cloud.conn.Authenticated() {
		if err := cloud.conn.Authenticate(); err != nil {
			return nil, fmt.Errorf("error authenticating to swift:%q", err)
		}
		defer cloud.conn.UnAuthenticate()
	}
	if len(Files) == 0 {
		return []models.Backup{}, nil
	}
	var backups []models.Backup
	for _, value := range Files {
		back, err := cloud.UploadFile(value, dstFolder)
		if err != nil {
			return nil, fmt.Errorf("Error in BackupClient:backup upload file %v", err)
		}
		backups = append(backups, *back)
	}
	return backups, nil
}
