// Create by code generator  2024-05:06 17:53:23.024
package model

import (
	"database/sql"
	"time"
)

const (
	USER_ID = "user_id" //
	DAY     = "day"     //
	NUM     = "num"     //

	TABLE_NAME = "test" // 表名
)

type TestModel struct {
	UserId sql.NullString `orm:"user_id" ` //
	Day    sql.NullTime   `orm:"day" `     //
	Num    sql.NullInt64  `orm:"num" `     //

}

func MapToStruct(m map[string]any) TestModel {
	model := TestModel{}

	if value, ok := m[USER_ID].(string); ok {
		model.UserId = sql.NullString{value, true}
	}

	if value, ok := m[DAY].(time.Time); ok {
		model.Day = sql.NullTime{value, true}
	}

	if value, ok := m[NUM].(int64); ok {
		model.Num = sql.NullInt64{value, true}
	}

	return model
}

func SliceToStructs(s []map[string]any) []TestModel {
	slices := make([]TestModel, 0)
	for _, m := range s {
		slices = append(slices, MapToStruct(m))
	}
	return slices
}

func (m *TestModel) ToMap(includeEmpty bool) map[string]any {
	view := make(map[string]any)

	if m.UserId.Valid {
		view[USER_ID] = m.UserId.String
	} else if includeEmpty {
		view[USER_ID] = nil
	}

	if m.Day.Valid {
		view[DAY] = m.Day.Time
	} else if includeEmpty {
		view[DAY] = nil
	}

	if m.Num.Valid {
		view[NUM] = m.Num.Int64
	} else if includeEmpty {
		view[NUM] = nil
	}

	return view
}
