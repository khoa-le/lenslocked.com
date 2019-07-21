package models

import (
	"github.com/jinzhu/gorm"
)

type Gallery struct {
	gorm.Model
	UserID uint `gorm:"not_null,index"`
	Title string `gorm:"not_null"`
}

type GalleryService interface{
	GalleryDB
}


type GalleryDB interface{
	ByID(id uint) (*Gallery, error)
	Create(gallery *Gallery) error
}

func NewGalleryService(db *gorm.DB) GalleryService{
	gg := &galleryGorm{db}
	gv := &galleryValidator{gg}
	return &galleryService{
		GalleryDB: gv,
	}
}
type galleryService struct{
	GalleryDB
}

type galleryValidator struct{
	GalleryDB
}

func (gv *galleryValidator) Create(gallery *Gallery) error{
	if err:=runGalleryValidatorFunction(gallery,
		gv.titleRequired,
		gv.userIDRequired,
		); err!=nil{
		return err
	}

	return gv.GalleryDB.Create(gallery)
}

func (gv *galleryValidator) userIDRequired(gallery *Gallery ) error{
	if gallery.UserID == 0{
		return ErrUserIDRequired
	}
	return nil
}

func (gv *galleryValidator) titleRequired(gallery *Gallery ) error{
	if gallery.Title == ""{
		return ErrGalleryTitleRequired
	}
	return nil
}


var _ GalleryDB = &galleryGorm{}

type galleryGorm struct{
	db *gorm.DB
}

func (gg *galleryGorm) Create(gallery *Gallery) error{
	return gg.db.Create(gallery).Error
}

func (gg *galleryGorm) ByID(id uint) (*Gallery, error){
	var gallery Gallery
	db := gg.db.Where("id=?", id)
	err := first(db, &gallery)
	return &gallery, err

}


type galleryValidatorFunction func(*Gallery) error

func runGalleryValidatorFunction(gallery *Gallery, fns ...galleryValidatorFunction) error{
	for _,fn :=range fns {
		if err := fn(gallery); err!=nil{
			return err
		}
	}
	return nil
}
