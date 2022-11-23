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
