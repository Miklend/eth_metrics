package metrics

import (
	"encoding/json"
	"io"
	"net/http"
)

// Структура с данными по каждому протоколу
type ProtocolTvl struct {
	Name     string  `json:"name"`
	Tvl      float64 `json:"tvl"`
	Category string  `json:"category"`
}

// Функция получения и парсинга данных
func GetDataTvl(urlProtocols, urlChain string) ([]ProtocolTvl, error) {
	// Первый запрос — /protocols
	resp, err := http.Get(urlProtocols)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var rawData []ProtocolTvl
	if err := json.Unmarshal(body, &rawData); err != nil {
		return nil, err
	}

	var filteredData []ProtocolTvl
	for _, item := range rawData {
		if item.Category != "Chain" {
			filteredData = append(filteredData, item)
		}
	}

	// Второй запрос — /chains
	respCh, err := http.Get(urlChain)
	if err != nil {
		return nil, err
	}
	defer respCh.Body.Close()

	bodyCh, err := io.ReadAll(respCh.Body)
	if err != nil {
		return nil, err
	}

	var rawDataCh []struct {
		Name string  `json:"name"`
		Tvl  float64 `json:"tvl"`
	}
	if err := json.Unmarshal(bodyCh, &rawDataCh); err != nil {
		return nil, err
	}

	// Преобразуем в ProtocolTvl и добавляем
	for _, item := range rawDataCh {
		filteredData = append(filteredData, ProtocolTvl{
			Name:     item.Name,
			Tvl:      item.Tvl,
			Category: "Chain",
		})
	}

	return filteredData, nil
}
