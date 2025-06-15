package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/NoobyTheTurtle/gophermart/internal/controller/http/v1/request"
	"github.com/NoobyTheTurtle/gophermart/internal/controller/http/v1/response"
	"github.com/NoobyTheTurtle/gophermart/internal/entity"
)

type authRoutes struct {
	authUseCase AuthUseCase
	logger      Logger
}

func newAuthRoutes(g *gin.RouterGroup, authUseCase AuthUseCase, logger Logger) {
	r := &authRoutes{
		authUseCase: authUseCase,
		logger:      logger,
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
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request format")
		return
	}

	authData := &entity.Auth{
		Login:    req.Login,
		Password: req.Password,
	}

	token, err := r.authUseCase.Register(c.Request.Context(), authData)
	if err != nil {
		if errors.Is(err, entity.ErrUserAlreadyExists) {
			response.ErrorResponse(c, http.StatusConflict, "login already taken")
			return
		}
		r.logger.Error("Failed to register user", "login", req.Login, "error", err)
		response.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
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
		response.ErrorResponse(c, http.StatusBadRequest, "invalid request format")
		return
	}

	authData := &entity.Auth{
		Login:    req.Login,
		Password: req.Password,
	}

	token, err := r.authUseCase.Login(c.Request.Context(), authData)
	if err != nil {
		if errors.Is(err, entity.ErrInvalidCredentials) {
			response.ErrorResponse(c, http.StatusUnauthorized, "invalid login/password pair")
			return
		}
		r.logger.Error("Failed to login user", "login", req.Login, "error", err)
		response.ErrorResponse(c, http.StatusInternalServerError, "internal server error")
		return
	}

	c.SetCookie("token", token, 3600*24, "/", "", false, true)
	c.JSON(http.StatusOK, response.Auth{Token: token})
}
