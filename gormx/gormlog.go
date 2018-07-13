package gormx

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
	"unicode"
)

// gorm zap GormLog
type GormLog struct {
	occurredAt time.Time
	source     string
	duration   time.Duration
	sql        string
	values     []string
	other      []string
}

func (l *GormLog) ToZapFields() []zapcore.Field {
	return []zapcore.Field{
		zap.Time("occurredAt", l.occurredAt),
		zap.String("source", l.source),
		zap.Duration("duration", l.duration),
		zap.String("sql", l.sql),
		zap.Strings("values", l.values),
		zap.Strings("other", l.other),
	}
}

func CreateGormLog(values []interface{}) *GormLog {
	ret := &GormLog{}
	ret.occurredAt = gorm.NowFunc()

	if len(values) > 1 {
		var level = values[0]
		ret.source = getSource(values)

		if level == "sql" {
			ret.duration = getDuration(values)
			ret.values = getFormattedValues(values)
			ret.sql = values[3].(string)
		} else {
			ret.other = append(ret.other, fmt.Sprint(values[2:]))
		}
	}

	return ret
}

func getSource(values []interface{}) string {
	return fmt.Sprint(values[1])
}

func getDuration(values []interface{}) time.Duration {
	return values[2].(time.Duration)
}

func getFormattedValues(values []interface{}) []string {
	rawValues := values[4].([]interface{})
	formattedValues := make([]string, 0, len(rawValues))
	for _, value := range rawValues {
		switch v := value.(type) {
		case time.Time:
			formattedValues = append(formattedValues, fmt.Sprint(v))
		case []byte:
			if str := string(v); isPrintable(str) {
				formattedValues = append(formattedValues, fmt.Sprint(str))
			} else {
				formattedValues = append(formattedValues, "<binary>")
			}
		default:
			str := "NULL"
			if v != nil {
				str = fmt.Sprint(v)
			}
			formattedValues = append(formattedValues, str)
		}
	}
	return formattedValues
}

func isPrintable(s string) bool {
	for _, r := range s {
		if !unicode.IsPrint(r) {
			return false
		}
	}
	return true
}
