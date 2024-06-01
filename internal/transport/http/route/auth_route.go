package route

import (
	"csv-analyzer-api/internal/config"
	"csv-analyzer-api/internal/service/user"
	"csv-analyzer-api/internal/transport/http/controller"

	"github.com/gin-gonic/gin"
)

func NewAuthRouter(cfg *config.Configuration, userService user.UserService, group *gin.RouterGroup) {
	ac := controller.AuthController{
		UserService: userService,
		Cfg:         cfg,
	}

	group.POST("/registration", ac.Register)
	group.POST("/login", ac.Login)
}
