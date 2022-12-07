package validate

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/vincentwijaya/go-ocr/constant/errs"
	"github.com/vincentwijaya/go-ocr/internal/app/domain"
	"github.com/vincentwijaya/go-ocr/internal/app/repo/face"
	"github.com/vincentwijaya/go-ocr/internal/app/repo/vehicle"
	recog "github.com/vincentwijaya/go-ocr/pkg/face-recog"
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
	var wg sync.WaitGroup

	for i, face := range faces {
		var faceDescriptor [512]byte
		copy(faceDescriptor[:], face.Descriptor)

		// If face descriptor 0 value will try to recognize face and save face descriptor to db
		if faceDescriptor == emptyByte {
			wg.Add(1)
			memberFacePhotoLocation := fmt.Sprintf("./files/images/face/%v-%s.jpeg", face.MemberID, face.FullName)

			go func(face domain.Face, memberFacePhotoLocation string) {
				defer func() {
					go func() {
						if err := utils.RemoveLocalFile(memberFacePhotoLocation); err != nil {
							logger.Error(err)
						}
					}()
				}()
				defer wg.Done()

				if err = utils.DownloadFile(face.PhotoURL, memberFacePhotoLocation); err != nil {
					logger.Errorf("Failed to download member face: %+v", err)
					return
				}

				faceDescriptor, err := recog.GetFaceDescriptor(ctx, memberFacePhotoLocation)
				if err != nil {
					return
				}

				face.Descriptor = []byte(faceDescriptor[:])
				faces[i] = face
				err = uc.faceRepo.SaveFaceData(&face)
				if err != nil {
					return
				}
			}(face, memberFacePhotoLocation)
		}
	}

	wg.Wait()

	faceDescriptor, err := recog.GetFaceDescriptor(ctx, facePhotoLocation)
	if err != nil {
		return
	}
	faceIDResult := recog.CompareFace(faceDescriptor, faces, 0)

	if faceIDResult < 1 {
		err = errs.FaceNotFound
		return
	}

	return nil
}
