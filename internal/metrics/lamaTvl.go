package metrics

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

// Структура с данными по каждому протоколу
type TvlChains struct {
	Name     string  `json:"name"`
	Tvl      float64 `json:"tvl"`
	Category string  `json:"category"`
}

// Функция получения и парсинга данных
func GetDataTvlChains(urlChain string) ([]TvlChains, error) {
	var filteredData []TvlChains

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

	for _, item := range rawDataCh {
		normalized := strings.ToLower(item.Name)
		if _, ok := chainSet[normalized]; !ok {
			continue
		}

		filteredData = append(filteredData, TvlChains{
			Name:     normalized, // можно оставить original Name, если нужно
			Tvl:      item.Tvl,
			Category: "Chain",
		})
	}

	return filteredData, nil
}

type ProtocolTvl struct {
	Name      string             `json:"name"`
	Tvl       float64            `json:"tvl"`
	Mcap      float64            `json:"mcap"`
	Category  string             `json:"category"`
	ChainTVLs map[string]float64 `json:"chainTvls"`
}

func GetDataTvlProtocols(url string) ([]ProtocolTvl, []ProtocolTvl, error) {
	// Список поддерживаемых сетей (основные + их варианты)
	supportedChains := map[string]bool{
		// Ethereum и его варианты
		"Ethereum": true, "Ethereum_borrowed": true, "Ethereum_pool2": true,
		"Ethereum_staking": true, "Ethereum_treasury": true, "Ethereum_vesting": true,
		"Ethereum_OwnTokens": true,

		// Solana и её варианты
		"Solana": true, "Solana_borrowed": true, "Solana_pool2": true,
		"Solana_staking": true, "Solana_vesting": true,

		// BSC (Binance Smart Chain) и её варианты
		"Binance": true, "Binance_borrowed": true, "Binance_pool2": true,
		"Binance_staking": true, "Binance_vesting": true,

		// Bitcoin и его варианты
		"Bitcoin": true, "Bitcoin_staking": true,

		// Tron и его варианты
		"Tron": true, "Tron_borrowed": true, "Tron_pool2": true, "Tron_staking": true,
	}

	resp, err := http.Get(url)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	var rawData []ProtocolTvl
	if err := json.Unmarshal(body, &rawData); err != nil {
		return nil, nil, err
	}

	var filteredDataMcap []ProtocolTvl
	var filteredDataTvl []ProtocolTvl

	for _, item := range rawData {
		// Создаем новый filteredChainTVLs только с поддерживаемыми сетями
		filteredChainTVLs := make(map[string]float64)

		for chain, value := range item.ChainTVLs {
			// Нормализуем название сети
			normalizedChain := strings.ReplaceAll(chain, "-", "_")
			normalizedChain = strings.ReplaceAll(normalizedChain, " ", "_")

			// Если сеть в списке поддерживаемых - добавляем
			if _, ok := supportedChains[normalizedChain]; ok {
				filteredChainTVLs[normalizedChain] = value
			}
		}

		// Если есть хотя бы одна поддерживаемая сеть
		if len(filteredChainTVLs) > 0 {
			normalizedItem := ProtocolTvl{
				Name:      item.Name,
				Tvl:       item.Tvl,
				Mcap:      item.Mcap,
				Category:  item.Category,
				ChainTVLs: filteredChainTVLs, // Только поддерживаемые сети
			}

			if item.Category != "Chain" && item.Tvl > 0 {
				filteredDataTvl = append(filteredDataTvl, normalizedItem)
			}
			if item.Mcap > 0 {
				filteredDataMcap = append(filteredDataMcap, normalizedItem)
			}
		}
	}
	return filteredDataTvl, filteredDataMcap, nil
}
