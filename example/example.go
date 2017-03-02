package example

//go:generate gog fo -t example
type example struct {
	Name    string `option:"Name"`
	Option1 string `option`
	Option2 xx
}

type xx struct {
	Name    string `option:"Name`
	Option1 string `option:"Option1`
}
