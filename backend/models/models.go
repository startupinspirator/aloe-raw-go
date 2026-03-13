package models

type User struct {
	ID        int    `json:"id"`
	GoogleID  string `json:"google_id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	Avatar    string `json:"avatar"`
	CreatedAt string `json:"created_at"`
}

type Product struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	Tagline       string `json:"tagline"`
	Description   string `json:"description"`
	Price         int    `json:"price"`
	OriginalPrice int    `json:"original_price"`
	Stock         int    `json:"stock"`
	Active        int    `json:"active"`
}

type CartItem struct {
	ID        int     `json:"id"`
	UserID    int     `json:"user_id"`
	ProductID int     `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Name      string  `json:"name"`
	Price     int     `json:"price"`
	Tagline   string  `json:"tagline"`
}

type Order struct {
	ID                int         `json:"id"`
	UserID            int         `json:"user_id"`
	RazorpayOrderID   string      `json:"razorpay_order_id"`
	RazorpayPaymentID string      `json:"razorpay_payment_id"`
	TotalAmount       int         `json:"total_amount"`
	Status            string      `json:"status"`
	ShippingName      string      `json:"shipping_name"`
	ShippingAddress   string      `json:"shipping_address"`
	ShippingCity      string      `json:"shipping_city"`
	ShippingPincode   string      `json:"shipping_pincode"`
	ShippingPhone     string      `json:"shipping_phone"`
	CreatedAt         string      `json:"created_at"`
	Items             []OrderItem `json:"items"`
}

type OrderItem struct {
	ID              int    `json:"id"`
	OrderID         int    `json:"order_id"`
	ProductID       int    `json:"product_id"`
	Quantity        int    `json:"quantity"`
	PriceAtPurchase int    `json:"price_at_purchase"`
	Name            string `json:"name"`
}

type ShippingDetails struct {
	Name    string `json:"name"    binding:"required"`
	Phone   string `json:"phone"   binding:"required"`
	Address string `json:"address" binding:"required"`
	City    string `json:"city"    binding:"required"`
	Pincode string `json:"pincode" binding:"required"`
}
