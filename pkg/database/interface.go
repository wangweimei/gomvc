package database

import "gorm.io/gorm"

type Database interface {
	Conn() *gorm.DB
}
