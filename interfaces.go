package storage

type IStorageService interface {
  UploadFile(srcPath string, dstFolder string) (*Backup, error)
}