package member

import (
	"github.com/vincentwijaya/go-ocr/internal/app/repo/face"
	"github.com/vincentwijaya/go-ocr/internal/app/repo/vehicle"
	"gorm.io/gorm"
)

type Member struct {
	gorm.Model
	FirstName   string
	LastName    string
	PhoneNumber string
	Email       string
	Vehicles    []vehicle.Vehicle
	Faces       []face.Face
}
