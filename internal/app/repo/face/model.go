package face

import "gorm.io/gorm"

type Face struct {
	gorm.Model
	MemberID uint
	FullName string
	PhotoURL string
}
