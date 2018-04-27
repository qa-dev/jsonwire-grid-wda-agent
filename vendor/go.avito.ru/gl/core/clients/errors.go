package clients

import (
	"fmt"
	"net/http"
)

// InvalidStatusError задает ошибку когда сервис вернул http-код, отличный от 200.
type InvalidStatusError struct {
	code    int
	message string
}

func (e *InvalidStatusError) String() string {
	return fmt.Sprintf("code: %v, message: %v", e.code, e.message)
}

func (e *InvalidStatusError) Error() string {
	return e.String()
}

// Code возвращает пришедший от сервиса http-код.
func (e *InvalidStatusError) Code() int {
	return e.code
}

// Message возвращает пришедший от сервиса http-статус.
func (e *InvalidStatusError) Message() string {
	return e.message
}

// NewInvalidStatusError возвращает указательн на новый InvalidStatusError с переданным http-кодом.
func NewInvalidStatusError(code int) *InvalidStatusError {
	return &InvalidStatusError{
		code:    code,
		message: http.StatusText(code),
	}
}

var (
	// ErrBadRequest описывает ошибку когда сервис вернул код 400.
	ErrBadRequest = NewInvalidStatusError(http.StatusBadRequest)
	// ErrUnauthorized описывает ошибку когда сервис вернул код 401.
	ErrUnauthorized = NewInvalidStatusError(http.StatusUnauthorized)
	// ErrForbidden описывает ошибку когда сервис вернул код 403.
	ErrForbidden = NewInvalidStatusError(http.StatusForbidden)
	// ErrNotFound описывает ошибку когда сервис вернул код 404.
	ErrNotFound = NewInvalidStatusError(http.StatusNotFound)
	// ErrServerError описывает ошибку когда сервис вернул код 500.
	ErrServerError = NewInvalidStatusError(http.StatusInternalServerError)
	// ErrGatewayTimeout описывает ошибку когда сервис не ответил по timeout.
	ErrGatewayTimeout = NewInvalidStatusError(http.StatusGatewayTimeout)
)

// StatusCodeToErrorMap задает отображение http-кодов на соответствующие ошибки.
var StatusCodeToErrorMap = map[int]*InvalidStatusError{
	http.StatusBadRequest:          ErrBadRequest,
	http.StatusUnauthorized:        ErrUnauthorized,
	http.StatusForbidden:           ErrForbidden,
	http.StatusNotFound:            ErrNotFound,
	http.StatusInternalServerError: ErrServerError,
	http.StatusGatewayTimeout:      ErrGatewayTimeout,
}

// GetInvalidStatusErrorByCode возвращает ошибку по http-коду.
func GetInvalidStatusErrorByCode(code int) *InvalidStatusError {
	err, ok := StatusCodeToErrorMap[code]

	if !ok {
		return ErrServerError
	}

	return err
}

// InvalidBodyError задает ошибку, связанную с телом http-ответа.
type InvalidBodyError struct {
	message string
	reason  error
}

func (e *InvalidBodyError) String() string {
	return fmt.Sprintf("message: %v, reason: %v", e.message, e.reason)
}

func (e *InvalidBodyError) Error() string {
	return e.String()
}

// Message возвращает текстовое описание ошибки.
func (e *InvalidBodyError) Message() string {
	return e.message
}

// Reason возвращает оригинальную ошибку.
func (e *InvalidBodyError) Reason() error {
	return e.reason
}

// NewInvalidBodyError возвращает новую проинициализированную InvalidBodyError.
func NewInvalidBodyError(message string, reason error) *InvalidBodyError {
	return &InvalidBodyError{
		message: message,
		reason:  reason,
	}
}

// GetStatusByError возвращает http status code и message, соответствующие переданной ошибке.
func GetStatusByError(err error) (int, string) {
	if err != nil {
		switch err.(type) {
		case *InvalidStatusError:
			invalidStatusErr := err.(*InvalidStatusError)
			if invalidStatusErr != nil {
				return invalidStatusErr.Code(), invalidStatusErr.Message()
			}

		default:
			return http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError)
		}
	}

	return http.StatusOK, http.StatusText(http.StatusOK)
}
