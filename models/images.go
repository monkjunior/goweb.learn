package models

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type ImageService interface {
	Create(galleryID uint, r io.ReadCloser, filename string) error
	ByGalleryID(galleryID uint) ([]string, error)
}

func NewImageService() ImageService {
	return &imageService{}
}

type imageService struct{}

func (i *imageService) Create(galleryID uint, r io.ReadCloser, filename string) error {
	defer r.Close()
	path, err := i.mkImagePath(galleryID)
	if err != nil {
		return err
	}
	dst, err := os.Create(path + filename)
	if err != nil {
		return err
	}
	_, err = io.Copy(dst, r)
	if err != nil {
		return err
	}
	return nil
}

func (i *imageService) ByGalleryID(galleryID uint) ([]string, error) {
	path := i.imagePath(galleryID)
	strings, err := filepath.Glob(path + "*")
	if err != nil {
		return nil, err
	}
	// Add a leading "/" to all image file paths
	for i := range strings {
		strings[i] = "/" + strings[i]
	}
	return strings, nil
}

func (i *imageService) imagePath(galleryID uint) string {
	return fmt.Sprintf("images/galleries/%v/", galleryID)
}

func (i *imageService) mkImagePath(galleryID uint) (string, error) {
	galleryPath := i.imagePath(galleryID)
	err := os.MkdirAll(galleryPath, 0755)
	if err != nil {
		return "", err
	}
	return galleryPath, nil
}
