package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"

	"github.com/vapor/federation/common"
	// "github.com/vapor/federation/config"
	"github.com/vapor/federation/database/orm"
)

func main2() {
	dsnTemplate := "%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local"
	dsn := fmt.Sprintf(dsnTemplate, "root", "toor", "127.0.0.1", 3306, "federation")
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		panic(err)
	}

	// db.LogMode(true)

	// stmt :=
	db.Debug().Model(&orm.CrossTransaction{}).Where("chain_id != ?", 1).
		Where(&orm.CrossTransaction{
			DestTxHash: sql.NullString{"tx.ID.String()", true},
			Status:     common.CrossTxSubmittedStatus,
		}).UpdateColumn(&orm.CrossTransaction{
		DestBlockHeight: sql.NullInt64{int64(2), true},
		DestBlockHash:   sql.NullString{"blockHash.String()", true},
		DestTxIndex:     sql.NullInt64{int64(3), true},
		Status:          common.CrossTxCompletedStatus,
	})

}
