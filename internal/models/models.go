package models

import (
	"fmt"
	"strings"
	"time"
)

type Conf struct {
	Port  int `yml:"port"`
	DbKey string
}

type Database struct {
	Host     string `yml:"host"`
	Port     string `yml:"port"`
	User     string `yml:"user"`
	Password string `yml:"password"`
	Name     string `yml:"name"`
}

type UserSub struct {
	Id           int        `db:"id"`
	ServiceName  string     `db:"service_name" validate:"required"`
	MonthlyPrice int    `db:"monthly_price" validate:"gt=0"`
	UserId       string     `db:"user_id" validate:"required,uuid4"`
	StartDate    MonthYear  `db:"start_date" validate:"required"`
	EndDate      *MonthYear `db:"end_date"`
}

type MonthYear time.Time

func (my *MonthYear) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	t, err := time.Parse("01-2006", s)
	if err != nil {
		return fmt.Errorf("неправильный формат даты (ожидался MM-YYYY): %w", err)
	}
	*my = MonthYear(t)
	return nil
}

func (my MonthYear) MarshalJSON() ([]byte, error) {
	t := time.Time(my)
	formatted := t.Format("01-2006")
	return []byte(`"` + formatted + `"`), nil
}

func (my MonthYear) ToTime() time.Time {
	return time.Time(my)
}
