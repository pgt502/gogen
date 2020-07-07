package dbgen

type Options struct {
	Pluralise   bool
	PackageName *string
	File        *string
}

type Option func(*Options)

func Pluralise(val bool) Option {
	return func(args *Options) {
		args.Pluralise = val
	}
}

func Package(val *string) Option {
	return func(args *Options) {
		args.PackageName = val
	}
}

func File(val *string) Option {
	return func(args *Options) {
		args.File = val
	}
}
