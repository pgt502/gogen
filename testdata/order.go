package testdata

import "time"

type Order struct {
	Product    string  `db:"name,pk"`
	Id         string  `json:"id" db:"id,pk"`
	Price      float64 `db:"price"`
	CreatedAt  int64   `db:"created_at"`
	CreatedBy  string  `db:"created_by"`
	IsNew      bool    `db:"isNew"`
	PriceLabel string
}

type TestType struct {
	Name      string
	Id        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

type TestInterface interface {
	Method1() error
	Method2(string, interface{}) (interface{}, error)
	Method3([]*TestType) []TestType
	Method4(inputs ...string)
	Method5(test *TestType) (ret *TestType, err error)
}
