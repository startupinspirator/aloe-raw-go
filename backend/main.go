package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/startupinspirator/aloe-raw/backend/database"
	"github.com/startupinspirator/aloe-raw/backend/handlers"
	"github.com/startupinspirator/aloe-raw/backend/middleware"
)

func main() {
	// Init DB
	database.Init()

	r := gin.Default()

	// ── Session store ──────────────────────────────────────
	secret := os.Getenv("SESSION_SECRET")
	if secret == "" {
		secret = "aloe-raw-secret-change-in-prod"
	}
	store := cookie.NewStore([]byte(secret))
	store.Options(sessions.Options{
		Path:     "/",
		MaxAge:   86400 * 7, // 7 days
		HttpOnly: true,
		Secure:   os.Getenv("GIN_MODE") == "release",
		SameSite: 4, // SameSiteNoneMode for cross-origin
	})
	r.Use(sessions.Sessions("aloe_session", store))

	// ── CORS ───────────────────────────────────────────────
	clientURL := os.Getenv("CLIENT_URL")
	if clientURL == "" {
		clientURL = "http://localhost:5173"
	}
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{clientURL},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept"},
		AllowCredentials: true,
	}))

	// ── Auth routes ────────────────────────────────────────
	auth := r.Group("/auth")
	{
		auth.GET("/google", handlers.GoogleLogin)
		auth.GET("/google/callback", handlers.GoogleCallback)
		auth.GET("/me", handlers.GetMe)
		auth.POST("/logout", handlers.Logout)
	}

	// ── API routes ─────────────────────────────────────────
	api := r.Group("/api")
	{
		api.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{"status": "ok"})
		})

		api.GET("/products", handlers.GetProducts)
		api.GET("/products/:id", handlers.GetProduct)

		cart := api.Group("/cart", middleware.RequireAuth)
		{
			cart.GET("", handlers.GetCart)
			cart.POST("", handlers.AddToCart)
			cart.DELETE("/:product_id", handlers.RemoveFromCart)
		}

		orders := api.Group("/orders", middleware.RequireAuth)
		{
			orders.GET("", handlers.GetOrders)
		}

		payment := api.Group("/payment", middleware.RequireAuth)
		{
			payment.POST("/create-order", handlers.CreateOrder)
			payment.POST("/verify", handlers.VerifyPayment)
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("🌿 Aloé Raw Go server running on :%s", port)
	r.Run(":" + port)
}
