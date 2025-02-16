package transport

import (
	"net/http"
	"strconv"

	"diianpro/coin-merch-store/internal/service"
	"diianpro/coin-merch-store/internal/transport/utils"

	"github.com/labstack/echo/v4"
)

type coinRoutes struct {
	coinService service.Coin
}

func newCoinRoutes(g *echo.Group, coinService service.Coin) {
	r := &coinRoutes{
		coinService: coinService,
	}

	g.POST("/sendCoin", r.send)
}

type Input struct {
	ToUser string `json:"toUser" validate:"required"`
	Amount int    `json:"amount" validate:"required"`
}

func (cr *coinRoutes) send(c echo.Context) error {
	var input Input
	if err := c.Bind(&input); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	toUser, err := strconv.Atoi(input.ToUser)
	if err != nil {
		return err
	}

	if err := c.Validate(input); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	fromUser := c.Get(userIdCtx).(int)

	if err := cr.coinService.TransferCoins(c.Request().Context(), fromUser, toUser, input.Amount); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, "Успешный ответ.")
}
