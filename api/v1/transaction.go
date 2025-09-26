package v1

import "time"

type GetTransactionResponse struct {
	TransactionId string    `json:"transactionId"`
	UserId        string    `json:"userId"`
	Amount        float64   `json:"amount"`
	Type          string    `json:"type"`
	Status        string    `json:"status"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}

type GetTransactionResponseData struct {
	Response
	Data GetTransactionResponse `json:"data"`
}

type GetTransactionsQuery struct {
	Page    int    `form:"page,default=1" example:"1"`
	Size    int    `form:"size,default=10" example:"10"`
	Sort    string `form:"sort,default=created_at desc" example:"created_at desc"`
	Filters string `form:"filters" example:"[[\"type\",\"=\",\"credit\"]]"`
}

type CreateTransactionRequest struct {
	UserId      string  `json:"userId" binding:"required" example:"user123"`
	Amount      float64 `json:"amount" binding:"required" example:"100.50"`
	Type        string  `json:"type" binding:"required" example:"credit"`
	Status      string  `json:"status" binding:"required" example:"pending"`
	Description string  `json:"description" example:"Payment for order #123"`
}

type CreateTransactionResponse struct {
	TransactionId string    `json:"transactionId"`
	UserId        string    `json:"userId"`
	Amount        float64   `json:"amount"`
	Type          string    `json:"type"`
	Status        string    `json:"status"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"createdAt"`
}

type CreateTransactionResponseData struct {
	Response
	Data CreateTransactionResponse `json:"data"`
}