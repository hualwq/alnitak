package api

import (
	"github.com/gin-gonic/gin"
	"interastral-peace.com/alnitak/internal/domain/dto"
	"interastral-peace.com/alnitak/internal/resp"
	"interastral-peace.com/alnitak/internal/service"
)

// AddSensitiveWord 添加敏感词
func AddSensitiveWord(ctx *gin.Context) {
	var req dto.AddSensitiveWordReq
	if err := ctx.Bind(&req); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.AddSensitiveWord(req); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.Ok(ctx)
}

// DeleteSensitiveWord 删除敏感词
func DeleteSensitiveWord(ctx *gin.Context) {
	var req dto.DeleteSensitiveWordReq
	if err := ctx.Bind(&req); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	if err := service.DeleteSensitiveWord(req); err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.Ok(ctx)
}

// GetSensitiveWordList 获取敏感词列表
func GetSensitiveWordList(ctx *gin.Context) {
	var req dto.SensitiveWordListReq
	if err := ctx.ShouldBindQuery(&req); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 10
	}

	result, err := service.GetSensitiveWordList(req)
	if err != nil {
		resp.FailWithMessage(ctx, err.Error())
		return
	}

	resp.OkWithData(ctx, result)
}

// CheckSensitiveWord 检测敏感词
func CheckSensitiveWord(ctx *gin.Context) {
	var req dto.CheckSensitiveWordReq
	if err := ctx.Bind(&req); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	result := service.CheckSensitiveWord(req)
	resp.OkWithData(ctx, result)
}

// ReplaceSensitiveWord 替换敏感词
func ReplaceSensitiveWord(ctx *gin.Context) {
	var req dto.ReplaceSensitiveWordReq
	if err := ctx.Bind(&req); err != nil {
		resp.FailWithMessage(ctx, "请求参数有误")
		return
	}

	result := service.ReplaceSensitiveWord(req)
	resp.OkWithData(ctx, result)
}
