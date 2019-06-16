package main

import (
	// "sync"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	// "github.com/vapor/errors"
	// "github.com/vapor/federation/config"
	// "github.com/vapor/federation/database"
	"github.com/vapor/federation/common"
	"github.com/vapor/federation/database/orm"
	// "github.com/vapor/federation/synchron"
)

func main() {
	dsnTemplate := "%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local"
	dsn := fmt.Sprintf(dsnTemplate, "root", "toor", "127.0.0.1", 3306, "federation")
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		log.Errorln(err, "open db cluster")
	}
	db.LogMode(true)

	// reqs := []*orm.CrossTransactionReq{}
	// db.Where(&orm.CrossTransactionReq{CrossTransactionID: 1}).Find(&reqs)
	// log.Info(reqs)

	txs := []*orm.CrossTransaction{}
	if err := db.Preload("Chain").Preload("Reqs").Model(&orm.CrossTransaction{}).Where("status = ?", common.CrossTxPendingStatus).Find(&txs).Error; err == gorm.ErrRecordNotFound {
		log.Warnln("ErrRecordNotFound")
	} else if err != nil {
		log.Warnln("collectUnsubmittedTx", err)
	}
}
