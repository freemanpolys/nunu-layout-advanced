package service

import (
	"context"
	v1 "github.com/go-nunu/nunu-layout-advanced/api/v1"
	"github.com/go-nunu/nunu-layout-advanced/internal/repository"
	"github.com/morkid/paginate"
	"net/http"
)

type TransactionService interface {
	GetTransaction(ctx context.Context, transactionId string) (*v1.GetTransactionResponse, error)
	GetTransactions(ctx context.Context, query *v1.GetTransactionsQuery, request *http.Request) *paginate.Page
}

func NewTransactionService(
	service *Service,
	transactionRepo repository.TransactionRepository,
) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		Service:         service,
	}
}

type transactionService struct {
	transactionRepo repository.TransactionRepository
	*Service
}

// GetTransaction retrieves a single transaction by ID
func (s *transactionService) GetTransaction(ctx context.Context, transactionId string) (*v1.GetTransactionResponse, error) {
	transaction, err := s.transactionRepo.GetByID(ctx, transactionId)
	if err != nil {
		return nil, err
	}

	return &v1.GetTransactionResponse{
		TransactionId: transaction.TransactionId,
		UserId:        transaction.UserId,
		Amount:        transaction.Amount,
		Type:          transaction.Type,
		Status:        transaction.Status,
		Description:   transaction.Description,
		CreatedAt:     transaction.CreatedAt,
		UpdatedAt:     transaction.UpdatedAt,
	}, nil
}

// GetTransactions retrieves transactions with pagination and filtering
func (s *transactionService) GetTransactions(ctx context.Context, query *v1.GetTransactionsQuery, request *http.Request) *paginate.Page {
	pg := paginate.New()
	page := s.transactionRepo.GetPaginated(ctx, query, pg, request)
	return page
}