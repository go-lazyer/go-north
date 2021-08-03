package generator

import (
	"bytes"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"time"
)

type Generator struct {
	dsn     string
	project string
}

func NewGenerator() *Generator {
	return &Generator{}
}
func (gen *Generator) Dsn(dsn string) *Generator {
	gen.dsn = dsn
	return gen
}
func (gen *Generator) Project(project string) *Generator {
	gen.project = project
	return gen
}

type Module struct {
	TableName           string //表名
	TableNameUpperCamel string //表名的大驼峰
	TableNameLowerCamel string //表名的小驼峰
	ModulePath          string //模块名用于生成文件名
	Fields              []Field
	PrimaryKeyFields    []Field //主键
	ModelFilePath       string  //全路径，不包含文件名
	ModelFileName       string  //只有文件名
	ModelPackageName    string  //只有包名，不包含文件名
	ModelPackagePath    string  //包含完整的包名

	ExtendFilePath    string //全路径，不包含文件名
	ExtendFileName    string //只有文件名
	ExtendPackageName string //只有包名，不包含文件名
	ExtendPackagePath string //包含完整的包名

	ViewFilePath    string
	ViewFileName    string
	ViewPackageName string
	ViewPackagePath string

	ParamFilePath    string
	ParamFileName    string
	ParamPackageName string
	ParamPackagePath string

	DaoFilePath    string
	DaoFileName    string
	DaoPackageName string
	DaoPackagePath string

	ServiceFilePath    string
	ServiceFileName    string
	ServicePackageName string
	ServicePackagePath string

	ControllerFilePath    string
	ControllerFileName    string
	ControllerPackageName string
	ControllerPackagePath string

	UpdateSql          string
	UpdateSelectiveSql string
	CreateTime         string
}

type Field struct {
	ColumnName           string         //msyql字段名 user_id
	ColumnNameLowerCamel string         //小驼峰 userId
	ColumnNameUpper      string         //字段名大写 USER_ID
	ColumnType           string         //msql 类型 varchat
	ColumnDefault        sql.NullString //默认值
	IsNullable           int            //允许为空
	IsPrimaryKey         int            //是否主键
	FieldName            string         //实体名称 大驼峰  UserId
	FieldNullType        string         //实体golang Null类型 sql.NullString
	FieldNullTypeValue   string         //实体golang Null类型 取值  String
	FieldType            string         //golang 类型  string
	FieldTypeDefault     string         //golang 类型  的默认值
	FieldOrmTag          string         //tag orm:
	FieldJsonTag         string         //tag json
	FieldFormTag         string         //tag form
	FieldDefaultTag      string         //tag 默认值
	Comment              string         //表中字段注释
}

var dbType = map[string]goType{
	"int":                {"int32", "sql.NullInt32", "Int32", "0"},
	"integer":            {"int32", "sql.NullInt32", "Int32", "0"},
	"tinyint":            {"int32", "sql.NullInt32", "Int32", "0"},
	"smallint":           {"int32", "sql.NullInt32", "Int32", "0"},
	"mediumint":          {"int32", "sql.NullInt32", "Int32", "0"},
	"bigint":             {"int32", "sql.NullInt32", "Int32", "0"},
	"int unsigned":       {"int32", "sql.NullInt32", "Int32", "0"},
	"integer unsigned":   {"int32", "sql.NullInt32", "Int32", "0"},
	"tinyint unsigned":   {"int32", "sql.NullInt32", "Int32", "0"},
	"smallint unsigned":  {"int32", "sql.NullInt32", "Int32", "0"},
	"mediumint unsigned": {"int32", "sql.NullInt32", "Int32", "0"},
	"bigint unsigned":    {"int32", "sql.NullInt32", "Int32", "0"},
	"bit":                {"int32", "sql.NullInt32", "Int32", "0"},
	"bool":               {"bool", "sql.NullBool", "Bool", "false"},
	"enum":               {"string", "sql.NullString", "String", "\"\""},
	"set":                {"string", "sql.NullString", "String", "\"\""},
	"varchar":            {"string", "sql.NullString", "String", "\"\""},
	"char":               {"string", "sql.NullString", "String", "\"\""},
	"tinytext":           {"string", "sql.NullString", "String", "\"\""},
	"mediumtext":         {"string", "sql.NullString", "String", "\"\""},
	"text":               {"string", "sql.NullString", "String", "\"\""},
	"longtext":           {"string", "sql.NullString", "String", "\"\""},
	"blob":               {"string", "sql.NullString", "String", "\"\""},
	"tinyblob":           {"string", "sql.NullString", "String", "\"\""},
	"mediumblob":         {"string", "sql.NullString", "String", "\"\""},
	"longblob":           {"string", "sql.NullString", "String", "\"\""},
	"date":               {"time.Time", "sql.NullTime", "Time", "nil"},
	"datetime":           {"time.Time", "sql.NullTime", "Time", "nil"},
	"timestamp":          {"time.Time", "sql.NullTime", "Time", "nil"},
	"time":               {"time.Time", "sql.NullTime", "Time", "nil"},
	"float":              {"float64", "sql.NullFloat64", "Float64", "0"},
	"double":             {"float64", "sql.NullFloat64", "Float64", "0"},
	"decimal":            {"float64", "sql.NullFloat64", "Float64", "0"},
	"binary":             {"string", "sql.NullString", "String", "\"\""},
	"varbinary":          {"string", "sql.NullString", "String", "\"\""},
}

type goType struct {
	baseType      string
	nullType      string
	nullTypeValue string
	defaultValue  string
}

func (gen *Generator) dialMysql() *sql.DB {
	if gen.dsn == "" {
		panic("dsn数据库配置缺失")
	}
	db, err := sql.Open("mysql", gen.dsn)
	if err != nil {
		panic(err)
	}
	return db
}

func getFields(tableName string, db *sql.DB) ([]Field, []Field, error) {
	var sqlStr = `select
					column_name name,
					data_type type,
					if('YES'=is_nullable,true,false) isNullable,
					if('PRI'=column_key,true,false) isPrimaryKey,
					column_comment comment,column_default 'default'
				from
					information_schema.COLUMNS 
				where
					table_schema = DATABASE() `
	sqlStr += fmt.Sprintf(" and table_name = '%s' order by isPrimaryKey desc", tableName)
	rows, err := db.Query(sqlStr)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	fields := make([]Field, 0)
	primaryKeyFields := make([]Field, 0)
	for rows.Next() {
		field := Field{}
		err = rows.Scan(&field.ColumnName, &field.ColumnType, &field.IsNullable, &field.IsPrimaryKey, &field.Comment, &field.ColumnDefault)
		if err != nil {
			panic(err)
		}
		field.FieldName = ToUpperCamelCase(field.ColumnName)
		field.ColumnNameLowerCamel = ToLowerCamelCase(field.ColumnName)
		field.ColumnNameUpper = strings.ToUpper(field.ColumnName)
		field.FieldType = dbType[field.ColumnType].baseType
		field.FieldTypeDefault = dbType[field.ColumnType].defaultValue
		field.FieldNullType = dbType[field.ColumnType].nullType
		field.FieldNullTypeValue = dbType[field.ColumnType].nullTypeValue
		field.FieldOrmTag = fmt.Sprintf("orm:\"%v\"", field.ColumnName)
		field.FieldJsonTag = fmt.Sprintf("json:\"%v\"", field.ColumnName)
		field.FieldFormTag = fmt.Sprintf("form:\"%v\"", field.ColumnName)
		val, _ := field.ColumnDefault.Value()
		if val != nil {
			field.FieldDefaultTag = fmt.Sprintf("default:\"%v\"", val)
		}
		if field.IsPrimaryKey == 1 {
			primaryKeyFields = append(primaryKeyFields, field)
		}
		fields = append(fields, field)
	}
	return fields, primaryKeyFields, nil
}

func (gen *Generator) Run(modules []Module) error {
	db := gen.dialMysql()

	for _, module := range modules {
		fields, primaryKeyFields, err := getFields(module.TableName, db)
		if err != nil {
			return err
		}
		tableName := strings.ReplaceAll(module.TableName, "p_", "")
		module.CreateTime = time.Now().Format("2006-01:02 15:04:05.006")
		module.Fields = fields
		module.PrimaryKeyFields = primaryKeyFields
		module.TableNameUpperCamel = ToUpperCamelCase(tableName)
		module.TableNameLowerCamel = ToLowerCamelCase(tableName)
		urls := strings.Split(module.ModulePath, gen.project)
		module.ModelPackageName = "model"
		module.ModelPackagePath = gen.project + urls[1] + "/" + module.ModelPackageName
		module.ModelFileName = tableName + "_" + module.ModelPackageName + ".go"
		module.ModelFilePath = module.ModulePath + "/" + module.ModelPackageName

		module.ExtendPackageName = "extend"
		module.ExtendPackagePath = gen.project + urls[1] + "/" + module.ModelPackageName
		module.ExtendFileName = tableName + "_" + module.ExtendPackageName + ".go"
		module.ExtendFilePath = module.ModulePath + "/" + module.ModelPackageName

		module.ViewPackageName = "view"
		module.ViewPackagePath = gen.project + urls[1] + "/" + module.ViewPackageName
		module.ViewFileName = tableName + "_" + module.ViewPackageName + ".go"
		module.ViewFilePath = module.ModulePath + "/" + module.ViewPackageName

		module.ParamPackageName = "param"
		module.ParamPackagePath = gen.project + urls[1] + "/" + module.ParamPackageName
		module.ParamFileName = tableName + "_" + module.ParamPackageName + ".go"
		module.ParamFilePath = module.ModulePath + "/" + module.ParamPackageName

		module.DaoPackageName = "dao"
		module.DaoPackagePath = gen.project + urls[1] + "/" + module.DaoPackageName
		module.DaoFileName = tableName + "_" + module.DaoPackageName + ".go"
		module.DaoFilePath = module.ModulePath + "/" + module.DaoPackageName

		module.ServicePackageName = "service"
		module.ServicePackagePath = gen.project + urls[1] + "/" + module.ServicePackageName
		module.ServiceFileName = tableName + "_" + module.ServicePackageName + ".go"
		module.ServiceFilePath = module.ModulePath + "/" + module.ServicePackageName

		module.ControllerPackageName = "controller"
		module.ControllerPackagePath = gen.project + urls[1] + "/" + module.ControllerPackageName
		module.ControllerFileName = tableName + "_" + module.ControllerPackageName + ".go"
		module.ControllerFilePath = module.ModulePath + "/" + module.ControllerPackageName

		genFile(&module, module.ModelPackageName)
		genFile(&module, module.ExtendPackageName)
		genFile(&module, module.ViewPackageName)
		genFile(&module, module.ParamPackageName)
		genFile(&module, module.DaoPackageName)
		genFile(&module, module.ServicePackageName)
		genFile(&module, module.ControllerPackageName)

	}
	return nil
}
func genSql(module *Module) {
	var sql bytes.Buffer
	var updateSelectiveSql bytes.Buffer
	sql.WriteString("update " + module.TableName + " set")
	updateSelectiveSql.WriteString("update " + module.TableName + " set")

	for i, field := range module.Fields {
		if field.IsPrimaryKey == 1 {
			continue
		}
		sql.WriteString(" `" + field.ColumnName + "` = ?")
		if i != len(module.Fields)-1 {
			sql.WriteString(",")
		}
	}
	sql.WriteString(" where ")
	for i, field := range module.Fields {
		if field.IsPrimaryKey != 1 {
			continue
		}
		if i != 0 {
			sql.WriteString(" and ")
			updateSelectiveSql.WriteString(" and ")
		}
		sql.WriteString("`" + field.ColumnName + "` = ?")
		updateSelectiveSql.WriteString("`" + field.ColumnName + "` = ?")
	}
	module.UpdateSql = sql.String()
	module.UpdateSelectiveSql = updateSelectiveSql.String()
}
func genFile(table *Module, packageName string) {

	var templateFile, filePath, file string
	if "model" == packageName {
		templateFile = "./template/model.template"
		filePath = table.ModelFilePath
		file = filePath + "/" + table.ModelFileName
	} else if "extend" == packageName {
		templateFile = "./template/extend.template"
		filePath = table.ExtendFilePath
		file = filePath + "/" + table.ExtendFileName
		if IsExist(file) { //extend 不覆盖
			return
		}
	} else if "view" == packageName {
		templateFile = "./template/view.template"
		filePath = table.ViewFilePath
		file = filePath + "/" + table.ViewFileName
		if IsExist(file) { //view 不覆盖
			return
		}
	} else if "param" == packageName {
		templateFile = "./template/param.template"
		filePath = table.ParamFilePath
		file = filePath + "/" + table.ParamFileName
		if IsExist(file) { //param 不覆盖
			return
		}
	} else if "dao" == packageName {
		templateFile = "./template/dao.template"
		filePath = table.DaoFilePath
		file = filePath + "/" + table.DaoFileName
		if IsExist(file) { //dao 不覆盖
			return
		}
	} else if "service" == packageName {
		templateFile = "./template/service.template"
		filePath = table.ServiceFilePath
		file = filePath + "/" + table.ServiceFileName
		if IsExist(file) { //service 不覆盖
			return
		}
	} else if "controller" == packageName {
		templateFile = "./template/controller.template"
		filePath = table.ControllerFilePath
		file = filePath + "/" + table.ControllerFileName
		if IsExist(file) { //controller 不覆盖
			return
		}
	}

	// 第一步，加载模版文件
	tmpl, err := template.ParseFiles(templateFile)
	if err != nil {
		fmt.Println("create template model, err:", err)
		return
	}
	// 第二步，创建文件目录
	err = CreateDir(filePath)
	if err != nil {
		fmt.Printf("create path:%v err", filePath)
		return
	}
	// 第三步，创建且打开文件
	f, err := os.OpenFile(file, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Can not write file")
		return
	}
	defer f.Close()

	// 第四步，写入数据
	tmpl.Execute(f, table)

	//第五步，格式化代码
	cmd := exec.Command("gofmt", "-w", file)
	cmd.Run()
}
