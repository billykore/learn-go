package learning

import (
	"fmt"
	"net/http"
	"testing"
)

type response struct {
	code    int
	message string
	data    any
}

type exceptionKind uint8

const (
	exceptionKindUnauthorized exceptionKind = iota + 1
	exceptionKindBadRequest
	exceptionKindNotFound
	exceptionKindServiceUnavailable
)

const (
	unauthorizedException       = "unauthorized"
	badRequestException         = "bad request"
	notFoundException           = "not found"
	serviceUnavailableException = "service unavailable"
)

type exception struct {
	kind    exceptionKind
	message string
}

var _ error = &exception{}

func (e *exception) Kind() exceptionKind {
	if e != nil {
		return e.kind
	}
	return 0
}

func (e *exception) Error() string {
	if e != nil {
		return e.message
	}
	return ""
}

type service struct{}

func (s *service) success() (*response, error) {
	return &response{}, nil
}

func (s *service) badRequest() (*response, error) {
	return nil, &exception{exceptionKindBadRequest, badRequestException}
}

func (s *service) unauthorized() (*response, error) {
	return nil, &exception{exceptionKindUnauthorized, unauthorizedException}
}

func (s *service) notFound() (*response, error) {
	return nil, &exception{exceptionKindNotFound, notFoundException}
}

func (s *service) serviceUnavailable() (*response, error) {
	return nil, &exception{exceptionKindServiceUnavailable, serviceUnavailableException}
}

func sendSuccessResponse(data any) *response {
	return &response{
		code:    http.StatusOK,
		message: "success",
		data:    data,
	}
}

func unauthorizedResponse() *response {
	return &response{
		code:    http.StatusUnauthorized,
		message: "unauthorized",
		data:    nil,
	}
}

func badRequestResponse() *response {
	return &response{
		code:    http.StatusBadRequest,
		message: "bad request",
		data:    nil,
	}
}

func notFoundResponse() *response {
	return &response{
		code:    http.StatusNotFound,
		message: "not found",
		data:    nil,
	}
}

func serviceUnavailableResponse() *response {
	return &response{
		code:    http.StatusServiceUnavailable,
		message: "service unavailable",
		data:    nil,
	}
}

func generalErrorResponse() *response {
	return &response{
		code:    http.StatusInternalServerError,
		message: "internal server error",
		data:    nil,
	}
}

func sendErrorResponse(err error) *response {
	switch err.(*exception).Kind() {
	case exceptionKindUnauthorized:
		return unauthorizedResponse()
	case exceptionKindBadRequest:
		return badRequestResponse()
	case exceptionKindNotFound:
		return notFoundResponse()
	case exceptionKindServiceUnavailable:
		return serviceUnavailableResponse()
	default:
		return generalErrorResponse()
	}
}

var myService = new(service)

func TestSuccess(t *testing.T) {
	res, err := myService.success()
	if err != nil {
		fmt.Println(sendErrorResponse(err))
	} else {
		fmt.Println(sendSuccessResponse(res))
	}
}

func TestUnauthorized(t *testing.T) {
	res, err := myService.unauthorized()
	if err != nil {
		fmt.Println(sendErrorResponse(err))
	} else {
		fmt.Println(sendSuccessResponse(res))
	}
}

func TestBadRequest(t *testing.T) {
	res, err := myService.badRequest()
	if err != nil {
		fmt.Println(sendErrorResponse(err))
	} else {
		fmt.Println(sendSuccessResponse(res))
	}
}

func TestNotFound(t *testing.T) {
	res, err := myService.notFound()
	if err != nil {
		fmt.Println(sendErrorResponse(err))
	} else {
		fmt.Println(sendSuccessResponse(res))
	}
}

func TestServiceUnavailable(t *testing.T) {
	res, err := myService.serviceUnavailable()
	if err != nil {
		fmt.Println(sendErrorResponse(err))
	} else {
		fmt.Println(sendSuccessResponse(res))
	}
}
