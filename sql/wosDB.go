package sql

import (
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func GetDB() *gorm.DB {
	return DB
}

func InitialDB() {
	var dsn, host string

	host = "dpg-cpfc9af109ks73bh74jg-a"
	// host += ".singapore-postgres.render.com"

	dsn = "host=%s user=admin password=ictd2Ts3IUs09e3kKeBAr2r7zj4FP7Ng dbname=noah_db port=5432 sslmode=disable TimeZone=Asia/Shanghai"

	DB = connectDB(fmt.Sprintf(dsn, host))
}

func connectDB(dsn string) *gorm.DB {
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("連接PostgreSQL失敗:", err)
		fmt.Println("========================================")
		fmt.Println("連接資料庫失敗,請聯繫後端")
		fmt.Println("========================================")
		return nil
	}

	db, err := conn.DB()
	if err != nil {
		fmt.Println("連接PostgreSQL成功,找不到DataBase:", err)
		fmt.Println("========================================")
		fmt.Println("連接資料庫成功,找不到DataBase,請聯繫後端")
		fmt.Println("========================================")
		return nil
	}

	db.SetConnMaxLifetime(time.Duration(10) * time.Second)
	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(10)

	return conn
}
