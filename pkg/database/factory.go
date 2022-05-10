package database

import "gorm.io/gorm"

type Factory struct {
}

func (f *Factory) Init(config map[string]interface{}) *gorm.DB {
	//mysql
	if _, ok := config["mysql"]; ok {
		return f.do(&Mysql{config: config["mysql"].(map[string]interface{})})
	}
	//sqlite
	if _, ok := config["sqlite"]; ok {
		return f.do(&Sqlite{config: config["sqlite"].(map[string]interface{})})
	}

	return nil
}

func (f *Factory) do(d Database) *gorm.DB {
	return d.Conn()
}
