package vehicle

import "gorm.io/gorm"

type vehicleRepo struct {
	db *gorm.DB
}

func NewVehicleRepo(db *gorm.DB) *vehicleRepo {
	return &vehicleRepo{
		db: db,
	}
}

func (vr *vehicleRepo) FindVehicleByPlateNumber(plateNumber string) (res *Vehicle, err error) {
	err = vr.db.Where("plate_number = ?", plateNumber).First(&res).Error

	return
}
