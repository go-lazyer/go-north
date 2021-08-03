package generator

import (
	"fmt"
	"testing"
)

func TestGenerator_BatchUpdateSql(t *testing.T) {

	f1 := map[string]interface{}{
		"name": "name1",
		"sex":  "sex1",
		"age":  "age1",
	}
	f2 := map[string]interface{}{
		"name": "name2",
		"sex":  "sex2",
		"age":  "age2",
	}
	f3 := map[string]interface{}{
		"name": "name3",
		"sex":  "sex3",
		"age":  "age3",
	}
	set := map[interface{}]map[string]interface{}{
		"1": f1,
		"2": f2,
		"3": f3,
	}
	dwids := []interface{}{
		123, 234, 345,
	}
	query := NewInQuery("dwid", dwids)
	gen := NewGenerator().Table("user1").Where(query).Primary("dwid").Updates(set)
	fmt.Print(gen.BatchUpdateSql(false))
}

func TestGenerator_BatchInsertSql(t *testing.T) {
	f1 := map[string]interface{}{
		"name": "name1",
		"sex":  "sex1",
		"age":  "age1",
	}
	f2 := map[string]interface{}{
		"name": "name2",
		"sex":  "sex2",
		"age":  "age2",
	}
	f3 := map[string]interface{}{
		"name": "name3",
		"sex":  "sex3",
		"age":  "age3",
	}
	dwids := []map[string]interface{}{
		f1, f2, f3,
	}
	gen := NewGenerator().Table("user1").Inserts(dwids)
	fmt.Print(gen.BatchInsertSql(false))

}

func TestGenerator_SelectSql(t *testing.T) {
	gen := NewGenerator().Table("user1")
	gen.SelectSql(false)
}
