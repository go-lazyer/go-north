// Create by go-generator  2024-05:06 17:39:04.024
package dao

import (
	"database/sql"
	"fmt"
	dbutil "github.com/go-lazyer/go-generator/db"
	generator "github.com/go-lazyer/go-generator/sql"
	"github.com/pkg/errors"
	"go-generator/test/model"
)

func getDatabase() *sql.DB {
	var database *sql.DB
	if database == nil {
		fmt.Println("db not allowed to be nil,need to instantiate yourself")
	}
	return database
}

// query first by primaryKey
func QueryByPrimaryKey(userId any) (*model.TestModel, error) {
	query := generator.NewEqualQuery(model.USER_ID, userId)

	gen := generator.NewGenerator().Table(model.TABLE_NAME).Where(query)
	sqlStr, params, err := gen.SelectSql(true)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	return QueryFirstBySql(sqlStr, params)
}

// query first by gen
func QueryFirstByGen(gen *generator.Generator) (*model.TestModel, error) {
	sqlStr, params, err := gen.SelectSql(true)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	return QueryFirstBySql(sqlStr, params)
}

// query first by sql
func QueryFirstBySql(sqlStr string, params []any) (*model.TestModel, error) {
	models, err := QueryBySql(sqlStr, params)

	if models == nil || len(models) == 0 || err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	return &(models[0]), nil
}

// query map by primaryKeys
func QueryMapByPrimaryKeys(primaryKeys []any) (map[string]model.TestModel, error) {
	gen := generator.NewGenerator().Table(model.TABLE_NAME).Where(generator.NewInQuery(model.USER_ID, primaryKeys))
	sqlStr, params, err := gen.SelectSql(true)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	return QueryMapBySql(sqlStr, params)
}

// query map by gen
func QueryMapByGen(gen *generator.Generator) (map[string]model.TestModel, error) {
	sqlStr, params, err := gen.SelectSql(true)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	return QueryMapBySql(sqlStr, params)
}

// query map by sql
func QueryMapBySql(sqlStr string, params []any) (map[string]model.TestModel, error) {
	tests := make([]model.TestModel, 0)
	err := dbutil.PrepareQuery(sqlStr, params, &tests, getDatabase())
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	if tests == nil || len(tests) == 0 {
		return nil, nil
	}
	testMap := make(map[string]model.TestModel, len(tests))
	for _, test := range tests {
		testMap[test.UserId.String] = test
	}
	return testMap, nil
}

// count by gen
func CountByGen(gen *generator.Generator) (int64, error) {
	sqlStr, params, err := gen.CountSql(true)
	if err != nil {
		err = errors.WithStack(err)
		return 0, err
	}
	return CountBySql(sqlStr, params)

}

// count by gen
func CountBySql(sqlStr string, params []any) (int64, error) {
	count, err := dbutil.PrepareCount(sqlStr, params, getDatabase())
	if err != nil {
		err = errors.WithStack(err)
		return 0, err
	}
	return count, nil
}

// query by gen
func QueryByGen(gen *generator.Generator) ([]model.TestModel, error) {
	sqlStr, params, err := gen.SelectSql(true)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	return QueryBySql(sqlStr, params)
}

// query by sql
func QueryBySql(sqlStr string, params []any) ([]model.TestModel, error) {
	tests := make([]model.TestModel, 0)
	err := dbutil.PrepareQuery(sqlStr, params, &tests, getDatabase())
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	return tests, nil
}

// query extend by gen
func QueryExtendByGen(gen *generator.Generator) ([]model.TestExtend, error) {
	sqlStr, params, err := gen.SelectSql(true)
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	return QueryExtendBySql(sqlStr, params)
}

// query extend by sql
func QueryExtendBySql(sqlStr string, params []any) ([]model.TestExtend, error) {
	testExtends := make([]model.TestExtend, 0)
	err := dbutil.PrepareQuery(sqlStr, params, &testExtends, getDatabase())
	if err != nil {
		err = errors.WithStack(err)
		return nil, err
	}
	return testExtends, nil
}

func Insert(m *model.TestModel) (int64, error) {
	gen := generator.NewGenerator().Table(model.TABLE_NAME).Insert(m.ToMap(false))
	return InsertByGen(gen)
}

func InsertByGen(gen *generator.Generator) (int64, error) {
	sqlStr, params, err := gen.InsertSql(true)
	if err != nil {
		err = errors.WithStack(err)
		return 0, err
	}
	return InsertBySql(sqlStr, params)
}

func InsertBySql(sqlStr string, params []any) (int64, error) {
	id, err := dbutil.PrepareInsert(sqlStr, params, getDatabase())
	if err != nil {
		err = errors.WithStack(err)
		return 0, err
	}
	return id, nil
}

//batch insert
func InsertByMaps(insertMaps []map[string]any) (int64, error) {
	gen := generator.NewGenerator().Table(model.TABLE_NAME).Inserts(insertMaps)
	sqlStr, params, err := gen.InsertSql(true)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return InsertBySql(sqlStr, params)
}
func Update(m *model.TestModel) (int64, error) {
	query := generator.NewEqualQuery(model.USER_ID, m.UserId.String)
	gen := generator.NewGenerator().Table(model.TABLE_NAME).Update(m.ToMap(false)).Where(query)
	return UpdateByGen(gen)
}

func UpdateByGen(gen *generator.Generator) (int64, error) {
	sqlStr, params, err := gen.UpdateSql(true)
	if err != nil {
		err = errors.WithStack(err)
		return 0, err
	}
	return UpdateBySql(sqlStr, params)
}

func UpdateBySql(sqlStr string, params []any) (int64, error) {
	count, err := dbutil.PrepareUpdate(sqlStr, params, getDatabase())
	if err != nil {
		err = errors.WithStack(err)
		return 0, err
	}
	return count, nil
}

// 批量更新，updateMaps中必须包含主键，联合主键的表不适应该方法
func UpdateByMaps(updateMaps []map[string]any) (int64, error) {
	if updateMaps == nil || len(updateMaps) == 0 {
		return 0, nil
	}
	ids := make([]any, 0)
	for _, updateMap := range updateMaps {
		if value, ok := updateMap[model.USER_ID]; ok {
			ids = append(ids, value)
		}
	}
	if ids == nil || len(ids) == 0 {
		return 0, errors.New("batch update primary not allowed to be nil")
	}
	query := generator.NewInQuery(model.USER_ID, ids)
	gen := generator.NewGenerator().Primary(model.USER_ID).Table(model.TABLE_NAME).Where(query).Updates(updateMaps)
	sqlStr, params, err := gen.UpdateSql(true)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return UpdateBySql(sqlStr, params)
}

func DeleteByPrimaryKey(userId any) (int64, error) {
	gen := generator.NewGenerator().Table(model.TABLE_NAME).Where(generator.NewEqualQuery(model.USER_ID, userId))
	sqlStr, params, err := gen.DeleteSql(true)
	if err != nil {
		err = errors.WithStack(err)
		return 0, err
	}
	return DeleteBySql(sqlStr, params)
}
func DeleteByPrimaryKeys(primaryKeys []any) (int64, error) {
	gen := generator.NewGenerator().Table(model.TABLE_NAME).Where(generator.NewInQuery(model.USER_ID, primaryKeys))
	sqlStr, params, err := gen.DeleteSql(true)
	if err != nil {
		err = errors.WithStack(err)
		return 0, err
	}
	return DeleteBySql(sqlStr, params)
}
func DeleteByGen(gen *generator.Generator) (int64, error) {
	sqlStr, params, err := gen.DeleteSql(true)
	if err != nil {
		err = errors.WithStack(err)
		return 0, err
	}
	return DeleteBySql(sqlStr, params)
}
func DeleteBySql(sqlStr string, params []any) (int64, error) {
	count, err := dbutil.PrepareDelete(sqlStr, params, getDatabase())
	if err != nil {
		err = errors.WithStack(err)
		return 0, err
	}
	return count, nil
}
