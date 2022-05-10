package database

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

type Mysql struct {
	config map[string]interface{}
}

func (m *Mysql) Conn() *gorm.DB {
	config := m.config
	dsn := config["username"].(string) + ":" + config["pwd"].(string) + "@tcp(" + config["host"].(string) + ":" + config["port"].(string) + ")/" + config["db"].(string) + "?charset=" + config["charset"].(string) + "&parseTime=True&loc=Local"
	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		Logger:                 logger.Default.LogMode(logger.Silent),
		PrepareStmt:            true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   config["prefix"].(string),
			SingularTable: true,
		},
	})
	if err != nil {
		panic(err)
	}
	sqlDB, _ := d.DB()
	// SetMaxIdleConns 用于设置连接池中空闲连接的最大数量。
	sqlDB.SetMaxIdleConns(int(config["maxIdle"].(float64)))
	return d
}
