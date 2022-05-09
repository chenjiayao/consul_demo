package controllers

import (
	"context"
	"encoding/json"
	"fake_service/globals"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/imroc/req/v3"
	"github.com/spf13/viper"
)

type IndexController struct {
}

func (c *IndexController) TestV1(ctx *gin.Context) {

	url := fmt.Sprintf("%s/v1", viper.GetString("upstream_service.host"))

	resp, err := req.R().Get(url)
	if err != nil {
		ctx.JSON(500, gin.H{
			"detail": err.Error(),
		})
		return
	}

	var body map[string]string
	json.Unmarshal(resp.Bytes(), &body)
	ctx.JSON(200, gin.H{
		"url": ctx.Request.URL.String(),
		"upstream_service": gin.H{
			"body":        body,
			"status_code": resp.StatusCode,
			"header":      resp.Header,
		},
	})
}

func (c *IndexController) TestV2(ctx *gin.Context) {
	url := fmt.Sprintf("%s/v2", viper.GetString("upstream_service.host"))

	resp, err := req.R().Get(url)
	if err != nil {
		ctx.JSON(500, gin.H{
			"detail": err.Error(),
		})
		return
	}

	var body map[string]string
	json.Unmarshal(resp.Bytes(), &body)
	ctx.JSON(200, gin.H{
		"url": ctx.Request.URL.String(),
		"upstream_service": gin.H{
			"body":        body,
			"status_code": resp.StatusCode,
			"header":      resp.Header,
		},
	})
}

func (c IndexController) Health(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "ok",
	})
}

func (c *IndexController) Redis(ctx *gin.Context) {

	res, err := globals.RedisClient.Ping(context.TODO()).Result()
	if err != nil {
		ctx.JSON(500, gin.H{
			"message": err.Error(),
		})
		return
	}

	ctx.JSON(200, gin.H{
		"message": res,
	})
}
