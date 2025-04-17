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
```
// добавим заголовок с авторизацией (accessToken)
ctx, err = client.WithAuthToken(ctx)
if err != nil {
	slog.Error("main", "WithAuthToken", err.Error())
	// если прошла ошибка, дальше работа бесполезна, не будет авторизации
	return
}

// Получение информации по конкретному аккаунту
accountId := "FINAM_ACCOUNT_ID"
res, err := client.AccountsService.GetAccount(ctx, &accounts_service.GetAccountRequest{AccountId: accountId})
if err != nil {
	slog.Error("accountService", "GetAccount", err.Error())
}
slog.Info("main", "Account", res)

// список позиций
//slog.Info("main", "Positions", res.Positions)
for row, pos := range res.Positions {
	slog.Info("positions",
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
* [ ] MarketDataService.SubscribeLatestTrades
