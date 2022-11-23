package vehicle

import (
	"github.com/vincentwijaya/go-ocr/internal/app/domain"
	"gorm.io/gorm"
)

type VehicleRepo struct {
	db *gorm.DB
}

func NewVehicleRepo(db *gorm.DB) *VehicleRepo {
	return &VehicleRepo{
		db: db,
	}
}

func (vr *VehicleRepo) FindVehicleByPlateNumber(plateNumber string) (res *domain.Vehicle, err error) {
	err = vr.db.Where("plate_number = ?", plateNumber).Preload("Member").First(&res).Error

	return
}
