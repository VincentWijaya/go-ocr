package handler

import (
	"context"
)

type (
	ValidateUsecase interface {
		ValidatePlateAndOwner(ctx context.Context, vehiclePhotoLocation, facePhotoLocation string) error
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
