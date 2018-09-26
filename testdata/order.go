package testdata

type Order struct {
	Product   string  `db:"name"`
	Id        string  `json:"id" db:"id,pk"`
	Price     float64 `db:"price"`
	CreatedAt int64   `db:"created_at"`
	CreatedBy string  `db:"created_by"`
}
