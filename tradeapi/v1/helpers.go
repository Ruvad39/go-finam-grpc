package proto

import (
	"fmt"
	"google.golang.org/genproto/googleapis/type/decimal"
	"strconv"
)

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

func DecimalToIntE(d *decimal.Decimal) (int64, error) {
	if d == nil {
		return 0, nil
	}
	val, err := DecimalToFloat64E(d)
	if err != nil {
		return 0, err
	}
	return int64(val), nil
}

func DecimalToInt(d *decimal.Decimal) int64 {
	result, _ := DecimalToIntE(d)
	return result
}
