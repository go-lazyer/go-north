package generator

import (
	"fmt"
	"testing"
)

func TestGenerator_CountSql(t *testing.T) {
	//select count(1) count from user where t.id>1000
	query := NewGreaterThanQuery("id", 1000)
	gen := NewGenerator().Table("user").Where(query)
	fmt.Println(gen.CountSql(false))
}

func TestGenerator_GroupSql(t *testing.T) {
	//select count(1) count from user where t.id>1000
	query := NewGreaterThanQuery("id", 1000)
	gen := NewGenerator().Table("user").Where(query)
	fmt.Println(gen.CountSql(false))
}

func TestGenerator_SelectSql1(t *testing.T) {
	//select * from user
	query1 := NewBoolQuery()
	gen := NewGenerator().Table("user").Where(query1)
	fmt.Println(gen.SelectSql(false))
}
func TestGenerator_SelectSql2(t *testing.T) {
	//select * from user where t.id=1000
	query := NewEqualQuery("id", 1000)
	gen := NewGenerator().Table("user").Where(query)
	fmt.Println(gen.SelectSql(false))
}
func TestGenerator_SelectSql3(t *testing.T) {
	// select * from user where id=1000 and age>20 order by age desc
	idQuery := NewEqualQuery("id", 1000)
	ageQuery := NewGreaterThanQuery("age", 20)
	boolQuery := NewBoolQuery().And(idQuery, ageQuery)
	gen := NewGenerator().Table("user").Where(boolQuery).AddOrderBy("age", "desc")
	fmt.Println(gen.SelectSql(false))
}
func TestGenerator_SelectSql4(t *testing.T) {
	// select id,name,age from user where (id=1000 and age>20) or age <=10 order by age desc ,id asc
	idQuery := NewEqualQuery("id", 1000)
	ageQuery := NewGreaterThanQuery("age", 20)
	boolQuery := NewBoolQuery().And(idQuery, ageQuery)
	ageQuery2 := NewLessThanOrEqualQuery("age", 10)
	gen := NewGenerator().Result("id", "name", "age").Table("user").Where(boolQuery, ageQuery2).AddOrderBy("age", "desc").AddOrderBy("id", "asc")
	fmt.Println(gen.SelectSql(false))
}
func TestGenerator_SelectSql5(t *testing.T) {
	// select user.id,order.id  from user join order on user.id=order.user_id where user.id='10000'
	idQuery := NewEqualQuery("id", 1000)

	join := NewJoin("order", INNER_JOIN).Condition("user", "id", "order", "user_id")
	gen := NewGenerator().Result("user.id", "order.id").Table("user").Join(join).Where(idQuery)
	fmt.Println(gen.SelectSql(false))
}

func TestGenerator_SelectSql6(t *testing.T) {
	// select user.sex,count(user.sex) count  from user group by user.sex

	gen := NewGenerator().Result("user.sex", "count(user.sex) count").Table("user").AddGroupBy("user.sex")
	fmt.Println(gen.SelectSql(false))
}

func TestGenerator_UpdateSql(t *testing.T) {
	// update user set age=21,name="lazeyr" where id="10000"
	query := NewEqualQuery("id", 1000)
	set := map[string]any{
		"age":  21,
		"name": "lazyer",
	}
	gen := NewGenerator().Table("user").Where(query).Update(set)
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
	fmt.Print(gen.UpdatesSql(false))
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
	gen := NewGenerator().Table("user").Inserts(dwids)
	fmt.Print(gen.InsertsSql(false))
}

func TestGenerator_InsertSql(t *testing.T) {
	f3 := map[string]any{
		"name": "hanmeimei",
		"sex":  "girl",
		"age":  "30",
	}
	gen := NewGenerator().Table("user").Insert(f3)
	fmt.Print(gen.InsertSql(false))
}
