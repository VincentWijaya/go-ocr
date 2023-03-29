package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"path"
	"text/template"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/vincentwijaya/go-ocr/constant/errs"
	"github.com/vincentwijaya/go-ocr/internal/entity"
	"github.com/vincentwijaya/go-pkg/log"
)

func writeResponse(w http.ResponseWriter, data interface{}, err error, ctx context.Context, renderHTML bool) {
	w.Header().Set("Content-Type", "application/json")
	response := &entity.Response{
		Success: false,
	}

	if err != nil {
		response.Message, response.Code = classifyError(err, ctx)
	} else {
		response.Success = true
		response.Code = errs.Success
		response.Message = "OK"
	}

	if data != nil && response.Success {
		response.Data = data
	}

	if renderHTML {
		w.Header().Set("Content-Type", "text/html")
		var filepath = path.Join("files/views", "index.html")
		var tmpl, err = template.ParseFiles(filepath)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, response)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		return
	}

	logger := log.WithFields(log.Fields{"request_id": middleware.GetReqID(ctx)})
	logger.Infof("Response: %+v", response)
	body, _ := json.Marshal(response)

	_, _ = w.Write(body)
}

//maps
type errResponse struct {
	message string
	code    string
	apiFail bool
}

var (
	errToResponse = map[error]errResponse{
		errs.BadRequest: {
			message: errs.GeneralErrorMessage,
			code:    errs.BadRequestCode,
			apiFail: true,
		},
		errs.NoData: {
			message: errs.NoDataMessage,
			code:    errs.NoDataCode,
			apiFail: false,
		},
		errs.Unauthorized: {
			message: errs.UnauthorizedMessage,
			code:    errs.UnauthorizedErrorCode,
			apiFail: false,
		},
		errs.FaceNotFound: {
			message: errs.FaceNotFoundMessage,
			code:    errs.FaceNotFoundErrorCode,
			apiFail: false,
		},
		errs.PlateNotRecognize: {
			message: errs.PlateNotRecognizeMessage,
			code:    errs.PlateNotRecognizeErrorCode,
			apiFail: false,
		},
		errs.DriverFaceNotDetected: {
			message: errs.DriverFaceNotDetectedMessage,
			code:    errs.DriverFaceNotDetectedErrorCode,
			apiFail: false,
		},
		errs.PlatenumberNotRegistered: {
			message: errs.PlatenumberNotRegisteredMessage,
			code:    errs.PlatenumberNotRegisteredErrorCode,
			apiFail: false,
		},
	}
)

func classifyError(err error, ctx context.Context) (string, string) {
	logger := log.WithFields(log.Fields{"request_id": middleware.GetReqID(ctx)})
	val, exist := errToResponse[err]
	if !exist {
		logger.Errorf("Unmapped error:%v", err.Error())
		return errs.GeneralErrorMessage, errs.UndefinedErrorCode
	}
	if val.apiFail {
		// on api fail, return general error message
		// log the error on log
		logger.Errorf("Error on API code [%s]:%s", val.code, err.Error())
	}
	return val.message, val.code
}
