/*
MarketDataServiceClient = клиент для работы с MarketDataService
методы:

//Получение исторических данных по инструменту (агрегированные свечи)
GetBars(ctx context.Context, symbol string, tf pb.TimeFrame, start, end time.Time)

// Получение последней котировки по инструменту
GetLastQuote(ctx context.Context, symbol string)

// получение списка последних сделок по инструменту
GetLatestTrades(ctx context.Context, symbol string)

// Получение текущего стакана по инструменту
 GetOrderBook(ctx context.Context, symbol string)

*/

package finam

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	marketdata_service "github.com/Ruvad39/go-finam-grpc/proto/grpc/tradeapi/v1/marketdata"
	"google.golang.org/genproto/googleapis/type/decimal"
	"google.golang.org/genproto/googleapis/type/interval"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// Bar аналог marketdata_service.Bar
// плюс доп. поля: Symbol + Timeframe
type Bar struct {
	// Инструмент
	Symbol string
	// Таймфрейм
	Timeframe marketdata_service.TimeFrame
	// Метка времени
	Timestamp *timestamppb.Timestamp
	// Цена открытия свечи
	Open *decimal.Decimal
	// Максимальная цена свечи
	High *decimal.Decimal
	// Минимальная цена свечи
	Low *decimal.Decimal
	// Цена закрытия свечи
	Close *decimal.Decimal
	// Объём торгов за свечу в шт.
	Volume *decimal.Decimal `protobuf:"bytes,6,opt,name=volume,proto3" json:"volume,omitempty"`
}

func (b *Bar) GetOpen() float64 {
	return DecimalToFloat64(b.Open)
}
func (b *Bar) GetHigh() float64 {
	return DecimalToFloat64(b.High)
}
func (b *Bar) GetLow() float64 {
	return DecimalToFloat64(b.Low)
}
func (b *Bar) GetClose() float64 {
	return DecimalToFloat64(b.Close)
}
func (b *Bar) GetVolume() int {
	return DecimalToInt(b.Volume)
}
func (b *Bar) GetTime() time.Time {
	return b.Timestamp.AsTime().In(TzMoscow)
}

func (b *Bar) String() string {
	str := fmt.Sprintf("%v,%v,%v, Open:%v, High:%v, Low:%v, CLose:%v, Volume:%v",
		b.GetTime().String(), b.Symbol, b.Timeframe.String(), b.GetOpen(), b.GetHigh(), b.GetLow(), b.GetClose(), b.GetVolume())
	return str
}

type BarFunc func(bar *Bar)

// MarketDataServiceClient клиент для работы MarketDataService
type MarketDataServiceClient struct {
	client            *Client
	MarketDataService marketdata_service.MarketDataServiceClient
}

func NewMarketDataServiceClient(c *Client) *MarketDataServiceClient {
	return &MarketDataServiceClient{client: c,
		MarketDataService: marketdata_service.NewMarketDataServiceClient(c.conn),
	}
}

// GetLastQuote Получение последней котировки по инструменту
func (s *MarketDataServiceClient) GetLastQuote(ctx context.Context, symbol string) (*marketdata_service.QuoteResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.client.opts.callTimeout)
	defer cancel()
	return s.MarketDataService.LastQuote(ctx, &marketdata_service.QuoteRequest{Symbol: symbol})
}

// GetBars Получение исторических данных по инструменту (агрегированные свечи)
// symbol string
// tf pb.TimeFrame
// start, end time.Time
func (s *MarketDataServiceClient) GetBars(ctx context.Context, symbol string, tf marketdata_service.TimeFrame, start, end time.Time) (*marketdata_service.BarsResponse, error) {
	i := &interval.Interval{
		StartTime: timestamppb.New(start),
		EndTime:   timestamppb.New(end),
	}
	ctx, cancel := context.WithTimeout(ctx, s.client.opts.callTimeout)
	defer cancel()
	return s.MarketDataService.Bars(ctx, &marketdata_service.BarsRequest{Symbol: symbol, Timeframe: tf, Interval: i})
}

// GetOrderBook Получение текущего стакана по инструменту
func (s *MarketDataServiceClient) GetOrderBook(ctx context.Context, symbol string) (*marketdata_service.OrderBookResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.client.opts.callTimeout)
	defer cancel()
	return s.MarketDataService.OrderBook(ctx, &marketdata_service.OrderBookRequest{Symbol: symbol})
}

// GetLatestTradesПолучение списка последних сделок по инструменту
func (s *MarketDataServiceClient) GetLatestTrades(ctx context.Context, symbol string) (*marketdata_service.LatestTradesResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.client.opts.callTimeout)
	defer cancel()
	return s.MarketDataService.LatestTrades(ctx, &marketdata_service.LatestTradesRequest{Symbol: symbol})
}

// GetHistoryBars Получение исторических данных по инструменту (агрегированные свечи)
//
// разбиваем период на интервалы (от глубины рынка) и делаем несколько запросов к брокеру
func (s *MarketDataServiceClient) GetHistoryBars(ctx context.Context,
	symbol string,
	tf marketdata_service.TimeFrame,
	start, end time.Time,
	limit int,
) (*marketdata_service.BarsResponse, error) {
	// найдем возможную глубину истории
	duration := SelectDuration(tf)

	intervals := SplitInterval(start, end, duration)

	// идем с конца слайса так как там более раннее время
	//log.Info("GetHistoryBar", "duration", duration, "len", len(intervals), "intervals", intervals)
	bars := make([]*marketdata_service.Bar, 0)
	requests := 0
	for index, value := range intervals {
		requests++
		log.Debug("GetHistoryBar", "requests", requests)
		log.Debug("GetHistoryBar", "index", index, "start", value.GetStartTime().AsTime(), "end", value.GetEndTime().AsTime())
		resp, err := s.MarketDataService.Bars(ctx, &marketdata_service.BarsRequest{Symbol: symbol, Timeframe: tf, Interval: value})
		if err != nil {
			return &marketdata_service.BarsResponse{}, err
		}

		// если уже лимит запросов = ждем 1 минуту (200 запросов в минуту)
		if requests == 190 {
			log.Debug("GetHistoryBar = time.Sleep(time.Second * 60)")
			time.Sleep(time.Second * 60)
			requests = 0
		}
		if len(resp.GetBars()) < 1 {
			continue
		}
		// ??? берем со второй свечи, так данная свеча уже была последней в прошлом запросе
		//bars = append(bars, resp.GetBars()[1:]...)
		bars = append(bars, resp.GetBars()...)
	}
	// обрежем до нужного кол-ва
	if limit != 0 {
		end := len(bars)
		if end > limit {
			start := end - limit
			kn := (bars)[start:]
			bars = kn
		}
	}

	return &marketdata_service.BarsResponse{Symbol: symbol, Bars: bars}, nil

}

// selectDuration вернем возможную глубину данных
//
// 1 минута = 7 дней.
// 5 минут = 30 дней.
// 15 минут = 30 дней.
// 30 минут = 30 дней.
// 1 час = 30 дней.
// 2 часа = 30 дней.
// 4 часа = 30 дней.
// 8 часов = 30 дней.
// День = 365 дней.
// Неделя = 365*5 дней.
// Месяц = 365*5 дней.
// Квартал = 365*5 дней.
func SelectDuration(tf marketdata_service.TimeFrame) time.Duration {
	var duration time.Duration
	switch tf {
	case marketdata_service.TimeFrame_TIME_FRAME_M1:
		duration = time.Hour * 24
	case marketdata_service.TimeFrame_TIME_FRAME_M5, marketdata_service.TimeFrame_TIME_FRAME_M15, marketdata_service.TimeFrame_TIME_FRAME_M30:
		duration = time.Hour * 24 * 30
	case marketdata_service.TimeFrame_TIME_FRAME_H1, marketdata_service.TimeFrame_TIME_FRAME_H2, marketdata_service.TimeFrame_TIME_FRAME_H4, marketdata_service.TimeFrame_TIME_FRAME_H8:
		duration = time.Hour * 24 * 30
	case marketdata_service.TimeFrame_TIME_FRAME_D:
		duration = time.Hour * 24 * 365
	case marketdata_service.TimeFrame_TIME_FRAME_W:
		duration = time.Hour * 24 * 365 * 5
	case marketdata_service.TimeFrame_TIME_FRAME_MN:
		duration = time.Hour * 24 * 365 * 5
	case marketdata_service.TimeFrame_TIME_FRAME_QR:
		duration = time.Hour * 24 * 365 * 5

	}
	return duration
}

// SplitInterval разобъем интервал на диапазоны по размеру duration
func SplitInterval(start, end time.Time, duration time.Duration) []*interval.Interval {
	if start.After(end) || duration <= 0 {
		return nil
	}

	intervals := make([]*interval.Interval, 0)
	currentStart := start

	for currentStart.Before(end) {
		currentEnd := currentStart.Add(duration)

		// Если следующий интервал выходит за границы, используем конечное время
		if currentEnd.After(end) {
			currentEnd = end
		}

		intervals = append(intervals, &interval.Interval{
			StartTime: timestamppb.New(currentStart),
			EndTime:   timestamppb.New(currentEnd),
		})

		currentStart = currentEnd

		// Выходим если достигли конца
		if currentStart.Equal(end) || currentStart.After(end) {
			break
		}
	}

	return intervals
}

// Метод записи в .csv файл исторических свечей в формате
//
// <TICKER>,<PER>,<DATE>,<TIME>,<OPEN>,<HIGH>,<LOW>,<CLOSE>,<VOL>
//
// AFKS,1,20080109,103200,41.260,41.260,41.260,41.260,1
func WriteBarsToFile(bars *marketdata_service.BarsResponse, tf marketdata_service.TimeFrame, filename string) error {
	file, err := os.Create(fmt.Sprintf("%v.csv", filename))
	if err != nil {
		return err
	}
	// заголовок
	_, err = fmt.Fprintf(file, "<TICKER>,<PER>,<DATE>,<TIME>,<OPEN>,<HIGH>,<LOW>,<CLOSE>,<VOL>\n")
	defer func() {
		err = file.Close()
		if err != nil {
			log.Error("WriteBarsToFile", "err", err.Error())
		}
	}()
	var b strings.Builder
	symbol := CleanSymbolFromMic(bars.GetSymbol())
	per := TimeFrameToString(tf)
	for _, bar := range bars.GetBars() {
		b.Reset()
		b.WriteString(symbol + ",")
		b.WriteString(per + ",")
		b.WriteString(bar.Timestamp.AsTime().In(TzMoscow).Format("2006.01.02") + ",")
		b.WriteString(bar.Timestamp.AsTime().In(TzMoscow).Format("15:04:05") + ",")
		b.WriteString(bar.GetOpen().GetValue() + "," + bar.GetHigh().GetValue() + "," + bar.GetLow().GetValue() + "," + bar.GetClose().GetValue() + ",")
		b.WriteString(IntToString(DecimalToInt(bar.GetVolume())))
		//b.WriteString("\n")
		_, err = fmt.Fprintln(file, b.String())
		if err != nil {
			return err
		}
	}
	return nil
}

var TimeFrameString = map[int32]string{
	0:  "UNSPECIFIED",
	1:  "M1",
	5:  "M5",
	9:  "M15",
	11: "30",
	12: "H1",
	13: "H2",
	15: "H4",
	17: "H8",
	19: "D1",
	20: "W",
	21: "MN",
	22: "QR",
}

func TimeFrameToString(tf marketdata_service.TimeFrame) string {
	res := TimeFrameString[int32(tf)]
	return res
}
