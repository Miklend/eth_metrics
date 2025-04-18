package coingecko

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type CurrencyData struct {
	Name          string
	USD           float64 `json:"usd"`
	USDMarketCap  float64 `json:"usd_market_cap"`
	USD24hVol     float64 `json:"usd_24h_vol"`
	LastUpdatedAt int64   `json:"last_updated_at"`
}

type CryptoData map[string]CurrencyData

func capitalizeFirstLetter(s string) string {
	if len(s) == 0 {
		return s
	}
	if s == "binancecoin" {
		return "BSC"
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func GetDataVolMcapCoingecko(url, apikey string) ([]CurrencyData, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("x-cg-demo-api-key", apikey)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var cryptoData CryptoData
	if err := json.Unmarshal(body, &cryptoData); err != nil {
		return nil, err
	}

	var filterdata []CurrencyData
	for chain, data := range cryptoData {
		dataWithName := CurrencyData{
			Name:          capitalizeFirstLetter(chain),
			USD:           data.USD,
			USDMarketCap:  data.USDMarketCap,
			USD24hVol:     data.USD24hVol,
			LastUpdatedAt: data.LastUpdatedAt,
		}
		filterdata = append(filterdata, dataWithName)
	}

	return filterdata, nil
}
