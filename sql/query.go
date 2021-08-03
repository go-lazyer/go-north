package generator

import (
	"bytes"
	"fmt"
	"go-generator/util"
)

type Query interface {
	Source(table string, prepare bool) (string, []interface{}, error)
}

type NullQuery struct {
	field string
}

func NewNullQuery(field string) *NullQuery {
	return &NullQuery{field: field}
}

func (q *NullQuery) Source(table string, prepare bool) (string, []interface{}, error) {
	return table + "." + q.field + " is null", nil, nil
}

type NotNullQuery struct {
	field string
}

func NewNotNullQuery(field string) *NotNullQuery {
	return &NotNullQuery{field: field}
}

func (q *NotNullQuery) Source(table string, prepare bool) (string, []interface{}, error) {
	return table + "." + q.field + " is not null", nil, nil
}

type BetweenQuery struct {
	field       string
	firstValue  interface{}
	secondValue interface{}
}

func NewBetweenQuery(field string, firstValue interface{}, secondValue interface{}) *BetweenQuery {
	return &BetweenQuery{field: field, firstValue: firstValue, secondValue: secondValue}
}

func (q *BetweenQuery) Source(table string, prepare bool) (string, []interface{}, error) {
	param := []interface{}{q.firstValue, q.secondValue}
	if prepare {
		return fmt.Sprintf("%s.%s between ? and ?", table, q.field), param, nil
	}
	if util.IsNumberType(q.firstValue) {
		return fmt.Sprintf("%s.%s between %v and %v", table, q.field, q.firstValue, q.secondValue), param, nil
	} else {
		return fmt.Sprintf("%s.%s between '%v' and '%v' ", table, q.field, q.firstValue, q.secondValue), param, nil
	}
}

type NotBetweenQuery struct {
	field       string
	firstValue  interface{}
	secondValue interface{}
}

func NewNotBetweenQuery(field string, firstValue interface{}, secondValue interface{}) *NotBetweenQuery {
	return &NotBetweenQuery{field: field, firstValue: firstValue, secondValue: secondValue}
}

func (q *NotBetweenQuery) Source(table string, prepare bool) (string, []interface{}, error) {
	param := []interface{}{q.firstValue, q.secondValue}
	if prepare {
		return fmt.Sprintf("%s.%s not between ? and ?", table, q.field), param, nil
	}
	if util.IsNumberType(q.firstValue) {
		return fmt.Sprintf("%s.%s not between %v and %v", table, q.field, q.firstValue, q.secondValue), param, nil
	} else {
		return fmt.Sprintf("%s.%s not between '%v' and '%v' ", table, q.field, q.firstValue, q.secondValue), param, nil
	}
}

type EqualQuery struct {
	field string
	value interface{}
}

func NewEqualQuery(field string, value interface{}) *EqualQuery {
	return &EqualQuery{field: field, value: value}
}

func (q *EqualQuery) Source(table string, prepare bool) (string, []interface{}, error) {
	if prepare {
		return fmt.Sprintf("%s.%s = ?", table, q.field), []interface{}{q.value}, nil
	}
	if util.IsNumberType(q.value) {
		return fmt.Sprintf("%s.%s = %v", table, q.field, q.value), []interface{}{q.value}, nil
	} else {
		return fmt.Sprintf("%s.%s = '%v'", table, q.field, q.value), []interface{}{q.value}, nil
	}
}

type NotEqualQuery struct {
	field string
	value interface{}
}

func NewNotEqualQuery(field string, value interface{}) *NotEqualQuery {
	return &NotEqualQuery{field: field, value: value}
}

func (q *NotEqualQuery) Source(table string, prepare bool) (string, []interface{}, error) {
	if prepare {
		return fmt.Sprintf("%s.%s != ?", table, q.field), []interface{}{q.value}, nil
	}
	if util.IsNumberType(q.value) {
		return fmt.Sprintf("%s.%s != %v", table, q.field, q.value), []interface{}{q.value}, nil
	} else {
		return fmt.Sprintf("%s.%s != '%v'", table, q.field, q.value), []interface{}{q.value}, nil
	}
}

type InQuery struct {
	field string
	value []interface{}
}

func NewInQuery(field string, value []interface{}) *InQuery {
	return &InQuery{field: field, value: value}
}

func (q *InQuery) Source(table string, prepare bool) (string, []interface{}, error) {
	var sql bytes.Buffer
	sql.WriteString(table + "." + q.field + " in (")
	for k, v := range q.value {
		if k != 0 {
			sql.WriteString(" ,")
		}
		if prepare {
			sql.WriteString(" ?")
		} else if util.IsNumberType(q.value) {
			sql.WriteString(fmt.Sprintf(" %v ", v))
		} else {
			sql.WriteString(fmt.Sprintf(" '%v' ", v))
		}

	}
	sql.WriteString(")")
	return sql.String(), q.value, nil
}

type NotInQuery struct {
	field string
	value []interface{}
}

func NewNotInQuery(field string, value []interface{}) *NotInQuery {
	return &NotInQuery{field: field, value: value}
}

func (q *NotInQuery) Source(table string, prepare bool) (string, []interface{}, error) {
	var sql bytes.Buffer
	sql.WriteString(table + "." + q.field + " not in (")
	for k, v := range q.value {
		if k != 0 {
			sql.WriteString(" ,")
		}
		if prepare {
			sql.WriteString(" ?")
		} else if util.IsNumberType(q.value) {
			sql.WriteString(fmt.Sprintf(" %v ", v))
		} else {
			sql.WriteString(fmt.Sprintf(" '%v' ", v))
		}

	}
	sql.WriteString(")")
	return sql.String(), q.value, nil
}

type LikeQuery struct {
	field string
	value interface{}
}

func NewLikeQuery(field string, value interface{}) *LikeQuery {
	return &LikeQuery{field: field, value: value}
}

func (q *LikeQuery) Source(table string, prepare bool) (string, []interface{}, error) {
	if prepare {
		return fmt.Sprintf("%s.%s like ?", table, q.field), []interface{}{q.value}, nil
	}
	if util.IsNumberType(q.value) {
		return fmt.Sprintf("%s.%s like '%v'", table, q.field, q.value), []interface{}{q.value}, nil
	} else {
		return fmt.Sprintf("%s.%s like '%v'", table, q.field, q.value), []interface{}{q.value}, nil
	}
}

type NotLikeQuery struct {
	field string
	value interface{}
}

func NewNotLikeQuery(field string, value interface{}) *NotLikeQuery {
	return &NotLikeQuery{field: field, value: value}
}

func (q *NotLikeQuery) Source(table string, prepare bool) (string, []interface{}, error) {
	if prepare {
		return fmt.Sprintf("%s.%s not like ?", table, q.field), []interface{}{q.value}, nil
	}
	if util.IsNumberType(q.value) {
		return fmt.Sprintf("%s.%s not like %v", table, q.field, q.value), []interface{}{q.value}, nil
	} else {
		return fmt.Sprintf("%s.%s not like '%v'", table, q.field, q.value), []interface{}{q.value}, nil
	}
}

type GreaterThanQuery struct {
	field string
	value interface{}
}

func NewGreaterThanQuery(field string, value interface{}) *GreaterThanQuery {
	return &GreaterThanQuery{field: field, value: value}
}

func (q *GreaterThanQuery) Source(table string, prepare bool) (string, []interface{}, error) {
	if prepare {
		return fmt.Sprintf("%s.%s > ?", table, q.field), []interface{}{q.value}, nil
	}
	if util.IsNumberType(q.value) {
		return fmt.Sprintf("%s.%s > %v", table, q.field, q.value), []interface{}{q.value}, nil
	} else {
		return fmt.Sprintf("%s.%s > '%v'", table, q.field, q.value), []interface{}{q.value}, nil
	}
}

type GreaterThanOrEqualQuery struct {
	field string
	value interface{}
}

func NewGreaterThanOrEqualQuery(field string, value interface{}) *GreaterThanOrEqualQuery {
	return &GreaterThanOrEqualQuery{field: field, value: value}
}

func (q *GreaterThanOrEqualQuery) Source(table string, prepare bool) (string, []interface{}, error) {
	if prepare {
		return fmt.Sprintf("%s.%s >= ?", table, q.field), []interface{}{q.value}, nil
	}
	if util.IsNumberType(q.value) {
		return fmt.Sprintf("%s.%s >= %v", table, q.field, q.value), []interface{}{q.value}, nil
	} else {
		return fmt.Sprintf("%s.%s >= '%v'", table, q.field, q.value), []interface{}{q.value}, nil
	}
}

type LessThanQuery struct {
	field string
	value interface{}
}

func NewLessThanQuery(field string, value interface{}) *LessThanQuery {
	return &LessThanQuery{field: field, value: value}
}

func (q *LessThanQuery) Source(table string, prepare bool) (string, []interface{}, error) {
	if prepare {
		return fmt.Sprintf("%s.%s < ?", table, q.field), []interface{}{q.value}, nil
	}
	if util.IsNumberType(q.value) {
		return fmt.Sprintf("%s.%s < %v", table, q.field, q.value), []interface{}{q.value}, nil
	} else {
		return fmt.Sprintf("%s.%s < '%v'", table, q.field, q.value), []interface{}{q.value}, nil
	}
}

type LessThanOrEqualQuery struct {
	field string
	value interface{}
}

func NewLessThanOrEqualQuery(field string, value interface{}) *LessThanOrEqualQuery {
	return &LessThanOrEqualQuery{field: field, value: value}
}

func (q *LessThanOrEqualQuery) Source(table string, prepare bool) (string, []interface{}, error) {
	if prepare {
		return fmt.Sprintf("%s.%s <= ?", table, q.field), []interface{}{q.value}, nil
	}
	if util.IsNumberType(q.value) {
		return fmt.Sprintf("%s.%s <= %v", table, q.field, q.value), []interface{}{q.value}, nil
	} else {
		return fmt.Sprintf("%s.%s <= '%v'", table, q.field, q.value), []interface{}{q.value}, nil
	}
}

type FieldEqualQuery struct {
	firstField  string
	secondField string
}

func NewFieldEqualQuery(firstField, secondField string) *FieldEqualQuery {
	return &FieldEqualQuery{firstField: firstField, secondField: secondField}
}

func (q *FieldEqualQuery) Source(table string, prepare bool) (string, []interface{}, error) {
	return fmt.Sprintf("%s.%s = %s.%s", table, q.firstField, table, q.secondField), []interface{}{}, nil
}

type BoolQuery struct {
	query []Query
}

func NewBoolQuery() *BoolQuery {
	return &BoolQuery{
		query: make([]Query, 0),
	}
}

func (q *BoolQuery) And(queries ...Query) *BoolQuery {
	q.query = append(q.query, queries...)
	return q
}

func (q *BoolQuery) Source(table string, prepare bool) (string, []interface{}, error) {
	params := make([]interface{}, 0)
	if q.query == nil || len(q.query) <= 0 {
		return "", params, nil
	}
	var sql bytes.Buffer
	sql.WriteString("(")
	if q.query != nil || len(q.query) > 0 {
		for k, query := range q.query {
			if k != 0 {
				sql.WriteString(" and")
			}
			source, param, _ := query.Source(table, prepare)
			params = append(params, param...)
			sql.WriteString(" " + source + " ")
		}
	}
	sql.WriteString(")")
	return sql.String(), params, nil
}
