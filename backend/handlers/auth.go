package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/startupinspirator/aloe-raw/backend/database"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func getOAuthConfig() *oauth2.Config {
	clientURL := os.Getenv("CLIENT_URL")
	if clientURL == "" {
		clientURL = "http://localhost:5173"
	}
	callbackURL := os.Getenv("GOOGLE_CALLBACK_URL")
	if callbackURL == "" {
		callbackURL = "http://localhost:8080/auth/google/callback"
	}
	return &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  callbackURL,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}
}

func GoogleLogin(c *gin.Context) {
	// MOCK BYPASS: If no Google Client ID is provided, redirect directly to callback with a mock code
	if os.Getenv("GOOGLE_CLIENT_ID") == "" || os.Getenv("GOOGLE_CLIENT_ID") == "your_google_client_id_here" {
		callbackURL := os.Getenv("GOOGLE_CALLBACK_URL")
		if callbackURL == "" {
			callbackURL = "http://localhost:8080/auth/google/callback"
		}
		c.Redirect(http.StatusTemporaryRedirect, callbackURL+"?code=MOCK_CODE")
		return
	}

	cfg := getOAuthConfig()
	url := cfg.AuthCodeURL("aloe-raw-state", oauth2.AccessTypeOnline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context) {
	clientURL := os.Getenv("CLIENT_URL")
	if clientURL == "" {
		clientURL = "http://localhost:5173"
	}

	code := c.Query("code")
	if code == "" {
		c.Redirect(http.StatusTemporaryRedirect, clientURL+"/?login=failed")
		return
	}

	var info struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Email   string `json:"email"`
		Picture string `json:"picture"`
	}

	// MOCK BYPASS
	if code == "MOCK_CODE" {
		info.ID = "mock_google_id_123"
		info.Name = "Test User"
		info.Email = "test@example.com"
		info.Picture = "https://ui-avatars.com/api/?name=Test+User&background=random"
	} else {
		cfg := getOAuthConfig()
		token, err := cfg.Exchange(context.Background(), code)
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, clientURL+"/?login=failed")
			return
		}

		client := cfg.Client(context.Background(), token)
		resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, clientURL+"/?login=failed")
			return
		}
		defer resp.Body.Close()

		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &info)
	}

	// Upsert user
	db := database.DB
	if code == "MOCK_CODE" {
		db.Exec(`
			INSERT INTO users (google_id, name, email, avatar, role)
			VALUES (?, ?, ?, ?, 'admin')
			ON CONFLICT(google_id) DO UPDATE SET name=excluded.name, avatar=excluded.avatar, role='admin'`,
			info.ID, info.Name, info.Email, info.Picture,
		)
	} else {
		db.Exec(`
			INSERT INTO users (google_id, name, email, avatar)
			VALUES (?, ?, ?, ?)
			ON CONFLICT(google_id) DO UPDATE SET name=excluded.name, avatar=excluded.avatar`,
			info.ID, info.Name, info.Email, info.Picture,
		)
	}

	var user struct {
		ID     int
		Name   string
		Email  string
		Avatar string
		Role   string
	}
	err := db.QueryRow("SELECT id, name, email, avatar, role FROM users WHERE google_id = ?", info.ID).
		Scan(&user.ID, &user.Name, &user.Email, &user.Avatar, &user.Role)
	if err != nil {
		fmt.Println("Error fetching mock user:", err)
		c.Redirect(http.StatusTemporaryRedirect, clientURL+"/?login=failed")
		return
	}

	// Save session
	session := sessions.Default(c)
	session.Set("user_id", user.ID)
	session.Set("user_name", user.Name)
	session.Set("user_email", user.Email)
	session.Set("user_avatar", user.Avatar)
	session.Set("user_role", user.Role)
	session.Save()

	c.Redirect(http.StatusTemporaryRedirect, fmt.Sprintf("%s/?login=success", clientURL))
}

func GetMe(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user_id")
	if userID == nil {
		c.JSON(http.StatusOK, gin.H{"user": nil})
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": gin.H{
		"id":     userID,
		"name":   session.Get("user_name"),
		"email":  session.Get("user_email"),
		"avatar": session.Get("user_avatar"),
		"role":   session.Get("user_role"),
	}})
}

func Logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Clear()
	session.Save()
	c.JSON(http.StatusOK, gin.H{"success": true})
}
