package utils

import (
	"fmt"
	ptime "github.com/yaa110/go-persian-calendar"
	"strconv"
	"strings"
	"time"
)

var DateType string

const (
	Year   = "year"
	Month  = "month"
	Day    = "day"
	Hour   = "hour"
	Minute = "minute"
	Second = "second"
)

func GetJalaliTimeDateFrom(t time.Time) (string, string) {
	pt := ptime.New(t)
	date := pt.Format("yyyy/MM/dd")

	h, m, s := pt.Clock()
	tm := fmt.Sprintf("%02d:%02d:%02d", h, m, s)

	return date, tm
}

func GetStringTimeJalali() (string, string) {
	pt := ptime.New(time.Now())
	date := pt.Format("yyyy/MM/dd")

	h, m, s := pt.Clock()
	time := fmt.Sprintf("%d:%d:%d", h, m, s)

	return date, time
}

func GetJalaliDate(inpTime time.Time) string {
	pt := ptime.New(inpTime)
	date := pt.Format("yyyy/MM/dd")

	return date
}

func GetJalaliTime(inpTime time.Time) string {
	pt := ptime.New(inpTime)
	h, m, s := pt.Clock()
	tm := fmt.Sprintf("%02d:%02d:%02d", h, m, s)

	return tm
}

func ConvertStringDateToTime(dateStr, timeStr, until string, location *time.Location) time.Time {

	switch until {
	case Day:
		arrayDate := strings.Split(dateStr, "/")
		year, _ := strconv.Atoi(arrayDate[0])
		month, _ := strconv.Atoi(arrayDate[1])
		day, _ := strconv.Atoi(arrayDate[2])
		return ptime.Date(year, ptime.Month(month), day, 0, 0, 0, 0, location).Time()

	case Minute:
		arrayDate := strings.Split(dateStr, "/")
		arrayTime := strings.Split(timeStr, ":")
		year, _ := strconv.Atoi(arrayDate[0])
		month, _ := strconv.Atoi(arrayDate[1])
		day, _ := strconv.Atoi(arrayDate[2])
		createdAtHour, _ := strconv.Atoi(arrayTime[0])
		createdAtMinute, _ := strconv.Atoi(arrayTime[1])
		return ptime.Date(year, ptime.Month(month), day, createdAtHour, createdAtMinute, 0, 0, location).Time()
	}

	return time.Time{}
}

func MonthDiff(start, end time.Time) int64 {
	// Extract year + month
	y1, m1 := start.Year(), int(start.Month())
	y2, m2 := end.Year(), int(end.Month())

	// Calculate total month difference
	months := (y2-y1)*12 + (m2 - m1)
	return int64(months)
}

func CombineTimes(datePart *time.Time, timePart *time.Time) *time.Time {
	combined := time.Date(
		datePart.Year(),
		datePart.Month(),
		datePart.Day(),
		timePart.Hour(),
		timePart.Minute(),
		timePart.Second(),
		timePart.Nanosecond(),
		datePart.Location(),
	)
	return &combined
}
