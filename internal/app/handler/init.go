package handler

import (
	"context"

	"github.com/vincentwijaya/go-ocr/internal/app/domain"
)

type (
	ValidateUsecase interface {
		ValidatePlateAndOwner(ctx context.Context, vehiclePhotoLocation, facePhotoLocation string) (*domain.Vehicle, error)
	}
)

type Module struct {
	validate ValidateUsecase
}

func New(validate ValidateUsecase) *Module {
	return &Module{
		validate: validate,
	}
}
