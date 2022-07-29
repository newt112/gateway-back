package main

import (
	"net/http"

	activityRoute "github.com/newt239/gateway-back/routes/activity"
	adminRoute "github.com/newt239/gateway-back/routes/admin"
	authRoute "github.com/newt239/gateway-back/routes/auth"
	exhibitRoute "github.com/newt239/gateway-back/routes/exhibit"
	guestRoute "github.com/newt239/gateway-back/routes/guest"
	reservationRoute "github.com/newt239/gateway-back/routes/reservation"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func Hello() echo.HandlerFunc {
	return func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello World")
	}
}

func main() {
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowCredentials: true,
		AllowOrigins:     []string{"http://localhost:8000", "https://gateway.sh-fes.com", "https://dev.gateway.sh-fes.com"},
		AllowHeaders: []string{
			echo.HeaderAccessControlAllowHeaders,
			echo.HeaderContentType,
			echo.HeaderContentLength,
			echo.HeaderAcceptEncoding,
			echo.HeaderXCSRFToken,
			echo.HeaderAuthorization,
		},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut, http.MethodPatch},
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	v1 := e.Group("/v1")
	v1.GET("/", Hello())

	auth := v1.Group("/auth")
	auth.POST("/login", authRoute.Login())
	auth.Use(middleware.JWT([]byte("secret")))
	auth.GET("/me", authRoute.Me())

	activity := v1.Group("/activity")
	activity.Use(middleware.JWT([]byte("secret")))
	activity.POST("/enter", activityRoute.Enter())
	activity.POST("/exit", activityRoute.Exit())
	activity.GET("/history/:from", activityRoute.History())

	guest := v1.Group("/guest")
	guest.Use(middleware.JWT([]byte("secret")))
	guest.GET("/info/:guest_id", guestRoute.Info())
	guest.GET("/activity/:guest_id", guestRoute.Activity())
	guest.POST("/register", guestRoute.Register())
	guest.POST("/revoke", guestRoute.Revoke())

	reservation := v1.Group("/reservation")
	reservation.Use(middleware.JWT([]byte("secret")))
	reservation.GET("/info/:reservation_id", reservationRoute.Info())

	exhibit := v1.Group("/exhibit")
	exhibit.Use(middleware.JWT([]byte("secret")))
	exhibit.GET("/list", exhibitRoute.ExhibitList())
	exhibit.GET("/info", exhibitRoute.InfoAllExhibit())
	exhibit.GET("/info/:exhibit_id", exhibitRoute.InfoEachExhibit())
	exhibit.GET("/current", exhibitRoute.CurrentAllExhibitData())
	exhibit.GET("/current/:exhibit_id", exhibitRoute.CurrentEachExhibit())
	exhibit.GET("/history/:exhibit_id/:day", exhibitRoute.HistoryEachExhibit())

	admin := v1.Group("/admin")
	admin.Use(middleware.JWT([]byte("secret")))
	admin.POST("/exhibit/create", adminRoute.CreateExhibit())
	admin.DELETE("/exhibit/delete/:exhibit_id", adminRoute.DeleteExhibit())

	e.Start(":3000")
}
