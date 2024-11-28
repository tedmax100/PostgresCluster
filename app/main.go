package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

func main() {
	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}

func setupRouter() *gin.Engine {
	primaryDsn := "host=localhost user=user password=password dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Taipei application_name=app"
	replicaDsn := "host=localhost user=user password=password dbname=postgres port=5433 sslmode=disable TimeZone=Asia/Taipei application_name=app"

	db, err := gorm.Open(postgres.Open(primaryDsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	err = db.Use(
		dbresolver.Register(dbresolver.Config{
			Sources:           []gorm.Dialector{postgres.Open(primaryDsn)},
			Replicas:          []gorm.Dialector{postgres.Open(replicaDsn)},
			Policy:            dbresolver.RandomPolicy{},
			TraceResolverMode: true,
		}).
			SetMaxIdleConns(2).
			SetMaxOpenConns(2).
			SetConnMaxIdleTime(10 * time.Minute).
			SetConnMaxLifetime(1 * time.Hour),
	)

	if err != nil {
		panic(fmt.Sprintf("failed to configure dbresolver: %v", err))
	}

	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	// Ping test
	r.GET("/products", func(c *gin.Context) {
		var products []Product
		if err := db.Preload("Categories").Find(&products).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var formattedProducts []map[string]interface{}
		for _, product := range products {
			categoryNames := []string{}
			for _, category := range product.Categories {
				categoryNames = append(categoryNames, category.Name)
			}

			formattedProducts = append(formattedProducts, map[string]interface{}{
				"id":          product.ID,
				"name":        product.Name,
				"description": product.Description,
				"picture":     product.Picture,
				"priceUsd": map[string]interface{}{
					"currencyCode": product.PriceCurrencyCode,
					"units":        product.PriceUnits,
					"nanos":        product.PriceNanos,
				},
				"categories": categoryNames,
			})
		}

		c.JSON(http.StatusOK, formattedProducts)
	})

	r.GET("/products/:id", func(c *gin.Context) {
		productID := c.Param("id")

		var product Product

		result := db.Preload("Categories").First(&product, "id = ?", productID)

		if result.Error != nil {
			if result.Error == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		categoryNames := []string{}
		for _, category := range product.Categories {
			categoryNames = append(categoryNames, category.Name)
		}

		formattedProduct := map[string]interface{}{
			"id":          product.ID,
			"name":        product.Name,
			"description": product.Description,
			"picture":     product.Picture,
			"priceUsd": map[string]interface{}{
				"currencyCode": product.PriceCurrencyCode,
				"units":        product.PriceUnits,
				"nanos":        product.PriceNanos,
			},
			"categories": categoryNames,
		}

		c.JSON(http.StatusOK, formattedProduct)
	})

	r.GET("/products/search", func(c *gin.Context) {
		query := c.Query("query")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter is required"})
			return
		}

		var products []Product

		result := db.Preload("Categories").Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", "%"+strings.ToLower(query)+"%", "%"+strings.ToLower(query)+"%").Find(&products)

		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": result.Error.Error()})
			return
		}

		var formattedProducts []map[string]interface{}
		for _, product := range products {
			categoryNames := []string{}
			for _, category := range product.Categories {
				categoryNames = append(categoryNames, category.Name)
			}

			formattedProducts = append(formattedProducts, map[string]interface{}{
				"id":          product.ID,
				"name":        product.Name,
				"description": product.Description,
				"picture":     product.Picture,
				"priceUsd": map[string]interface{}{
					"currencyCode": product.PriceCurrencyCode,
					"units":        product.PriceUnits,
					"nanos":        product.PriceNanos,
				},
				"categories": categoryNames,
			})
		}

		c.JSON(http.StatusOK, formattedProducts)
	})

	return r
}
