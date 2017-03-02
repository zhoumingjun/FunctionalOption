package fo

const tmplOptions = `
package {{.PackageName}}
// Option is used to set options for the logger.
type Option interface {
    apply(*{{.TypeName}})
}

// optionFunc wraps a func so it satisfies the Option interface.
type optionFunc func(*{{.TypeName}})

func (f optionFunc) apply(p *{{.TypeName}}) {
    f(p)
}

func New(options ...Option) *{{.TypeName}} {
 
    p := &{{.TypeName}}{

    }
    return p.WithOptions(options...)
}

func (p *{{.TypeName}}) WithOptions(opts ...Option) *{{.TypeName}} {

    for _, opt := range opts {
        opt.apply(p)
    }
    return p
}


{{range $index, $option := .Options }}
func {{$option.OptionName}}(v {{$option.Type}}) Option {
    return optionFunc(func(p *{{$.TypeName}}) {
        p.{{$option.FieldName}} = v
    })
}

{{end}}
`
