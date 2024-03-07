package main

import (
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"time"
)

const (
	minPort = 49152
	maxPort = 65535
)

func RequestSSHPort(ctx *gin.Context) {
	//todo 用户校验
	src := rand.NewSource(time.Now().UnixNano())
	r := rand.New(src)
	randomPort := r.Intn(maxPort-minPort+1) + minPort
	ctx.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{"port": randomPort},
		"msg":  nil,
	})
}
