package interfaces

import (
	"github.com/dadosjusbr/storage/models"
)

type IStorageRepository interface {
	UploadFile(srcPath string, dstFolder string) (*models.Backup, error)
} 
