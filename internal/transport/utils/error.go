package utils

import (
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

var (
	ErrInvalidAuthHeader = fmt.Errorf("invalid auth header")
	ErrCannotParseToken  = fmt.Errorf("cannot parse token")
)

func NewErrorResponse(c echo.Context, errStatus int, message string) {
	err := errors.New(message)
	var HTTPError *echo.HTTPError
	ok := errors.As(err, &HTTPError)
	if !ok {
		report := echo.NewHTTPError(errStatus, err.Error())
		_ = c.JSON(errStatus, report)
	}
	c.Error(errors.New("internal server error"))
}
