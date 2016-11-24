package main

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/garyburd/redigo/redis"
)

func validateName(name string) bool {
	if len(name) < 6 || len(name) > 60 {
		return false
	}

	if strings.Contains(name, ",") {
		return false
	}

	return true
}

func validateTel(tel string) bool {
	p := `^\d{11}$`
	r := regexp.MustCompile(p)

	if r.MatchString(tel) {
		return true
	}
	return false
}

func validatePeriod(c redis.Conn, grade, campus, period string) bool {
	valid := false
	var err error
	exists := false
	k := ""
	score := ""

	k = fmt.Sprintf("%v/%v", campus, grade)
	if exists, err = redis.Bool(c.Do("EXISTS", k)); err != nil {
		goto end
	}

	if !exists {
		goto end
	}

	if score, err = redis.String(c.Do("ZSCORE", k, period)); err != nil {
		goto end
	}

	if score == "" {
		goto end
	}

	valid = true
end:
	return valid
}
