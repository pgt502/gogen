# Gogen: code generators for Go [![Build Status](https://travis-ci.com/pgt502/gogen.svg?branch=master)](https://travis-ci.com/pgt502/gogen)

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
The following options can be overwritten by providing the respective flags:
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
* store interface,
* table implementation,
* store factory.

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
The following options can be overwritten by providing the respective flags:
* `-o`: output folder (default: `.`)
* `-t`: templates folder (default: `./templates/repo`)
* `-p`: pluralise struct name in table name

### Struct tags
The generator looks at the struct tags (`db`) to determine the primary key and the names of the columns. The first part of the tag is the name of the column in the table and the second one (optional) indicates if the column is part of the primary key. The fields with the `db:"-"` tag will be ignored.

### Templates
In the repository, there are 2 sets of templates included: 
* `./templates/repo`
* `./templates/store` 

Both of them have different styles of implementing access to the database (repository vs store). The repository style has an extra layer of abstraction which hides the `database/sql` package dependency.

### Example using the "repo" templates
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
package tables

import (
	"database/sql"
	testdata "github.com/pgt502/gogen/testdata"
)

type OrderTable interface {
	Insert(tx *sql.Tx, order testdata.Order) (err error)
	Update(tx *sql.Tx, order testdata.Order) (err error)
	GetAll() (orders []*testdata.Order, err error)
	Get(product string, id string) (order testdata.Order, err error)
	Delete(tx *sql.Tx, product string, id string) (err error)
}
```

* table interface implementation for postgres
```go
package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	core "github.com/pgt502/gogen/core"
	testdata "github.com/pgt502/gogen/testdata"

	"github.com/pkg/errors"
)

type pgOrderTable struct {
	tableName string
	db        core.Queryable
	columns   []string
	values    string
}

func NewPgOrderTable(q core.Queryable) (t tables.OrderTable) {
	return &pgOrderTable{
		tableName: "orders",
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

func (t *pgOrderTable) Insert(tx *sql.Tx, el testdata.Order) (err error) {
	ownTx := tx == nil
	if ownTx {
		tx, err = t.db.Begin()
		if err != nil {
			err = errors.Wrap(err, "error creating tx")
			return
		}
	}
	sqlStatement := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s)", t.tableName, strings.Join(t.columns, ","), t.values)

	var stmt *sql.Stmt
	stmt, err = tx.Prepare(sqlStatement)
	if err != nil {
		err = errors.Wrap(err, "error preparing statement")
		if ownTx {
			tx.Rollback()
		}
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(sqlStatement,
		el.Product,
		el.Id,
		el.Price,
		el.CreatedAt,
		el.CreatedBy,
		el.IsNew,
	)

	if err != nil {
		err = errors.Wrap(err, "error inserting")
		if ownTx {
			tx.Rollback()
		}
		return
	}
	if ownTx {
		err = tx.Commit()
		if err != nil {
			err = errors.Wrap(err, "error committing tx")
			return
		}
	}

	return
}

func (t *pgOrderTable) Update(tx *sql.Tx, el testdata.Order) (err error) {
	ownTx := tx == nil
	if ownTx {
		tx, err = t.db.Begin()
		if err != nil {
			err = errors.Wrap(err, "error creating tx")
			return
		}
	}
	primaryKey := "name=$1 AND id=$2"
	valueSet := "price=$3,created_at=$4,created_by=$5,isNew=$6"

	sqlStatement := fmt.Sprintf("UPDATE %s SET %s WHERE %s", t.tableName, valueSet, primaryKey)

	var stmt *sql.Stmt
	stmt, err = tx.Prepare(sqlStatement)
	if err != nil {
		err = errors.Wrap(err, "error preparing statement")
		if ownTx {
			tx.Rollback()
		}
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(sqlStatement,
		el.Product,
		el.Id,
		el.Price,
		el.CreatedAt,
		el.CreatedBy,
		el.IsNew,
	)

	if err != nil {
		err = errors.Wrap(err, "error updating")
		if ownTx {
			tx.Rollback()
		}
		return
	}
	if ownTx {
		err = tx.Commit()
		if err != nil {
			err = errors.Wrap(err, "error committing tx")
			return
		}
	}

	return
}

func (t *pgOrderTable) GetAll() (ret []*testdata.Order, err error) {
	sqlStatement := fmt.Sprintf(`SELECT %s
        FROM %s
    `, strings.Join(t.columns, ","), t.tableName)

	rows, err := t.db.Query(sqlStatement)
	if err != nil && err != sql.ErrNoRows {
		err = errors.Wrap(err, "error querying all")
		return
	}

	ret, err = t.ReadRows(rows)
	if err != nil {
		err = errors.Wrap(err, "error reading all")
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
		err = errors.Wrap(err, "error fetching")
		return
	}
	return
}

func (t *pgOrderTable) Delete(tx *sql.Tx, product string, id string) (err error) {
	ownTx := tx == nil
	if ownTx {
		tx, err = t.db.Begin()
		if err != nil {
			err = errors.Wrap(err, "error creating tx")
			return
		}
	}
	where := "name=$1 AND id=$2"
	sqlStatement := fmt.Sprintf(`DELETE 
        FROM %s
        WHERE %s`,
		t.tableName,
		where,
	)

	var stmt *sql.Stmt
	stmt, err = tx.Prepare(sqlStatement)
	if err != nil {
		err = errors.Wrap(err, "error preparing statement")
		if ownTx {
			tx.Rollback()
		}
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(sqlStatement,
		product,
		id,
	)

	if err != nil {
		err = errors.Wrap(err, "error deleting")
		if ownTx {
			tx.Rollback()
		}
		return
	}
	if ownTx {
		err = tx.Commit()
		if err != nil {
			err = errors.Wrap(err, "error committing tx")
			return
		}
	}

	return
}

func (t *pgOrderTable) ReadRows(rows core.ScannableExt) (items []*testdata.Order, err error) {
	for rows.Next() {
		var item testdata.Order
		item, err = t.ReadRow(rows)
		if err != nil {
			err = errors.Wrap(err, "error reading row from db")
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
	Create(order testdata.Order) (err error)
	Update(order testdata.Order) (err error)
	GetAll() (orders []*testdata.Order, err error)
	Get(product string, id string) (order testdata.Order, err error)
	Delete(product string, id string) (err error)
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
	err = r.db.Insert(nil, el)
	return
}

func (r *orderRepo) Update(el testdata.Order) (err error) {
	err = r.db.Update(nil, el)
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

func (r *orderRepo) Delete(product string, id string) (err error) {
	err = r.db.Delete(nil, product, id)
	return
}
```

* repository factory
```go
package repos

import (
	testdata "github.com/pgt502/gogen/testdata"
)

type RepoFactory interface {
	GetOrderRepo() OrderRepo
}

type repoFactory struct {
	storeType core.StorageType
	db        core.Queryable
}

func NewRepoFactory(storeType core.StorageType, q core.Queryable) RepoFactory {
	return &repoFactory{
		db:        q,
		storeType: storeType,
	}
}

func (f *repoFactory) GetOrderRepo() OrderRepo {
	switch f.storeType {
	case core.STORETYPE_POSTGRES:
		table := postgres.NewPgOrderTable(f.db)
		r := NewOrderRepo(table)
		return r
	default:
		// not supported
	}
	return nil
}
```

### Example using the "store" templates
The same struct as above will generate:

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

* store interface
```go
package stores

import (
	testdata "github.com/pgt502/gogen/testdata"
)

type OrderStore interface {
	Add(testdata.Order) error
	Update(testdata.Order) error
	GetAll() ([]*testdata.Order, error)
	Get(product string, id string) (testdata.Order, error)
	Delete(product string, id string) error
}
```

* store interface implementation for postgres
```go
package postgres

import (
	"database/sql"
	"fmt"
	"strings"

	core "github.com/pgt502/gogen/core"
	testdata "github.com/pgt502/gogen/testdata"

	"github.com/pkg/errors"
)

type pgOrderTable struct {
	tableName string
	db        core.Queryable
	columns   []string
	values    string
}

func NewPgOrderTable(q core.Queryable) (t stores.OrderStore) {
	return &pgOrderTable{
		tableName: "orders",
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

func (t *pgOrderTable) Add(el testdata.Order) (err error) {
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
		err = errors.Wrap(err, "error adding")
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
		err = errors.Wrap(err, "error updating")
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
		err = errors.Wrap(err, "error querying all")
		return
	}

	ret, err = ReadRows(rows)
	if err != nil {
		err = errors.Wrap(err, "error reading rows")
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
		err = errors.Wrap(err, "error getting")
		return
	}
	return
}

func (t *pgOrderTable) Delete(product string, id string) (err error) {
	where := "name=$1 AND id=$2"
	sqlStatement := fmt.Sprintf(`DELETE 
        FROM %s
        WHERE %s`,
		t.tableName,
		where,
	)

	_, err = t.db.Exec(sqlStatement,
		product,
		id,
	)

	if err != nil {
		err = errors.Wrap(err, "error deleting")
		return
	}

	return
}

func (t *pgOrderTable) ReadRows(rows core.ScannableExt) (items []*testdata.Order, err error) {
	for rows.Next() {
		var item testdata.Order
		item, err = t.ReadRow(rows)
		if err != nil {
			err = errors.Wrap(err, "error reading row from db")
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
* store factory 
```go
package stores

import (
	testdata "github.com/pgt502/gogen/testdata"
)

type StoreFactory interface {
	GetOrderStore() OrderStore
}

type storeFactory struct {
	storeType core.StorageType
	db        core.Queryable
}

func NewStoreFactory(storeType core.StorageType, q core.Queryable) StoreFactory {
	return &storeFactory{
		db:        q,
		storeType: storeType,
	}
}

func (f *storeFactory) GetOrderStore() OrderStore {
	switch f.storeType {
	case core.STORETYPE_POSTGRES:
		table := postgress.NewPgOrderTable(f.db)
		return table
	default:
		// not supported
	}
	return nil
}

```