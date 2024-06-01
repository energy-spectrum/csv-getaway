package controller

import (
	"net/http"
	"strconv"

	"csv-analyzer-api/internal/config"
	"csv-analyzer-api/internal/entity"
	"csv-analyzer-api/internal/service/template"
	"csv-analyzer-api/internal/value"

	"github.com/gin-gonic/gin"
	_ "github.com/santosh/gingo/docs"
	log "github.com/sirupsen/logrus"
)

type TemplateController struct {
	TemplateService template.TemplateService
	Cfg             *config.Configuration
}

// @Summary CreateTemplate
// @Description
// @Tags Template
// @Accept json
// @Produce json
// @Param email formData string true "Email address of the user"
// @Success 200 {object} entity.Template
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /template [post]
func (tc *TemplateController) Create(c *gin.Context) {
	var req entity.Template
	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf("%d err: %v", http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	template, err := tc.TemplateService.Create(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, template)
}

// @Summary GetByID
// @Description get Template
// @Tags Template
// @Accept json
// @Produce json
// @Success 200 {object} TemplateResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /template/{id} [get]
func (tc *TemplateController) GetByID(c *gin.Context) {
	id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse("Template not found"))
		return
	}

	template, err := tc.TemplateService.GetByID(c, value.TemplateID(id))
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse("Wrong login or password"))
		return
	}

	c.JSON(http.StatusOK, template)
}

// @Summary UpdateTemplate
// @Description
// @Tags Template
// @Accept json
// @Produce json
// @Success 200 {object} entity.Template
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /template [put]
func (tc *TemplateController) Update(c *gin.Context) {
	var req entity.Template
	err := c.ShouldBind(&req)
	if err != nil {
		log.Errorf("%d err: %v", http.StatusBadRequest, err.Error())
		c.JSON(http.StatusBadRequest, errorResponse(err.Error()))
		return
	}

	err = tc.TemplateService.Update(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, struct{}{})
}

// @Summary DeleteTemplate
// @Description
// @Tags Template
// @Accept json
// @Produce json
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /template/{id} [delete]
func (tc *TemplateController) Delete(c *gin.Context) {
	id, err := strconv.ParseInt(c.Params.ByName("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusNotFound, errorResponse("Template not found"))
		return
	}

	arg := template.DeleteArg{
		ID: value.TemplateID(id),
	}
	err = tc.TemplateService.Delete(c, &arg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errorResponse(err.Error()))
		return
	}

	c.JSON(http.StatusOK, successResponse("Template is deleted"))
}
