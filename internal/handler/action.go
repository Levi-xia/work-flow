package handler

import (
	"net/http"
	"workflow/internal/action/service"
	"workflow/internal/common"
	"workflow/internal/constants"
	"workflow/internal/dto"
	"workflow/internal/utils"

	"github.com/gin-gonic/gin"
)


func CreateActionDefine(c *gin.Context) {
	form := &dto.CreateActionDefineRequest{}
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
	define, err := service.NewActionDefine(uid, form.Name, form.Code, form.Protocol, form.Content, form.InputStructs, form.OutputChecks)
	if err != nil {
		c.JSON(http.StatusOK, rsp.Error(common.ServiceError, err.Error()))
		return
	}
	c.JSON(http.StatusOK, rsp.Success(&dto.CreateActionDefineResponse{
		ActionDefineID: define.Meta.ID,
	}))
}