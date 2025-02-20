package north

import (
	"database/sql"
	"fmt"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

const (
	USER_ID = "user_id" //
	DAY     = "day"     //
	NUM     = "num"     //
)

type Test struct {
	UserId sql.NullString `orm:"user_id" ` //
	Day    sql.NullTime   `orm:"day" `     //
	Num    sql.NullInt64  `orm:"num" `     //
}
type Test1 struct {
	Test
}

func TestFull(t *testing.T) {
	username := "root"
	password := "Daoway_Mysql_iO12"
	host := "test.daoway.cn"
	port := "3306"
	dbname := "daowei"
	connStr := fmt.Sprintf("%s:%s@(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local", username, password, host, port, dbname)
	// 创建数据库连接
	ds, err := Open(DRIVER_NAME_MYSQL, connStr, &Config{MaxOpenConns: 10, MaxIdleConns: 10})
	if err != nil {
		fmt.Println(err)
	}
	sql1 := "select * from test t where t.user_id='a'"
	params := make([]any, 0)
	ts, err := Query[Test1](sql1, params, ds)
	if err != nil {
		fmt.Println(err)
	}
	for _, testInstance := range ts {
		fmt.Print(testInstance)
		// 如果有更多字段，可以继续访问它们
	}
}
func TestGenerator_CountSql(t *testing.T) {
	//select count(1) count from user where t.id>1000
	query := NewGreaterThanQuery("id", 1000)
	gen := CreateNorth().Table("user").Where(query)
	fmt.Println(gen.CountSql(false))
}

func TestGenerator_GroupSql(t *testing.T) {
	//select count(1) count from user where t.id>1000
	query := NewGreaterThanQuery("id", 1000)
	gen := CreateNorth().Table("user").Where(query)
	fmt.Println(gen.CountSql(false))
}

func TestGenerator_SelectSql1(t *testing.T) {
	//select * from user
	query1 := NewBoolQuery()
	gen := CreateNorth().Table("user").Where(query1)
	fmt.Println(gen.SelectSql(false))
}
func TestGenerator_SelectSql2(t *testing.T) {
	//select * from user where t.id=1000
	query := NewEqualQuery("id", 1000)
	gen := CreateNorth().Table("user").Where(query)
	fmt.Println(gen.SelectSql(false))
}
func TestGenerator_SelectSql3(t *testing.T) {
	// select * from user where id=1000 and age>20 order by age desc
	idQuery := NewEqualQuery("id", 1000)
	ageQuery := NewGreaterThanQuery("age", 20)
	boolQuery := NewBoolQuery().And(idQuery, ageQuery)
	gen := CreateNorth().Table("user").Where(boolQuery).AddOrderBy("age", "desc")
	fmt.Println(gen.SelectSql(false))
}
func TestGenerator_SelectSql4(t *testing.T) {
	// select id,name,age from user where (id=1000 and age>20) or age <=10 order by age desc ,id asc
	idQuery := NewEqualQuery("id", 1000)
	ageQuery := NewGreaterThanQuery("age", 20)
	boolQuery := NewBoolQuery().And(idQuery, ageQuery)
	ageQuery2 := NewLessThanOrEqualQuery("age", 10)
	gen := CreateNorth().Result("id", "name", "age").Table("user").Where(boolQuery, ageQuery2).AddOrderBy("age", "desc").AddOrderBy("id", "asc")
	fmt.Println(gen.SelectSql(false))
}
func TestGenerator_SelectSql5(t *testing.T) {
	// select user.id,order.id  from user join order on user.id=order.user_id where user.id='10000'
	idQuery := NewEqualQuery("id", 1000)

	join := NewAliasJoin("order", "o", INNER_JOIN).Condition("u", "id", "o", "user_id")
	gen := CreateNorth().Result("u.id", "o.id").TableAlias("user", "u").Join(join).Where(idQuery)
	fmt.Println(gen.SelectSql(false))
}

func TestGenerator_SelectSql6(t *testing.T) {
	// select user.sex,count(user.sex) count  from user group by user.sex

	gen := CreateNorth().Result("user.sex", "count(user.sex) count").Table("user").AddGroupBy("user", "sex")
	fmt.Println(gen.SelectSql(false))
}

func TestGenerator_SelectSql7(t *testing.T) {
	// select user.id,order.id  from user join order on user.id=order.user_id and order.create_time>user.create_time where user.id='10000'
	idQuery := NewEqualQuery("id", 1000)

	join := NewAliasJoin("order", "o", INNER_JOIN).Condition("u", "id", "o", "user_id").Where(NewFieldGreaterThanQuery("o", "create_time", "u", "create_time"))
	gen := CreateNorth().Result("u.id", "o.id").TableAlias("user", "u").Join(join).Where(idQuery)
	fmt.Println(gen.SelectSql(true))
}

func TestGenerator_UpdateSql(t *testing.T) {
	// update user set age=21,name="lazeyr" where id="10000"
	query := NewEqualQuery("id", 1000)
	set := map[string]any{
		"age":  21,
		"name": "lazyer",
	}
	gen := CreateNorth().Table("user").Where(query).Update(set)
	fmt.Println(gen.UpdateSql(false))
}

func TestGenerator_UpdatesSql(t *testing.T) {
	// update
	// `user`
	// set
	// 	sex = case dwid
	// 	when 10001 then boy
	// 	when 10002 then boy
	// 	when 10003 then girl
	// 	end,
	// 	age = case dwid
	// 	when 10001 then 10
	// 	when 10002 then 20
	// 	when 10003 then 30
	// 	end,
	// 	name = case dwid
	// 	when 10001 then lilie
	// 	when 10002 then lining
	// 	when 10003 then hanmeimei
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
	set := []map[string]any{
		f1, f2, f3,
	}
	dwids := []any{
		10001, 10002, 10003,
	}
	query := NewInQuery("dwid", dwids)
	gen := CreateNorth().Table("user").Where(query).Primary("dwid").Updates(set)
	fmt.Print(gen.UpdateSql(false))
}

func TestGenerator_InsertsSql(t *testing.T) {
	//insert into `user` ( age , name , sex ) values( '10' , 'lilie' , 'boy' ),
	//( '20' , 'lining' , 'boy' ),
	//( '30' , 'hanmeimei' , 'girl' )
	f1 := map[string]any{
		"sex":  "boy",
		"name": "lilie",
		"age":  "10",
	}
	f2 := map[string]any{
		"name": "lining",
		"age":  "20",
		"sex":  "boy",
	}
	f3 := map[string]any{
		"name": "hanmeimei",
		"sex":  "girl",
		"age":  "30",
	}
	dwids := []map[string]any{
		f1, f2, f3,
	}
	gen := CreateNorth().Table("user").Inserts(dwids)
	fmt.Print(gen.InsertSql(false))
}

func TestGenerator_InsertSql(t *testing.T) {
	f3 := map[string]any{
		"name": "hanmeimei",
		"sex":  "girl",
		"age":  "30",
	}
	gen := CreateNorth().Table("user").Insert(f3)
	fmt.Print(gen.InsertSql(false))
}
