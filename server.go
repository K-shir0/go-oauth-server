package main

import (
	"echo-oauth/config"
	"fmt"
	"github.com/gomodule/oauth1/oauth"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo"
	"github.com/labstack/echo-contrib/session"
	"net/http"
)

type Server struct {
	config *config.Config
	e      *echo.Echo
}

type OAuth struct {
	client oauth.Client
	config *config.Config
}

func NewOAuthHandler(config *config.Config) *OAuth {
	return &OAuth{
		client: oauth.Client{
			TemporaryCredentialRequestURI: config.Twitter.RequestURI,
			ResourceOwnerAuthorizationURI: config.Twitter.AuthorizationURI,
			TokenRequestURI:               config.Twitter.TokenRequestURI,
			Credentials: oauth.Credentials{
				Token:  config.Twitter.Token,
				Secret: config.Twitter.Secret,
			},
		},
		config: config,
	}
}

func main() {

	e := echo.New()

	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret"))))

	e.GET("/", hello)
	e.GET("/twitter/callback", hello)

	e.Logger.Fatal(e.Start(":1323"))
}

func hello(c echo.Context) error {
	config := config.New()

	o := NewOAuthHandler(config)

	// get token
	credentials, err := o.client.RequestTemporaryCredentials(nil, o.config.Twitter.CallbackURI, nil)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	//fmt.Print(credentials)

	s, err := session.Get("test_session", c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	fmt.Print(s)

	s.Options.HttpOnly = true

	s.Values["tokenKey"] = credentials.Token
	s.Values["secretKey"] = credentials.Secret

	if err := s.Save(c.Request(), c.Response()); err != nil {

	}

	return c.JSON(http.StatusOK, o.client.AuthorizationURL(credentials, nil))
}
