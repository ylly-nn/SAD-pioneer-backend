package timeparsing

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
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
