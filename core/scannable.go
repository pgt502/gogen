package core

type Scannable interface {
	Scan(dest ...interface{}) error
}

type ScannableExt interface {
	Scannable
	Next() bool
	Close() error
}
