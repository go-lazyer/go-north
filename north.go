package north

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"reflect"
	"runtime"
	"strings"
)

const (
	DRIVER_NAME_POSTGRES = "postgres"
	DRIVER_NAME_MYSQL    = "mysql"
	PLACE_HOLDER_GO      = "ⒼⓄ" //
)

type North struct {
	DataSources []*DataSource
}

func New() *North {
	return &North{}
}

func (north *North) Add(dataSource *DataSource) {
	if north.DataSources == nil {
		north.DataSources = make([]*DataSource, 0)
	}
	north.DataSources = append(north.DataSources, dataSource)
}
func (north *North) getDataSource(filePath string) *DataSource {
	// if len(north.DataSources) == 1 {
	// 	return north.DataSources[0]
	// }
	for _, dataSource := range north.DataSources {
		for _, daoFilePath := range dataSource.DaoFilePaths {
			if strings.Contains(filePath, daoFilePath) {
				return dataSource
			}
		}
	}
	return nil
}

func (north *North) Count(sql string, params []any) (int64, error) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return 0, errors.New("failed to get caller infof")
	}
	ds := north.getDataSource(file)
	if ds.Db == nil {
		return 0, errors.New("db not allowed to be nil,need to instantiate yourself")
	}
	serverMode := os.Getenv("sql.log")
	if serverMode == "stdout" {
		fmt.Printf("sql is %v\n", sql)
		fmt.Printf("params is %v\n", params)
	}

	rows, err := ds.Db.Query(sql, params...)
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

func (north *North) PrepareCount(sql string, params []any) (int64, error) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return 0, errors.New("failed to get caller infof")
	}
	ds := north.getDataSource(file)
	if ds.Db == nil {
		return 0, errors.New("db not allowed to be nil,need to instantiate yourself")
	}
	sql = prepareConvert(sql, ds.DriverName)

	serverMode := os.Getenv("sql.log")
	if serverMode == "stdout" {
		fmt.Printf("sql is %v\n", sql)
		fmt.Printf("params is %v\n", params)
	}

	stmt, err := ds.Db.Prepare(sql)
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

// 普通查询
func (north *North) Query(sql string, params []any) ([]map[string]any, error) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return nil, errors.New("failed to get caller infof")
	}
	ds := north.getDataSource(file)
	if ds.Db == nil {
		return nil, errors.New("db not allowed to be nil,need to instantiate yourself")
	}
	serverMode := os.Getenv("sql.log")
	if serverMode == "stdout" {
		fmt.Printf("sql is %v\n", sql)
		fmt.Printf("params is %v\n", params)
	}
	rows, err := ds.Db.Query(sql, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return RowsToMapSlice(rows)
}

// 预处理查询
func (north *North) PrepareQuery(sql string, params []any) ([]map[string]any, error) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return nil, errors.New("failed to get caller infof")
	}
	ds := north.getDataSource(file)
	if ds.Db == nil {
		return nil, errors.New("db not allowed to be nil,need to instantiate yourself")
	}
	sql = prepareConvert(sql, ds.DriverName)
	serverMode := os.Getenv("sql.log")
	if serverMode == "stdout" {
		fmt.Printf("sql is %v\n", sql)
		fmt.Printf("params is %v\n", params)
	}
	stmt, err := ds.Db.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return RowsToMapSlice(rows)
}

// 预处理插入 返回批量自增ID
func (north *North) PrepareInsert(sql string, params []any) (int64, error) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return 0, errors.New("failed to get caller infof")
	}
	ds := north.getDataSource(file)
	if ds.Db == nil {
		return 0, errors.New("db not allowed to be nil,need to instantiate yourself")
	}
	sql = prepareConvert(sql, ds.DriverName)
	serverMode := os.Getenv("sql.log")
	if serverMode == "stdout" {
		fmt.Printf("sql is %v\n", sql)
		fmt.Printf("params is %v\n", params)
	}
	stmt, err := ds.Db.Prepare(sql)
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

func (north *North) PrepareUpdate(sql string, params []any) (int64, error) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return 0, errors.New("failed to get caller infof")
	}
	ds := north.getDataSource(file)
	if ds.Db == nil {
		return 0, errors.New("db not allowed to be nil,need to instantiate yourself")
	}
	sql = prepareConvert(sql, ds.DriverName)
	serverMode := os.Getenv("sql.log")
	if serverMode == "stdout" {
		fmt.Printf("sql is %v\n", sql)
		fmt.Printf("params is %v\n", params)
	}
	ret, err := ds.Db.Exec(sql, params...)
	if err != nil {
		return 0, err
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		return 0, err
	}
	return n, nil
}
func (north *North) PrepareSave(sql string, params []any) (int64, error) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return 0, errors.New("failed to get caller infof")
	}
	ds := north.getDataSource(file)
	if ds.Db == nil {
		return 0, errors.New("db not allowed to be nil,need to instantiate yourself")
	}
	sql = prepareConvert(sql, ds.DriverName)
	serverMode := os.Getenv("sql.log")
	if serverMode == "stdout" {
		fmt.Printf("sql is %v\n", sql)
		fmt.Printf("params is %v\n", params)
	}
	ret, err := ds.Db.Exec(sql, params...)
	if err != nil {
		return 0, err
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		return 0, err
	}
	return n, nil
}
func (north *North) PrepareDelete(sql string, params []any) (int64, error) {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return 0, errors.New("failed to get caller infof")
	}
	ds := north.getDataSource(file)
	if ds.Db == nil {
		return 0, errors.New("db not allowed to be nil,need to instantiate yourself")
	}
	sql = prepareConvert(sql, ds.DriverName)
	serverMode := os.Getenv("sql.log")
	if serverMode == "stdout" {
		fmt.Printf("sql is %v\n", sql)
		fmt.Printf("params is %v\n", params)
	}
	ret, err := ds.Db.Exec(sql, params...)
	if err != nil {
		return 0, err
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		return 0, err
	}
	return n, nil
}

// 把查询结果映射为实体
// 支持list
func RowsToResults[T any](rows *sql.Rows) ([]T, error) {
	// if reflect.TypeOf(results).Kind() != reflect.Slice {
	// 	return errors.New("results Must be slice")
	// }

	// fmt.Println("type:", reflect.TypeOf(results))
	// fmt.Println("type:", reflect.TypeOf(results).Kind())

	return RowsToStruct[T](rows)
}

func RowsToStruct[T any](rows *sql.Rows) ([]T, error) {
	//创建一个切片的类型
	sliceType := reflect.SliceOf(reflect.TypeOf(new(T)).Elem())
	//创建一个具有给定长度和容量的切片
	sliceValue := reflect.MakeSlice(sliceType, 0, 0)

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// 构建字段到列索引的映射
	fieldToColIndex := make(map[string]int)
	structType := reflect.TypeOf(new(T)).Elem()
	for i, columnName := range columns {
		for j := 0; j < structType.NumField(); j++ {
			field := structType.Field(j)
			if tagValue, ok := getFieldTagValue(field, "orm"); ok && tagValue == columnName {
				fieldToColIndex[field.Name] = i
				break
			}
		}
	}

	for rows.Next() {
		var elem T
		elemValue := reflect.ValueOf(&elem).Elem()
		scanArgs := make([]any, len(columns))
		for columnName, colIndex := range fieldToColIndex {
			field := elemValue.FieldByName(columnName)
			if !field.IsValid() {
				return nil, fmt.Errorf("field %s not found in type %T", columnName, elem)
			}
			scanArgs[colIndex] = field.Addr().Interface()

		}

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}
		sliceValue = reflect.Append(sliceValue, elemValue)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sliceValue.Interface().([]T), nil
}

// getFieldTagValue 获取结构体字段的tag值
func getFieldTagValue(field reflect.StructField, tagName string) (string, bool) {
	tag, ok := field.Tag.Lookup(tagName)
	if ok {
		parts := strings.Split(tag, ",")
		return parts[0], true
	}
	return "", false
}

// func RowsToStructs(rows *sql.Rows, results any) (err error) {
// 	columns, err := rows.Columns()
// 	if err != nil {
// 		return err
// 	}
// 	strusRV := reflect.Indirect(reflect.ValueOf(results))
// 	elemRT := strusRV.Type().Elem()

// 	fieldInfo := getFieldInfo(elemRT)
// 	for rows.Next() {
// 		var struRV reflect.Value
// 		var struField reflect.Value
// 		if elemRT.Kind() == reflect.Ptr {
// 			struRV = reflect.New(elemRT.Elem())
// 			struField = reflect.Indirect(struRV)
// 		} else {
// 			struRV = reflect.Indirect(reflect.New(elemRT))
// 			struField = struRV
// 		}
// 		var values []any
// 		for _, column := range columns {
// 			idx, ok := fieldInfo[strings.ToLower(column)]
// 			var v any
// 			if !ok {
// 				var i any
// 				v = &i
// 			} else {
// 				v = struField.FieldByIndex(idx).Addr().Interface()
// 			}
// 			values = append(values, v)
// 		}
// 		err = rows.Scan(values...)
// 		if err != nil {
// 			return err
// 		}
// 		strusRV = reflect.Append(strusRV, struRV)
// 	}
// 	if err = rows.Err(); err != nil {
// 		return err
// 	}
// 	reflect.Indirect(reflect.ValueOf(results)).Set(strusRV)
// 	return nil
// }
// func getFieldInfo(typ reflect.Type) map[string][]int {
// 	if typ.Kind() == reflect.Ptr {
// 		typ = typ.Elem()
// 	}
// 	finfo := make(map[string][]int)

// 	for i := 0; i < typ.NumField(); i++ {
// 		f := typ.Field(i)
// 		tag := f.Tag.Get("orm")

// 		// Skip unexported fields or fields marked with "-"
// 		if f.PkgPath != "" || tag == "-" {
// 			continue
// 		}

// 		// Handle embedded structs
// 		if f.Anonymous && f.Type.Kind() == reflect.Struct {
// 			for k, v := range getFieldInfo(f.Type) {
// 				finfo[k] = append(f.Index, v...)
// 			}
// 			continue
// 		}

// 		// Use field name for untagged fields
// 		if tag == "" {
// 			tag = f.Name
// 		}

// 		tag = strings.ToLower(tag)

// 		finfo[tag] = f.Index
// 	}

// 	return finfo
// }

// func RowsToCnts(rows *sql.Rows, cnts any) (err error) {
// 	cntsRV := reflect.Indirect(reflect.ValueOf(cnts))
// 	elemRT := cntsRV.Type().Elem()

// 	for rows.Next() {
// 		var values []any
// 		var cntRV reflect.Value
// 		if elemRT.Kind() == reflect.Ptr {
// 			cntRV = reflect.New(elemRT.Elem())
// 			values = append(values, cntRV.Interface())
// 		} else {
// 			cntRV = reflect.Indirect(reflect.New(elemRT))
// 			values = append(values, cntRV.Addr().Interface())
// 		}
// 		err = rows.Scan(values...)
// 		if err != nil {
// 			return
// 		}
// 		cntsRV = reflect.Append(cntsRV, cntRV)
// 	}
// 	if err = rows.Err(); err != nil {
// 		return
// 	}
// 	reflect.Indirect(reflect.ValueOf(cnts)).Set(cntsRV)

// 	return
// }

// func RowsToCnt(rows *sql.Rows, cnt any) (err error) {
// 	cntRT := reflect.TypeOf(cnt).Elem()

//		cntsPtrRV := reflect.New(reflect.SliceOf(cntRT))
//		err = RowsToCnts(rows, cntsPtrRV.Interface())
//		if err != nil {
//			return
//		}
//		cntsRV := reflect.Indirect(cntsPtrRV)
//		if cntsRV.Len() == 0 {
//			err = sql.ErrNoRows
//			return
//		}
//		reflect.Indirect(reflect.ValueOf(cnt)).Set(cntsRV.Index(0))
//		return
//	}
func RowsToMapSlice(rows *sql.Rows) ([]map[string]any, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	results := make([]map[string]any, 0)
	values := make([]any, len(columns))
	scanArgs := make([]any, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}

	for rows.Next() {
		rowMap := make(map[string]any)
		err := rows.Scan(scanArgs...)
		if err != nil {
			return nil, err
		}

		for i, col := range columns {
			var value any = values[i]
			if b, ok := value.([]byte); ok {
				value = string(b)
			}
			rowMap[col] = value
		}
		results = append(results, rowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func prepareConvert(sqlStr, driverName string) string {
	if driverName == DRIVER_NAME_MYSQL {
		return strings.ReplaceAll(sqlStr, PLACE_HOLDER_GO, "?")
	}
	n := 1
	for strings.Index(sqlStr, PLACE_HOLDER_GO) > 0 {
		sqlStr = strings.Replace(sqlStr, PLACE_HOLDER_GO, fmt.Sprintf("$%v", n), 1)
		n = n + 1
	}
	return sqlStr
}
