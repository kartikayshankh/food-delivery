package handler

import (
	healthService "assignment/service/health"
	"assignment/service/rider"
	healthEndpoint "assignment/transport/health"
	restaurantEndpoint "assignment/transport/restaurant"
	riderEndpoint "assignment/transport/rider"
	userEndpoint "assignment/transport/user"
	"net/http"
	"strings"
	"time"

	customMiddleware "assignment/handler/middleware"
	"assignment/service/restaurant"
	"assignment/service/user"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/labstack/gommon/random"
	"github.com/spf13/viper"
)

func MakeHTTPHandler(configViper *viper.Viper, healthService healthService.Service, userService user.UserService, riderService rider.RiderService, restaurantService restaurant.RestaurantService) *echo.Echo {

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return random.String(32)
		},
		TargetHeader: echo.HeaderXRequestID,
	}))

	allowOriginsStr := configViper.GetString("service.allow_origins")
	if len(allowOriginsStr) == 0 {
		log.Fatalf("No allowed origins set")
	}

	log.Printf("\n Allowed origins picked are : %s \n", allowOriginsStr)

	allowOrigins := strings.Split(allowOriginsStr, ",")

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     allowOrigins,
		AllowCredentials: true,
		AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodOptions, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
		AllowHeaders:     []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAuthorization, echo.HeaderAccept, echo.HeaderContentDisposition},
		ExposeHeaders:    []string{"Authorization", "X-Captcha-Id", echo.HeaderContentType, "X-XSS-Protection", " X-Content-Type-Options", "X-Frame-Options"}, // we can add debug here --- "Debug" ---
	}))

	config := middleware.RateLimiterConfig{
		Skipper: middleware.DefaultSkipper,
		Store: middleware.NewRateLimiterMemoryStoreWithConfig(
			middleware.RateLimiterMemoryStoreConfig{Rate: 100, Burst: 200, ExpiresIn: 3 * time.Minute},
		),
		IdentifierExtractor: func(ctx echo.Context) (string, error) {
			id := ctx.RealIP()
			return id, nil
		},
		ErrorHandler: func(context echo.Context, err error) error {
			return context.JSON(http.StatusForbidden, nil)
		},
		DenyHandler: func(context echo.Context, identifier string, err error) error {
			return context.JSON(http.StatusTooManyRequests, nil)
		},
	}

	e.Use(middleware.RateLimiterWithConfig(config))

	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize: 1 << 10, // 1 KB
		LogLevel:  log.ERROR,
	}))

	e.Use(middleware.TimeoutWithConfig(middleware.TimeoutConfig{
		Timeout: 20 * time.Second,
	}))

	healthEndpoint := healthEndpoint.NewEndpoint(healthService)
	userEndpoint := userEndpoint.NewUserEndpoint(userService)
	riderEndpoint := riderEndpoint.NewRiderEndpoint(riderService)
	restaurantEndpoint := restaurantEndpoint.NewRestaurantEndpoint(restaurantService)
	e.GET("/health", healthEndpoint.Health, customMiddleware.SecurityMiddleware)

	//users Endpoint
	user := e.Group("/user")
	user.POST("/v1/register", userEndpoint.Register, customMiddleware.AuthMiddleware)
	user.GET("/v1/:id/orders", userEndpoint.GetUserOrderHistory, customMiddleware.AuthMiddleware)

	rider := e.Group("/rider")
	rider.POST("/v1/register", riderEndpoint.Register, customMiddleware.AuthMiddleware)
	rider.PUT("/v1/location/:id", riderEndpoint.UpdateRiderLocation, customMiddleware.AuthMiddleware)
	rider.GET("/v1/:id/orders", riderEndpoint.GetRiderOrderHistory, customMiddleware.AuthMiddleware)

	restaurant := e.Group("/restaurant")
	restaurant.POST("/v1/register", restaurantEndpoint.Register, customMiddleware.AuthMiddleware)
	restaurant.POST("/v1/nearest-rider/:id", restaurantEndpoint.NearestRider, customMiddleware.AuthMiddleware)
	restaurant.POST("/v1/riderlocation", restaurantEndpoint.GetRiderLocation, customMiddleware.AuthMiddleware)
	restaurant.POST("/v1/suggest/restaurant", restaurantEndpoint.SuggestRestaurant, customMiddleware.AuthMiddleware)
	restaurant.GET("/v1/:id/menu", restaurantEndpoint.GetRestaurantMenu, customMiddleware.AuthMiddleware)
	restaurant.POST("/v1/order", restaurantEndpoint.AcceptOrder, customMiddleware.AuthMiddleware)
	return e
}
