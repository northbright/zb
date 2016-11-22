package main

import (
	//"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getZB(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "长宁美校转班申请",
	})
}

func postZB(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "长宁美校转班申请",
	})
}
