# Неофициальный Go SDK для Finam Trade API

**gRPC-клиент создан на базе [proto файлов](https://github.com/FinamWeb/finam-trade-api)**

⚠️ **Важное предупреждение**  
Это **неофициальная** библиотека для работы с [Finam Trade API](https://tradeapi.finam.ru/docs/guides/grpc/).  
Возможны ошибки как в реализации SDK, так и в самом Trade API.  
**Используйте этот код на свой страх и риск.**  
Разработчик не несет ответственности за потери.


## Установка

```bash
go get github.com/Ruvad39/go-finam-grpc
```

## Примеры

### Пример создание клиента. Получение данных по токену
```go
ctx := context.Background()
token, _ := "FINAM_TOKEN"
client, err := finam.NewClient(ctx, token)
if err != nil {
    slog.Error("NewClient", "err", err.Error())
return
}
defer client.Close()

// Получение информации о токене сессии
res, err := client.GetTokenDetails(ctx)
if err != nil {
    slog.Error("main", "AuthService.TokenDetails", err.Error())
}
slog.Info("main", "res", res)


```

### Получить информацию по торговому счету
```go

// создаем клиент для доступа к AccountService
accountService := client.NewAccountServiceClient()

// запрос по заданному счету
res, err := accountServicet.GetAccount(ctx,  "FINAM_ACCOUNT_ID")
if err != nil {
    slog.Error("AccountsService.GetAccount", "GetAccount", err.Error())
    return
}
slog.Info("AccountsService.GetAccount",
    "AccountId", res.AccountId,
    "Type", res.Type,
    "Status", res.Status,
    "Equity", fmt.Sprintf("%.2f", finam.DecimalToFloat64(res.Equity)),
    "UnrealizedProfit", fmt.Sprintf("%.2f", finam.DecimalToFloat64(res.UnrealizedProfit)),
    "Cash", res.Cash,
)

// список позиций
for row, pos := range res.Positions {
    slog.Info("AccountsService.GetAccount.Positions",
        "row", row,
        "Symbol", pos.Symbol,
        "Quantity", finam.DecimalToFloat64(pos.Quantity),
        "AveragePrice", finam.DecimalToFloat64(pos.AveragePrice),
        "CurrentPrice", finam.DecimalToFloat64(pos.CurrentPrice),
    )
}	
```

### Подробные примеры смотрите [тут](/_examples)


### Методы реализованные на текущие момент

```go
// AssetServiceClient 
// вернем текущее время сервера (в TzMoscow)
GetTime(ctx context.Context) (time.Time, error)
// Получение списка доступных бирж, названия и mic коды
GetExchanges(ctx context.Context) (*pb.ExchangesResponse, error)
// Получение списка доступных инструментов, их описание
GetAssets(ctx context.Context) (*pb.AssetsResponse, error)
// Получение информации по конкретному инструменту
GetAsset(ctx context.Context, accountId, symbol string) (*pb.GetAssetResponse, error)
// Получение торговых параметров по инструменту
GetAssetParams(ctx context.Context, accountId, symbol string) (*pb.GetAssetParamsResponse, error)
// Получение расписания торгов для инструмента
GetSchedule(ctx context.Context, symbol string) (*pb.ScheduleResponse, error)
// TODO OptionsChain


// AccountServiceClient 
// Получение информации по конкретному счету
GetAccount(ctx context.Context, accountId string) (*pb.GetAccountResponse, error)
// Получение истории по сделкам заданного счета
GetTrades(ctx context.Context, accountId string, start, end time.Time, limit int32) (*pb.TradesResponse, error)
// Получение списка транзакций по счету
GetTransactions(ctx context.Context, accountId string, start, end time.Time, limit int32) (*pb.TransactionsResponse, error)


// MarketDataServiceClient 
//Получение исторических данных по инструменту (агрегированные свечи)
GetBars(ctx context.Context, symbol string, tf pb.TimeFrame, start, end time.Time)
// Получение последней котировки по инструменту
GetLastQuote(ctx context.Context, symbol string) (*pb.QuoteResponse, error)
// получение списка последних сделок по инструменту
GetLatestTrades(ctx context.Context, symbol string) (*pb.BarsResponse, error)
// Получение текущего стакана по инструменту
GetOrderBook(ctx context.Context, symbol string) (*pb.OrderBookResponse, error)


// OrderServiceClient
// Получение списка заявок по заданному счету
GetOrders(ctx context.Context, accountId string) (*pb.OrdersResponse, error)
// Получение информации о конкретном ордере
GetOrder(ctx context.Context, accountId, orderId string) (*pb.OrderState, error)
// Отмена биржевой заявки
CancelOrder(ctx context.Context, accountId, orderId string) (*pb.OrderState, error)
// Выставление биржевой заявки
PlaceOrder(ctx context.Context, order *pb.Order) (*pb.OrderState, error)

// Вспомогательные методы для создания ордера
// создать ордер на покупку по рынку
NewBuyOrder(accountId, symbol string, quantity int) *pb.Order {
// создать ордер на продажу по рынку
NewSellOrder(accountId, symbol string, quantity int) *pb.Order {
// создать ордер на покупку по лимитной цене
NewBuyLimitOrder(accountId, symbol string, quantity int, price float64) *pb.Order
// создать ордер на продажу по лимитной цене
NewSellLimitOrder(accountId, symbol string, quantity int, price float64) *pb.Order

// Потоки данных (stream)
// Подписка на собственные заявки и сделки
NewOrderTradeStream(parent context.Context, accountId string, callbackOrder func(*orders_service.OrderState),callbackTrade func(*v1.AccountTrade),) *OrderTradeStream
// создание стрима на стакан
NewOrderBookStream(parent context.Context, symbol string, callback func(book []*pb.StreamOrderBook)) *OrderBookStream
// создание стрима на агрегированные свечи
NewBarStream(parent context.Context, symbol string, timeframe marketdata_service.TimeFrame, callback func(bar *Bar)) *BarStream
```

## TODO
* [ ] SubscribeQuote
