package repository

import (
	"context"
	"errors"
	v1 "github.com/go-nunu/nunu-layout-advanced/api/v1"
	"github.com/go-nunu/nunu-layout-advanced/internal/model"
	"github.com/morkid/paginate"
	"gorm.io/gorm"
	"net/http"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *model.Transaction) error
	GetByID(ctx context.Context, id string) (*model.Transaction, error)
	GetPaginated(ctx context.Context, query *v1.GetTransactionsQuery, pg *paginate.Pagination, request *http.Request) *paginate.Page
}

func NewTransactionRepository(
	r *Repository,
) TransactionRepository {
	return &transactionRepository{
		Repository: r,
	}
}

type transactionRepository struct {
	*Repository
}

// Create creates a new transaction
func (r *transactionRepository) Create(ctx context.Context, transaction *model.Transaction) error {
	if err := r.DB(ctx).Create(transaction).Error; err != nil {
		return err
	}
	return nil
}

// GetByID retrieves a transaction by its ID
func (r *transactionRepository) GetByID(ctx context.Context, transactionId string) (*model.Transaction, error) {
	var transaction model.Transaction
	if err := r.DB(ctx).Preload("User").Where("transaction_id = ?", transactionId).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, v1.ErrNotFound
		}
		return nil, err
	}
	return &transaction, nil
}

// GetPaginated retrieves transactions with pagination and filtering
func (r *transactionRepository) GetPaginated(ctx context.Context, query *v1.GetTransactionsQuery, pg *paginate.Pagination, request *http.Request) *paginate.Page {
	db := r.DB(ctx).Model(&model.Transaction{}).Preload("User")
	
	// Apply filters
	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}
	if query.Search != "" {
		db = db.Where("description LIKE ?", "%"+query.Search+"%")
	}
	
	// Use paginate package as shown in example
	page := pg.With(db).Request(request).Response(&[]model.Transaction{})
	
	return &page
}