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
	return RowsToStruct[T](rows)
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
