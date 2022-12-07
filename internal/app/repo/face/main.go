package face

import (
	"github.com/vincentwijaya/go-ocr/internal/app/domain"
	"gorm.io/gorm"
)

type FaceRepo struct {
	db *gorm.DB
}

func NewFaceRepo(db *gorm.DB) *FaceRepo {
	return &FaceRepo{
		db: db,
	}
}

func (fr *FaceRepo) FindFaceByMemberID(memberID uint) (res []domain.Face, err error) {
	result := fr.db.Where("member_id = ?", memberID).Find(&res)

	if result.Error != nil {
		err = result.Error
		return
	}

	return
}

func (fr *FaceRepo) SaveFaceData(face *domain.Face) (err error) {
	err = fr.db.Save(face).Error

	return
}
