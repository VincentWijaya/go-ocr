package handler

import (
	"encoding/json"
	"net/http"

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
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil && err.Error() != "EOF" {
		w.WriteHeader(400)
		log.Error("Failed to decode register wallet request to json")
		writeResponse(w, nil, err)
		return
	}

	log.Infof("Request Register: %+v", request)

	err = m.validate.ValidatePlateAndOwner(ctx, request)

	writeResponse(w, response, err)
}
