package transport

import (
	"net/http"

	"diianpro/coin-merch-store/internal/service"
	serviceErr "diianpro/coin-merch-store/internal/service/utils"
	"diianpro/coin-merch-store/internal/transport/utils"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
)

type authRoutes struct {
	authService service.Auth
}

func newAuthRoutes(g *echo.Group, authService service.Auth) {
	r := &authRoutes{
		authService: authService,
	}

	g.POST("/auth", r.auth)
}

type AuthInput struct {
	Username string `json:"username" validate:"required,min=4,max=32"`
	Password string `json:"password" validate:"required,password"`
}

// @Summary Sign up
// @Description Sign up
// @Tags auth
// @Accept json
// @Produce json
// @Param input body signUpInput true "input"
// @Success 201 {object} v1.authRoutes.signUp.response
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /auth/sign-up [post]
func (r *authRoutes) auth(c echo.Context) error {
	var input AuthInput

	if err := c.Bind(&input); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	id, err := r.authService.CreateUser(c.Request().Context(), service.AuthCreateUserInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		if errors.Is(err, serviceErr.ErrAccountAlreadyExists) {
			utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
			return err
		}
		utils.NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Id int64 `json:"id"`
	}

	return c.JSON(http.StatusCreated, response{
		Id: int64(id),
	})
}

type signInInput struct {
	Username string `json:"username" validate:"required,min=4,max=32"`
	Password string `json:"password" validate:"required,password"`
}

// @Summary Sign in
// @Description Sign in
// @Tags auth
// @Accept json
// @Produce json
// @Param input body signInInput true "input"
// @Success 200 {object} v1.authRoutes.signIn.response
// @Failure 400 {object} echo.HTTPError
// @Failure 500 {object} echo.HTTPError
// @Router /auth/sign-in [post]
func (r *authRoutes) signIn(c echo.Context) error {
	var input signInInput

	if err := c.Bind(&input); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, "invalid request body")
		return err
	}

	if err := c.Validate(input); err != nil {
		utils.NewErrorResponse(c, http.StatusBadRequest, err.Error())
		return err
	}

	token, err := r.authService.GenerateToken(c.Request().Context(), service.AuthGenerateTokenInput{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		if errors.Is(err, serviceErr.ErrUserNotFound) {
			utils.NewErrorResponse(c, http.StatusBadRequest, "invalid username or password")
			return err
		}
		utils.NewErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return err
	}

	type response struct {
		Token string `json:"token"`
	}

	return c.JSON(http.StatusOK, response{
		Token: token,
	})
}
