package models

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"net/url"
)

type Image struct {
	GalleryID uint
	FileName  string
}

func (i *Image) RelativePath() string{
	return fmt.Sprintf("images/galleries/%v/%v", i.GalleryID, i.FileName)
}

func (i *Image) Path() string {
	url := url.URL{
		Path: "/" + i.RelativePath(),
	}

	return url.String()
}

type ImageService interface {
	Create(galleryID uint, r io.ReadCloser, filename string) error
	ByGalleryID(galleryID uint) ([]Image, error)
	Delete(i *Image) error
}

func NewImageService() ImageService {
	return &imageService{}
}

type imageService struct{}

func (is *imageService) Create(galleryID uint, r io.ReadCloser, filename string) error {
	defer r.Close()
	path, err := is.makeImagePath(galleryID)
	if err != nil {
		return err
	}
	dst, err := os.Create(path + filename)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, r)
	if err != nil {
		return err
	}

	return nil
}
func (is *imageService) ByGalleryID(galleryID uint) ([]Image, error) {
	path := is.imagePath(galleryID)
	imageStrings, err := filepath.Glob(path + "*")
	if err != nil {
		return nil, err
	}

	ret := make([]Image, len(imageStrings))

	for i := range imageStrings {
		imageStrings[i] = strings.Replace(imageStrings[i], path, "", 1)
		ret[i] = Image{GalleryID: galleryID, FileName: imageStrings[i]}
	}
	return ret, nil

}

func (is *imageService) Delete(i *Image) error{

	return os.Remove(i.RelativePath())
}

func (is *imageService) imagePath(galleryID uint) string {
	return fmt.Sprintf("images/galleries/%v/", galleryID)
}

func (is *imageService) makeImagePath(galleryID uint) (string, error) {
	galleryPath := is.imagePath(galleryID)
	err := os.MkdirAll(galleryPath, 0755)
	if err != nil {
		return "", err
	}
	return galleryPath, nil
}
