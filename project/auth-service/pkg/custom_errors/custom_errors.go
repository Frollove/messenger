package custom_errors

import (
	"github.com/valyala/fasthttp"
)

type ErrHttp struct {
	Code    int
	Message string
}

var (
	ErrPageNotFound   = &ErrHttp{Code: fasthttp.StatusNotFound, Message: "page not found"}
	ErrWrongInputData = &ErrHttp{Code: fasthttp.StatusUnprocessableEntity, Message: "wrong format of input data"}
	ErrInternal       = &ErrHttp{Code: fasthttp.StatusInternalServerError, Message: "something went wrong"}
	ErrWrongMethod    = &ErrHttp{Code: fasthttp.StatusMethodNotAllowed, Message: "wrong method"}
	ErrNotFound       = &ErrHttp{Code: fasthttp.StatusNotFound, Message: "no records in DB"}
	ErrUserExist      = &ErrHttp{Code: fasthttp.StatusUnprocessableEntity, Message: "user already exist"}
	ErrUnauthorized   = &ErrHttp{Code: fasthttp.StatusUnauthorized, Message: "unauthorized"}
)

func (e *ErrHttp) Error() string {
	if e == nil {
		return ""
	}

	return e.Message
}
