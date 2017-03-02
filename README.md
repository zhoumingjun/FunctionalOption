# gog
gog is short for "golang generator"
It is a tool for generating code skeleton with specific pattern

# Installation

go get -u github.com/zhoumingjun/gog

# Quick Start
Sample user-defined struct:

```
//go:generate gog fo -t Example
type Example struct {
	Name    string `option:"Name"`
	Age     int    `option`
	Address string `option:"Addr"`
}
```

gog will generate functional options for the field marked with "option" 
```
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

```

Then create the instance in the style of functional options
```
	xx := New(Name("haha"), Age(11), Addr("some address"))
```