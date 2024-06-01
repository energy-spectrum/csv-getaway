package route

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"csv-analyzer-api/internal/config"
	"csv-analyzer-api/internal/service/user"
	"csv-analyzer-api/internal/service/template"

	"csv-analyzer-api/internal/transport/http/middleware"
)

func Setup(cfg *config.Configuration,
	userService user.UserService,
	templateService template.TemplateService,
	router *gin.RouterGroup) {
	// CORS
	router.Use(middleware.CORSMiddleware())

	router.GET("/ping", func(с *gin.Context) {
		с.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	authRouter := router.Group("auth")
	NewAuthRouter(cfg, userService, authRouter)

	templateRouter := router.Group("template")
	NewTemplateRouter(cfg, templateService, templateRouter)
}
