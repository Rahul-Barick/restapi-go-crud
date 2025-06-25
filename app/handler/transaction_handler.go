package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"restapi-go-crud/app/dto"
	"restapi-go-crud/app/models"
	"restapi-go-crud/app/utils"

	"github.com/google/uuid"
)

func CreateTransaction(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req dto.CreateTransactionRequest

		// Strict JSON binding (disallow unknown fields)
		if ok := utils.StrictJSONBind(c, &req); !ok {
			return
		}

		// Validate payload fields
		if err := validate.Struct(req); err != nil {
			if utils.MapValidationErrors(c, err) {
				return
			}
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Get idempotency key from header
		refID, err := utils.GetReferenceID(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		parsedRefID, err := uuid.Parse(refID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid referenceId format"})
			return
		}

		// Begin transaction block
		err = db.Transaction(func(tx *gorm.DB) error {
			// Check if transaction with same reference ID already exists
			var existingTxn models.Transaction
			if err := tx.Where("reference_id = ?", parsedRefID).First(&existingTxn).Error; err == nil {
				c.JSON(http.StatusOK, gin.H{"message": "Transaction already processed"})
				return nil
			}

			// Fetch both accounts with locking
			var accounts []models.Account
			accountIDs := []int{req.SourceAccountID, req.DestinationAccountID}

			if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
				Where("account_id IN ?", accountIDs).
				Find(&accounts).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error while fetching accounts"})
				return err
			}

			// Verify both accounts exist
			if len(accounts) != 2 {
				found := map[int]bool{}
				for _, acc := range accounts {
					found[acc.AccountID] = true
				}

				if !found[req.SourceAccountID] {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Source account not found"})
					return fmt.Errorf("source account not found")
				}
				if !found[req.DestinationAccountID] {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Destination account not found"})
					return fmt.Errorf("destination account not found")
				}
			}

			var source, dest models.Account
			for _, acc := range accounts {
				if acc.AccountID == req.SourceAccountID {
					source = acc
				} else if acc.AccountID == req.DestinationAccountID {
					dest = acc
				}
			}

			// Check sufficient balance
			if source.Balance.LessThan(req.Amount) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Insufficient balance in source account"})
				return fmt.Errorf("insufficient funds")
			}

			// Create transaction record
			txn := models.Transaction{
				ID:                   uuid.New(),
				SourceAccountID:      req.SourceAccountID,
				DestinationAccountID: req.DestinationAccountID,
				Amount:               req.Amount,
				ReferenceID:          parsedRefID,
			}
			if err := tx.Create(&txn).Error; err != nil {
				return err
			}

			// Update balances
			source.Balance = source.Balance.Sub(req.Amount)
			dest.Balance = dest.Balance.Add(req.Amount)

			if err := tx.Save(&source).Error; err != nil {
				return err
			}
			if err := tx.Save(&dest).Error; err != nil {
				return err
			}

			// Create debit and credit ledger entries
			debitEntry := models.LedgerEntries{
				AccountID:     source.AccountID,
				TransactionID: txn.ID,
				Amount:        req.Amount.Neg(),
				Type:          models.DEBIT,
			}
			creditEntry := models.LedgerEntries{
				AccountID:     dest.AccountID,
				TransactionID: txn.ID,
				Amount:        req.Amount,
				Type:          models.CREDIT,
			}

			if err := tx.Create(&debitEntry).Error; err != nil {
				return err
			}
			if err := tx.Create(&creditEntry).Error; err != nil {
				return err
			}

			// Response on success
			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("Transaction successfully created with transaction ID: %s", txn.ID.String()),
			})

			return nil
		})

		if err != nil {
			log.Println("Transaction failed:", err)
			if c.IsAborted() {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process transaction"})
				c.Abort()
			}
			return
		}
	}
}
