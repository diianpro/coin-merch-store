package transport

import (
	"net/http"

	"diianpro/coin-merch-store/internal/service"
	"diianpro/coin-merch-store/internal/transport/utils"

	"github.com/labstack/echo/v4"
)

type merchRoutes struct {
	merchService service.Merch
}

func newMerchRoutes(g *echo.Group, merchService service.Merch) {
	r := &merchRoutes{
		merchService: merchService,
	}

	g.GET("/buy/{item}", r.buy)
}

type ItemInput struct {
	Item string `json:"item" validate:"required"`
}

func (m *merchRoutes) buy(c echo.Context) error {
	var input ItemInput

	if err := c.Bind(&input); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	userID := c.Get(userIdCtx).(int)

	if err := m.merchService.OrderMerch(c.Request().Context(), userID, input.Item); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, "Успешный ответ.")
}
