package controllers

import (
	"github.com/gin-gonic/gin"
)

type IndexController struct {
}

func (c IndexController) TestV1(ctx *gin.Context) {

	ctx.JSON(200, gin.H{
		"message": "hello world from upstream service",
		"version": "v1",
	})
}

func (c IndexController) TestV2(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "hello world from upstream service",
		"version": "v2",
	})
}

func (c IndexController) Health(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "ok",
	})
}
