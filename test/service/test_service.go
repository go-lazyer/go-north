// Create by code generator  2024-05:06 17:39:04.024
package service

import (
	"go-generator/test/dao"
	"go-generator/test/model"
	"go-generator/test/param"

	generator "github.com/go-lazyer/go-generator/sql"
)

func QueryByPrimaryKey(userId any) (*model.TestModel, error) {
	test, err := dao.QueryByPrimaryKey(userId)
	if err != nil {
		return nil, err
	}
	return test, nil
}

func QueryByParam(testParam *param.TestParam) ([]model.TestModel, error) {
	query := generator.NewBoolQuery()
	gen := generator.NewGenerator().PageNum(testParam.PageNum).PageStart(testParam.PageStart).PageSize(testParam.PageSize).Table(model.TABLE_NAME).Where(query)
	tests, err := dao.QueryByGen(gen)
	if err != nil {
		return nil, err
	}
	return tests, nil
}
