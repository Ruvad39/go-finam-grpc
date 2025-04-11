package finam

import (
	"fmt"
	"google.golang.org/genproto/googleapis/type/decimal"
	"google.golang.org/genproto/googleapis/type/interval"
	"google.golang.org/genproto/googleapis/type/money"
	"google.golang.org/protobuf/types/known/timestamppb"
	"math"
	"strconv"
	"time"
)

var TzMoscow = initMoscow()

func initMoscow() *time.Location {
	var loc, err = time.LoadLocation("Europe/Moscow")
	if err != nil {
		loc = time.FixedZone("MSK", int(3*time.Hour/time.Second))
	}
	return loc
}

// Float64ToDecimal конвертируем  float64 в google.Decimal
func Float64ToDecimal(f float64) *decimal.Decimal {
	// Конвертируем float64 в строку с нужной точностью (например, 6 знаков после точки)
	// Можно использовать fmt.Sprintf("%.Nf", f) для фиксации количества знаков
	return &decimal.Decimal{
		Value: strconv.FormatFloat(f, 'f', -1, 64),
	}
}

// DecimalToFloat64E конвертируем google.Decimal в float64
// с обработкой ошибки
func DecimalToFloat64E(d *decimal.Decimal) (float64, error) {
	if d == nil {
		return 0, fmt.Errorf("decimal is nil")
	}
	return strconv.ParseFloat(d.Value, 64)
}

// DecimalToFloat64 конвертируем google.Decimal в float64
// БЕЗ обработки ошибки
func DecimalToFloat64(d *decimal.Decimal) float64 {
	result, _ := DecimalToFloat64E(d)
	return result
}

func DecimalToIntE(d *decimal.Decimal) (int, error) {
	if d == nil {
		return 0, nil
	}
	val, err := DecimalToFloat64E(d)
	if err != nil {
		return 0, err
	}
	return int(val), nil
}

func DecimalToInt(d *decimal.Decimal) int {
	result, _ := DecimalToIntE(d)
	return result
}

// MoneyToFloat64 конвертируем google.money в float64
func MoneyToFloat64(m *money.Money) float64 {
	if m == nil {
		return 0
	}
	return float64(m.Units) + float64(m.Nanos)/1e9
}

func Float64ToMoney(value float64, currency string) *money.Money {
	units := int64(value)
	nanos := int32(math.Round((value - float64(units)) * 1e9))

	return &money.Money{
		CurrencyCode: currency,
		Units:        units,
		Nanos:        nanos,
	}
}

// NewInterval создадим google.interval
// start time.Time = StartTime
// end time.Time = EndTime
func NewInterval(start, end time.Time) *interval.Interval {
	result := &interval.Interval{
		StartTime: timestamppb.New(start),
		EndTime:   timestamppb.New(end),
	}
	return result
}

// Проверка: входит ли t в интервал [start, end]
func IsWithinInterval(t time.Time, iv *interval.Interval) bool {
	if iv == nil || iv.StartTime == nil || iv.EndTime == nil {
		return false
	}

	start := iv.StartTime.AsTime()
	end := iv.EndTime.AsTime()

	return !t.Before(start) && !t.After(end)
}

//func TimestampAsTime(x *Timestamp) AsTime() time.Time {
//	return time.Unix(int64(x.GetSeconds()), int64(x.GetNanos())).UTC()
//}

// AsTime converts x to a time.Time.
//func (x *Timestamp) AsTime() time.Time {
//	return time.Unix(int64(x.GetSeconds()), int64(x.GetNanos())).UTC()
//}
