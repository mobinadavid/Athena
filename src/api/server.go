package api

import (
	"athena/src/api/http/routes/admins"
	users "athena/src/api/http/routes/users"

	"athena/src/api/http/middlewares"
	"athena/src/api/http/routes"
	"athena/src/config"
	"fmt"
	"log"

	"github.com/gin-contrib/secure"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
)

var (
	configs      = config.GetInstance()
	isProduction = configs.Get("APP_ENV") == "production"
	g            errgroup.Group
)

func Init() (err error) {
	g.Go(func() error {
		return initUserServer()
	})

	g.Go(func() error {
		return initAdminServer()
	})

	if err = g.Wait(); err != nil {
		log.Fatalln(err)
		return err
	}

	return err
}

func getNewRouter() *gin.Engine {
	// set gin to release mode.
	gin.SetMode(gin.ReleaseMode)

	// Initialize new app.
	router := gin.New()

	// Attach CORS middleware.
	router.Use(middlewares.Cors())

	// Attach logger middleware.
	router.Use(gin.Logger())

	// Attach recovery middleware.
	router.Use(gin.Recovery())

	// Attach request id middleware.
	router.Use(middlewares.RequestID)

	// Attach i18n middleware.
	router.Use(middlewares.I18n)

	if isProduction {

		router.Use(secure.New(secure.Config{
			AllowedHosts:          []string{configs.Get("APP_HOST")},
			SSLRedirect:           true,
			SSLHost:               configs.Get("APP_HOST"),
			STSSeconds:            315360000,
			STSIncludeSubdomains:  true,
			FrameDeny:             true,
			ContentTypeNosniff:    true,
			BrowserXssFilter:      true,
			ContentSecurityPolicy: "default-src 'self'",
			IENoOpen:              true,
			ReferrerPolicy:        "strict-origin-when-cross-origin",
			SSLProxyHeaders:       map[string]string{"X-Forwarded-Proto": "https"},
		}))

		// Trusted proxies.
		_ = router.SetTrustedProxies([]string{"https://" + configs.Get("APP_HOST")})
	}

	return router
}

func initUserServer() error {
	router := getNewRouter()

	v1 := router.Group("api/v1")
	{
		users.AuthenticationRouter(v1)
		users.RegisterAccessTokenRouter(v1)
		routes.BlockchainRouter(v1)
		routes.WalletAddressRouter(v1)
		routes.BlockchainExplorerRouter(v1)
		routes.IpgRouter(v1)

	}
	// Run App.
	if err := router.RunTLS(
		fmt.Sprintf(":%s", configs.Get("USER_APP_PORT")),
		configs.Get("SSL_CERT_PATH"),
		configs.Get("SSL_KEY_PATH"),
	); err != nil {
		return err
	}

	return nil
}

func initAdminServer() error {
	router := getNewRouter()

	v1 := router.Group("api/v1")
	{
		admins.AuthenticationRouter(v1)
		admins.RegisterAccessTokensRouter(v1)
		admins.RegisterUserRouter(v1)
		admins.RegisterAuthorizationRouter(v1)
		admins.RegisterAdminRouter(v1)
		routes.BlockchainRouter(v1)
		routes.WalletAddressRouter(v1)
		routes.BlockchainExplorerRouter(v1)
		routes.IpgRouter(v1)
	}

	// Run App.
	if err := router.RunTLS(
		fmt.Sprintf(":%s", configs.Get("ADMIN_APP_PORT")),
		configs.Get("SSL_CERT_PATH"),
		configs.Get("SSL_KEY_PATH"),
	); err != nil {
		return err
	}
	return nil
}
