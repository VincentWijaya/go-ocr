package validate

import (
	"context"

	"github.com/go-chi/chi/middleware"
	"github.com/vincentwijaya/go-ocr/internal/app/repo/member"
	"github.com/vincentwijaya/go-ocr/pkg/utils"
)

type ValidateUC interface {
}

type validateUC struct {
	memberRepo member.MemberRepo
}

func New(memberRepo member.MemberRepo) *validateUC {
	return &validateUC{
		memberRepo: memberRepo,
	}
}

func (uc *validateUC) ValidatePlateAndOwner(ctx context.Context, request ValidatePlateAndOwnerRequest) error {
	vehicleFileName := "./files/image/vehicle/" + middleware.GetReqID(ctx)

	if err := utils.Base64ToPNG(vehicleFileName, request.VehiclePhoto); err != nil {
		return err
	}

	_ = utils.RemoveLocalFile(vehicleFileName)

	return nil
}
