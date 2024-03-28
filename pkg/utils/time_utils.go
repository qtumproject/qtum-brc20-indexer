package utils

import (
	"time"
)

const (
	SecondFormatLayout = "2006-01-02 15:04:05"
	DayFormatLayout    = "2006-01-02"
)

func IsSameDay(t1 int64, t2 int64, tz *time.Location) bool {
	if tz == nil {
		tz = time.FixedZone("CST", 8*3600)
		//tz = time.FixedZone("UTC", 0)
	}
	return time.Unix(t1, 0).In(tz).Format(DayFormatLayout) == time.Unix(t2, 0).In(tz).Format(DayFormatLayout)
}

func IsSameUTCDay(timestamp1, timestamp2 int64) bool {
	timeZone, _ := time.LoadLocation("UTC")
	return IsSameDay(timestamp1, timestamp2, timeZone)
	//utcTime1 := time.Unix(timestamp1, 0).UTC()
	//utcTime2 := time.Unix(timestamp2, 0).UTC()
	//return utcTime1.Year() == utcTime2.Year() &&
	//	utcTime1.Month() == utcTime2.Month() &&
	//	utcTime1.Day() == utcTime2.Day()
}

func TimestampUnixToFormat(unixTimestamp int64, tz *time.Location) string {
	if tz == nil {
		tz = time.FixedZone("CST", 8*3600)
		//tz = time.FixedZone("UTC", 0)
	}
	if unixTimestamp == 0 {
		return time.Now().In(tz).Format(SecondFormatLayout)
	}
	return time.Unix(unixTimestamp, 0).In(tz).Format(SecondFormatLayout)
}

func GetUTCDayFormat(timestamp int64) string {
	currentTime := time.Now().UTC()
	if timestamp > 0 {
		currentTime = time.Unix(timestamp, 0).UTC()
	}
	return currentTime.Format(DayFormatLayout)
}

func TimestampFormatToUnix(strTimestamp string) (int64, error) {
	timestamp, err := time.Parse(SecondFormatLayout, strTimestamp)
	if err != nil {
		return 0, err
	}
	return timestamp.Unix(), nil
}

// TodayZeroTimeStamp 获取指定时区当天的00:00:00的时间戳
// *time.Location 可以通过time.LoadLocation("UTC")来获取
func TodayZeroTimeStamp(timeZone *time.Location) int64 {
	t := time.Now().In(timeZone)
	addTime := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, timeZone)
	return addTime.Unix()
}
