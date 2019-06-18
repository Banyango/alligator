package endpoint

type Endpoint struct {
	Name string
	Host string
}

var GOOGLE = &Endpoint{
	Name:"Google",
	Host:"google.com",
}