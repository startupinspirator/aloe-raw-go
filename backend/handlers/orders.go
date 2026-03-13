package handlers

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/startupinspirator/aloe-raw/backend/database"
	"github.com/startupinspirator/aloe-raw/backend/models"
)

func GetOrders(c *gin.Context) {
	userID := c.MustGet("user_id")
	rows, err := database.DB.Query(`
		SELECT id,razorpay_order_id,razorpay_payment_id,total_amount,status,
		       shipping_name,shipping_address,shipping_city,shipping_pincode,shipping_phone,created_at
		FROM orders WHERE user_id=? ORDER BY created_at DESC`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		rows.Scan(&o.ID, &o.RazorpayOrderID, &o.RazorpayPaymentID, &o.TotalAmount,
			&o.Status, &o.ShippingName, &o.ShippingAddress, &o.ShippingCity,
			&o.ShippingPincode, &o.ShippingPhone, &o.CreatedAt)

		// Get items
		itemRows, _ := database.DB.Query(`
			SELECT oi.id, oi.product_id, oi.quantity, oi.price_at_purchase, p.name
			FROM order_items oi JOIN products p ON oi.product_id=p.id
			WHERE oi.order_id=?`, o.ID)
		for itemRows.Next() {
			var item models.OrderItem
			itemRows.Scan(&item.ID, &item.ProductID, &item.Quantity, &item.PriceAtPurchase, &item.Name)
			o.Items = append(o.Items, item)
		}
		itemRows.Close()
		if o.Items == nil { o.Items = []models.OrderItem{} }
		orders = append(orders, o)
	}
	if orders == nil { orders = []models.Order{} }
	c.JSON(http.StatusOK, orders)
}
