package metrics

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Структура одного протокола
type Protocol struct {
	Name        string                        `json:"name"`
	Category    string                        `json:"category"`
	Total24h    float64                       `json:"total24h"`
	Breakdown24 map[string]map[string]float64 `json:"breakdown24h"`
}

// Упрощённая структура для вывода
type ProtocolInfo struct {
	Name        string
	Category    string
	Total24h    float64
	Breakdown24 map[string]float64 // Chain -> amount
}

// Множество разрешённых сетей
var chainSet = map[string]struct{}{
	"ethereum": {},
	"bitcoin":  {},
	"solana":   {},
	"bsc":      {},
	"tron":     {},
}

// Получение и парсинг данных
func GetDataVFR(url string) ([]ProtocolInfo, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения тела: %w", err)
	}

	// Структура с корневым элементом "protocols"
	var raw struct {
		Protocols []Protocol `json:"protocols"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return nil, fmt.Errorf("ошибка парсинга JSON: %w", err)
	}

	// Обработка и сбор нужной информации
	var result []ProtocolInfo
	for _, p := range raw.Protocols {
		breakdown := make(map[string]float64)
		hasAllowedChain := false
		for chain, apps := range p.Breakdown24 {
			// Заменяем "-" на "_"
			chain = strings.ReplaceAll(chain, "-", "_")

			if _, ok := chainSet[chain]; !ok {
				continue
			}
			hasAllowedChain = true

			for _, value := range apps {
				breakdown[chain] += value
			}
		}

		if !hasAllowedChain || p.Total24h <= 0 {
			continue
		}

		result = append(result, ProtocolInfo{
			Name:        p.Name,
			Category:    p.Category,
			Total24h:    p.Total24h,
			Breakdown24: breakdown,
		})
	}

	return result, nil
}
