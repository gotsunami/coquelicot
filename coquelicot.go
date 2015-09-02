package coquelicot

type Storage struct {
	output    string
	verbosity int
}

// FIXME: global for now
var makeThumbnail bool

func (s *Storage) StorageDir() string {
	return s.output
}

func NewStorage(storageDir string) *Storage {
	return &Storage{output: storageDir}
}
