package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
	//"fmt"
)

func main() {
	router := gin.Default()
	router.Static("/static", "./static")

	router.LoadHTMLGlob("templates/*")
	//router.LoadHTMLFiles("templates/information.tmpl")
	router.GET("/information", func(c *gin.Context) {
		c.HTML(http.StatusOK, "information.tmpl", gin.H{
			"title": "Vitess benchmark",
		})
	})
	router.Run(":8080")

}
