package north

import "fmt"

type Join struct {
	tableName  string
	tableAlias string
	condition  string
	joinType   string //inner  left  right
	querys     []BaseQuery
}

func (s *Join) Where(query ...BaseQuery) *Join {
	if s.querys == nil {
		s.querys = make([]BaseQuery, 0)
	}
	s.querys = append(s.querys, query...)
	return s
}

// Condition Join 条件
func (s *Join) Condition(firstTable string, firstField string, secondTable string, secondField string) *Join {
	s.condition = fmt.Sprintf("%v.%v=%v.%v", firstTable, firstField, secondTable, secondField)
	return s
}
func NewJoin(from, joinType string) *Join {
	return &Join{
		tableName: from,
		joinType:  joinType,
	}
}
func NewAliasJoin(from, alias, joinType string) *Join {
	return &Join{
		tableName:  from,
		tableAlias: alias,
		joinType:   joinType,
	}
}
