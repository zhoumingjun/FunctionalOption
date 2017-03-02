package example

//go:generate gog fo -t Example
type Example struct {
	Name    string `option:"Name"`
	Age     int    `option`
	Address string `option:"Addr"`
}
