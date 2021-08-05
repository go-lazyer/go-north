package dbutil

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
)

func PrepareCount(sql string, params []interface{}, db *sql.DB) (int64, error) {
	serverMode := os.Getenv("server.mode")
	if serverMode == "dev" {
		fmt.Printf("sql is %v", sql)
		fmt.Printf("params is %v", params)
	}

	stmt, err := db.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(params...)
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	var count int64
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return 0, err
		}
	}
	return count, nil
}
func PrepareFirst(sql string, params []interface{}, structs interface{}, db *sql.DB) error {
	serverMode := os.Getenv("server.mode")
	if serverMode == "dev" {
		fmt.Printf("sql is %v", sql)
		fmt.Printf("params is %v", params)
	}
	stmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	rows, err := stmt.Query(params...)
	if err != nil {
		return err
	}
	defer rows.Close()
	err = RowsToStruct(rows, structs)
	if err != nil {
		return err
	}
	return nil
}

func PrepareQuery(sql string, params []interface{}, results interface{}, db *sql.DB) error {
	if db == nil {
		return errors.New("db not allowed to be nil,need to instantiate yourself")
	}
	serverMode := os.Getenv("server.mode")
	if serverMode == "dev" {
		fmt.Printf("sql is %v", sql)
		fmt.Printf("params is %v", params)
	}
	stmt, err := db.Prepare(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	rows, err := stmt.Query(params...)
	if err != nil {
		return err
	}
	defer rows.Close()
	err = RowsToStructs(rows, results)
	if err != nil {
		return err
	}
	return nil
}

// 预处理插入 返回批量自增ID
func PrepareInsert(sql string, params []interface{}, db *sql.DB) (int64, error) {
	if db == nil {
		return 0, errors.New("db not allowed to be nil,need to instantiate yourself")
	}
	serverMode := os.Getenv("server.mode")
	if serverMode == "dev" {
		fmt.Printf("sql is %v", sql)
		fmt.Printf("params is %v", params)
	}
	stmt, err := db.Prepare(sql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	ret, err := stmt.Exec(params...)
	if err != nil {
		return 0, err
	}
	id, err := ret.LastInsertId() // 新插入数据的id
	if err != nil {
		return 0, err
	}
	return id, nil
}

func PrepareUpdate(sql string, params []interface{}, db *sql.DB) (int64, error) {
	if db == nil {
		return 0, errors.New("db not allowed to be nil,need to instantiate yourself")
	}
	serverMode := os.Getenv("server.mode")
	if serverMode == "dev" {
		fmt.Printf("sql is %v", sql)
		fmt.Printf("params is %v", params)
	}
	ret, err := db.Exec(sql, params...)
	if err != nil {
		return 0, err
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		return 0, err
	}
	return n, nil
}

func getFieldInfo(typ reflect.Type) map[string][]int {
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	finfo := make(map[string][]int)

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		tag := f.Tag.Get("orm")

		// Skip unexported fields or fields marked with "-"
		if f.PkgPath != "" || tag == "-" {
			continue
		}

		// Handle embedded structs
		if f.Anonymous && f.Type.Kind() == reflect.Struct {
			for k, v := range getFieldInfo(f.Type) {
				finfo[k] = append(f.Index, v...)
			}
			continue
		}

		// Use field name for untagged fields
		if tag == "" {
			tag = f.Name
		}

		tag = strings.ToLower(tag)

		finfo[tag] = f.Index
	}

	return finfo
}

func RowsToStructs(rows *sql.Rows, results interface{}) (err error) {
	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	strusRV := reflect.Indirect(reflect.ValueOf(results))
	elemRT := strusRV.Type().Elem()

	fieldInfo := getFieldInfo(elemRT)
	for rows.Next() {
		var struRV reflect.Value
		var struField reflect.Value
		if elemRT.Kind() == reflect.Ptr {
			struRV = reflect.New(elemRT.Elem())
			struField = reflect.Indirect(struRV)
		} else {
			struRV = reflect.Indirect(reflect.New(elemRT))
			struField = struRV
		}
		var values []interface{}
		for _, column := range columns {
			idx, ok := fieldInfo[strings.ToLower(column)]
			var v interface{}
			if !ok {
				var i interface{}
				v = &i
			} else {
				v = struField.FieldByIndex(idx).Addr().Interface()
			}
			values = append(values, v)
		}
		err = rows.Scan(values...)
		if err != nil {
			return err
		}
		strusRV = reflect.Append(strusRV, struRV)
	}
	if err = rows.Err(); err != nil {
		return err
	}
	reflect.Indirect(reflect.ValueOf(results)).Set(strusRV)
	return nil
}
func RowsToStruct(rows *sql.Rows, result interface{}) (err error) {
	struRT := reflect.TypeOf(result).Elem()

	strusPtrRV := reflect.New(reflect.SliceOf(struRT))
	err = RowsToStructs(rows, strusPtrRV.Interface())
	if err != nil {
		return err
	}
	strusRV := reflect.Indirect(strusPtrRV)
	if strusRV.Len() == 0 {
		return
	}
	reflect.Indirect(reflect.ValueOf(result)).Set(strusRV.Index(0))
	return
}
func RowsToCnts(rows *sql.Rows, cnts interface{}) (err error) {
	cntsRV := reflect.Indirect(reflect.ValueOf(cnts))
	elemRT := cntsRV.Type().Elem()

	for rows.Next() {
		var values []interface{}
		var cntRV reflect.Value
		if elemRT.Kind() == reflect.Ptr {
			cntRV = reflect.New(elemRT.Elem())
			values = append(values, cntRV.Interface())
		} else {
			cntRV = reflect.Indirect(reflect.New(elemRT))
			values = append(values, cntRV.Addr().Interface())
		}
		err = rows.Scan(values...)
		if err != nil {
			return
		}
		cntsRV = reflect.Append(cntsRV, cntRV)
	}
	if err = rows.Err(); err != nil {
		return
	}
	reflect.Indirect(reflect.ValueOf(cnts)).Set(cntsRV)

	return
}

func RowsToCnt(rows *sql.Rows, cnt interface{}) (err error) {
	cntRT := reflect.TypeOf(cnt).Elem()

	cntsPtrRV := reflect.New(reflect.SliceOf(cntRT))
	err = RowsToCnts(rows, cntsPtrRV.Interface())
	if err != nil {
		return
	}
	cntsRV := reflect.Indirect(cntsPtrRV)
	if cntsRV.Len() == 0 {
		err = sql.ErrNoRows
		return
	}
	reflect.Indirect(reflect.ValueOf(cnt)).Set(cntsRV.Index(0))
	return
}
