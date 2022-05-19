package database

import (
	"fmt"

	"xorm.io/xorm"
)

func CreateDBEngine() (*xorm.Engine, error) {
	engine, err := xorm.NewEngine("sqlite3", fmt.Sprintf("./database.db"))
	if err != nil {
		return nil, err
	}

	if err = engine.Ping(); err != nil {
		return nil, err
	}

	return engine, nil
}
