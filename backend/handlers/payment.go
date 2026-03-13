package handlers

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/razorpay/razorpay-go"
	"github.com/startupinspirator/aloe-raw/backend/database"
	"github.com/startupinspirator/aloe-raw/backend/models"
)

func CreateOrder(c *gin.Context) {
	userID := c.MustGet("user_id")

	var body struct {
		Shipping models.ShippingDetails `json:"shipping" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Shipping details required"})
		return
	}

	// Get cart
	rows, err := database.DB.Query(`
		SELECT cart.quantity, products.id, products.price, products.name
		FROM cart JOIN products ON cart.product_id=products.id
		WHERE cart.user_id=?`, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "DB error"})
		return
	}
	defer rows.Close()

	type cartRow struct{ qty, productID, price int; name string }
	var cartItems []cartRow
	totalAmount := 0
	for rows.Next() {
		var r cartRow
		rows.Scan(&r.qty, &r.productID, &r.price, &r.name)
		totalAmount += r.price * r.qty
		cartItems = append(cartItems, r)
	}
	if len(cartItems) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	// Create order in DB
	s := body.Shipping
	result, err := database.DB.Exec(`
		INSERT INTO orders (user_id,total_amount,status,shipping_name,shipping_address,shipping_city,shipping_pincode,shipping_phone)
		VALUES (?,?,?,?,?,?,?,?)`,
		userID, totalAmount, "pending", s.Name, s.Address, s.City, s.Pincode, s.Phone)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create order"})
		return
	}
	orderID, _ := result.LastInsertId()

	for _, item := range cartItems {
		database.DB.Exec(`INSERT INTO order_items (order_id,product_id,quantity,price_at_purchase) VALUES (?,?,?,?)`,
			orderID, item.productID, item.qty, item.price)
	}

	// Create Razorpay order
	client := razorpay.NewClient(os.Getenv("RAZORPAY_KEY_ID"), os.Getenv("RAZORPAY_KEY_SECRET"))
	data := map[string]interface{}{
		"amount":   totalAmount * 100,
		"currency": "INR",
		"receipt":  fmt.Sprintf("order_%d", orderID),
	}
	rzpOrder, err := client.Order.Create(data, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Razorpay error"})
		return
	}

	rzpOrderID := fmt.Sprintf("%v", rzpOrder["id"])
	database.DB.Exec("UPDATE orders SET razorpay_order_id=? WHERE id=?", rzpOrderID, orderID)

	c.JSON(http.StatusOK, gin.H{
		"orderId":         orderID,
		"razorpayOrderId": rzpOrderID,
		"amount":          totalAmount,
		"keyId":           os.Getenv("RAZORPAY_KEY_ID"),
	})
}

func VerifyPayment(c *gin.Context) {
	userID := c.MustGet("user_id")
	var body struct {
		RazorpayOrderID   string `json:"razorpay_order_id"`
		RazorpayPaymentID string `json:"razorpay_payment_id"`
		RazorpaySignature string `json:"razorpay_signature"`
		OrderID           int    `json:"orderId"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid body"})
		return
	}

	// Verify signature
	secret := os.Getenv("RAZORPAY_KEY_SECRET")
	message := body.RazorpayOrderID + "|" + body.RazorpayPaymentID
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(message))
	expected := hex.EncodeToString(mac.Sum(nil))

	if expected != body.RazorpaySignature {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Payment verification failed"})
		return
	}

	database.DB.Exec("UPDATE orders SET status='paid', razorpay_payment_id=? WHERE id=?",
		body.RazorpayPaymentID, body.OrderID)
	database.DB.Exec("DELETE FROM cart WHERE user_id=?", userID)

	c.JSON(http.StatusOK, gin.H{"success": true, "orderId": body.OrderID})
}
