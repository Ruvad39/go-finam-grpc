# Changelog

## [0.6.0] - 2025-12-04
### New
- BarStream: отдельный метод для запуска Start
- BarStream: создание стрима с указанием callback функции NewBarStreamWithCallback
- BarStream: создание стрима с возвратом канала NewBarStreamWithChannel
- QuoteStream: стрим котировок
- QuoteStream: создание стрима с указанием callback функции NewQuoteStreamWithCallback
- QuoteStream: создание стрима с возвратом канала NewQuoteStreamWithChannel
- OrderStream: стрим получение ордеров NewOrderStream (NewOrderStreamWithCallback)
- TradeStream: стрим получение сделок NewTradeStream (NewTradeStreamWithCallback)

### Update
- обновление proto файлов до 2.9.0



## [0.5.0] - 2025-10-31
### New
- clientOptions: WithCallTimeout(value time.Duration) WithLogger(logger *slog.Logger)

### Update
- обновление proto файлов до 2.8.0
- добавление в методы вызова grpc context.WithTimeout


## [0.4.0] - 2025-08-29
### New
- TokenAgent (соответсвует интерфейсу PerRPCCredentials) (Спасибо Артему)
- при создание grpc.NewClient добавил WithPerRPCCredentials(tokenAgent)
- clientOptions: WithJwtRefreshInterval(value time.Duration) WithKeepaliveTime(value time.Duration) WithKeepaliveTimeout
- BarStream = подписка на свечи
- OrderTradeStream = Подписка на собственные заявки и сделки
- OrderBookStream = подписка на стакан

### Update
- убрал client.WithAuthToken


## [0.3.0] - 2025-08-28
### Update
- переместил proto-файлы в каталог \proto\grpc\tradeapi\v1\...
- обновление proto файлов до 2.7.0


## [0.2.0] - 2025-06-23
### Update
- Полная реорганизация проекта
- proto файлы все в одном каталоге tradeapi\v1
- для работы с сервисами теперь создаем отдельные "клиенты" 

### Added
- NewAccountServiceClient = создаем клиент для работы с AccountService
- NewAssetServiceClient = создаем клиент для работы с AssetService
- NewMarketDataServiceClient = создаем клиент для работы с MarketDataService
- NewOrderServiceClient  = создаем клиент для работы с OrdersService


## [0.1.0] - 2025-04-17
### Added
- реализован (почти) весь функционал для работы с api.
