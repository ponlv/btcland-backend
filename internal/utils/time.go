package utils

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func ParseDate(dateStr string) (time.Time, error) {
	// Định dạng layout cho chuỗi thời gian "dd/MM/yyyy"
	const layout = "02/01/2006"

	// Parse chuỗi thời gian thành time.Time
	parsedTime, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, fmt.Errorf("error parsing date: %v", err)
	}

	return parsedTime, nil
}

func FormatTime(t time.Time) string {
	// Định dạng layout cho "yyyy-MM-dd HH:mm:ss.sss"
	const layout = "2006-01-02 15:04:05.000"
	return t.Format(layout)
}

func TimeToString(t time.Time) string {
	return fmt.Sprintf("%02d-%02d-%d", t.Day(), t.Month(), t.Year())
}

func ParseDateComponents(dateStr string) (time.Time, error) {
	// Tách chuỗi theo dấu "/"
	parts := strings.Split(dateStr, "/")
	if len(parts) != 3 {
		return time.Time{}, fmt.Errorf("định dạng ngày không hợp lệ")
	}

	// Chuyển các phần tử thành số nguyên
	day, err := strconv.Atoi(parts[0])
	if err != nil {
		return time.Time{}, fmt.Errorf("ngày không hợp lệ: %v", err)
	}

	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return time.Time{}, fmt.Errorf("tháng không hợp lệ: %v", err)
	}

	year, err := strconv.Atoi(parts[2])
	if err != nil {
		return time.Time{}, fmt.Errorf("năm không hợp lệ: %v", err)
	}

	// Tạo time.Time từ các thành phần
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return t, nil
}

// NormalizeDate use for convert date to time.Time in Golang
// dataStr is string of date exp: 10/09/1998 , 10-09-1998
// separate is symbol
func NormalizeDate(dataStr string, separate string) (time.Time, error) {
	normalize := strings.Replace(dataStr, separate, "", -1)

	if len(normalize) != 8 {
		return time.Time{}, errors.New("datetime not match")
	}

	day, err := strconv.Atoi(normalize[0:2])
	if err != nil {
		return time.Time{}, fmt.Errorf("ngày không hợp lệ: %v", err)
	}

	month, err := strconv.Atoi(normalize[2:4])
	if err != nil {
		return time.Time{}, fmt.Errorf("tháng không hợp lệ: %v", err)
	}

	year, err := strconv.Atoi(normalize[4:8])
	if err != nil {
		return time.Time{}, fmt.Errorf("năm không hợp lệ: %v", err)
	}

	// Tạo time.Time từ các thành phần
	t := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return t, nil
}
