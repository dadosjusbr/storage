package file_storage

import (
	"github.com/dadosjusbr/storage/models"
)

type Interface interface {
	UploadFile(srcPath string, dstFolder string) (*models.Backup, error)
	GetFile(dstFolder string) (*models.Backup, error)
}
