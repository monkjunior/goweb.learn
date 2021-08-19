package models

import "gorm.io/gorm"

// Gallery is our image container resources
type Gallery struct {
	gorm.Model
	UserID uint   `gorm:"not_null;index"`
	Title  string `gorm:"not_null"`
}

type GalleryService interface {
	GalleryDB
}

type GalleryDB interface {
	// Methods for querying for a single gallery
	ByID(id uint) (*Gallery, error)
	ByUserID(userID uint) ([]Gallery, error)

	// Methods for altering galleries
	Create(gallery *Gallery) error
	Update(gallery *Gallery) error
	Delete(galleryID uint) error
}

func NewGalleryService(db *gorm.DB) GalleryService {
	return &galleryService{
		GalleryDB: &galleryValidator{
			GalleryDB: &galleryGorm{
				db: db,
			},
		},
	}
}

type galleryService struct {
	GalleryDB
}

type galleryValFunc func(*Gallery) error

func runGalleryValFuncs(gallery *Gallery, fns ...galleryValFunc) error {
	for _, fn := range fns {
		if err := fn(gallery); err != nil {
			return err
		}
	}
	return nil
}

type galleryValidator struct {
	GalleryDB
}

func (gv *galleryValidator) Create(gallery *Gallery) error {
	err := runGalleryValFuncs(gallery,
		gv.titleRequired,
		gv.userIDRequired,
	)
	if err != nil {
		return err
	}
	return gv.GalleryDB.Create(gallery)
}

func (gv *galleryValidator) Update(gallery *Gallery) error {
	err := runGalleryValFuncs(gallery,
		gv.titleRequired,
		gv.userIDRequired,
	)
	if err != nil {
		return err
	}
	return gv.GalleryDB.Update(gallery)
}

// Delete will delete the gallery with the provided ID
func (gv *galleryValidator) Delete(ID uint) error {
	var gallery Gallery
	gallery.ID = ID
	if err := runGalleryValFuncs(&gallery, gv.idGreaterThan(0)); err != nil {
		return err
	}
	return gv.GalleryDB.Delete(gallery.ID)
}

func (gv *galleryValidator) titleRequired(gallery *Gallery) error {
	if gallery.Title == "" {
		return ErrTitleRequired
	}
	return nil
}

func (gv *galleryValidator) userIDRequired(gallery *Gallery) error {
	if gallery.UserID <= 0 {
		return ErrUserIDRequired
	}
	return nil
}

func (gv *galleryValidator) idGreaterThan(n uint) galleryValFunc {
	return func(g *Gallery) error {
		if g.ID <= n {
			return ErrIDInvalid
		}
		return nil
	}
}

type galleryGorm struct {
	db *gorm.DB
}

// ByID will look up by the provided ID.
func (gg *galleryGorm) ByID(id uint) (*Gallery, error) {
	var gallery Gallery
	db := gg.db.Where("id = ?", id)
	err := first(db, &gallery)
	if err != nil {
		return nil, err
	}
	return &gallery, err
}

// ByUserID will list all galleries that belong to the user provided ID.
func (gg *galleryGorm) ByUserID(userID uint) ([]Gallery, error) {
	var galleries []Gallery
	if err := gg.db.Where("user_id = ?", userID).Find(&galleries).Error; err != nil {
		return nil, err
	}
	return galleries, nil
}

// Create will create the provided gallery and backfill data
// like the ID, CreatedAt, and UpdatedAt fields.
func (gg *galleryGorm) Create(gallery *Gallery) error {
	return gg.db.Create(gallery).Error
}

// Update will update the provided gallery with all of the data
// in the provided gallery object.
func (gg *galleryGorm) Update(gallery *Gallery) error {
	return gg.db.Save(gallery).Error
}

// Delete will delete the gallery with the provided ID
func (gg *galleryGorm) Delete(galleryID uint) error {
	gallery := Gallery{
		Model: gorm.Model{
			ID: galleryID,
		},
	}
	return gg.db.Delete(&gallery).Error
}
