package utils

import (
	"fmt"
	"time"
)

// ParseDateRange แปลง startDate, endDate (YYYY-MM-DD) เป็น time.Time ครอบคลุมทั้งวัน
func ParseDateRange(startDate, endDate string) (time.Time, time.Time, error) {
	layout := "2006-01-02"
	loc, _ := time.LoadLocation("Asia/Bangkok")

	start, err := time.ParseInLocation(layout, startDate, loc)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid startDate format: %v", err)
	}
	end, err := time.ParseInLocation(layout, endDate, loc)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid endDate format: %v", err)
	}

	// เพิ่มเวลาจนถึง 23:59:59
	end = end.Add(time.Hour*23 + time.Minute*59 + time.Second*59)

	return start, end, nil
}
