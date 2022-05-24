package emaStrategy

type EMAStrategy struct {
	Id         string  `json:"id" bson:"_id"`
	EntryTime  float64 `json:"entryTime" bson:"entryTime"`
	Type       string  `json:"type" bson:"type"` // Entry Long | Exit Long | Entry Short | Exit Short
	Quantity   float64 `json:"quantity" bson:"quantity"`
	EntryPrice float64 `json:"entryPrice" bson:"entryPrice"`
	ExitTime   float64 `json:"exitTime" bson:"exitTime"`
	ExitPrice  float64 `json:"exitPrice" bson:"exitPrice"`
	Profit     float64 `json:"profit" bson:"profit"`
	CumProfit  float64 `json:"cumProfit" bson:"cumProfit"`
}

type EMAStrategyRepository interface {
	Create(EMAStrategy) (*EMAStrategy, error)
	GetAll() ([]EMAStrategy, error)
	GetById(string) (*EMAStrategy, error)
	UpdateById(string, EMAStrategy) (*EMAStrategy, error)
	DeleteById(string) (int, error)
}
