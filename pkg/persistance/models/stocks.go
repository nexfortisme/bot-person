package persistance

type Stock struct {
	StockTicker string  `json:"stockTicker"`
	StockCount  float64 `json:"stockCount"`
}