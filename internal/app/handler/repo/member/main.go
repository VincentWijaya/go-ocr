package member

import "gorm.io/gorm"

type memberRepo struct {
	db *gorm.DB
}

func NewMemberRepo(db *gorm.DB) *memberRepo {
	return &memberRepo{
		db: db,
	}
}

func (mr *memberRepo) CreateMember(member Member) error {
	if result := mr.db.Create(&member); result.Error != nil {
		return result.Error
	}

	return nil
}
