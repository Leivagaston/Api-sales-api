package sale

import "errors"

// ErrNotFound is returned when a sale with the given ID is not found.
var ErrNotFound = errors.New("sale not found")

// ErrEmptyID is returned when trying to store a sale with an empty ID.
var ErrEmptyID = errors.New("empty sale ID")

// Storage is the main interface for our storage layer.
type Storage interface {
	Set(sale *Sale) error
	Read(id string) (*Sale, error)
	ReadAllByUserID(id string) []Sale
	Delete(id string) error
}

// LocalStorage provides an in-memory implementation for storing sales.
type LocalStorage struct {
	m map[string]*Sale
}

// NewLocalStorage instantiates a new LocalStorage with an empty map.
func NewLocalStorage() *LocalStorage {
	return &LocalStorage{
		m: map[string]*Sale{},
	}
}

// Set stores or updates a sale in the local storage.
// Returns ErrEmptyID if the sale has an empty ID.
func (l *LocalStorage) Set(sale *Sale) error {
	if sale.Id == "" {
		return ErrEmptyID
	}

	l.m[sale.Id] = sale
	return nil
}

// Read retrieves a sale from the local storage by ID.
// Returns ErrNotFound if the sale is not found.
func (l *LocalStorage) Read(id string) (*Sale, error) {
	u, ok := l.m[id]
	if !ok {
		return nil, ErrNotFound
	}

	return u, nil
}

// Read retrieves all sales from the local storage by user ID.
// Returns ErrNotFound if the sale is not found.
func (l *LocalStorage) ReadAllByUserID(id string) []Sale {
	var result []Sale
	for _, sale := range l.m {
		if sale != nil && sale.UserID == id {
			result = append(result, *sale)
		}
	}
	return result
}

// Delete removes a sale from the local storage by ID.
// Returns ErrNotFound if the sale does not exist.
func (l *LocalStorage) Delete(id string) error {
	_, err := l.Read(id)
	if err != nil {
		return err
	}

	delete(l.m, id)
	return nil
}
