package handler

import (
	"context"

	"github.com/vincentwijaya/go-ocr/internal/app/usecase/validate"
)

type (
	ValidateUsecase interface {
		ValidatePlateAndOwner(ctx context.Context, request validate.ValidatePlateAndOwnerRequest) error
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
