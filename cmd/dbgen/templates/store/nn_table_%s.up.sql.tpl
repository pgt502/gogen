CREATE TABLE IF NOT EXISTS {{.Schema}}.{{.Table}}(
    {{range $i, $col := .Fields}} "{{$col.Column}}" {{$col.DBType}} NOT NULL,
    {{end}}
    PRIMARY KEY ({{range $i, $f := .PKFields}}{{if $i}},{{end}}"{{$f.Column}}"{{end}})
);