package database

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Sqlite struct {
	config map[string]interface{}
}

func (s *Sqlite) Conn() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(s.config["file"].(string)), &gorm.Config{})
	return db
}
