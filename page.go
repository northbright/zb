package main

import (
	//"fmt"
	//"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func getZB(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "转班申请",
	})
}
