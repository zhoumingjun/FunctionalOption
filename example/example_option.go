package example

// Option is used to set options for the logger.
type Option interface {
	apply(*example)
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*example)

func (f optionFunc) apply(p *example) {
	f(p)
}

func New(options ...Option) *example {

	p := &example{}
	return p.WithOptions(options...)
}

func (p *example) WithOptions(opts ...Option) *example {

	for _, opt := range opts {
		opt.apply(p)
	}
	return p
}

func Name(v string) Option {
	return optionFunc(func(p *example) {
		p.Name = v
	})
}

func Option1(v string) Option {
	return optionFunc(func(p *example) {
		p.Option1 = v
	})
}
