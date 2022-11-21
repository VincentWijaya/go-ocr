package member

import "gorm.io/gorm"

type MemberRepo struct {
	db *gorm.DB
}

func NewMemberRepo(db *gorm.DB) *MemberRepo {
	return &MemberRepo{
		db: db,
	}
}

func (mr *MemberRepo) CreateMember(member Member) error {
	if result := mr.db.Create(&member); result.Error != nil {
		return result.Error
	}

	return nil
}
