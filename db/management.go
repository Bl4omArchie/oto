package db

import (
    "fmt"
    "path/filepath"

    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
)


func OpenDatabase(dbPath string) (*gorm.DB, error) {
	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		return nil, fmt.Errorf("Incorrect path : %s", dbPath)
	}
	db, err := gorm.Open(sqlite.Open(absPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
    return db, nil
}

func Migrate(odb *gorm.DB, models ...any) error {
    err := odb.AutoMigrate(models...)
    return err
}

func GetBy[T any](odb *gorm.DB, key string, value string) (*T, error) {
	var model T
	if err := odb.First(&model, fmt.Sprintf("%s = ?", key), value).Error; err != nil {
		return nil, fmt.Errorf("invalid inputs: %w", err)
	}
	return &model, nil
}

func GetTable[T any](odb *gorm.DB) ([]*T, error) {
	var model []*T
	if err := odb.Find(&model).Error; err != nil {
		return nil, fmt.Errorf("invalid inputs: %w", err)
	}
	return model, nil
}

func UpdateTable[T any](odb *gorm.DB, key string, value string, newColumn string, newValue string) (*T, error) {
	var model T
	if err := odb.Model(&model).
		Where(fmt.Sprintf("%s = ?", key), value).
		Update(newColumn, newValue).Error; err != nil {
		return nil, fmt.Errorf("couldn't update row %s = %s: %w", key, value, err)
	}
	return &model, nil
}
