package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

// LockData presents passwords for lock stage
type LockData struct {
	Code     string `json:"code" validate:"omitempty,alphanum,lte=8"`
	LockedBy string `json:"locked_by" validate:"required,alphaunicode,gt=0,lte=100"`
	Comment  string `json:"comment,omitempty" validate:"lte=500"`
}

// StandartJSONResponse implement standart json-response
type StandartJSONResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Stage presents stage server model
type Stage struct {
	gorm.Model
	Name     string `gorm:"type:varchar(80);unique;not null" json:"name" validate:"required,alphanum,gt=0,lte=15"`
	LockCode string `gorm:"type:varchar(80)" json:"code"`
	Locked   bool   `gorm:"default:false;not null" json:"locked"`
	LockedBy string `gorm:"type:varchar(80)" json:"locked_by"`
	Comment  string `gorm:"type:varchar(500)" json:"comment"`
}

// CustomValidator for echo
type CustomValidator struct {
	validator *validator.Validate
}

// Validate func
func (cv *CustomValidator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

// CustomContext for echo with db connection
type CustomContext struct {
	echo.Context
	dbConn *gorm.DB
}

// CreateConn func create new connection to database
func CreateConn() *gorm.DB {
	dbPath := viper.GetString("settings.db")
	db, err := gorm.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("failed to connect db: %s", err.Error())
	}

	return db
}

// GetAndValidateLockData func bind lockData and validate it
func GetAndValidateLockData(ctx *CustomContext, lockData *LockData) error {
	if err := ctx.Bind(lockData); err != nil {
		return err
	}

	if err := ctx.Validate(lockData); err != nil {
		return err
	}

	return nil

}

// GetStagesHandler handler
func GetStagesHandler(ctx echo.Context) error {

	context := ctx.(*CustomContext)
	db := context.dbConn

	stages := []*Stage{}

	db.Find(&stages)

	for _, stage := range stages {
		stage.LockCode = "hidden"
	}

	return context.JSON(http.StatusOK, &StandartJSONResponse{
		Status:  "success",
		Message: "Ok",
		Data:    stages,
	})
}

// LockStageHandler intended for locking stage instance
func LockStageHandler(ctx echo.Context) error {

	context := ctx.(*CustomContext)
	db := context.dbConn
	stageName := context.Param("name")
	lockData := &LockData{}

	if err := GetAndValidateLockData(context, lockData); err != nil {
		return context.JSON(http.StatusBadRequest, &StandartJSONResponse{
			Status:  "error",
			Message: err.Error(),
		})
	}

	stage := &Stage{}

	if err := db.Where("name = ?", stageName).First(stage).Error; err != nil {
		return context.JSON(http.StatusBadRequest, &StandartJSONResponse{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if stage.Locked {
		return context.JSON(http.StatusForbidden, &StandartJSONResponse{
			Status:  "error",
			Message: fmt.Sprintf("Stage %s already locked", stage.Name),
		})
	}

	stage.Locked = true
	stage.LockCode = lockData.Code
	stage.Comment = lockData.Comment
	stage.LockedBy = lockData.LockedBy

	if err := db.Save(stage).Error; err != nil {
		return context.JSON(http.StatusBadRequest, &StandartJSONResponse{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return context.JSON(http.StatusOK, &StandartJSONResponse{
		Status:  "success",
		Message: fmt.Sprintf("Stage %s was successfully locked", stage.Name),
	})
}

// UnLockStageHandler intended for unlocking stage instance
func UnLockStageHandler(ctx echo.Context) error {

	context := ctx.(*CustomContext)
	db := context.dbConn
	stageName := context.Param("name")
	lockData := &LockData{}

	if err := GetAndValidateLockData(context, lockData); err != nil {
		return context.JSON(http.StatusBadRequest, &StandartJSONResponse{
			Status:  "error",
			Message: err.Error(),
		})
	}

	stage := &Stage{}

	if err := db.Where("name = ?", stageName).First(stage).Error; err != nil {
		return context.JSON(http.StatusBadRequest, &StandartJSONResponse{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if stage.Locked == false {
		return context.JSON(http.StatusBadRequest, &StandartJSONResponse{
			Status:  "error",
			Message: fmt.Sprintf("Stage %s not locked", stage.Name),
		})
	}

	if lockData.Code != stage.LockCode {
		return context.JSON(http.StatusBadRequest, &StandartJSONResponse{
			Status:  "error",
			Message: "Invalid lock code",
		})
	}

	stage.Locked = false
	stage.LockCode = ""
	stage.Comment = ""
	stage.LockedBy = ""

	if err := db.Save(stage).Error; err != nil {
		return context.JSON(http.StatusBadRequest, &StandartJSONResponse{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return context.JSON(http.StatusOK, &StandartJSONResponse{
		Status:  "success",
		Message: fmt.Sprintf("Stage %s was successfully unlocked", stage.Name),
	})
}

// AddStageHandler intened for created new stage server record
func AddStageHandler(ctx echo.Context) error {

	context := ctx.(*CustomContext)
	db := context.dbConn

	stage := &Stage{}

	if err := context.Bind(stage); err != nil {
		return context.JSON(http.StatusBadRequest, &StandartJSONResponse{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := context.Validate(stage); err != nil {
		return context.JSON(http.StatusBadRequest, &StandartJSONResponse{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if err := db.Create(stage).Error; err != nil {
		return context.JSON(http.StatusBadRequest, &StandartJSONResponse{
			Status:  "error",
			Message: err.Error(),
		})
	}

	return context.JSON(http.StatusOK, &StandartJSONResponse{
		Status:  "success",
		Message: "New stage created",
		Data:    stage,
	})

}

// DeleteStageHandler intended for deleting stage instance
func DeleteStageHandler(ctx echo.Context) error {

	context := ctx.(*CustomContext)
	db := context.dbConn

	stageName := context.Param("name")

	stage := &Stage{}

	if err := db.Where("name = ?", stageName).First(stage).Error; err != nil {
		return context.JSON(http.StatusBadRequest, &StandartJSONResponse{
			Status:  "error",
			Message: err.Error(),
		})
	}

	if stage.DeletedAt != nil {
		return context.JSON(http.StatusBadRequest, &StandartJSONResponse{
			Status:  "error",
			Message: fmt.Sprintf("Stage %s already deleted", stageName),
		})
	}

	if stage.Locked {
		return context.JSON(http.StatusBadRequest, &StandartJSONResponse{
			Status:  "error",
			Message: fmt.Sprintf("Stage %s locked. Unlock it first", stageName),
		})
	}

	db.Delete(stage)

	return context.JSON(http.StatusOK, &StandartJSONResponse{
		Status:  "success",
		Message: fmt.Sprintf("Stage %s was successfully deleted", stageName),
	})
}

func main() {
	appConfig := flag.String("config", "./config.toml", "Path to config file (toml)")

	flag.Parse()

	viper.SetConfigFile(*appConfig)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Can't open config file: %s", err.Error())
	}

	db := CreateConn()
	defer db.Close()

	db.DB().SetMaxIdleConns(10)
	db.DB().SetMaxOpenConns(50)
	db.DB().SetConnMaxLifetime(time.Hour)

	db.AutoMigrate(&Stage{})

	app := echo.New()
	app.Debug = viper.GetBool("settings.debug")
	app.Validator = &CustomValidator{validator: validator.New()}

	app.Use(middleware.Logger())
	app.Use(middleware.CORS())
	app.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			customContext := &CustomContext{c, db}
			return next(customContext)
		}
	})

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

	app.GET("/stages", GetStagesHandler)
	app.POST("/stages/:name/lock", LockStageHandler)
	app.POST("/stages/:name/unlock", UnLockStageHandler)
	app.POST("/stages/add", AddStageHandler, authMiddleware)
	app.DELETE("/stages/:name", DeleteStageHandler, authMiddleware)

	serverPort := fmt.Sprintf(":%d", viper.GetInt("settings.port"))

	app.Logger.Fatal(app.Start(serverPort))

}
