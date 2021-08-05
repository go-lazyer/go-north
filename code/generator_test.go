package generator

import (
	"testing"
)

func TestGenModel(t *testing.T) {
	dsn := "root:123@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=true&loc=Local"
	var tables = []Module{
		{
			TableName:  "user",
			ModulePath: "/Users/Hch/Workspace/lazyer/api/user", //相对路径，包含项目名
		},
		// {
		// 	TableName:  "order",
		// 	ModulePath: "/Users/Hch/Workspace/lazyer/api/order", //相对路径，包含项目名
		// },
	}

	NewGenerator().Dsn(dsn).Project("lazyer").Run(tables)
}
