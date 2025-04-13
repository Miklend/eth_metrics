package metrics

import (
	"encoding/json"
	"io"
	"net/http"
)

// Структура с данными по каждому протоколу
type ProtocolVFR struct {
	Name     string  `json:"name"`
	Total24h float64 `json:"total24h"`
	Category string  `json:"category"`
}

// Обертка для списка протоколов
type DexOverviewVFR struct {
	Protocols []ProtocolVFR `json:"protocols"`
}

// Функция получения и парсинга данных
func GetDataVFR(url string) (*DexOverviewVFR, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data DexOverviewVFR
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}
	return &data, nil
}
