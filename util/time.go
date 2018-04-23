package util

import (
	"time"
)

var (
	format      = "2006-01-02 15:04:05"
	formatOfDay = "2006-01-02"
)

//判断某一时间是否在2个时间之间
func IsBetweenTime(since time.Time, till time.Time, t time.Time) (res bool) {
	if t.After(since) && t.Before(till) {
		res = true
	}
	return
}

//把当前时间utc转换为某一时区的时间
func CovnNOWUTC2Location(location string) (dt time.Time, err error) {
	now := time.Now().UTC().Format(format)
	loc, err := time.LoadLocation(location)
	if err != nil {
		return
	}
	dt, err = time.ParseInLocation(format, now, loc)
	return
}

//把任一时间转换为某一时区的时间(通常为UTC) 传入的时间格式为『2006-01-02 15:04:05』
func CovnTime2Location(t string, location string) (dt time.Time, err error) {
	loc, err := time.LoadLocation(location)
	if err != nil {
		return
	}
	dt, err = time.ParseInLocation(format, t, loc)
	return
}

//任一时区的时间转为其他时区的时间。
func Zone2Zone(t time.Time, zone string) (dt time.Time, err error) {
	utc := t.UTC()
	timeStr := utc.Format(format)
	return CovnTime2Location(timeStr, zone)
}

//某一时区一天的开始和结束
func BgeinAndEndDAYOfZone(zone string) (begin time.Time, end time.Time, err error) {
	now, err := CovnNOWUTC2Location(zone)
	if err != nil {
		return
	}
	nowStr := now.Format(formatOfDay) + " 00:00:00"
	begin, err = CovnTime2Location(nowStr, zone)
	if err != nil {
		return
	}

	oneDay, err := time.ParseDuration("24h") //一天后
	if err != nil {
		return
	}
	end = begin.Add(oneDay)
	return
}
