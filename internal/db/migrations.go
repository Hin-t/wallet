package db

//import (
//	"gorm.io/driver/mysql"
//	"gorm.io/gorm"
//	"wallet/internal/repository"
//)

//func migrations() {
//	database, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:3306)/wasmchain?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{})
//	if err != nil {
//		return
//	}
//	err = database.AutoMigrate(&repository.Account{})
//	//err = db.AutoMigrate(&repository.Transaction{})
//	if err != nil {
//		return
//	}
//}
