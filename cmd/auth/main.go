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
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func main() {
	router := gin.Default() 

	router.GET("/", func(ctx *gin.Context) {
		ctx.HTML(http.StatusOK, "/Users/anton.bondarchuk/revise/templates/index.html", nil)
	})

	router.POST("/auth/google/login", OauthGoogleByLogin)
	router.POST("/auth/google/callback", OauthGoogleCallback)
}


var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8000/auth/google/callback",
	ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

func OauthGoogleByLogin(c *gin.Context) {
	// ctx := c.Request.Context()

	oauthState := generateStateOauthCookie(c)
	redirectUrl  := googleOauthConfig.AuthCodeURL(oauthState)
	c.Redirect(http.StatusTemporaryRedirect, redirectUrl)
}

func OauthGoogleCallback(c *gin.Context) {
	oauthState, _ := c.Cookie("oauthstate")

	if c.Request.FormValue("state") != oauthState {
		log.Println("invalid oauth google state")
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	data, err := getUserDataFromGoogle(c.Request.FormValue("code"))
	if err != nil {
		log.Println(err.Error())
		c.Redirect(http.StatusTemporaryRedirect, "/")
		return
	}

	// TODO: GetOrCreate User in your db.
	c.String(http.StatusOK, "user info: %s\n", data)
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

func getUserDataFromGoogle(code string) ([]byte, error) {
	// TODO: invistigate Exchange fn logic
	// : check context correct handling
	// : move oauthGoogleUrlAPI to config
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	resp, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
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

