# 🌿 Aloé Raw — Go + React E-Commerce

**Frontend:** React 18 + Vite → GitHub Pages  
**Backend:** Go + Gin → Render  
**Database:** SQLite | **Auth:** Google OAuth | **Payments:** Razorpay

---

## Architecture

```
startupinspirator.github.io/aloe-raw   →   React Frontend
aloe-raw-api.onrender.com              →   Go Backend (Gin)
```

---

## Local Development

### Prerequisites
- [Go 1.21+](https://go.dev/dl/)
- [Node.js 18+](https://nodejs.org/)
- [GCC](https://gcc.gnu.org/) (for SQLite — on Windows install [TDM-GCC](https://jmeubank.github.io/tdm-gcc/))

### 1. Clone the repo
```bash
git clone https://github.com/startupinspirator/aloe-raw.git
cd aloe-raw
```

### 2. Set up the Go backend
```bash
cd backend
cp .env.example .env      # Fill in your keys (instructions inside)
go mod tidy               # Download dependencies
go run .                  # Starts on http://localhost:8080
```

### 3. Set up the React frontend (new terminal)
```bash
cd frontend
npm install
npm run dev               # Starts on http://localhost:5173
```

Open **http://localhost:5173** — Vite proxies all `/api` and `/auth` calls to Go.

---

## Get Your API Keys

### Google OAuth
1. Go to [console.cloud.google.com](https://console.cloud.google.com)
2. Create project → APIs & Services → Credentials
3. Create OAuth 2.0 Client ID → Web application
4. Add redirect URIs:
   - `http://localhost:8080/auth/google/callback`
   - `https://aloe-raw-api.onrender.com/auth/google/callback`
5. Copy Client ID + Secret → `backend/.env`

### Razorpay
1. Go to [dashboard.razorpay.com](https://dashboard.razorpay.com)
2. Settings → API Keys → Generate Test Key
3. Copy Key ID + Secret → `backend/.env`

---

## Deploy

### Backend → Render (free)
1. Push this repo to GitHub
2. Go to [render.com](https://render.com) → New → Web Service
3. Connect your GitHub repo → select the **backend/** folder (or root)
4. Render auto-detects Go. Use these settings:
   - **Build:** `go build -o server .`
   - **Start:** `./server`
5. Add all env vars from `backend/.env.example` under Environment
6. Deploy — your API will be at `https://aloe-raw-api.onrender.com`

### Frontend → GitHub Pages (automatic)
1. In your GitHub repo → Settings → Secrets → Actions:
   - Add secret: `VITE_API_URL` = `https://aloe-raw-api.onrender.com`
2. Go to Settings → Pages → Source: `gh-pages` branch
3. Push any change to `frontend/` → GitHub Actions builds + deploys automatically
4. Site live at: **https://startupinspirator.github.io/aloe-raw**

---

## Project Structure
```
aloe-raw/
├── .github/workflows/deploy.yml   # Auto-deploy frontend on push
├── backend/
│   ├── main.go                    # Gin server entry point
│   ├── go.mod                     # Go dependencies
│   ├── render.yaml                # One-click Render deploy config
│   ├── .env.example               # Environment variables template
│   ├── database/db.go             # SQLite init, tables, seed
│   ├── models/models.go           # Go structs
│   ├── middleware/auth.go         # Session auth guard
│   └── handlers/
│       ├── auth.go                # Google OAuth (golang.org/x/oauth2)
│       ├── products.go            # GET /api/products
│       ├── cart.go                # Cart CRUD
│       ├── orders.go              # Order history
│       └── payment.go             # Razorpay create + verify
└── frontend/
    ├── src/
    │   ├── assets/images.js       # All 3 product photos (base64)
    │   ├── config.js              # API URL switching dev/prod
    │   ├── context/
    │   │   ├── AuthContext.jsx
    │   │   └── CartContext.jsx
    │   ├── components/Navbar.jsx
    │   ├── pages/
    │   │   ├── Home.jsx
    │   │   ├── Product.jsx
    │   │   ├── Cart.jsx
    │   │   ├── Checkout.jsx
    │   │   └── Orders.jsx
    │   ├── App.jsx
    │   ├── main.jsx
    │   └── index.css
    ├── vite.config.js             # base: "/aloe-raw/" for GitHub Pages
    ├── .env.example
    └── package.json               # includes gh-pages deploy script
```
