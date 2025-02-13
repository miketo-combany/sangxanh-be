package api

import (
	"SangXanh/pkg/common/errors"
	"SangXanh/pkg/common/query"
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"net/http"
)

var validate = validator.New()

type Response interface {
}

type responseMeta struct {
	*query.Pagination
	Message string `json:"message"`
	Debug   any    `json:"debug,omitempty"`
	Error   string `json:"error,omitempty"`
}

type response struct {
	Meta responseMeta `json:"meta"`
	Data any          `json:"data,omitempty"`
}

func Success(data any) Response {
	return response{
		Meta: responseMeta{
			Message: "success",
		},
		Data: data,
	}
}

func SuccessPagination(data any, p *query.Pagination) Response {
	return response{
		Meta: responseMeta{
			Pagination: p,
			Message:    "success",
		},
		Data: data,
	}
}

type API[Req any] func(e echo.Context, req Req) (Response, error)

func Execute[Req any](c echo.Context, f func(e context.Context, req Req) (Response, error)) error {
	var req Req
	if err := c.Bind(&req); err != nil {
		return Serve(c, nil, err)
	}
	if err := validate.Struct(req); err != nil {
		return Serve(c, nil, errors.BadRequest(err.Error()))
	}
	resp, err := f(c.Request().Context(), req)
	return Serve(c, resp, err)
}

func Serve(c echo.Context, resp Response, err error) error {
	if err != nil {
		var badRequest *errors.BadRequestError
		ok := errors.As(err, &badRequest)
		if ok {
			return c.JSON(http.StatusBadRequest, response{
				Meta: responseMeta{
					Message: "",
					Debug:   badRequest.Debug,
					Error:   badRequest.Wrapper.Error(),
				},
			})
		}
		return c.JSON(http.StatusInternalServerError, response{
			Meta: responseMeta{
				Message: "",
				Error:   err.Error(),
			},
		})
	}
	return c.JSON(http.StatusOK, resp)
}
