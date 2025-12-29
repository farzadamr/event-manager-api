package helper

import (
	"net/http"

	"github.com/farzadamr/event-manager-api/pkg/service_errors"
)

var StatusCodeMapping = map[string]int{
	//User
	service_errors.EmailExists:         409,
	service_errors.StudentNumberExists: 409,
	service_errors.RecordNotFound:      404,
	service_errors.PermissionDenied:    403,
}

func TranslateErrorToStatusCode(err error) int {
	value, ok := StatusCodeMapping[err.Error()]
	if !ok {
		return http.StatusInternalServerError
	}
	return value
}
