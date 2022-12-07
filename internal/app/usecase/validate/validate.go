package validate

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"

	recog "github.com/vincentwijaya/go-ocr/pkg/face-recog"

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
	defer func() {
		go func() {
			if err := utils.RemoveLocalFile(vehiclePhotoLocation); err != nil {
				logger.Error(err)
			}
			if err := utils.RemoveLocalFile(facePhotoLocation); err != nil {
				logger.Error(err)
			}
		}()
	}()

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
	if len(faces) < 1 {
		err = errors.New("no face data found")
		return
	}

	var emptyByte [512]byte

	for _, face := range faces {
		// If face descriptor 0 value will try to recognize face
		if face.Descriptor == emptyByte {
			recog.GetFaceDescriptor(ctx)
		}
	}

	return nil
}

func downloadAndSaveFile(URL, fileName string) error {
	//Get the response bytes from the url
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return errors.New("received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return err
	}

	return nil
}
