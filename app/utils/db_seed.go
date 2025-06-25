package utils

import (
	"restapi-go-crud/app/models"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

func EnsureSystemAccountExists(db *gorm.DB) error {
	var count int64
	err := db.Model(&models.Account{}).Where("account_id = ?", 0).Count(&count).Error
	if err != nil {
		return err
	}

	if count == 0 {
		systemAccount := models.Account{
			AccountID: 0,
			Balance:   decimal.NewFromInt(0),
		}
		return db.Create(&systemAccount).Error
	}
	return nil
}
