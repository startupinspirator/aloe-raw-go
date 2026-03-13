package database

import (
	"database/sql"
	"log"
	"os"
	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func Init() {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./aloe_raw.db"
	}
	var err error
	DB, err = sql.Open("sqlite3", dbPath+"?_journal_mode=WAL")
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
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS products (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			tagline TEXT,
			description TEXT,
			price INTEGER NOT NULL,
			original_price INTEGER,
			stock INTEGER DEFAULT 100,
			active INTEGER DEFAULT 1
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
	}
	for _, q := range queries {
		if _, err := DB.Exec(q); err != nil { log.Fatal("Table creation failed:", err) }
	}
}

func seedProduct() {
	var count int
	DB.QueryRow("SELECT COUNT(*) FROM products").Scan(&count)
	if count > 0 { return }
	DB.Exec(`INSERT INTO products (name,tagline,description,price,original_price,stock) VALUES (?,?,?,?,?,?)`,
		"Aloé Raw Cubes",
		"Farm Fresh, 100% Pure. Zero Processing.",
		"10 individually sealed cubes of freshly harvested aloe vera gel. Cut from mature leaves on our certified organic farm and shipped within 24 hours of harvest. No preservatives, no fillers — just pure living aloe. Suitable for skin, hair, and internal use. Refrigerate up to 7 days, freeze up to 3 months.",
		299, 450, 100,
	)
	log.Println("Product seeded")
}
