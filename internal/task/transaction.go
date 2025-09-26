package task

import (
	"context"
	"github.com/go-nunu/nunu-layout-advanced/internal/model"
	"github.com/go-nunu/nunu-layout-advanced/internal/repository"
	"go.uber.org/zap"
)

type TransactionTask interface {
	CreateSampleTransactions(ctx context.Context) error
}

func NewTransactionTask(
	task *Task,
	transactionRepo repository.TransactionRepository,
	userRepo repository.UserRepository,
) TransactionTask {
	return &transactionTask{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
		Task:            task,
	}
}

type transactionTask struct {
	transactionRepo repository.TransactionRepository
	userRepo        repository.UserRepository
	*Task
}

// CreateSampleTransactions creates sample transaction data for testing
func (t *transactionTask) CreateSampleTransactions(ctx context.Context) error {
	t.logger.Info("Creating sample transactions...")

	// Create sample transactions
	sampleTransactions := []*model.Transaction{
		{
			TransactionId: "tx001",
			UserId:        "user001",
			Amount:        100.50,
			Type:          "credit",
			Status:        "completed",
			Description:   "Sample credit transaction",
		},
		{
			TransactionId: "tx002",
			UserId:        "user001",
			Amount:        50.25,
			Type:          "debit",
			Status:        "completed",
			Description:   "Sample debit transaction",
		},
		{
			TransactionId: "tx003",
			UserId:        "user002",
			Amount:        200.00,
			Type:          "credit",
			Status:        "pending",
			Description:   "Pending credit transaction",
		},
		{
			TransactionId: "tx004",
			UserId:        "user002",
			Amount:        75.30,
			Type:          "debit",
			Status:        "failed",
			Description:   "Failed debit transaction",
		},
		{
			TransactionId: "tx005",
			UserId:        "user001",
			Amount:        300.00,
			Type:          "credit",
			Status:        "completed",
			Description:   "Large credit transaction",
		},
	}

	for _, transaction := range sampleTransactions {
		if err := t.transactionRepo.Create(ctx, transaction); err != nil {
			t.logger.Error("Failed to create transaction", zap.Error(err), zap.String("transactionId", transaction.TransactionId))
			return err
		}
		t.logger.Info("Created transaction", zap.String("transactionId", transaction.TransactionId))
	}

	t.logger.Info("Sample transactions created successfully")
	return nil
}