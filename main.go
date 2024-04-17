package main

import (
	"fmt"
	"net/http"
	"io"
	"encoding/json"
	"strconv"
	"sync"
	"time"
	"math"
)

func floatConvert(dataToConvert interface{}) float64 {

	closePriceStr, ok := dataToConvert.(string)
		if !ok {
			return 0.0
		}

	closePrice, err := strconv.ParseFloat(closePriceStr, 64)
	if err != nil {
		return 0.0
	}

	return closePrice
}

func fetchCoinbasePrices() ([]float64, error) {
	url := "https://api.exchange.coinbase.com/products/SOL-USD/ticker"

	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	var parsedData map[string]interface{}
	parsingErr := json.Unmarshal([]byte(body), &parsedData)
	if parsingErr != nil {
		fmt.Println("Error unmarshalling JSON:", parsingErr)
		return nil, err
	}

	var allPrices []float64
	allPrices = append(allPrices, floatConvert(parsedData["price"]))
	allPrices = append(allPrices, floatConvert(parsedData["ask"]))
	allPrices = append(allPrices, floatConvert(parsedData["bid"]))

	return allPrices, nil
}

func fetchKrakenPrices() ([]float64, error) {
	url := "https://api.kraken.com/0/public/Ticker?pair=SOLUSD"

	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
	}

	var parsedData map[string]interface{}
	parsingErr := json.Unmarshal([]byte(body), &parsedData)
	if parsingErr != nil {
		fmt.Println("Error unmarshalling JSON:", parsingErr)
		return nil, err
	}

	result, ok := parsedData["result"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("error asserting type for 'result'")
	}

	solUsd, ok := result["SOLUSD"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("error asserting type for 'SOLUSD'")
	}

	c, ok := solUsd["c"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("error asserting type for 'c'")
	}

	a, ok := solUsd["a"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("error asserting type for 'a'")
	}

	b, ok := solUsd["b"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("error asserting type for 'b'")
	}

	var allPrices []float64

	allPrices = append(allPrices, floatConvert(c[0]))
	allPrices = append(allPrices, floatConvert(a[0]))
	allPrices = append(allPrices, floatConvert(b[0]))

	return allPrices, nil
}

func getAllPrices() ([]float64, []float64) {
	var wg sync.WaitGroup
	var coinbasePrices []float64
	var krakenPrices []float64
	var err1, err2 error

	wg.Add(2)

	go func() {
		defer wg.Done()
		coinbasePrices, err1 = fetchCoinbasePrices()
	}()

	go func() {
		defer wg.Done()
		krakenPrices, err2 = fetchKrakenPrices()
	}()

	if err1 != nil {
		fmt.Println("Error fetching from Coinbase:", err1)
		return nil, nil
	}
	if err2 != nil {
		fmt.Println("Error fetching from Kraken:", err2)
		return nil, nil
	}

	wg.Wait()
	return coinbasePrices, krakenPrices
}

func isArbitrageOpportunity(price_a, price_b float64) bool {
	threshold := 1.0
	lastPriceDifference := math.Abs(((price_a - price_b) / price_b) * 100)
	fmt.Println(lastPriceDifference)
	if lastPriceDifference >= threshold {
		return true
	} else {
		return false
	}
	
}
 
func main() {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			//Each contains in order: last trade price, bid price, ask price
			coinbasePrices, krakenPrices := getAllPrices()
		
			fmt.Println("Arbitrage opportunity?: ", isArbitrageOpportunity(coinbasePrices[0], krakenPrices[0]))
		}
	}
}