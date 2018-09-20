package messages

type Subscribe struct {
	Resize   bool
	Keyboard bool
	Input    bool
}

func (message *Subscribe) Tag() string {
	return "subscribe"
}
