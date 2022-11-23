package validate

import (
	"context"
	"errors"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/vincentwijaya/go-ocr/internal/app/repo/member"
	"github.com/vincentwijaya/go-ocr/pkg/log"
	"github.com/vincentwijaya/go-ocr/pkg/recognizer"
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

func (uc *validateUC) ValidatePlateAndOwner(ctx context.Context, vehiclePhotoLocation, facePhotoLocation string) (err error) {
	logger := log.WithFields(log.Fields{"request_id": middleware.GetReqID(ctx)})

	recognizeRes, err := recognizer.Recognize(vehiclePhotoLocation)
	if err != nil {
		return
	}

	logger.Infof("%+v", recognizeRes)
	if recognizeRes.Results[0].Confidence < 90 {
		err = errors.New("Please check vehicle photo")
		return
	}

	go func() {
		if err := utils.RemoveLocalFile(vehiclePhotoLocation); err != nil {
			logger.Error(err)
		}
		if err := utils.RemoveLocalFile(facePhotoLocation); err != nil {
			logger.Error(err)
		}
	}()

	return nil
}
