package validate

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/vincentwijaya/go-ocr/internal/app/repo/face"
	"github.com/vincentwijaya/go-ocr/internal/app/repo/vehicle"
	"github.com/vincentwijaya/go-ocr/pkg/log"
	"github.com/vincentwijaya/go-ocr/pkg/recognizer"
	"github.com/vincentwijaya/go-ocr/pkg/utils"
)

type ValidateUC interface {
}

type validateUC struct {
	vehicleRepo vehicle.VehicleRepo
	faceRepo    face.FaceRepo
}

func New(vehicleRepo vehicle.VehicleRepo, faceRepo face.FaceRepo) *validateUC {
	return &validateUC{
		vehicleRepo: vehicleRepo,
		faceRepo:    faceRepo,
	}
}

func (uc *validateUC) ValidatePlateAndOwner(ctx context.Context, vehiclePhotoLocation, facePhotoLocation string) (err error) {
	logger := log.WithFields(log.Fields{"request_id": middleware.GetReqID(ctx)})

	recognizeRes, err := recognizer.DirectRecognize(vehiclePhotoLocation)
	if err != nil {
		return
	}

	logger.Infof("%+v", recognizeRes)
	recognizeData := recognizeRes.Plates[0]
	if recognizeData.BestPlate == "" {
		err = errors.New("plate not recognize")
		return
	}

	res, err := uc.vehicleRepo.FindVehicleByPlateNumber(recognizeData.BestPlate)
	if err != nil {
		return
	}

	faces, err := uc.faceRepo.FindFaceByMemberID(res.Member.ID)
	if err != nil {
		return
	}

	fmt.Println("%+v", faces)

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
