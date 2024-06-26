

# go-north

### 一、安装教程

```
go get github.com/go-lazyer/go-north
```
### 二、使用说明
go-north sql生成和struct生成的工具，减少开发者自行拼接sql，和table to struct的工作量，可以专注于核心业务代码，该工具是模仿 mybatis-generator 实现的，同时又借鉴了go elastic工具的包github.com/olivere/elastic 的查询实现方式，所以熟悉 mybatis-generator和olivere 及容易上手。

go-north 分为两个模块，sql-generator（生成sql）和code-gengrator(生成代码)，两个模块都可以单独使用，也可以配合使用

注意：该工具只支持mysql

### 三、sql-generator

sql-generator 可以生成普通sql和预处理sql，配合golang官方提供的sql.DB,可以轻松实现增删改查.

#### 1、引入

```go
import generator "github.com/go-lazyer/go-north/sql"
```



#### 2、统计

``` go
//select count(1) count from user where t.id>1000

query := generator.NewGreaterThanQuery("id", 1000)

gen := generator.NewGenerator().Table("user").Where(query)

fmt.Println(gen.CountSql(false))
  
```
#### 3、基础查询

```go
//select * from user

gen := NewGenerator().Table("user")

fmt.Println(gen.SelectSql(false))

//select * from user where t.id=1000

query := NewEqualQuery("id", 1000)

gen := NewGenerator().Table("user").Where(query)

fmt.Println(gen.SelectSql(false))
```

#### 4、排序查询

```go
// select * from user where id=1000 and age>20 order by age desc,id asc

idQuery := NewEqualQuery("id", 1000)

ageQuery := NewGreaterThanQuery("age", 20)

boolQuery := NewBoolQuery().And(idQuery, ageQuery)

gen := NewGenerator().Table("user").Where(boolQuery).AddOrderBy("age", "desc").AddOrderBy("id", "asc")

fmt.Println(gen.SelectSql(false))
```

#### 5、复杂查询

```go
// select id,name,age from user where (id=1000 and age>20) or age <=10 order by age desc

idQuery := NewEqualQuery("id", 1000)

ageQuery := NewGreaterThanQuery("age", 20)

boolQuery := NewBoolQuery().And(idQuery, ageQuery)

ageQuery2 := NewLessThanOrEqualQuery("age", 10)

gen := NewGenerator().Result("id", "name", "age").Table("user").Where(boolQuery, ageQuery2).AddOrderBy("age", "desc")

fmt.Println(gen.SelectSql(false))
```

#### 6、联表查询

 ```go
 // select user.id,order.id  from user join order on user.id=order.user_id where user.id='10000'
 
 idQuery = NewEqualQuery("id", 1000)
 
 join := NewJoin("order", INNER_JOIN).Condition("user", "id", "order", "user_id")
 
 gen = NewGenerator().Result("user.id", "order.id").Table("user").Join(join).Where(idQuery)
 
 fmt.Println(gen.SelectSql(false))
 ```

#### 7、更新

```go
// update user set age=21,name="lazeyr" where id="10000"	

query := NewEqualQuery("id", 1000)

set := map[string]any{
  "age":  21,
  "name": "lazyer",
}

gen := NewGenerator().Table("user").Where(query).Update(set)

fmt.Println(gen.UpdateSql(false))
```

#### 8、批量更新(只支持主键更新)

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
// 	user.dwid in('10001', '10002', '10003')

f1 := map[string]any{
  "name": "lilie",
  "sex":  "boy",
  "age":  "10",
}
f2 := map[string]any{
  "name": "lining",
  "sex":  "boy",
  "age":  "20",
}
f3 := map[string]any{
  "name": "hanmeimei",
  "sex":  "girl",
  "age":  "30",
}
set := map[any]map[string]any{
  "10001": f1,
  "10002": f2,
  "10003": f3,
}

dwids := []any{
  10001, 10002, 10003,
}

query := NewInQuery("dwid", dwids)

gen := NewGenerator().Table("user").Where(query).Primary("dwid").Updates(set)

fmt.Print(gen.UpdateSql(false))
```

#### 9、单条插入

```go
// insert into `user` ( age , name , sex ) values ( '10' , 'lilie' , 'boy' ),

m := map[string]any{
  "name": "lilie",
  "sex":  "boy",
  "age":  "10",
}

gen := generator.NewGenerator().Table(model.TABLE_NAME).Insert(m)

fmt.Println(gen.UpdateSql(false))
```

#### 10、批量插入

```go
//insert into `user` ( age , name , sex ) 
//values
//( '10' , 'lilie' , 'boy' ),
//( '20' , 'lining' , 'boy' ),
//( '30' , 'hanmeimei' , 'girl' )
f1 := map[string]any{
  "name": "lilie",
  "sex":  "boy",
  "age":  "10",
}
f2 := map[string]any{
  "name": "lining",
  "sex":  "boy",
  "age":  "20",
}
f3 := map[string]any{
  "name": "hanmeimei",
  "sex":  "girl",
  "age":  "30",
}

dwids := []map[string]any{
  f1, f2, f3,
}

gen := NewGenerator().Table("user").Inserts(dwids)

fmt.Print(gen.InsertsSql(false))
```

### 四、code-gengrator

code-gengrator 模块主要用于生成数据库表对应的struct，以及dao文件，同时会生成相关的附属类文件

#### 文件介绍

1. model 文件,struct 所在文件，每次都会更新。
2. extend文件，model的扩展文件，用于接收联表查询的返回值，只生成一次。
3. view文件，提供接口时，接口中的返回值，和model 独立，只生成一次。
4. param文件，提供接口时，用于接收接口的参数，只生成一次。
5. dao文件，提供常用的增删改查方法，只生成一次。

#### 使用教程

新建main方法，配置需要生成代码的表(支持配置多个)，运行即可生成代码

```go
package main

import (
	generator "github.com/go-lazyer/go-north/code"
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

