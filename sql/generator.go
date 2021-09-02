package generator

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

const (
	INNER_JOIN = "inner join" // inner  join
	LEFT_JOIN  = "left join"  // left  join
	RIGHT_JOIN = "right join" // right join
)

type Generator struct {
	orderBy   []string //排序字段
	pageStart int
	pageSize  int
	pageNum   int
	querys    []Query
	update    map[string]interface{}
	updates   map[interface{}]map[string]interface{}
	insert    map[string]interface{}
	inserts   []map[string]interface{}
	joins     []*Join
	tableName string
	primary   string //主键
	columns   []string
}

func NewGenerator() *Generator {
	return new(Generator)
}

func (s *Generator) Where(query ...Query) *Generator {
	if s.querys == nil {
		s.querys = make([]Query, 0)
	}
	s.querys = append(s.querys, query...)
	return s
}

func (s *Generator) Update(m map[string]interface{}) *Generator {
	s.update = m
	return s
}
func (s *Generator) Updates(m map[interface{}]map[string]interface{}) *Generator {
	s.updates = m
	return s
}
func (s *Generator) Insert(m map[string]interface{}) *Generator {
	s.insert = m
	return s
}
func (s *Generator) Inserts(m []map[string]interface{}) *Generator {
	s.inserts = m
	return s
}

func (s *Generator) Join(join ...*Join) *Generator {
	if s.joins == nil {
		s.joins = make([]*Join, 0)
	}
	s.joins = append(s.joins, join...)
	return s
}
func (s *Generator) Table(tableName string) *Generator {
	s.tableName = tableName
	return s
}
func (s *Generator) Primary(primary string) *Generator {
	s.primary = primary
	return s
}
func (s *Generator) Result(columns ...string) *Generator {
	s.columns = columns
	return s
}
func (s *Generator) PageNum(pageNum int) *Generator {
	s.pageNum = pageNum
	return s
}
func (s *Generator) PageStart(pageStart int) *Generator {
	s.pageStart = pageStart
	return s
}
func (s *Generator) PageSize(pageSize int) *Generator {
	s.pageSize = pageSize
	return s
}
func (s *Generator) OrderBy(orderBy []string) *Generator {
	s.orderBy = orderBy
	return s
}
func (s *Generator) AddOrderBy(name string, orderByType string) *Generator {
	if s.orderBy == nil {
		s.orderBy = make([]string, 0)
	}
	s.orderBy = append(s.orderBy, name+" "+orderByType)
	return s
}

func (s *Generator) CountSql(prepare bool) (string, []interface{}, error) {
	params := make([]interface{}, 0, 10)
	var sql bytes.Buffer
	sql.WriteString("select count(*) count ")
	if len(s.tableName) > 0 {
		sql.WriteString(" from  `" + s.tableName + "`")
	}

	if s.joins != nil && len(s.joins) > 0 {
		for _, join := range s.joins {
			sql.WriteString(fmt.Sprintf(" %v %v on %v", join.joinType, join.tableName, join.condition))
			for i, query := range join.querys {
				if i == 0 {
					sql.WriteString(" and ")
				} else {
					sql.WriteString(" or ")
				}
				source, param, _ := query.Source(join.tableName, prepare)
				sql.WriteString(" " + source + " ")
				params = append(params, param...)
			}
		}
	}

	if s.querys != nil && len(s.querys) > 0 {
		sql.WriteString(" where   ")
		for i, query := range s.querys {
			if i != 0 {
				sql.WriteString(" or ")
			}
			source, param, _ := query.Source(s.tableName, prepare)
			sql.WriteString(" " + source + " ")
			params = append(params, param...)
		}
	}
	return sql.String(), params, nil
}

func (s *Generator) SelectSql(prepare bool) (string, []interface{}, error) {
	params := make([]interface{}, 0, 10)
	var sql bytes.Buffer
	sql.WriteString("select ")
	if s.columns == nil {
		sql.WriteString(" * ")
	} else {
		sql.WriteString(strings.Join(s.columns, ","))
	}
	if len(s.tableName) > 0 {
		sql.WriteString(" from  `" + s.tableName + "`")
	}

	if s.joins != nil && len(s.joins) > 0 {
		for _, join := range s.joins {
			sql.WriteString(fmt.Sprintf(" %v %v on %v", join.joinType, join.tableName, join.condition))
			for i, query := range join.querys {
				if i == 0 {
					sql.WriteString(" and ")
				} else {
					sql.WriteString(" or ")
				}
				source, param, _ := query.Source(join.tableName, prepare)
				sql.WriteString(" " + source + " ")
				params = append(params, param...)
			}
		}
	}

	if s.querys != nil && len(s.querys) > 0 {
		var source string
		var param []interface{}
		for i, query := range s.querys {
			if i != 0 {
				sql.WriteString(" or ")
			}
			source, param, _ = query.Source(s.tableName, prepare)
			params = append(params, param...)
		}
		if strings.TrimSpace(source) != "" {
			sql.WriteString(" where   " + source + " ")
		}
	}

	if s.orderBy != nil && len(s.orderBy) > 0 {
		sql.WriteString(" order by   ")
		for n, v := range s.orderBy {
			if n != 0 {
				sql.WriteString(", ")
			}
			sql.WriteString(v)
		}
	}
	if s.pageSize > 0 {
		if s.pageNum > 0 {
			s.pageStart = (s.pageNum - 1) * s.pageSize
		}
		params = append(params, s.pageStart, s.pageSize)
		if prepare {
			sql.WriteString(fmt.Sprintf(" limit ?,?"))
		} else {
			sql.WriteString(fmt.Sprintf(" limit %d,%d", s.pageStart, s.pageSize))
		}
	}

	return sql.String(), params, nil
}

func (s *Generator) DeleteSql(prepare bool) (string, []interface{}, error) {
	if s.tableName == "" {
		return "", nil, errors.New("tableName is not null")
	}
	if s.querys == nil || len(s.querys) == 0 {
		return "", nil, errors.New("warn: query is not null")
	}
	params := make([]interface{}, 0, 10)
	var sql bytes.Buffer
	sql.WriteString("delete from `" + s.tableName + "` ")

	sql.WriteString(" where   ")
	for i, query := range s.querys {
		if i != 0 {
			sql.WriteString(" or ")
		}
		source, param, _ := query.Source(s.tableName, prepare)
		sql.WriteString(" " + source + " ")
		params = append(params, param...)
	}

	return sql.String(), params, nil
}
func (s *Generator) InsertSql(prepare bool) (string, []interface{}, error) {

	if s.tableName == "" {
		return "", nil, errors.New("tableName is not null")
	}
	if s.insert == nil || len(s.insert) == 0 {
		return "", nil, errors.New("The insert is not null")
	}
	//把所有要修改的字段提取出来
	fields := make([]string, 0)
	for field, _ := range s.insert {
		fields = append(fields, field)
	}

	var sql bytes.Buffer
	sql.WriteString("insert into `" + s.tableName + "` ")
	sql.WriteString("(")
	n := 0

	for _, field := range fields {
		if n != 0 {
			sql.WriteString(",")
		}
		sql.WriteString(" " + field + " ")
		n++
	}
	sql.WriteString(") values")
	n = 0
	params := make([]interface{}, 0)
	sql.WriteString("(")
	m := 0
	for _, field := range fields {
		if m != 0 {
			sql.WriteString(",")
		}
		params = append(params, s.insert[field])
		if prepare {
			sql.WriteString(" ? ")
		} else {
			sql.WriteString(fmt.Sprintf(" '%v' ", s.insert[field]))
		}
		m++
	}
	sql.WriteString(")")
	return sql.String(), params, nil
}
func (s *Generator) InsertsSql(prepare bool) (string, []interface{}, error) {

	if s.tableName == "" {
		return "", nil, errors.New("tableName is not null")
	}
	if s.inserts == nil || len(s.inserts) == 0 {
		return "", nil, errors.New("The inserts is not null")
	}
	//把所有要修改的字段提取出来
	fields := make([]string, 0)
	for field, _ := range s.inserts[0] {
		fields = append(fields, field)
	}
	var sql bytes.Buffer
	sql.WriteString("insert into `" + s.tableName + "` ")
	sql.WriteString("(")
	n := 0

	for _, field := range fields {
		if n != 0 {
			sql.WriteString(",")
		}
		sql.WriteString(" " + field + " ")
		n++
	}
	sql.WriteString(") values")
	n = 0
	params := make([]interface{}, 0)
	for _, maps := range s.inserts {
		if n != 0 {
			sql.WriteString(",")
		}
		sql.WriteString("(")
		m := 0
		for _, field := range fields {
			if m != 0 {
				sql.WriteString(",")
			}
			params = append(params, maps[field])
			if prepare {
				sql.WriteString(" ? ")
			} else {
				sql.WriteString(fmt.Sprintf(" '%v' ", maps[field]))
			}
			m++
		}
		sql.WriteString(")")
		n++
	}
	return sql.String(), params, nil
}

func (s *Generator) UpdateSql(prepare bool) (string, []interface{}, error) {
	if s.tableName == "" {
		return "", nil, errors.New("tableName is not null")
	}
	if s.update == nil || len(s.update) <= 0 {
		return "", nil, errors.New("update is not null")
	}
	params := make([]interface{}, 0, 10)
	var sql bytes.Buffer
	sql.WriteString("update `" + s.tableName + "` set ")

	n := 0
	for name, value := range s.update {
		if n != 0 {
			sql.WriteString(",")
		}
		if prepare {
			sql.WriteString(fmt.Sprintf("%v=?", name))
		} else {
			sql.WriteString(fmt.Sprintf("%v='%v'", name, value))
		}
		params = append(params, value)
		n++
	}
	if s.querys != nil && len(s.querys) > 0 {
		sql.WriteString(" where   ")
		for i, query := range s.querys {
			if i != 0 {
				sql.WriteString(" or ")
			}
			source, param, _ := query.Source(s.tableName, prepare)
			sql.WriteString(" " + source + " ")
			params = append(params, param...)
		}
	}

	return sql.String(), params, nil
}

func (s *Generator) UpdatesSql(prepare bool) (string, []interface{}, error) {

	if s.tableName == "" {
		return "", nil, errors.New("tableName is not null")
	}
	if s.primary == "" {
		return "", nil, errors.New("primary is not null")
	}

	if s.querys == nil || len(s.querys) != 1 {
		return "", nil, errors.New("the querys size must be 1")
	}

	if s.updates == nil || len(s.updates) <= 0 {
		return "", nil, errors.New("batchSet is not null")
	}
	params := make([]interface{}, 0, 10)
	var sql bytes.Buffer
	sql.WriteString("update `" + s.tableName + "` set ")

	//把所有要修改的字段提取出来
	fields := make(map[string]string)
	for _, setMap := range s.updates {
		for name, _ := range setMap {
			fields[name] = ""
		}
	}
	n := 0
	for field, _ := range fields {
		if n != 0 {
			sql.WriteString(",")
		}
		sql.WriteString(fmt.Sprintf("%v = CASE %v", field, s.primary))
		for id, setMap := range s.updates {
			v, ok := setMap[field]
			if !ok {
				continue
			}
			params = append(params, id, v)
			if prepare {
				sql.WriteString(" WHEN ? THEN ?")
			} else {
				sql.WriteString(fmt.Sprintf(" WHEN '%v' THEN '%v'", id, v))
			}

		}
		sql.WriteString(" END ")
		n++
	}
	source, param, _ := s.querys[0].Source(s.tableName, prepare)
	sql.WriteString("where " + source + " ")
	params = append(params, param...)

	return sql.String(), params, nil
}
