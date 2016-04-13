package macedon

import "errors"

var (
	BadConfigError       = errors.New("Bad Config")
	NoContentError       = errors.New("No Content")
	BadRequestError      = errors.New("Bad Request")
	InternalServerError  = errors.New("Interanal Server Error")
	BadGatewayError      = errors.New("Bad Gateway")
)
