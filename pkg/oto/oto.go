package oto

import (
	"gorm.io/gorm"

	"github.om/Bl4omArchie/oto/db"
	"github.om/Bl4omArchie/oto/models"
)

type Oto struct {
	Database *gorm.DB
}

func OpenOto(dbPath string) (*Oto, error) {
	otoDb, err := db.OpenDatabase(dbPath)
	if err != nil {
		return nil, err
	}
	oto := &Oto{
		Database: otoDb,
	}
	oto.RefreshOto()
	return oto, nil
}

func (oto *Oto) RefreshOto() {
	db.Migrate(oto.Database, &models.Executable{}, &models.Parameter{}, &models.Command{})
}
