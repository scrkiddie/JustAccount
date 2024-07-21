package adapter

import (
	"awesomeProject12/internal/model"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type FileAdapter struct {
}

func NewFileAdapter() *FileAdapter {
	return &FileAdapter{}
}

func (s *FileAdapter) StoreFile(storePath string, f *model.File) (*model.File, error) {

	file, err := f.FileHeader.Open()
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer file.Close()
	uniqueFilename := uuid.New().String() + filepath.Ext(f.FileHeader.Filename)
	fullPath := filepath.Join(storePath, uniqueFilename)

	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		return nil, err
	}
	return &model.File{Name: uniqueFilename}, nil
}

func (s *FileAdapter) DeleteFile(storePath string, f *model.File) error {
	filePath := filepath.Join(storePath, f.Name)

	var err error
	for attempt := 0; attempt < 5; attempt++ {
		err = os.Remove(filePath)
		if err == nil {
			return nil
		}
		log.Printf("attempt %d: failed to delete file %s: %s\n", attempt+1, filePath, err.Error())
		delay := 5 * time.Second * time.Duration(1<<attempt)
		time.Sleep(delay)
	}

	return fmt.Errorf("failed to delete file %s after multiple attempts: %s", filePath, err.Error())
}

func (s *FileAdapter) DeleteFileAsync(storePath string, f *model.File) {
	go func() {
		err := s.DeleteFile(storePath, f)
		if err != nil {
			log.Printf("async file deletion failed for %s: %v\n", f.Name, err)
		}
	}()
}
