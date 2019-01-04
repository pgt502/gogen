package repos

import (
    {{.Package}} "{{.PackagePath}}"
    {{range $path, $name := .Imports}}
    {{$name}} "{{$path}}"{{end}}
)

type RepoFactory interface{
    Get{{.Name}}Repo() {{.Name}}Repo
}

type repoFactory struct{
    storeType core.StorageType
    db core.Queryable
}

func NewRepoFactory(storeType core.StorageType, q core.Queryable) RepoFactory {
    return &repoFactory{
        db: q,
        storeType: storeType,
    }    
}

func (f *repoFactory) Get{{.Name}}Repo() {{.Name}}Repo {
    switch f.storeType {
    case core.STORETYPE_POSTGRES:
        table := {{.PGTablePackage}}.NewPg{{.Name}}Table(f.db)
        r := New{{.Name}}Repo(table)
        return r
    default:
    // not supported
    }
    return nil
}