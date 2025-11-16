package common

import (
	"github.com/LordPax/godular/core"
	"gorm.io/gorm"
)

// ModelMock is a generic struct that implements the IModel interface for any type T.
type ModelMock[T any] struct {
	*core.Provider
	*core.Mockable
}

// NewModelMock creates a new instance of ModelMock for the specified type T.
func NewModelMock[T any](name string) *ModelMock[T] {
	return &ModelMock[T]{
		Provider: core.NewProvider(name),
		Mockable: core.NewMockable(),
	}
}

// Migrate performs automatic migration for the model's schema.
func (m *ModelMock[T]) Migrate() error {
	m.MethodCalled("Migrate")
	m.CallFunc("Migrate")
	return nil
}

// GetDB returns the current gorm.DB instance.
func (m *ModelMock[T]) GetDB() *gorm.DB {
	return nil
}

// SetDB sets the gorm.DB instance for the model.
func (m *ModelMock[T]) SetDB(db *gorm.DB) {}

// FindByID retrieves a record by its ID.
func (m *ModelMock[T]) FindByID(id string) (T, error) {
	m.MethodCalled("FindByID", id)
	return m.CallFunc("FindByID").(T), nil
}

// FindOneBy retrieves a record by a specified field and value.
func (m *ModelMock[T]) FindOneBy(field string, value any) (T, error) {
	m.MethodCalled("FindOneBy", field, value)
	return m.CallFunc("FindOneBy").(T), nil
}

// Create inserts a new record into the database.
func (m *ModelMock[T]) Create(entity T) error {
	m.MethodCalled("Create", entity)
	m.CallFunc("Create")
	return nil
}

// DeleteByID deletes a record by its ID.
func (m *ModelMock[T]) DeleteByID(id string) error {
	m.MethodCalled("DeleteByID", id)
	m.CallFunc("DeleteByID")
	return nil
}

// DeleteBy deletes records matching a specified field and value.
func (m *ModelMock[T]) DeleteBy(field string, value any) error {
	m.MethodCalled("DeleteBy", field, value)
	m.CallFunc("DeleteBy")
	return nil
}

// UpdateByID updates a record by its ID with the provided updates.
func (m *ModelMock[T]) UpdateByID(id string, updates T) error {
	m.MethodCalled("UpdateByID", id, updates)
	m.CallFunc("UpdateByID")
	return nil
}

// CountBy counts the number of records matching a specified field and value.
func (m *ModelMock[T]) CountBy(field string, value any) (int64, error) {
	m.MethodCalled("CountBy", field, value)
	return m.CallFunc("CountBy").(int64), nil
}

func (m *ModelMock[T]) ClearTable() error {
	m.MethodCalled("ClearTable")
	m.CallFunc("ClearTable")
	return nil
}
