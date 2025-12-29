package helper

import "github.com/farzadamr/event-manager-api/api/validation"

type BaseHttpResponse struct {
	Result           any                           `json:"result"`
	Success          bool                          `json:"success"`
	ValidationErrors *[]validation.ValidationError `json:"validationErrors"`
	Error            any                           `json:"error"`
}

func GenerateBaseResponse(result any, success bool) *BaseHttpResponse {
	return &BaseHttpResponse{
		Success: success,
		Result:  result,
	}
}

func GenerateBaseResponseWithError(result any, success bool, err error) *BaseHttpResponse {
	return &BaseHttpResponse{
		Result:  result,
		Success: success,
		Error:   err.Error(),
	}
}

func GenerateBaseResponseWithAnyError(result any, success bool, err any) *BaseHttpResponse {
	return &BaseHttpResponse{
		Result:  result,
		Success: success,
		Error:   err,
	}
}

func GenerateBaseResponseWithValidationError(result any, success bool, err error) *BaseHttpResponse {
	return &BaseHttpResponse{
		Result:           result,
		Success:          success,
		ValidationErrors: validation.GetValidationErrors(err),
	}
}
