// Create by code generator  2024-05:06 17:39:04.024
package view

import (
	"go-generator/test/model"
	"time"
)

type TestView struct {
	UserId string    `json:"user_id"` //
	Day    time.Time `json:"day"`     //
	Num    int64     `json:"num"`     //

}

func Convert(m *model.TestModel) *TestView {
	return &TestView{
		UserId: m.UserId.String,
		Day:    m.Day.Time,
		Num:    m.Num.Int64,
	}
}
func Converts(models []model.TestModel) []TestView {
	views := make([]TestView, 0, len(models))
	for _, model := range models {
		views = append(views, *Convert(&model))
	}
	return views
}

func ConvertExtend(m *model.TestExtend) *TestView {
	view := Convert(&m.TestModel)
	return view
}
func ConvertExtends(extends []model.TestExtend) []TestView {
	views := make([]TestView, 0, len(extends))
	for _, extend := range extends {
		views = append(views, *ConvertExtend(&extend))
	}
	return views
}
