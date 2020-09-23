package utils

import (
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

// CustomContext for echo with db connection
type CustomContext struct {
	echo.Context
	DB *gorm.DB
}

// SetCustomContext create custom request context with db connection
func SetCustomContext(app *echo.Echo, db *gorm.DB) {
	app.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			customContext := &CustomContext{c, db}
			return next(customContext)
		}
	})
}
