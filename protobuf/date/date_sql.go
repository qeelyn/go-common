package date

import (
	"time"
	"errors"
	"database/sql/driver"
)

// Scan implements the Scanner interface.
func (dt *Date) Scan(value interface{}) error {
	t, ok := value.(time.Time)
	if !ok {
		return errors.New("not a time.Time type")
	}
	dt.Year = int32(t.Year())
	dt.Month = int32(t.Month())
	dt.Day = int32(t.Day())
	return nil
}

// Value implements the driver Valuer interface.
func (dt Date) Value() (driver.Value, error) {
	t := time.Date(int(dt.Year), time.Month(dt.Month), int(dt.Day), 0, 0, 0, 0, time.Local)
	return t, nil
}
