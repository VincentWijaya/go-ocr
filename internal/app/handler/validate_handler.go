package handler

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/h2non/bimg"
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

	response, err = m.validate.ValidatePlateAndOwner(ctx, "./files/images/vehicle/"+vehiclePhotoLocation, "./files/images/face/"+driverPhotoLocation)

	writeResponse(w, response, err, ctx, true)
}

func readAndSaveFile(r *http.Request, dir, formKey string) (fileLocation string, err error) {
	file, handler, err := r.FormFile(formKey)
	if err != nil {
		return
	}
	defer file.Close()

	fileName := handler.Filename

	if formKey == "driverPhoto" {
		fileName = strings.Split(handler.Filename, ".")[0]
		fileName = fileName + ".jpeg"
	}
	log.Info(dir+fileName, " ", fileName, " ", formKey)
	f, err := os.OpenFile(dir+fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return
	}

	if formKey == "vehiclePhoto" {
		data, err := ioutil.ReadAll(file)
		if err != nil {
			return fileName, err
		}

		converted, err := bimg.NewImage(data).Convert(bimg.JPEG)
		if err != nil {
			return fileName, err
		}

		processed, err := bimg.NewImage(converted).Process(bimg.Options{Quality: 70})
		if err != nil {
			return fileName, err
		}

		rotated, err := bimg.NewImage(processed).AutoRotate()
		if err != nil {
			return fileName, err
		}

		writeError := bimg.Write(fmt.Sprintf("./files/images/vehicle"+"/%s", fileName), rotated)
		if writeError != nil {
			return fileName, writeError
		}

		return fileName, nil
	}

	io.Copy(f, file)

	return fileName, nil
}
