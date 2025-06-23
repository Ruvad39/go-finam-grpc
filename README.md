# go-finam-grpc

**gRPC-клиент на Go для работы с API Финама**  
[tradeapi.finam](https://tradeapi.finam.ru/docs/about/)


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


### Примеры смотрите [тут](/_examples)


## TODO
* [ ] MarketDataServiceClient
* [ ] OrderServiceClient
* [ ] streams
