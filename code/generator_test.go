package generator

import (
	"testing"
)

func TestGenModel(t *testing.T) {
	dsn := "root:Daoway_Mysql_iO12@tcp(analysis.daoway.cn:3306)/daowei?charset=utf8mb4&parseTime=true&loc=Local"
	var tables = []Module{
		{
			TableName:  "s_tech_info",
			ModulePath: "/Users/Hch/Workspace/analysis/api/s_tech_info", //相对路径，包含项目名
		},
		// {
		// 	TableName:  "technician",
		// 	ModulePath: "/Users/Hch/Workspace/analysis/api/technician", //相对路径，包含项目名
		// },
	}

	NewGenerator().Dsn(dsn).Project("analysis").Run(tables)
}
