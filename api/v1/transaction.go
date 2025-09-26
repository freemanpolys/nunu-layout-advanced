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
	Page     int    `form:"page,default=1" example:"1"`
	Size     int    `form:"size,default=10" example:"10"`
	Sort     string `form:"sort,default=created_at desc" example:"created_at desc"`
	Type     string `form:"type" example:"credit"`
	Status   string `form:"status" example:"completed"`
	Search   string `form:"search" example:"description search"`
}