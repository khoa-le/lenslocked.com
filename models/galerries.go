package models

import "github.com/jinzhu/gorm"

type Gallery struct {
	gorm.Model
	UserID uint `gorm:"not_null,index"`
	Title string `gorm:"not_null"`
}

type GalleryService interface{
	GalleryDB
}


type GalleryDB interface{
	Create(gallery *Gallery) error
}

func NewGalleryService(db *gorm.DB) GalleryService{
	gg := &galleryGorm{db}
	gv := &GalleryValidator{gg}
	return &galleryService{
		GalleryDB: gv,
	}
}
type galleryService struct{
	GalleryDB
}

type GalleryValidator struct{
	GalleryDB
}

var _ GalleryDB = &galleryGorm{}


type galleryGorm struct{
	db *gorm.DB
}

func (gg *galleryGorm) Create(gallery *Gallery) error{
	return gg.db.Create(gallery).Error
}
