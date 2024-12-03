package handler

import (
	"net/http"
	"workflow/internal/common"
	"workflow/internal/constants"
	"workflow/internal/dto"
	"workflow/internal/form/service"
	"workflow/internal/utils"

	"github.com/gin-gonic/gin"
)

func CreateFormDefine(c *gin.Context) {
	form := &dto.CreateFormDefineRequest{}
	rsp := &common.Result{}
	if err := c.ShouldBind(form); err != nil {
		c.JSON(http.StatusOK, rsp.Error(common.ParamError, common.GetErrorMsg(form, err)))
		return
	}
	uid, err := utils.StringToInt(c.GetString(constants.ACCESSTOKENUSERIDKEY))
	if err != nil {
		c.JSON(http.StatusOK, rsp.Error(common.ParamError, "get uid failed"))
		return
	}
	define, err := service.NewFormDefine(uid, form.Name, form.Code, form.FormStructure, form.ComponentStructure)
	if err != nil {
		c.JSON(http.StatusOK, rsp.Error(common.ServiceError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, rsp.Success(&dto.CreateFormDefineResponse{
		FormDefineId: define.Meta.ID,
	}))
}
