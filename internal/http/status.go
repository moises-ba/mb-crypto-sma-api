package http

import (
	"net/http"

	"github.com/moises-ba/mb-crypto-mms-api/internal/errors"
)

var mapKindHttpStatus = map[errors.Kind]int{
	errors.NotFound:   http.StatusNotFound,
	errors.Invalid:    http.StatusBadRequest,
	errors.Internal:   http.StatusInternalServerError,
	errors.Unexpected: http.StatusInternalServerError,
}

func EvaluateStatus(err errors.ApiError) int {
	if err == nil {
		return http.StatusOK
	}

	return mapKindHttpStatus[err.Kind()]
}
