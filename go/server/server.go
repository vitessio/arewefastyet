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

	//Home page
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Vitess benchmark",
		})
	})

	//Request benchmark page
	router.GET("/search_compare", func(c *gin.Context) {
		c.HTML(http.StatusOK, "search_compare.tmpl", gin.H{
			"title": "Vitess benchmark",
		})
	})

	//Request benchmark page
	router.GET("/request_benchmark", func(c *gin.Context) {
		c.HTML(http.StatusOK, "request_benchmark.tmpl", gin.H{
			"title": "Vitess benchmark",
		})
	})




	router.Run(":8080")

}
