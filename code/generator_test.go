package generator

import (
	"testing"
)

func TestGenModel(t *testing.T) {
	dsn := "root:Daoway_Mysql_iO12@tcp(analysis.daoway.cn:3306)/daowei?charset=utf8mb4&parseTime=true&loc=Local"
	var tables = []Module{
		{
			TableName:  "city",
			ModulePath: "/Users/Hch/Workspace/lazyer/api/city", //相对路径，包含项目名
		},
		// {
		// 	TableName:  "technician",
		// 	ModulePath: "/Users/Hch/Workspace/analysis/api/technician", //相对路径，包含项目名
		// },
	}

	NewGenerator().Dsn(dsn).Project("lazyer").Run(tables)
}
