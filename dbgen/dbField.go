package dbgen

import (
	"log"
	"strings"

	"github.com/fatih/structtag"
	"github.com/pgt502/gogen/generate"
)

type dbField struct {
	generate.Field
	tags *structtag.Tags
}

const (
	TN_TAG_NAME_DB = "db"
	TN_TAG_NAME_PK = "pk"
)

type DbField interface {
	generate.Field
	Column() string
	IsPK() bool
}

func NewDbField(f generate.Field) DbField {
	dbf := &dbField{
		Field: f,
	}
	dbf.parseTag(f.Tag())
	return dbf
}

func (f *dbField) parseTag(tag string) {
	tags, err := structtag.Parse(tag)
	if err != nil {
		log.Printf("error parsing tags: %s\n", err)
		return
	}
	f.tags = tags
}

func (f *dbField) Column() string {
	if f.tags == nil {
		return strings.ToLower(f.Name())
	}
	tag, err := f.tags.Get(TN_TAG_NAME_DB)
	if err != nil {
		log.Printf("error getting tag '%s': %s\n", TN_TAG_NAME_DB, err)
		return strings.ToLower(f.Name())
	}
	return tag.Name
}

func (f *dbField) DBType() string {
	switch t := f.Type(); t {
	case "string":
		return "TEXT"
	case "float32":
		return "REAL"
	case "float64":
		return "DOUBLE PRECISION"

	case "int8":
		fallthrough
	case "uint8":
		fallthrough
	case "uint16":
		fallthrough
	case "int16":
		return "SMALLINT"
	case "int32":
		fallthrough
	case "uint32":
		fallthrough
	case "uint":
		fallthrough
	case "int":
		return "INTEGER"
	case "uint64":
		fallthrough
	case "int64":
		return "BIGINT"
	default:
		log.Printf("%s type not supported", t)
		return t
	}
}

func (f *dbField) ColumnIndex() int {
	return f.Field.Index() + 1
}

func (f *dbField) IsPK() bool {
	if f.tags == nil {
		return false
	}
	tag, err := f.tags.Get(TN_TAG_NAME_DB)
	if err != nil {
		log.Printf("error getting tag '%s': %s\n", TN_TAG_NAME_DB, err)
		return false
	}
	return tag.HasOption(TN_TAG_NAME_PK)
}
