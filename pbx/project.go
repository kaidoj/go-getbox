package pbx

type Project struct {
	Id     string
	Name   string
	Render Render
}

type Render struct {
	Status int
	URL    string
}
