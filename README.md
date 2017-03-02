# gog
gog is short for "golang generator"
It is a tool for generating code skeleton with specific pattern

# Installation

go get -u github.com/zhoumingjun/gog

# Quick Start
Sample user-defined struct:

```
//go:generate gog fo -t example
type example struct {
	Name    string `option:"Name"`
	Option1 string `option`
	Option2 xx
}
```

gog will generate functional options for the field marked with "option" 
```
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
```

Then create the instance in the style of functional options
```
xx := New(Name("haha"), Option1("ok"))
```