package handler

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/vincentwijaya/go-ocr/internal/app/usecase/validate"
	"github.com/vincentwijaya/go-ocr/pkg/log"
)

func (m *Module) ValidateVehicleAndOwner(w http.ResponseWriter, r *http.Request) {
	var (
		request  validate.ValidatePlateAndOwnerRequest
		err      error
		response interface{}
	)

	ctx := r.Context()
	logger := log.WithFields(log.Fields{"request_id": middleware.GetReqID(ctx)})
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil && err.Error() != "EOF" {
		w.WriteHeader(400)
		log.Error("Failed to decode register wallet request to json")
		writeResponse(w, nil, err, ctx)
		return
	}

	logger.Infof("Request Register: %+v", request)

	err = m.validate.ValidatePlateAndOwner(ctx, request)

	writeResponse(w, response, err, ctx)
}
