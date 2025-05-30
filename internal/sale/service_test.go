package sale

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestService_Create_Simple(t *testing.T) {
	s := NewService(NewLocalStorage(), nil)

	input := &Sale{}

	fields := &CreateFields{
		UserID: strPtr("user-123"),
		Amount: floatPtr(5500.00),
	}

	err := s.Create(input, fields)

	require.Nil(t, err)
	require.NotEmpty(t, input.Id)
	require.NotEmpty(t, input.CreatedAt)
	require.NotEmpty(t, input.UpdatedAt)
	require.Equal(t, 1, input.Version)

	s = NewService(&mockStorage{
		mockSet: func(sale *Sale) error {
			return errors.New("fake error trying to set sale")
		},
	}, nil)

	err = s.Create(input, fields)
	require.NotNil(t, err)
	require.EqualError(t, err, "fake error trying to set sale")
}

func TestService_Create(t *testing.T) {
	type fields struct {
		storage Storage
	}

	type args struct {
		sale   *Sale
		fields *CreateFields
	}

	tests := []struct {
		name     string
		fields   fields
		args     args
		wantErr  func(t *testing.T, err error)
		wantSale func(t *testing.T, sale *Sale)
	}{
		{
			name: "error",
			fields: fields{
				storage: &mockStorage{
					mockSet: func(sale *Sale) error {
						return errors.New("fake error trying to set sale")
					},
				},
			},
			args: args{
				sale: &Sale{},
				fields: &CreateFields{
					UserID: strPtr("user-123"),
					Amount: floatPtr(1234.56),
				},
			},
			wantErr: func(t *testing.T, err error) {
				require.NotNil(t, err)
				require.EqualError(t, err, "fake error trying to set sale")
			},
			wantSale: nil,
		},
		{
			name: "success",
			fields: fields{
				storage: NewLocalStorage(),
			},
			args: args{
				sale: &Sale{
					UserID: "user-123",
					Amount: 5500.00,
				},
				fields: &CreateFields{
					UserID: strPtr("user-123"),
					Amount: floatPtr(5500.00),
				},
			},
			wantErr: func(t *testing.T, err error) {
				require.Nil(t, err)
			},
			wantSale: func(t *testing.T, input *Sale) {
				require.NotEmpty(t, input.Id)
				require.NotEmpty(t, input.CreatedAt)
				require.NotEmpty(t, input.UpdatedAt)
				require.Equal(t, 1, input.Version)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				storage: tt.fields.storage,
			}

			err := s.Create(tt.args.sale, tt.args.fields)
			if tt.wantErr != nil {
				tt.wantErr(t, err)
			}

			if tt.wantSale != nil {
				tt.wantSale(t, tt.args.sale)
			}
		})
	}
}

type mockStorage struct {
	mockSet             func(sale *Sale) error
	mockRead            func(id string) (*Sale, error)
	mockReadAllByUserID func(id string) []Sale
	mockDelete          func(id string) error
}

func (m *mockStorage) Set(sale *Sale) error {
	return m.mockSet(sale)
}

func (m *mockStorage) Read(id string) (*Sale, error) {
	if m.mockRead != nil {
		return m.mockRead(id)
	}
	return nil, errors.New("mockRead not implemented")
}

func (m *mockStorage) Delete(id string) error {
	if m.mockDelete != nil {
		return m.mockDelete(id)
	}
	return errors.New("mockDelete not implemented")
}

func (m *mockStorage) ReadAllByUserID(id string) []Sale {
	if m.mockReadAllByUserID != nil {
		return m.mockReadAllByUserID(id)
	}
	return nil
}

func strPtr(s string) *string {
	return &s
}

func floatPtr(f float64) *float64 {
	return &f
}
