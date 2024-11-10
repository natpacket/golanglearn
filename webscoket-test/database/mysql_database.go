package database

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func init() {

	var err error
	DB, err = gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:123456@tcp(192.168.0.161:3306)/wx_finder?charset=utf8&parseTime=True&loc=Local", // data source name
		DefaultStringSize:         256,                                                                                   // default size for string fields
		DisableDatetimePrecision:  true,                                                                                  // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,                                                                                  // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,                                                                                  // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,                                                                                 // auto configure based on currently MySQL version
	}), &gorm.Config{})
	if err != nil {
		panic(fmt.Sprintf("【MySQL】连接失败，ERROR：%v", err.Error()))
	}
}
