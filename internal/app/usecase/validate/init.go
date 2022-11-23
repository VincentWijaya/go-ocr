package validate

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/vincentwijaya/go-ocr/internal/app/repo/vehicle"
	"github.com/vincentwijaya/go-ocr/pkg/log"
	"github.com/vincentwijaya/go-ocr/pkg/recognizer"
	"github.com/vincentwijaya/go-ocr/pkg/utils"
)

type ValidateUC interface {
}

type validateUC struct {
	vehicleRepo vehicle.VehicleRepo
}

func New(vehicleRepo vehicle.VehicleRepo) *validateUC {
	return &validateUC{
		vehicleRepo: vehicleRepo,
	}
}

func (uc *validateUC) ValidatePlateAndOwner(ctx context.Context, vehiclePhotoLocation, facePhotoLocation string) (err error) {
	logger := log.WithFields(log.Fields{"request_id": middleware.GetReqID(ctx)})

	recognizeRes, err := recognizer.Recognize(vehiclePhotoLocation)
	if err != nil {
		return
	}

	logger.Infof("%+v", recognizeRes)
	recognizeData := recognizeRes.Results[0]
	if recognizeData.Confidence < 90 {
		err = errors.New("Please check vehicle photo")
		return
	}

	res, err := uc.vehicleRepo.FindVehicleByPlateNumber(recognizeData.Plate)
	if err != nil {
		return
	}

	fmt.Println("%+v", res.Member.ID)

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
