package repositories

import (
	"fmt"
	"math/rand"
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID       int
	Name     string
	Quantity int
}

type ProductRepository interface {
	FindAll(useCache bool) ([]Product, error)
}

type productRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) ProductRepository {
	db.AutoMigrate(&Product{})
	mock(db)
	return productRepository{db}
}

func (r productRepository) FindAll(useCache bool) (products []Product, err error) {

	err = r.db.Order("quantity desc").Limit(30).Find(&products).Error
	return
}

func mock(db *gorm.DB) error {

	var count int64
	db.Model(&Product{}).Count(&count)
	if count > 0 {
		// if data existed, return
		return nil
	}

	seed := rand.NewSource(time.Now().UnixNano())
	random := rand.New(seed)

	products := []Product{}
	for i := 0; i < 5000; i++ {
		products = append(products, Product{
			Name:     fmt.Sprintf("Product %v", i),
			Quantity: random.Intn(5000),
		})
	}

	return db.Create(&products).Error
}
