package btcmarketsmodule

import (
	"encoding/json"
	"fmt"

	"btcmwatch.local/btcmarketsclient"
)

type handle struct {
	btcclient *btcmarketsclient.Config
}

type CurrencyBalance struct {
	Balance        float64 `json:"balance"`
	PendingFunds   float64 `json:"pendingFunds"`
	Currency       string  `json:"currency"`
	CurrentPricing MarketData
}

type MarketData struct {
	BestBid    float64 `json:"bestBid"`
	BestAsk    float64 `json:"bestAsk"`
	LastPrice  float64 `json:"lastPrice"`
	Currency   string  `json:"currency"`
	Instrument string  `json:"instument"`
	Timestamp  int     `json:"timestamp"`
	Volume24h  float64 `json:"volume24h"`
	Price24h   float64 `json:"price24h"`
	Low24h     float64 `json:"low24h"`
	High24h    float64 `json:"high24h"`
}

func CreateInstance(key string, secret string) *handle {
	h := handle{
		btcclient: btcmarketsclient.Configuration(key, secret),
	}

	return &h
}

func (h handle) PrintCurrencies() error {
	currencyBalances, err := h.getAccountBalance()
	if err != nil {
		return err
	}

	for _, cur := range currencyBalances {
		currencyData, err := h.getCurrentPrices(cur.Currency)
		if err != nil {
			return err
		}
		cur.CurrentPricing = currencyData
	}

	h.printValues(currencyBalances)

	return nil
}

func (h handle) printValues(currencyBalances []*CurrencyBalance) error {

	var total float64

	fmt.Println("Cur\tBalance\t\tLast Price\tValue")
	for _, cur := range currencyBalances {
		if cur.Balance == 0 || cur.Currency == "AUD" {
			continue
		}
		balance := cur.Balance / 100000000
		value := balance * cur.CurrentPricing.LastPrice
		total += value
		valueStr := fmt.Sprintf("%.2f", value)
		lastPriceStr := fmt.Sprintf("%.2f", cur.CurrentPricing.LastPrice)
		fmt.Printf("%s\t%f\t%s\t%s\n", cur.Currency, balance, lastPriceStr, valueStr)
	}

	fmt.Printf("Total: %s\n", fmt.Sprintf("%.2f", total))

	return nil
}

func (h handle) getAccountBalance() ([]*CurrencyBalance, error) {
	data, err := h.btcclient.MakeRequest("/account/balance")
	if err != nil {
		return []*CurrencyBalance{}, err
	}

	currencies := make([]*CurrencyBalance, 0)
	if err := json.Unmarshal(data, &currencies); err != nil {
		fmt.Println("=== There is an error")
		return []*CurrencyBalance{}, err
	}

	return currencies, nil
}

func (h handle) getCurrentPrices(currency string) (MarketData, error) {
	data, err := h.btcclient.MakeRequest("/market/" + currency + "/AUD/tick")
	if err != nil {
		return MarketData{}, err
	}

	currencyData := MarketData{}
	if err := json.Unmarshal(data, &currencyData); err != nil {
		fmt.Println("=== There is an error getting currency data")
		return currencyData, err
	}

	return currencyData, nil
}
