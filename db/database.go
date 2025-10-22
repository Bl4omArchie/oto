package oto

import (
    "fmt"
    "path/filepath"

    "gorm.io/gorm"
    "gorm.io/driver/sqlite"
)

type OtoDB struct {
    otoDb *gorm.DB
}

func OpenDatabase(dbPath string) (*OtoDB, error) {
	absPath, err := filepath.Abs(dbPath)
	if err != nil {
		return nil, fmt.Errorf("Incorrect path : %s", dbPath)
	}
	db, err := gorm.Open(sqlite.Open(absPath), &gorm.Config{})
	if err != nil {
		return nil, err
	}
    return &OtoDB{otoDb: db}, nil
}


func (db *OtoDB) Migrate(values ...any) {
	db.otoDb.AutoMigrate(values)
}


// get a specific value from a given row and column 
func (db *OtoDB) GetBy(key string, value string) (*any, error) {
	var table any
    if err := db.otoDb.First(&table, fmt.Sprintf("%s = ?", key), value).Error; err != nil {
        return nil, fmt.Errorf("Invalid inputs : %w", err)
    }
    return &table, nil
}

// for a given row, update a value from a given column
func (db *OtoDB) UpdateTabl(key string, value string, newColumn string, newValue string) (*any, error) {
    var table any
    if err := db.otoDb.Model(&table).
        Where(fmt.Sprintf("%s = ?", key), value).
        Update(newColumn, newValue); err  != nil {
            return nil, fmt.Errorf("Couldn't update the row : %s with value %s", key, value)
        }
        return &table, nil
}
