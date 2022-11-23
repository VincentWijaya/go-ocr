package handler

import (
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/vincentwijaya/go-ocr/pkg/log"
)

func (m *Module) ValidateVehicleAndOwner(w http.ResponseWriter, r *http.Request) {
	var (
		err      error
		response interface{}
	)

	r.ParseMultipartForm(32 << 20)
	ctx := r.Context()
	logger := log.WithFields(log.Fields{"request_id": middleware.GetReqID(ctx)})

	logger.Infof("Request Register: %+v", nil)

	var wg sync.WaitGroup
	wg.Add(2)
	var vehiclePhotoLocation, driverPhotoLocation string

	go func(r *http.Request, dir, formKey string) {
		fileLocation, err := readAndSaveFile(r, dir, formKey)
		if err != nil {
			logger.Error(err)
		}

		vehiclePhotoLocation = fileLocation
		wg.Done()
	}(r, "./files/images/vehicle/", "vehiclePhoto")

	go func(r *http.Request, dir, formKey string) {
		fileLocation, err := readAndSaveFile(r, dir, formKey)
		if err != nil {
			logger.Error(err)
		}

		driverPhotoLocation = fileLocation
		wg.Done()
	}(r, "./files/images/face/", "driverPhoto")

	wg.Wait()

	err = m.validate.ValidatePlateAndOwner(ctx, "./files/images/vehicle/"+vehiclePhotoLocation, "./files/images/face/"+driverPhotoLocation)

	writeResponse(w, response, err, ctx)
}

func readAndSaveFile(r *http.Request, dir, formKey string) (fileLocation string, err error) {
	file, handler, err := r.FormFile(formKey)
	if err != nil {
		return
	}
	defer file.Close()

	f, err := os.OpenFile(dir+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return
	}
	io.Copy(f, file)
	return handler.Filename, nil
}
