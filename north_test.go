package north

import (
	"database/sql"
	"fmt"
	"testing"
	// _ "github.com/go-sql-driver/mysql"
)

const (
	USER_ID = "user_id" //
	DAY     = "day"     //
	NUM     = "num"     //
)

type Test struct {
	UserId sql.NullString `orm:"user_id" ` //
	Day    sql.NullTime   `orm:"day" `     //
	Num    sql.NullInt64  `orm:"num" `     //
}

func TestFull(t *testing.T) {
	username := "root"
	password := "XXXXXX"
	host := "test.daoway.cn"
	port := "3306"
	dbname := "daowei"
	connStr := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname)
	// 创建数据库连接
	ds, err := Open(DRIVER_NAME_MYSQL, connStr, &Config{MaxOpenConns: 10, MaxIdleConns: 10})
	if err != nil {
		fmt.Println(err)
	}
	sql1 := "select t.num from test t where t.user_id='a'"
	params := make([]any, 0)
	fmt.Println(PrepareQuery[Test](sql1, params, ds))
}
