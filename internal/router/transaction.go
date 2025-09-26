package router

import (
	"github.com/gin-gonic/gin"
)

func InitTransactionRouter(
	deps RouterDeps,
	r *gin.RouterGroup,
) {
	// Public routes (no authentication required)
	publicRouter := r.Group("/")
	{
		publicRouter.GET("/transaction/:id", deps.TransactionHandler.GetTransaction)
		publicRouter.GET("/transactions", deps.TransactionHandler.GetTransactions)
	}
}