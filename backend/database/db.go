package database

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	_ "modernc.org/sqlite"
)

var DB *sql.DB

func Init() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./aloe_raw.db"
	}

	// Ensure directory exists if DB_PATH is in a subdirectory (e.g., /app/data/aloe_raw.db)
	if dir := filepath.Dir(dbPath); dir != "." {
		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Printf("Warning: Failed to create DB directory %s: %v", dir, err)
		}
	}

	var err error
	DB, err = sql.Open("sqlite", dbPath)
	if err != nil { log.Fatal("Failed to open database:", err) }
	if err = DB.Ping(); err != nil { log.Fatal("Failed to ping database:", err) }
	createTables()
	seedProduct()
	log.Println("Database ready")
}

func createTables() {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			google_id TEXT UNIQUE NOT NULL,
			name TEXT NOT NULL,
			email TEXT UNIQUE NOT NULL,
			avatar TEXT,
			role TEXT DEFAULT 'customer',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS categories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			slug TEXT UNIQUE NOT NULL,
			description TEXT
		)`,
		`CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			category_id INTEGER,
			name TEXT NOT NULL,
			tagline TEXT,
			description TEXT,
			price INTEGER NOT NULL,
			original_price INTEGER,
			stock INTEGER DEFAULT 100,
			active INTEGER DEFAULT 1,
			FOREIGN KEY(category_id) REFERENCES categories(id)
		)`,
		`CREATE TABLE IF NOT EXISTS cart (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			quantity INTEGER DEFAULT 1,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(product_id) REFERENCES products(id),
			UNIQUE(user_id, product_id)
		)`,
		`CREATE TABLE IF NOT EXISTS orders (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			razorpay_order_id TEXT,
			razorpay_payment_id TEXT,
			total_amount INTEGER NOT NULL,
			status TEXT DEFAULT 'pending',
			shipping_name TEXT,
			shipping_address TEXT,
			shipping_city TEXT,
			shipping_pincode TEXT,
			shipping_phone TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS order_items (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			order_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			quantity INTEGER NOT NULL,
			price_at_purchase INTEGER NOT NULL,
			FOREIGN KEY(order_id) REFERENCES orders(id),
			FOREIGN KEY(product_id) REFERENCES products(id)
		)`,
		`CREATE TABLE IF NOT EXISTS reviews (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			product_id INTEGER NOT NULL,
			rating INTEGER NOT NULL CHECK(rating >= 1 AND rating <= 5),
			comment TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id),
			FOREIGN KEY(product_id) REFERENCES products(id)
		)`,
	}
	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil { log.Fatal("Table creation failed:", err) }
	}
}

func seedProduct() {
	// Seed Categories
	var catCount int
	DB.QueryRow("SELECT COUNT(*) FROM categories").Scan(&catCount)
	if catCount == 0 {
		DB.Exec(`INSERT INTO categories (name,slug,description) VALUES (?,?,?)`, "Health & Wellness", "health-wellness", "Pure, organic health products for your wellbeing.")
		DB.Exec(`INSERT INTO categories (name,slug,description) VALUES (?,?,?)`, "Skincare", "skincare", "Natural, restorative skincare routines.")
		DB.Exec(`INSERT INTO categories (name,slug,description) VALUES (?,?,?)`, "Supplements", "supplements", "Daily vitamins and herbal supplements.")
		log.Println("Categories seeded")
	}

	// Seed Products
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if count == 0 {
		// Category 1: Health & Wellness
		DB.Exec(`INSERT INTO products (category_id,name,tagline,description,price,original_price,stock) VALUES (?,?,?,?,?,?,?)`,
			1, "Aloé Raw Cubes", "Farm Fresh, 100% Pure. Zero Processing.", "10 individually sealed cubes of freshly harvested aloe vera gel. Cut from mature leaves on our certified organic farm and shipped within 24 hours of harvest. No preservatives, no fillers — just pure living aloe. Suitable for skin, hair, and internal use. Refrigerate up to 7 days, freeze up to 3 months.", 299, 450, 100)
		DB.Exec(`INSERT INTO products (category_id,name,tagline,description,price,original_price,stock) VALUES (?,?,?,?,?,?,?)`,
			1, "Aloé Detox Juice", "Start Your Day Fresh.", "A 500ml bottle of pure aloe vera detox juice. Infused with a hint of lemon and ginger, it's the perfect morning cleanse for your digestive system.", 399, 500, 50)
		
		// Category 2: Skincare
		DB.Exec(`INSERT INTO products (category_id,name,tagline,description,price,original_price,stock) VALUES (?,?,?,?,?,?,?)`,
			2, "Aloé Glow Gel", "Hydration in a jar.", "A lightweight, fast-absorbing pure aloe vera gel designed to hydrate, soothe, and calm irritated skin. Ideal for after-sun care or daily moisturization.", 450, 600, 75)
		DB.Exec(`INSERT INTO products (category_id,name,tagline,description,price,original_price,stock) VALUES (?,?,?,?,?,?,?)`,
			2, "Aloé Face Wash", "Gentle, Natural Cleansing.", "Formulated with 80% raw aloe vera and chamomile extract. Washes away impurities without stripping your skin's natural oils.", 350, 450, 120)

		// Category 3: Supplements
		DB.Exec(`INSERT INTO products (category_id,name,tagline,description,price,original_price,stock) VALUES (?,?,?,?,?,?,?)`,
			3, "Aloé Immunity Boost capsules", "Your daily defense.", "60 capsules of freeze-dried organic aloe powder. Rich in acemannan to support healthy immune function and digestion.", 899, 1200, 30)

		log.Println("Products seeded")
	}

	// Seed Mock Users
	var userCount int
	DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&userCount)
	if userCount == 0 {
		DB.Exec(`INSERT INTO users (google_id, name, email, avatar, role) VALUES ('mock_admin', 'Admin User', 'admin@aloeraw.com', 'https://ui-avatars.com/api/?name=Admin+User', 'admin')`)
		DB.Exec(`INSERT INTO users (google_id, name, email, avatar, role) VALUES ('mock_user_1', 'Alice Green', 'alice@example.com', 'https://ui-avatars.com/api/?name=Alice+Green', 'customer')`)
		DB.Exec(`INSERT INTO users (google_id, name, email, avatar, role) VALUES ('mock_user_2', 'Bob Smith', 'bob@example.com', 'https://ui-avatars.com/api/?name=Bob+Smith', 'customer')`)
		log.Println("Users seeded")
	}

	// Seed Mock Reviews
	var reviewCount int
	DB.QueryRow("SELECT COUNT(*) FROM reviews").Scan(&reviewCount)
	if reviewCount == 0 {
		var u1, u2 int
		DB.QueryRow("SELECT id FROM users WHERE google_id = 'mock_user_1'").Scan(&u1)
		DB.QueryRow("SELECT id FROM users WHERE google_id = 'mock_user_2'").Scan(&u2)
		
		// Reviews for Product 1
		DB.Exec(`INSERT INTO reviews (user_id, product_id, rating, comment) VALUES (?, ?, ?, ?)`, u1, 1, 5, "Absolutely love these aloe cubes! So fresh.")
		DB.Exec(`INSERT INTO reviews (user_id, product_id, rating, comment) VALUES (?, ?, ?, ?)`, u2, 1, 4, "Great product, fast shipping. A bit pricey but worth it.")
		// Reviews for Product 3
		DB.Exec(`INSERT INTO reviews (user_id, product_id, rating, comment) VALUES (?, ?, ?, ?)`, u1, 3, 5, "My skin has never felt better! Will buy again.")
		
		log.Println("Reviews seeded")
	}

	// Seed Orders
	var orderCount int
	DB.QueryRow("SELECT COUNT(*) FROM orders").Scan(&orderCount)
	if orderCount == 0 {
		var u1 int
		DB.QueryRow("SELECT id FROM users WHERE google_id = 'mock_user_1'").Scan(&u1)

		// Create a mock order
		res, _ := DB.Exec(`INSERT INTO orders (user_id, total_amount, status, shipping_name, shipping_address, shipping_city, shipping_pincode, shipping_phone) VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			u1, 749, "delivered", "Alice Green", "123 Mock Street", "Dream City", "100001", "9876543210")
		orderID, _ := res.LastInsertId()

		// Add order items
		if orderID > 0 {
			DB.Exec(`INSERT INTO order_items (order_id, product_id, quantity, price_at_purchase) VALUES (?, ?, ?, ?)`, orderID, 1, 1, 299)
			DB.Exec(`INSERT INTO order_items (order_id, product_id, quantity, price_at_purchase) VALUES (?, ?, ?, ?)`, orderID, 3, 1, 450)
		}

		log.Println("Orders seeded")
	}
}
