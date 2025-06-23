# Changelog


## [0.2.0] - 2025-06-23
### Update
- Полная реорганизация проекта
- proto файлы все в одном каталоге tradeapi\v1
- для работы с сервисами теперь создаем отдельные "клиенты" 

### Added
- NewAccountServiceClient = создаем клиент для доступа к AccountService
- NewAssetServiceClient = создаем клиент для доступа к AssetService
- NewMarketDataServiceClient = создаем клиент для работы MarketDataService


## [0.1.0] - 2025-04-17
### Added
- реализован (почти) весь функционал для работы с api.