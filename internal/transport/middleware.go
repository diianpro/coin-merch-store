package transport

import (
	"log/slog"
	"net/http"
	"strings"

	"diianpro/coin-merch-store/internal/service"
	"diianpro/coin-merch-store/internal/transport/utils"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

const (
	userIdCtx = "userId"
)

type AuthMiddleware struct {
	authService service.Auth
}

func (h *AuthMiddleware) UserIdentity(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		token, ok := bearerToken(c.Request())
		if !ok {
			slog.Error("AuthMiddleware.UserIdentity: bearerToken: %v", utils.ErrInvalidAuthHeader)
			utils.NewErrorResponse(c, http.StatusUnauthorized, utils.ErrInvalidAuthHeader.Error())
			return nil
		}

		userId, err := h.authService.ParseToken(token)
		if err != nil {
			log.Errorf("AuthMiddleware.UserIdentity: h.authService.ParseToken: %v", err)
			utils.NewErrorResponse(c, http.StatusUnauthorized, utils.ErrCannotParseToken.Error())
			return err
		}

		c.Set(userIdCtx, userId)

		return next(c)
	}
}

func bearerToken(r *http.Request) (string, bool) {
	const prefix = "Bearer "

	header := r.Header.Get(echo.HeaderAuthorization)
	if header == "" {
		return "", false
	}

	if len(header) > len(prefix) && strings.EqualFold(header[:len(prefix)], prefix) {
		return header[len(prefix):], true
	}

	return "", false
}
