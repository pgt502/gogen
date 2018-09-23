package testdata

type Order struct {
	Id        string  `json:"id" db:"id,pk"`
	Product   string  `db:"name,pk"`
	Price     float64 `db:"price,pk"`
	CreatedAt int64   `db:"created_at,pk"`
	CreatedBy string  `db:"created_by,pk"`
}
