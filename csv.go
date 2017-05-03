package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/garyburd/redigo/redis"
	"github.com/northbright/pathhelper"
)

var (
	csvFiles = []string{
		"csv/zb-weining.csv",
		"csv/zb-zhongshan.csv",
	}
	gradeScores = map[string]int{
		"幼小":      1,
		"幼中":      2,
		"幼大":      3,
		"一年级":     4,
		"二年级":     5,
		"三年级":     6,
		"四年级":     7,
		"五年级":     8,
		"六年级":     9,
		"初一":      10,
		"初二":      11,
		"初三":      12,
		"高中":      13,
		"成人":      14,
		"国画":      15,
		"书法":      16,
		"幼儿书法彩墨班": 17,
		"花鸟画":     18,
	}
	weekDayScores = map[string]int{
		"周一": 1,
		"周二": 2,
		"周三": 3,
		"周四": 4,
		"周五": 5,
		"周六": 6,
		"周日": 7,
	}
)

func loadRecordsFromCSV() (records [][]string, err error) {
	totalRecords := [][]string{}

	for _, csvFile := range csvFiles {
		p := ""
		var f *os.File
		records := [][]string{}

		if p, err = pathhelper.GetAbsPath(csvFile); err != nil {
			goto end
		}
		if f, err = os.Open(p); err != nil {
			goto end
		}
		r := csv.NewReader(f)

		records, err = r.ReadAll()
		if err != nil {
			goto end
		}

		totalRecords = append(totalRecords, records[1:]...)
	}

	log.Printf("totalRecords: %v\n", totalRecords)
end:
	if err != nil {
		log.Printf("loadRecordsFromCSV() error: %v\n", err)
		return [][]string{}, err
	}
	return totalRecords, nil
}

func parsePeriod(period string) (beginHour, beginMin, endHour, endMin int, err error) {
	hours := [2]int64{}
	mins := [2]int64{}
	arr := strings.Split(period, "-")
	if len(arr) != 2 {
		err = fmt.Errorf("Split '-' failed.")
		goto end
	}

	for i, time := range arr {
		subArr := strings.Split(time, ":")
		if len(subArr) != 2 {
			err = fmt.Errorf("Split ':' failed.")
			goto end
		}

		if hours[i], err = strconv.ParseInt(subArr[0], 10, 32); err != nil {
			goto end
		}

		if mins[i], err = strconv.ParseInt(subArr[1], 10, 32); err != nil {
			goto end
		}
	}
end:
	if err != nil {
		log.Printf("parsePeriod() error: %v\n", err)
		return 0, 0, 0, 0, err
	}
	return int(hours[0]), int(mins[0]), int(hours[1]), int(mins[1]), nil
}

func initPeriods(records [][]string) (err error) {
	var c redis.Conn

	if c, err = GetRedisConn(redisAddr, redisPassword); err != nil {
		goto end
	}
	defer c.Close()

	// record format:
	// 0: Campus, 1: grade, 2: week day, 3: period.
	for _, r := range records {
		score := 0
		weekDayScore := 0
		k := ""
		ok := false

		if len(records) < 4 {
			err = fmt.Errorf("Length of record < 4.")
			goto end
		}
		campus, grade, weekDay, period := r[0], r[1], r[2], r[3]

		// Get grade score
		if score, ok = gradeScores[grade]; !ok {
			score = 0
		}

		// Set key: "grade", members: all grades.
		if _, err = c.Do("ZADD", "grades", score, grade); err != nil {
			goto end
		}

		// Set key: grade:"campuses", memebers: available campuses.
		k = fmt.Sprintf("%v:campuses", grade)
		if _, err = c.Do("SADD", k, campus); err != nil {
			goto end
		}

		campus_grade := fmt.Sprintf("%v/%v", campus, grade)
		weekDay_period := fmt.Sprintf("%v/%v", weekDay, period)

		// Get week day - period score
		var beginHour, beginMin = 0, 0
		if beginHour, beginMin, _, _, err = parsePeriod(period); err != nil {
			goto end
		}

		if weekDayScore, ok = weekDayScores[weekDay]; !ok {
			weekDayScore = 0
		}

		score = weekDayScore*1000 + beginHour*10 + beginMin

		// Set key: campus/grade, members: weekday/period.
		if _, err = c.Do("ZADD", campus_grade, score, weekDay_period); err != nil {
			goto end
		}
	}

end:
	if err != nil {
		log.Printf("initMaps() error: %v\n", err)
		return err
	}
	return nil
}
