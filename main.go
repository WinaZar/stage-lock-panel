package main

import (
	"flag"
	"fmt"
	"stage-lock-panel/auth"
	"stage-lock-panel/handlers"
	"stage-lock-panel/models"
	"stage-lock-panel/utils"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// CreateConn func create new connection to database
func CreateConn(logger echo.Logger) *gorm.DB {
	dbPath := viper.GetString("settings.db")
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})

	if err != nil {
		logger.Fatalf("failed to connect db: %s", err.Error())
	}

	return db
}

func main() {

	appConfig := flag.String("config", "./config.toml", "Path to config file (toml)")

	flag.Parse()

	app := echo.New()

	viper.SetConfigFile(*appConfig)

	if err := viper.ReadInConfig(); err != nil {
		app.Logger.Fatalf("Can't open config file: %s", err.Error())
	}

	db := CreateConn(app.Logger)
	sqlDB, _ := db.DB()

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	db.AutoMigrate(&models.Stage{}, &models.StageHistoryRecord{})

	app.Debug = viper.GetBool("settings.debug")
	app.Validator = utils.Validator

	app.Use(middleware.Logger())
	app.Use(middleware.CORS())
	utils.SetCustomContext(app, db)

	if staticPath := viper.GetString("settings.static-path"); len(staticPath) > 0 {
		app.Static("/static", staticPath)
		app.File("/", fmt.Sprintf("%s/index.html", staticPath))
	}

	auth.GoogleOauthConfig.ClientID = viper.GetString("settings.google.client_id")
	auth.GoogleOauthConfig.ClientSecret = viper.GetString("settings.google.client_secret")

	authMiddleware := middleware.KeyAuthWithConfig(middleware.KeyAuthConfig{
		KeyLookup: "header:X-Admin-Auth",
		Validator: func(key string, c echo.Context) (bool, error) {
			return key == viper.GetString("settings.admin-auth-key"), nil
		},
	})

	app.GET("/stages", handlers.GetStages)
	app.GET("/stages/:name/history", handlers.GetStageHistory)
	app.POST("/stages/:name/lock", handlers.LockStage)
	app.POST("/stages/:name/unlock", handlers.UnLockStage)
	app.POST("/stages/add", handlers.AddStage, authMiddleware)
	app.DELETE("/stages/:name", handlers.DeleteStage, authMiddleware)

	app.GET("/auth/google", handlers.BeginGoogleAuth)
	app.GET("/auth/google/callback", handlers.CompleteGoogleAuth)

	serverPort := fmt.Sprintf(":%d", viper.GetInt("settings.port"))

	app.Logger.Fatal(app.Start(serverPort))

}
