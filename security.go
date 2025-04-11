/*
	локальное хранение списка торговых инструментов
*/

package finam

import (
	"context"
	assets_service "github.com/Ruvad39/go-finam-grpc/trade_api/v1/assets"
)

// Security Информация об инструменте
type Security struct {
	Ticker string `json:"ticker,omitempty"` // Тикер инструмента
	Symbol string `json:"symbol,omitempty"` // Символ инструмента ticker@mic
	Name   string `json:"name,omitempty"`   // Наименование инструмента
	Mic    string `json:"mic,omitempty"`    // Идентификатор биржи
	Type   string `json:"type,omitempty"`   // Тип инструмента
	Id     string `json:"id,omitempty"`     // Идентификатор инструмента
}

// TODO
type SecurityStory struct {
	Securities map[string]Security
}

// загрузим список торговых инструментов
// TODO  (1) с сервера 	assets, err := client.AssetsService.Assets(ctx, finam.NewAssetsRequest())
// TODO (2) с локального файла json (какой путь? какой формат? откуда файл взялся?)
func (c *Client) LoadSecurities(ctx context.Context) error {
	// ??? заполним s.Securities[symbol]
	return nil
}

// GetSecurity Вернем параметры инструмента (из локальных данных Client)
func (c *Client) GetSecurity(symbol string) (Security, bool) {
	sec, exists := c.Securities[symbol]
	if !exists {
		return Security{}, false
	}
	return sec, true
}

func ToSecurity(assets []*assets_service.Asset) []Security {
	result := make([]Security, len(assets))
	for i, asset := range assets {
		result[i] = Security{
			Ticker: asset.Ticker,
			Symbol: asset.Symbol,
			Name:   asset.Name,
			Mic:    asset.Mic,
			Type:   asset.Type,
			Id:     asset.Id,
		}
	}
	return result
}

// ToSecurityMap
func ToSecurityMap(assets []*assets_service.Asset) map[string]Security {
	result := make(map[string]Security, len(assets))
	for _, asset := range assets {
		s := Security{
			Ticker: asset.Ticker,
			Symbol: asset.Symbol,
			Name:   asset.Name,
			Mic:    asset.Mic,
			Type:   asset.Type,
			Id:     asset.Id,
		}
		result[s.Symbol] = s
	}
	return result
}

// SecurityToMap
func SecurityToMap(sec []Security) map[string]Security {
	result := make(map[string]Security, len(sec))
	for _, s := range sec {
		result[s.Symbol] = s
	}
	return result
}
