package dbgen

type Options struct {
	Pluralise bool
}

type Option func(*Options)

func Pluralise(val bool) Option {
	return func(args *Options) {
		args.Pluralise = val
	}
}
