package response

import (
	"fmt"
	"time"
)

type JsonTime time.Time

func (j JsonTime) MarshalJSON() ([]byte, error) {
	var stmp = fmt.Sprintf("\"%s\"", time.Time(j).Format("2016-01-05"))
	return []byte(stmp), nil
}

type UserResponse struct {
	Id       int32    `json:"Id"`
	NickName string   `json:"name"`
	Birthday JsonTime `json:"birthday;"`
	Gender   string   `json:"gender"`
	Mobile   string   `json:"mobile"`
}
