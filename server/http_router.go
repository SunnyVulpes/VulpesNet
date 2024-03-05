package server

import "github.com/gin-gonic/gin"

func InitHttpRouter() *gin.Engine {
	r := gin.Default()

	return r
}
