package models

type Transaction struct {
    ID          int     `json:"id"`
    Amount      float64 `json:"amount"`
    Type        string  `json:"type"` // "income" or "expense"
    Description string  `json:"description"`
}
