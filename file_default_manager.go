package coquelicot

import (
	"os"
)

type FileDefaultManager struct {
	*FileBaseManager
	Size int64
}

func (fdm *FileDefaultManager) Convert(src string, convert string) error {
	if err := fdm.copyFile(src, fdm.Filepath()); err != nil {
		return err
	}

	f, err := os.Open(fdm.Filepath())
	if err != nil {
		return err
	}
	fi, err := f.Stat()
	if err != nil {
		return err
	}
	fdm.Size = fi.Size()

	return nil
}

func (fdm *FileDefaultManager) ToJson() map[string]interface{} {
	return map[string]interface{}{
		"url":      fdm.Url(),
		"filename": fdm.Filename,
		"size":     fdm.Size,
	}

}
