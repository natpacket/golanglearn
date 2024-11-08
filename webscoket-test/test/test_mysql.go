package test

import (
	"encoding/json"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

// Model Struct
type DeviceInfo struct {
	Name           string
	Username       string
	Status         int
	Interval       int
	LastActiveTime time.Time
}

func (DeviceInfo) TableName() string {
	return "acq_finder_pc_device"
}

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

func TestMysql() {

	deviceInfo := &DeviceInfo{}
	DB.First(&deviceInfo, "username = ?", "faker-username")
	arr, _ := json.Marshal(deviceInfo)
	println(string(arr))

	type user struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
	}
	type user1 struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	newUser1 := user1{
		Id:   1,
		Name: "杉杉",
		Age:  18,
	}

	var newInterface1 interface{}

	//第一种使用interface
	newInterface1 = newUser1

	resByre, resByteErr := json.Marshal(newInterface1)
	if resByteErr != nil {
		fmt.Printf("%v", resByteErr)
		return
	}
	var newUser user
	err := json.Unmarshal(resByre, &newUser)
	if err != nil {
		fmt.Printf("%v", err)
		return
	}
	fmt.Printf("使用 json: %v", newUser)
}
