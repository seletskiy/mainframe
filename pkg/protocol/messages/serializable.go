package messages

type Arg struct {
	Name  string
	Value interface{}
}

type Serializable interface {
	Tagged

	Serialize() []Arg
}
