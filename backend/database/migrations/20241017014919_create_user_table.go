package migrations

import (
	"fmt"
	"jazz/backend/models"

	"gorm.io/gorm"
)

// Up is executed when this migration is applied
func Up(db *gorm.DB) {
	fmt.Println("Applying migration: create sample table")
	// Automigrate the model
	db.AutoMigrate(&models.User{})
}

// Down is executed when this migration is reverted
func Down(db *gorm.DB) {
	fmt.Println("Reverting migration: drop sample table")
	// Drop the table associated with the model
	db.Migrator().DropTable(&models.User{})
}
