package main

import (
	"fmt"
	"log"
	"restapi-go-crud/app/handler"
	"restapi-go-crud/app/middleware"
	"restapi-go-crud/app/utils"
	appValidator "restapi-go-crud/app/validator"
	"restapi-go-crud/config"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	validatorPkg "github.com/go-playground/validator/v10"
)

func initValidators() error {
	if v, ok := binding.Validator.Engine().(*validatorPkg.Validate); ok {
		fmt.Println("Custom validator registered ****")
		return appValidator.RegisterCustomValidators(v)
	}
	fmt.Println("Could not register validator *****")
	return nil
}

func main() {
	// Connect Postgres DB
	db := config.ConnectDb()

	if err := utils.EnsureSystemAccountExists(db); err != nil {
		log.Fatalf("Failed to ensure system account: %v", err)
	}

	if err := initValidators(); err != nil {
		panic("Failed to register custom validators: " + err.Error())
	}
	// Routes
	r := gin.Default()
	api := r.Group("/api")
	api.Use(middleware.IdempotencyMiddleware()) // sets referenceID in context
	{
		api.POST("/accounts", handler.CreateAccount(db))
		api.GET("/accounts/:account_id", handler.GetAccount(db))
		api.POST("/transactions", handler.CreateTransaction(db))
	}

	log.Println("Server is starting at ****** http://localhost:3000...")
	if err := r.Run(":3000"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}

}
