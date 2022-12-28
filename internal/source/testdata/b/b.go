package b

type BRunner interface {
	IsBRunner()
}

type B struct {
	Name string
}

type BB struct {
	Name string
}

func (b *B) IsBRunner()  {}
func (b *BB) IsBRunner() {}
