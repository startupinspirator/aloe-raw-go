package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/startupinspirator/aloe-raw/backend/database"
	"github.com/startupinspirator/aloe-raw/backend/models"
)

func GetCategories(c *gin.Context) {
	rows, err := database.DB.Query("SELECT id, name, slug, description FROM categories")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch categories"})
		return
	}
	defer rows.Close()

	var categories []models.Category
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.Slug, &cat.Description); err != nil {
			continue
		}
		categories = append(categories, cat)
	}

	c.JSON(http.StatusOK, gin.H{"categories": categories})
}
