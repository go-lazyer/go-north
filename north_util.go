package north

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

const (
	DRIVER_NAME_POSTGRES = "postgres"
	DRIVER_NAME_MYSQL    = "mysql"
	PLACE_HOLDER_GO      = "ⒼⓄ" //
)

type DataSource struct {
	Db           *sql.DB
	DriverName   string
	DaoFilePaths []string
}
type Config struct {
	MaxOpenConns int
	MaxIdleConns int
}

func Open(driverName string, dsn string, config *Config) (*DataSource, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)

	return &DataSource{
		Db:         db,
		DriverName: driverName,
	}, nil
}

func Count(sql string, params []any, ds *DataSource) (int64, error) {
	if ds.Db == nil {
		return 0, errors.New("db not allowed to be nil,need to instantiate yourself")
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

func PrepareCount(sql string, params []any, ds *DataSource) (int64, error) {
	if ds.Db == nil {
		return 0, errors.New("db not allowed to be nil,need to instantiate yourself")
	}
	sql = prepareConvert(sql, ds.DriverName)

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
func Query[T any](sql string, params []any, ds *DataSource) ([]*T, error) {
	if ds.Db == nil {
		return nil, errors.New("db not allowed to be nil,need to instantiate yourself")
	}
	rows, err := ds.Db.Query(sql, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return RowsToStructPtrs[T](rows)
}

// 预处理查询func RowsToStruct[T any](rows *sql.Rows) ([]T, error) {
func PrepareQuery[T any](sql string, params []any, ds *DataSource) ([]*T, error) {
	if ds.Db == nil {
		return nil, errors.New("db not allowed to be nil,need to instantiate yourself")
	}
	sql = prepareConvert(sql, ds.DriverName)
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
	return RowsToStruct[T](rows)
}

// 预处理插入 返回批量自增ID
func PrepareInsert(sql string, params []any, ds *DataSource) (int64, error) {
	if ds.Db == nil {
		return 0, errors.New("db not allowed to be nil,need to instantiate yourself")
	}
	sql = prepareConvert(sql, ds.DriverName)
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

func PrepareUpdate(sql string, params []any, ds *DataSource) (int64, error) {
	if ds.Db == nil {
		return 0, errors.New("db not allowed to be nil,need to instantiate yourself")
	}
	sql = prepareConvert(sql, ds.DriverName)
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
func PrepareSave(sql string, params []any, ds *DataSource) (int64, error) {
	if ds.Db == nil {
		return 0, errors.New("db not allowed to be nil,need to instantiate yourself")
	}
	sql = prepareConvert(sql, ds.DriverName)
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
func PrepareDelete(sql string, params []any, ds *DataSource) (int64, error) {
	if ds.Db == nil {
		return 0, errors.New("db not allowed to be nil,need to instantiate yourself")
	}
	sql = prepareConvert(sql, ds.DriverName)
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
func RowsToStruct[T any](rows *sql.Rows) ([]*T, error) {
	// 创建一个切片的类型，元素为指针类型
	sliceType := reflect.SliceOf(reflect.PtrTo(reflect.TypeOf(new(T)).Elem()))
	// 创建一个具有给定长度和容量的切片
	sliceValue := reflect.MakeSlice(sliceType, 0, 0)

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// 构建字段到列索引的映射
	fieldToColIndex := make(map[string]int)

	structType := reflect.TypeOf(new(T)).Elem()
	if structType.Kind() == reflect.Ptr {
		return nil, errors.New("struct must be a non-pointer type")
	}
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
		// 将 elem 的地址（即指针）添加到切片中
		sliceValue = reflect.Append(sliceValue, reflect.ValueOf(&elem))
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return sliceValue.Interface().([]*T), nil
}
func RowsToStructPtrs[T any](rows *sql.Rows) ([]*T, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	sliceType := reflect.SliceOf(reflect.PtrTo(reflect.TypeOf(new(T)).Elem()))
	sliceValue := reflect.MakeSlice(sliceType, 0, 0)

	structType := reflect.TypeOf(new(T)).Elem()
	if structType.Kind() == reflect.Ptr {
		return nil, errors.New("struct must be a non-pointer type")
	}

	// 递归获取所有字段及其对应的 orm 标签
	fieldToColIndex, err := getAllFieldToColIndex(structType, columns)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		elemPtr := reflect.New(structType)
		elemValue := elemPtr.Elem()

		scanArgs := make([]any, len(columns))
		for i := range scanArgs {
			var temp interface{}
			scanArgs[i] = &temp
		}

		for columnName, colIndex := range fieldToColIndex {
			field := elemValue.FieldByName(columnName)
			if !field.IsValid() || !field.CanAddr() {
				return nil, fmt.Errorf("field %s not found or not addressable in type %T", columnName, elemValue.Interface())
			}
			scanArgs[colIndex] = field.Addr().Interface()
		}

		if err := rows.Scan(scanArgs...); err != nil {
			return nil, err
		}

		sliceValue = reflect.Append(sliceValue, elemPtr)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return sliceValue.Interface().([]*T), nil
}

// 递归获取所有字段及其对应的 orm 标签
func getAllFieldToColIndex(structType reflect.Type, columns []string) (map[string]int, error) {
	fieldToColIndex := make(map[string]int)

	for i, columnName := range columns {
		var found bool
		err := forEachField(structType, func(field reflect.StructField) error {
			tagValue := field.Tag.Get("orm")
			if tagValue == columnName {
				fieldToColIndex[field.Name] = i
				found = true
				return errors.New("found") // 用于退出循环
			}
			return nil
		})
		if err != nil && err.Error() != "found" {
			return nil, err
		}
		if !found {
			return nil, fmt.Errorf("column %s not found in struct fields", columnName)
		}
	}

	return fieldToColIndex, nil
}

// 遍历结构体的所有字段，包括嵌套字段
func forEachField(structType reflect.Type, fn func(reflect.StructField) error) error {
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			// 如果是嵌套结构体，递归处理
			if err := forEachField(field.Type, fn); err != nil {
				return err
			}
		} else {
			if err := fn(field); err != nil {
				return err
			}
		}
	}
	return nil
}

// getFieldTagValue 获取结构体字段的tag值
func getFieldTagValue(field reflect.StructField, tagName string) (string, bool) {
	if tag, ok := field.Tag.Lookup(tagName); ok {
		parts := strings.Split(tag, ",")
		return parts[0], true
	}
	return "", false
}
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
