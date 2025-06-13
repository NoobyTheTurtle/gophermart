package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/NoobyTheTurtle/gophermart/internal/controller/http/v1/request"
	"github.com/NoobyTheTurtle/gophermart/internal/controller/http/v1/response"
	"github.com/NoobyTheTurtle/gophermart/internal/entity"
	"github.com/NoobyTheTurtle/gophermart/internal/usecase/auth"
)

type authRoutes struct {
	authUseCase auth.AuthUseCase
}

func newAuthRoutes(g *gin.RouterGroup, authUseCase auth.AuthUseCase) {
	r := &authRoutes{
		authUseCase: authUseCase,
	}

	g.POST("/register", r.register)
	g.POST("/login", r.login)
}

// register handles user registration
// @Summary Register user
// @Description Register new user with login and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.Auth true "User credentials"
// @Success 200 {object} response.Auth
// @Failure 400 {object} response.Error
// @Failure 409 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /api/user/register [post]
func (r *authRoutes) register(c *gin.Context) {
	var req request.Auth
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request format")
		return
	}

	authData := &entity.Auth{
		Login:    req.Login,
		Password: req.Password,
	}

	token, err := r.authUseCase.Register(c.Request.Context(), authData)
	if err != nil {
		if errors.Is(err, auth.ErrUserAlreadyExists) {
			errorResponse(c, http.StatusConflict, "login already taken")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.SetCookie("token", token, 3600*24, "/", "", false, true)
	c.JSON(http.StatusOK, response.Auth{Token: token})
}

// login handles user authentication
// @Summary Login user
// @Description Authenticate user with login and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body request.Auth true "User credentials"
// @Success 200 {object} response.Auth
// @Failure 400 {object} response.Error
// @Failure 401 {object} response.Error
// @Failure 500 {object} response.Error
// @Router /api/user/login [post]
func (r *authRoutes) login(c *gin.Context) {
	var req request.Auth
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid request format")
		return
	}

	authData := &entity.Auth{
		Login:    req.Login,
		Password: req.Password,
	}

	token, err := r.authUseCase.Login(c.Request.Context(), authData)
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			errorResponse(c, http.StatusUnauthorized, "invalid login/password pair")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.SetCookie("token", token, 3600*24, "/", "", false, true)
	c.JSON(http.StatusOK, response.Auth{Token: token})
}
