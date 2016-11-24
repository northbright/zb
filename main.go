package main

import (
	"fmt"
	"log"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/northbright/pathhelper"
)

var (
	redisAddr     = ":6379"
	redisPassword = ""
	serverRoot    = ""
	templatesPath = ""
	staticPath    = ""
)

func main() {
	var err error
	r := gin.Default()

	records := [][]string{}
	if records, err = loadRecordsFromCSV(); err != nil {
		goto end

	}
	log.Printf("len of records: %v\n", len(records))

	if err = initPeriods(records); err != nil {
		goto end
	}

	serverRoot, _ = pathhelper.GetCurrentExecDir()
	templatesPath = path.Join(serverRoot, "templates")
	staticPath = path.Join(serverRoot, "static")

	// Serve Static files.
	r.Static("/static/", staticPath)

	// Load Templates.
	r.LoadHTMLGlob(fmt.Sprintf("%v/*", templatesPath))

	// Pages
	r.GET("/", getZB)
	r.POST("/zb", postZB)

	// APIs
	r.GET("/grades/", getGrades)
	r.GET("/campuses/:grade", getCampuses)
	r.GET("/periods/:campus/:grade/", getPeriods)

	r.Run(":80")
end:
	if err != nil {
		log.Printf("main() error: %v\n", err)
		return
	}
}
