package stores

import (
    {{.Package}} "{{.PackagePath}}"
    {{range $path, $name := .Imports}}
    {{$name}} "{{$path}}"{{end}}
)

type StoreFactory interface{
    Get{{.Name}}Store() {{.Name}}Store
}

type storeFactory struct{
    storeType core.StorageType
    db core.Queryable
}

func NewStoreFactory(storeType core.StorageType, q core.Queryable) StoreFactory {
    return &storeFactory{
        db: q,
        storeType: storeType,
    }    
}

func (f *storeFactory) Get{{.Name}}Store() {{.Name}}Store {
    switch f.storeType {
    case core.STORETYPE_POSTGRES:
        table := {{.PGTablePackage}}.NewPg{{.Name}}Table(f.db)        
        return table
    default:
    // not supported
    }
    return nil
}