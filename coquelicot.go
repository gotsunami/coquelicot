package coquelicot

type Storage struct {
	output string
}

func (s *Storage) StorageDir() string {
	return s.output
}

func NewStorage(storageDir string) *Storage {
	return &Storage{storageDir}
}
