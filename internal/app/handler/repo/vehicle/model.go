package vehicle

import "gorm.io/gorm"

type Vehicle struct {
	gorm.Model
	MemberID    uint
	PlateNumber string
	Brand       string
	ModelName   string
}
