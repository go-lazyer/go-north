package dbutil

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	USER_ID = "user_id" //
	DAY     = "day"     //
	NUM     = "num"     //

	TABLE_NAME = "tip_off" // 表名
)

type Test struct {
	UserId sql.NullString `orm:"user_id" ` //
	Day    sql.NullTime   `orm:"day" `     //
	Num    sql.NullInt64  `orm:"day" `     //
}

func ToStruct(m map[string]any) Test {
	model := Test{}
	if value, ok := m[USER_ID].(string); ok {
		model.UserId = sql.NullString{String: value, Valid: true}
	}
	if value, ok := m[NUM].(int64); ok {
		model.Num = sql.NullInt64{Int64: value, Valid: true}
	}
	if value, ok := m[DAY].(time.Time); ok {
		// 如果已经是 time.Time 类型，则直接使用
		model.Day = sql.NullTime{value, true}
	}
	return model
}

func ToStructs(s []map[string]any) []Test {
	slices := make([]Test, 0)
	for _, m := range s {
		slices = append(slices, ToStruct(m))
	}
	return slices
}

func TestRowsToResults(t *testing.T) {
	// s := make([]any, 0)

	username := "root"
	password := "xxxx"
	host := "test.daoway.cn"
	port := "3306"
	dbname := "daowei"
	connStr := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname)
	// 创建数据库连接
	db, err := sql.Open("mysql", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sql1 := "select * from test"

	params := make([]any, 0)

	rows, err := db.Query(sql1, params...)
	if err != nil {
		fmt.Println(err)
	}

	// results, _ := RowsToResults[s](rows)
	// results, _ := RowsToResults[Test](rows)
	results, _ := RowsToMapSlice(rows)
	fmt.Println(ToStructs(results))
	// fmt.Println("type:", reflect.TypeOf(results))
	// fmt.Println("type:", reflect.TypeOf(results).Kind())
}

func TestQuery(t *testing.T) {
	type args struct {
		sql     string
		params  []any
		results any
		db      *sql.DB
	}
	tests := []struct {
		name    string
		args    args
		want    []any
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Query(tt.args.sql, tt.args.params, tt.args.results, tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("Query() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Query() = %v, want %v", got, tt.want)
			}
		})
	}
}
