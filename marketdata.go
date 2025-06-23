package finam

import (
	"time"
)

// Информация о котировке
type Quote struct {
	Symbol    string // Символ инструмента
	Timestamp int64  // Метка времени
	Ask       float64
	Bid       float64
	Last      float64 // Цена последней сделки
}

func (q *Quote) Reset() {
	*q = Quote{} // копируем в q значение "пустого" объекта
}

func (q *Quote) Time() time.Time {
	return time.Unix(0, q.Timestamp).In(TzMoscow)
}
