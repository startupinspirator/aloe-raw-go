package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/startupinspirator/aloe-raw/backend/database"
	"github.com/startupinspirator/aloe-raw/backend/models"
)

func GetCart(c *gin.Context) {
	userID := c.MustGet("user_id")
	rows, err := database.DB.Query(`
		SELECT cart.id, cart.quantity, products.id, products.name, products.price, products.tagline
		FROM cart JOIN products ON cart.product_id = products.id
		WHERE cart.user_id = ?`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}
	defer rows.Close()
	var items []models.CartItem
	for rows.Next() {
		var item models.CartItem
		rows.Scan(&item.ID, &item.Quantity, &item.ProductID, &item.Name, &item.Price, &item.Tagline)
		items = append(items, item)
	}
	if items == nil { items = []models.CartItem{} }
	c.JSON(http.StatusOK, items)
}

func AddToCart(c *gin.Context) {
	userID := c.MustGet("user_id")
	var body struct {
		ProductID int `json:"product_id" binding:"required"`
		Quantity  int `json:"quantity"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "product_id required"})
		return
	}
	if body.Quantity < 1 { body.Quantity = 1 }
	database.DB.Exec(`
		INSERT INTO cart (user_id, product_id, quantity) VALUES (?,?,?)
		ON CONFLICT(user_id, product_id) DO UPDATE SET quantity=excluded.quantity`,
		userID, body.ProductID, body.Quantity)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func RemoveFromCart(c *gin.Context) {
	userID := c.MustGet("user_id")
	productID := c.Param("product_id")
	database.DB.Exec("DELETE FROM cart WHERE user_id=? AND product_id=?", userID, productID)
	c.JSON(http.StatusOK, gin.H{"success": true})
}
