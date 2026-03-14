package handlers

import (
	"log"
	"net/http"
	"strconv"
	
	"github.com/gin-gonic/gin"
	"github.com/startupinspirator/aloe-raw/backend/database"
	"github.com/startupinspirator/aloe-raw/backend/models"
)

func GetReviews(c *gin.Context) {
	productID := c.Param("id")
	rows, err := database.DB.Query(`
		SELECT r.id, r.user_id, r.product_id, r.rating, r.comment, r.created_at, u.name, u.avatar
		FROM reviews r
		JOIN users u ON r.user_id = u.id
		WHERE r.product_id = ?
		ORDER BY r.created_at DESC`, productID)

	if err != nil {
		log.Println("Error fetching reviews:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch reviews"})
		return
	}
	defer rows.Close()

	var reviews []models.Review
	for rows.Next() {
		var r models.Review
		if err := rows.Scan(&r.ID, &r.UserID, &r.ProductID, &r.Rating, &r.Comment, &r.CreatedAt, &r.UserName, &r.UserAvatar); err != nil {
			continue
		}
		reviews = append(reviews, r)
	}

	c.JSON(http.StatusOK, gin.H{"reviews": reviews})
}

func SubmitReview(c *gin.Context) {
	userID := c.MustGet("user_id").(int)
	productIDStr := c.Param("id")
	
	productID, err := strconv.Atoi(productIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	var body struct {
		Rating  int    `json:"rating" binding:"required,min=1,max=5"`
		Comment string `json:"comment"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already reviewed this product
	var count int
	database.DB.QueryRow("SELECT COUNT(*) FROM reviews WHERE user_id = ? AND product_id = ?", userID, productID).Scan(&count)
	if count > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You have already reviewed this product"})
		return
	}

	_, err = database.DB.Exec(`
		INSERT INTO reviews (user_id, product_id, rating, comment)
		VALUES (?, ?, ?, ?)`,
		userID, productID, body.Rating, body.Comment,
	)

	if err != nil {
		log.Println("Error submitting review:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to submit review"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Review submitted successfully"})
}
