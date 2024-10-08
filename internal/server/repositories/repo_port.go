package repositories

import (
	"gorm.io/gorm"

	"arczed/internal/entities/models"
)

type ConstRepository interface {
}

type ProductRepository interface {
	SaveProduct() error
	FindProduct(req *models.AddProduct) error
}

func NewConstRepository(db *gorm.DB) ConstRepository {
	return &beseDB{
		baseRequest: &baseRequest{db: db},
	}
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	return &beseDB{
		baseRequest: &baseRequest{db: db},
	}
}
