package main

import (
	"fmt"
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

var (
	redisAddr     = ":6379"
	redisPassword = ""
)

func getCampuses(c *gin.Context) {
	var err error
	var conn redis.Conn
	err_msg := ""
	exist := false
	campuses := []string{}

	if conn, err = GetRedisConn(redisAddr, redisPassword); err != nil {
		err_msg = "Failed connect to db."
		goto end
	}
	defer conn.Close()

	if exist, err = redis.Bool(conn.Do("EXISTS", "campus")); err != nil {
		err_msg = "Internal server error."
		goto end
	}

	if !exist {
		err_msg = "Campus does not exist in db."
		goto end
	}

	if campuses, err = redis.Strings(conn.Do("SMEMBERS", "campus")); err != nil {
		err_msg = "Internal server error."
		goto end
	}

	c.JSON(200, gin.H{
		"success":  true,
		"err_msg":  "",
		"campuses": campuses,
	})

	return
end:
	c.JSON(200, gin.H{
		"success": false,
		"err_msg": err_msg,
	})
}

func getGrades(c *gin.Context) {
	var err error
	var conn redis.Conn
	campus := ""
	err_msg := ""
	exist := false
	grades := []string{}

	campus = c.Param("campus")
	if campus == "" {
		err_msg = "Empty campus name."
		goto end
	}

	if conn, err = GetRedisConn(redisAddr, redisPassword); err != nil {
		err_msg = "Failed connect to db."
		goto end
	}
	defer conn.Close()

	if exist, err = redis.Bool(conn.Do("EXISTS", campus)); err != nil {
		err_msg = "Internal server error."
		goto end
	}

	if !exist {
		err_msg = fmt.Sprintf("Campus: %v does not exist in db.", campus)
		goto end
	}

	if grades, err = redis.Strings(conn.Do("ZRANGE", campus, 0, -1)); err != nil {
		err_msg = "Internal server error."
		goto end
	}

	c.JSON(200, gin.H{
		"success":  true,
		"err_msg":  "",
		"campuses": grades,
	})

	return
end:
	c.JSON(200, gin.H{
		"success": false,
		"err_msg": err_msg,
	})
}

func getPeriods(c *gin.Context) {
	var err error
	var conn redis.Conn
	campus := ""
	grade := ""
	k := ""
	err_msg := ""
	exist := false
	periods := []string{}

	campus = c.Param("campus")
	if campus == "" {
		err_msg = "Empty campus name."
		goto end
	}

	grade = c.Param("grade")
	if grade == "" {
		err_msg = "Empty grade name."
		goto end
	}

	if conn, err = GetRedisConn(redisAddr, redisPassword); err != nil {
		err_msg = "Failed connect to db."
		goto end
	}
	defer conn.Close()

	k = fmt.Sprintf("%v/%v", campus, grade)
	if exist, err = redis.Bool(conn.Do("EXISTS", k)); err != nil {
		err_msg = "Internal server error."
		goto end
	}

	if !exist {
		err_msg = fmt.Sprintf("Campus/Grade: %v does not exist in db.", k)
		goto end
	}

	if periods, err = redis.Strings(conn.Do("ZRANGE", k, 0, -1)); err != nil {
		err_msg = "Internal server error."
		goto end
	}

	c.JSON(200, gin.H{
		"success":  true,
		"err_msg":  "",
		"campuses": periods,
	})

	return
end:
	c.JSON(200, gin.H{
		"success": false,
		"err_msg": err_msg,
	})
}

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
