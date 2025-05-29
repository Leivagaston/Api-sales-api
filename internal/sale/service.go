package sale

import (
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// Service provides high-level sale management operations on a LocalStorage backend.
type Service struct {
	// storage is the underlying persistence for User entities.
	storage Storage

	// logger is our observability component to log.
	logger *zap.Logger
}

// NewService creates a new Service.
func NewService(storage Storage, logger *zap.Logger) *Service {
	if logger == nil {
		logger, _ = zap.NewProduction()
		defer logger.Sync() // flushes buffer, if any
	}

	return &Service{
		storage: storage,
		logger:  logger,
	}
}

// Create adds a brand-new sale to the system.
// It sets CreatedAt and UpdatedAt to the current time and initializes Version to 1.
// Returns ErrEmptyID if sale.ID is empty.
func (s *Service) Create(sale *Sale, newSale *CreateFields) error {
	if newSale == nil || newSale.UserID == nil || newSale.Amount == nil {
		return errors.New("missing required fields: UserID or Amount")
	}

	statuses := []string{"pending", "completed", "cancelled"}

	sale.UserID = *newSale.UserID
	sale.Amount = *newSale.Amount
	sale.Id = uuid.NewString()
	now := time.Now()
	sale.CreatedAt = now
	sale.UpdatedAt = now

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	sale.Status = statuses[r.Intn(len(statuses))]
	sale.Version = 1

	if err := s.storage.Set(sale); err != nil {
		if s.logger != nil {
			s.logger.Error("failed to set sale", zap.Error(err), zap.Any("sale", sale))
		}
		return err
	}

	return nil
}

// Get retrieves a sale by its ID.
// Returns ErrNotFound if no sale exists with the given ID.
func (s *Service) Get(user_id string, status string) (*Sale, error) {
	return s.storage.Read(user_id)

	//una vez que recibo todas las ventas asociadas al user_id hago un filtro por status
}

// Update modifies an existing sale's data.
// It updates Name, Address, NickName, sets UpdatedAt to now and increments Version.
// Returns ErrNotFound if the sale does not exist, or ErrEmptyID if sale.ID is empty.
func (s *Service) Update(id string, sale *UpdateFields) (*Sale, error) {
	existing, err := s.storage.Read(id)
	if err != nil {
		return nil, err
	}

	if sale.Status != nil && strings.EqualFold(*sale.Status, "pending") {
		existing.Status = *sale.Status
	}

	existing.UpdatedAt = time.Now()
	existing.Version++

	if err := s.storage.Set(existing); err != nil {
		return nil, err
	}

	return existing, nil
}

// Delete removes a sale from the system by its ID.
// Returns ErrNotFound if the sale does not exist.
func (s *Service) Delete(id string) error {
	return s.storage.Delete(id)
}
