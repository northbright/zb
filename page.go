package main

import (
	//"fmt"
	//"log"
	"net/http"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

func getZB(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title": "转班申请",
	})
}

type Record struct {
	Data string
}

func admin(c *gin.Context) {
	var conn redis.Conn
	var err error
	msg := ""
	fields := []string{}
	records := []Record{}

	if conn, err = GetRedisConn(redisAddr, redisPassword); err != nil {
		msg = "连接数据库失败."
		goto end
	}
	defer conn.Close()

	if fields, err = redis.Strings(conn.Do("ZRANGE", "idx:time", 0, -1)); err != nil {
		msg = "获取时间索引失败."
		goto end
	}

	for _, f := range fields {
		r := Record{}
		if r.Data, err = redis.String(conn.Do("HGET", "records", f)); err != nil {
			msg = "获取记录失败."
			goto end
		}
		records = append(records, r)
	}
end:
	c.HTML(http.StatusOK, "admin.tmpl", gin.H{
		"title":   "转班申请",
		"msg":     msg,
		"count":   len(records),
		"Records": records,
	})
}

func success(c *gin.Context) {
	c.HTML(http.StatusOK, "success.tmpl", gin.H{
		"title": "提交成功",
		"msg":   "提交成功。请等待学校电话通知处理结果。",
	})
}
