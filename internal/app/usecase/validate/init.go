package validate

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/vincentwijaya/go-ocr/internal/app/repo/member"
	"github.com/vincentwijaya/go-ocr/pkg/log"
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

	pwd, _ := os.Getwd()
	logger.Infof(fmt.Sprintln("docker", "run", "--rm", "-v", pwd+"/files/images/vehicle:/data:ro", "openalpr/openalpr", "-j", "-c", "sg", path.Base(vehiclePhotoLocation)))
	cmd := exec.Command("docker", "run", "--rm", "--platform", "linux/amd64", "-v", pwd+"/files/images/vehicle:/data:ro", "openalpr/openalpr", "-j", "-c", "sg", path.Base(vehiclePhotoLocation))
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err = cmd.Run()
	if err != nil {
		err = errors.New(stderr.String())
		return
	}

	logger.Infof(stdout.String())

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
