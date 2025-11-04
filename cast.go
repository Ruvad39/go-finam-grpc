package finam

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/genproto/googleapis/type/decimal"
	"google.golang.org/genproto/googleapis/type/interval"
	"google.golang.org/genproto/googleapis/type/money"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var TzMoscow = initMoscow()

func initMoscow() *time.Location {
	var loc, err = time.LoadLocation("Europe/Moscow")
	if err != nil {
		loc = time.FixedZone("MSK", int(3*time.Hour/time.Second))
	}
	return loc
}

// IntToDecimal конвертируем int в google.Decimal
func IntToDecimal(i int) *decimal.Decimal {
	return &decimal.Decimal{
		Value: strconv.FormatInt(int64(i), 10),
	}
}

// Float64ToDecimal конвертируем  float64 в google.Decimal
func Float64ToDecimal(f float64) *decimal.Decimal {
	// Конвертируем float64 в строку с нужной точностью (например, 6 знаков после точки)
	// Можно использовать fmt.Sprintf("%.Nf", f) для фиксации количества знаков
	return &decimal.Decimal{
		Value: strconv.FormatFloat(f, 'f', -1, 64),
	}
}

func DecimalToString(d *decimal.Decimal) string {
	return d.GetValue()
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
// TODO в параметры добавить зону
func IsWithinInterval(t time.Time, iv *interval.Interval) bool {
	if iv == nil || iv.StartTime == nil || iv.EndTime == nil {
		return false
	}

	start := iv.StartTime.AsTime()
	end := iv.EndTime.AsTime()

	return !t.Before(start) && !t.After(end)
}

// ???
func CleanSymbolFromMic(symbol string) string {
	if idx := strings.Index(symbol, "@"); idx != -1 {
		return symbol[:idx]
	}
	return symbol
}

func valueOrZero(v *decimal.Decimal) string {
	if v == nil {
		return "0"
	}
	return v.Value
}

// TimestampToUnixMilli переведем в миллисекунды с эпохи Unix
func TimestampToUnixMilli(ts *timestamppb.Timestamp) int64 {
	if ts == nil {
		return 0
	}
	return ts.Seconds*1000 + int64(ts.Nanos)/1_000_000
}

// DateToTimeToTime преобразует google.type.Date в time.Time
func DateToTime(d *date.Date) time.Time {
	if d == nil {
		return time.Time{} // нулевое значение
	}

	return time.Date(
		int(d.Year),
		time.Month(d.Month),
		int(d.Day),
		0, 0, 0, 0, // время всегда 00:00:00
		time.UTC, // используем UTC для дат без времени
	)
}

// DateToTimeWithLocation преобразует с указанием локации
func DateToTimeWithLocation(d *date.Date, loc *time.Location) time.Time {
	if d == nil || loc == nil {
		return time.Time{}
	}

	return time.Date(
		int(d.Year),
		time.Month(d.Month),
		int(d.Day),
		0, 0, 0, 0,
		loc,
	)
}

// Float64ToString
func Float64ToString(in float64, inPrec ...int) string {
	prec := -1
	if len(inPrec) != 0 {
		prec = inPrec[0]
	}
	return strconv.FormatFloat(in, 'f', prec, 64) // 'f' - формат, prec - точность
}

func Float32ToString(in float32, inPrec ...int) string {
	prec := -1
	if len(inPrec) != 0 {
		prec = inPrec[0]
	}
	return strconv.FormatFloat(float64(in), 'f', prec, 32) // 'f' - формат, prec - точность
}

func Int64ToString(in int64) string {
	return strconv.FormatInt(in, 10)
}

func Int32ToString(in int32) string {
	return strconv.FormatInt(int64(in), 10)
}

func IntToString(in int) string {
	return strconv.FormatInt(int64(in), 10)
}
