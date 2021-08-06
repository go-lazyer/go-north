

# go-generator

### 一、安装教程

```
go get github.com/go-lazyer/go-generator
```
### 二、使用说明
go-generator sql生成和struct生成的工具，减少开发者自行拼接sql，和table to struct的工作量，可以专注于核心业务代码，该工具是模仿 mybatis-generator 实现的，同时又借鉴了go elastic工具的包github.com/olivere/elastic 的查询实现方式，所以熟悉 mybatis-generator和olivere 及容易上手。

go-generator 分为两个模块，sql-generator（生成sql）和code-gengrator(生成代码)，两个模块都可以单独使用，也可以配合使用

注意：该工具只支持mysql

### 三、sql-generator

sql-generator 可以生成普通sql和预处理sql，配合golang官方提供的sql.DB,可以轻松实现增删改查.

1. 引入

```go
import generator "github.com/go-lazyer/go-generator/sql"
```



1.   统计
``` go
//select count(1) count from user where t.id>1000

query := generator.NewGreaterThanQuery("id", 1000)

gen := generator.NewGenerator().Table("user").Where(query)

fmt.Println(gen.CountSql(false))
  
```
2. 基础查询

```go
//select * from user

gen := NewGenerator().Table("user")

fmt.Println(gen.SelectSql(false))

//select * from user where t.id=1000

query := NewEqualQuery("id", 1000)

gen := NewGenerator().Table("user").Where(query)

fmt.Println(gen.SelectSql(false))
```

3. 排序查询

```go
// select * from user where id=1000 and age>20 order by age desc,id asc

idQuery := NewEqualQuery("id", 1000)

ageQuery := NewGreaterThanQuery("age", 20)

boolQuery := NewBoolQuery().And(idQuery, ageQuery)

gen := NewGenerator().Table("user").Where(boolQuery).AddOrderBy("age", "desc").AddOrderBy("id", "asc")

fmt.Println(gen.SelectSql(false))
```

4. 复杂查询

```go
// select id,name,age from user where (id=1000 and age>20) or age <=10 order by age desc

idQuery := NewEqualQuery("id", 1000)

ageQuery := NewGreaterThanQuery("age", 20)

boolQuery := NewBoolQuery().And(idQuery, ageQuery)

ageQuery2 := NewLessThanOrEqualQuery("age", 10)

gen := NewGenerator().Result("id", "name", "age").Table("user").Where(boolQuery, ageQuery2).AddOrderBy("age", "desc")

fmt.Println(gen.SelectSql(false))
```

5. 联表查询

 ```go
 // select user.id,order.id  from user join order on user.id=order.user_id where user.id='10000'
 
 idQuery = NewEqualQuery("id", 1000)
 
 join := NewJoin("order", INNER_JOIN).Condition("user", "id", "order", "user_id")
 
 gen = NewGenerator().Result("user.id", "order.id").Table("user").Join(join).Where(idQuery)
 
 fmt.Println(gen.SelectSql(false))
 ```

6. 更新

```go
// update user set age=21,name="lazeyr" where id="10000"	

query := NewEqualQuery("id", 1000)

set := map[string]interface{}{
  "age":  21,
  "name": "lazyer",
}

gen := NewGenerator().Table("user").Where(query).Set(set)

fmt.Println(gen.UpdateSql(false))
```

7. 批量更新

```go
// update `user`
// set
// 	name = case dwid
// 		when 10001 then boy
// 		when 10002 then boy
// 		when 10003 then girl
// 		end,
// 	age = case dwid
// 		when 10001 then 10
// 		when 10002 then 20
// 		when 10003 then 30
// 	end,
// 	name = case dwid
// 		when 10001 then lilie
// 		when 10002 then lining
// 		when 10003 then hanmeimei
// 	end
// where
// 		user.dwid in('10001', '10002', '10003')

f1 := map[string]interface{}{
  "name": "lilie",
  "sex":  "boy",
  "age":  "10",
}
f2 := map[string]interface{}{
  "name": "lining",
  "sex":  "boy",
  "age":  "20",
}
f3 := map[string]interface{}{
  "name": "hanmeimei",
  "sex":  "girl",
  "age":  "30",
}
set := map[interface{}]map[string]interface{}{
  "10001": f1,
  "10002": f2,
  "10003": f3,
}

dwids := []interface{}{
  10001, 10002, 10003,
}

query := NewInQuery("dwid", dwids)

gen := NewGenerator().Table("user").Where(query).Primary("dwid").Updates(set)

fmt.Print(gen.BatchUpdateSql(false))
```

8. 批量插入

```go
//insert into `user` ( age , name , sex ) 
//values
//( '10' , 'lilie' , 'boy' ),
//( '20' , 'lining' , 'boy' ),
//( '30' , 'hanmeimei' , 'girl' )
f1 := map[string]interface{}{
  "name": "lilie",
  "sex":  "boy",
  "age":  "10",
}
f2 := map[string]interface{}{
  "name": "lining",
  "sex":  "boy",
  "age":  "20",
}
f3 := map[string]interface{}{
  "name": "hanmeimei",
  "sex":  "girl",
  "age":  "30",
}

dwids := []map[string]interface{}{
  f1, f2, f3,
}

gen := NewGenerator().Table("user").Inserts(dwids)

fmt.Print(gen.BatchInsertSql(false))
```

### 四、code-gengrator

code-gengrator 模块主要用于生成数据库表对应的struct，以及dao文件，同时会生成相关的附属类文件

#### 文件介绍

1. model 文件

```go
// Create by code generator  2021-08:06 16:27:45.021
package model

import (
	"bytes"
	"database/sql"
	"errors"
)

const (
	ID   = "id"   // 编号
	NAME = "name" // 姓名
	AGE  = "age"  // 年龄
	SEX  = "sex"  // 性别

	TABLE_NAME = "test" // 表名
)

type TestModel struct {
	Id   sql.NullInt32  `orm:"id" default:"0"` // 编号
	Name sql.NullString `orm:"name" `          // 姓名
	Age  sql.NullInt32  `orm:"age" `           // 年龄
	Sex  sql.NullString `orm:"sex" `           // 性别

}

func (m *TestModel) UpdateSql() (string, []interface{}, error) {

	if !m.Id.Valid {
		return "", nil, errors.New("Id is not null")
	}

	params := make([]interface{}, 0, 4)
	var sql bytes.Buffer
	sql.WriteString("update `test` ")
	sql.WriteString("set `name` = ?,`age` = ?,`sex` = ? ")

	nameV, err := m.Name.Value()
	if nameV == nil || err != nil {
		params = append(params, nil)
	} else {
		params = append(params, nameV)
	}

	ageV, err := m.Age.Value()
	if ageV == nil || err != nil {
		params = append(params, nil)
	} else {
		params = append(params, ageV)
	}

	sexV, err := m.Sex.Value()
	if sexV == nil || err != nil {
		params = append(params, nil)
	} else {
		params = append(params, sexV)
	}

	sql.WriteString(" where  `id` = ? ")
	params = append(params, m.Id.Int32)
	return sql.String(), params, nil
}

func (m *TestModel) UpdateSqlBySelective() (string, []interface{}, error) {
	if !m.Id.Valid {
		return "", nil, errors.New("Id is not null")
	}

	params := make([]interface{}, 0, 4)
	var sql bytes.Buffer
	sql.WriteString("update `test` ")

	sql.WriteString(" set ")

	if m.Id.Valid {
		sql.WriteString(" `id` = ? ")
		params = append(params, m.Id.Int32)
	}

	if m.Name.Valid {
		sql.WriteString(", `name` = ? ")
		params = append(params, m.Name.String)
	}

	if m.Age.Valid {
		sql.WriteString(", `age` = ? ")
		params = append(params, m.Age.Int32)
	}

	if m.Sex.Valid {
		sql.WriteString(", `sex` = ? ")
		params = append(params, m.Sex.String)
	}

	sql.WriteString(" where  `id` = ? ")

	params = append(params, m.Id.Int32)
	return sql.String(), params, nil
}

func (m *TestModel) InsertSql() (string, []interface{}, error) {
	params := make([]interface{}, 0, 4)
	var sql bytes.Buffer
	sql.WriteString("insert into `test` ")
	sql.WriteString(" ( `id` ,`name` ,`age` ,`sex`)")
	sql.WriteString("values ( ? ,? ,? ,?)")

	idV, err := m.Id.Value()
	if idV == nil || err != nil {
		params = append(params, "0")
	} else {
		params = append(params, idV)
	}

	nameV, err := m.Name.Value()
	if nameV == nil || err != nil {
		params = append(params, nil)
	} else {
		params = append(params, nameV)
	}

	ageV, err := m.Age.Value()
	if ageV == nil || err != nil {
		params = append(params, nil)
	} else {
		params = append(params, ageV)
	}

	sexV, err := m.Sex.Value()
	if sexV == nil || err != nil {
		params = append(params, nil)
	} else {
		params = append(params, sexV)
	}

	return sql.String(), params, nil
}

```



#### 使用教程

新建main方法，配置需要生成代码的表(支持配置多个)，运行即可生成代码

```go
package main

import (
	generator "github.com/go-lazyer/go-generator/code"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "root:123@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=true&loc=Local"
	var tables = []Module{
		{//最小配置
		 	TableName:  "user",
		 	ModulePath: "/Users/Hch/Workspace/lazyer/api/user",
		},
    { //完整配置
			TableName:             "user",                                 //表名
			ModulePath:            "/Users/Hch/Workspace/lazyer/api/user", //相对路径，包含项目名
			Model:                 true,                                   //是否生成Model层代码
			ModelPackageName:      "model",                                //Model层包名
			ModelFileName:         "user_model.go",                        //Model层文件名
			Extend:                true,                                   //是否生成层代码
			ExtendPackageName:     "extend",                               //Extend包名
			ExtendFileName:        "user_extend.go",                       //Extend文件名
			Param:                 true,                                   //是否生成Param代码
			ParamPackageName:      "param",                                //Param包名
			ParamFileName:         "user_param.go",                        //Param文件名
			View:                  true,                                   //是否生成View代码
			ViewPackageName:       "view",                                 //View包名
			ViewFileName:          "user_view.go",                         //View文件名
			Dao:                   true,                                   //是否生成Dao代码
			DaoPackageName:        "dao",                                  //Dao层包名
			DaoFileName:           "user_dao.go",                          //Dao层文件名
			Service:               true,                                   //是否生成Service层代码
			ServicePackageName:    "service",                              //Service层包名
			ServiceFileName:       "user_service.go",                      //Service层文件名
			Controller:            true,                                   //是否生成Controller层代码
			ControllerPackageName: "controller",                           //Controller层包名
			ControllerFileName:    "user_controller.go",                   //Controller层文件名
		},
	}

	generator.NewGenerator().Dsn(dsn).Project("lazyer").Gen(tables)
}

```

