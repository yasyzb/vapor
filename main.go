package main

import (
	// "sync"
	// "database/sql"
	"fmt"

	btmTypes "github.com/bytom/protocol/bc/types"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"

	// "github.com/vapor/errors"
	// "github.com/vapor/federation/config"
	// "github.com/vapor/federation/database"
	"github.com/vapor/federation/common"
	"github.com/vapor/federation/database/orm"
	"github.com/vapor/federation/service"
	"github.com/vapor/federation/util"
	vaporTypes "github.com/vapor/protocol/bc/types"
	// "github.com/vapor/federation/synchron"
)

func main() {
	str := "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff"
	asset, _ := util.StringToAssetID(str)
	fmt.Println(asset)

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

	node := service.NewNode("http://127.0.0.1")
	tx1 := &btmTypes.Tx{}
	tx2 := &vaporTypes.Tx{}
	node.SubmitTx(tx1)
	node.SubmitTx(tx2)

	// ormTx := &orm.CrossTransaction{
	// 	ChainID:              1,
	// 	SourceBlockHeight:    2,
	// 	SourceBlockHash:      "blockHash.String()",
	// 	SourceTxIndex:        3,
	// 	SourceMuxID:          "muxID.String()",
	// 	SourceTxHash:         "tx.ID.String()",
	// 	SourceRawTransaction: "string(rawTx)",
	// 	DestBlockHeight:      sql.NullInt64{Valid: false},
	// 	DestBlockHash:        sql.NullString{Valid: false},
	// 	DestTxIndex:          sql.NullInt64{Valid: false},
	// 	DestTxHash:           sql.NullString{Valid: false},
	// 	Status:               common.CrossTxPendingStatus,
	// }
	// if err := db.Create(ormTx).Error; err != nil {
	// 	log.Warnln(err)
	// }
}
