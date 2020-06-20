package main

import (
	"flag"
	"fmt"
	"stage-lock-panel/handlers"
	"stage-lock-panel/models"
	"stage-lock-panel/utils"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
)

// CreateConn func create new connection to database
func CreateConn(logger echo.Logger) *gorm.DB {
	dbPath := viper.GetString("settings.db")
	db, err := gorm.Open("sqlite3", dbPath)
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
	defer db.Close()

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(50)
	db.DB().SetConnMaxLifetime(time.Hour)

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

	serverPort := fmt.Sprintf(":%d", viper.GetInt("settings.port"))

	app.Logger.Fatal(app.Start(serverPort))

}
