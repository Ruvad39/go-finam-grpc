package proto

// Bar
// OpenAsFloat  CloseAsFloat HighAsFloat LowAsFloat VolumeAsFloat
func (bar *Bar) OpenAsFloat() float64 {
	return DecimalToFloat64(bar.Open)
}

func (bar *Bar) CloseAsFloat() float64 {
	return DecimalToFloat64(bar.Close)
}

func (bar *Bar) HighAsFloat() float64 {
	return DecimalToFloat64(bar.High)
}

func (bar *Bar) LowAsFloat() float64 {
	return DecimalToFloat64(bar.Low)
}

func (bar *Bar) VolumeAsInt() int64 {
	return DecimalToInt(bar.Volume)
}

// Quote
// Ask AskSize Bid  BidSize Last LastSize
func (q *Quote) AskAsFloat() float64 {
	return DecimalToFloat64(q.Ask)
}
func (q *Quote) BidAsFloat() float64 {
	return DecimalToFloat64(q.Bid)
}
func (q *Quote) LastAsFloat() float64 {
	return DecimalToFloat64(q.Last)
}
func (q *Quote) OpenAsFloat() float64 {
	return DecimalToFloat64(q.Open)
}
func (q *Quote) CloseAsFloat() float64 {
	return DecimalToFloat64(q.Close)
}
func (q *Quote) HighAsFloat() float64 {
	return DecimalToFloat64(q.High)
}
func (q *Quote) LowAsFloat() float64 {
	return DecimalToFloat64(q.Low)
}
func (q *Quote) VolumeAsInt() int64 {
	return DecimalToInt(q.Volume)
}

// OrderBook_Row PriceAsFloat SellSizeAsInt BuySizeAsInt
func (q *OrderBook_Row) PriceAsFloat() float64 {
	return DecimalToFloat64(q.Price)
}
func (q *OrderBook_Row) SellSizeAsInt() int64 {
	return DecimalToInt(q.GetSellSize())
}
func (q *OrderBook_Row) BuySizeAsInt() int64 {
	return DecimalToInt(q.GetBuySize())
}

// StreamOrderBook_Row PriceAsFloat
func (q *StreamOrderBook_Row) PriceAsFloat() float64 {
	return DecimalToFloat64(q.Price)
}
func (q *StreamOrderBook_Row) SellSizeAsInt() int64 {
	return DecimalToInt(q.GetSellSize())
}
func (q *StreamOrderBook_Row) BuySizeAsInt() int64 {
	return DecimalToInt(q.GetBuySize())
}
