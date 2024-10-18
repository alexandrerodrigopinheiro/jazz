// database/database_text.go
package database

import (
	"jazz/backend/pkg/logger"

	"gorm.io/gorm"
)

// DatabaseText provides basic operations for a text-based table in the database.
type DatabaseText struct {
	db *gorm.DB
}

// NewDatabaseText creates a new instance of DatabaseText.
func NewDatabaseText() *DatabaseText {
	db := GetDBInstance()
	if db == nil {
		logger.Logger.Fatal("Failed to initialize database instance for DatabaseText")
	}
	return &DatabaseText{db: db}
}

// InsertText inserts a new text record into the specified table.
func (dt *DatabaseText) InsertText(tableName string, textRecord interface{}) error {
	if err := dt.db.Table(tableName).Create(textRecord).Error; err != nil {
		logger.Logger.Errorw("Failed to insert text record", "table", tableName, "error", err)
		return err
	}
	logger.Logger.Infof("Successfully inserted text record into table %s", tableName)
	return nil
}

// UpdateText updates an existing text record in the specified table.
func (dt *DatabaseText) UpdateText(tableName string, condition interface{}, updatedValues interface{}) error {
	if err := dt.db.Table(tableName).Where(condition).Updates(updatedValues).Error; err != nil {
		logger.Logger.Errorw("Failed to update text record", "table", tableName, "error", err)
		return err
	}
	logger.Logger.Infof("Successfully updated text record in table %s", tableName)
	return nil
}

// DeleteText deletes a text record from the specified table.
func (dt *DatabaseText) DeleteText(tableName string, condition interface{}) error {
	if err := dt.db.Table(tableName).Where(condition).Delete(nil).Error; err != nil {
		logger.Logger.Errorw("Failed to delete text record", "table", tableName, "error", err)
		return err
	}
	logger.Logger.Infof("Successfully deleted text record from table %s", tableName)
	return nil
}

// FindText finds a text record from the specified table.
func (dt *DatabaseText) FindText(tableName string, condition interface{}, result interface{}) error {
	if err := dt.db.Table(tableName).Where(condition).First(result).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			logger.Logger.Warnw("No text record found", "table", tableName, "condition", condition)
			return nil
		}
		logger.Logger.Errorw("Failed to find text record", "table", tableName, "error", err)
		return err
	}
	logger.Logger.Infof("Successfully found text record in table %s", tableName)
	return nil
}

// FindAllText retrieves all text records from the specified table.
func (dt *DatabaseText) FindAllText(tableName string, condition interface{}, results interface{}) error {
	if err := dt.db.Table(tableName).Where(condition).Find(results).Error; err != nil {
		logger.Logger.Errorw("Failed to find text records", "table", tableName, "error", err)
		return err
	}
	logger.Logger.Infof("Successfully retrieved all text records from table %s", tableName)
	return nil
}
