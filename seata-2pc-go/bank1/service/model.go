package service

type AccountInfo struct {
	Id              int64   `json:"id"`
	AccountName     string  `json:"account_name"`
	AccountNo       string  `json:"account_no"`
	AccountPassword string  `json:"account_password"`
	AccountBalance  float64 `json:"account_balance"`
}

func (m *AccountInfo) TableName() string {
	return "account_info"
}

type AccountEvent struct {
	AccountNo int64   `json:"account_no"`
	Amount    float64 `json:"amount"`
}
