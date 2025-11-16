package common

import (
	"fmt"

	"github.com/LordPax/godular/core"
	"gorm.io/gorm"
)

// IModel defines the interface for a generic model with common database operations.
type IModel[T any] interface {
	FindByID(id string) (T, error)
	FindOneBy(field string, value any) (T, error)
	Create(entity T) error
	DeleteByID(id string) error
	DeleteBy(field string, value any) error
	UpdateByID(id string, updates T) error
	CountBy(field string, value any) (int64, error)
	GetDB() *gorm.DB
	SetDB(db *gorm.DB)
	Migrate() error
	ClearTable() error
}

// Model is a generic struct that implements the IModel interface for any type T.
type Model[T any] struct {
	*core.Provider
	db *gorm.DB
}

// NewModel creates a new instance of Model for the specified type T.
func NewModel[T any](name string) *Model[T] {
	return &Model[T]{
		Provider: core.NewProvider(name),
	}
}

// Migrate performs automatic migration for the model's schema.
func (m *Model[T]) Migrate() error {
	fmt.Printf("Migrating %s\n", m.GetName())
	return m.db.AutoMigrate(new(T))
}

// GetDB returns the current gorm.DB instance.
func (m *Model[T]) GetDB() *gorm.DB {
	return m.db
}

// SetDB sets the gorm.DB instance for the model.
func (m *Model[T]) SetDB(db *gorm.DB) {
	m.db = db
}

// FindByID retrieves a record by its ID.
func (m *Model[T]) FindByID(id string) (T, error) {
	var item T
	err := m.db.Model(new(T)).First(&item, "id = ?", id).Error
	return item, err
}

// FindOneBy retrieves a record by a specified field and value.
func (m *Model[T]) FindOneBy(field string, value any) (T, error) {
	var item T
	err := m.db.Model(new(T)).Where(field, value).First(&item).Error
	return item, err
}

// Create inserts a new record into the database.
func (m *Model[T]) Create(entity T) error {
	return m.db.Model(new(T)).Create(entity).Error
}

// DeleteByID deletes a record by its ID.
func (m *Model[T]) DeleteByID(id string) error {
	return m.db.Model(new(T)).Delete(new(T), "id = ?", id).Error
}

// DeleteBy deletes records matching a specified field and value.
func (m *Model[T]) DeleteBy(field string, value any) error {
	return m.db.Model(new(T)).Where(field, value).Delete(new(T)).Error
}

// UpdateByID updates a record by its ID with the provided updates.
func (m *Model[T]) UpdateByID(id string, updates T) error {
	return m.db.Model(new(T)).Where("id = ?", id).Updates(updates).Error
}

// CountBy counts the number of records matching a specified field and value.
func (m *Model[T]) CountBy(field string, value any) (int64, error) {
	var count int64 = 0
	err := m.db.Model(new(T)).Where(field, value).Count(&count).Error
	return count, err
}

func (m *Model[T]) ClearTable() error {
	return m.db.
		Session(&gorm.Session{AllowGlobalUpdate: true}).
		Delete(new(T)).Error
}
