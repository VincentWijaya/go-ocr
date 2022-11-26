package domain

import "gorm.io/gorm"

type Face struct {
	gorm.Model
	MemberID   uint
	FullName   string
	PhotoURL   string
	Descriptor []byte
}
