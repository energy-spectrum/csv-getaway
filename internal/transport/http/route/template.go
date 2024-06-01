package route

import (
	"csv-analyzer-api/internal/config"
	"csv-analyzer-api/internal/transport/http/controller"
	"csv-analyzer-api/internal/service/template"

	"github.com/gin-gonic/gin"
)

func NewTemplateRouter(cfg *config.Configuration, templateService template.TemplateService, group *gin.RouterGroup) {
	tc := controller.TemplateController{
		TemplateService: templateService,
		Cfg:             cfg,
	}

	group.POST("/", tc.Create)
	group.GET("/:id", tc.GetByID)
	group.PUT("/", tc.Update)
	group.DELETE("/:id", tc.Delete)
}
