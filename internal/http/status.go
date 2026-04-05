package http

import (
	"net/http"

	"github.com/moises-ba/mb-crypto-sma-api/internal/errors"
)

var mapKindHttpStatus = map[errors.Kind]int{
	errors.NotFound:      http.StatusNotFound,
	errors.Invalid:       http.StatusBadRequest,
	errors.Internal:      http.StatusInternalServerError,
	errors.Unexpected:    http.StatusInternalServerError,
	errors.Unprocessable: http.StatusUnprocessableEntity,
}

func EvaluateStatus(err errors.ApiError) int {
	if err == nil {
		return http.StatusOK
	}

	if httpStatus, ok := mapKindHttpStatus[err.Kind()]; ok {
		return httpStatus
	}

	return http.StatusInternalServerError
}
