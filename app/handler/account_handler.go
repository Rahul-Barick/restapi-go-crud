// Package handler contains the API route handlers for accounts creation and fetching account details.

package handler

import (
	"fmt"
	"net/http"
	"restapi-go-crud/app/constants"
	"restapi-go-crud/app/dto"
	"restapi-go-crud/app/models"
	"restapi-go-crud/app/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var validate = validator.New()

func CreateAccount(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract idempotency key from context (used to prevent duplicate submissions)
		refID, err := utils.GetReferenceID(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		var req dto.CreateAccountRequest
		if ok := utils.StrictJSONBind(c, &req); !ok {
			return
		}

		// Validate request fields using struct-level rules
		if err := validate.Struct(req); err != nil {
			if utils.MapValidationErrors(c, err) {
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Ensure account_id is unique to prevent duplicate account creation
		var existingAccount models.Account
		if err := db.Where("account_id = ?", req.AccountID).First(&existingAccount).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "account_id already exists"})
			return
		} else if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Check idempotency key in transaction table
		var existingTxn models.Transaction
		if err := db.Where("reference_id = ?", refID).First(&existingTxn).Error; err == nil {
			c.JSON(http.StatusOK, gin.H{"message": "account creation already processed"})
			return
		} else if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// Wrapping all DB operations in a db transaction to ensure atomicity
		err = db.Transaction(func(tx *gorm.DB) error {
			parsedRefID, err := uuid.Parse(refID)
			if err != nil {
				return fmt.Errorf("invalid referenceId format")
			}

			account := models.Account{
				AccountID: req.AccountID,
				Balance:   req.InitialBalance,
			}
			// Create new account
			if err := tx.Create(&account).Error; err != nil {
				return err
			}

			// Create a new account with the given initial balance
			txn := models.Transaction{
				SourceAccountID:      constants.SystemAccountID,
				DestinationAccountID: req.AccountID,
				Amount:               req.InitialBalance,
				ReferenceID:          parsedRefID,
			}
			if err := tx.Create(&txn).Error; err != nil {
				return err
			}

			// Add ledger entry
			ledger := models.LedgerEntries{
				AccountID:     req.AccountID,
				TransactionID: txn.ID,
				Amount:        req.InitialBalance,
				Type:          models.CREDIT,
			}
			if err := tx.Create(&ledger).Error; err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			if err.Error() == "invalid referenceId format" {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create account"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{})
	}
}

func GetAccount(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		accountIDParam, err := strconv.Atoi(c.Param("account_id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid account_id"})
			return
		}

		var account models.Account
		if err := db.Where("account_id = ?", accountIDParam).First(&account).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusNotFound, gin.H{"error": "Account not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"account_id": account.AccountID,
			"balance":    account.Balance.StringFixed(5),
		})
	}
}
