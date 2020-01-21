package exchanges

// Order contains order information
type GeminiOrder struct {
	OrderID           int64    `json:"order_id,string"`
	ID                int64    `json:"id,string"`
	ClientOrderID     string   `json:"client_order_id"`
	Symbol            string   `json:"symbol"`
	Exchange          string   `json:"exchange"`
	Price             float64  `json:"price,string"`
	AvgExecutionPrice float64  `json:"avg_execution_price,string"`
	Side              string   `json:"side"`
	Type              string   `json:"type"`
	Timestamp         int64    `json:"timestamp,string"`
	TimestampMS       int64    `json:"timestampms"`
	IsLive            bool     `json:"is_live"`
	IsCancelled       bool     `json:"is_cancelled"`
	IsHidden          bool     `json:"is_hidden"`
	Options           []string `json:"options"`
	WasForced         bool     `json:"was_forced"`
	ExecutedAmount    float64  `json:"executed_amount,string"`
	RemainingAmount   float64  `json:"remaining_amount,string"`
	OriginalAmount    float64  `json:"original_amount,string"`
	Message           string   `json:"message"`
}

type GeminiError struct {
	Result  string `json:"result"`
	Reason  string `json:"reason"`
	Message string `json:"message"`
}

type GeminiBalance struct {
	Currency               string `json:"currency"`
	Amount                 string `json:"amount"`
	Available              string `json:"available"`
	AvailableForWithdrawal string `json:"availableForWithdrawal"`
	Type                   string `json:"type"`
}
