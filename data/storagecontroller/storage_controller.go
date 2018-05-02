package storagecontroller

import (
	"fmt"
	"os"

	"github.com/KyberNetwork/reserve-data/common/archive"
)

type StorageController struct {
	Runner StorageControllerRunner
	Arch   archive.Archive
}

func NewStorageController(storageControllerRunner StorageControllerRunner, arch archive.Archive) (StorageController, error) {
	storageController := StorageController{
		storageControllerRunner, arch,
	}
	return storageController, nil
}

func (self StorageController) BackupFile(fileName, remoteFolderPath, remoteBucketName string) error {
	err := self.Arch.UploadFile(remoteFolderPath, fileName, remoteBucketName)
	if err != nil {
		return err
	}
	intergrity, err := self.Arch.CheckFileIntergrity(remoteFolderPath, fileName, remoteFolderPath)
	if err != nil {
		return err
	}
	if intergrity {
		return os.Remove(fileName)
	} else {
		return fmt.Errorf("StorageController: Upload File  %s: corrupted", fileName)
	}

	return nil
}
