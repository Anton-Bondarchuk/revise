package main

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"revise/internal/config"
	"revise/internal/handlers"
	"revise/internal/service/document"
	"revise/internal/storage"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)



func main() {
	config := config.MustLoad()
	db, err := storage.New(config.StorageConfig)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}
	defer db.Close()

	docService := service.New(log.Default(), db)
	docHandlers := handlers.New(*docService)

	router := gin.Default()

	router.GET("/documents", docHandlers.GetDocuments)
	router.POST("/documents", docHandlers.SaveDocument)

	/* 
	* Auth
	*/
	oauthConfig := NewGoogleOAuthConfig(config.OauthConfig)
	oauthHandlers := New(*oauthConfig, config.OauthConfig.GoogleAuthURL)

	router.GET("/auth/google/login", oauthHandlers.OauthGoogleLogin)      // Changed to GET for OAuth flow
	router.GET("/auth/google/callback", oauthHandlers.OauthGoogleCallback)

	// TODO: add server config and server port
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("ListenAndServe: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	log.Println("Shutting down server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}
	log.Println("Server exiting")
}


type OauthHandlers struct {
	config oauth2.Config
	googleAuthURL string
}

func New(config oauth2.Config, googleAuthURL string) *OauthHandlers {
	return &OauthHandlers{
		config: config,
		googleAuthURL: googleAuthURL,
	}
}

func NewGoogleOAuthConfig(cfg config.OauthConfig) *oauth2.Config {
    return &oauth2.Config{
		// TODO: add server config
        RedirectURL:  "http://localhost:8000" + cfg.RedirectURI,
        ClientID:     cfg.ClientID,
        ClientSecret: cfg.ClientSecret,
        Scopes:       cfg.Scopes,
        Endpoint:     google.Endpoint,
    }
}


func (h *OauthHandlers) OauthGoogleLogin(c *gin.Context) {
	oauthState := generateStateOauthCookie(c)
	redirectUrl := h.config.AuthCodeURL(oauthState)
	c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}

func (h *OauthHandlers) OauthGoogleCallback(c *gin.Context) {
	oauthState, _ := c.Cookie("oauthstate")

	if c.Request.FormValue("state") != oauthState {
		log.Println("invalid oauth google state")
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	// TODO: rework googleAuthURL -> userInfoURL
	data, err := getUserDataFromGoogle(c, &h.config, h.googleAuthURL)
	if err != nil {
		log.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	// TODO: GetOrCreate User in your db.
	// TODO: Process user data and create session/JWT
	c.JSON(http.StatusOK, gin.H{"user": string(data)})
}


func generateStateOauthCookie(c *gin.Context) string {
	// TODO: add time expiration to config
	// rewrite logic set coockie
	expiration := time.Now().Add(120 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	c.SetCookie("oauthstate", state, expiration.Minute(), "/", "/", false, false)

	return state
}

func getUserDataFromGoogle(c *gin.Context, config *oauth2.Config, userInfoURL string) ([]byte, error) {
	// TODO: invistigate Exchange fn logic
	code := c.Query("code")
	token, err := config.Exchange(c, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}

	client := config.Client(c, token)
	resp, err := client.Get(userInfoURL)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer resp.Body.Close()
	
	var content []byte 
	err = json.NewDecoder(resp.Body).Decode(&content)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return content, nil
}
