package api

import (
	"fmt"
	"log"

	"athena/src/api/http/middlewares"
	"athena/src/api/http/routes"
	"athena/src/api/http/routes/admins"
	users "athena/src/api/http/routes/users"
	"athena/src/config"
	"athena/src/providers"
	"athena/src/services"

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
	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(middlewares.Cors())
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middlewares.RequestID)
	router.Use(middlewares.I18n)

	globalLimiter := providers.ProvideRateLimiterMiddleware(
		providers.ProvideRateLimiterService(),
	).SetLimiter(services.DefaultLimiter()).SetKey(services.DefaultKeyGetter)
	router.Use(globalLimiter.Middleware)

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
		routes.IpgRouter(v1)
		users.RegisterPaymentTrackingRouter(v1)
	}

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
		admins.RegisterBlockchainRouter(v1)
		admins.RegisterBlockchainExplorerRouter(v1)
		admins.RegisterWalletAddressRouter(v1)
		admins.RegisterPaymentRouter(v1)
	}

	if err := router.RunTLS(
		fmt.Sprintf(":%s", configs.Get("ADMIN_APP_PORT")),
		configs.Get("SSL_CERT_PATH"),
		configs.Get("SSL_KEY_PATH"),
	); err != nil {
		return err
	}
	return nil
}
