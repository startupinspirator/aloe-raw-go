package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/startupinspirator/aloe-raw/backend/database"
	"github.com/startupinspirator/aloe-raw/backend/models"
)

func GetAdminOrders(c *gin.Context) {
	rows, err := database.DB.Query(`
		SELECT id, user_id, razorpay_order_id, razorpay_payment_id, total_amount, status,
		       shipping_name, shipping_address, shipping_city, shipping_pincode, shipping_phone, created_at
		FROM orders ORDER BY created_at DESC`)
	if err != nil {
		log.Println("Error fetching all orders:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var o models.Order
		var rpO, rpP, shipN, shipA, shipC, shipPin, shipPh *string
		rows.Scan(&o.ID, &o.UserID, &rpO, &rpP, &o.TotalAmount, &o.Status,
			&shipN, &shipA, &shipC, &shipPin, &shipPh, &o.CreatedAt)

		if rpO != nil { o.RazorpayOrderID = *rpO }
		if rpP != nil { o.RazorpayPaymentID = *rpP }
		if shipN != nil { o.ShippingName = *shipN }
		if shipA != nil { o.ShippingAddress = *shipA }
		if shipC != nil { o.ShippingCity = *shipC }
		if shipPin != nil { o.ShippingPincode = *shipPin }
		if shipPh != nil { o.ShippingPhone = *shipPh }

		// Get Items for this order
		iRows, _ := database.DB.Query(`
			SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price_at_purchase, p.name 
			FROM order_items oi JOIN products p ON oi.product_id=p.id 
			WHERE oi.order_id=?`, o.ID)
		for iRows.Next() {
			var i models.OrderItem
			iRows.Scan(&i.ID, &i.OrderID, &i.ProductID, &i.Quantity, &i.PriceAtPurchase, &i.Name)
			o.Items = append(o.Items, i)
		}
		iRows.Close()

		orders = append(orders, o)
	}

	c.JSON(http.StatusOK, gin.H{"orders": orders})
}
