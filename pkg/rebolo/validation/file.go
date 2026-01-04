package validation

import (
	"io"
	"mime/multipart"
	"os"
)

// File holds information regarding an uploaded file
type File struct {
	multipart.File
	*multipart.FileHeader
}

// Valid returns true if there is an actual uploaded file
func (f File) Valid() bool {
	return f.File != nil
}

// String returns the filename if a file is present
func (f File) String() string {
	if f.FileHeader == nil {
		return ""
	}
	return f.FileHeader.Filename
}

// Size returns the size of the uploaded file
func (f File) Size() int64 {
	if f.FileHeader == nil {
		return 0
	}
	return f.FileHeader.Size
}

// ContentType returns the content type of the uploaded file
func (f File) ContentType() string {
	if f.FileHeader == nil {
		return ""
	}
	return f.FileHeader.Header.Get("Content-Type")
}

// Save saves the uploaded file to the given path
func (f File) Save(path string) error {
	if !f.Valid() {
		return nil // Nothing to save
	}

	// Read all data from the file
	data, err := io.ReadAll(f.File)
	if err != nil {
		return err
	}
	defer f.File.Close()

	// Write to destination
	return writeFile(path, data)
}

// writeFile is a helper to write file data to disk
func writeFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}
