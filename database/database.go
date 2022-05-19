package database

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/th3khan/Go-Authentication-with-Fiber/database/models"
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

	if err := engine.Sync2(new(models.User)); err != nil {
		return nil, err
	}

	return engine, nil
}
