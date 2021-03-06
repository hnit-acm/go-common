package main

import (
	"fmt"
	"github.com/hnit-acm/hfunc/hapi"

	"github.com/gin-gonic/gin"
)

func Hello(c *gin.Context) {
	fmt.Println("runing hello world")
	c.JSON(200, "hello world")
	return
}

func Router(g *gin.Engine, mw ...gin.HandlerFunc) *gin.Engine {
	g.Use(gin.Recovery())
	g.Use(mw...)
	v1 := g.Group("/v1")
	{
		helloGroup := v1.Group("hello")
		{
			helloGroup.GET("world", Hello)
		}
	}
	return g
}

func Middleware(c *gin.Context) {
	fmt.Println("before")
	c.Next()
	fmt.Println("after")
}

func main() {
	hapi.Serve("8080", nil, func(c *gin.Engine) {
		Router(c, Middleware)
	})
}
