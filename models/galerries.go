package models

import (
	"github.com/jinzhu/gorm"
)

type Gallery struct {
	gorm.Model
	UserID uint     `gorm:"not_null,index"`
	Title  string   `gorm:"not_null"`
	Images []Image `gorm:"-"`
}

func (g *Gallery) ImageSplitN(n int) [][]Image {
	ret := make([][]Image, n)
	for i := 0; i < n; i++ {
		ret[i] = make([]Image, 0)
	}
	for i, img := range g.Images {
		bucket := i % n
		ret[bucket] = append(ret[bucket], img)
	}
	return ret
}

type GalleryService interface {
	GalleryDB
}

type GalleryDB interface {
	ByID(id uint) (*Gallery, error)
	Create(gallery *Gallery) error
	Update(gallery *Gallery) error
	Delete(id uint) error
	ByUserID(userID uint) ([]Gallery, error)
}

func NewGalleryService(db *gorm.DB) GalleryService {
	gg := &galleryGorm{db}
	gv := &galleryValidator{gg}
	return &galleryService{
		GalleryDB: gv,
	}
}

type galleryService struct {
	GalleryDB
}

type galleryValidator struct {
	GalleryDB
}

func (gv *galleryValidator) Create(gallery *Gallery) error {
	if err := runGalleryValidatorFunction(gallery,
		gv.titleRequired,
		gv.userIDRequired,
	); err != nil {
		return err
	}

	return gv.GalleryDB.Create(gallery)
}

func (gv *galleryValidator) Update(gallery *Gallery) error {
	if err := runGalleryValidatorFunction(gallery,
		gv.titleRequired,
		gv.userIDRequired, ); err != nil {
		return err
	}
	return gv.GalleryDB.Update(gallery)
}

func (gv *galleryValidator) Delete(id uint) error {
	gallery := &Gallery{Model: gorm.Model{ID: id}}
	if err := runGalleryValidatorFunction(gallery,
		gv.idGreaterThan(0), ); err != nil {
		return err
	}
	return gv.GalleryDB.Delete(id)
}

func (gv *galleryValidator) userIDRequired(gallery *Gallery) error {
	if gallery.UserID == 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (gv *galleryValidator) titleRequired(gallery *Gallery) error {
	if gallery.Title == "" {
		return ErrGalleryTitleRequired
	}
	return nil
}

func (gv *galleryValidator) idGreaterThan(n uint) galleryValidatorFunction {
	return galleryValidatorFunction(func(gallery *Gallery) error {
		if gallery.ID <= n {
			return ErrIDInvalid
		}
		return nil
	})
}

var _ GalleryDB = &galleryGorm{}

type galleryGorm struct {
	db *gorm.DB
}

func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}

func (gg *galleryGorm) Update(gallery *Gallery) error {
	return gg.db.Save(gallery).Error
}

func (gg *galleryGorm) ByID(id uint) (*Gallery, error) {
	var gallery Gallery
	db := gg.db.Where("id=?", id)
	err := first(db, &gallery)
	return &gallery, err
}

func (gg *galleryGorm) Delete(id uint) error {
	gallery := Gallery{Model: gorm.Model{ID: id}}
	return gg.db.Delete(gallery).Error
}

func (gg *galleryGorm) ByUserID(userID uint) ([]Gallery, error) {
	var galleries []Gallery
	gg.db.Where("user_id=?", userID).Find(&galleries)
	return galleries, nil
}

type galleryValidatorFunction func(*Gallery) error

func runGalleryValidatorFunction(gallery *Gallery, fns ...galleryValidatorFunction) error {
	for _, fn := range fns {
		if err := fn(gallery); err != nil {
			return err
		}
	}
	return nil
}
