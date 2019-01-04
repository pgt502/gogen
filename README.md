# Gogen: code generators for Go

Package gogen consists of 2 code generators with multiple templates:
* DB code generator
* Mock generator

It is a simple code generator aiming to cater for 90% of the cases. The idea is that the code gets generated at the beginning and then a user is supposed to make any necessary changes. It provides a foundation and not a complete solution. If you use a slightly different style or types then you can adjust the templates as you like without having to change the code. 

## Credits: 
Based on and used some of the code from [github.com/ernesto-jimenez/gogen](http://github.com/ernesto-jimenez/gogen)


## Installation
```text
$ go get github.com/pgt502/gogen
```

## Mock generator
It was inspired by [this](https://www.youtube.com/watch?v=uFXfTXSSt4I) talk where a mock is a struct which implements all the methods of the passed interface. Each method has a field in the mock struct named after its method with a postfix `Call`. All the input arguments are in a sub-struct `Receives` and all return arguments are available from the `Returns` sub-struct. This way test setups are very easy.
Additionally, each method will have a `GetsCalled` field which will be incremented every time a method is called.

### Usage
To build the mock generator run:

```text
$ cd cmd/mockgen
$ go build
```

To generate a mock file for an interface, provide a (full) package name in `-pkg` parameter followed by the interface name. For example, to generate a mock for  `TestInterface` in `github.com/pgt502/gogen/testdata` package run:
```text
$ ./mockgen -pkg=github.com/pgt502/gogen/testdata TestInterface
```

#### Other options
The following options can be overwritten by providing the relative flags:
* `-o`: output folder (default: `.`)
* `-t`: templates folder (default: `./templates`)

### For example
Given the following interface:
```go
type Tester interface{
  Test(input string) (err error)
}
```
the following mock will be generated
```go
package mocks

type MockTester struct{
  TestCall struct{
    Receives struct{
      Input string
    }
    Returns struct{
      Err error
    }
    GetsCalled struct{
      Times int
    }
  }
}

func(m *MockTester)Test(input string) (err error) {
  m.TestCall.GetsCalled.Times++
  m.TestCall.Receives.Input = input
  err = m.TestCall.Returns.Err
  return
}
```

### Notes
The mock generator uses the names of the arguments to generate meaningful struct fields. If the arguments dont have names they will be named `param0`, `param1`... for input parameters and `ret0`, `ret1`... for output parameters.  

## DB generator
Generates code from a struct necessary to build and access data from database. At the moment only `PostgreSQL` is supported.
The generated files include:
* migration scripts, 
* table interface,
* repository.

The supported operations are the most common ones: insert, update, get all and get which fetches a single row from database by providing the primary key as parameters. These are to act as templates and is some cases will need to be tweaked to provide other operations, support joins or transactions.

### Usage
To build the db code generator run:
```text
$ cd cmd/dbgen
$ go build
```

To generate db code for a struct, provide a (full) package name in `-pkg` parameter followed by the struct name. For example, for the `Order` struct in `github.com/pgt502/gogen/testdata` package run:
```text
$ ./dbgen -pkg=github.com/pgt502/gogen/testdata Order
```

#### Other options
The following options can be overwritten by providing the relative flags:
* `-o`: output folder (default: `.`)
* `-t`: templates folder (default: `./templates`)
* `-p`: pluralise struct name in table name

### Struct tags
The generator looks at the struct tags (`db`) to determine the primary key and the names of the columns. The first part of the tag is the name of the column in the table and the second one (optional) indicates if the column is part of the primary key. The fields with the `db:"-"` tag will be ignored.

### For example
The following struct (from github.com/pgt502/gogen/testdata package)

```go
type Order struct {
	Product   string  `db:"name,pk"` // column name will be "name", part of Primary Key
	Id        string  `json:"id" db:"id,pk"` // column name "id", part of Primary Key
	Price     float64 `db:"price"` 
	CreatedAt int64   `db:"created_at"`
	CreatedBy string  `db:"created_by"`
	IsNew     bool    `db:"isNew"`
	PriceLabel string `db:"-"` // this field will be ignored
}
```
will generate:
* migration up script
```sql
CREATE TABLE IF NOT EXISTS public.order(
     "name" TEXT NOT NULL,
     "id" TEXT NOT NULL,
     "price" DOUBLE PRECISION NOT NULL,
     "created_at" BIGINT NOT NULL,
     "created_by" TEXT NOT NULL,
     "isNew" BOOLEAN NOT NULL,
    
    PRIMARY KEY ("name","id")
);
```
* migration down script
```sql
DROP TABLE IF EXISTS public.order;
```

* table interface
```go
package table

import (
	testdata "github.com/pgt502/gogen/testdata"
)

type OrderTable interface {
	Insert(testdata.Order) error
	Update(testdata.Order) error
	GetAll() ([]*testdata.Order, error)
	Get(product string, id string) (testdata.Order, error)
}
```

* table interface implementation for postgres
```go
package postgress

import (
	"database/sql"
	"fmt"
	"strings"

	core "github.com/pgt502/gogen/core"
	testdata "github.com/pgt502/gogen/testdata"
)

type pgOrderTable struct {
	tableName string
	db        core.Queryable
	columns   []string
	values    string
}

func NewPgOrderTable(q core.Queryable) (t tables.OrderTable) {
	return &pgOrderTable{
		tableName: "order",
		db:        q,
		columns: []string{
			"name",
			"id",
			"price",
			"created_at",
			"created_by",
			"isNew",
		},
		values: "$1,$2,$3,$4,$5,$6",
	}
}

func (t *pgOrderTable) Insert(el testdata.Order) (err error) {
	sqlStatement := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", t.tableName, strings.Join(t.columns, ","), t.values)

	_, err = t.db.Exec(sqlStatement,
		el.Product,
		el.Id,
		el.Price,
		el.CreatedAt,
		el.CreatedBy,
		el.IsNew,
	)

	if err != nil {
		//logging.LogErrore(err)
		return
	}

	return
}

func (t *pgOrderTable) Update(el testdata.Order) (err error) {
	primaryKey := "name=$1 AND id=$2"
	valueSet := "price=$3,created_at=$4,created_by=$5,isNew=$6"

	sqlStatement := fmt.Sprintf("UPDATE %s SET %s WHERE %s", t.tableName, valueSet, primaryKey)
	_, err = t.db.Exec(sqlStatement,
		el.Product,
		el.Id,
		el.Price,
		el.CreatedAt,
		el.CreatedBy,
		el.IsNew,
	)

	if err != nil {
		//logging.LogErrore(err)
		return
	}

	return
}

func (t *pgOrderTable) GetAll() (ret []*testdata.Order, err error) {
	sqlStatement := fmt.Sprintf(`SELECT %s
        FROM %s
    `, strings.Join(t.columns, ","), t.tableName)

	rows, err := t.db.Query(sqlStatement)
	if err != nil && err != sql.ErrNoRows {
		//logging.LogErrore(err)
		return
	}

	ret, err = ReadRows(rows)
	if err != nil {
		//logging.LogErrore(err)
		return
	}

	return
}

func (t *pgOrderTable) Get(product string, id string) (ret testdata.Order, err error) {
	where := "name=$1 AND id=$2"
	sqlStatement := fmt.Sprintf(`SELECT %s
        FROM %s
        WHERE %s`,
		strings.Join(t.columns, ","),
		t.tableName,
		where,
	)

	row := t.db.QueryRow(sqlStatement,
		product,
		id,
	)
	ret, err = t.ReadRow(row)
	if err != nil && err != sql.ErrNoRows {
		//logging.LogErrore(err)
		return
	}
	return
}

func (t *pgOrderTable) ReadRows(rows core.ScannableExt) (items []*testdata.Order, err error) {
	for rows.Next() {
		var item testdata.Order
		item, err = t.ReadRow(rows)
		if err != nil {
			//logging.LogErrore(err)
			return
		}
		items = append(items, &item)
	}
	return
}

func (t *pgOrderTable) ReadRow(row core.Scannable) (item testdata.Order, err error) {
	err = row.Scan(
		&item.Product,
		&item.Id,
		&item.Price,
		&item.CreatedAt,
		&item.CreatedBy,
		&item.IsNew,
	)
	return
}

```
* repository 
```go
package repos

import (
	testdata "github.com/pgt502/gogen/testdata"
)

type OrderRepo interface {
	Create(testdata.Order) error
	Update(testdata.Order) error
	GetAll() ([]*testdata.Order, error)
	Get(product string, id string) (testdata.Order, error)
}

type orderRepo struct {
	db tables.OrderTable
}

func NewOrderRepo(tb tables.OrderTable) OrderRepo {
	return &orderRepo{
		db: tb,
	}
}

func (r *orderRepo) Create(el testdata.Order) (err error) {
	err = r.db.Insert(el)
	return
}

func (r *orderRepo) Update(el testdata.Order) (err error) {
	err = r.db.Update(el)
	return
}

func (r *orderRepo) GetAll() (ret []*testdata.Order, err error) {
	ret, err = r.db.GetAll()
	return
}

func (r *orderRepo) Get(product string, id string) (ret testdata.Order, err error) {
	ret, err = r.db.Get(product, id)
	return
}

```