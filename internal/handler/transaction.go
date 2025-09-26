package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/go-nunu/nunu-layout-advanced/api/v1"
	"github.com/go-nunu/nunu-layout-advanced/internal/service"
	"go.uber.org/zap"
	"net/http"
)

type TransactionHandler struct {
	*Handler
	transactionService service.TransactionService
}

func NewTransactionHandler(handler *Handler, transactionService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		Handler:            handler,
		transactionService: transactionService,
	}
}

// GetTransaction godoc
// @Summary 获取单个交易
// @Schemes
// @Description 根据交易ID获取交易详情
// @Tags 交易模块
// @Accept json
// @Produce json
// @Param id path string true "Transaction ID"
// @Success 200 {object} v1.GetTransactionResponseData
// @Router /transaction/{id} [get]
func (h *TransactionHandler) GetTransaction(ctx *gin.Context) {
	transactionId := ctx.Param("id")
	if transactionId == "" {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	transaction, err := h.transactionService.GetTransaction(ctx, transactionId)
	if err != nil {
		h.logger.WithContext(ctx).Error("transactionService.GetTransaction error", zap.Error(err))
		if err == v1.ErrNotFound {
			v1.HandleError(ctx, http.StatusNotFound, err, nil)
		} else {
			v1.HandleError(ctx, http.StatusInternalServerError, v1.ErrInternalServerError, nil)
		}
		return
	}

	v1.HandleSuccess(ctx, transaction)
}

// GetTransactions godoc
// @Summary 获取交易列表
// @Schemes
// @Description 分页获取交易列表，支持筛选和搜索
// @Tags 交易模块
// @Accept json
// @Produce json
// @Param page query int false "页码" default(1)
// @Param size query int false "每页数量" default(10)
// @Param sort query string false "排序字段" default("created_at desc")
// @Param type query string false "交易类型"
// @Param status query string false "交易状态"
// @Param search query string false "搜索关键词"
// @Success 200 {object} map[string]interface{}
// @Router /transactions [get]
func (h *TransactionHandler) GetTransactions(ctx *gin.Context) {
	var query v1.GetTransactionsQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		v1.HandleError(ctx, http.StatusBadRequest, v1.ErrBadRequest, nil)
		return
	}

	page := h.transactionService.GetTransactions(ctx, &query, ctx.Request)

	v1.HandleSuccess(ctx, page)
}