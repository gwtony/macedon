package errors

import "errors"

var (
	BadConfigError       = errors.New("Bad Config")

	ForbiddenError       = errors.New("Forbidden")
	NoContentError       = errors.New("No Content")
	BadRequestError      = errors.New("Bad Request")
	InternalServerError  = errors.New("Interanal Server Error")
	BadGatewayError      = errors.New("Bad Gateway")
)
