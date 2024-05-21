package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/redis/go-redis/v9"
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
	rd *redis.Client
}

func NewProductRepository(db *gorm.DB, rd *redis.Client) ProductRepository {
	db.AutoMigrate(&Product{})
	mock(db)
	return productRepository{db, rd}
}

func (r productRepository) FindAll(useCache bool) (products []Product, err error) {

	key := "products"

	// Check and try getting cache in Redis
	if useCache {
		productJson, err := r.rd.Get(context.Background(), key).Result()
		if err == nil {
			err := json.Unmarshal([]byte(productJson), &products)
			if err == nil {
				println("cache")
				return products, nil
			}
		}
	}

	// Query from database
	err = r.db.Order("quantity desc").Limit(30).Find(&products).Error
	println("database")

	// Set to Redis
	if useCache {
		data, err := json.Marshal(products)
		if err == nil {
			r.rd.Set(context.Background(), key, string(data), time.Second*10)
		}
	}

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
