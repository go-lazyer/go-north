// Create by code generator  2024-05:06 17:39:04.024
package param

import (
	"time"
)

type TestParam struct {
	UserId string    `form:"user_id" json:"user_id"` //
	Day    time.Time `form:"day" json:"day"`         //
	Num    int64     `form:"num" json:"num"`         //

	PageNum   int `form:"page" json:"page"`
	PageStart int `form:"start" json:"start"`
	PageSize  int `form:"size" json:"size"`
}
