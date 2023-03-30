package external

import (
	"context"
	"errors"
	"fmt"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
)

func GetStockData(ticker string, apiKey string) (float32, error) {

	// I don't like reinitializing this every time its called but it will stay for now
	// TODO - Find a better way to do this
	cfg := finnhub.NewConfiguration()
	cfg.AddDefaultHeader("X-Finnhub-Token", apiKey)
	finnhubClient := finnhub.NewAPIClient(cfg).DefaultApi

	// Quote
	quote, _, _ := finnhubClient.Quote(context.Background()).Symbol(ticker).Execute()
	fmt.Printf("Price: %f\n", *quote.C)

	if *quote.C == 0 {
		return 0, errors.New("no data found")
	} else {
		return *quote.C, nil
	}
}
