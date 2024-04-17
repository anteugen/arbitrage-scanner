func fetchClosePriceBinance() (float64, error) {
	url := "https://fapi.binance.com/fapi/v1/klines?symbol=SOLUSDT&interval=1m&limit=1"
	fmt.Println(url)

	response, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
		return 0.0, err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return 0.0, err
	}

	var parsedData [][]interface{}
	if err := json.Unmarshal(body, &parsedData); err != nil {
		return 0.0, fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	return floatConvert(parsedData[0][4]), nil
}