package validate

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/vincentwijaya/go-ocr/constant/errs"
	"github.com/vincentwijaya/go-ocr/internal/app/domain"
	"github.com/vincentwijaya/go-ocr/internal/app/repo/face"
	"github.com/vincentwijaya/go-ocr/internal/app/repo/vehicle"
	recog "github.com/vincentwijaya/go-ocr/pkg/face-recog"
	"github.com/vincentwijaya/go-ocr/pkg/log"
	"github.com/vincentwijaya/go-ocr/pkg/mailer"
	"github.com/vincentwijaya/go-ocr/pkg/recognizer"
	"github.com/vincentwijaya/go-ocr/pkg/utils"
)

type ValidateUC interface {
}

type validateUC struct {
	vehicleRepo vehicle.VehicleRepo
	faceRepo    face.FaceRepo
	mailjet     mailer.MailJetClient
	secretKey   string
}

func New(vehicleRepo vehicle.VehicleRepo, faceRepo face.FaceRepo, mailjet mailer.MailJetClient, secretKey string) *validateUC {
	return &validateUC{
		vehicleRepo: vehicleRepo,
		faceRepo:    faceRepo,
		mailjet:     mailjet,
		secretKey:   secretKey,
	}
}

func (uc *validateUC) ValidatePlateAndOwner(ctx context.Context, vehiclePhotoLocation, facePhotoLocation string) (res *domain.Vehicle, err error) {
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

	recognizeRes, err := recognizer.Recognize(vehiclePhotoLocation, uc.secretKey)
	if err != nil {
		err = errs.PlateNotRecognize
		return
	}
	if len(recognizeRes.Results) < 1 {
		err = errs.PlateNotRecognize
		return
	}

	logger.Infof("%+v", recognizeRes)
	recognizeData := recognizeRes.Results[0]
	if recognizeData.Plate == "" {
		err = errs.PlateNotRecognize
		return
	}

	res, err = uc.vehicleRepo.FindVehicleByPlateNumber(strings.ToUpper(recognizeData.Plate))
	if err != nil {
		err = errs.PlatenumberNotRegistered
		return
	}

	driverFaceDescriptor, err := recog.GetFaceDescriptor(ctx, facePhotoLocation)
	if err != nil {
		err = errs.DriverFaceNotDetected
		return
	}

	logger.Infof("%+v", recog.BBytesToDescriptor(driverFaceDescriptor))

	faces, err := uc.faceRepo.FindFaceByMemberID(res.Member.ID)
	if err != nil {
		err = errs.FaceNotFound
		return
	}
	if len(faces) < 1 {
		err = errs.FaceNotFound
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

	faceIDResult := recog.CompareFace(ctx, driverFaceDescriptor, faces, 0.2)

	logger.Info(faceIDResult)
	if faceIDResult < 1 {
		err = errs.FaceNotFound
		go func() {
			uc.mailjet.SendNotifUnidentifiedFace(res.Member.Email, fmt.Sprintf("%s %s", res.Member.FirstName, res.Member.LastName), res.PlateNumber)
		}()
		return
	}

	return
}
