package model

type ResponseBody struct {
	Id            string  `json:"id"`
	PhotoURi      *string `json:"photouri"`
	Latitude      float64 `json:"latitude"`
	Longitude     float64 `json:"longitude"`
	DisplayPrice  int     `json:"displayPrice"`
	Status        string  `json:"status"`
	Bedrooms      int     `json:"badrooms"`
	FullBathrooms int     `json:"fullBathrooms"`
	HalfBathrooms int     `json:"halfBathrooms"`
	SquareFeet    *int    `json:"squareFeet"`
	Address       string  `json:"address"`
	Unit          *string `json:"unit"`
	City          string  `json:"city"`
	State         string  `json:"state"`
	Zip           int     `json:"zip"`
}
