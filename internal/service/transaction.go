package service

import (
	"context"
	v1 "github.com/go-nunu/nunu-layout-advanced/api/v1"
	"github.com/go-nunu/nunu-layout-advanced/internal/model"
	"github.com/go-nunu/nunu-layout-advanced/internal/repository"
	"github.com/morkid/paginate"
	"net/http"
	"time"
)

type TransactionService interface {
	CreateTransaction(ctx context.Context, req *v1.CreateTransactionRequest) (*v1.CreateTransactionResponse, error)
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

// CreateTransaction creates a new transaction
func (s *transactionService) CreateTransaction(ctx context.Context, req *v1.CreateTransactionRequest) (*v1.CreateTransactionResponse, error) {
	// Generate transaction ID
	transactionId, err := s.sid.GenString()
	if err != nil {
		return nil, v1.ErrInternalServerError
	}

	// Create transaction model
	transaction := &model.Transaction{
		TransactionId: transactionId,
		UserId:        req.UserId,
		Amount:        req.Amount,
		Type:          req.Type,
		Status:        req.Status,
		Description:   req.Description,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	// Save to database
	if err := s.transactionRepo.Create(ctx, transaction); err != nil {
		return nil, v1.ErrInternalServerError
	}

	return &v1.CreateTransactionResponse{
		TransactionId: transaction.TransactionId,
		UserId:        transaction.UserId,
		Amount:        transaction.Amount,
		Type:          transaction.Type,
		Status:        transaction.Status,
		Description:   transaction.Description,
		CreatedAt:     transaction.CreatedAt,
	}, nil
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