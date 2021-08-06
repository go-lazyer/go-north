# go-generator

### 介绍
go-generator sql生成和struct生成的工具，减少开发者自行拼接sql，和table to struct的工作量，可以专注于核心业务代码，该工具是模仿 mybatis-generator 实现的，同时又借鉴了go elastic工具的包github.com/olivere/elastic 的查询实现方式，所以熟悉 mybatis-generator和olivere 及容易上手。

### 安装教程

```
go get github.com/go-lazyer/go-generator
```

### 使用说明
go-generator 分为两个模块，sql-generator（生成sql）和code-gengrator(生成代码)，两个模块都可以单独使用，也可以配合使用
#### sql-generator
sql-generator 可以生成普通sql和预处理sql，配合golang官方提供的sql.DB,可以轻松实现增删改查
1.  生成统计sql
``` go
	//select count(1) count from user where t.id>1000
	query := NewGreaterThanQuery("id", 1000)
	gen := NewGenerator().Table("user").Where(query)
	fmt.Println(gen.CountSql(false))
```
3.  xxxx

#### code-gengrator
1.  xxxx
2.  xxxx
3.  xxxx
