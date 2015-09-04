package coquelicot

import (
	"os"
)

// Attachment contain info about directory, base mime type and all files saved.
type Attachment struct {
	OriginalFile *OriginalFile
	Dir          *DirManager
	Versions     map[string]FileManager
}

// Function receive root directory, original file, convertion parameters.
// Return Attachment saved. The final chunk is deleted if delChunk is true.
func Create(storage string, ofile *OriginalFile, converts map[string]string, delChunk bool) (*Attachment, error) {
	dm, err := CreateDir(storage, ofile.BaseMime)
	if err != nil {
		return nil, err
	}

	attachment := &Attachment{
		OriginalFile: ofile,
		Dir:          dm,
		Versions:     make(map[string]FileManager),
	}

	if ofile.BaseMime == "image" {
		converts["thumbnail"] = "120x90"
	}

	makeVersion := func(a *Attachment, version, convert string) error {
		fm, err := attachment.CreateVersion(version, convert)
		if err != nil {
			return err
		}
		attachment.Versions[version] = fm
		return nil
	}

	if err := makeVersion(attachment, "original", ""); err != nil {
		return nil, err
	}

	if makeThumbnail {
		if err := makeVersion(attachment, "thumbnail", converts["thumbnail"]); err != nil {
			return nil, err
		}
	}

	if delChunk {
		return attachment, os.Remove(attachment.OriginalFile.Filepath)
	}
	return attachment, nil
}

// Directly save single version and return FileManager.
func (attachment *Attachment) CreateVersion(version string, convert string) (FileManager, error) {
	fm := NewFileManager(attachment.Dir, attachment.OriginalFile.BaseMime, version)
	fm.SetFilename(attachment.OriginalFile)

	if err := fm.Convert(attachment.OriginalFile.Filepath, convert); err != nil {
		return nil, err
	}

	return fm, nil
}

func (attachment *Attachment) ToJson() map[string]interface{} {
	data := make(map[string]interface{})
	data["type"] = attachment.OriginalFile.BaseMime
	data["dir"] = attachment.Dir.Path
	data["name"] = attachment.OriginalFile.Filename
	versions := make(map[string]interface{})
	for version, fm := range attachment.Versions {
		versions[version] = fm.ToJson()
	}
	data["versions"] = versions

	return data
}
