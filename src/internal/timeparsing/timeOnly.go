package timeparsing

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// TimeOnly представляет время суток с часовым поясом, без даты.
type TimeOnly time.Time

// Scan реализует sql.Scanner для чтения из БД.
func (t *TimeOnly) Scan(src any) error {
	if src == nil {
		*t = TimeOnly(time.Time{})
		return nil
	}
	var s string
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fmt.Errorf("cannot scan type %T into TimeOnly", src)
	}
	// Ожидаемый формат от PostgreSQL: "15:04:05-07" (например, "10:00:00+04")
	parsed, err := time.Parse("15:04:05-07", s)
	if err != nil {
		return err
	}
	*t = TimeOnly(parsed)
	return nil
}

// MarshalJSON реализует json.Marshaler.
// Возвращает время со смещением в формате "15:04:05-07:00".
func (t TimeOnly) MarshalJSON() ([]byte, error) {
	// Приводим к time.Time и форматируем с исходным смещением
	formatted := time.Time(t).Format("15:04:05-07:00")
	return json.Marshal(formatted)
}

// Value реализует driver.Valuer для возможной записи в БД (опционально).
func (t TimeOnly) Value() (driver.Value, error) {
	return time.Time(t).Format("15:04:05-07:00"), nil
}

// UnmarshalJSON реализует json.Unmarshaler
func (t *TimeOnly) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "" || s == "null" {
		*t = TimeOnly(time.Time{})
		return nil
	}

	// Форматы
	formats := []string{
		"15:04:05-07:00", // 10:00:00+03:00
		"15:04:05-07",    // 10:00:00+03
		"15:04:05Z07:00", // 10:00:00Z03:00
		"15:04:05",       // 10:00:00
	}

	for _, format := range formats {
		parsed, err := time.Parse(format, s)
		if err == nil {
			*t = TimeOnly(parsed)
			return nil
		}
	}

	return fmt.Errorf("invalid time format: %s", s)
}

func ParseLocation(tz string) (*time.Location, error) {
	// Пытаемся как IANA-зону
	if loc, err := time.LoadLocation(tz); err == nil {
		return loc, nil
	}
	// Пытаемся как фиксированное смещение
	offset, err := ParseOffset(tz)
	if err != nil {
		return nil, err
	}
	return time.FixedZone(tz, offset), nil
}

// parseOffset разбирает строку вида "+03:00" или "-05:00" и возвращает смещение в секундах.
func ParseOffset(s string) (int, error) {
	if len(s) < 3 {
		return 0, fmt.Errorf("invalid offset")
	}
	sign := 1
	if s[0] == '-' {
		sign = -1
	} else if s[0] != '+' {
		return 0, fmt.Errorf("invalid offset sign")
	}
	s = s[1:]
	parts := strings.SplitN(s, ":", 2)
	if len(parts) != 2 {
		return 0, fmt.Errorf("offset must be HH:MM")
	}
	hours, err := strconv.Atoi(parts[0])
	if err != nil || hours < 0 || hours > 23 {
		return 0, fmt.Errorf("invalid hour")
	}
	mins, err := strconv.Atoi(parts[1])
	if err != nil || mins < 0 || mins > 59 {
		return 0, fmt.Errorf("invalid minute")
	}
	return sign * (hours*3600 + mins*60), nil
}
