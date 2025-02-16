package transport

import (
	"diianpro/coin-merch-store/internal/service"

	"github.com/labstack/echo/v4"
)

func NewRouter(handler *echo.Echo, services *service.Services) {
	handler.GET("/health", func(c echo.Context) error { return c.NoContent(200) })

	auth := handler.Group("api/auth")
	{
		newAuthRoutes(auth, services.Auth)
	}

	authMiddleware := &AuthMiddleware{services.Auth}
	v1 := handler.Group("/api/", authMiddleware.UserIdentity)
	{
		newAuthRoutes(v1.Group("/info"), services.Auth)
		newCoinRoutes(v1.Group("/sendCoin"), services.Coin)
		newMerchRoutes(v1.Group("/buy/{item}"), services.Merch)
	}
}
