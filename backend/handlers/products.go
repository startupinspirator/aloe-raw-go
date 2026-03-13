package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/startupinspirator/aloe-raw/backend/database"
	"github.com/startupinspirator/aloe-raw/backend/models"
)

func GetProducts(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id,name,tagline,description,price,original_price,stock FROM products WHERE active=1")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}
	defer rows.Close()
	var products []models.Product
	for rows.Next() {
		var p models.Product
		rows.Scan(&p.ID, &p.Name, &p.Tagline, &p.Description, &p.Price, &p.OriginalPrice, &p.Stock)
		products = append(products, p)
	}
	c.JSON(http.StatusOK, products)
}

func GetProduct(c *gin.Context) {
	id := c.Param("id")
	var p models.Product
	err := database.DB.QueryRow(
		"SELECT id,name,tagline,description,price,original_price,stock FROM products WHERE id=? AND active=1", id,
	).Scan(&p.ID, &p.Name, &p.Tagline, &p.Description, &p.Price, &p.OriginalPrice, &p.Stock)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
		return
	}
	c.JSON(http.StatusOK, p)
}
