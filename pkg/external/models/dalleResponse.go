package external

type DalleResponse struct {
	Created int           `json:"created"`
	Data    []DalleImages `json:"data"`
}