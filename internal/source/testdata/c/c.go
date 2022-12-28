package c

type CRunner interface {
	IsCRunner()
}

type C struct {
	Name string
}

type CC struct {
	Name string
}

func (c C) IsCRunner()  {}
func (c CC) IsCRunner() {}
