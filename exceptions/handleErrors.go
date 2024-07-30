package exceptions

import (
	"fmt"
	"net/http"
)

type ErrorResponse struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message"`
}

func HandleExceptions(err error) ErrorResponse {
	customErr, ok := err.(*Error)
	if !ok {
		fmt.Println(err.Error())
		if err.Error() == "code=404, message=Not Found" {
			return ErrorResponse{
				Message: "Not Found",
			}
		}
		return ErrorResponse{
			Code:    500,
			Message: "Internal server error",
		}
	}

	switch customErr.Err {
	case ErrChannelNotFound:
		return ErrorResponse{
			Code:    http.StatusNotFound,
			Message: customErr.Err.Error(),
		}
	case ErrInvalidPayload:
		return ErrorResponse{
			Code:    http.StatusUnprocessableEntity,
			Message: customErr.Err.Error(),
		}
	case
		ErrChannelAlreadyExists,
		ErrInvalidID,
		ErrInvalidNameField,
		ErrInvalidMembersField,
		ErrInvalidAdminsField,
		ErrInvalidUserIdSent:
		return ErrorResponse{
			Code:    http.StatusBadRequest,
			Message: customErr.Err.Error(),
		}
	case ErrDatabaseFailure:
		return ErrorResponse{
			Code:    http.StatusInternalServerError,
			Message: customErr.Err.Error(),
		}
	default:
		return ErrorResponse{
			Code:    500,
			Message: "Internal server error",
		}
	}
}
