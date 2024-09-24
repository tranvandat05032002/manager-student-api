package Common

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

/*
	{
		status: number,
		message: ".....",
		error: error
	}
*/
type Type string

const (
	Authorization Type = "AUTHORIZATION"
	BadRequest    Type = "BADREQUEST"
	Conflict      Type = "CONFLICT"
	Internal      Type = "INTERNAL"
	NotFound      Type = "NOTFOUND"
)

type ErrorRes struct {
	Status  int         `json:"status"`
	Message string      `json:"message,omitempty"`
	Error   interface{} `json:"error,omitempty"`
}

func NewErrorResponse(c *gin.Context, status int, message string, err interface{}) {
	errorResponse := ErrorRes{
		Status:  status,
		Message: message,
		Error:   err,
	}
	c.JSON(status, errorResponse)
}

type Error struct {
	Type    Type   `json:"type"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

// status is a mapping errors to status codes
func (e *Error) Status() int {
	switch e.Type {
	case Authorization:
		return http.StatusUnauthorized
	case BadRequest:
		return http.StatusBadRequest
	case Conflict:
		return http.StatusConflict
	case Internal:
		return http.StatusInternalServerError
	case NotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

// status checks the runtime type
func Status(err error) int {
	var e *Error
	if errors.As(err, &e) {
		return e.Status()
	}
	return http.StatusInternalServerError
}

/*
* Error "Factories
 */
// NewAuthorization to create a 401
func NewAuthorization(reason string) *Error {
	return &Error{
		Type:    Authorization,
		Message: reason,
	}
}

// NewBadRequest to create 400 errors (validation, for example)
func NewBadRequest(reason string) *Error {
	return &Error{
		Type:    BadRequest,
		Message: fmt.Sprintf("Yêu cầu không hợp lệ: %v", reason),
	}
}

// NewConflict to create an error for 409
func NewConflict(name string, value string) *Error {
	return &Error{
		Type:    Conflict,
		Message: fmt.Sprintf("resource: %v with value: %v already exists", name, value),
	}
}

// NewInternal for 500 errors and unknown errors
func NewInternal() *Error {
	return &Error{
		Type:    Internal,
		Message: fmt.Sprintf("Lỗi hệ thống!"),
	}
}

// NewNotFound to create an error for 404
func NewNotFound(name string, value string) *Error {
	return &Error{
		Type:    NotFound,
		Message: fmt.Sprintf("Tài nguyên: %v vơ giá trị: %v không tìm thấy", name, value),
	}
}
