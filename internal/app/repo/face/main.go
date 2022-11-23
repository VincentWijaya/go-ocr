package face

import "gorm.io/gorm"

type FaceRepo struct {
	db *gorm.DB
}

func NewFaceRepo(db *gorm.DB) *FaceRepo {
	return &FaceRepo{
		db: db,
	}
}

func (fr *FaceRepo) FindFaceByUserID(userID uint) (res *[]FaceRepo, err error) {
	result := fr.db.Where("user_id = ?", userID).Find(&res)

	if result.Error != nil {
		err = result.Error
		return
	}

	return
}
