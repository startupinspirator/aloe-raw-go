package routes

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/startupinspirator/aloe-raw/backend/handlers"
	"github.com/startupinspirator/aloe-raw/backend/middleware"
)

func SetupRouter() *gin.Engine {
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
	allowedOrigins := []string{"http://localhost:5173", "https://aloe-raw-go.pages.dev"}
	if clientURL != "" {
		allowedOrigins = append(allowedOrigins, clientURL)
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     allowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
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

		api.GET("/categories", handlers.GetCategories)
		api.GET("/products", handlers.GetProducts)
		api.GET("/products/:id", handlers.GetProduct)
		api.GET("/products/:id/reviews", handlers.GetReviews)

		// Authenticated Routes
		api.POST("/products/:id/reviews", middleware.RequireAuth, handlers.SubmitReview)

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

		admin := api.Group("/admin", middleware.RequireAuth, middleware.RequireAdmin)
		{
			admin.GET("/orders", handlers.GetAdminOrders)
		}
	}

	return r
}
