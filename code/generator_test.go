package generator

import (
	"fmt"
	"testing"
	// _ "github.com/go-sql-driver/mysql"
)

func TestGenModel(t *testing.T) {
	dsn := "root:xxxx@tcp(test.daoway.cn:3306)/daowei?charset=utf8mb4&parseTime=true&loc=Local"
	var tables = []Module{
		// { //最小配置
		// 	TableName:  "user",
		// 	ModulePath: "/Users/Hch/Workspace/test/api/user",
		// },
		{ //完整配置
			TableName:             "test",                                      //表名
			ModulePath:            "/Users/Lazyer/Workspace/go-generator/test", //相对路径，包含项目名
			Model:                 true,                                        //是否生成Model层代码
			ModelPackageName:      "model",                                     //Model层包名
			ModelFileName:         "user_model.go",                             //Model层文件名
			Extend:                true,                                        //是否生成层代码
			ExtendPackageName:     "extend",                                    //Extend层包名
			ExtendFileName:        "user_extend.go",                            //Extend层文件名
			Param:                 true,                                        //是否生成Param层代码
			ParamPackageName:      "param",                                     //Param层包名
			ParamFileName:         "user_param.go",                             //Param层文件名
			View:                  true,                                        //是否生成View层代码
			ViewPackageName:       "view",                                      //View层包名
			ViewFileName:          "user_view.go",                              //View层文件名
			Dao:                   true,                                        //是否生成Dao层代码
			DaoPackageName:        "dao",                                       //Dao层包名
			DaoFileName:           "user_dao.go",                               //Dao层文件名
			Service:               true,                                        //是否生成Service层代码
			ServicePackageName:    "service",                                   //Service层包名
			ServiceFileName:       "user_service.go",                           //Service层文件名
			Controller:            true,                                        //是否生成Controller层代码
			ControllerPackageName: "controller",                                //Controller层包名
			ControllerFileName:    "user_controller.go",                        //Controller层文件名
		},
	}

	fmt.Println(NewGenerator().Dsn(dsn).Project("go-generator").Gen(tables))
}
