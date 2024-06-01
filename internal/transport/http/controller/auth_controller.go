package controller

import (
	"mime/multipart"
	"net/http"

	"csv-analyzer-api/internal/config"
	"csv-analyzer-api/internal/service/user"
	"csv-analyzer-api/internal/util"
	"csv-analyzer-api/internal/value"

	"github.com/gin-gonic/gin"
	_ "github.com/santosh/gingo/docs"
	log "github.com/sirupsen/logrus"
)

type AuthController struct {
	UserService user.UserService
	Cfg         *config.Configuration
}

type RegisterRequest struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required,min=4,max=25"`
	Name     string `form:"name" binding:"required,min=1,max=255"`
	Role     string `form:"role" binding:"required"`

	AvatarImage *multipart.FileHeader `form:"avatarImage"`
}

type AuthResponse struct {
	AccessToken string `json:"accessToken"`
}

// @Summary Register
// @Description Register a new user with email, password, fullname, nickname, and avatar image
// @Tags Auth
// @Accept multipart/form-data
// @Produce json
// @Param email formData string true "Email address of the user"
// @Param password formData string true "Password for the user account"
// @Param fullname formData string true "Full name of the user"
// @Param nickname formData string true "Nickname of the user"
// @Param avatarImage formData file true "Avatar image for the user"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/registration [post]
func (ac *AuthController) Register(c *gin.Context) {
	var req RegisterRequest
	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf("%d err: %v", http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	arg := user.CreateArg{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
		Role:     value.Role(req.Role),
	}

	user, err := ac.UserService.Create(c, &arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	accessToken, err := util.CreateAccessToken(user, ac.Cfg.HTTPServer.AccessTokenSecret,
		ac.Cfg.HTTPServer.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	authResponse := AuthResponse{
		AccessToken: accessToken,
	}

	c.JSON(http.StatusOK, authResponse)
}

type LoginRequest struct {
	Email    string `form:"email" binding:"required,email"`
	Password string `form:"password" binding:"required,min=4,max=25"`
}

// @Summary Login
// @Description Authenticate a user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param email formData string true "Email address of the user"
// @Param password formData string true "Password for the user account"
// @Success 200 {object} AuthResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /auth/login [post]
func (ac *AuthController) Login(c *gin.Context) {
	var req LoginRequest
	err := c.ShouldBind(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	user, err := ac.UserService.GetByEmail(c, req.Email)
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse("Wrong login or password"))
		return
	}

	if err := util.CheckPassword(req.Password, user.HashedPassword); err != nil {
		log.Info(err)
		c.JSON(http.StatusUnauthorized, errorResponse("Wrong login or password"))
		return
	}

	accessToken, err := util.CreateAccessToken(user, ac.Cfg.HTTPServer.AccessTokenSecret,
		ac.Cfg.HTTPServer.AccessTokenExpiryHour)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	authResponse := AuthResponse{
		AccessToken: accessToken,
	}

	c.JSON(http.StatusOK, authResponse)
}
