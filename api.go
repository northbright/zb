package main

import (
	"fmt"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
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
		"success": true,
		"err_msg": "",
		"grades":  grades,
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
		"success": true,
		"err_msg": "",
		"periods": periods,
	})

	return
end:
	c.JSON(200, gin.H{
		"success": false,
		"err_msg": err_msg,
	})
}

func postZB(c *gin.Context) {
	var err error
	var conn redis.Conn
	var n, periodNum int64
	var k0, k1 = "", ""
	var t time.Time
	tm := ""
	campuses := []string{}
	err_msg := ""
	k := ""
	field := ""
	valid := false
	record := ""
	name := c.DefaultPostForm("name", "")
	tel := c.DefaultPostForm("tel", "")
	grade := c.DefaultPostForm("grade", "")
	currentCampus := c.DefaultPostForm("currentCampus", "")
	currentPeriod := c.DefaultPostForm("currentPeriod", "")
	wantedCampus := c.DefaultPostForm("wantedCampus", "")
	wantedPeriod := c.DefaultPostForm("wantedPeriod", "")

	log.Printf("name: %v\n", name)
	log.Printf("tel: %v\n", tel)
	log.Printf("grade: %v\n", grade)
	log.Printf("currentCampus: %v\n", currentCampus)
	log.Printf("currentPeriod: %v\n", currentPeriod)
	log.Printf("wantedCampus: %v\n", wantedCampus)
	log.Printf("wantedPeriod: %v\n", wantedPeriod)

	if name == "" || tel == "" || grade == "" || currentCampus == "" || currentPeriod == "" || wantedCampus == "" || wantedPeriod == "" {
		err_msg = "信息不完整，请返回重新填写."
		goto end
	}

	if valid = validateName(name); !valid {
		err_msg = "学生姓名有误，请重新填写."
		goto end
	}

	if valid = validateTel(tel); !valid {
		err_msg = "联系手机号码有误，请重新填写."
		goto end
	}

	if conn, err = GetRedisConn(redisAddr, redisPassword); err != nil {
		err_msg = "连接数据库失败."
		goto end
	}
	defer conn.Close()

	if valid = validatePeriod(conn, grade, currentCampus, currentPeriod); !valid {
		err_msg = "当前时段有误：校区，年级与时段不匹配."
		goto end
	}

	if valid = validatePeriod(conn, grade, wantedCampus, wantedPeriod); !valid {
		err_msg = "期望时段有误：校区，年级与时段不匹配."
		goto end
	}

	k = fmt.Sprintf("%v:campuses", grade)
	if campuses, err = redis.Strings(conn.Do("SMEMBERS", k)); err != nil {
		err_msg = "获取年级对应校区信息失败."
		goto end
	}

	for _, campus := range campuses {
		k = fmt.Sprintf("%v/%v", campus, grade)
		if n, err = redis.Int64(conn.Do("ZCARD", k)); err != nil {
			err_msg = "获取年级校区对应时段数量失败."
			goto end
		}
		log.Printf("k: %v, n: %v\n", k, n)
		periodNum += n
	}

	k0 = fmt.Sprintf("%v/%v", currentCampus, currentPeriod)
	k1 = fmt.Sprintf("%v/%v", wantedCampus, wantedPeriod)

	if periodNum <= 1 {
		err_msg = "无可选时间段，不能转班."
		goto end
	} else {
		if k0 == k1 {
			err_msg = "当前时段与期望时段一致，请重新选择."
			goto end
		}
	}

	k = "records"
	// Field: name:tel
	field = fmt.Sprintf("%v:%v", name, tel)
	t = time.Now().Local()
	tm = fmt.Sprintf("%04d/%02d/%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
	record = fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v", name, tel, grade, currentCampus, currentPeriod, wantedCampus, wantedPeriod, tm)
	log.Printf("record: %v\n", record)

	conn.Send("MULTI")
	conn.Send("HSET", k, field, record)
	conn.Send("ZADD", "idx:time", t.Unix(), field)

	if _, err = conn.Do("EXEC"); err != nil {
		err_msg = "数据写入错误."
		goto end
	}

end:
	if err_msg != "" {
		c.JSON(200, gin.H{
			"success": false,
			"err_msg": err_msg,
		})
	} else {
		c.JSON(200, gin.H{
			"success": true,
			"err_msg": err_msg,
		})
	}
}
