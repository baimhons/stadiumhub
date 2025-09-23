package utils

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseRepository[T any] interface {
	GetAll(items *[]T, pagination *PaginationQuery) error
	GetBy(field string, value string, item *T) error
	GetByID(item *T, id uuid.UUID) error
	Create(item *T) error
	Update(item *T) error
	Delete(item *T) error
	Begin() *gorm.DB
}

type baseRepository[T any] struct {
	DB *gorm.DB
}

func NewBaseRepository[T any](db *gorm.DB) BaseRepository[T] {
	return &baseRepository[T]{DB: db}
}

func (r *baseRepository[T]) GetAll(items *[]T, pagination *PaginationQuery) error {
	query := r.DB

	if pagination.Page != nil && pagination.PageSize != nil {
		offset := *pagination.Page * *pagination.PageSize
		query = query.Offset(offset).Limit(*pagination.PageSize)
	}

	if pagination.Sort != nil && pagination.Order != nil {
		orderClause := *pagination.Sort + " " + *pagination.Order
		query = query.Order(orderClause)
	}

	if pagination.Type != nil && *pagination.Type != "" {
		query = query.Where("type = ?", *pagination.Type)
	}

	return query.Find(items).Error
}

func (r *baseRepository[T]) GetBy(field string, value string, item *T) error {
	return r.DB.Where(field+" = ?", value).First(item).Error
}

func (r *baseRepository[T]) GetByID(item *T, id uuid.UUID) error {
	return r.DB.Where("id = ?", id).First(item).Error
}

func (r *baseRepository[T]) Create(item *T) error {
	return r.DB.Create(item).Error
}

func (r *baseRepository[T]) Update(item *T) error {
	return r.DB.Updates(item).Error
}

func (r *baseRepository[T]) Delete(item *T) error {
	return r.DB.Delete(item).Error
}

func (r *baseRepository[T]) Begin() *gorm.DB {
	return r.DB.Begin()
}
