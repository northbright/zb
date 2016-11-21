package main

import (
	"log"
)

var (
	redisAddr     = ":6379"
	redisPassword = ""
)

func main() {
	var err error

	records := [][]string{}
	if records, err = loadRecordsFromCSV(); err != nil {
		goto end
	}
	log.Printf("len of records: %v\n", len(records))

	if err = initPeriods(records); err != nil {
		goto end
	}

end:
	if err != nil {
		log.Printf("main() error: %v\n", err)
		return
	}
}
