package north

import (
	"database/sql"
	"errors"
	"fmt"
	"maps"
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
	Tx           *sql.Tx
	DriverName   string
	DaoFilePaths []string
}
type Config struct {
	MaxOpenConns int
	MaxIdleConns int
}

func Open(driverName string, dsn string, config Config) (DataSource, error) {
	db, err := sql.Open(driverName, dsn)
	if err != nil {
		return DataSource{}, err
	}
	err = db.Ping()
	if err != nil {
		return DataSource{}, err
	}
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)

	return DataSource{
		Db:         db,
		DriverName: driverName,
	}, nil
}

func Count(sql string, params []any, ds DataSource) (int64, error) {
	if ds.Db == nil {
		return 0, errors.New("db not allowed to be nil,need to instantiate yourself")
	}
	sql = prepareConvert(sql, ds.DriverName)
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

func PrepareCount(sql string, params []any, ds DataSource) (int64, error) {
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
func Query[T any](sql string, params []any, ds DataSource) ([]T, error) {
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
func PrepareQuery[T any](sql string, params []any, ds DataSource) ([]T, error) {
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
func PrepareInsert(sqlStr string, params []any, ds DataSource) (int64, error) {
	if ds.Db == nil && ds.Tx == nil {
		return 0, errors.New("db and tx not allowed to be all nil,need to instantiate yourself")
	}
	sqlStr = prepareConvert(sqlStr, ds.DriverName)

	var stmt *sql.Stmt
	var err error
	if ds.Tx != nil {
		stmt, err = ds.Tx.Prepare(sqlStr)
	}
	if ds.Db != nil {
		stmt, err = ds.Db.Prepare(sqlStr)
	}
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

func PrepareUpdate(sqlStr string, params []any, ds DataSource) (int64, error) {
	if ds.Db == nil && ds.Tx == nil {
		return 0, errors.New("db and tx not allowed to be all nil,need to instantiate yourself")
	}
	sqlStr = prepareConvert(sqlStr, ds.DriverName)

	var stmt *sql.Stmt
	var err error
	if ds.Tx != nil {
		stmt, err = ds.Tx.Prepare(sqlStr)
	}
	if ds.Db != nil {
		stmt, err = ds.Db.Prepare(sqlStr)
	}
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	ret, err := stmt.Exec(params...)

	if ds.Db != nil {
		ret, err = ds.Db.Exec(sqlStr, params...)
	}
	if err != nil {
		return 0, err
	}
	n, err := ret.RowsAffected() // 操作影响的行数
	if err != nil {
		return 0, err
	}
	return n, nil
}
func PrepareDelete(sqlStr string, params []any, ds DataSource) (int64, error) {
	if ds.Db == nil && ds.Tx == nil {
		return 0, errors.New("db and tx not allowed to be all nil,need to instantiate yourself")
	}
	sqlStr = prepareConvert(sqlStr, ds.DriverName)

	var stmt *sql.Stmt
	var err error
	if ds.Tx != nil {
		stmt, err = ds.Tx.Prepare(sqlStr)
	}
	if ds.Db != nil {
		stmt, err = ds.Db.Prepare(sqlStr)
	}
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	ret, err := stmt.Exec(params...)

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
func RowsToStruct[T any](rows *sql.Rows) ([]T, error) {
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	// 获取类型 T 的反射信息
	structType := reflect.TypeOf(new(T)).Elem()
	if structType.Kind() == reflect.Ptr {
		return nil, errors.New("t must be a non-pointer type")
	}

	// 创建值类型的切片（[]T）
	sliceType := reflect.SliceOf(structType)
	sliceValue := reflect.MakeSlice(sliceType, 0, 0)

	// 递归获取字段与列的映射（此处省略具体实现）

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

		for fieldName, index := range fieldToColIndex {
			field := elemValue.FieldByName(fieldName)
			if !field.IsValid() || !field.CanAddr() {
				return nil, fmt.Errorf("field %s not found or not addressable in type %T", fieldName, elemValue.Interface())
			}
			scanArgs[index] = field.Addr().Interface()
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

// 获取表字段和 struct 字段的并集
func getAllFieldToColIndex(structType reflect.Type, columns []string) (map[string]int, error) {
	fieldMap := forEachField(structType)
	fieldToColIndex := make(map[string]int)
	for i, columnName := range columns {
		if fieldName, ok := fieldMap[columnName]; ok {
			fieldToColIndex[fieldName] = i
		} else {
			fmt.Printf("table column %s not found in struct fields\n", columnName)
		}
	}
	return fieldToColIndex, nil
}

// 遍历结构体的所有字段和tag 中的 orm 标签，包括继承的 struct
func forEachField(structType reflect.Type) map[string]string {
	fields := make(map[string]string, 0)

	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		if field.Anonymous && field.Type.Kind() == reflect.Struct {
			// 如果是嵌套结构体，递归处理
			maps.Copy(fields, forEachField(field.Type))
		} else {
			tagValue := field.Tag.Get("orm")
			if tagValue != "" {
				fields[field.Tag.Get("orm")] = field.Name
			}
		}
	}
	return fields
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
func (ds DataSource) Transaction(fc func(tx *sql.Tx) error) (err error) {
	tx, err := ds.Db.Begin()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	if err = fc(tx); err == nil {
		return tx.Commit()
	}
	return
}
