package main

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

func getCampuses(c *gin.Context) {
	var err error
	var conn redis.Conn
	err_msg := ""
	exist := false
	campuses := []string{}
	grade := ""
	k := ""

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

	k = fmt.Sprintf("%v:campuses", grade)
	if exist, err = redis.Bool(conn.Do("EXISTS", k)); err != nil {
		err_msg = "Internal server error."
		goto end
	}

	if !exist {
		err_msg = fmt.Sprintf("Grade: %v does not exist in db.", grade)
		goto end
	}

	if campuses, err = redis.Strings(conn.Do("SMEMBERS", k)); err != nil {
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
	err_msg := ""
	exist := false
	grades := []string{}

	if conn, err = GetRedisConn(redisAddr, redisPassword); err != nil {
		err_msg = "Failed connect to db."
		goto end
	}
	defer conn.Close()

	if exist, err = redis.Bool(conn.Do("EXISTS", "grades")); err != nil {
		err_msg = "Internal server error."
		goto end
	}

	if !exist {
		err_msg = "Grades does not exist in db."
		goto end
	}

	if grades, err = redis.Strings(conn.Do("ZRANGE", "grades", 0, -1)); err != nil {
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
