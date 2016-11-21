package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

var (
	redisAddr     = ":6379"
	redisPassword = ""
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

	r.GET("/campuses/", getCampuses)
	r.GET("/grades/:campus/", getGrades)
	r.GET("/periods/:campus/:grade/", getPeriods)

	r.Run(":8000")
end:
	if err != nil {
		log.Printf("main() error: %v\n", err)
		return
	}
}
