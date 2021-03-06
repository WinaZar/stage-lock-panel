package handlers

import (
	"fmt"
	"net/http"
	"stage-lock-panel/auth"
	"stage-lock-panel/models"
	"stage-lock-panel/utils"
	"time"

	"github.com/labstack/echo/v4"
)

// GetStages handler
func GetStages(ctx echo.Context) error {

	context := ctx.(*utils.CustomContext)
	db := context.DB

	stages := []*models.Stage{}

	db.Find(&stages)

	for _, stage := range stages {
		stage.LockCode = "hidden"
	}

	return context.JSON(http.StatusOK, &utils.StandartJSONResponse{
		Status:  "success",
		Message: "Ok",
		Data:    stages,
	})
}

// LockStage handler intended for locking stage instance
func LockStage(ctx echo.Context) error {

	context := ctx.(*utils.CustomContext)
	db := context.DB
	stageName := context.Param("name")
	lockData := &utils.LockData{}

	if inputError := utils.GetAndValidateLockData(context, lockData); len(inputError) > 0 {
		return context.JSON(http.StatusBadRequest, &utils.StandartJSONResponse{
			Status:  "error",
			Message: "Error occurred",
			Data:    inputError,
		})
	}

	stage := &models.Stage{}

	if err := db.Where("name = ?", stageName).First(stage).Error; err != nil {
		context.Logger().Error(err)
		return context.JSON(http.StatusBadRequest, &utils.StandartJSONResponse{
			Status:  "error",
			Message: "Stage not found",
		})
	}

	if stage.Locked {
		return context.JSON(http.StatusForbidden, &utils.StandartJSONResponse{
			Status:  "error",
			Message: fmt.Sprintf("Stage %s already locked", stage.Name),
		})
	}

	txn := db.Begin()

	stage.Locked = true
	stage.LockCode = lockData.Code
	stage.Comment = lockData.Comment
	stage.LockedBy = lockData.LockedBy

	if err := txn.Save(stage).Error; err != nil {
		txn.Rollback()
		context.Logger().Error(err)
		return context.JSON(http.StatusBadRequest, &utils.StandartJSONResponse{
			Status:  "error",
			Message: "Something went wrong",
		})
	}

	historyRecord := &models.StageHistoryRecord{
		Action:   "lock",
		LockedBy: lockData.LockedBy,
		Comment:  lockData.Comment,
		Stage:    stage.Name,
	}

	if err := txn.Create(historyRecord).Error; err != nil {
		txn.Rollback()
		context.Logger().Error(err)
		return context.JSON(http.StatusBadRequest, &utils.StandartJSONResponse{
			Status:  "error",
			Message: "Something went wrong",
		})
	}

	txn.Commit()

	return context.JSON(http.StatusOK, &utils.StandartJSONResponse{
		Status:  "success",
		Message: fmt.Sprintf("Stage %s was successfully locked", stage.Name),
	})
}

// UnLockStage intended for unlocking stage instance
func UnLockStage(ctx echo.Context) error {

	context := ctx.(*utils.CustomContext)
	db := context.DB
	stageName := context.Param("name")
	lockData := &utils.LockData{}

	if inputError := utils.GetAndValidateLockData(context, lockData); len(inputError) > 0 {
		return context.JSON(http.StatusBadRequest, &utils.StandartJSONResponse{
			Status:  "error",
			Message: "Error occurred",
			Data:    inputError,
		})
	}

	stage := &models.Stage{}

	if err := db.Where("name = ?", stageName).First(stage).Error; err != nil {
		context.Logger().Error(err)
		return context.JSON(http.StatusBadRequest, &utils.StandartJSONResponse{
			Status:  "error",
			Message: "Stage not found",
		})
	}

	if stage.Locked == false {
		return context.JSON(http.StatusBadRequest, &utils.StandartJSONResponse{
			Status:  "error",
			Message: fmt.Sprintf("Stage %s not locked", stage.Name),
		})
	}

	if lockData.Code != stage.LockCode {
		return context.JSON(http.StatusBadRequest, &utils.StandartJSONResponse{
			Status:  "error",
			Message: "Invalid lock code",
		})
	}

	txn := db.Begin()

	stage.Locked = false
	stage.LockCode = ""
	stage.Comment = ""
	stage.LockedBy = ""

	if err := txn.Save(stage).Error; err != nil {
		txn.Rollback()
		context.Logger().Error(err)
		return context.JSON(http.StatusBadRequest, &utils.StandartJSONResponse{
			Status:  "error",
			Message: "Something went wrong",
		})
	}

	historyRecord := &models.StageHistoryRecord{
		Action: "unlock",
		Stage:  stage.Name,
	}

	if err := txn.Create(historyRecord).Error; err != nil {
		txn.Rollback()
		context.Logger().Error(err)
		return context.JSON(http.StatusBadRequest, &utils.StandartJSONResponse{
			Status:  "error",
			Message: "Something went wrong",
		})
	}

	txn.Commit()

	return context.JSON(http.StatusOK, &utils.StandartJSONResponse{
		Status:  "success",
		Message: fmt.Sprintf("Stage %s was successfully unlocked", stage.Name),
	})
}

// AddStage intened for created new stage server record
func AddStage(ctx echo.Context) error {

	context := ctx.(*utils.CustomContext)
	db := context.DB

	stage := &models.Stage{}

	if err := context.Bind(stage); err != nil {
		context.Logger().Error(err)
		return context.JSON(http.StatusBadRequest, &utils.StandartJSONResponse{
			Status:  "error",
			Message: "Can't bind context data",
		})
	}

	if err := context.Validate(stage); err != nil {
		return context.JSON(http.StatusBadRequest, &utils.StandartJSONResponse{
			Status:  "error",
			Message: "Error occurred",
			Data:    utils.ParseValidatorErrors(err),
		})
	}

	if err := db.Create(stage).Error; err != nil {
		context.Logger().Error(err)
		return context.JSON(http.StatusBadRequest, &utils.StandartJSONResponse{
			Status:  "error",
			Message: "Database error occured",
		})
	}

	return context.JSON(http.StatusOK, &utils.StandartJSONResponse{
		Status:  "success",
		Message: "New stage created",
		Data:    stage,
	})
}

// DeleteStage intended for deleting stage instance
func DeleteStage(ctx echo.Context) error {

	context := ctx.(*utils.CustomContext)
	db := context.DB

	stageName := context.Param("name")

	stage := &models.Stage{}

	if err := db.Where("name = ?", stageName).First(stage).Error; err != nil {
		context.Logger().Error(err)
		return context.JSON(http.StatusBadRequest, &utils.StandartJSONResponse{
			Status:  "error",
			Message: "Stage not found",
		})
	}

	if stage.DeletedAt != nil {
		return context.JSON(http.StatusBadRequest, &utils.StandartJSONResponse{
			Status:  "error",
			Message: fmt.Sprintf("Stage %s already deleted", stageName),
		})
	}

	if stage.Locked {
		return context.JSON(http.StatusBadRequest, &utils.StandartJSONResponse{
			Status:  "error",
			Message: fmt.Sprintf("Stage %s locked. Unlock it first", stageName),
		})
	}

	db.Delete(stage)

	return context.JSON(http.StatusOK, &utils.StandartJSONResponse{
		Status:  "success",
		Message: fmt.Sprintf("Stage %s was successfully deleted", stageName),
	})
}

// GetStageHistory return stage history (locks/unlocks)
func GetStageHistory(ctx echo.Context) error {
	context := ctx.(*utils.CustomContext)
	db := context.DB
	stageName := context.Param("name")
	pagination := &utils.PaginationData{}
	sort := &utils.SortData{}

	if err := context.Bind(pagination); err != nil {
		context.Logger().Error(err)
		return context.JSON(http.StatusBadRequest, &utils.StandartJSONResponse{
			Status:  "error",
			Message: "Context bind error occurred",
		})
	}

	if err := context.Validate(pagination); err != nil {
		return context.JSON(http.StatusBadRequest, &utils.StandartJSONResponse{
			Status:  "error",
			Message: "Error occurred",
			Data:    utils.ParseValidatorErrors(err),
		})
	}

	if err := context.Bind(sort); err != nil {
		context.Logger().Error(err)
		return context.JSON(http.StatusBadRequest, &utils.StandartJSONResponse{
			Status:  "error",
			Message: "Context bind error occurred",
		})
	}

	if err := context.Validate(sort); err != nil {
		return context.JSON(http.StatusBadRequest, &utils.StandartJSONResponse{
			Status:  "error",
			Message: "Error occurred",
			Data:    utils.ParseValidatorErrors(err),
		})
	}

	history := []models.StageHistoryRecord{}
	stage := &models.Stage{}

	if err := db.Where("name = ?", stageName).First(stage).Error; err != nil {
		context.Logger().Error(err)
		return context.JSON(http.StatusBadRequest, &utils.StandartJSONResponse{
			Status:  "error",
			Message: "Stage not found",
		})
	}
	offset := (pagination.Page - 1) * pagination.PerPage

	type result struct {
		Pagination utils.PaginationData        `json:"pagination"`
		History    []models.StageHistoryRecord `json:"history"`
	}

	pagination.TotalItems = int(db.Model(stage).Association("History").Count())

	order := fmt.Sprintf("%s %s", sort.SortBy, sort.SortOrder)

	db.Order(order).Offset(offset).Limit(pagination.PerPage).Model(stage).Association("History").Find(&history)

	return context.JSON(http.StatusOK, &utils.StandartJSONResponse{
		Status:  "success",
		Message: "Ok",
		Data:    &result{Pagination: *pagination, History: history},
	})

}

// BeginGoogleAuth handler
func BeginGoogleAuth(ctx echo.Context) error {
	context := ctx.(*utils.CustomContext)

	// Create oauthState cookie
	oauthState := auth.GenerateOauthState()

	cookie := http.Cookie{Name: "oauthstate", Value: oauthState, Expires: time.Now().Add(365 * 24 * time.Hour)}
	context.SetCookie(&cookie)

	/*
		AuthCodeURL receive state that is a token to protect the user from CSRF attacks. You must always provide a non-empty string and
		validate that it matches the the state query parameter on your redirect callback.
	*/
	redirectURL := auth.GoogleOauthConfig.AuthCodeURL(oauthState)

	return context.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

// CompleteGoogleAuth handler
func CompleteGoogleAuth(ctx echo.Context) error {
	context := ctx.(*utils.CustomContext)

	oauthState, _ := context.Cookie("oauthstate")

	if context.FormValue("state") != oauthState.Value {
		return context.JSON(http.StatusBadRequest, &utils.StandartJSONResponse{
			Status:  "error",
			Message: "Invalid OAuth state",
		})
	}

	data, err := auth.GetUserDataFromGoogle(context.FormValue("code"))
	if err != nil {
		context.Logger().Error(err)
		return context.JSON(http.StatusBadRequest, &utils.StandartJSONResponse{
			Status:  "error",
			Message: "Auth error",
		})
	}

	// GetOrCreate User in your db.
	// Redirect or response with a token.
	// More code .....
	fmt.Printf("UserInfo: %s\n", data)

	return context.JSON(http.StatusOK, &utils.StandartJSONResponse{
		Status:  "success",
		Message: "Ok",
	})
}
