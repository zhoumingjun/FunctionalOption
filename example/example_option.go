package example

// Option is used to set options for the Example.
type Option interface {
	apply(*Example)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*Example)

func (f optionFunc) apply(p *Example) {
	f(p)
}

// the New funciton with options
func New(options ...Option) *Example {

	p := &Example{}
	return p.WithOptions(options...)
}

func (p *Example) WithOptions(opts ...Option) *Example {

	for _, opt := range opts {
		opt.apply(p)
	}
	return p
}

// options

func Name(v string) Option {
	return optionFunc(func(p *Example) {
		p.Name = v
	})
}

func Age(v int) Option {
	return optionFunc(func(p *Example) {
		p.Age = v
	})
}

func Addr(v string) Option {
	return optionFunc(func(p *Example) {
		p.Address = v
	})
}
