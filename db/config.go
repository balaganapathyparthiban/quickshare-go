package db

import (
	"log"

	"github.com/syndtr/goleveldb/leveldb"
)

var Store *leveldb.DB

func InitDB() {
	db, err := leveldb.OpenFile("./db/data", nil)
	if err != nil {
		log.Panic(err)
		return
	}

	Store = db
}
